CREATE TABLE users (
    id TEXT PRIMARY KEY,
    username TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE folders (
    id TEXT PRIMARY KEY,
    parent_folder_id TEXT REFERENCES folders(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE
);
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
CREATE TABLE shared_files (
    id TEXT PRIMARY KEY,
    file_id TEXT NOT NULL,
    owner_id TEXT NOT NULL,
    shared_with TEXT NOT NULL,
    permissions TEXT NOT NULL
);
-- TEST-----------------
-- seed a user first (required by uploaded_by foreign key)
INSERT INTO users (id, username, email) VALUES
    ('user_1', 'henry', 'henry@example.com');

-- seed a test file referencing that user
INSERT INTO files (id, folder_id, name, size, file_type, description, uploaded_by, tags, checksum) VALUES
    ('file_1', 'NULL', 'henry.txt', 10, 'text/plain', 'descr', 'user_1', 'test', 'abc12345');
