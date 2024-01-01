# Project README

This README provides an overview of the structure and configuration of the "API Fiber" project. The project is implemented in Go using the Fiber web framework and involves various components for handling requests, managing database connections, and securing endpoints with JSON Web Tokens (JWT).

## Project Structure

### Files
   ```bash
   .
   ├── API-FIBER-DEMO.postman_collection.json
   ├── Makefile
   ├── README.me
   ├── config
   │   ├── config.go
   │   ├── debugAPI.go
   │   └── fiberConfig.go
   ├── database
   │   ├── connect.go
   │   └── createTables.go
   ├── docker-compose.yml
   ├── go.mod
   ├── go.sum
   ├── handlers
   │   ├── handler.hello.go
   │   └── handler.login.go
   ├── main.go
   ├── middleware
   │   └── middleware.go
   ├── models
   │   └── model.login.go
   ├── router
   │   ├── route.notFound.go
   │   ├── router.hello.go
   │   ├── router.login.go
   │   └── router.middleware.go
   └── utils
      └── utils.go
   ```

### Configuration

- **.env**: Contains environment variables for configuration, including server settings, JWT expiration, and MySQL database connection details.

### Usage

1. Set up your `.env` file with appropriate configurations.

2. Run the following command to start the API:
   ```bash
   go run main.go```


## Dependencies

- **[github.com/gofiber/fiber/v2](https://github.com/gofiber/fiber)**: Fiber web framework for handling HTTP requests.
  
- **[github.com/golang-jwt/jwt](https://github.com/golang-jwt/jwt)**: Package for handling JSON Web Tokens (JWT).

- **[github.com/joho/godotenv](https://github.com/joho/godotenv)**: Package for loading environment variables from a file.

- **[github.com/go-sql-driver/mysql](https://github.com/go-sql-driver/mysql)**: MySQL driver for Go's `database/sql` package.


## How to Run

Follow these steps to set up and run the API:

1. **Install Go:**
   Make sure Go is installed on your machine. If not, you can download and install it from [https://golang.org/dl/](https://golang.org/dl/).

2. **Clone the Repository:**
   ```bash
   git clone https://github.com/lordbasex/api-fiber.git
   cd api-fiber```

3. **Create a `.env` file based on the provided `.env` template and update the configurations.**

4. **Run the application:**
   ```bash
   go run main.go```


## Contributing

If you find issues or have suggestions, please feel free to open an issue or create a pull request. Your contributions are welcome!

## License

This project is licensed under the MIT License - see the LICENSE file for details.

