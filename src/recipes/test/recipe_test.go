package recipes_test

import (
	"testing"

	"github.com/dawsonc/recipes/src/recipes"
)

func TestAddTagToRecipe(t *testing.T) {
	// Make a copy of the test recipe
	recipe := testRecipe1

	// Add a tag to the recipe
	new_tag := "new_tag"
	recipes.AddTagToRecipe(&recipe, new_tag)

	// Make sure the new tag was added
	found := false
	for _, tag := range recipe.Tags {
		if tag == new_tag {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Tag %s not found in recipe", new_tag)
	}

	// Add the tag again (nothing should change)
	recipes.AddTagToRecipe(&recipe, new_tag)

	// Make sure the new tag was not added twice
	count := 0
	for _, tag := range recipe.Tags {
		if tag == new_tag {
			count++
		}
	}
}

// TestRemoveTagFromRecipe tests the RemoveTagFromRecipe function
func TestRemoveTagFromRecipe(t *testing.T) {
	// Make a copy of the test recipe
	recipe := testRecipe1

	// Remove a tag from the recipe
	tag_to_remove := "Test Tag 1"
	recipes.RemoveTagFromRecipe(&recipe, tag_to_remove)

	// Make sure the tag was removed
	for _, tag := range recipe.Tags {
		if tag == tag_to_remove {
			t.Errorf("Tag %s was not removed from recipe %v", tag_to_remove, recipe)
		}
	}

	// Remove the tag again (nothing should change)
	recipes.RemoveTagFromRecipe(&recipe, tag_to_remove)
}

// TestMergeIngredientQuantities tests the MergeIngredientQuantities function
func TestMergeIngredientQuantities(t *testing.T) {
	// Make some test ingredients
	same_ingredients := []recipes.Ingredient{
		{
			Name:     "Test Ingredient",
			Quantity: "1/2 cup",
		},
		{
			Name:     "Test Ingredient",
			Quantity: "1/4 cup",
		},
	}

	different_ingredients := []recipes.Ingredient{
		{
			Name:     "Test Ingredient",
			Quantity: "1/2 cup",
		},
		{
			Name:     "Test Ingredient 2",
			Quantity: "1/4 cup",
		},
	}

	// Merging the same ingredients should be successful
	merged, err := recipes.MergeIngredientQuantities(same_ingredients)
	if err != nil {
		t.Errorf("Error merging ingredients: %v", err)
	}

	// Make sure the ingredients were merged correctly
	expected := "1/2 cup, 1/4 cup"
	if merged.Quantity != expected {
		t.Errorf("Expected merged quantity %s, got %s", expected, merged.Quantity)
	}

	// Merging different ingredients should fail
	_, err = recipes.MergeIngredientQuantities(different_ingredients)
	if err == nil {
		t.Errorf("Merging different ingredients should fail")
	}
}

// TestMergeIngredients tests the MergeIngredients function
func TestMergeIngredients(t *testing.T) {
	// Make some test ingredients
	ingredients := []recipes.Ingredient{
		{
			Name:     "Test Ingredient",
			Quantity: "1/2 cup",
		},
		{
			Name:     "Test Ingredient",
			Quantity: "1/4 cup",
		},
		{
			Name:     "Test Ingredient 2",
			Quantity: "1/4 cup",
		},
	}

	// Merge the ingredients
	merged_ingredients := recipes.MergeIngredients(ingredients)

	// Make sure the ingredients were merged correctly
	if len(merged_ingredients) != 2 {
		t.Errorf("Expected 2 merged ingredients, got %d", len(merged_ingredients))
	}

	// The merged ingredients should be have the correct names and quantities
	expected := map[string]string{
		"Test Ingredient":   "1/2 cup, 1/4 cup",
		"Test Ingredient 2": "1/4 cup",
	}
	for _, ingredient := range merged_ingredients {
		if ingredient.Quantity != expected[ingredient.Name] {
			t.Errorf("Expected quantity %s for ingredient %s, got %s", expected[ingredient.Name], ingredient.Name, ingredient.Quantity)
		}
	}
}
