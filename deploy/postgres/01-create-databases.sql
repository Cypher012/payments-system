-- Runs one, when the Postgres data volumes is first created
-- One database per service: each service has its own data.
CREATE DATABASE wallet;
CREATE DATABASE auth;
CREATE DATABASE transfer;
