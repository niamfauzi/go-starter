CREATE TABLE IF NOT EXISTS products (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    tenant_id VARCHAR(100) NOT NULL,
    name VARCHAR(150) NOT NULL,
    description TEXT NULL,
    price BIGINT NOT NULL DEFAULT 0,
    stock INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    KEY idx_products_tenant_id (tenant_id),
    KEY idx_products_tenant_deleted (tenant_id, deleted_at),
    KEY idx_products_name (name)
);
