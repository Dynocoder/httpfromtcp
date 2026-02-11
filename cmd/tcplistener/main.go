package main

import (
	"fmt"
	"httpfromtcp/internal/request"
	"log"
	"net"
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

		request, err := request.RequestFromReader(conn)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Request Line: ")
		fmt.Printf("- Method: %s\n", request.RequestLine.Method)
		fmt.Printf("- Target: %s\n", request.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", request.RequestLine.HttpVersion)

		fmt.Println("Headers: ")
		for name, value := range request.Headers {
			fmt.Printf("- %s: %s\n", name, value)
		}

		fmt.Println("Body: ")
		fmt.Println(string(request.Body))

	}
}
