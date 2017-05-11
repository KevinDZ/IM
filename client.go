// chatrrom-client

package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
func MessageSend(conn net.Conn) {
	var input string
	for {
		reader := bufio.NewReader(os.Stdin)
		data, _, _ := reader.ReadLine()
		input = string(data)

		if strings.ToUpper(input) == "EXIT" {
			conn.Close()
			break
		}

		_, err := conn.Write([]byte(input))
		if err != nil {
			conn.Close()
			fmt.Println("client connect failure :" + err.Error())
		}
	}
}
func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	CheckError(err)
	defer conn.Close()

	go MessageSend(conn)

	buf := make([]byte, 1024)
	for {
		data, err := conn.Read(buf)
		if err != nil {
			fmt.Println("您已经退出，欢迎再次使用!")
			os.Exit(0)
		}
		fmt.Println("receiver server message content:" + string(buf[0:data]))
	}

}
