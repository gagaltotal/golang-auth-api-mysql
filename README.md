# Golang Auth API - Complete REST API

Production-ready REST API dengan Golang, MySQL, JWT authentication, role-based access control, audit logging, dan product management.

## Features

### Authentication & Authorization
- Register user baru
- Login dengan JWT + refresh token
- Refresh token rotation (automatic revoke old token)
- Logout dengan revoke refresh token
- Password hashing dengan bcrypt
- Role-based access control (Admin & User)
- JWT middleware untuk protected routes
- Rate limiting untuk prevent abuse

### User Management
- Get user profile (`/me`)
- Update profile (name, email)
- Change password (dengan validasi old password)
- Delete user (Admin only)
- List all users dengan pagination (Admin only)
- User tidak bisa akses data user lain

### Product Management
- Create product (Admin only)
- Get product by ID (Public)
- List products dengan pagination & filter (Public)
- Update product (Admin only)
- Delete product (Admin only)
- Upload product image dengan validasi:
  - MIME type checking (jpeg, png, gif, webp)
  - Size limit (max 5MB)
  - Random filename generation (security)
  - Auto delete old image saat update

### Audit Logging
- Auto log semua user & admin actions
- Store ke database MySQL
- Track: user, action, resource, method, path, IP, payload
- Get audit logs dengan filter (Admin only)
- Middleware audit otomatis

### Testing
- Unit tests dengan testify/mock
- 100% mock repository & service
- Test coverage untuk:
  - Auth service (register, login, refresh, logout)
  - User service (profile, update, change password)
  - Product service (CRUD operations)
- Test-driven development ready

### API Documentation
- Swagger YAML documentation
- Semua endpoint terdokumentasi
- Request/response examples
- Security schemes (Bearer JWT)
- Multiple server configs

### Docker Support
- MySQL container
- API container
- Docker Compose untuk easy setup
- Production-ready Dockerfile
- Environment variables

---

## Tech Stack

- **Language**: Go 1.21 - 1.22+
- **Framework**: Gin Web Framework
- **Database**: MySQL 8.0
- **ORM**: GORM
- **Authentication**: JWT (golang-jwt/jwt)
- **Testing**: Testify
- **Documentation**: OpenAPI 3.0 (Swagger)
- **Containerization**: Docker & Docker Compose

---

## Installation

### Prerequisites
- Go 1.21 - 1.22 atau lebih tinggi
- MySQL 8.0
- Docker & Docker Compose (optional)

### 1. Clone Repository
```bash
git clone https://github.com/gagaltotal/golang-auth-api-mysql.git
cd golang-auth-api-mysql
```

### 2. Install Dependencies
```bash
go mod download
```

### 3. Setup Environment
```bash
cp .env.example .env
```

Edit `.env` sesuai konfigurasi:
```env
ENVIRONMENT=development
SERVER_PORT=8080

# Database
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=auth_api

# JWT
JWT_SECRET=your-super-secret-jwt-key-change-in-production
JWT_EXPIRATION=15m
REFRESH_EXPIRATION=7d

# Upload
UPLOAD_PATH=./uploads
```

### 4. Create Database
```bash
mysql -u root -p
```
```sql
CREATE DATABASE auth_api CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 5. Run Application
```bash
go run cmd/api/main.go
```

Server akan berjalan di `http://localhost:8080`

---

## Docker Setup

### 1. Setup Environment
```bash
cp .env.example .env
```

### 2. Run with Docker Compose
```bash
docker-compose -f docker/docker-compose.yml up -d
```

### 3. Check Status
```bash
docker-compose -f docker/docker-compose.yml ps
```

### 4. View Logs
```bash
docker-compose -f docker/docker-compose.yml logs -f api
```

### 5. Stop Containers
```bash
docker-compose -f docker/docker-compose.yml down
```

---

## Running Tests

### Run All Tests
```bash
go test ./tests/... -v
```

### Run Specific Test
```bash
go test ./tests/auth_service_test.go -v
```

