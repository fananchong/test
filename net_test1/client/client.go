package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"sync/atomic"
	"time"
)

var ip = flag.String("ip", "127.0.0.1", "ip address")
var port = flag.String("port", "54321", "port")
var msgsize = flag.Int("message_size", 100, "message size")
var playerNum = flag.Int("playerNum", 500, "player num")

type session struct {
	playerIndex int
	conn        net.Conn
	c           chan []byte
	counter     int64
}

func (s *session) readloop() {
	for {
		msg := make([]byte, *msgsize)
		io.ReadFull(s.conn, msg)
		if msg[1] != 'a' || msg[*msgsize-1] != 'b' {
			panic("data error!")
		} else {
			atomic.AddInt64(&s.counter, 1)
		}
	}
}

func (s *session) sendloop() {
	for {
		msg := <-s.c
		s.conn.Write(msg)
	}
}

func (s *session) update() {
	tick := time.NewTicker(time.Second / 30)
	defer tick.Stop()
	for {
		<-tick.C

		v := atomic.LoadInt64(&s.counter)
		if v%(30*5) == 0 && v != 0 {
			fmt.Println(v)
		}

		msg := make([]byte, *msgsize)
		msg[1] = 'a'
		msg[*msgsize-1] = 'b'
		s.c <- msg
	}
}

func main() {
	flag.Parse()

	for i := 0; i < *playerNum; i++ {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", *ip, *port))
		if err != nil {
			fmt.Printf("conn server failed, err:%v\n", err)
			return
		}
		sess := &session{
			playerIndex: i,
			conn:        conn,
			c:           make(chan []byte),
		}
		go sess.readloop()
		go sess.sendloop()
		go sess.update()
	}

	select {}
}
