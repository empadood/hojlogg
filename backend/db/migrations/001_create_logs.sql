-- 001_create_logs.sql
-- Initial schema for hojlogg

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS logs (
    id           UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    odometer_km  NUMERIC(10,2) NOT NULL CHECK (odometer_km >= 0),
    fuel_level   NUMERIC(5,2)  CHECK (fuel_level IS NULL OR (fuel_level >= 0 AND fuel_level <= 100)),
    notes        TEXT          NOT NULL DEFAULT '',
    image_path   TEXT          NOT NULL DEFAULT '',
    parsed_by_ocr BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at   TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_logs_created_at ON logs (created_at DESC);
