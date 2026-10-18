# GoEP core

GoEp Core is a modular template for building applications in Go.

Designed with clean architecture principles, it allows teams to extend, branch and evolve the system fully customized with long-term maintainability in mind.

<img width="400" height="204" alt="image" src="https://github.com/user-attachments/assets/7418bed6-0bd3-402c-ab4d-6c0ef5b01977" />
<img width="400" height="204" alt="image" src="https://github.com/user-attachments/assets/14705738-f83e-47bb-9b80-a1a9e21f170a" />
<img width="400" height="204" alt="image" src="https://github.com/user-attachments/assets/8cffade1-ed25-48d7-af27-164de220f5e5" />
<img width="400" height="204" alt="image" src="https://github.com/user-attachments/assets/e3e54b89-72c5-428f-a44a-be50d878679a" />

## Architecture

- **Language:** Go (Golang)
- **Architecture:** Hexagonal (Ports & Adapters)
- **Database:** MongoDB
- **Storage:** S3-compatible (MinIO / AWS S3)
- **API:** REST
- **Email:** Resend or Gmail (SMTP)
- **Frontend:** Templ + Datastar + TypeScript (extensible to React separately)

See [the backend module guide](docs/agregar-modulo.md) for a detailed example of adding a module, its manifest, factory, dependencies, and tests.


## Authentication & Authorization

- JWT (Access + Refresh)
- HttpOnly Secure Cookies
- Role-Based Access Control
- Permission-based authorization
- Authorization handled via HTTP middleware.
- Permission-based admin panel

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

6. Run the project, visit http://localhost:8000.

   ```bash
   go tool task run
   ```

   Notes:
   This task runs in watch mode — it stays attentive to changes in Templ, Tailwind, and TypeScript, recompiling and reloading automatically as you      edit. Just save your files and let it do the rest.

---

### If you want to use MinIO as file storage

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
   FILE_STORAGE_PROVIDER=minio
   MINIO_ENDPOINT=localhost:9000
   MINIO_SECURE=false

   MINIO_ACCESS_KEY=appuser
   MINIO_SECRET_KEY=appusersecret
   ```

### Recommended VS Code extensions
- [https://marketplace.visualstudio.com/items?itemName=starfederation.datastar-vscode](Datastar)
- [https://marketplace.visualstudio.com/items?itemName=a-h.templ](Templ)
