package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int

	Name string
	conn net.Conn

	//choose chating mode
	flag int
}

func (this *Client) menu() bool {
	var flag int

	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更换用户名")
	fmt.Println("0.退出")

	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3 {
		this.flag = flag
		return true
	} else {
		fmt.Println("请输入合法数字")
		return false
	}
}

func NewClient(sip string, sp int) *Client {
	c := &Client{
		ServerIp:   sip,
		ServerPort: sp,
		flag:       999,
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", c.ServerIp, c.ServerPort))

	if err != nil {
		fmt.Println("dial error", err)
		return nil
	}

	c.conn = conn

	return c

}

func (this *Client) UpdateName() bool {
	fmt.Println("输入用户名")
	fmt.Scanln(&this.Name)

	sendmsg := "rename|" + this.Name + "\n"
	_, err := this.conn.Write([]byte(sendmsg))

	if err != nil {
		fmt.Println("err :", err)
		return false
	}
	return true

}

func (this *Client) DealResponse() {
	//这句话等价于下面的注释，都是从conn套接字中读取缓冲区中数据
	//stdout 永久阻塞监听
	io.Copy(os.Stdout, this.conn)
	/*for{
		buf := make([]byte ,4096)
		n, _ := this.conn.Read(buf)
		msg := buf[:n-1]
	}*/

}

func (this *Client) SelectUsers() {
	msg := "who\n"
	_, err := this.conn.Write([]byte(msg))

	if err != nil {
		fmt.Println("err is ", err)
		return
	}
}

func (this *Client) PrivateChat() {

	var (
		chatmsg    string
		remotename string
	)
	this.SelectUsers()
	fmt.Println("请输入要发送给的对象, exit 退出")
	fmt.Scanln(&remotename)

	//if input exit , public system exit
	for remotename != "exit" {
		fmt.Println("请输入要发送的内容")
		fmt.Scanln(&chatmsg)

		for chatmsg != "exit" {
			if len(chatmsg) != 0 {
				sendmsg := "to|" + remotename + "|" + chatmsg + "\n"
				_, err := this.conn.Write([]byte(sendmsg))
				if err != nil {
					fmt.Println("err is", err)
					break
				}
			}
			chatmsg = ""
			fmt.Println("请输入要发送的内容 exit 退出")
			fmt.Scanln(&chatmsg)
		}
	}
	remotename = ""
	this.SelectUsers()
	fmt.Println("请输入要发送给的对象, exit 退出")
	fmt.Scanln(&remotename)
}

func (this *Client) PublicChat() {

	var chatmsg string
	fmt.Println("请输入要发送的内容")
	fmt.Scanln(&chatmsg)
	//if input exit , public system exit
	for chatmsg != "exit" {
		if len(chatmsg) != 0 {
			sendmsg := chatmsg + "\n"
			_, err := this.conn.Write([]byte(sendmsg))
			if err != nil {
				fmt.Println("err is", err)
				break
			}
		}
		chatmsg = ""
		fmt.Println("请输入要发送的内容")
		fmt.Scanln(&chatmsg)
	}
}

func (this *Client) Run() {
	for this.flag != 0 {
		for this.menu() != true {
		}

		switch this.flag {
		case 1:
			fmt.Println("公聊")
			this.PublicChat()
			break

		case 2:
			fmt.Println("私聊")
			this.PrivateChat()
			break

		case 3:
			fmt.Println("修改用户名")
			this.UpdateName()
			break

		case 0:
			fmt.Println("退出")
			break
		}

	}
}

//命令行解析 输入参数
var serverip string
var serverport int

func init() {
	flag.StringVar(&serverip, "ip", "127.0.0.1", "设置服务器地址（默认为127.0.0.1）")
	flag.IntVar(&serverport, "port", 1111, "设置端口地址（默认为12345）")
}

func main() {
	flag.Parse()

	var client *Client = NewClient(serverip, serverport)
	if client == nil {
		fmt.Println("some thing is wrong")
		return
	}

	fmt.Println("connect succeed")

	// message recieve thread
	go client.DealResponse()
	//run client
	client.Run()
}
