package recipes_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/dawsonc/recipes/src/recipes"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	testDBName, testURI string
)

func init() {
	// Get the database name and uri from environment variables using os
	testDBName = os.Getenv("TEST_DB_NAME")
	testURI = os.Getenv("TEST_DB_URI")

	// If the environment variables are not set, use the default values
	if testDBName == "" {
		testDBName = "recipes_test"
	}
	if testURI == "" {
		testURI = "mongodb://localhost:27017"
	}
}

// Define functions to set up and tear down the test database before and after each
// test
func setupTestDB() error {
	// Set client options
	clientOptions := options.Client().ApplyURI(testURI)

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return fmt.Errorf("failed to connect to test database: %w", err)
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("failed to ping test database: %w", err)
	}

	// Drop the existing test database (if it exists)
	err = client.Database(testDBName).Drop(context.Background())
	if err != nil && err.Error() != "mongo: database not found" {
		return fmt.Errorf("failed to drop test database: %w", err)
	}

	return nil
}

func teardownTestDB() {
	// Set client options
	clientOptions := options.Client().ApplyURI(testURI)

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to test database: %v", err)
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatalf("Failed to ping test database: %v", err)
	}

	// Drop the test database
	err = client.Database(testDBName).Drop(context.Background())
	if err != nil {
		log.Fatalf("Failed to drop test database: %v", err)
	}
}

// TestAddRecipe tests the AddRecipe function
func TestAddRecipe(t *testing.T) {
	// Set up the test database
	err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	defer teardownTestDB()

	// Create a new recipe manager
	recipeManager, err := recipes.CreateMongoRecipeManager(testURI, testDBName, "recipes")
	if err != nil {
		t.Fatalf("Failed to create recipe manager: %v", err)
	}

	// Add the first test recipe
	recipeID, err := recipeManager.AddRecipe(testRecipe1)
	if err != nil {
		t.Fatalf("Failed to add test recipe: %v", err)
	}

	// Check that the recipe was added correctly
	recipe, err := recipeManager.GetRecipeByID(recipeID)
	if err != nil {
		t.Fatalf("Failed to get test recipe: %v", err)
	}
	// The ID may have been updated by the database, so set it to the expected value
	// before comparing
	recipe.ID = testRecipe1.ID
	if !reflect.DeepEqual(recipe, testRecipe1) {
		t.Fatalf("Test recipe was not added correctly")
	}

	// Add the second test recipe
	recipeID, err = recipeManager.AddRecipe(testRecipe2)
	if err != nil {
		t.Fatalf("Failed to add test recipe: %v", err)
	}

	// Check that the recipe was added correctly
	recipe, err = recipeManager.GetRecipeByID(recipeID)
	if err != nil {
		t.Fatalf("Failed to get test recipe: %v", err)
	}
	// The ID may have been updated by the database, so set it to the expected value
	// before comparing
	recipe.ID = testRecipe1.ID
	if !reflect.DeepEqual(recipe, testRecipe2) {
		t.Fatalf("Test recipe was not added correctly")
	}
}

// TestDeleteRecipe tests the DeleteRecipe function
func TestDeleteRecipe(t *testing.T) {
	// Set up the test database
	err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	defer teardownTestDB()

	// Create a new recipe manager
	recipeManager, err := recipes.CreateMongoRecipeManager(testURI, testDBName, "recipes")
	if err != nil {
		t.Fatalf("Failed to create recipe manager: %v", err)
	}

	// Add the first test recipe
	recipeID, err := recipeManager.AddRecipe(testRecipe1)
	if err != nil {
		t.Fatalf("Failed to add test recipe: %v", err)
	}

	// Delete the recipe
	err = recipeManager.DeleteRecipe(recipeID)
	if err != nil {
		t.Fatalf("Failed to delete test recipe: %v", err)
	}

	// Check that the recipe was deleted
	_, err = recipeManager.GetRecipeByID(recipeID)
	if err == nil {
		t.Fatalf("Test recipe was not deleted")
	}
}

