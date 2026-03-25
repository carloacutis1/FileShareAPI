# BASEMENT - RESTful File Storage API

A comprehensive file storage API backend built with Go and Gin Web Framework, designed to replicate Google Drive functionality with user account management, hierarchical folder structures, file uploads/downloads, and file sharing capabilities.

**Project:** SPR-26 Backend Web-Dev 
**Repository:** `gitlab.franciscan.edu/spr-26-backend/basement`

## Development Team

| Developer     | Primary Responsibility        |
|---------------|-------------------------------|
| Matthew Suess | Users domain & API            |
| Marley Hall   | Folders domain & Docker setup |
| Henry Cooper  | Files domain & API            |
| Joshua Fouch  | Shared Files domain & API     |

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
- **Go 1.25.0+** (for building and running the application)
- **Git** (for version control)
- **PostgreSQL client tools** (optional, for direct database interaction)

## Project Structure

```
basement/
├── go.mod                  # Go module definition with dependencies
├── README.md              # This file
├── cmd/
│   └── api/
│       └── main.go        # Application entry point and route registration
├── internal/
│   ├── database/
│   │   ├── init.sql       # Master init file—imports all other SQL files
│   │   ├── user.sql       # Users table schema
│   │   ├── folder.sql     # Folders table schema
│   │   ├── flie.sql       # Files table schema (note: filename typo preserved)
│   │   ├── shared_file.sql # Shared_files table schema
│   │   └── file-sample-data.sql # Sample data for testing
│   ├── handlers/
│   │   ├── user.go        # HTTP handler for user endpoints
│   │   ├── folder.go      # HTTP handler for folder endpoints
│   │   ├── file.go        # HTTP handler for file endpoints
│   │   └── shared_file.go # HTTP handler for shared file endpoints
│   └── models/
│       ├── user.go        # User data model and mock data
│       ├── folder.go      # Folder data model and mock data
│       ├── file.go        # File data model and mock data
│       └── shared_file.go # SharedFile data model and mock data
```

### Directory Responsibilities

#### `cmd/api/`
The entry point of the application. Contains the main server initialization, route registration, and HTTP listener setup. The Gin router is instantiated here with all API endpoints configured.

#### `internal/database/`
Contains PostgreSQL table schema definitions and sample data:
- **init.sql**: Master initialization script that imports all other SQL files in sequence
- **user.sql**: Defines the `users` table with user account information
- **folder.sql**: Defines the `folders` table with hierarchical folder structure
- **flie.sql**: Defines the `files` table (note: filename contains typo "flie" instead of "file")
- **shared_file.sql**: Defines the `shared_files` table for file sharing relationships
- **file-sample-data.sql**: Contains INSERT statements with example data for testing

#### `internal/handlers/`
HTTP request handlers that bridge Gin routes and the model layer. Each handler:
1. Extracts data from the models package
2. Serializes data to JSON
3. Returns HTTP responses (currently all return 200 OK)

Follows the pattern: Route → Handler → Model → JSON Response

#### `internal/models/`
Contains Go structures and mock data:
- Defines Go structs matching database schemas
- Includes mock data slices for initial testing (before database integration)
- Provides getter functions to retrieve all data for each domain

**Current State:** All models are using in-memory mock data. As development progresses, these will be replaced with actual database queries.

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
    file_id TEXT NOT NULL REFERENCES files(id) ON DELETE CASCADE,
    shared_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE
);
```
- **Purpose:** Implements file sharing—tracks which users have access to which files
- **Key Concept:** A row in this table grants `user_id` access to `file_id`
- **Cascade Behavior:** Deleting a file or user removes all sharing relationships
- **Example:**
  - File 101 (owned by user 1) can be shared with users 5, 6, 7 via separate rows

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

Returns a list of all shared file relationships—which files are shared with which users.

**Response (200 OK):**
```json
[
  {
    "id": 1,
    "file_id": 101,
    "owner_id": 1,
    "shared_with": "alice",
    "shared_with_id": 5
  },
  {
    "id": 2,
    "file_id": 102,
    "owner_id": 2,
    "shared_with": "bob",
    "shared_with_id": 6
  }
]
```

**Mock Data:** 5 sample shared file entries are defined in [models/shared_file.go](internal/models/shared_file.go)

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
    ID           int    `json:"id"`
    FileID       int    `json:"file_id"`
    OwnerID      int    `json:"owner_id"`
    SharedWith   string `json:"shared_with"`
    SharedWithID int    `json:"shared_with_id"`
}
```

