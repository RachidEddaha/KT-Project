# 🎬 KT Backend API (Go)

A backend API for managing films and user favorites in an online streaming platform.

## 🚀 Features

- User authentication via JWT
- CRUD operations for films
- Role-based access control (only creators can modify/delete)
- Secure password handling
- Database support (MySQL, PostgreSQL, MongoDB, etc.)
- Environment variable configuration
- Optional Docker support

## 📂 Project Structure

```
📂 project-root
├── 📂 cmd                  # Application entry point
│   ├── config-example.env  # Example environment config
│   ├── config.env          # Actual environment config
│   └── main.go             # Main application file
├── 📂 internal             # Business logic & API controllers
│   ├── 📂 controllers      # HTTP request handlers
│   │   ├── authentication
│   │   └── films
│   ├── 📂 dto              # Data Transfer Objects (DTOs)
│   │   ├── authentication.go
│   │   └── films.go
│   ├── 📂 models           # Data models
│   ├── 📂 repositories     
│   ├── 📂 services         # Business logic
├── 📂 pkg                  # Utility packages
│   ├── 📂 configuration    # Configuration management
│   ├── 📂 customerror      # Custom error handling
│   ├── 📂 database         # Database connection setup
│   ├── 📂 logger           # Logging utilities
│   ├── 📂 middlewares      # Middleware functions
│   ├── 📂 utils            # Utility functions
│   └── 📂 webutils         # Web request utilities
├── 📂 scripts              # 
│   ├── 📂 docker           # Docker setup
│   ├── 📂 sql              # script init
├── go.mod                  # Go module dependencies
├── go.sum                  # Go dependency checksums
└── README.md               # Project documentation
```

## 🛠️ Setup & Installation

### 1️⃣ Clone the Repository
```sh
git clone https://github.com/RachidEddaha/KT-Project
cd your-repo
```

### 2️⃣ Set Up Environment Variables
Copy the example config and modify it as needed:
```sh
cp cmd/config-example.env cmd/config.env
```
Update database credentials, JWT secrets, and other settings.

### 3️⃣ Install Dependencies
```sh
go mod tidy
```

### 4️⃣ Run Docker compose
```sh
cd scripts/docker
docker-compose up
```

### 5️⃣ Start the Server
```sh
go run cmd/main.go
```

## 📌 API Endpoints

| Method | Endpoint          | Description                     | Auth Required |
|--------|------------------|---------------------------------|--------------|
| POST   | `/register`      | Register a new user            | ❌ |
| POST   | `/login`         | Login and get JWT token        | ❌ |
| POST   | `/films`         | Create a film                  | ✅ |
| GET    | `/films`         | Get list of films              | ✅ |
| GET    | `/films/:id`     | Get film details               | ✅ |
| PUT    | `/films/:id`     | Update a film (creator only)   | ✅ |
| DELETE | `/films/:id`     | Delete a film (creator only)   | ✅ |

## ✅ Testing
Run tests using:
```sh
go test ./...
```