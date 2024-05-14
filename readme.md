# API-Tinder
### Prerequisites

- You need to install [Golang 1.19 or higher](https://golang.org/)
- You need to install [MongoDB](https://www.mongodb.com/) in your local machine
- You need to install [Visual Studio Code](https://code.visualstudio.com/) as IDE
- You need to install [Git](https://git-scm.com/) in your local machine
- You need to install [Postman](https://www.postman.com/) in your local machine

### Installation
1. Clone the repo
   ```sh
   git clone https://github.com/roby-aw/TechTestAppTinder.git
   ```
2. Go to project directory
   ```sh
   cd TechTestAppTinder
   ```
3. Copy .env.example to .env
   ```sh
    cp .env.example .env
   ```
4. Install dependencies
   ```sh
   go mod tidy
   ```
5. Run the project
   ```sh
    go run main.go
   ```
6. Open Postman and import [API Documentation](https://is3.cloudhost.id/projectvm/PostManAPITinder.json) to Postman
7. You can use the API

## Tech Stack
- [Golang](https://golang.org/)
- [MongoDB](https://www.mongodb.com/)
- [Docker](https://www.docker.com/)
- [JWT](https://jwt.io/)
- [Fiber](https://gofiber.io/)
- [Validator](https://github.com/go-playground/validator)
- [validator-v10](https://github.com/go-playground/validator)
- [godotenv](https://github.com/joho/godotenv)
- [google/uuid](https://github.com/google/uuid)
- [mongo-driver](https://github.com/mongodb/mongo-go-driver)
- [primitive](https://github.com/mongodb/mongo-go-driver/blob/master/bson/primitive/primitive.go)

## Deployment
This project is deployed on [Server Private](http://103.139.193.139:3080)

## API Documentation
You can download the API documentation Postman in [here](https://is3.cloudhost.id/projectvm/PostManAPITinder.json)