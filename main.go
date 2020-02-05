package main

// a simple server for testing
import (
	"github.com/gin-gonic/gin"
)

func main() {
	lb := NewLoadBalancer([]string{"localhost:8081", "localhost:8082"}, 10, 2)
	go lb.Run()
	r1 := gin.Default()
	r2 := gin.Default()

	r1.GET("/", Hello1)
	r2.GET("/", Hello2)
	go r1.Run(":8081")
	r2.Run(":8082")
}

func Hello1(c *gin.Context) {
	c.JSON(200, gin.H{"result": true, "hello": "1"})
	return
}

func Hello2(c *gin.Context) {
	c.JSON(200, gin.H{"result": true, "hello": "2"})
	return
}
