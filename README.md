# Forum Project

This project is a web-based forum application where users can create, view, and participate in discussions on various topics.

## Features

- User authentication and authorization
- Create, edit, and delete posts
- Comment on posts
- Upvote and downvote posts and comments
- Search functionality
- User profiles

## Technologies Used

- Frontend: JavaScript, HTML, CSS
- Backend: Golang
- Database: SQL
- Third Party Dependencies: [go-sqlite3](https://github.com/mattn/go-sqlite3), [bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt)

## Setup Instructions

### Running Directly with Go
1. Clone the repository:
    ```sh
    git clone https://learn.zone01kisumu.ke/git/quochieng/forum.git
    ```
2. Navigate to the project directory:
    ```sh
    cd forum
    ```
3. Run the project:
    ```sh
    go run cmd/main.go
    ```

### Running with Docker

You can also run the project using Docker. Two options are available:

#### 1. Using the Provided Script
The run-forum.sh script builds the Docker image and spins up the container for you.
1. Make sure the script is executable:
    ```sh
    chmod +x run-forum.sh
    ```
2. Run the script:
    ```sh
    ./run-forum.sh
    ```

#### 2. Manually with Docker Commands
1. Build the Docker image:
    ```sh
    docker build -t forum-app .
    ```
2. Run the Docker container (ensure no container named `forum-app-container` is already running; otherwise, remove it with `docker rm -f forum-app-container`):
    ```sh
    docker run -p 8080:8080 --name forum-app-container forum-app
    ```

## Figma Design

For the Figma file presentation, check this [link](https://www.figma.com/design/f8NcgzW17Vox3M1i7v2MCg/forum-wireframe?node-id=35-0&p=f&t=9S7FLtb4xRrrIeaa-0).

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request.

## License

This project is licensed under the MIT License. See the [LICENSE](./LICENSE) file for details.