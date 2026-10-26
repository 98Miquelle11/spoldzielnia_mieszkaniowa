-- MySQL 8 / MariaDB. Importuj przed uruchomieniem backendu.
-- Bez DROP, bez usuwania danych. Schemat dla obecnego etapu projektu.
CREATE DATABASE IF NOT EXISTS mirmilowo
  CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE mirmilowo;

CREATE TABLE IF NOT EXISTS users (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  email VARCHAR(254) NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  role ENUM('mieszkaniec', 'admin') NOT NULL DEFAULT 'mieszkaniec',
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY users_email_unique (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS sessions (
  token_hash CHAR(64) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
  user_id BIGINT UNSIGNED NOT NULL,
  expires_at DATETIME(6) NOT NULL COMMENT 'UTC',
  PRIMARY KEY (token_hash),
  KEY sessions_expiry (expires_at),
  KEY sessions_user (user_id),
  CONSTRAINT sessions_user_fk FOREIGN KEY (user_id)
    REFERENCES users (id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
