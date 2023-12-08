package routes

import (
	"notes/controllers"

	"github.com/gin-gonic/gin"
)

func NoteRoute(router *gin.Engine) {
	//All routes related to users comes here
	router.POST("/note", controllers.CreateNote())
	router.GET("/note/:id", controllers.GetNotes())
}
