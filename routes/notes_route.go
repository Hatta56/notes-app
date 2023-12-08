package routes

import (
	"notes/controllers"

	"github.com/gin-gonic/gin"
)

func NoteRoute(router *gin.Engine) {
	router.POST("/note", controllers.CreateNote())
	router.GET("/note/:id", controllers.GetNotes())
	router.GET("/note", controllers.GetAllNote())
	router.PUT("/note/:id", controllers.EditNote())
	router.DELETE("/note/:id", controllers.DeleteNote())
}
