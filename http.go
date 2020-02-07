package go_simple_load_balancer

//package main

import (
	"log"
	"net/http"
	"net/http/httputil"
)

func (lb *LoadBalancer) RunHttp(proxyIP string) {

	proxy := &httputil.ReverseProxy{Director: func(req *http.Request) {
		req.URL.Scheme = "http"
		serverIP := lb.RRAlgorithm()
		req.URL.Host = serverIP

		log.Println(req)
	},
	}
	http.DefaultTransport.(*http.Transport).DisableKeepAlives = true
	go http.ListenAndServe(proxyIP, proxy) //ignore err
	select {}
}
