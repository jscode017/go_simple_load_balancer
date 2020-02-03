package main

import (
	"bufio"
	"io"
	"log"
	"net"
)

var (
	seed int
)

type LoadBalancer struct {
	ServersIp []string
	Address   string
}

func NewLoadBalancer() *LoadBalancer {
	serversIp := []string{"localhost:8081", "localhost:8082"}
	return &LoadBalancer{
		ServersIp: serversIp,
		Address:   "localhost:8085",
	}
}

func (lb *LoadBalancer) Run() {
	listener, err := net.Listen("tcp", lb.Address)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("accept a connection")
		go lb.ReverseProxy(conn)
		seed++
	}
}

func (lb *LoadBalancer) ReverseProxy(conn net.Conn) {
	defer conn.Close()
	dst, err := net.Dial("tcp", lb.ServersIp[seed%2])
	if err != nil {
		log.Fatal(err)
	}

	clientReader := bufio.NewReader(conn)
	log.Println("directing")
	go io.Copy(conn, dst)
	io.Copy(dst, clientReader)

}
