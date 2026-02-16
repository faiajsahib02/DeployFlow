# DeployFlow

DeployFlow is a streamlined platform designed to empower data scientists to deploy machine learning models to production environments effortlessly. Say goodbye to the complexities of production infrastructure, server management, and deployment pipelines. With DeployFlow, focus on what you do best—building models—while we handle the rest.

## The Core Problem: "The Model Handover Gap"

In most companies, there is a massive wall between Data Scientists and DevOps Engineers.

Data Scientists write code in Python (Jupyter Notebooks). They care about accuracy, math, and experimentation. They are bad at Docker, networking, and latency.

Production Engineers care about speed, uptime, and security. They hate messy Python scripts that crash servers.

**The Pain Point**: A Data Scientist builds a great model, emails the file to an Engineer, and the Engineer spends 2 weeks trying to make it run on a server. It is slow, manual, and frustrating.

## The Solution: DeployFlow

You built DeployFlow to automate this handover.

It is a platform where a Data Scientist can just upload their model.pkl file, and DeployFlow automatically:

- Wraps it in a Docker container.
- Spins up a high-performance Go proxy in front of it (for speed/security).
- Deploys it to a live URL (e.g., api.deployflow.com/v1/predict).

## The Technical Motivation (Why Go + Python?)

You realized that while Python is the king of AI, it is slow for handling HTTP requests and concurrency.

**Motivation**: You wanted the best of both worlds.

**Architecture**: You used Go for the "Infrastructure Layer" (handling traffic, routing, load balancing) because it is fast. You used Python only for the mathematical inference.

## 🛠️ Tech Stack

### Infrastructure & Containerization
- **Docker** - Container runtime for model isolation
- **Docker Compose** - Multi-container orchestration
- **Docker SDK** - Docker API integration for automated deployments

### Backend
- **Go 1.25** - High-performance backend service
- **Chi Router** - Lightweight HTTP router
- **PostgreSQL** - Relational database for deployment metadata
- **CORS** - Cross-origin resource sharing support

### Frontend
- **React 19** - Modern UI framework
- **TypeScript** - Type-safe JavaScript
- **Vite** - Fast build tool
- **Tailwind CSS** - Utility-first CSS framework
- **Radix UI** - Accessible component library
- **Axios** - HTTP client for API communication

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
│   │   ├── handler/             # HTTP request handlers for deployment operations
│   │   ├── proxy/               # HTTP proxy for model serving
│   │   ├── runtime/             # Docker client for container management
│   │   └── storage/
│   │       └── postgres/        # Database layer for deployment metadata
│   ├── core/
│   │   ├── domain/              # Domain models for projects and deployments
│   │   ├── ports/               # Interface definitions
│   │   └── services/            # Business logic for deployment workflows
│   └── utils/                   # Utility functions (e.g., tarball creation)
├── web/                         # React frontend
│   ├── src/
│   │   ├── components/          # React components for dashboard and deployment UI
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
- `GET /api/v1/projects` - List all ML model projects
- `POST /api/v1/projects` - Create new model project
- `GET /api/v1/deployments` - List model deployments
- `POST /api/v1/deployments` - Deploy a new model

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
- `000001_init_schema.up.sql` - Initial schema setup for projects and deployments

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


