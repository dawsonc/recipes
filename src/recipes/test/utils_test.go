package recipes_test

import "github.com/dawsonc/recipes/src/recipes"

// Define some test recipes
var testRecipe1 = recipes.Recipe{
	Name:        "Test Recipe",
	Description: "This is a test recipe",
	Ingredients: []recipes.Ingredient{
		{Name: "Test Ingredient 1", Quantity: 1, Unit: "bunch"},
		{Name: "Test Ingredient 2", Quantity: 1, Unit: "can"},
	},
	Steps: []string{"Test Step 1", "Test Step 2"},
	Tags:  []string{"Test Tag 1", "Test Tag 2"},
}
var testRecipe2 = recipes.Recipe{
	Name:        "Test Recipe 2",
	Description: "This is a second test recipe",
	Ingredients: []recipes.Ingredient{
		{Name: "Test Ingredient 1", Quantity: 1, Unit: "bundle"},
		{Name: "Test Ingredient 3", Quantity: 2, Unit: "bits"},
	},
	Steps: []string{"Test Step 1", "Test Step 2"},
	Tags:  []string{"Test Tag 1", "Test Tag 3"},
}
