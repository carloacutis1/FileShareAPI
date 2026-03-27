-- put sample data in here
-- seed a test file referencing that user
INSERT INTO files (id, folder_id, name, size, file_type, description, uploaded_by, tags, checksum) VALUES
    ('file_1', 'NULL', 'henry.txt', 10, 'text/plain', 'descr', 'user_1', 'test', 'abc12345'),
    ('file_2', 'NULL', 'budget.xlsx', 24576, 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet', 'Q1 budget sheet', 'user_2', 'finance,budget', 'bgt99871'),
    ('file_3', 'NULL', 'team-photo.jpg', 3145728, 'image/jpeg', 'Team retreat photo', 'user_3', 'image,team', 'img55632'),
    ('file_4', 'NULL', 'notes.md', 2048, 'text/markdown', 'Project notes', 'user_4', 'notes,project', 'nts44021'),
    ('file_5', 'NULL', 'diagram.pdf', 189440, 'application/pdf', 'Architecture diagram', 'user_5', 'design,architecture', 'pdf77190'),
    ('file_6', 'NULL', 'demo.mp4', 12582912, 'video/mp4', 'Demo recording', 'user_1', 'video,demo', 'vid33218');