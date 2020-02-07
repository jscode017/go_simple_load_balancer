package main

// a simple server for testing
import (
	"github.com/gin-gonic/gin"
	loadbalancer "github.com/jscode017/go_simple_load_balancer"
)

func main() {
	lb := loadbalancer.NewLoadBalancer([]string{"localhost:8081", "localhost:8082", "localhost:8083"}, 30, 2)
	go lb.Run()
	go lb.RunHttp()
	r1 := gin.Default()
	r2 := gin.Default()
	r3 := gin.Default()

	r1.GET("/", Hello1)
	r2.GET("/", Hello2)
	r3.GET("/hello3", Hello3)
	go r1.Run(":8081")
	go r2.Run(":8082")
	r3.Run(":8083")
}

func Hello1(c *gin.Context) {
	c.JSON(200, gin.H{"result": true, "hello": "1"})
	return
}

func Hello2(c *gin.Context) {
	c.JSON(200, gin.H{"result": true, "hello": "2"})
	return
}

func Hello3(c *gin.Context) {
	c.JSON(200, gin.H{"result": true, "hello": "3"})
	return
}
