CREATE TABLE shared_files (
    id INTEGER PRIMARY KEY,
    file_id INTEGER NOT NULL,
    owner_id INTEGER NOT NULL,
    shared_with TEXT NOT NULL,
    shared_with_id INTEGER NOT NULL
);
