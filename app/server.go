package main

import (
	"os"

	"github.com/gin-gonic/gin"

	"github.com/dawsonc/recipes/src/recipes"
)

func main() {
	// Make a router
	router := gin.Default()

	// Create the recipes manager and connect to MongoDB
	dbURI := os.Getenv("DB_URI")
	// If the environment variables are not set, use default values
	if dbURI == "" {
		dbURI = "mongodb://localhost:27017"
	}
	dbName := "recipes"
	recipe_manager, err := recipes.CreateMongoRecipeManager(dbURI, dbName, "recipes")
	if err != nil {
		panic(err)
	}

	// Provide a RESTful API for recipes
	AddRecipesAPI(router, recipe_manager)

	// Serve frontend files
	router.Static("/app", "./frontend")

	// Redirect the homepage to the app
	router.GET("/", func(c *gin.Context) {
		c.Redirect(301, "/app")
	})

	// Run the server
	router.Run(":8080")
}
