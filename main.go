package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"net"
	"os"
)

func main() {
	log.SetFlags(0)
	log.SetOutput(os.Stdin)
	listener, err := net.Listen("tcp", ":56000")
	if err != nil {
		log.Fatal(err)
	}

	go client("tcp", ":56000")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go handleConnection(conn)
	}

}

func client(network string, port string) {
	log.Print("Client: ")
	reader := bufio.NewReader(os.Stdin)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			log.Fatal(err)
		}

		conn, err := net.Dial(network, port)
		if err != nil {
			log.Fatal(err)
		}
		_, err = conn.Write([]byte(line))
		if err != nil {
			log.Fatal(err)
		}
		conn.Close()
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Server: \n" + string(bytes)+"\n")
	log.Print("Client: ")
}
