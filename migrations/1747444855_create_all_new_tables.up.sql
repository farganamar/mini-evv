-- User table: users
CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    phone_number TEXT NOT NULL UNIQUE,
    roles TEXT NOT NULL DEFAULT 'CAREGIVER',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);
-- Add index for roles field
CREATE INDEX idx_users_roles ON users(roles);

-- Client table: clients
CREATE TABLE IF NOT EXISTS clients (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    address TEXT NOT NULL,
    latitude REAL NOT NULL,  -- SQLite doesn't have a POINT type, use separate lat/lng columns
    longitude REAL NOT NULL,
    phone_number TEXT NOT NULL UNIQUE,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);
-- Add index for coordinates
CREATE INDEX idx_clients_coordinates ON clients(latitude, longitude);

CREATE TABLE IF NOT EXISTS appointments (
    id TEXT PRIMARY KEY,
    client_id TEXT NOT NULL,
    caregiver_id TEXT NOT NULL,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    status TEXT NOT NULL DEFAULT 'SCHEDULED', -- SCHEDULED, IN_PROGRESS, COMPLETED, CANCELLED
    verification_code TEXT NOT NULL,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    FOREIGN KEY (client_id) REFERENCES clients(id),
    FOREIGN KEY (caregiver_id) REFERENCES users(id)
);

-- ADD index
CREATE INDEX idx_appointments_client_id ON appointments(client_id);
CREATE INDEX idx_appointments_caregiver_id ON appointments(caregiver_id);
CREATE INDEX idx_appointments_start_time ON appointments(start_time);
CREATE INDEX idx_appointments_end_time ON appointments(end_time);
CREATE INDEX idx_appointments_status ON appointments(status);

CREATE TABLE IF NOT EXISTS appointment_logs (
    id TEXT PRIMARY KEY,
    appointment_id TEXT NOT NULL,
    caregiver_id TEXT NOT NULL,
    log_type TEXT NOT NULL, -- CHECK-IN, CHECK-OUT, NOTE
    log_data TEXT NOT NULL DEFAULT '{}', -- JSON stored as text in SQLite
    latitude REAL,  -- Add explicit latitude/longitude for location logging
    longitude REAL,
    timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (appointment_id) REFERENCES appointments(id),
    FOREIGN KEY (caregiver_id) REFERENCES users(id)
);

-- ADD index
CREATE INDEX idx_appointment_logs_appointment_id ON appointment_logs(appointment_id);
CREATE INDEX idx_appointment_logs_caregiver_id ON appointment_logs(caregiver_id);
CREATE INDEX idx_appointment_logs_log_type ON appointment_logs(log_type);
