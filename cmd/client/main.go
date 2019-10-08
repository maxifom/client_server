package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	reader := bufio.NewReader(os.Stdin)
	conn, err := net.Dial("tcp", "0.0.0.0:9866")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	for {
		buf := make([]byte, binary.MaxVarintLen64)
		line, _, err := reader.ReadLine()
		if err != nil {
			log.Fatal(err)
		}
		if string(line) == "" {
			continue
		}
		bytesLength := len(line)

		binary.PutVarint(buf, int64(bytesLength))
		buf = append(buf, line...)
		_, err = conn.Write(buf)
		if err != nil {
			if err == io.ErrClosedPipe || err == io.EOF {
				log.Println("closed pipe or eof")
				conn, err = net.Dial("tcp", "0.0.0.0:9866")
				if err != nil {
					log.Fatal(err)
				}
				_, err = conn.Write(buf)
				if err != nil {
					log.Fatal(err)
				}
			} else {
				log.Fatal(err)
			}
		}
		log.Printf("sent %s: %d bytes", string(buf), bytesLength)
		buf1 := make([]byte, binary.MaxVarintLen64)
		readBytes, err := conn.Read(buf1)
		if err != nil {
			if err == io.ErrClosedPipe || err == io.EOF {
				log.Println("closed pipe or eof")
				conn, err = net.Dial("tcp", "0.0.0.0:9866")
				if err != nil {
					log.Fatal(err)
				}
				_, err = conn.Write(buf)
				if err != nil {
					log.Fatal(err)
				}
				readBytes, err = conn.Read(buf1)
				if err != nil {
					log.Fatal(err)
				}
			} else {
				log.Fatal(err)
			}
		}
		bytesToRead, err := binary.ReadVarint(bytes.NewReader(buf1[:readBytes]))
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Message size from server: %d bytes\n", bytesToRead)

		readBuffer := make([]byte, bytesToRead)

		_, err = conn.Read(readBuffer)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Received from server: %s", string(readBuffer))
	}
}
