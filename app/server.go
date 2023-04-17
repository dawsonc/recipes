package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/dawsonc/recipes/src/recipes"
)

type PageData struct {
	Recipes []recipes.Recipe
	Tags    map[string]bool
}

func main() {
	// Make a router
	router := gin.Default()

	// Create the recipes manager and connect to MongoDB
	recipe_manager, err := recipes.CreateMongoRecipeManager("mongodb://localhost:27017", "recipes", "recipes")
	if err != nil {
		panic(err)
	}

	// Provide a RESTful API for recipes
	recipesAPI := router.Group("/api/recipes")
	{
		// GET /api/recipes - get all recipes
		recipesAPI.GET("/", func(c *gin.Context) {
			// Get all recipes from the database
			recipes, err := recipe_manager.GetAllRecipes()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, recipes)
		})

		// GET /api/recipes/id/:id - get a recipe by ID
		recipesAPI.GET("/id/:id", func(c *gin.Context) {
			// Get the ID from the URL
			id := c.Param("id")

			// Get the recipe from the database
			recipe, err := recipe_manager.GetRecipeByID(id)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, recipe)
		})

		// GET /api/recipes/tags/:tag - get all recipes with a given list of comma-separated tags
		recipesAPI.GET("/tags/:tags", func(c *gin.Context) {
			// Get the tag from the URL, splitting along commas
			tags := c.Param("tags")
			tag_list := strings.Split(tags, ",")

			// Get all recipes from the database
			recipes, err := recipe_manager.GetRecipesByTags(tag_list)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, recipes)
		})

		// POST /api/recipes - create a new recipe
		recipesAPI.POST("/", func(c *gin.Context) {
			// Get the recipe from the request
			var recipe recipes.Recipe
			if err := c.BindJSON(&recipe); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			// Insert the recipe into the database
			// TODO

			c.JSON(http.StatusOK, gin.H{"message": "Recipe created successfully"})
		})

		// PUT /api/recipes/:id - update a recipe by ID
		recipesAPI.PUT("/:id", func(c *gin.Context) {
			// Get the ID from the URL
			id := c.Param("id")

			// Get the recipe from the request
			var recipe recipes.Recipe
			if err := c.BindJSON(&recipe); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			// TODO make update take a string ID

			// Update the recipe in the database
			// TODO

			c.JSON(http.StatusOK, gin.H{"message": "Recipe updated successfully"})
		})

		// DELETE /api/recipes/:id - delete a recipe by ID
		recipesAPI.DELETE("/:id", func(c *gin.Context) {
			// Get the ID from the URL
			id := c.Param("id")

			// Delete the recipe from the database
			// TODO

			c.JSON(http.StatusOK, gin.H{"message": "Recipe deleted successfully"})
		})
	}

	// Run the server
	router.Run(":8080")
}
