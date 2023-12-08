package controllers

import (
	"context"
	"net/http"
	"notes/configs"
	"notes/models"
	"notes/responses"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var noteCollection *mongo.Collection = configs.GetCollection(configs.DB, "notes")
var validate = validator.New()

func GetNotes() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		noteID := c.Param("id")
		var note models.Note
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(noteID)

		err := noteCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&note)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.NoteResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusOK, responses.NoteResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": note}})
	}
}

func CreateNote() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var note models.Note
		defer cancel()

		if err := c.BindJSON(&note); err != nil {

			c.JSON(http.StatusBadRequest, responses.NoteResponse{Status: http.StatusBadRequest, Message: "error" + err.Error()})
		}

		if validationErr := validate.Struct(&note); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.NoteResponse{Status: http.StatusBadRequest, Message: "error :" + validationErr.Error()})
		}

		currentTime := primitive.NewDateTimeFromTime(time.Now())
		note.CreatedAt = currentTime
		note.UpdatedAt = currentTime
		newNote := models.Note{
			ID:        primitive.NewObjectID(),
			Title:     note.Title,
			Content:   note.Content,
			CreatedAt: note.CreatedAt,
			UpdatedAt: note.UpdatedAt,
		}

		result, err := noteCollection.InsertOne(ctx, newNote)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.NoteResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
		}

		c.JSON(http.StatusCreated, responses.NoteResponse{Status: http.StatusCreated, Message: "Succes Created", Data: map[string]interface{}{"data": result}})
	}
}
