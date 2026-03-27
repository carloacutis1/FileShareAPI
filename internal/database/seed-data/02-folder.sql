-- Sample folders for local dev / query tests (depends on users from 01-user.sql)
-- Root first, then children (parent_folder_id -> folders.id)

INSERT INTO folders (id, parent_folder_id, name, user_id) VALUES
    ('fld_home', NULL, 'Home', 'user_1'),
    ('fld_documents', 'fld_home', 'Documents', 'user_1'),
    ('fld_pictures', 'fld_home', 'Pictures', 'user_1'),
    ('fld_work', 'fld_documents', 'Work', 'user_1'),
    ('fld_archive', 'fld_documents', 'Archive', 'user_1');
