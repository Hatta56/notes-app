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

func DeleteNote() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		noteID := c.Param("id")
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(noteID)

		result, err := noteCollection.DeleteOne(ctx, bson.M{"id": objId})

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.NoteResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if result.DeletedCount < 1 {
			c.JSON(http.StatusNotFound,
				responses.NoteResponse{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "Note with specified ID not found!"}},
			)
			return
		}

		c.JSON(http.StatusOK,
			responses.NoteResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "Note successfully deleted!"}},
		)
	}
}

func EditNote() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		noteID := c.Param("id")
		var note models.Note
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(noteID)

		//validate the request body
		if err := c.BindJSON(&note); err != nil {
			c.JSON(http.StatusBadRequest, responses.NoteResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&note); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.NoteResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		update := bson.M{"title": note.Title, "content": note.Content}
		result, err := noteCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.NoteResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//get updated user details
		var updatedUser models.Note
		if result.MatchedCount == 1 {
			err := noteCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedUser)
			if err != nil {
				c.JSON(http.StatusInternalServerError, responses.NoteResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
		}

		c.JSON(http.StatusOK, responses.NoteResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": updatedUser}})
	}
}

func GetAllNote() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var notes []models.Note
		defer cancel()

		results, err := noteCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusBadRequest, responses.NoteResponse{Status: http.StatusBadRequest, Message: "error" + err.Error()})
		}
		defer results.Close(ctx)

		for results.Next(ctx) {
			var singleNote models.Note
			if err = results.Decode(&singleNote); err != nil {
				c.JSON(http.StatusBadRequest, responses.NoteResponse{Status: http.StatusBadRequest, Message: "error" + err.Error()})
			}

			notes = append(notes, singleNote)
		}

		c.JSON(http.StatusOK, responses.NoteResponse{Status: http.StatusOK, Message: "Succes Get All", Data: map[string]interface{}{"data": notes}})

	}
}

func GetNotes() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		noteID := c.Param("id")
		var note models.Note
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(noteID)

		err := noteCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&note)
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
			c.JSON(http.StatusBadRequest, responses.NoteResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
		}

		if validationErr := validate.Struct(&note); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.NoteResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
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
