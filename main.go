package main

import (
	"fmt"
	"os"
	"time"

	"github.com/devmatteo/go-project/env"
	"github.com/devmatteo/go-project/models"
	"github.com/devmatteo/go-project/routes/auth"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func worker(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Println("worker", id, "started  job", j)
		fmt.Println("worker", id, "finished job", j)
		results <- j * 2
	}
}

func main() {
	burstyLimiter := make(chan time.Time, 3)

	for i := 0; i < 3; i++ {
		burstyLimiter <- time.Now()
	}

	go func() {
		for t := range time.Tick(200 * time.Millisecond) {
			burstyLimiter <- t
		}
	}()

	burstyRequests := make(chan int, 5)

	for i := 1; i <= 5; i++ {
		burstyRequests <- i
	}
	close(burstyRequests)

	for req := range burstyRequests {
		<-burstyLimiter
		fmt.Println("request", req, time.Now())
	}

	err := env.ValidateEnv()

	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(mysql.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&models.User{})

	if err != nil {
		panic(err)
	}

	router := gin.Default()

	auth.Routes(router, db)

	router.Run()
}