// TestUpdateRecipe tests the UpdateRecipe function
func TestUpdateRecipe(t *testing.T) {
	// Set up the test database
	err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	defer teardownTestDB()

	// Create a new recipe manager
	recipeManager, err := recipes.CreateMongoRecipeManager(testURI, testDBName, "recipes")
	if err != nil {
		t.Fatalf("Failed to create recipe manager: %v", err)
	}

	// Add the first test recipe
	recipeID, err := recipeManager.AddRecipe(testRecipe1)
	if err != nil {
		t.Fatalf("Failed to add test recipe: %v", err)
	}

	// Update the recipe
	var testRecipe1Copy = testRecipe1
	testRecipe1Copy.ID, _ = primitive.ObjectIDFromHex(recipeID)
	testRecipe1Copy.Name = "Updated Test Recipe"
	err = recipeManager.UpdateRecipe(testRecipe1Copy)
	if err != nil {
		t.Fatalf("Failed to update test recipe: %v", err)
	}

	// Check that the recipe was updated
	recipe, err := recipeManager.GetRecipeByID(recipeID)
	if err != nil {
		t.Fatalf("Failed to get test recipe: %v", err)
	}
	// The ID should not have been updated when updating the recipe
	if !reflect.DeepEqual(recipe, testRecipe1Copy) {
		t.Fatalf("Test recipe was not updated correctly."+
			"Expected: %v, got: %v", testRecipe1Copy, recipe)
	}
}

// TestGetAllRecipes tests the GetAllRecipes function
func TestGetAllRecipes(t *testing.T) {
	// Set up the test database
	err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	defer teardownTestDB()

	// Create a new recipe manager
	recipeManager, err := recipes.CreateMongoRecipeManager(testURI, testDBName, "recipes")
	if err != nil {
		t.Fatalf("Failed to create recipe manager: %v", err)
	}

	// Add the first test recipe
	recipeID1, err := recipeManager.AddRecipe(testRecipe1)
	if err != nil {
		t.Fatalf("Failed to add test recipe: %v", err)
	}

	// Add the second test recipe
	recipeID2, err := recipeManager.AddRecipe(testRecipe2)
	if err != nil {
		t.Fatalf("Failed to add test recipe: %v", err)
	}

	// Get all recipes
	recipes, err := recipeManager.GetAllRecipes()
	if err != nil {
		t.Fatalf("Failed to get all recipes: %v", err)
	}

	// Check that the recipes were retrieved correctly
	if len(recipes) != 2 {
		t.Fatalf("Incorrect number of recipes retrieved")
	}
	objID1, _ := primitive.ObjectIDFromHex(recipeID1)
	objID2, _ := primitive.ObjectIDFromHex(recipeID2)
	if recipes[0].ID != objID1 && recipes[1].ID != objID1 {
		t.Fatalf("Test recipe 1 was not retrieved")
	}
	if recipes[0].ID != objID2 && recipes[1].ID != objID2 {
		t.Fatalf("Test recipe 2 was not retrieved")
	}
}

// TestGetRecipesByTags tests the GetRecipesByTags function
func TestGetRecipesByTags(t *testing.T) {
	// Set up the test database
	err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	defer teardownTestDB()

	// Create a new recipe manager
	recipeManager, err := recipes.CreateMongoRecipeManager(testURI, testDBName, "recipes")
	if err != nil {
		t.Fatalf("Failed to create recipe manager: %v", err)
	}

	// Add the first test recipe
	recipeID1, err := recipeManager.AddRecipe(testRecipe1)
	if err != nil {
		t.Fatalf("Failed to add test recipe: %v", err)
	}

	// Add the second test recipe
	recipeID2, err := recipeManager.AddRecipe(testRecipe2)
	if err != nil {
		t.Fatalf("Failed to add test recipe: %v", err)
	}

	// Get recipes by a tag that both recipes have
	tags := []string{"Test Tag 1"}
	recipes, err := recipeManager.GetRecipesByTags(tags)
	if err != nil {
		t.Fatalf("Failed to get recipes with tags %v: %v", tags, err)
	}

	// This should retrieve both recipes
	if len(recipes) != 2 {
		t.Fatalf("Incorrect number of recipes retrieved with tags %v", tags)
	}
	objID1, _ := primitive.ObjectIDFromHex(recipeID1)
	objID2, _ := primitive.ObjectIDFromHex(recipeID2)
	if recipes[0].ID != objID1 && recipes[1].ID != objID1 {
		t.Fatalf("Test recipe 1 was not retrieved with tags %v", tags)
	}
	if recipes[0].ID != objID2 && recipes[1].ID != objID2 {
		t.Fatalf("Test recipe 2 was not retrieved with tags %v", tags)
	}

	// Get recipes by a tag that only one recipe has
	tags = []string{"Test Tag 2"}
	recipes, err = recipeManager.GetRecipesByTags(tags)
	if err != nil {
		t.Fatalf("Failed to get recipes with tags %v: %v", tags, err)
	}

	// This should retrieve only one recipe
	if len(recipes) != 1 {
		t.Fatalf("Incorrect number of recipes retrieved with tags %v", tags)
	}
	if recipes[0].ID != objID1 {
		t.Fatalf("Test recipe 1 was not retrieved with tags %v", tags)
	}
}

