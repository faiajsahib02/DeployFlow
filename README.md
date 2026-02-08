# DepoyFlow

A modern, full-stack application for managing Docker container deployments with ease. DepoyFlow provides a web-based interface to deploy, manage, and monitor containerized applications with a clean, intuitive user experience.

## 🚀 Features

- **Docker Integration**: Seamlessly manage Docker containers and deployments
- **Project Management**: Organize deployments by projects
- **Real-time Monitoring**: Track deployment status and container health
- **REST API**: Full-featured API for programmatic deployment management
- **Web Dashboard**: Modern React-based UI for managing deployments
- **Database Persistence**: PostgreSQL backend for reliable data storage
- **Docker Compose Support**: Easy local development setup

## 🛠️ Tech Stack

### Infrastructure & Containerization
- **Docker** - Container runtime
- **Docker Compose** - Multi-container orchestration
- **Docker SDK** - Docker API integration

### Backend
- **Go 1.25** - Fast, compiled backend service
- **Chi Router** - Lightweight HTTP router
- **PostgreSQL** - Relational database
- **CORS** - Cross-origin resource sharing support

### Frontend
- **React 19** - Modern UI framework
- **TypeScript** - Type-safe JavaScript
- **Vite** - Next-generation build tool
- **Tailwind CSS** - Utility-first CSS framework
- **Radix UI** - Accessible component library
- **Axios** - HTTP client

## 📋 Prerequisites

- Go 1.25+
- Node.js 18+
- Docker & Docker Compose
- PostgreSQL 15+ (provided via Docker Compose)

## 🚀 Quick Start

### 1. Clone the Repository
```bash
git clone https://github.com/sahib002/deployflow.git
cd deployflow
```

### 2. Start PostgreSQL with Docker Compose
```bash
docker-compose up -d
```

This will start a PostgreSQL container with the following credentials:
- **User**: postgres
- **Password**: postgres
- **Database**: deployflow
- **Port**: 5435

### 3. Run the Backend Server
```bash
go run ./cmd/server/main.go
```

The API will be available at `http://localhost:8080`

### 4. Run the Frontend Development Server
```bash
cd web
npm install
npm run dev
```

The web UI will be available at `http://localhost:3000`

## 📁 Project Structure

```
deployflow/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── adapter/
│   │   ├── handler/             # HTTP request handlers
│   │   ├── proxy/               # HTTP proxy functionality
│   │   ├── runtime/             # Docker client integration
│   │   └── storage/
│   │       └── postgres/        # Database layer
│   ├── core/
│   │   ├── domain/              # Domain models
│   │   ├── ports/               # Interface definitions
│   │   └── services/            # Business logic
│   └── utils/                   # Utility functions
├── web/                         # React frontend
│   ├── src/
│   │   ├── components/          # React components
│   │   ├── services/            # API service layer
│   │   └── lib/                 # Utility functions
│   └── public/                  # Static assets
├── docker-compose.yml           # Docker services configuration
├── go.mod & go.sum             # Go dependencies
└── README.md                    # This file
```

## 📚 API Endpoints

The API is available at `http://localhost:8080/api/v1/`

### Key Endpoints
- `GET /api/v1/projects` - List all projects
- `POST /api/v1/projects` - Create new project
- `GET /api/v1/deployments` - List deployments
- `POST /api/v1/deployments` - Create new deployment

## 🔧 Configuration

### Backend Environment Variables
The backend uses PostgreSQL connection string:
```
postgres://postgres:postgres@localhost:5435/deployflow?sslmode=disable
```

Modify this in `cmd/server/main.go` for production use.

### Frontend Configuration
Update API base URL in `web/src/services/api.ts` for different environments.

## 📦 Building for Production

### Backend
```bash
go build -o bin/deployflow ./cmd/server
```

### Frontend
```bash
cd web
npm run build
```

The built files will be in `web/dist/`

## 🧪 Testing

```bash
# Run Go tests
go test ./...

# Run frontend linting
cd web && npm run lint
```

## 📝 Database Migrations

Migrations are located in `internal/adapter/storage/postgres/migrations/`

Currently includes:
- `000001_init_schema.up.sql` - Initial schema setup

## 🤝 Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 💡 Support

For issues, questions, or suggestions, please open an issue on GitHub.

---

**Built with ❤️ by Sahib**
