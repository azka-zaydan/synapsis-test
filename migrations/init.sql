CREATE TABLE IF NOT EXISTS user (
    id CHAR(36) PRIMARY KEY NOT NULL,
    username VARCHAR(255),
    password_hash TEXT NOT NULL,
    created_by CHAR(36) NOT NULL,
    meta_created_at TIMESTAMP NOT NULL,
    updated_by CHAR(36) NOT NULL,
    meta_updated_at TIMESTAMP NOT NULL,
    deleted_by CHAR(36),
    meta_deleted_at TIMESTAMP,
    INDEX idx_username (username)
);

CREATE TABLE IF NOT EXISTS `order` (
    id CHAR(36) PRIMARY KEY NOT NULL,
    order_detail CHAR(36) NOT NULL,
    user_id CHAR(36) NOT NULL,
    product_id CHAR(36) NOT NULL,
    total_price DECIMAL(10, 2) NOT NULL,
    status INT NOT NULL,
    order_at TIMESTAMP NOT NULL,
    payment_at TIMESTAMP,
    completed_at TIMESTAMP,
    created_by CHAR(36) NOT NULL,
    meta_created_at TIMESTAMP NOT NULL,
    updated_by CHAR(36) NOT NULL,
    meta_updated_at TIMESTAMP NOT NULL,
    deleted_by CHAR(36),
    meta_deleted_at TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_product_id (product_id),
    INDEX idx_status (status),
    INDEX idx_order_detail (order_detail),
    INDEX idx_created_by (created_by)
);

CREATE TABLE IF NOT EXISTS category (
    id CHAR(36) PRIMARY KEY NOT NULL,
    name VARCHAR(266) NOT NULL,
    created_by CHAR(36) NOT NULL,
    meta_created_at TIMESTAMP NOT NULL,
    updated_by CHAR(36) NOT NULL,
    meta_updated_at TIMESTAMP NOT NULL,
    deleted_by CHAR(36),
    meta_deleted_at TIMESTAMP,
    INDEX idx_created_by (created_by)
);

CREATE TABLE IF NOT EXISTS product (
    id CHAR(36) PRIMARY KEY NOT NULL,
    category_id CHAR(36) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description VARCHAR(255) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    stock INT NOT NULL,
    created_by CHAR(36) NOT NULL,
    meta_created_at TIMESTAMP NOT NULL,
    updated_by CHAR(36) NOT NULL,
    meta_updated_at TIMESTAMP NOT NULL,
    deleted_by CHAR(36),
    meta_deleted_at TIMESTAMP,
    INDEX idx_category_id (category_id),
    INDEX idx_name (name),
    INDEX idx_created_by (created_by)
);

CREATE TABLE IF NOT EXISTS cart (
    id CHAR(36) PRIMARY KEY NOT NULL,
    user_id CHAR(36) NOT NULL,
    total_items INT NOT NULL,
    total_price DECIMAL(10, 2) NOT NULL,
    created_by CHAR(36) NOT NULL,
    meta_created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by CHAR(36) NOT NULL,
    meta_updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_by CHAR(36),
    meta_deleted_at TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_created_by (created_by)
);


CREATE TABLE IF NOT EXISTS cart_item (
    id CHAR(36) PRIMARY KEY NOT NULL,
    cart_id CHAR(36) NOT NULL,
    product_id CHAR(36) NOT NULL,
    quantity INT NOT NULL,
    created_by CHAR(36) NOT NULL,
    meta_created_at TIMESTAMP NOT NULL,
    updated_by CHAR(36) NOT NULL,
    meta_updated_at TIMESTAMP NOT NULL,
    deleted_by CHAR(36),
    meta_deleted_at TIMESTAMP,
    INDEX idx_cart_id (cart_id),
    INDEX idx_product_id (product_id),
    INDEX idx_created_by (created_by)
);

CREATE TABLE IF NOT EXISTS payment (
    id CHAR(36) PRIMARY KEY NOT NULL,
    payment_method CHAR(36) NOT NULL,
    status INT NOT NULL,
    created_by CHAR(36) NOT NULL,
    meta_created_at TIMESTAMP NOT NULL,
    updated_by CHAR(36) NOT NULL,
    meta_updated_at TIMESTAMP NOT NULL,
    deleted_by CHAR(36),
    meta_deleted_at TIMESTAMP,
    INDEX idx_payment_method (payment_method),
    INDEX idx_status (status),
    INDEX idx_created_by (created_by)
);

CREATE TABLE IF NOT EXISTS order_detail (
    id CHAR(36) PRIMARY KEY NOT NULL,
    product_id CHAR(36) NOT NULL,
    total_items INT NOT NULL,
    subtotal_product_price DECIMAL(10, 2) NOT NULL,
    created_by CHAR(36) NOT NULL,
    meta_created_at TIMESTAMP NOT NULL,
    updated_by CHAR(36) NOT NULL,
    meta_updated_at TIMESTAMP NOT NULL,
    deleted_by CHAR(36),
    meta_deleted_at TIMESTAMP,
    INDEX idx_product_id (product_id),
    INDEX idx_created_by (created_by)
);
