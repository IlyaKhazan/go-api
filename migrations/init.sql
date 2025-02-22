CREATE TABLE IF NOT EXISTS flights (
                                       id SERIAL PRIMARY KEY,
                                       destination_from TEXT NOT NULL,
                                       destination_to TEXT NOT NULL,
                                       deleted_at TIMESTAMP DEFAULT NULL
);
