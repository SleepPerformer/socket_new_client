package main

import (
	"bufio"
	"fmt"
	"net"
)

var quitSemaphore chan bool

func main() {
	var tcpAddr *net.TCPAddr
	tcpAddr, _ = net.ResolveTCPAddr("tcp", "172.16.100.210:6060")

	conn, _ := net.DialTCP("tcp", nil, tcpAddr)
	defer conn.Close()
	fmt.Println("connected!")

	go onMessageRecived(conn)

	b := []byte("download#")
	conn.Write(b)

	<-quitSemaphore
}

func onMessageRecived(conn *net.TCPConn) {
	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('#')
		fmt.Println(msg)
		if err != nil {
			quitSemaphore <- true
			break
		}
		//		time.Sleep(time.Second)
		// b := []byte(msg)
		// conn.Write(b)
	}
}
