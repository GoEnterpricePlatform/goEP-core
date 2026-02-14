# Go CMS Template

CMS backend template built with **Golang**, designed to be used as a base for scalable and maintainable applications.

This project is designed with scalability and clean architecture in mind, allowing easy integration with web frontends and other clients.

| Imagen 1 | Imagen 2 |
|----------|----------|
| <img width="400" src="https://github.com/user-attachments/assets/1d9986cf-728d-4583-a7c8-b0dfb5509a8e" /> | <img src="https://github.com/user-attachments/assets/69602576-efce-4635-ab1e-ec2398c76b12" width="400"> |
| <img src="https://github.com/user-attachments/assets/0e5f99fd-f523-4e85-83d1-e365d1d5e068" width="400"> |  <img src="https://github.com/user-attachments/assets/920fa934-78c2-401f-9ff5-8d358dde96db" width="400"> |

## Architecture
- **Language:** Go (Golang)
- **Architecture:** Hexagonal Architecture 
- **Database:** MongoDB
- **Storage:** S3-compatible (MinIO / AWS S3)
- **API Style:** REST
- **Frontend:** Can be extended with React or any web frontend
- **Email:** Can use Resend or Gmail to send emails

## Authentication & Authorization

- JWT (Access + Refresh)
- HttpOnly Secure Cookies
- Role-Based Access Control
- Permission-based authorization
- Authorization handled via HTTP middleware.

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/amorindev/go-cms-tmpl
   cd go-cms-tmpl
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
