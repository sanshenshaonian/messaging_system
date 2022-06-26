package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int

	//load online user
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	Message chan string
}

func NewServer(ip string, port int) *Server {
	nserver := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return nserver
}

func (this *Server) ListenMessage() {
	for {
		msg := <-this.Message

		//map lock
		this.mapLock.Lock()
		for _, usr := range this.OnlineMap {
			usr.C <- msg
		}

		this.mapLock.Unlock()
	}
}

//send message to SERVER.CHANNEL
func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg

	this.Message <- sendMsg
}

func (this *Server) Hnadler(connfd net.Conn) {
	fmt.Println("和客户端链接已经建立， 描述符号为connfd")

	var user *User = NewUser(connfd, this)

	user.Online()

	go func() {
		buf := make([]byte, 4096)

		for {
			n, err := connfd.Read(buf)

			if n == 0 {
				user.Offline()
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("非法操作")
				return
			}

			msg := string(buf[:n-1])
			user.DoMessage(msg)
		}
	}()

	//block
	select {}
}

func (this *Server) Start() {

	//create listener thread
	listenfd, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("listen error is : ", err)
		return
	}
	defer listenfd.Close()

	//listen MESSAGE CHANNLE
	go this.ListenMessage()

	for {
		connfd, err := listenfd.Accept()
		if err != nil {
			fmt.Println("accept is err :", err)
			continue
		}

		go this.Hnadler(connfd)
	}

}
