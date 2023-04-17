package recipes

import (
	"fmt"
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
}

type Ingredient struct {
	Name     string `bson:"name"`
	Quantity string `bson:"quantity"`
}

type Comments struct {
	Comment string    `bson:"comment"`
	Author  string    `bson:"author"`
	Date    time.Time `bson:"date"`
}

// Define functions for operating on recipes

// AddTagToRecipe adds a tag to a recipe if it is not already present
func AddTagToRecipe(recipe *Recipe, tag string) {
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
func RemoveTagFromRecipe(recipe *Recipe, tag string) {
	for i, t := range recipe.Tags {
		if t == tag {
			recipe.Tags = append(recipe.Tags[:i], recipe.Tags[i+1:]...)
			break
		}
	}
}

// MergeIngredientQuantities merges multiple ingredients of the same type together
func MergeIngredientQuantities(ingredients []Ingredient) (Ingredient, error) {
	// Only merge ingredients of the same type
	if len(ingredients) == 0 {
		return Ingredient{}, fmt.Errorf("cannot merge zero ingredients")
	}
	for i := 1; i < len(ingredients); i++ {
		if ingredients[i].Name != ingredients[0].Name {
			return Ingredient{}, fmt.Errorf(
				"cannot merge ingredients of different types: "+
					"%s and %s", ingredients[i].Name, ingredients[0].Name)
		}
	}

	// Merge the ingredients
	merged := ingredients[0]
	for i := 1; i < len(ingredients); i++ {
		merged.Quantity += ", " + ingredients[i].Quantity
	}

	return merged, nil
}

// MergeIngredients takes a list of ingredients and merges any that are the same
func MergeIngredients(ingredients []Ingredient) []Ingredient {
	// Make a map of ingredient names to lists of ingredients
	ingredient_map := make(map[string][]Ingredient)

	// Fill that map to track which ingredients are the same
	for _, ingredient := range ingredients {
		ingredient_map[ingredient.Name] = append(ingredient_map[ingredient.Name], ingredient)
	}

	// Merge all ingredients with the same name
	merged_ingredients := make([]Ingredient, 0)
	for _, ingredient_list := range ingredient_map {
		merged, err := MergeIngredientQuantities(ingredient_list)
		if err != nil {
			// This should never happen, since all ingredients in the list
			// should have the same name
			panic(err)
		}
		merged_ingredients = append(merged_ingredients, merged)
	}

	return merged_ingredients
}
