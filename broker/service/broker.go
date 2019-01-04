//  Copyright © 2018 Sunface <CTO@188.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package service

import (
	"fmt"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"net/http"
	_ "net/http/pprof"

	"github.com/bwmarrin/snowflake"
	"github.com/mafanr/g"
	"github.com/mafanr/meq/proto/websocket"
	"github.com/sunface/talent"
	"go.uber.org/zap"
)

type Broker struct {
	wg          *sync.WaitGroup
	running     bool
	runningTime time.Time
	listener    net.Listener

	clients map[uint64]*client

	store   Storage
	router  *Router
	cluster *cluster

	subtrie   *SubTrie
	subSynced bool

	idgen  *snowflake.Node
	topics *Topics

	conf *Config
	sync.RWMutex
}

func NewBroker(path string) *Broker {
	b := &Broker{
		wg:      &sync.WaitGroup{},
		clients: make(map[uint64]*client),
		subtrie: NewSubTrie(),
	}
	// init base config
	b.conf = initConfig(path)
	g.InitLogger(b.conf.Common.LogLevel)
	g.L.Info("base configuration loaded")

	return b
}

func (b *Broker) Start() {
	b.running = true
	b.runningTime = time.Now()

	// tcp listener
	b.startTcp()
	// websocket listenser
	b.startWS()

	// admin
	admin := &Admin{}
	admin.Init(b)

	// init store
	switch b.conf.Store.Engine {
	case "memory":
		b.store = &MemStore{
			bk: b,
		}
	case "fdb":
		b.store = &FdbStore{
			bk: b,
		}
	}

	b.store.Init()

	b.topics = &Topics{}
	b.topics.Init(b)

	// init cluster
	b.cluster = &cluster{
		bk: b,
	}
	go b.cluster.Init()

	// init Router
	b.router = &Router{
		bk: b,
	}
	b.router.Init()

	// init messsage id generator
	StartIDGenerator(b)

	go func() {
		log.Println(http.ListenAndServe(":6070", nil))
	}()
}
func (b *Broker) Shutdown() {
	b.running = false
	b.listener.Close()

	for _, c := range b.clients {
		c.conn.Close()
	}
	b.cluster.Close()
	b.store.Close()
	b.router.Close()

	g.L.Sync()
	b.wg.Wait()
}

var uid uint64

func (b *Broker) Accept() {
	tmpDelay := ACCEPT_MIN_SLEEP

	for b.running {
		conn, err := b.listener.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				g.L.Error("Temporary Client Accept Error ", zap.Error(err))
				time.Sleep(tmpDelay)
				tmpDelay *= 2
				if tmpDelay > ACCEPT_MAX_SLEEP {
					tmpDelay = ACCEPT_MAX_SLEEP
				}
			} else if b.running {
				g.L.Error("Client Accept Error", zap.Error(err))
			}
			continue
		}
		tmpDelay = ACCEPT_MIN_SLEEP
		atomic.AddUint64(&uid, 1)
		go b.process(conn, uid, false)
	}
}

func (b *Broker) process(conn net.Conn, id uint64, isWs bool) {
	defer func() {
		b.Lock()
		delete(b.clients, id)
		b.Unlock()
		conn.Close()
		g.L.Info("client closed", zap.Uint64("conn_id", id))
	}()

	cli := initClient(id, conn, b)

	b.Lock()
	b.clients[id] = cli
	b.Unlock()

	err := cli.waitForConnect()
	if err != nil {
		fmt.Println("cant receive connect packet from client", err, zap.Uint64("cid", id))
		return
	}

	g.L.Info("new user online", zap.Uint64("conn_id", id), zap.String("username", string(cli.username)), zap.String("ip", conn.RemoteAddr().String()))

	go cli.writeLoop()
	err = cli.readLoop(isWs)
	if err != nil {
		if !talent.IsEOF(err) {
			g.L.Info("client read loop error", zap.Uint64("cid", cli.cid), zap.Error(err))
		}
	}
}

func (b *Broker) startTcp() {
	addr := net.JoinHostPort(b.conf.Broker.Host, b.conf.Broker.TcpPort)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		g.L.Fatal("Fatal error when listening tcp", zap.Error(err), zap.String("addr", addr))
	}
	b.listener = l

	g.L.Info("tcp listening at :", zap.String("addr", addr))
	go b.Accept()
}

func (b *Broker) startWS() {
	addr := net.JoinHostPort(b.conf.Broker.Host, b.conf.Broker.WsPort)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		g.L.Fatal("Fatal error when listening websocket", zap.Error(err), zap.String("addr", addr))
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if conn, ok := websocket.TryUpgrade(w, r); ok {
			atomic.AddUint64(&uid, 1)
			go b.process(conn, uid, true)
		}
	})

	go http.Serve(lis, mux)
	g.L.Info("websocket listening at :", zap.String("addr", addr))
}