### Run with Coverage
```bash
go test ./tests/... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

---

## API Documentation

### Base URL
```
http://localhost:8080/api/v1
```

### Authentication
Semua protected endpoints memerlukan Bearer token di header:
```
Authorization: Bearer <your_access_token>
```

### Endpoints

#### Authentication
| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| POST | `/auth/register` | Register user baru | ❌ |
| POST | `/auth/login` | Login user | ❌ |
| POST | `/auth/refresh` | Refresh access token | ❌ |
| POST | `/auth/logout` | Logout & revoke token | ❌ |

#### Users
| Method | Endpoint | Description | Auth | Role |
|--------|----------|-------------|------|------|
| GET | `/users/me` | Get profile | ✅ | User |
| PUT | `/users/me` | Update profile | ✅ | User |
| PUT | `/users/me/password` | Change password | ✅ | User |
| GET | `/users` | Get all users | ✅ | Admin |
| DELETE | `/users/:id` | Delete user | ✅ | Admin |

#### Products
| Method | Endpoint | Description | Auth | Role |
|--------|----------|-------------|------|------|
| GET | `/products` | Get all products | ❌ | Public |
| GET | `/products/:id` | Get product by ID | ❌ | Public |
| POST | `/products` | Create product | ✅ | Admin |
| PUT | `/products/:id` | Update product | ✅ | Admin |
| DELETE | `/products/:id` | Delete product | ✅ | Admin |
| POST | `/products/:id/image` | Upload image | ✅ | Admin |

#### Audit Logs
| Method | Endpoint | Description | Auth | Role |
|--------|----------|-------------|------|------|
| GET | `/audit-logs` | Get audit logs | ✅ | Admin |

---

## Usage Examples

### 1. Register User
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "password123"
  }'
```

### 2. Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

Response:
```json
{
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "550e8400-e29b-41d4-a716...",
    "user": {
      "id": 1,
      "name": "John Doe",
      "email": "john@example.com",
      "role": "user"
    }
  }
}
```

### 3. Get Profile
```bash
curl -X GET http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### 4. Create Product (Admin)
```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "New Product",
    "description": "Product description",
    "price": 15000,
    "stock": 50
  }'
```

### 5. Upload Product Image (Admin)
```bash
curl -X POST http://localhost:8080/api/v1/products/1/image \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -F "image=@/path/to/image.jpg"
```

### 6. Get Products with Filter
```bash
curl -X GET "http://localhost:8080/api/v1/products?page=1&limit=10&search=product&min_price=10000&max_price=50000"
```

### 7. Refresh Token
```bash
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "550e8400-e29b-41d4-a716..."
  }'
