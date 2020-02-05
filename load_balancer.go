package main

import (
	"bufio"
	//"context"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

type LoadBalancer struct {
	ServerIps         []string
	Address           string
	Seed              uint32
	HeartBeatDuration time.Duration
	ServerTimeOut     time.Duration
	sync.RWMutex
}

func NewLoadBalancer(serverIPs []string, heartBeatDuration int, serverTimeOut int) *LoadBalancer {
	//serversIp := []string{"localhost:8081", "localhost:8082"}
	return &LoadBalancer{
		ServerIps:         serverIPs,
		Address:           "localhost:8085",
		HeartBeatDuration: time.Duration(heartBeatDuration) * time.Second,
		ServerTimeOut:     time.Duration(serverTimeOut) * time.Second,
	}
}

func (lb *LoadBalancer) Run() {
	listener, err := net.Listen("tcp", lb.Address)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Println(err)
				continue
			}
			log.Println("accept a connection")
			go lb.ReverseProxy(conn)

		}
	}()
	ticker := time.NewTicker(lb.HeartBeatDuration)
	defer ticker.Stop()
	select {
	case <-ticker.C:
		serverIpIndex := 0
		for serverIpIndex < len(lb.ServerIps) {
			log.Println("sending heart beat to: ", lb.ServerIps[serverIpIndex])
			err = lb.SendHeartBeat(serverIpIndex)
			if err != nil {
				log.Println("error from server ", lb.ServerIps[serverIpIndex], err)
				lb.Lock()
				lb.ServerIps = append(lb.ServerIps[:serverIpIndex], lb.ServerIps[serverIpIndex+1:]...)
				lb.Unlock()
			} else {
				log.Println(lb.ServerIps[serverIpIndex], " is healthy")
				serverIpIndex++
			}
		}
	}
}

func (lb *LoadBalancer) ReverseProxy(conn net.Conn) {
	defer conn.Close()
	lb.Lock()
	serverIpIndex := lb.Seed % uint32(len(lb.ServerIps)) //though is a read op, but still want to get an unique seed, so use lock instead of rlock
	serverIP := lb.ServerIps[serverIpIndex]
	lb.Seed++ //do not need to handle overflow, it would simply become 0 if overflow
	lb.Unlock()
	dst, err := net.Dial("tcp", serverIP)
	if err != nil {
		log.Fatal(err)
	}

	lbReader := bufio.NewReader(conn)
	log.Println("directing")
	//ctx,cancel := context.WithCancel(context.Background())
	go io.Copy(conn, dst)
	io.Copy(dst, lbReader)

}

func (lb *LoadBalancer) SendHeartBeat(index int) error {
	serverIP := lb.ServerIps[index]
	conn, err := net.DialTimeout("tcp", serverIP, lb.ServerTimeOut)
	defer conn.Close()
	if err != nil {
		return err
	}

	return nil

}
