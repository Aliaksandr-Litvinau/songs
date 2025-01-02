-- up.sql
CREATE TABLE groups (
                        id SERIAL PRIMARY KEY,
                        name VARCHAR(255) UNIQUE NOT NULL
);

CREATE TABLE songs (
                       id SERIAL PRIMARY KEY,
                       group_id INTEGER NOT NULL,
                       title VARCHAR(255) NOT NULL,
                       release_date TIMESTAMP NOT NULL,
                       text TEXT,
                       link VARCHAR(255),
                       FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE
);
