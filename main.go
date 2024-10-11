package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type App struct {
	router      *gin.Engine
	client      *mongo.Client
	userCollection *mongo.Collection
	clientCollection *mongo.Collection
	projectCollection *mongo.Collection
	// ... other collections
}

func NewApp() *App {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("MONGODB_URI environment variable is not set")
	}

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database("freelance-platform")
	userCollection := db.Collection("users")
	clientCollection := db.Collection("clients")
	projectCollection := db.Collection("projects")

	return &App{
		router:      gin.Default(),
		client:      client,
		userCollection: userCollection,
		clientCollection: clientCollection,
		projectCollection: projectCollection,
		// ... other collections
	}
}

// Authentication Middleware
func (app *App) authMiddleware(c *gin.Context) {
	// ... Implement your authentication logic here ...
	// Example: Check for a valid JWT token in the Authorization header

	c.Next()
}

func (app *App) Run() {
	app.router.Use(app.authMiddleware)
	app.setupRoutes()
	err := app.router.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}

func (app *App) setupRoutes() {
	// Client Routes
	app.router.GET("/clients", app.getAllClients)
	app.router.GET("/clients/:id", app.getClientByID)
	app.router.POST("/clients", app.createClient)
	app.router.PUT("/clients/:id", app.updateClient)
	app.router.DELETE("/clients/:id", app.deleteClient)

	// ... Other routes
}

// Client Handlers
func (app *App) getAllClients(c *gin.Context) {
	var clients []models.Client
	cursor, err := app.clientCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var client models.Client
		err := cursor.Decode(&client)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		clients = append(clients, client)
	}

	c.JSON(http.StatusOK, clients)
}

func (app *App) getClientByID(c *gin.Context) {
	id := c.Param("id")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var client models.Client
	err = app.clientCollection.FindOne(context.TODO(), bson.M{"_id": objectId}).Decode(&client)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Client not found"})
		return
	}

	c.JSON(http.StatusOK, client)
}

func (app *App) createClient(c *gin.Context) {
	var client models.Client
	if err := c.ShouldBindJSON(&client); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := app.clientCollection.InsertOne(context.TODO(), client)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": result.InsertedID})
}

func (app *App) updateClient(c *gin.Context) {
	id := c.Param("id")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var client models.Client
	if err := c.ShouldBindJSON(&client); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filter := bson.M{"_id": objectId}
	update := bson.M{"$set": client}
	result, err := app.clientCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Client not found"})
		return
	}

	c.Status(http.StatusOK)
}

func (app *App) deleteClient(c *gin.Context) {
	id := c.Param("id")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := app.clientCollection.DeleteOne(context.TODO(), bson.M{"_id": objectId})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Client not found"})
		return
	}

	c.Status(http.StatusOK)
}

// User Handlers
// ... (Implement User CRUD handlers similarly to Client handlers)

// Project Handlers
// ... (Implement Project CRUD handlers similarly to Client handlers)

func main() {
	app := NewApp()
	app.Run()
}