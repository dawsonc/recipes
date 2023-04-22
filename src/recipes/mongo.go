package recipes

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Define a MongoDB recipe manager that implements the RecipeManager interface
type MongoRecipeManager struct {
	client         *mongo.Client
	dbName         string
	collectionName string
}

// CreateMongoRecipeManager creates a new MongoDB recipe manager
func CreateMongoRecipeManager(uri, dbName, collectionName string) (*MongoRecipeManager, error) {
	// Set client options
	clientOptions := options.Client().ApplyURI(uri)

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	// Create a new recipe manager
	recipeManager := &MongoRecipeManager{
		client:         client,
		dbName:         dbName,
		collectionName: collectionName,
	}

	return recipeManager, nil
}

// AddRecipe adds a recipe to the recipe manager and returns the ID of the new recipe
func (m *MongoRecipeManager) AddRecipe(recipe Recipe) (string, error) {
	// Get the collection handle
	collection := m.client.Database(m.dbName).Collection(m.collectionName)

	// Use a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Add the recipe
	result, err := collection.InsertOne(ctx, recipe)
	if err != nil {
		return "", err
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

// DeleteRecipe deletes a recipe from the recipe manager
func (m *MongoRecipeManager) DeleteRecipe(id string) error {
	// Get the collection handle
	collection := m.client.Database(m.dbName).Collection(m.collectionName)

	// Use a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Delete the recipe
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}

	return nil
}

// UpdateRecipe updates a recipe in the recipe manager
func (m *MongoRecipeManager) UpdateRecipe(recipe Recipe) error {
	// Get the collection handle
	collection := m.client.Database(m.dbName).Collection(m.collectionName)

	// Use a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Update the recipe
	_, err := collection.UpdateOne(ctx, bson.M{"_id": recipe.ID}, bson.M{"$set": recipe})
	if err != nil {
		return err
	}

	return nil
}

// GetAllRecipes returns all recipes in the recipe manager
func (m *MongoRecipeManager) GetAllRecipes() ([]Recipe, error) {
	// Get the collection handle
	collection := m.client.Database(m.dbName).Collection(m.collectionName)

	// Use a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Find all documents in the collection
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Decode each document into a Recipe object
	var recipes []Recipe
	for cursor.Next(ctx) {
		var recipe Recipe
		if err := cursor.Decode(&recipe); err != nil {
			return nil, err
		}
		recipes = append(recipes, recipe)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return recipes, nil
}

// GetRecipeByID returns a recipe with the given ID
func (m *MongoRecipeManager) GetRecipeByID(id string) (Recipe, error) {
	// Get the collection handle
	collection := m.client.Database(m.dbName).Collection(m.collectionName)

	// Use a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Find the document with the given ID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return Recipe{}, err
	}
	var recipe Recipe
	if err := collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&recipe); err != nil {
		return recipe, err
	}

	return recipe, nil
}

// GetRecipesByTags returns all recipes with the given tags
func (m *MongoRecipeManager) GetRecipesByTags(tags []string) ([]Recipe, error) {
	// Get the collection handle
	collection := m.client.Database(m.dbName).Collection(m.collectionName)

	// Use a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// create a filter to match recipes with the given tags
	filter := bson.M{"tags": bson.M{"$all": tags}}

	// execute the find operation and get the result cursor
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	// iterate over the cursor and decode each document into a Recipe object
	var recipes []Recipe
	for cursor.Next(ctx) {
		var recipe Recipe
		err := cursor.Decode(&recipe)
		if err != nil {
			return nil, err
		}
		recipes = append(recipes, recipe)
	}

	// check if there were any errors during iteration
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	// close the cursor and return the recipes
	err = cursor.Close(ctx)
	if err != nil {
		return nil, err
	}
	return recipes, nil
}

// GetTags returns all tags in the recipe manager
func (m *MongoRecipeManager) GetTags() ([]string, error) {
	// Get the collection handle
	collection := m.client.Database(m.dbName).Collection(m.collectionName)

	// Use a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Find all documents in the collection
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Decode each document into a Recipe object and collect tags as keys in a map
	tags := make(map[string]bool)
	for cursor.Next(ctx) {
		var recipe Recipe
		if err := cursor.Decode(&recipe); err != nil {
			return nil, err
		}

		for _, tag := range recipe.Tags {
			tags[tag] = true
		}
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	// Convert to slice
	var return_tags []string
	for tag := range tags {
		return_tags = append(return_tags, tag)
	}

	return return_tags, nil
}

// SearchRecipes returns all recipes that match the given query string and tags
func (m *MongoRecipeManager) SearchRecipes(query string, tags []string) ([]RecipeSummary, error) {
	// Get the collection handle
	collection := m.client.Database(m.dbName).Collection(m.collectionName)

	// Use a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// create a filter to match recipes with the given tags and the search query
	query_filter := bson.M{"$regex": query, "$options": "i"}
	var tags_filter bson.M
	if len(tags) > 0 {
		// Only search for tags if we've been given a list of tags to search
		tags_filter = bson.M{"tags": bson.M{"$all": tags}}
	} else {
		tags_filter = bson.M{"tags": bson.M{"$exists": true}}
	}
	filter := bson.M{
		"$and": []bson.M{
			tags_filter,
			{"$or": []bson.M{
				{"title": query_filter},
				{"description": query_filter},
				{"comments": bson.M{"$elemMatch": bson.M{"comment": query_filter}}},
			}},
		}}

	// execute the find operation and get the result cursor
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	// iterate over the cursor and decode each document into a RecipeSummary
	var recipe_summaries []RecipeSummary
	for cursor.Next(ctx) {
		var recipe Recipe
		err := cursor.Decode(&recipe)
		if err != nil {
			return nil, err
		}
		recipe_summaries = append(
			recipe_summaries,
			RecipeSummary{recipe.ID.Hex(), recipe.Name, recipe.Tags},
		)
	}

	// check if there were any errors during iteration
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	// close the cursor and return the recipes
	err = cursor.Close(ctx)
	if err != nil {
		return nil, err
	}
	return recipe_summaries, nil
}
