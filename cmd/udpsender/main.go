package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println("Address: ", addr.String())

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatal(err)
		return
	}

	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
			return
		}

		conn.Write([]byte(input))
	}

}
