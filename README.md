# BASEMENT - RESTful File Storage API

A comprehensive file storage API backend built with Go and Gin Web Framework, designed to replicate Google Drive functionality with user account management, hierarchical folder structures, file uploads/downloads, and file sharing capabilities.

**Project:** SPR-26 Backend Web-Dev
**Repository:** `github.com/carloacutis1/FileShareAPI`

## Development Team

| Developer     | Primary Responsibility    |
| ------------- | ------------------------- |
| Matthew Suess | Users domain & API        |
| Marley Hall   | Folders domain & Docker   |
| Henry Cooper  | Files domain & API        |
| Joshua Fouch  | Shared Files domain & API |

## Table of Contents

- [Overview](#overview)
- [Technology Stack](#technology-stack)
- [Prerequisites](#prerequisites)
- [Project Structure](#project-structure)
- [Database Schema](#database-schema)
- [API Endpoints](#api-endpoints)
- [Data Models](#data-models)
- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)

---

## Overview

**BASEMENT** is a backend API service for a cloud-based file storage platform. The system enables users to:

- Create and manage user accounts with unique emails
- Organize files in hierarchical folder structures
- Store files with metadata (size, type, tags, checksums)
- Share files with other users through a shared files mechanism
- Track file access and uploads with timestamps

The project is built with Go 1.25.0 and uses PostgreSQL as the primary data store. API requests and responses are formatted as JSON.

## Technology Stack

- **Language:** Go 1.25.0
- **Web Framework:** Gin v1.11.0
- **Database:** PostgreSQL (via `lib/pq`)
- **Runtime:** HTTP server listening on port 8080
- **Containerization:** Docker

**Key Dependencies:**

- `github.com/gin-gonic/gin`: REST API framework
- `github.com/lib/pq`: PostgreSQL driver
- `go.uber.org/mock`: Mocking utilities for testing

## Prerequisites

Ensure the following are installed on your development machine:

- **Docker Desktop** (for running PostgreSQL and containerized services)
  TO SETUP DOCKER CONTAINER:
    - sudo systemctl start docker
    - sudo docker compose up -d --build
  TO DISABLE CONTAINER:
    - sudo docker compose down -v
- **Go 1.25.0+** (for building and running the application)
- **Git** (for version control)
- **PostgreSQL client tools** (optional, for direct database interaction)

## Project Structure

```
basement/
├── go.mod                  # Go module definition with dependencies
├── go.sum                  # Dependency checksums
├── README.md              # This file
├── Dockerfile             # Container image definition
├── docker-compose.yml     # Multi-container Docker application setup
├── cmd/
│   └── api/
│       └── main.go        # Application entry point and server initialization
├── internal/
│   ├── database/
│   │   ├── 00-init.sql    # Master schema initialization file
│   │   ├── db.go          # Database connection and initialization logic
│   │   ├── queries/
│   │   │   └── shared-file.go # Database queries for shared files
│   │   └── seed-data/
│   │       ├── 01-user.sql       # Sample user data
│   │       ├── 02-folder.sql     # Sample folder data
│   │       ├── 03-file.sql       # Sample file data
│   │       └── 04-shared_file.sql # Sample shared file data
│   ├── handlers/
│   │   ├── user.go        # HTTP handler for user endpoints
│   │   ├── folder.go      # HTTP handler for folder endpoints
│   │   ├── file.go        # HTTP handler for file endpoints
│   │   └── shared_file.go # HTTP handler for shared file endpoints
│   ├── models/
│   │   ├── user.go        # User data model
│   │   ├── folder.go      # Folder data model
│   │   ├── file.go        # File data model
│   │   └── shared_file.go # SharedFile data model
│   └── routes/
│       ├── routes.go      # Main route registration and router setup
│       ├── user.go        # User route definitions
│       ├── folder.go      # Folder route definitions
│       ├── file.go        # File route definitions
│       └── shared_file.go # Shared file route definitions
```

### Directory Responsibilities

#### `cmd/api/`

The application entry point. Contains the main server initialization and startup logic. The server configuration and HTTP listener are set up here, delegating route registration to the `internal/routes` package.

#### `internal/database/`

Manages database schema, connections, and queries:

- **00-init.sql**: Master schema initialization script that defines all database tables
- **db.go**: Database connection pool management and initialization functions
- **queries/**: Contains domain-specific database query functions (currently shared-file queries)
- **seed-data/**: SQL files with INSERT statements for development and testing data, numbered for execution order

#### `internal/handlers/`

HTTP request handlers that process API requests and generate responses. Each handler:

1. Receives HTTP requests via Gin context
2. Interacts with models and database queries
3. Serializes data to JSON
4. Returns appropriate HTTP responses with status codes

Follows the pattern: Route → Handler → Database/Model → JSON Response

**Current Implementation:**
- **Shared files:** Uses database queries from `internal/database/queries/`
- **Users, files, folders:** Still using mock data from `internal/models/`

#### `internal/models/`

Contains Go data structures that represent the application's core entities:

- Defines Go structs matching database table schemas
- Includes JSON struct tags for API serialization
- Provides type definitions for data passed between handlers and database layers

#### `internal/routes/`

Centralized route definition and registration:

- **routes.go**: Main router setup and initialization, registers all domain routes
- **Domain route files**: Each file defines HTTP endpoints for a specific domain (users, folders, files, shared files)
- Maps HTTP methods and paths to their corresponding handlers

Separates routing configuration from handler implementation for better organization.

---

## Database Schema

The PostgreSQL database consists of four core tables with defined relationships:

### Users Table

```sql
CREATE TABLE users (
    id TEXT PRIMARY KEY,
    username TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

- **Purpose:** Stores user account information
- **Key Fields:**
  - `id`: Unique user identifier (text-based)
  - `username`: Human-readable user name
  - `email`: Unique email address for authentication
  - `created_at`: Account creation timestamp (auto-set)

### Folders Table

```sql
CREATE TABLE folders (
    id TEXT PRIMARY KEY,
    parent_folder_id TEXT REFERENCES folders(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE
);
```

- **Purpose:** Implements hierarchical folder structure for file organization
- **Key Features:**
  - Self-referential: `parent_folder_id` references another folder's `id`, enabling nested structures
  - Cascade delete: Deleting a folder deletes all subfolders and files
  - User-scoped: Each folder belongs to exactly one user
- **Example Hierarchy:**
  ```
  user_1
  └── fld_1 (Home, parent_folder_id = NULL)
      ├── fld_2 (Documents)
      └── fld_3 (Pictures)
  ```

### Files Table

```sql
CREATE TABLE files (
    id TEXT PRIMARY KEY,
    folder_id TEXT REFERENCES folders(id) ON DELETE SET NULL,
    name TEXT NOT NULL,
    size INTEGER NOT NULL,
    file_type TEXT,
    uploaded_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    description TEXT,
    uploaded_by TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tags TEXT,
    last_accessed TIMESTAMP,
    checksum TEXT,
    custom_attributes JSONB
);
```

- **Purpose:** Stores file metadata and references
- **Key Fields:**
  - `id`: Unique file identifier
  - `folder_id`: Optional—file may exist at root level (NULL) or in a folder
  - `name`: Filename (not full path)
  - `size`: File size in bytes
  - `file_type`: MIME type or file extension (e.g., "PDF")
  - `uploaded_by`: User who uploaded the file (cascade delete)
  - `description`: Optional text description
  - `tags`: Comma-separated tags for classification
  - `last_accessed`: Timestamp of last view/download
  - `checksum`: Hash for integrity verification
  - `custom_attributes`: JSON field for extensible metadata
- **Cascade Behavior:** Deleting `folder_id` sets it to NULL (not cascade), preserving the file; deleting the uploader user cascade-deletes the file

### Shared Files Table

```sql
CREATE TABLE shared_files (
    id TEXT PRIMARY KEY,
    file_id TEXT NOT NULL,
    owner_id TEXT NOT NULL,
    shared_with TEXT NOT NULL,
    permissions TEXT NOT NULL
);
```

- **Purpose:** Implements file sharing—tracks which users have access to which files and their permission levels
- **Key Fields:**
  - `id`: Unique sharing relationship identifier
  - `file_id`: Reference to the file being shared
  - `owner_id`: User who owns/is sharing the file
  - `shared_with`: Email address of the user receiving access
  - `permissions`: Access level granted (e.g., "read", "write", "admin")
- **Key Concept:** A row in this table grants the user identified by `shared_with` access to `file_id` with specified `permissions`
- **Example:**
  - File 'f1' (owned by user 'u1') can be shared with bob@example.com (read access) and charlie@example.com (write access) via separate rows

### Schema Relationships

```
users (1)
  ├─ (1:N)─ folders (user_id FK)
  │           └─ (1:N self-ref)─ folders (parent_folder_id FK)
  ├─ (1:N)─ files (uploaded_by FK)
  │           └─ (0:1)─ folders (folder_id FK, nullable)
  └─ (1:N)─ shared_files (user_id FK)
              └─ (N:1)─ files (file_id FK)
```

---

## API Endpoints

The API is currently in Iteration 1 with basic endpoints returning mock data. All endpoints return JSON responses with HTTP status 200 (OK). The server runs on **port 8080**.

### Users Endpoint

**GET /users**

Returns a list of all user models with basic user information.

**Response (200 OK):**

```json
[
  {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "createdat": "2024-01-15T10:30:00Z"
  },
  {
    "id": 2,
    "name": "Jane Doe",
    "email": "jane@example.com",
    "createdat": "2024-01-16T14:22:00Z"
  }
]
```

**Mock Data:** 5 sample users are defined in [models/user.go](internal/models/user.go)

### Folders Endpoint

**GET /folders**

Returns a list of all folder objects, including parent-child relationships and ownership.

**Response (200 OK):**

```json
[
  {
    "id": "fld_1",
    "parent_folder_id": "",
    "name": "Home",
    "created_at": "2024-03-21T08:45:00Z",
    "user_id": "user_1"
  },
  {
    "id": "fld_2",
    "parent_folder_id": "fld_1",
    "name": "Documents",
    "created_at": "2024-03-22T12:30:00Z",
    "user_id": "user_1"
  }
]
```

**Mock Data:** 3 sample folders are defined in [models/folder.go](internal/models/folder.go)

### Files Endpoint

**GET /files**

Returns a list of all file objects with ownership and location information.

**Response (200 OK):**

```json
[
  {
    "id": 1,
    "file_id": 1001,
    "owner_id": 42
  },
  {
    "id": 2,
    "file_id": 1002,
    "owner_id": 42
  }
]
```

**Mock Data:** 5 sample files are defined in [models/file.go](internal/models/file.go)

**Note:** Current mock model is simplified. Full file metadata will be implemented during database integration.

### Shared Files Endpoint

**GET /shared-files**

Returns a list of all shared file relationships—which files are shared with which users and their permission levels.

**Response (200 OK):**

```json
[
  {
    "id": "sf1",
    "file_id": "f1",
    "owner_id": "u1",
    "shared_with": "bob@example.com",
    "permissions": "read"
  },
  {
    "id": "sf2",
    "file_id": "f1",
    "owner_id": "u1",
    "shared_with": "charlie@example.com",
    "permissions": "write"
  }
]
```

**Database Integration:** This endpoint queries the live PostgreSQL database via `internal/database/queries/shared-file.go`

**Error Response (500 Internal Server Error):**

```json
{
  "error": "Failed to fetch shared files"
}
```

---

## Data Models

Each domain area (Users, Folders, Files, Shared Files) has a corresponding Go struct that mirrors the database table schema. These models use JSON struct tags for API serialization.

### User Model

**File:** [internal/models/user.go](internal/models/user.go)

```go
type UserModel struct {
    ID        int    `json:"id"`
    Name      string `json:"name"`
    Email     string `json:"email"`
    CreatedAt string `json:"createdat"`
}
```

**Mock Data Function:** `GetAllUserModels() []UserModel`

### Folder Model

**File:** [internal/models/folder.go](internal/models/folder.go)

```go
type Folder struct {
    ID             string    `json:"id"`
    ParentFolderID string    `json:"parent_folder_id"`
    Name           string    `json:"name"`
    CreatedAt      time.Time `json:"created_at"`
    UserID         string    `json:"user_id"`
}
```

**Key Feature:** `ParentFolderID` enables hierarchical folder structures. An empty string indicates a root folder.

**Mock Data Function:** `GetFolders() []Folder`

### File Model

**File:** [internal/models/file.go](internal/models/file.go)

```go
type File struct {
    ID      int `json:"id"`
    FileID  int `json:"file_id"`
    OwnerID int `json:"owner_id"`
}
```

**Note:** This is a simplified initial model. The database schema in [database/flie.sql](internal/database/flie.sql) includes additional fields: folder_id, name, size, file_type, uploaded_at, description, tags, last_accessed, checksum, and custom_attributes.

**Mock Data Function:** `GetFiles() []File`

### SharedFile Model

**File:** [internal/models/shared_file.go](internal/models/shared_file.go)

```go
type SharedFile struct {
    ID           string `json:"id"`
    FileID       string `json:"file_id"`
    OwnerID      string `json:"owner_id"`
    SharedWith   string `json:"shared_with"`
    Permissions  string `json:"permissions"`
}
```

**Purpose:** Represents a file sharing relationship between a file owner and another user, with specific permission levels.

**Database Integration:** This model uses live database queries via `internal/database/queries/shared-file.go` instead of mock data.

---

## Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/carloacutis1/FileShareAPI.git
cd basement
```

### 2. Install Go Dependencies

```bash
go mod download
```

This fetches all dependencies defined in [go.mod](go.mod), including Gin and PostgreSQL driver.

### 3. Set Up PostgreSQL Database (Docker)

The easiest way to set up the database is using Docker Compose:

```bash
docker-compose up -d
```

This will:
- Start a PostgreSQL container with the correct configuration
- Automatically initialize the database schema from `internal/database/00-init.sql`
- Load sample data for the shared_files table
- Start the API server connected to the database

**Alternative Manual Setup:**

If you prefer to run PostgreSQL manually:

```bash
docker run --name basement-postgres \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=database \
  -p 5432:5432 \
  -v basement-data:/var/lib/postgresql/data \
  -d postgres:16-alpine
```

Then initialize the schema:

```bash
docker exec -i basement-postgres psql -U postgres -d database < internal/database/00-init.sql
```

### 4. Configure Environment

The application uses the `DATABASE_URL` environment variable to connect to PostgreSQL.

**Default (local development):**
```
postgres://postgres:postgres@localhost:5432/database?sslmode=disable
```

**Docker Compose:** The DATABASE_URL is automatically configured in `docker-compose.yml`

**Custom configuration:** Set the `DATABASE_URL` environment variable before running the application:

```bash
export DATABASE_URL="postgres://user:password@host:port/dbname?sslmode=disable"
```

### 5. Build and Run

**Using Docker Compose (Recommended):**

```bash
docker-compose up
```

The API server starts on `http://localhost:8080` with the database automatically connected.

**Manual Build and Run:**

If you prefer to run the Go application directly (ensure PostgreSQL is running separately):

```bash
go run cmd/api/main.go
```

The server starts on `http://localhost:8080`.

### 6. Test the API

```bash
# Test users endpoint
curl http://localhost:8080/users

# Test folders endpoint
curl http://localhost:8080/folders

# Test files endpoint
curl http://localhost:8080/files

# Test shared files endpoint
curl http://localhost:8080/shared-files
```

---

## Development Workflow

### Code Organization Principles

1. **Separation of Concerns:**
   - Models (`internal/models/`) define data structures
   - Queries (`internal/database/queries/`) handle database operations
   - Handlers (`internal/handlers/`) manage HTTP request/response logic
   - Routes (`internal/routes/`) define endpoint configurations
   - Database schema (`internal/database/`) defines table structures and initialization
   - Main entrypoint (`cmd/api/`) initializes the application

2. **Adding a New Endpoint:**
   1. Define the Go struct in `internal/models/[domain].go`
   2. Create database query functions in `internal/database/queries/[domain].go`
   3. Create the handler in `internal/handlers/[domain].go`
   4. Register routes in `internal/routes/[domain].go`
   5. Add route setup call in `internal/routes/routes.go`
   6. (Schema should already exist in `internal/database/00-init.sql`)

3. **Naming Conventions:**
   - Files: lowercase with underscores (e.g., `shared_file.go`)
   - Functions: PascalCase for exported, camelCase for private
   - JSON tags: snake_case (e.g., `"parent_folder_id"`)
   - Database IDs: Prefixed identifiers (e.g., `u1`, `f1`, `sf1`)

4. **Comments:**
   - Add package-level comments at the top of files
   - Document purpose of functions and non-obvious logic
   - Note future implementation plans (e.g., "will be replaced with database queries")

### Current Status: Iteration 2

✅ **Completed:**

- Basic project structure with Go modules
- Gin web framework setup
- 4 data models with handler endpoints
- PostgreSQL schema definitions
- Docker and Docker Compose configuration
- Modular routes layer (extracted from main.go)
- Database connection pool management
- **Shared files domain:** Full database integration with queries layer
- Error handling for database operations (shared files)

🔄 **In Progress:**

- Database integration for remaining domains (users, files, folders)
- Expanding queries layer for all models

⏳ **Planned (Iteration 3+):**

- Complete database integration (replace remaining mock data)
- Authentication & authorization
- File upload/download endpoints with actual file storage
- CRUD operations (Create, Update, Delete) for all resources
- Input validation and comprehensive error handling
- Unit and integration tests
- API documentation (Swagger/OpenAPI)
- Pagination for list endpoints

### Debugging Tips

- **Check logs:** Look at console output for Gin debug information
- **Verify routes:** Comment out route registrations to isolate issues
- **Test endpoints:** Use `curl`, Postman, or VS Code REST Client
- **Database connection:** Confirm PostgreSQL is running with `docker ps`
- **Go errors:** Run `go vet ./...` to check for static analysis issues

---

**Last Updated:** March 26, 2026  
**Version:** 0.2.0 (Iteration 2 - Partial Database Integration)
