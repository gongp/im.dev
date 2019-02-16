
MeQ [mi:kju]
------------
A modern  platform for Message Push、Instant Messaging etc, based on [MQTT protocol](https://github.com/mafanr/meq/tree/master/proto/mqtt). 

MeQ is written in pure go. Our goal is to be the best messaging platform in the world.

- Homepage: http://meq.io
- <a href="ReadmeCn.md">中文Readme</a>
<p align="left">
    <a href="http://meq.io">
</p>

### Project status
Still under  developing, 2019-01-03


Example of Chatting room
------------
### Install the foundationdb client(required!)
Mac OS: https://www.foundationdb.org/downloads/5.1.7/macOS/installers/FoundationDB-5.1.7.pkg

CentOS7: https://www.foundationdb.org/downloads/5.1.7/rhel7/installers/foundationdb-clients-5.1.7-1.el7.x86_64.rpm

CentOS6: https://www.foundationdb.org/downloads/5.1.7/rhel6/installers/foundationdb-clients-5.1.7-1.el6.x86_64.rpm

Other Os: https://www.foundationdb.org/download/ (Please choose version 5.1.7)

### Install foundationdb server(optional,only needed when you use fdb store mode)
CentOS7: https://www.foundationdb.org/downloads/5.1.7/rhel7/installers/foundationdb-server-5.1.7-1.el7.x86_64.rpm

### Start MeQ
```bash
> go get github.com/mafanr/meq
> cd $GOPATH/src/github.com/mafanr/meq/broker
> go run main.go
```
### Start example
```bash
> cd ../demos/chatting
> npm install
> npm run dev
```
### Chatting now!
Open your browser and accsess http://localhost:8080 in two seperate pages,in one page input username A,in other page input username B, then you can start chatting!

### More(persistent storage)
Broker uses memory storage by default, if you want to use persistent storage, please install [foundationDB server V5.1.7](https://www.foundationdb.org/download/) , and then edit the broker.yaml,change store.engine from memory to fdb.

Features
------------
### Extremely performant
- Written in go, as fast as c/c++ network programing
- Zero allocation and copy
- Hot path and algorithm optimized

### Easy
- Easy to study and use
- Detailed docs and examples

### Robust
- Robust is one of our most important goal throughout the developing
- Default message persistent,also you can implement the persistent inteface in your way

### Various scenario supported
- Message push
- Group chatting
- IM
- IoT messaging
- Real time Web interaction,like dashboard、online collaboration etc
### Group chat
- Join and leave the group
- Each user has a separate unread message number,e.g. you have 97 unread messages in ethereum-welcom group
- You can retrive your message even after sended
- History messages playback
- Query all the users or online users in the group
- Every message will store only once, all the users share the message
### Advanced feature
- Topic wild match, e.g. you can  push to /china/+/city1, then **all** the **china** province which has city named **city1**, will receive the message.
- Mqtt and websocket,tcp etc
- Low latency, message push deliver to user as soon as possible
- Fertile administrator features and ui
- Monitor and message trace

Why MeQ? 
------------
There were already many messaging platform out there, but none of them can satisfy all the messaging scenario, e.g., if you want build a online chatting platform, you need to build it from scratch, nearly none of the existing platform can satisfy your need. So here meq comes, cover your messaging scenario in the most easy way.

Contribute
------------
Want to be part of the meq project? Great! All are welcome! We will get there quicker together :) Whether you find a bug, have a great feature request or you fancy owning a task from the road map, you can leave a issue.


Contributors(Sort by contribution)
------------
- <a href="https://github.com/sunface" target="_blank">Sunface</a> 
- <a href="https://github.com/shaocongcong" target="_blank">Cong</a>