**Purpose:** Represents a file sharing relationship between a file owner and another user.

**Mock Data Function:** `GetSharedFiles() []SharedFile`

---

## Getting Started

### 1. Clone the Repository

```bash
git clone https://gitlab.franciscan.edu/spr-26-backend/basement.git
cd basement
```

### 2. Install Go Dependencies

```bash
go mod download
```

This fetches all dependencies defined in [go.mod](go.mod), including Gin and PostgreSQL driver.

### 3. Set Up PostgreSQL Database (Docker)

```bash
docker run --name basement-postgres \
  -e POSTGRES_USER=basement \
  -e POSTGRES_PASSWORD=your_password \
  -e POSTGRES_DB=basement \
  -p 5432:5432 \
  -v basement-data:/var/lib/postgresql/data \
  -d postgres:latest
```

Then initialize the schema:

```bash
docker exec -i basement-postgres psql -U basement -d basement < internal/database/init.sql
```

**Alternative:** Use the [docker-compose.yml](docker-compose.yml) file (if present) for simplified setup.

### 4. Configure Environment (Future)

Create a `.env` file in the project root for database credentials (implementation pending).

### 5. Build and Run

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
   - Models (`internal/models/`) handle data structures and mock data
   - Handlers (`internal/handlers/`) manage HTTP request/response logic
   - Database schema (`internal/database/`) defines table structures
   - Main entrypoint (`cmd/api/`) configures routes only

2. **Adding a New Endpoint:**
   1. Define the Go struct in `internal/models/[domain].go`
   2. Create mock data and a getter function
   3. Create the handler in `internal/handlers/[domain].go`
   4. Register the route in `cmd/api/main.go`
   5. Create the SQL schema in `internal/database/[domain].sql`
   6. Add import to `internal/database/init.sql`

3. **Naming Conventions:**
   - Files: lowercase with hyphens (e.g., `shared_file.go`)
   - Functions: PascalCase for exported, camelCase for private
   - JSON tags: snake_case (e.g., `"parent_folder_id"`)
   - Database IDs: Prefixed identifiers (e.g., `user_`, `fld_`, `file_`)

4. **Comments:**
   - Add package-level comments at the top of files
   - Document purpose of functions and non-obvious logic
   - Note future implementation plans (e.g., "will be replaced with database queries")

### Current Status: Iteration 1

✅ **Completed:**
- Basic project structure with Go modules
- Gin web framework setup
- 4 mock data models with handler endpoints
- PostgreSQL schema definitions
- Simple HTTP JSON responses

🔄 **In Progress:**
- Testing environment validation
- Docker setup completion

⏳ **Planned (Iteration 2+):**
- Database integration (replace mock data with queries)
- Authentication & authorization
- File upload/download endpoints
- CRUD operations (Create, Read, Update, Delete)
- Error handling and validation
- Unit and integration tests
- API documentation (Swagger/OpenAPI)

### Debugging Tips

- **Check logs:** Look at console output for Gin debug information
- **Verify routes:** Comment out route registrations to isolate issues
- **Test endpoints:** Use `curl`, Postman, or VS Code REST Client
- **Database connection:** Confirm PostgreSQL is running with `docker ps`
- **Go errors:** Run `go vet ./...` to check for static analysis issues

---

**Last Updated:** March 25, 2026  
**Version:** 0.1.0 (Iteration 1)
