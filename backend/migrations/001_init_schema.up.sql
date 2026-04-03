-- 001_init_schema.up.sql
-- 政策雷达数据库初始化脚本

-- 政策分类表
CREATE TABLE IF NOT EXISTS ${prefix}_policy_category (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(50) UNIQUE,
    description VARCHAR(500),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 国标政策表
CREATE TABLE IF NOT EXISTS ${prefix}_national_standard (
    id INT AUTO_INCREMENT PRIMARY KEY,
    link1 VARCHAR(500),
    link2 VARCHAR(500),
    standard_no VARCHAR(100),
    standard_name VARCHAR(500),
    english_name VARCHAR(500),
    publish_date VARCHAR(20),
    implement_date VARCHAR(20),
    status VARCHAR(20),
    nature VARCHAR(20),
    category VARCHAR(20),
    is_adopted VARCHAR(10),
    ccs_code VARCHAR(20),
    ics_code VARCHAR(50),
    department VARCHAR(100),
    technical_dept VARCHAR(100),
    publisher VARCHAR(200),
    description TEXT,
    standard_type VARCHAR(20) DEFAULT '国标',
    download_url VARCHAR(500),
    check_status INT DEFAULT 0,
    category_id INT,
    source VARCHAR(20) DEFAULT 'manual',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    INDEX idx_standard_no (standard_no),
    INDEX idx_standard_name (standard_name(191)),
    INDEX idx_category_id (category_id),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 行标政策表
CREATE TABLE IF NOT EXISTS ${prefix}_industry_standard (
    id INT AUTO_INCREMENT PRIMARY KEY,
    detail_link VARCHAR(500),
    standard_no VARCHAR(100) UNIQUE,
    standard_name VARCHAR(500),
    publish_date VARCHAR(20),
    implement_date VARCHAR(20),
    revision_type VARCHAR(20),
    ccs_code VARCHAR(20),
    ics_code VARCHAR(50),
    technical_owner VARCHAR(200),
    approve_dept VARCHAR(100),
    industry_class VARCHAR(100),
    standard_class VARCHAR(50),
    status VARCHAR(20),
    replace_standard VARCHAR(200),
    standard_type VARCHAR(20) DEFAULT '行标',
    download_url VARCHAR(500),
    check_status INT DEFAULT 0,
    category_id INT,
    source VARCHAR(20) DEFAULT 'manual',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    INDEX idx_standard_no (standard_no),
    INDEX idx_category_id (category_id),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 地标政策表
CREATE TABLE IF NOT EXISTS ${prefix}_local_standard (
    id INT AUTO_INCREMENT PRIMARY KEY,
    detail_link VARCHAR(500),
    standard_no VARCHAR(100) UNIQUE,
    standard_name VARCHAR(500),
    publish_date VARCHAR(20),
    implement_date VARCHAR(20),
    status VARCHAR(20),
    nature VARCHAR(20),
    ccs_code VARCHAR(20),
    ics_code VARCHAR(50),
    department VARCHAR(100),
    publisher VARCHAR(200),
    description TEXT,
    standard_type VARCHAR(20) DEFAULT '地标',
    download_url VARCHAR(500),
    check_status INT DEFAULT 0,
    category_id INT,
    source VARCHAR(20) DEFAULT 'manual',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    INDEX idx_standard_no (standard_no),
    INDEX idx_category_id (category_id),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 推荐政策表
CREATE TABLE IF NOT EXISTS ${prefix}_policy_recommend (
    id INT AUTO_INCREMENT PRIMARY KEY,
    policy_id INT NOT NULL,
    policy_type VARCHAR(20) NOT NULL,
    policy_name VARCHAR(500),
    title VARCHAR(200),
    content TEXT,
    sort INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    INDEX idx_policy_id (policy_id),
    INDEX idx_policy_type (policy_type),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 爬虫日志表
CREATE TABLE IF NOT EXISTS ${prefix}_crawler_log (
    id INT AUTO_INCREMENT PRIMARY KEY,
    script TEXT,
    query VARCHAR(200),
    result TEXT,
    status VARCHAR(20),
    count INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;