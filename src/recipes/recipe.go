package recipes

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Define structs for a recipe and its ingredients

type Recipe struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name"`
	Description string             `bson:"description"`
	Ingredients []Ingredient       `bson:"ingredients"`
	Steps       []string           `bson:"steps"`
	Tags        []string           `bson:"tags"`
	Comments    []Comments         `bson:"comments"`
	Author      string             `bson:"author"`
}

type Ingredient struct {
	Name     string  `bson:"name"`
	Quantity float32 `bson:"quantity"`
	Unit     string  `bson:"unit"`
}

type Comments struct {
	Comment string    `bson:"comment"`
	Author  string    `bson:"author"`
	Date    time.Time `bson:"date"`
}

// Define functions for operating on recipes

// SetID sets the ID of a recipe
func (recipe *Recipe) SetID(id string) error {
	new_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	recipe.ID = new_id
	return nil
}

// AddTagToRecipe adds a tag to a recipe if it is not already present
func (recipe *Recipe) AddTagToRecipe(tag string) {
	// Only add the tag if it is not in the list already
	for _, t := range recipe.Tags {
		if t == tag {
			return
		}
	}

	// OK to add the tag
	recipe.Tags = append(recipe.Tags, tag)
}

// RemoveTagFromRecipe removes a tag from a recipe if it is present
func (recipe *Recipe) RemoveTagFromRecipe(tag string) {
	for i, t := range recipe.Tags {
		if t == tag {
			recipe.Tags = append(recipe.Tags[:i], recipe.Tags[i+1:]...)
			break
		}
	}
}
