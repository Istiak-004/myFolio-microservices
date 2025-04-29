auth-service/
├── cmd/
│   └── main.go                     # Application entry point
├── internal/
│   ├── api/                        # Interface adapters
│   │   ├── http/                   # HTTP handlers
│   │   │   ├── handlers.go         # REST endpoint handlers
│   │   │   ├── middlewares.go      # HTTP middlewares
│   │   │   └── router.go           # Route definitions
│   │   └── grpc/                   # gRPC interface
│   │       ├── auth.proto          # Protocol buffer definition
│   │       ├── server.go           # gRPC server implementation
│   │       └── service.go          # gRPC service
│   ├── domain/                     # Core business logic
│   │   ├── models/                 # Domain entities
│   │   │   └── user.go             # User model
│   │   ├── ports/                  # Interfaces (ports)
│   │   │   ├── repository.go       # Repository interface
│   │   │   ├── service.go          # Service interface
│   │   │   └── token.go            # Token interface
│   │   └── service/                # Business logic
│   │       └── auth.go             # Auth service implementation
│   └── infrastructure/             # Infrastructure adapters
│       ├── persistence/            # Database implementations
│       │   ├── postgres/           # Postgres repository
│       │   │   ├── repository.go   # Postgres implementation
│       │   │   └── migrations/     # Database migrations
│       │   └── redis/              # Redis repository
│       │       └── repository.go   # Redis implementation
│       └── oauth/                  # OAuth providers
│           ├── google.go           # Google OAuth implementation
│           └── provider.go         # OAuth provider interface
├── pkg/                            # Shared utilities
│   ├── config/                     # Configuration
│   │   ├── config.go               # Config structs
│   │   └── load.go                 # Config loading
│   └── utils/                      # Utility functions
│       ├── crypto.go               # Cryptographic functions
│       └── validator.go            # Validation helpers
├── migrations/                     # Database migrations
│   ├── 0001_init.up.sql            # Initial migration
│   └── 0001_init.down.sql          # Rollback migration
├── deployments/                    # Deployment configs
│   └── docker-compose.yaml         # Docker setup
├── .env.example                    # Environment variables
├── go.mod                          # Go module file
├── Makefile                        # Build commands
└── README.md                       # Project documentation