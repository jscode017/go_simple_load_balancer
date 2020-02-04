package main

import (
	"bufio"
	"context"
	"io"
	"log"
	"math"
	"net"
	"sync"
)

type LoadBalancer struct {
	ServerIps []string
	Address   string
	Seed      uint32
	sync.RWMutex
}

func NewLoadBalancer(serverIPs []string) *LoadBalancer {
	//serversIp := []string{"localhost:8081", "localhost:8082"}
	return &LoadBalancer{
		ServersIp: serverIPs,
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

	}
}

func (lb *LoadBalancer) ReverseProxy(conn net.Conn) {
	defer conn.Close()
	lb.Lock()
	curSeed := atomic.LoadUint32(&lb.Seed) //though is a read op, but still want to get an unique seed
	serverIPIndex := lb.Seed % len(lb.ServerIps)
	serverIP := lb.ServerIps[serverIPIndex]
	lb.Seed++ //do not need to handle overflow, it would simply become 0 if overflow
	lb.Unlock()
	dst, err := net.Dial("tcp", serverIP)
	if err != nil {
		log.Fatal(err)
	}

	lbReader := bufio.NewReader(conn)
	log.Println("directing")
	ctx := context.WithCancel(context.Background())
	go io.Copy(conn, dst)
	io.Copy(dst, lbReader)

}
