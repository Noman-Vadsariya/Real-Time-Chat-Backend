# Go Chat - Real-time Chat Application

![Go Chat Logo](https://raw.githubusercontent.com/yourusername/go-chat/main/assets/logo.png)

A modern, real-time chat application built with Go, WebSocket, and a beautiful UI. This application provides instant messaging capabilities with features like user authentication, real-time status updates, and message persistence.

## 🌟 Features

- **Real-time Messaging**: Instant message delivery using WebSocket technology
- **User Authentication**: Secure login and registration system
- **Online Status**: Real-time user online/offline status updates
- **Message History**: Persistent message storage and retrieval
- **Modern UI**: Clean and responsive interface built with Tailwind CSS
- **Cross-platform**: Works on desktop and mobile browsers

## 🚀 Quick Start

### Prerequisites

- Go 1.16 or higher
- SQLite (default) or PostgreSQL
- Modern web browser

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/go-chat.git
   cd go-chat
   ```

2. Install dependencies:

   ```bash
   go mod download
   ```

3. Configure environment variables:

   ```bash
   # Create a .env file
   DB_TYPE=sqlite
   PORT=3000
   JWT_SECRET=your-super-secret-key-123
   ```

4. Run the application:

   ```bash
   go run .
   ```

5. Open your browser and navigate to:
   ```
   http://localhost:3000/client/
   ```

## 📸 Screenshots

### Login Screen

![Login Screen](https://raw.githubusercontent.com/yourusername/go-chat/main/assets/login.png)

### Chat Interface

![Chat Interface](https://raw.githubusercontent.com/yourusername/go-chat/main/assets/chat.png)

### User List

![User List](https://raw.githubusercontent.com/yourusername/go-chat/main/assets/users.png)

## 🔧 Configuration

### Environment Variables

| Variable     | Description                     | Default |
| ------------ | ------------------------------- | ------- |
| `DB_TYPE`    | Database type (sqlite/postgres) | sqlite  |
| `PORT`       | Server port                     | 3000    |
| `JWT_SECRET` | Secret key for JWT tokens       | -       |

### Database Configuration

The application supports both SQLite and PostgreSQL:

```bash
# For SQLite (default)
DB_TYPE=sqlite

# For PostgreSQL
DB_TYPE=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=chat_db
```

## 🛠️ API Endpoints

### Authentication

- `POST /register` - Register a new user
- `POST /login` - Login and get JWT token

### Users

- `GET /users` - Get list of users (requires authentication)
- `GET /users/:id` - Get user details (requires authentication)

### Messages

- `GET /messages/:user_id` - Get chat history with a user (requires authentication)
- WebSocket `/ws` - Real-time messaging endpoint

## 🔒 Security

- JWT-based authentication
- Password hashing using bcrypt
- WebSocket connection validation
- CORS protection
- Input sanitization

## 🎨 UI Features

- Responsive design
- Real-time status indicators
- Message bubbles with timestamps
- User online/offline status
- Clean and modern interface
- Dark mode support (coming soon)

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [Gorilla WebSocket](https://github.com/gorilla/websocket)
- [Tailwind CSS](https://tailwindcss.com)
- [GORM](https://gorm.io)

## 📞 Support

For support, email support@gochat.com or open an issue in the GitHub repository.

---

Made with ❤️ by [Your Name]
