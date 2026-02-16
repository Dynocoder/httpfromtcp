package main

import (
	"fmt"
	"httpfromtcp/internal/headers"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"httpfromtcp/internal/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const port = 42069

func main() {
	server, err := server.Serve(
		port,
		func(w *response.Writer, req *request.Request) {
			headers := headers.NewHeaders()
			statusCode := response.StatusOK
			body := response200()

			switch req.RequestLine.RequestTarget {
			case "/yourproblem":
				statusCode = response.StatusBadRequest
				body = response400()
			case "/myproblem":
				statusCode = response.StatusInternalServerError
				body = response500()
			default:
				body = response200()
			}

			headers.Replace("Content-Length", fmt.Sprintf("%d", len(body)))
			headers.Replace("Content-Type", "text/html")

			w.ReturnResponse(statusCode, &headers, body)

		},
	)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()

	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")

}

// These should probably be a user provided file or smt
// but we can keep it out of scope for now and just serve
// this as an example
func response400() []byte {
	return []byte(`
		<html>
			<head>
				<title>400 Bad Request</title>
			</head>
			<body>
				<h1>Bad Request</h1>
				<p>Your request honestly kinda sucked.</p>
			</body>
		</html>
	`)
}

func response500() []byte {
	return []byte(`
	<html>
		<head>
			<title>500 Internal Server Error</title>
		</head>
		<body>
			<h1>Internal Server Error</h1>
			<p>Okay, you know what? This one is on me.</p>
		</body>
	</html>
	`)
}

func response200() []byte {
	return []byte(`
		<html>
			<head>
				<title>200 OK</title>
			</head>
			<body>
				<h1>Success!</h1>
				<p>Your request was an absolute banger.</p>
			</body>
		</html>
	`)
}
