package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	recipes "github.com/dawsonc/recipes/src"
)

type PageData struct {
	Recipes []recipes.Recipe
	Tags    map[string]bool
}

func main() {
	// Make a router and load the HTML templates
	router := gin.Default()
	router.LoadHTMLGlob("templates/*.html")

	// Connect to a local MongoDB instance
	client, err := recipes.ConnectToMongoDB("mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}

	// The homepage should list all recipes
	router.GET("/", func(c *gin.Context) {
		// Get all recipes from the database
		recipes, err := recipes.GetAllRecipes(client, "recipes", "recipes")
		if err != nil {
			// Display an error page
			c.HTML(http.StatusInternalServerError, "error.html", nil)
		}

		// Get a slice of all unique tags
		tags := make(map[string]bool)
		for _, recipe := range recipes {
			for _, tag := range recipe.Tags {
				tags[tag] = true
			}
		}

		data := PageData{
			Recipes: recipes,
			Tags:    tags,
		}
		c.HTML(http.StatusOK, "index.html", data)
	})

	// A page for creating new recipes
	router.GET("/add_recipe", func(c *gin.Context) {
		c.HTML(http.StatusOK, "new_recipe.html", nil)
	})

	// Define a REST API for recipes
	recipesAPI := router.Group("/api/recipes")
	{
		// GET /api/recipes - get all recipes
		recipesAPI.GET("/", func(c *gin.Context) {
			// Get all recipes from the database
			recipes, err := recipes.GetAllRecipes(client, "recipes", "recipes")
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, recipes)
		})

		// GET /api/recipes/:id - get a recipe by ID
		recipesAPI.GET("/:id", func(c *gin.Context) {
			// Get the ID from the URL
			id := c.Param("id")

			// Get the recipe from the database
			recipe, err := recipes.GetRecipeByID(client, "recipes", "recipes", id)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, recipe)
		})

		// GET /api/recipes/tags/:tag - get all recipes with a given tag
		recipesAPI.GET("/tags/:tag", func(c *gin.Context) {
			// Get the tag from the URL
			tag := c.Param("tag")

			// Get all recipes from the database
			recipes, err := recipes.GetRecipesByTag(client, "recipes", "recipes", tag)
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
			if err := recipes.InsertRecipe(client, "recipes", "recipes", recipe); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

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

			// Update the recipe in the database
			if err := recipes.UpdateRecipe(client, "recipes", "recipes", id, recipe); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{"message": "Recipe updated successfully"})
		})

		// DELETE /api/recipes/:id - delete a recipe by ID
		recipesAPI.DELETE("/:id", func(c *gin.Context) {
			// Get the ID from the URL
			id := c.Param("id")

			// Delete the recipe from the database
			if err := recipes.DeleteRecipe(client, "recipes", "recipes", id); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{"message": "Recipe deleted successfully"})
		})
	}

	// Run the server
	router.Run(":8080")
}