// TestSearchRecipes tests the SearchRecipes function
func TestSearchRecipes(t *testing.T) {
	// Set up the test database
	err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	defer teardownTestDB()

	// Create a new recipe manager
	recipeManager, err := recipes.CreateMongoRecipeManager(testURI, testDBName, "recipes")
	if err != nil {
		t.Fatalf("Failed to create recipe manager: %v", err)
	}

	// Add the first test recipe
	recipeID1, err := recipeManager.AddRecipe(testRecipe1)
	if err != nil {
		t.Fatalf("Failed to add test recipe: %v", err)
	}

	// Add the second test recipe
	recipeID2, err := recipeManager.AddRecipe(testRecipe2)
	if err != nil {
		t.Fatalf("Failed to add test recipe: %v", err)
	}

	// Search for recipes with the query "recipe" with no tags (should match both)
	query := "recipe"
	tags := []string{}
	recipes, err := recipeManager.SearchRecipes(query, tags)
	if err != nil {
		t.Fatalf("Failed to search for recipes (query %v, tags %v): %v", query, tags, err)
	}

	// This should retrieve both recipes
	if len(recipes) != 2 {
		t.Fatalf("Incorrect number of recipes retrieved (query %v, tags %v): "+
			"expected %v, got %v", query, tags, 2, len(recipes))
	}
	objID1, _ := primitive.ObjectIDFromHex(recipeID1)
	objID2, _ := primitive.ObjectIDFromHex(recipeID2)
	if recipes[0].ID != objID1 && recipes[1].ID != objID1 {
		t.Fatalf("Test recipe 1 was not retrieved (query %v, tags %v)", query, tags)
	}
	if recipes[0].ID != objID2 && recipes[1].ID != objID2 {
		t.Fatalf("Test recipe 2 was not retrieved (query %v, tags %v)", query, tags)
	}

	// Now filter the search to only recipes with the tag "Test Tag 2" (should only
	// match the first recipe)
	tags = []string{"Test Tag 2"}
	recipes, err = recipeManager.SearchRecipes(query, tags)
	if err != nil {
		t.Fatalf("Failed to search for recipes (query %v, tags %v): %v", query, tags, err)
	}

	// This should retrieve only one recipe
	if len(recipes) != 1 {
		t.Fatalf("Incorrect number of recipes retrieved (query %v, tags %v): "+
			"expected %v, got %v", query, tags, 1, len(recipes))
	}
	if recipes[0].ID != objID1 {
		t.Fatalf("Test recipe 1 was not retrieved (query %v, tags %v)", query, tags)
	}
}

// TestGetTags tests the GetTags function
func TestGetTags(t *testing.T) {
	// Set up the test database
	err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	defer teardownTestDB()

	// Create a new recipe manager
	recipeManager, err := recipes.CreateMongoRecipeManager(testURI, testDBName, "recipes")
	if err != nil {
		t.Fatalf("Failed to create recipe manager: %v", err)
	}

	// Add the first test recipe
	_, err = recipeManager.AddRecipe(testRecipe1)
	if err != nil {
		t.Fatalf("Failed to add test recipe: %v", err)
	}

	// Add the second test recipe
	_, err = recipeManager.AddRecipe(testRecipe2)
	if err != nil {
		t.Fatalf("Failed to add test recipe: %v", err)
	}

	// Get the tags
	tags, err := recipeManager.GetTags()
	if err != nil {
		t.Fatalf("Failed to get tags: %v", err)
	}

	// There should be 3 tags
	if len(tags) != 3 {
		t.Fatalf("Incorrect number of tags retrieved: expected %v, got %v", 3, len(tags))
	}

	// The tags "Test Tag 1", "Test Tag 2", and "Test Tag 3" should all be in this slice
	if !isMember(tags, "Test Tag 1") {
		t.Fatalf("Tag \"Test Tag 1\" was not retrieved, got %v", tags)
	}
	if !isMember(tags, "Test Tag 2") {
		t.Fatalf("Tag \"Test Tag 2\" was not retrieved, got %v", tags)
	}
	if !isMember(tags, "Test Tag 3") {
		t.Fatalf("Tag \"Test Tag 2\" was not retrieved, got %v", tags)
	}
}

// isMember returns true if the given value is in the given slice
func isMember(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}
