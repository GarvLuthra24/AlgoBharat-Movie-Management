# Environment Variables Setup

This project now uses environment variables for configuration. Follow the steps below to set up your environment.

## Backend Setup

1. Copy the example environment file:
   ```bash
   cd backend
   cp .env.example .env
   ```

2. Edit the `.env` file with your actual values:
   ```env
   # Database Configuration
   DB_HOST=localhost
   DB_PORT=3306
   DB_USER=root
   DB_PASSWORD=your_mysql_password
   DB_NAME=algoBharat
   DB_DRIVER=mysql

   # Server Configuration
   PORT=8080
   CORS_ORIGIN=http://localhost:5173

   # JWT Configuration
   JWT_SECRET=your-very-secure-secret-key-here
   ```

3. Make sure MySQL is running and create the database:
   ```sql
   CREATE DATABASE algoBharat;
   ```

4. Install dependencies and run:
   ```bash
   go mod tidy
   go run main.go
   ```

## Frontend Setup

1. Copy the example environment file:
   ```bash
   cd frontend
   cp .env.example .env
   ```

2. Edit the `.env` file with your backend URL:
   ```env
   # Backend API Configuration
   VITE_API_BASE_URL=http://localhost:8080
   ```

3. Install dependencies and run:
   ```bash
   npm install
   npm run dev
   ```

## Environment Variables Reference

### Backend (.env)

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `DB_HOST` | Database host | localhost | No |
| `DB_PORT` | Database port | 3306 | No |
| `DB_USER` | Database username | root | No |
| `DB_PASSWORD` | Database password | password | No |
| `DB_NAME` | Database name | algoBharat | No |
| `DB_DRIVER` | Database driver (mysql/sqlite3) | sqlite3 | No |
| `PORT` | Server port | 8080 | No |
| `CORS_ORIGIN` | CORS allowed origin | http://localhost:5173 | No |
| `JWT_SECRET` | JWT signing secret | my_secret_key | No |

### Frontend (.env)

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `VITE_API_BASE_URL` | Backend API base URL | http://localhost:8080 | No |

## Database Migration

The application supports both SQLite (default) and MySQL. To switch to MySQL:

1. Set `DB_DRIVER=mysql` in your backend `.env` file
2. Configure the MySQL connection parameters
3. Ensure MySQL is running and the database exists
4. The application will automatically create tables on startup

## Security Notes

- **Never commit `.env` files to version control**
- Use strong, unique JWT secrets in production
- Use environment-specific database credentials
- Consider using a secrets management service for production deployments

## Production Deployment

For production deployments:

1. Set environment variables directly on your server/container
2. Use strong, randomly generated JWT secrets
3. Use production database credentials
4. Set appropriate CORS origins
5. Consider using HTTPS for API_BASE_URL
