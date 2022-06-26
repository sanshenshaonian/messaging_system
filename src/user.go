package main

import (
	"net"
)

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn
}

// create a user class
func NewUser(conn net.Conn) *User {
	useraddr := conn.RemoteAddr().String()
	user := &User{
		Name: useraddr,
		Addr: useraddr,
		C:    make(chan string),
		conn: conn,
	}

	go user.ListenMessage()

	return user
}

// listen user channel
func (this *User) ListenMessage() {
	for {
		msg := <-this.C

		this.conn.Write([]byte(msg + "\n"))
	}
}
