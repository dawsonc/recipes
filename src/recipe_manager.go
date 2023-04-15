package recipes

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Recipe struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name"`
	Description string             `bson:"description"`
	Ingredients []Ingredient       `bson:"ingredients"`
	Steps       []string           `bson:"steps"`
	Tags        []string           `bson:"tags"`
}

type Ingredient struct {
	Name     string `bson:"name"`
	Quantity string `bson:"quantity"`
}

func GetAllRecipes(client *mongo.Client, dbName string, collectionName string) ([]Recipe, error) {
	var recipes []Recipe

	// Get the collection handle
	collection := client.Database(dbName).Collection(collectionName)

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

func GetRecipeByID(client *mongo.Client, dbName string, collectionName string, id string) (Recipe, error) {
	var recipe Recipe

	// Get the collection handle
	collection := client.Database(dbName).Collection(collectionName)

	// Use a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Find the document with the given ID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return recipe, err
	}

	if err := collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&recipe); err != nil {
		return recipe, err
	}

	return recipe, nil
}

func GetRecipesByTag(client *mongo.Client, dbName string, collectionName string, tag string) ([]Recipe, error) {
	// Get the collection handle
	collection := client.Database(dbName).Collection(collectionName)

	// Use a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// create a filter to match recipes with the given tag
	filter := bson.M{"tags": tag}

	// execute the find operation and get the result cursor
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
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

func UpdateRecipe(client *mongo.Client, dbName string, collectionName string, id string, recipe Recipe) error {
	// Get the collection handle
	collection := client.Database(dbName).Collection(collectionName)

	// Use a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Update the recipe
	_, err := collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": recipe})
	if err != nil {
		return err
	}

	return nil
}

func InsertRecipe(client *mongo.Client, dbName string, collectionName string, recipe Recipe) error {
	// Get the collection handle
	collection := client.Database(dbName).Collection(collectionName)

	// Use a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Insert the recipe
	_, err := collection.InsertOne(ctx, recipe)
	if err != nil {
		return err
	}

	return nil
}

func DeleteRecipe(client *mongo.Client, dbName string, collectionName string, id string) error {
	// Get the collection handle
	collection := client.Database(dbName).Collection(collectionName)

	// Use a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Delete the recipe
	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	return nil
}

func ConnectToMongoDB(uri string) (*mongo.Client, error) {
	// Set client options
	clientOptions := options.Client().ApplyURI(uri)

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}
