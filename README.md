# Go CMS Template

CMS backend template built with **Golang**, designed to be used as a base for scalable and maintainable applications.

This project is designed with scalability and clean architecture in mind, allowing easy integration with web frontends and other clients.


## Architecture

- Language: Go (Golang)
- Architecture: Hexagonal Architecture
- Database: MongoDB
- Storage: Object storage (MinIO / S3 compatible)
- API Style: REST
- Environment-based configuration

## Installation

1. Clone the repository:
    ```bash
    git clone https://github.com/amorindev/go-cms-tmpl
    cd go-tmpl
    ```

2. Download dependencies:
    ```bash
    go mod tidy
    ```

3. Set environment variables, add a `.env` file based on `env.example`

4. Run the project:
    ```bash
    make run
    ```
