package main

// a simple server for testing
import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	loadbalancer "github.com/jscode017/go_simple_load_balancer"
	"net/http"
)

func main() {
	lb := loadbalancer.NewLoadBalancer([]string{"localhost:8081", "localhost:8082", "localhost:8083"}, "localhost:8085", 10, 2)
	go lb.Run()
	go lb.RunHttp(":8091")
	r1 := gin.Default()
	r2 := gin.Default()
	r3 := gin.Default()

	r1.GET("/", Hello1)
	r2.GET("/", Hello2)
	r3.GET("/", Hello3)
	r1.GET("/hello", Hello2)
	r2.GET("/hello", Hello1)
	r3.GET("/hello", Hello3)
	r1.POST("/p", Hi)
	r2.POST("/p", Hi)
	r3.POST("/p", Hi)
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
func Hi(c *gin.Context) {
	fmt.Println(c.Request.Body)
	var requestBody map[string]interface{}
	err := json.NewDecoder(c.Request.Body).Decode(&requestBody)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"result": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"result": true, "body": fmt.Sprintf("%v", requestBody)})
}
