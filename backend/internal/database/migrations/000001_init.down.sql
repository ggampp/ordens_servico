DROP TABLE IF EXISTS service_order_history;
DROP TABLE IF EXISTS service_orders;
DROP TABLE IF EXISTS employee_positions;
ALTER TABLE IF EXISTS users DROP CONSTRAINT IF EXISTS fk_users_employee;
DROP TABLE IF EXISTS employees;
DROP TABLE IF EXISTS users;
