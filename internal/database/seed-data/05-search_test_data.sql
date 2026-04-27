-- Seed data for testing search functionality
-- Using existing user_1 from 01-user.sql

-- Folders
INSERT INTO folders (id, parent_folder_id, name, user_id) VALUES
('f_work', NULL, 'Work Documents', 'user_1'),
('f_personal', NULL, 'Personal Stuff', 'user_1'),
('f_projects', 'f_work', 'Project Alpha', 'user_1'),
('f_archive', 'f_work', 'Old Archive', 'user_1'),
('f_photos', 'f_personal', 'Vacation Photos', 'user_1'),
('f_finances', 'f_personal', 'Financial Records 2024', 'user_1'),
('f_source', 'f_projects', 'Source Code', 'user_1'),
('f_docs', 'f_projects', 'Design Documents', 'user_1'),
('f_secret', NULL, 'Top Secret Folder', 'user_1'),
('f_tmp', NULL, 'Temporary Workspace', 'user_1');

-- Files (Starting from file_100 to avoid conflicts with 03-file.sql)
INSERT INTO files (id, folder_id, name, size, file_type, uploaded_by, tags) VALUES
('file_101', 'f_projects', 'Project_Alpha_Spec.pdf', 1024500, 'PDF', 'user_1', 'work,important,spec'),
('file_102', 'f_projects', 'Alpha_Logo.png', 54200, 'PNG', 'user_1', 'work,design,image'),
('file_103', 'f_source', 'main.go', 2400, 'GO', 'user_1', 'code,backend'),
('file_104', 'f_source', 'utils.go', 1200, 'GO', 'user_1', 'code,helper'),
('file_105', 'f_photos', 'Hawaii_Sunset.jpg', 3500000, 'JPG', 'user_1', 'personal,photo,vacation'),
('file_106', 'f_photos', 'Beach_Day.jpg', 4200000, 'JPG', 'user_1', 'personal,photo'),
('file_107', 'f_finances', 'Tax_Return_2023.pdf', 850000, 'PDF', 'user_1', 'finance,tax,important'),
('file_108', 'f_finances', 'Budget_2024.xlsx', 45000, 'XLSX', 'user_1', 'finance,budget'),
('file_109', 'f_archive', 'Legacy_System_Backup.zip', 50000000, 'ZIP', 'user_1', 'archive,old'),
('file_110', 'f_docs', 'UI_Mockups.fig', 1500000, 'FIG', 'user_1', 'design,ui'),
('file_111', NULL, 'README_FIRST.txt', 500, 'TXT', 'user_1', 'note'),
('file_112', 'f_secret', 'Secret_Passwords.txt', 200, 'TXT', 'user_1', 'private,security'),
('file_113', 'f_work', 'Annual_Report_Draft.docx', 120000, 'DOCX', 'user_1', 'work,report'),
('file_114', 'f_projects', 'Alpha_Architecture.drawio', 35000, 'DRAWIO', 'user_1', 'work,diagram'),
('file_115', 'f_source', 'go.mod', 150, 'MOD', 'user_1', 'code,config'),

-- Data for Jane Doe (user_2)
('file_201', NULL, 'Jane_Resume.pdf', 45000, 'PDF', 'user_2', 'career,personal'),
('file_202', NULL, 'Budget_Sheet_Jane.xlsx', 15000, 'XLSX', 'user_2', 'finance'),

-- Data for John Programming (user_3)
('file_301', NULL, 'Compiler_Notes.txt', 1024, 'TXT', 'user_3', 'education,programming'),
('file_302', NULL, 'Rust_Tutorial.mp4', 85000000, 'MP4', 'user_3', 'learning,video'),

-- Data for Noname (user_4)
('file_401', NULL, 'Battle_Plan.md', 5000, 'MD', 'user_4', 'strategy,secret'),

-- Data for Just Sayin (user_5)
('file_501', NULL, 'Inspirational_Quote.png', 250000, 'PNG', 'user_5', 'social');

-- Shared Files (Cross-user access)
INSERT INTO shared_files (id, file_id, owner_id, shared_with, permissions) VALUES
('sf_3', 'file_103', 'user_1', 'user_3', 'read'), -- user_1 shares code with user_3
('sf_4', 'file_401', 'user_4', 'user_1', 'write'), -- user_4 shares battle plan with user_1
('sf_5', 'file_301', 'user_3', 'user_2', 'read'); -- user_3 shares notes with user_2
