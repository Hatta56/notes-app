package main

import (
	"notes/configs" //add this
	"notes/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	//run database
	configs.ConnectDB()
	//routes
	routes.NoteRoute(router)
	router.Run("localhost:6000")
}