```

---

## Database Schema

### Users Table
```sql
CREATE TABLE users (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  email VARCHAR(100) UNIQUE NOT NULL,
  password VARCHAR(255) NOT NULL,
  role VARCHAR(20) DEFAULT 'user',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL
);
```

### Products Table
```sql
CREATE TABLE products (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(200) NOT NULL,
  description TEXT,
  price DECIMAL(10,2) NOT NULL,
  stock INT DEFAULT 0,
  image_url VARCHAR(500),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL
);
```

### Refresh Tokens Table
```sql
CREATE TABLE refresh_tokens (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  user_id BIGINT UNSIGNED NOT NULL,
  token VARCHAR(500) UNIQUE NOT NULL,
  expires_at TIMESTAMP NOT NULL,
  revoked BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
```

### Audit Logs Table
```sql
CREATE TABLE audit_logs (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  user_id BIGINT UNSIGNED,
  action VARCHAR(50) NOT NULL,
  resource VARCHAR(50),
  resource_id BIGINT UNSIGNED,
  method VARCHAR(10),
  path VARCHAR(255),
  ip_address VARCHAR(50),
  user_agent VARCHAR(500),
  payload TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
);
```

---

## Security Features

1. **Password Hashing**: Bcrypt dengan cost 10
2. **JWT Token**: Short-lived access token (15 minutes)
3. **Refresh Token**: Long-lived (7 days) dengan rotation
4. **Token Revocation**: Blacklist refresh tokens
5. **Rate Limiting**: 60 requests per minute per IP
6. **File Upload Security**:
   - MIME type validation
   - Size limit (5MB)
   - Random filename generation
   - Stored outside webroot
7. **SQL Injection**: Prevented by GORM ORM
8. **Role-Based Access Control**: Middleware untuk admin/user
9. **Audit Logging**: Track semua sensitive operations

---

## Project Structure Details

```
golang-auth-api-mysql/
├── cmd/api/main.go              # Entry point aplikasi
├── config/config.go             # Configuration management
├── database/                    # Database related
│   ├── mysql.go                # MySQL connection
│   ├── migrate.go              # Auto migration
│   └── seed.go                 # Seed data
├── internal/
│   ├── domain/                 # Domain models & DTOs
│   │   ├── user.go             # User model
│   │   ├── product.go          # Product model
│   │   ├── audit.go            # Audit log model
│   │   └── refresh_token.go    # Refresh token model
│   ├── repository/             # Data access layer
│   │   └── *_impl.go           # Repository implementations
│   ├── service/                # Business logic layer
│   │   └── *_impl.go           # Service implementations
│   ├── handler/                # HTTP handlers
│   │   ├── handler.go          # Main handlers
│   │   └── upload.go           # File upload handler
│   ├── middleware/             # HTTP middlewares
│   │   ├── jwt_middleware.go   # JWT authentication
│   │   ├── role_middleware.go  # Role authorization
│   │   ├── rate_limit.go       # Rate limiting
│   │   └── audit_middleware.go # Audit logging
│   └── routes/routes.go        # Route definitions
├── pkg/                        # Reusable packages
│   ├── hash/bcrypt.go          # Password hashing
│   └── jwt/jwt.go              # JWT utilities
├── scripts/                    # Script install swagger
│   ├── generate-swagger.sh     # file generate swagger
│   └── install-swag.sh         # file install swagger
├── tests/                      # Unit tests
│   ├── auth_service_test.go    # Auth tests
│   ├── product_service_test.go # Product tests
│   └── user_service_test.go   # User tests
├── docker/                    # Docker files
│   ├── Dockerfile             # API Dockerfile
│   └── docker-compose.yml     # Compose config
├── docs/                      # Documentation
│   ├── swagger.yaml           # API documentation
│   └── docs.go                # Swagger docs
├── uploads/                   # Upload directory
├── .air.toml                  # Air - Live reload for Go apps
├── .env.example               # Environment template
├── .gitignore                 # Git ignore rules
├── go.mod                     # Go modules
├── Makefile                   # Make file
└── README.md                  # This file
```

---

## Deployment

### Production Checklist

- [ ] Change `JWT_SECRET` ke random strong secret
- [ ] Set `ENVIRONMENT=production`
- [ ] Use strong database password
- [ ] Enable HTTPS/TLS
- [ ] Setup firewall rules
- [ ] Configure proper CORS
- [ ] Setup monitoring & logging
- [ ] Regular database backups
- [ ] Rate limiting tuning
- [ ] File upload size limits

### Deploy ke Cloud

#### 1. **DigitalOcean / AWS / GCP**
- Build Docker image
- Push ke container registry
- Deploy ke VM atau Kubernetes
- Setup load balancer
- Configure SSL certificate

#### 2. **Heroku**
```bash
heroku create your-app-name
heroku addons:create jawsdb:kitefin  # MySQL addon
git push heroku main
```

---

## Contributing

1. Fork repository
2. Create feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to branch (`git push origin feature/AmazingFeature`)
5. Open Pull Request

---

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## Author

**Tot666 - Gagaltotal666**
- Email: gagaltotal666@gtrtech666.my.id
- GitHub: [@gagaltotal](https://github.com/gagaltotal)

---

## Acknowledgments

- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [GORM](https://gorm.io/)
- [JWT-Go](https://github.com/golang-jwt/jwt)
- [Testify](https://github.com/stretchr/testify)

---

## Support

Jika ada pertanyaan atau issue:
- Create GitHub issue
- Email: gagaltotal666@gtrtech666.my.id

---

## Roadmap

- [x] Basic authentication
- [x] Refresh token rotation
- [x] Role-based access control
- [x] Audit logging
- [x] Product management
- [x] Image upload
- [x] Unit tests
- [x] Swagger documentation
- [x] Docker support
- [ ] Email verification
- [ ] Forgot password
- [ ] 2FA authentication
- [ ] WebSocket support
- [ ] GraphQL API
- [ ] Microservices architecture

---

**Happy Coding!**