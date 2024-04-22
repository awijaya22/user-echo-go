CREATE TABLE users (
	id serial PRIMARY KEY,
	full_name VARCHAR ( 64 ) NOT NULL,
  phone_number VARCHAR ( 16 ) UNIQUE NOT NULL,
  password VARCHAR ( 255 ) NOT NULL,
  sucessful_login_count INT DEFAULT 0
);



