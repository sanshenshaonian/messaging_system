package main

import (
	"net"
)

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn

	//call server fucn
	server *Server
}

// create a user class
func NewUser(conn net.Conn, server *Server) *User {
	useraddr := conn.RemoteAddr().String()
	user := &User{
		Name:   useraddr,
		Addr:   useraddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
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

func (this *User) Online() {
	//save user in map
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Addr] = this
	this.server.mapLock.Unlock()

	//send msg to server.chan then server will broadcase to every user
	this.server.BroadCast(this, "上线了！！！！！！")
}

func (this *User) Offline() {
	//delete user in map
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Addr)
	this.server.mapLock.Unlock()

	//send msg to server.chan then server will broadcase to every user
	this.server.BroadCast(this, "下线了！！！！！！")
}

func (this *User) DoMessage(msg string) {
	this.server.BroadCast(this, msg)
}
