package main

import (
	"FORUM/app/midwares"
	"FORUM/config/database"
	"FORUM/config/router"
	"FORUM/config/session"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	database.Init()
	r := gin.Default()
	r.NoMethod(midwares.HandleNotFound)
	r.NoRoute(midwares.HandleNotFound)
	session.Init(r)
	router.Init(r)

	err := r.Run()
	if err != nil {
		log.Fatal("Sener start error", err)
	}
}



