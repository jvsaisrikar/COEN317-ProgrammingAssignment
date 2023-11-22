package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

var Directory = "./"
var Port int
var pacificTimeZone *time.Location

func init() {
	var err error
	pacificTimeZone, err = time.LoadLocation("America/Los_Angeles")
	if err != nil {
		log.Fatal("Failed to load Pacific timezone:", err)
	}
}

func processClientRequest(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return
	}
	request := string(buffer[:n]) // Only capture the received bytes
	requestHeaders := extractMethodLine(request)
	handleRequest(request, requestHeaders, conn)
}

func extractMethodLine(request string) string {
	lines := strings.Split(request, "\n")
	if len(lines) == 0 {
		return ""
	}
	methodLine := lines[0]
	return methodLine
}

func logRequest(statusCode int, requestURI string) {
	log.Printf("[%s] %d: %s", time.Now().In(pacificTimeZone).Format(time.RFC3339), statusCode, requestURI)
}

func handleRequest(fullRequest string, requestHeaders string, conn net.Conn) {
	parts := strings.SplitN(requestHeaders, " ", 3)
	if len(parts) < 3 {
		logRequest(400, fullRequest)
		conn.Write([]byte("HTTP/1.1 400 Bad Request\r\nContent-Type: text/plain\r\n\r\nBad Request"))
		return
	}

	requestURI := parts[1]

	if requestURI == "/" {
		requestURI = "/index.html"
	}

	if strings.HasPrefix(requestURI, "/private") {
		logRequest(403, requestURI)
		conn.Write([]byte("HTTP/1.1 403 Forbidden\r\nContent-Type: text/plain\r\n\r\nAccess Denied"))
		return
	}

	filepath := Directory + requestURI
	content, err := ioutil.ReadFile(filepath)

	if err != nil {
		if os.IsNotExist(err) {
			logRequest(404, requestURI)
			conn.Write([]byte("HTTP/1.1 404 NOT FOUND\r\nContent-Type: text/plain\r\n\r\nFile Not Found"))
		} else {
			logRequest(500, requestURI)
			conn.Write([]byte("HTTP/1.1 500 Internal Server Error\r\nContent-Type: text/plain\r\n\r\nInternal Server Error"))
		}
		return
	}

	now := time.Now().UTC()
	date := now.Format(http.TimeFormat)

	var contentType string
	switch {
	case strings.HasSuffix(requestURI, ".html"):
		contentType = "text/html"
	case strings.HasSuffix(requestURI, "error.html"):
		contentType = "text/html"
	case strings.HasSuffix(requestURI, ".txt"):
		contentType = "text/plain"
	case strings.HasSuffix(requestURI, ".jpg"):
		contentType = "image/jpeg"
	case strings.HasSuffix(requestURI, ".gif"):
		contentType = "image/gif"
	default:
		contentType = "application/octet-stream"
	}

	headers := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: %s\r\nContent-Length: %d\r\nDate: %s\r\n\r\n",
		contentType, len(content), date)

	logRequest(200, requestURI)
	conn.Write([]byte(headers))
	conn.Write(content)
}

func main() {
	flag.StringVar(&Directory, "document_root", "./", "Root directory of server")
	flag.IntVar(&Port, "port", 8888, "Port to run server on")
	flag.Parse()

	ln, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", Port))
	if err != nil {
		log.Println("Error starting server:", err)
		return
	}
	defer ln.Close()

	fmt.Printf("Server is running on: %d\nServer Directory: %s\n", Port, Directory)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}
		go processClientRequest(conn)
	}
}
