CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       username VARCHAR(50) UNIQUE NOT NULL,
                       password VARCHAR(255) NOT NULL,
                       created_at          timestamptz  NOT NULL DEFAULT now(),
                       updated_at          timestamptz  NOT NULL DEFAULT now()
);

CREATE TABLE films (
                       id SERIAL PRIMARY KEY,
                       title VARCHAR(255) UNIQUE NOT NULL,
                       director VARCHAR(100),
                       release_date DATE,
                       genre VARCHAR(50),
                       synopsis TEXT,
                       user_id INT NOT NULL,
                       created_at          timestamptz  NOT NULL DEFAULT now(),
                       updated_at          timestamptz  NOT NULL DEFAULT now(),
                       FOREIGN KEY (user_id) REFERENCES users(id)
);