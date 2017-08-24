package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"
)

var quitSemaphore chan bool

func main() {
	var tcpAddr *net.TCPAddr
	var totalTime int64

	tcpAddr, _ = net.ResolveTCPAddr("tcp", "10.0.0.1:6060")

	conn, _ := net.DialTCP("tcp", nil, tcpAddr)
	defer conn.Close()
	fmt.Println("connected!")

	go onMessageRecived(conn, &totalTime)
	var count int
	for {
		totalTime = 0
		fmt.Println("测试次数:")
		fmt.Scanln(&count)
		if count == -1 {
			break
		}
		for i := 0; i < count; i++ {
			b := []byte("download#")
			conn.Write(b)
			// <-quitSemaphore
			time.Sleep(500 * time.Millisecond)
		}
		<-quitSemaphore
	}
	// fmt.Println("测试次数:")
	// fmt.Scanln(&count)
	// for i := 0; i < count; i++ {
	// 	b := []byte("download#")
	// 	conn.Write(b)
	// 	// <-quitSemaphore
	// 	time.Sleep(500 * time.Millisecond)
	// }
	// <-quitSemaphore
}

func onMessageRecived(conn *net.TCPConn, totalTime *int64) {
	reader := bufio.NewReader(conn)
	for {
		startTime := strconv.FormatInt(time.Now().UnixNano(), 10)
		_, err := reader.ReadString('#')
		if err != nil {
			quitSemaphore <- true
			break
		}
		var mutex sync.Mutex
		mutex.Lock()
		defer mutex.Unlock()
		d_time := subTime(startTime)
		*totalTime += d_time
		fmt.Printf("下载文件耗时%d纳秒\n总耗时%d纳秒\n\n", d_time, *totalTime)
	}
}
func subTime(startTime string) int64 {
	currentTime := time.Now().UnixNano() //获取当前时间，类型是Go的时间类型Time
	start, _ := strconv.ParseInt(startTime, 10, 64)
	return currentTime - start
}
