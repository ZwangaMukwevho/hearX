-- 01_create_tasks_table.sql
CREATE TABLE IF NOT EXISTS tasks (
  id           BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  title        VARCHAR(255) NOT NULL,                    -- index target
  description  TEXT        NULL,
  completed    BOOLEAN     NOT NULL DEFAULT FALSE,

  created_at   TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at   TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP
                            ON UPDATE CURRENT_TIMESTAMP,
  deleted_at   TIMESTAMP   NULL DEFAULT NULL,

  INDEX idx_tasks_title (title)                          -- added index
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
