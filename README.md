# Golang Backend API

A RESTful API built with Go, Gin, and MySQL.

## Features

- User management (CRUD operations)
- Database health monitoring
- WebSocket support
- CORS enabled for frontend integration

## API Endpoints

### User Management

#### Create User
- **POST** `/api/users/`
- **Body:**
```json
{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "password123"
}
```
- **Response:** `201 Created`
```json
{
  "message": "User created successfully",
  "user": {
    "id": 1,
    "username": "john_doe",
    "email": "john@example.com",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

#### Get All Users
- **GET** `/api/users/`
- **Response:** `200 OK`
```json
{
  "users": [
    {
      "id": 1,
      "username": "john_doe",
      "email": "john@example.com",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

#### Get User by ID
- **GET** `/api/users/{id}`
- **Response:** `200 OK`
```json
{
  "user": {
    "id": 1,
    "username": "john_doe",
    "email": "john@example.com",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

#### Update User
- **PUT** `/api/users/{id}`
- **Body:**
```json
{
  "username": "john_updated",
  "email": "john_updated@example.com"
}
```
- **Response:** `200 OK`
```json
{
  "message": "User updated successfully",
  "user": {
    "id": 1,
    "username": "john_updated",
    "email": "john_updated@example.com",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

#### Update Password
- **PATCH** `/api/users/{id}/password`
- **Body:**
```json
{
  "password": "newpassword123"
}
```
- **Response:** `200 OK`
```json
{
  "message": "Password updated successfully"
}
```

#### Delete User
- **DELETE** `/api/users/{id}`
- **Response:** `200 OK`
```json
{
  "message": "User deleted successfully"
}
```

### Health Check
- **GET** `/health`
- **Response:** `200 OK`
```json
{
  "status": "up",
  "message": "It's healthy",
  "open_connections": "1",
  "in_use": "0",
  "idle": "1"
}
```

### WebSocket
- **GET** `/websocket`
- Establishes a WebSocket connection that sends timestamps every 2 seconds

## Environment Variables

Create a `.env` file with the following variables:

```env
PORT=8080
MYSQL_DB_HOST=localhost
MYSQL_DB_PORT=3306
MYSQL_DB_USERNAME=user
MYSQL_DB_PASSWORD=password
MYSQL_DB_DATABASE=database
```

## Running the Application

1. **Install dependencies:**
```bash
go mod download
```

2. **Run the application:**
```bash
go run cmd/api/main.go
```

3. **Run tests:**
```bash
go test ./...
```

## Database Schema

### Users Table
```sql
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_email (email),
    INDEX idx_username (username)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

## Error Handling

The API returns appropriate HTTP status codes and error messages:

- `400 Bad Request`: Invalid request data
- `404 Not Found`: Resource not found
- `409 Conflict`: Resource already exists (e.g., duplicate email/username)
- `500 Internal Server Error`: Server error

## Validation

- Username: Required, unique
- Email: Required, valid email format, unique
- Password: Required, minimum 6 characters
- User ID: Must be a valid integer
