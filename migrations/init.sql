CREATE TABLE IF NOT EXISTS flights (
                                       id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                                       destination_from TEXT NOT NULL,
                                       destination_to TEXT NOT NULL,
                                       deleted_at TIMESTAMP DEFAULT NULL
);
