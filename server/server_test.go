package main

import (
	"bufio"
	"io/ioutil"
	"net"
	"strings"
	"testing"
	"time"
)
func Test_Server(t *testing.T) {
	go func() {
		err := startServer("localhost:2001")
		if err != nil {
			t.Fatalf("can't start server %v", err)
		}
	}()
	time.Sleep(time.Second * 2)
	conn, err := net.Dial("tcp", "localhost:2001")
	if err != nil {
		t.Fatalf("can't connect to srver %v", err)
	}
	defer conn.Close()
	writer := bufio.NewWriter(conn)
	_, err = writer.WriteString("GET / HTTP/1.1\r\nHost: localhost\r\n\r\n")
	if err != nil {
		t.Fatalf("can't write to server %v", err)
	}
	err = writer.Flush()
	if err != nil {
		t.Fatalf("can't write to server %v", err)
	}
	bytes, err := ioutil.ReadAll(conn)
	if err != nil {
		t.Fatalf("can't read from server %v", err)
	}
	response := string(bytes)
	if !strings.Contains(response, "200 OK") {
		t.Fatalf("non-success response: %s", response)
	}

}