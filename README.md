# GoPaster

GoPaster is a simple, web-based pastebin application written in Go. It allows users to create, read, and delete pastes securely using SSL. This application is designed for simplicity and ease of use, making it perfect for quickly sharing snippets of text or code.

## Features

- **Create Pastes**: Easily create and share pastes with a unique URL.
- **Read Pastes**: Access your pastes via their unique URL.
- **Delete Pastes**: Remove pastes when they are no longer needed.
- **SSL Support**: Secure communication with SSL (self-signed certificates can be generated automatically).
- **Docker Support**: Easily deployable with Docker and Docker Compose.

## Getting Started

### Local Run

1. **Create SSL certificates if you don't have them:**

    ```sh
    mkdir certs
    openssl req -x509 -newkey rsa:4096 -keyout certs/server.key -out certs/server.crt -days 365 -nodes -subj "/CN=localhost"
    ```

2. **Run the application:**

    ```sh
    go run .
    ```

### Docker

1. **Build Docker image:**

    ```sh
    docker build -t gopaster -f Dockerfile .
    ```

2. **Run Docker image:**

    ```sh
    docker run -d -p 443:443 -v $(pwd)/data:/app/data -v $(pwd)/certs:/app/certs gopaster
    ```

### Docker Compose

1. **Run Docker Compose:**

    ```sh
    docker-compose up --build
    ```

## Contributing

Feel free to submit issues and pull requests. For major changes, please open an issue first to discuss what you would like to change.

## License

This project is licensed under the GPL-3.0 license - see the [LICENSE](LICENSE) file for details.