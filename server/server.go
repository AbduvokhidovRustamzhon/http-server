package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
)


const wayInFiles = "server/files/"
const textOrHtml = "server/files/"
const imageOrPng  = "server/files/"
const imageOrImg  = "server/files/"
const pdf = "server/files/"

func main() {
host := "0.0.0.0"
port, ok := os.LookupEnv("PORT")
if !ok {
port = "9998"
}
err := start(fmt.Sprintf("%s:%s", host, port))
if err != nil {
log.Fatal(err)
}
}

func start(addr string) (err error) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("can't listen %s: %w", addr, err)
	}
	defer func() {
		err := listener.Close()
		if err != nil {
			log.Fatalf("Can't close conn: %v", err)
		}
	}()
	for {
		conn, err := listener.Accept()
		log.Print("accept connection")
		if err != nil {
			log.Printf("can't accept: %v", err)
			continue
		}
		log.Print("handle connection")
		handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Fatalf("Can't close conn: %v", err)
		}
	}()
	reader := bufio.NewReader(conn)
	requestLine, err := reader.ReadString('\n')
	if err != nil {
		return
	}
	log.Print(requestLine)
	parts := strings.Split(strings.TrimSpace(requestLine), " ")
	if len(parts) != 3 {
		return
	}

	method, request, protocol := parts[0], parts[1], parts[2]
	contentTypeDownload := ""
	if method == "GET" {
		if protocol == "HTTP/1.1" {
			if strings.Contains(request, "?download") {
				contentTypeDownload = "application/octet-stream"
				request = strings.Replace(request, "?download", "", -1)
			}

			switch request {
			case "/":
				optimization(conn,"main.html", textOrHtml, contentTypeDownload, request)
			case "/index-file.html":
				optimization(conn, "file.html"/*"index-file.html"*/, textOrHtml, contentTypeDownload, request)
			case "/photo.png":
				optimization(conn,"photo.png", imageOrPng, contentTypeDownload, request)
			case "/photo.jpg":
				optimization(conn,"photo.jpg", imageOrImg, contentTypeDownload, request)
			case "/text.txt":
				optimization(conn,"text.txt", textOrHtml, contentTypeDownload, request)
			case "/pdf.pdf":
				optimization(conn,"pdf.pdf", pdf, contentTypeDownload, request)
			case "/html5.png":
				optimization(conn,"html5.png", imageOrPng, contentTypeDownload, request)
			default:
				return
				// optimization(conn,"error404.jpg", imageOrImg, contentTypeDownload, request)
			}
			return
		}
	}
}

func optimization(conn net.Conn, fileName, contentType, contentTypeDownload, request string) {
	fileName = wayInFiles + fileName
	if contentTypeDownload != "" {
		contentType = contentTypeDownload
	}
	writeHeader(conn, fileName, contentType, request)
}

func writeHeader(conn net.Conn, fileName, contentType, request string) {
	writer := bufio.NewWriter(conn)
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Printf("can't read file: %v", err)
	}
	_, err = writer.WriteString("HTTP/1.1 200 OK\r\n")
	if err != nil {
		log.Printf("can't write: %v", err)
	}
	_, err = writer.WriteString(fmt.Sprintf("Content-Length: %d\r\n", len(bytes)))
	if err != nil {
		log.Printf("can't write: %v", err)
	}
	_, err = writer.WriteString("Content-Type:" + " " + contentType + "\r\n")
	if err != nil {
		log.Printf("can't write: %v", err)
	}
	_, err = writer.WriteString("Connection: Close\r\n")
	if err != nil {
		log.Printf("can't write: %v", err)
	}
	_, err = writer.WriteString("\r\n")
	if err != nil {
		log.Printf("can't write: %v", err)
	}
	_, err = writer.Write(bytes)
	if err != nil {
		log.Printf("can't write: %v", err)
	}
	err = writer.Flush()
	if err != nil {
		log.Printf("can't sent response: %v", err)
	}
	log.Printf("response on: %s", request)
}







