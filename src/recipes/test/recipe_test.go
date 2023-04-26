package recipes_test

import (
	"testing"
)

func TestAddTagToRecipe(t *testing.T) {
	// Make a copy of the test recipe
	recipe := testRecipe1

	// Add a tag to the recipe
	new_tag := "new_tag"
	recipe.AddTagToRecipe(new_tag)

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
	recipe.AddTagToRecipe(new_tag)

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
	recipe.RemoveTagFromRecipe(tag_to_remove)

	// Make sure the tag was removed
	for _, tag := range recipe.Tags {
		if tag == tag_to_remove {
			t.Errorf("Tag %s was not removed from recipe %v", tag_to_remove, recipe)
		}
	}

	// Remove the tag again (nothing should change)
	recipe.RemoveTagFromRecipe(tag_to_remove)
}
