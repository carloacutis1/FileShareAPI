-- put sample data in here
-- password_hash values are bcrypt hashes of "password"
INSERT INTO users (id, username, email, dob, password_hash) VALUES
('user_1', 'John Doe', 'jdoe001@programming.com', '01/01/1900', '$2b$10$EnxcfsfGyZZPwkeWIkMLze2qRS1vhBsUeWaxCX6AYDhHbdQ9RKxY6'),
('user_2', 'Jane Doe', 'jdoe002@programming.com', '12/31/1999', '$2b$10$EnxcfsfGyZZPwkeWIkMLze2qRS1vhBsUeWaxCX6AYDhHbdQ9RKxY6'),
('user_3', 'John Programming', 'jprogramming001@programming.com', '04/05/2026', '$2b$10$EnxcfsfGyZZPwkeWIkMLze2qRS1vhBsUeWaxCX6AYDhHbdQ9RKxY6'),
('user_4', 'Noname the Barbarian', 'nbarbarian001@programming.com', '12/25/2025', '$2b$10$EnxcfsfGyZZPwkeWIkMLze2qRS1vhBsUeWaxCX6AYDhHbdQ9RKxY6'),
('user_5', 'Just Sayin', 'jsayin001@programming.com', '02/14/2012', '$2b$10$EnxcfsfGyZZPwkeWIkMLze2qRS1vhBsUeWaxCX6AYDhHbdQ9RKxY6');
