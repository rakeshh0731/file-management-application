# File Hub Application

A full-stack file management application built with a Go backend and a React frontend, designed for efficient and authenticated file handling.

## 🚀 Technology Stack

### Backend
- **Go**: Core application language
- **Gorilla/Mux**: HTTP router
- **GORM**: ORM for PostgreSQL user database
- **PostgreSQL**: User authentication and metadata
- **MongoDB**: File metadata storage
- **JWT**: For handling user authentication tokens

### Frontend
- React 18 with TypeScript
- TanStack Query (React Query) for data fetching
- Axios for API communication
- Tailwind CSS for styling

### Infrastructure
- Docker for running PostgreSQL and MongoDB databases

## 📋 Prerequisites

Before you begin, ensure you have installed:
- **Docker**: To run the databases.
- **Go** (1.21 or higher): For running the backend.
- **Node.js** (18.x or higher): For running the frontend.

## 🛠️ Installation & Setup

### 1. Start Databases with Docker

Run the following commands to start the PostgreSQL and MongoDB containers.

**PostgreSQL (for Users)**
```sh
docker run --name filehub-psql -d \
  -p 5432:5432 \
  -e POSTGRES_USER=admin \
  -e POSTGRES_PASSWORD=psql \
  -e POSTGRES_DB=filehub_users \
  postgres:latest
```

**MongoDB (for Files)**
```sh
docker run --name filehub-mongo -d -p 27017:27017 mongo:latest
```

### 2. Backend Setup (Go)

1. **Navigate to the backend directory**
   ```sh
   cd go-backend
   ```

2. **Install dependencies**
   ```sh
   go mod tidy
   ```

3. **Create environment file**
   Copy the example environment file. The default values should work with the Docker setup above.
   ```sh
   cp .env.example .env
   ```

4. **Start the development server**
   ```sh
   go run .
   ```
   The backend server will be running on `http://localhost:8000`.

### 3. Frontend Setup (React)

1. **Navigate to the frontend directory**
   ```sh
   cd frontend
   ```

2. **Install dependencies**
   ```sh
   npm install
   ```

3. **Start the development server**
   ```sh
   npm start
   ```

## 🌐 Accessing the Application

- Frontend Application: http://localhost:3000
- Backend API: http://localhost:8000/api

## 🗄️ Project Structure

```
file-management/
├── go-backend/            # Go backend
│   ├── api/               # API handlers and router
│   ├── config/            # Environment configuration
│   ├── database/          # Database connections (PSQL, Mongo)
│   ├── models/            # Data models (User, File)
│   ├── go.mod             # Go dependencies
│   └── main.go            # Application entrypoint
├── frontend/              # React frontend
│   ├── src/
│   │   ├── components/    # React components
│   │   ├── contexts/      # React contexts (e.g., AuthContext)
│   │   ├── pages/         # Page components
│   │   ├── services/      # API services
│   │   └── types/         # TypeScript types
│   └── package.json      # Node.js dependencies
└── .gitignore             # Files and folders to ignore in git
```
