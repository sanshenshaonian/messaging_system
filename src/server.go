package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

func NewServer(ip string, port int) *Server {
	nserver := &Server{
		Ip:   ip,
		Port: port,
	}
	return nserver
}

func (this *Server) Hnadler(connfd net.Conn) {
	fmt.Println("和客户端链接已经建立， 描述符号为connfd")
}

func (this *Server) Start() {

	//create listener thread
	listenfd, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("listen error is : ", err)
	}
	defer listenfd.Close()

	for {
		connfd, err := listenfd.Accept()
		if err != nil {
			fmt.Println("accept is err :", err)
		}

		go this.Hnadler(connfd)
	}

}
