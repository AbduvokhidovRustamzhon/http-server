package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"strings"
)


var ContentType = make(map[string]string)


func Init() {
	ContentType["/"] = "Content-Type: file/html\r\n"
	ContentType["jpg"] = "Content-Type: image/jpg\r\n"
	ContentType["jpg?download"] = "Content-Disposition: attachment; filename = photo.jpg\r\n"
	ContentType["png"] = "Content-Type: image/png\r\n"
	ContentType["png?download"] = "Content-Disposition: attachment; filename = photo.png\r\n"
	ContentType["html"] = "Content-Type: text/html\r\n"
	ContentType["html?download"] = "Content-Disposition: attachment; filename = file.html\r\n"
	ContentType["txt"] = "Content-Type: text/html\r\n"
	ContentType["txt?download"] = "Content-Disposition: attachment; filename = text.txt\r\n"
	ContentType["pdf"] = "Content-Type: application/pdf\r\n"
	ContentType["jpg?download"] = "Content-Disposition: attachment; filename=error404.jpg\r\n"
}

func AnswerToHttp(http string, conn net.Conn) {
	var index string
	download := ""
	var file string
	if strings.HasSuffix(http, "?download") {
		http = http[1 : len(http)-9]
		file = http

		download = "?download"
	} else {
		file = http[1:]
	}
	if strings.HasSuffix(http, ".jpg") {
		index = "jpg"
	} else
	if strings.HasSuffix(http, ".png") {
		index = "png"
	} else
	if strings.HasSuffix(http, ".html") {
		index = "html"
	} else
	if strings.HasSuffix(http, ".txt") {
		index = "txt"
	} else
	if strings.HasSuffix(http, ".pdf") {
		index = "pdf"
	} else {
		index = ""
	}
	index += download
	if file == "" {
		file = "C:/Users/Admin/Desktop/http-server/server/files/main.html"
	}

	all, err := ioutil.ReadFile(file)
	if err != nil {
		log.Printf("can't read file %v", err)
		file = "C:/Users/Admin/Desktop/http-server/server/files/error404.jpg"
		all, err = ioutil.ReadFile(file)
		if err != nil {
			log.Printf("can't read file %v", err)
			return
		}
	}
	writer := bufio.NewWriter(conn)
	_, err = writer.WriteString("HTTP/1.1 200 OK\r\n")
	if err != nil {
		log.Printf("can't send string %v", err)
		return
	}
	_, err = writer.WriteString(fmt.Sprintf("Content-Length: %d\r\n", len(all)))
	if err != nil {
		log.Printf("can't send string %v", err)
		return
	}
	_, err = writer.WriteString(ContentType[index])
	if err != nil {
		log.Printf("can't send string %v", err)
		return
	}
	_, err = writer.WriteString("Connection: Close\r\n")
	if err != nil {
		log.Printf("can't send string %v", err)
		return
	}
	_, err = writer.WriteString("\r\n")
	if err != nil {
		log.Printf("can't send string %v", err)
		return
	}
	_, err = writer.Write(all)
	if err != nil {
		log.Printf("can't send string %v", err)
		return
	}
	err = writer.Flush()
	if err != nil {
		log.Printf("can't send string %v", err)
		return
	}
}
