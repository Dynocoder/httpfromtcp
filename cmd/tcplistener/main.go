package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func main() {

	listener, err := net.Listen("tcp", ":42069")

	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Connection Established")

		c := getLinesChannel(conn)

		for line := range c {
			fmt.Println(line)
		}

	}
}

func getLinesChannel(conn net.Conn) <-chan string {

	c := make(chan string)

	go readLine(conn, c)

	return c

}

func readLine(conn net.Conn, c chan string) {
	buf := make([]byte, 8)
	line := ""
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Connection Closed")
				close(c)
				conn.Close()
				return
			} else {
				fmt.Println("Error", err)
				close(c)
				return
			}
		}
		data := string(buf[:n])

		parts := strings.Split(data, "\n")

		if len(parts) == 1 {
			line += parts[0]
		} else {
			line += parts[0]
			c <- line
			line = parts[1]
		}
	}
}
