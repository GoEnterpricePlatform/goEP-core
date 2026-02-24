# GoEP core

GoEp Core is a modular enterprise system template built with Golang.

It provides a scalable foundation composed of independent business modules that a modern enterprise platform requires — such as identity, authorization, content & Information Management.

Designed with clean architecture principles, it enables teams to extend, fork, and evolve the system into fully customized enterprise solutions with long-term maintainability in mind.

| Imagen 1                                                                                                  | Imagen 2                                                                                                |
| --------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------- |
| <img width="400" src="https://github.com/user-attachments/assets/1d9986cf-728d-4583-a7c8-b0dfb5509a8e" /> | <img src="https://github.com/user-attachments/assets/69602576-efce-4635-ab1e-ec2398c76b12" width="400"> |
| <img src="https://github.com/user-attachments/assets/0e5f99fd-f523-4e85-83d1-e365d1d5e068" width="400">   | <img src="https://github.com/user-attachments/assets/920fa934-78c2-401f-9ff5-8d358dde96db" width="400"> |

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
   git clone https://github.com/GoEnterpricePlatform/goEP-core
   cd goEP-core
   ```

2. Download dependencies:

   ```bash
   go mod tidy
   ```

3. Get Gmail Credentials
   - Go to your [Google Account](https://myaccount.google.com/)
   - Navigate to **Security**
   - Enable **2-Step Verification** (required to create App Passwords)

   > You can check the official guide here:  
   > [Google Help Center](https://support.google.com/mail/answer/185833)
   - Go to [Create and manage app passwords](https://myaccount.google.com/apppasswords)
     - Or access it from **Security → Signing in to Google → App passwords** (choose the method you prefer)

   - Enter a name for your app
   - Click **Create**
   - Google will generate a password
   - Copy the password and remove the spaces
   - Use it as `GMAIL_PASS`

   ```env
   GMAIL_USERNAME=my-email@gmail.com
   GMAIL_PASS=xxxxxxxx
   ```

4. Set environment variables, add a `.env` file based on `env.example`, MINIO_ACCESS_KEY and MINIO_SECRET_KEY will be assigned in the following steps.
5. Starts the development environment using Docker Compose.

   ```bash
   make compose-dev
   ```

   Notes:
   - On Windows, Docker Desktop must be open.
   - Skips this step if the containers are already running.

6. MinIO Configuration

   Newer versions of MinIO don’t allow creating credentials via the UI, so we’ll use the (MinIO Client)[https://github.com/minio/mc]

   ```bash
   go install github.com/minio/mc@latest
   ```

   Establish a connection to the MinIO using the MINIO_ROOT_USER and MINIO_ROOT_PASSWORD from the .env file:

   ```bash
   mc alias set local http://localhost:9000 user "123456()Secret"
   ```

   Creates a new user

   ```bash
   mc admin user add local appuser appusersecret
   ```

   Grant read/write permissions on all buckets to the user (appuser) for simplicity.

   ```bash
   mc admin policy attach local readwrite --user=appuser
   ```

   Set the new MinIO user credentials in the .env file:

   ```
   MINIO_ACCESS_KEY=appuser
   MINIO_SECRET_KEY=appusersecret
   ```

7. Run the project:
   ```bash
   make run
   ```
