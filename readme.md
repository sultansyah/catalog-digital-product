# Catalog Digital Product

## Prerequisites
- Go
- MySQL
- migrate CLI tool

## Setup Instructions

### 1. Clone Repository
```bash
git clone https://github.com/sultansyah/note-api.git
cd note-api
```

### 2. Environment Variables
Copy `.env.example` to `.env`:
```bash
cp .env.example .env
```

Update the values in `.env`:
```env
JWT_SECRET_KEY=your_jwt_secret
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=your_db_name
DB_HOST=localhost
DB_PORT=3306
```

### 3. Database Migration
Run MySQL migrations:
```bash
migrate -database "mysql://username:password@tcp(localhost:3306)/database_name" -path database/migrations up
```

### 4. Install Dependencies
```bash
go mod tidy
```

### 5. Run Server
```bash
go run main.go
```