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
	"net"
	"net/http"

	"github.com/mafanr/g"
	"github.com/sunface/talent"
	"go.uber.org/zap"

	"github.com/apple/foundationdb/bindings/go/src/fdb"
	"github.com/labstack/echo"
)

type Admin struct {
	bk *Broker
}

func (ad *Admin) Init(bk *Broker) {
	ad.bk = bk
	go func() {
		e := echo.New()
		e.POST("/op/clearStore", ad.clearStore)
		e.GET("/status/users", ad.statusUsers)

		addr := net.JoinHostPort(ad.bk.conf.Broker.Host, ad.bk.conf.Admin.Port)
		e.Logger.Fatal(e.Start(addr))
		g.L.Info("http listening at :", zap.String("addr", addr))
	}()
}

func (ad *Admin) clearStore(c echo.Context) error {
	token := c.FormValue("token")
	if token != ad.bk.conf.Broker.Token {
		return c.String(http.StatusOK, "invalid admin token")
	}

	f, ok := ad.bk.store.(*FdbStore)
	if !ok {
		return c.String(http.StatusOK, "not fdb store engine,ignore")
	}

	i := getcounts % uint64(f.bk.conf.Store.FDB.Threads)
	d := f.dbs[i]

	_, err := d.db.Transact(func(tr fdb.Transaction) (ret interface{}, err error) {
		tr.ClearRange(d.msgsp)
		tr.ClearRange(d.normalCountSP)
		tr.ClearRange(d.chatroomSP)
		return
	})
	if err != nil {
		return c.String(http.StatusOK, fmt.Sprintf("error happens: %v", err))
	}

	return c.String(http.StatusOK, "clear store ok")
}

func (ad *Admin) statusUsers(c echo.Context) error {
	// token := c.FormValue("token")
	// if token != ad.bk.conf.Broker.Token {
	// 	return c.String(http.StatusOK, "invalid admin token")
	// }

	var ret = make(map[uint64]string)
	ad.bk.RLock()
	for _, cli := range ad.bk.clients {
		ret[cli.cid] = fmt.Sprintf("user %s login at %s", string(cli.username), talent.Time2String(cli.initTime))
	}
	ad.bk.RUnlock()
	return c.JSONPretty(200, ret, "  ")
}
