# Go Gin Template

A production-ready RESTful API boilerplate built with Go, Gin, GORM, and Clean Architecture principles.

## 🏗 Architecture

This project follows **Clean Architecture** to ensure separation of concerns, testability, and maintainability:

- **`cmd/`**: Entry points for the application (Server).
- **`internal/delivery/http/`**: Transport layer (Gin Handlers & Routing).
- **`internal/service/`**: Business logic layer.
- **`internal/repository/`**: Data access layer (GORM).
- **`internal/entity/`**: Domain models.
- **`internal/dto/`**: Data Transfer Objects for requests and responses.
- **`internal/middleware/`**: Custom Gin middlewares (JWT, Logging, CORS).
- **`internal/config/`**: Configuration management using Viper and `.env`.
- **`pkg/`**: Internal libraries and utilities (Logger, Standard Response).

## 🚀 Features

- **Authentication**: Secure JWT-based auth with Access & Refresh Token support.
- **Role-Based Access Control (RBAC)**: Middleware for restricting access based on user roles (`admin` vs `user`).
- **Account Verification**: System to `Activate` or `Deactivate` users (Admin only).
- **Advanced Pagination & Search**: Generic pagination with metadata (`total_pages`, `has_next`, etc.) and search filtering.
- **Production Logging**: Scalable logging using Uber's **Zap**. In production, logs are automatically separated into files by level (`info.log`, `warn.log`, `error.log`).
- **Standardized Response**: Unified JSON response format for success and error states.

## 🛣 API Endpoints

### Authentication
| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `POST` | `/api/v1/auth/register` | Register a new user |
| `POST` | `/api/v1/auth/login` | Login and get JWT tokens |
| `POST` | `/api/v1/auth/refresh` | Refresh access token using refresh token |

### User Profile (Authenticated)
| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `GET` | `/api/v1/me` | Get current user profile |
| `PUT` | `/api/v1/me` | Update personal profile |

### Admin Management (Admin Only)
| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `GET` | `/api/v1/users` | List users (Paginated & Searchable) |
| `PUT` | `/api/v1/users/:id` | Update any user data |
| `DELETE` | `/api/v1/users/:id` | Soft delete a user |
| `PATCH` | `/api/v1/users/:id/activate` | Verify/Activate user account |
| `PATCH` | `/api/v1/users/:id/deactivate` | Deactivate/Unverify user account |

### Health Check
| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `GET` | `/health` | Check if server is running |

## 🛠 Getting Started

1. **Clone the repository**
2. **Setup environment variables**
   ```bash
   cp .env.example .env
   ```
3. **Run the application**
   ```bash
   go run cmd/server/main.go
   ```

## 📝 Logging (Production Mode)
When `APP_ENV=production` is set in `.env`, logs will be stored in:
- `logs/info.log`
- `logs/warn.log`
- `logs/error.log`

## 📝 Lisensi

MIT License

Copyright (c) 2026 M. Darma Putra Ramadhan

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
