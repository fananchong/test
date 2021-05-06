package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	_ "net/http/pprof"
)

var port = flag.String("port", "54321", "port")
var msgsize = flag.Int("message_size", 100, "message size")

type session struct {
	conn net.Conn
	c    chan []byte
}

func (s *session) readloop() {
	for {
		msg := make([]byte, *msgsize)
		io.ReadFull(s.conn, msg)
		s.c <- msg
	}
}

func (s *session) sendloop() {
	for {
		msg := <-s.c
		s.conn.Write(msg)
	}
}

func main() {
	flag.Parse()

	listen, err := net.Listen("tcp", fmt.Sprintf(":%s", *port))
	if err != nil {
		fmt.Printf("listen failed, err:%v\n", err)
		return
	}

	go func() {
		http.ListenAndServe(":54320", nil)
	}()

	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Printf("accept failed, err:%v\n", err)
			return
		}
		sess := &session{
			conn: conn,
			c:    make(chan []byte),
		}
		go sess.readloop()
		go sess.sendloop()
	}
}
