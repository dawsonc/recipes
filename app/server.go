package main

import (
	"github.com/gin-gonic/gin"

	"github.com/dawsonc/recipes/src/recipes"
)

func main() {
	// Make a router
	router := gin.Default()

	// Create the recipes manager and connect to MongoDB
	recipe_manager, err := recipes.CreateMongoRecipeManager("mongodb://localhost:27017", "recipes", "recipes")
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
