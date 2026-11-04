-- Uruchom po wybraniu istniejącej bazy: mysql NAZWA_BAZY < ten_plik
-- Bezpieczne dodanie modułu do bazy z tabelami users i sessions.
CREATE TABLE IF NOT EXISTS cultural_events (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  title VARCHAR(200) NOT NULL,
  description TEXT NOT NULL,
  starts_at DATETIME NOT NULL COMMENT 'UTC',
  ends_at DATETIME NOT NULL COMMENT 'UTC',
  location VARCHAR(255) NOT NULL,
  is_published BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  KEY cultural_events_public_start (is_published, starts_at),
  KEY cultural_events_public_end (is_published, ends_at),
  CONSTRAINT cultural_events_dates_check CHECK (ends_at >= starts_at),
  CONSTRAINT cultural_events_publication_check CHECK (is_published IN (0, 1))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

