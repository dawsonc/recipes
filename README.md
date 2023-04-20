# Recipes

![Build Status](https://github.com/dawsonc/recipes/actions/workflows/go_build.yml/badge.svg?branch=main)
![Code Health](https://github.com/dawsonc/recipes/actions/workflows/go_fmt.yml/badge.svg?branch=main)

This is a web application for managing recipes. It includes a backend API written in Go and a frontend UI written in React.

## Getting started

To get started, make sure that you have installed [Go](https://go.dev/doc/install) and [MongoDB](https://www.mongodb.com/docs/manual/administration/install-community/). Make sure MongoDB is running as a local service before proceeding.

1. Clone the repository into your `GOPATH`:

```bash
mkdir $GOPATH/src/github.com/dawsonc
cd $GOPATH/src/github.com
git clone https://github.com/dawsonc/recipes.git
```

2. Install dependencies with `go mod download`
3. Run the app with `go run app/*`
4. Visit `localhost:8080` to try it out!

You can also run the unit tests with `go test src/recipes/test/*`

## Technologies Used

- Go
- React
- Bootstrap
- MongoDB

## Contributing

Contributions are welcome! Please see the [CONTRIBUTING.md](CONTRIBUTING.md) file for more information.

## License

This project is licensed under the [MIT License](LICENSE).
