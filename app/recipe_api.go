package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/dawsonc/recipes/src/recipes"
)

func AddRecipesAPI(router *gin.Engine, recipe_manager recipes.RecipeManager) {
	// Provide a RESTful API for recipes
	recipesAPI := router.Group("/api/recipes")
	{
		// GET /api/recipes - get recipes, possibly with filters
		// e.g. /api/recipes?id=ID
		// e.g. /api/recipes?tags=tag1,tag2
		// e.g. /api/recipes?q=search_term&tags=tag1,tag2
		recipesAPI.GET("/", func(c *gin.Context) {
			// Allow optional query strings
			id := c.Query("id")
			tags := strings.Split(c.Query("tags"), ",")
			search_term := c.Query("q")

			// If the tags string is empty, then tags will have a single element that
			// is an empty string. Remove this element.
			if len(tags) == 1 && tags[0] == "" {
				tags = []string{}
			}

			// Get recipes based on the provided queries
			var queried_recipes []recipes.Recipe
			var err error
			switch {
			case id != "":
				fmt.Println("ID: ", id)
				// Get the recipe from the database with the given ID
				var recipe recipes.Recipe
				recipe, err = recipe_manager.GetRecipeByID(id)
				// Wrap the recipe in a single element slice
				queried_recipes = []recipes.Recipe{recipe}
			case search_term != "":
				fmt.Println("Search term: ", search_term, " Tags: ", tags)
				// Get all recipes that match the given tags and search term
				queried_recipes, err = recipe_manager.SearchRecipes(search_term, tags)
			case len(tags) > 0:
				// Get all recipes that match the given tags
				queried_recipes, err = recipe_manager.GetRecipesByTags(tags)
			default:
				fmt.Println("No query strings")
				// Get all recipes from the database
				queried_recipes, err = recipe_manager.GetAllRecipes()
			}

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, queried_recipes)
		})

		// GET /api/recipes/tags - get all tags
		recipesAPI.GET("/tags", func(c *gin.Context) {
			// Get all tags from the database
			tags, err := recipe_manager.GetTags()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, tags)
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
			id, err := recipe_manager.AddRecipe(recipe)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{"message": "Recipe created successfully", "id": id})
		})

		// GET /api/recipes/:id - get a recipe by ID
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

		// PUT /api/recipes - update a recipe
		recipesAPI.PUT("/id/:id", func(c *gin.Context) {
			// Get the ID from the URL
			id := c.Param("id")

			// Get the recipe from the request
			var recipe recipes.Recipe
			if err := c.BindJSON(&recipe); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			// Make sure the ID in the URL matches the ID in the recipe
			recipe.SetID(id)

			// Update the recipe
			err := recipe_manager.UpdateRecipe(recipe)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{"message": "Recipe updated successfully"})
		})

		// DELETE /api/recipes/:id - delete a recipe by ID
		recipesAPI.DELETE("/id/:id", func(c *gin.Context) {
			// Get the ID from the URL
			id := c.Param("id")

			// Delete the recipe from the database
			err := recipe_manager.DeleteRecipe(id)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{"message": "Recipe deleted successfully"})
		})
	}
}
