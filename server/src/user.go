package main

import (
	"fmt"
	"net"
	"strings"
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

func (this *User) SendMessage(msg string) {
	this.conn.Write([]byte(msg))
}

func (this *User) DoMessage(msg string) {

	if msg == "who" {
		this.server.mapLock.Lock()
		for _, usr := range this.server.OnlineMap {
			onlinemsg := "[" + usr.Addr + "]" + usr.Name + " : 在线 ... \n"
			this.SendMessage(onlinemsg)
		}
		this.server.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		newname := strings.Split(msg, "|")[1]

		//check whether name is exist
		_, ok := this.server.OnlineMap[newname]

		if ok {
			this.SendMessage("用户名已经存在！！！")
		} else {
			this.server.mapLock.Lock()
			delete(this.server.OnlineMap, this.Name)
			this.server.OnlineMap[newname] = this
			this.server.mapLock.Unlock()
			this.Name = newname

			this.SendMessage("已经更新用户名为 ：" + newname)
		}

	} else if len(msg) > 4 && msg[:3] == "to|" {
		remotename := strings.Split(msg, "|")[1]
		if remotename == "" {
			fmt.Println("消息格式错误！！！！！")
		}

		remoteuser, ok := this.server.OnlineMap[remotename]
		if !ok {
			fmt.Println("该用户不存在")
		}

		content := strings.Split(msg, "|")[2]
		if content == "" {
			fmt.Println("消息为空，请重新发送")
		}

		remoteuser.SendMessage(this.Name + "对您说:" + content)

	} else {
		this.server.BroadCast(this, msg)
	}

}
