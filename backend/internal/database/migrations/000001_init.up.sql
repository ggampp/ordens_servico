-- Geospatial support
CREATE EXTENSION IF NOT EXISTS postgis;

-- =====================================================================
-- Users: authentication & authorization
-- =====================================================================
CREATE TABLE IF NOT EXISTS users (
    id            BIGSERIAL PRIMARY KEY,
    name          VARCHAR(150) NOT NULL,
    email         VARCHAR(150) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    role          VARCHAR(20) NOT NULL DEFAULT 'operator'
                  CHECK (role IN ('admin', 'supervisor', 'operator')),
    employee_id   BIGINT,
    active        BOOLEAN NOT NULL DEFAULT TRUE,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- =====================================================================
-- Employees
-- =====================================================================
CREATE TABLE IF NOT EXISTS employees (
    id          BIGSERIAL PRIMARY KEY,
    code        VARCHAR(50) NOT NULL UNIQUE,
    name        VARCHAR(150) NOT NULL,
    email       VARCHAR(150),
    phone       VARCHAR(30),
    role        VARCHAR(80),
    status      VARCHAR(10) NOT NULL DEFAULT 'active'
                CHECK (status IN ('active', 'inactive')),
    deleted     BOOLEAN NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Link users.employee_id now that employees exists
ALTER TABLE users
    ADD CONSTRAINT fk_users_employee
    FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE SET NULL;

-- =====================================================================
-- Employee position history
-- =====================================================================
CREATE TABLE IF NOT EXISTS employee_positions (
    id          BIGSERIAL PRIMARY KEY,
    employee_id BIGINT NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    latitude    DOUBLE PRECISION NOT NULL,
    longitude   DOUBLE PRECISION NOT NULL,
    geom        GEOGRAPHY(Point, 4326),
    recorded_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_positions_employee ON employee_positions(employee_id);
CREATE INDEX IF NOT EXISTS idx_positions_recorded ON employee_positions(recorded_at DESC);
CREATE INDEX IF NOT EXISTS idx_positions_geom ON employee_positions USING GIST(geom);

-- =====================================================================
-- Service Orders
-- =====================================================================
CREATE TABLE IF NOT EXISTS service_orders (
    id               BIGSERIAL PRIMARY KEY,
    number           VARCHAR(30) NOT NULL UNIQUE,
    title            VARCHAR(200) NOT NULL,
    description      TEXT,
    priority         VARCHAR(10) NOT NULL DEFAULT 'medium'
                     CHECK (priority IN ('low', 'medium', 'high', 'urgent')),
    status           VARCHAR(20) NOT NULL DEFAULT 'open'
                     CHECK (status IN ('open', 'assigned', 'in_progress', 'completed', 'cancelled')),
    employee_id      BIGINT REFERENCES employees(id) ON DELETE SET NULL,
    address          VARCHAR(300),
    latitude         DOUBLE PRECISION,
    longitude        DOUBLE PRECISION,
    geom             GEOGRAPHY(Point, 4326),
    opened_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    due_at           TIMESTAMPTZ,
    completed_at     TIMESTAMPTZ,
    notes            TEXT,
    deleted          BOOLEAN NOT NULL DEFAULT FALSE,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_orders_status ON service_orders(status);
CREATE INDEX IF NOT EXISTS idx_orders_employee ON service_orders(employee_id);
CREATE INDEX IF NOT EXISTS idx_orders_priority ON service_orders(priority);
CREATE INDEX IF NOT EXISTS idx_orders_geom ON service_orders USING GIST(geom);

-- =====================================================================
-- Service Order status history
-- =====================================================================
CREATE TABLE IF NOT EXISTS service_order_history (
    id               BIGSERIAL PRIMARY KEY,
    service_order_id BIGINT NOT NULL REFERENCES service_orders(id) ON DELETE CASCADE,
    old_status       VARCHAR(20),
    new_status       VARCHAR(20) NOT NULL,
    changed_by       BIGINT REFERENCES users(id) ON DELETE SET NULL,
    note             TEXT,
    changed_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_history_order ON service_order_history(service_order_id);
