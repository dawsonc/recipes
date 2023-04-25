package recipes

// Define an interface for a generic recipe manager
type RecipeManager interface {
	// AddRecipe adds a recipe to the recipe manager and returns the ID of the new recipe
	AddRecipe(recipe Recipe) (string, error)
	// DeleteRecipe deletes a recipe from the recipe manager
	DeleteRecipe(id string) error
	// UpdateRecipe updates a recipe in the recipe manager
	UpdateRecipe(recipe Recipe) error
	// GetAllRecipes returns all recipes in the recipe manager
	GetAllRecipes() ([]Recipe, error)
	// GetRecipeByID returns a recipe with the given ID
	GetRecipeByID(id string) (Recipe, error)
	// GetRecipesByTags returns all recipes with the given tags
	GetRecipesByTags(tags []string) ([]Recipe, error)
	// GetTags returns all tags in the recipe manager
	GetTags() ([]string, error)
	// GetAuthors returns all authors in the recipe manager
	GetAuthors() ([]string, error)
	// SearchRecipes returns all recipes that match the given query string and tags
	SearchRecipes(query string, tags []string, authors []string) ([]RecipeSummary, error)
}

// Define a struct for summarizing a recipe
type RecipeSummary struct {
	ID     string
	Name   string
	Author string
	Tags   []string
}
