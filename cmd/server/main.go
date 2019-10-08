package main

import (
	bytes2 "bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	listener, err := net.Listen("tcp", "0.0.0.0:9866")
	if err != nil {
		log.Fatal(err)
	}

	for {
		log.Println("wait connection")
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go handleConnection(conn)
	}

}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	for {
		err := conn.SetReadDeadline(time.Now().Add(time.Second * 10))
		if err != nil {
			log.Println(err)
			return
		}
		err = conn.SetWriteDeadline(time.Now().Add(time.Second * 10))
		if err != nil {
			log.Println(err)
			return
		}

		log.Println("connection")
		buf := make([]byte, binary.MaxVarintLen64)
		readBytes, err := conn.Read(buf)
		if err != nil {
			log.Println(err)
			return
		}
		bytesToRead, err := binary.ReadVarint(bytes2.NewReader(buf[:readBytes]))
		if err != nil {
			log.Fatal(err)
			return
		}

		log.Printf("Message size: %d bytes\n", bytesToRead)

		readBuffer := make([]byte, bytesToRead)

		_, err = conn.Read(readBuffer)
		if err != nil {
			log.Println(err)
			return
		}
		err = conn.SetReadDeadline(time.Now().Add(time.Second * 10))
		if err != nil {
			log.Println(err)
			return
		}
		/*
			bytes, err := ioutil.ReadAll(reader)
			if err != nil {
				log.Fatal(err)
			}*/

		s := string(readBuffer)
		header := ""
		switch s[0] {
		case '0':
			header = "int"
		case '1':
			header = "string"
		case '2':
			header = "float"
		default:
			header = "unknown"
		}

		log.Printf("Header: %c", s[0])
		line := fmt.Sprintf("Received message %s with %s header", s[1:], header)
		log.Println(line)
		bytesLength := len(line)
		buf1 := make([]byte, binary.MaxVarintLen64)
		binary.PutVarint(buf1, int64(bytesLength))
		buf1 = append(buf1, []byte(line)...)
		_, err = conn.Write(buf1)
		if err != nil {
			log.Println(err)
			return
		}
		err = conn.SetWriteDeadline(time.Now().Add(time.Second * 10))
		if err != nil {
			log.Println(err)
			return
		}

		log.Printf("sent back %s: %d bytes", line, bytesLength)
	}
}
