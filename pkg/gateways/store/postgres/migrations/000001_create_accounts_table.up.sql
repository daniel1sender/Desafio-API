CREATE TABLE IF NOT EXISTS accounts(
id UUID PRIMARY KEY,
name TEXT NOT NULL,
cpf TEXT UNIQUE NOT NULL,
secret TEXT NOT NULL,
balance BIGINT NOT NULL,
created_at TIMESTAMP WITH TIME ZONE NOT NULL
);
