package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {

	file, err := os.OpenFile("server-log.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Printf("can't open log file %e", err)
		return
	}
	defer func() {
		err := file.Close()
		if err != nil {
			fmt.Printf("can't close log file %e", err)
		}
	}()

	log.SetOutput(file)
	log.Print("start application\n")

	host := "0.0.0.0"
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "9990"
	}
	err = startServer(fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		log.Fatal(err)
	}




}

func startServer(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Printf("can't listen 0.0.0.0:2000 %v", err)
		return err
	}
	defer func() {
		err = listener.Close()
		if err != nil {
			log.Printf("Can't close Listener: %v", err)
		}
	}()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("can't connect client")
			continue
		}
		fmt.Println("some connection")
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	readString, _ := reader.ReadString('\n')
	split := strings.Split(strings.TrimSpace(readString), " ")
	if len(split) != 3 {
		log.Print("incorrect request")
		return
	}
	method, request, protocol := split[0], split[1], split[2]
	if method == "GET" && protocol == "HTTP/1.1" {
		AnswerToHttp(request, conn)
	}

}
