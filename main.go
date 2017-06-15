package main

import (
	"./util"
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

var running = true

func sender(conn net.Conn, name string) {
	reader := bufio.NewReader(os.Stdin)
	for {
		input, _, _ := reader.ReadLine()
		if string(input) == "\\q" {
			running = false
			break
		}
		_, err := conn.Write(input)
		util.ChkErr(err, "sender write")
	}
}

func receiver(conn net.Conn, name string) {
	buf := make([]byte, 560)
	for running == true {
		n, err := conn.Read(buf)
		util.ChkErr(err, "Receiver read")
		fmt.Println(string(buf[:n]))
		buf = make([]byte, 560)
	}
}

func Close(c io.Closer) {
	err := c.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
	}
}

func main() {
	fmt.Print("Please input your name: ")
	reader := bufio.NewReader(os.Stdin)
	name, _, err := reader.ReadLine()

	host := "127.0.0.1:7777"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", host)
	util.ChkErr(err, "tcpAddr")

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	util.ChkErr(err, "DialTCP")
	defer Close(conn)

	_, err = conn.Write(name)
	util.ChkErr(err, "Write name")

	go receiver(conn, string(name))
	go sender(conn, string(name))

	for running {
		time.Sleep(1 * 1e9)
	}
}
