# Contributing

Welcome! It's awesome that you are interested in contributing to our project! This file provides some tips for getting started and outlines the expectations for contributions.

## Code Structure

The code is structured into three directories:

- `frontend`: contains the HTML and JS files for the frontend. Files in this directory are served statically. One day we will have a legit build system that generates these static files (transpiling the JSX), but for now we just transpile the JSX on the client side.
- `app`: contains the Go files for the backend server (all as part of the `main` package).
    - `recipe_api.go` defines the endpoints for a REST API for managing recipes (see the `recipes.postman_collection.json` file for an example of using these APIs).
    - `server.go` launches a server that both serves the static frontend files and the REST APIs defined in `recipe_api.go`.
- `src`: contains the packages used by the backend (these are loaded by the `main` package).
    - `recipes`: a package for managing recipes, including relevant data types, an interface for a recipe manager (for loading/editing/searching/etc. recipes), and an implementation of that interface using MongoDB.
        - `tests`: contains the `recipes_test` package used to unit test the recipes interface.


## GitHub

We use [GitHub issues](https://github.com/dawsonc/recipes/issues) to track project development, including new features and bugs to fix. Submitting an issue or contributing a fix for an open issue is a great way to get started!

Once you have completed your feature/fix/etc., please submit your contributions via a [pull request](https://github.com/dawsonc/recipes/pulls). After a (hopefully constructive and helpful) code review, your change can get merged into the main code! To help code review go smoothly, please follow some of the guidelines below.

## Contributor guidelines

Here are some things to keep in mind while working on this project.

- Any new features or code paths should be covered by unit tests. Any changes to existing code should come with accompanying updates to the unit tests (either to catch a previously-uncaught bug or to test the new behavior).
- Code should prioritize readability and be well-commented.
- TODO@dawsonc style guide?
- Code should be auto-formatted.
    - Go code should be formatted with `gofmt`.
    - HTML and JS should be formatted with `TODO@dawsonc`.
- Update the documentation sooner rather than later.

These points mirror the checklist in the PR template, so addressing them early and keeping them in mind will help set your code up for a successfull code review!
