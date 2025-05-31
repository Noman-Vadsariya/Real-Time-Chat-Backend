# Real-Time Chat Backend

A real-time chat backend system built with Go, featuring WebSocket support, JWT authentication, and PostgreSQL database integration.

## Features

- User Authentication (JWT-based)
- Real-time messaging using WebSocket
- Message persistence in PostgreSQL
- Online/offline user status
- Private messaging between users
- RESTful API endpoints

## Prerequisites

- Go 1.21 or higher
- PostgreSQL
- Make (optional, for using Makefile commands)

## Environment Variables

Create a `.env` file in the root directory with the following variables:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=chat_db
JWT_SECRET=your-secret-key-here
PORT=8080
```

## Installation

1. Clone the repository:

```bash
git clone <repository-url>
cd chat-backend
```

2. Install dependencies:

```bash
go mod download
```

3. Create the database:

```sql
CREATE DATABASE chat_db;
```

4. Run the application:

```bash
go run .
```

## API Endpoints

### Authentication

- `POST /register` - Register a new user

  ```json
  {
    "username": "user1",
    "password": "password123",
    "email": "user1@example.com"
  }
  ```

- `POST /login` - Login user
  ```json
  {
    "username": "user1",
    "password": "password123"
  }
  ```

### WebSocket

- `GET /ws` - WebSocket connection endpoint
  - Requires Authorization header with JWT token
  - Messages are sent and received in JSON format

### REST API

- `GET /messages/:userId` - Get chat history with a specific user

  - Requires Authorization header with JWT token
  - Returns last 50 messages

- `GET /users` - Get list of all users with online status
  - Requires Authorization header with JWT token

## WebSocket Message Format

### Sending a Message

```json
{
  "type": "message",
  "content": "Hello!",
  "data": {
    "receiver_id": 2
  }
}
```

### User Status Update

```json
{
  "type": "user_status",
  "data": {
    "user_id": 1,
    "online": true
  }
}
```

## Security

- Passwords are hashed using bcrypt
- JWT tokens are used for authentication
- WebSocket connections are protected with JWT authentication
- CORS is enabled for development (should be configured for production)

## Error Handling

The API returns appropriate HTTP status codes and error messages in JSON format:

```json
{
  "error": "Error message description"
}
```

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request
