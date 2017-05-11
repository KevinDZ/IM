//chatroom server
package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

const (
	LOG_DIRECTORY = "./test.log"
)

var onLineConns = make(map[string]net.Conn)
var messageQueue = make(chan string, 1000)
var QuitChan = make(chan bool)
var logFile *os.File
var logger *log.Logger

func CheckError(err error) {
	if err != nil {
		logger.Println(err)
	}
}

func ProcessInfo(conn net.Conn) {
	buf := make([]byte, 1024)
	//稍微好点的错误处理
	defer func(conn net.Conn) {
		addr := fmt.Sprintf("%s", conn.RemoteAddr())
		delete(onLineConns, addr)
		conn.Close()

		for i := range onLineConns {
			fmt.Println("now online conns:" + i)
		}
	}(conn)

	for {
		numOfBytes, err := conn.Read(buf)
		if err != nil {
			break
		}
		if numOfBytes != 0 {
			message := string(buf[0:numOfBytes])
			messageQueue <- message
		}
	}
}

func ConsumeMessage() {
	for {
		select {
		case message := <-messageQueue:
			doProcessMessage(message)
		case <-QuitChan:
			break
		}
	}
}

func doProcessMessage(message string) {
	contents := strings.Split(message, "#")
	if len(contents) > 1 {
		addr := contents[0]
		sendMessage := strings.Join(contents[1:], "#")
		addr = strings.Trim(addr, " ")
		if conn, ok := onLineConns[addr]; ok {
			_, err := conn.Write([]byte(sendMessage))
			if err != nil {
				fmt.Println("online conns send failure!")
			}
		}
	} else {
		contents := strings.Split(message, "*")
		if strings.ToUpper(contents[1]) == "LIST" {
			var ips string = ""
			for i := range onLineConns {
				ips = ips + "|" + i
			}
			if conn, ok := onLineConns[contents[0]]; ok {
				_, err := conn.Write([]byte(ips))
				if err != nil {
					fmt.Println("online conns send failure!")
				}
			}
		}
	}
}

func main() {
	logFile, err := os.OpenFile(LOG_DIRECTORY, os.O_RDWR|os.O_CREATE, 0)
	if err != nil {
		fmt.Println("log file create failure!")
		os.Exit(-1)
	}
	defer logFile.Close()
	logger = log.New(logFile, "\r\n", log.Ldate|log.Ltime|log.Llongfile)

	listen_socket, err := net.Listen("tcp", "127.0.0.1:8080")
	CheckError(err)
	defer listen_socket.Close()
	fmt.Println("Server is waiting....")

	logger.Println("i am writing log!")
	go ConsumeMessage()

	for {
		conn, err := listen_socket.Accept()
		CheckError(err)

		//conn储存到onLineConns映射表
		addr := fmt.Sprintf("%s", conn.RemoteAddr())
		onLineConns[addr] = conn

		for i := range onLineConns {
			fmt.Println(i)
		}
		go ProcessInfo(conn)
	}
}
