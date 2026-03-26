-- ============================================================
-- Migration: 008_create_operation_logs
-- Description: 创建后台操作日志表
-- Author: antblog
-- Date: 2026-03-05
-- Depends: 001_create_users
-- ============================================================

-- +migrate Up

CREATE TABLE IF NOT EXISTS `operation_logs` (
    `id`          BIGINT       UNSIGNED NOT NULL AUTO_INCREMENT        COMMENT '日志ID',
    `user_id`     BIGINT       UNSIGNED NOT NULL                       COMMENT '操作用户ID',
    `module`      VARCHAR(64)  NOT NULL DEFAULT ''                     COMMENT '所属模块，如 article / category / tag / comment / media / user',
    `action`      VARCHAR(64)  NOT NULL DEFAULT ''                     COMMENT '操作类型，如 create / update / delete / publish / restore',
    `target_id`   BIGINT       UNSIGNED                                COMMENT '操作目标记录 ID（NULL=无具体目标）',
    `description` VARCHAR(512) NOT NULL DEFAULT ''                     COMMENT '操作描述，可包含变更摘要',
    `ip`          VARCHAR(64)  NOT NULL DEFAULT ''                     COMMENT '操作者 IP',
    `user_agent`  VARCHAR(512) NOT NULL DEFAULT ''                     COMMENT '操作者客户端 UA',
    `created_at`  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP      COMMENT '操作时间',

    PRIMARY KEY (`id`),
    -- 按操作用户查询
    KEY `idx_oplogs_user`   (`user_id`, `created_at` DESC),
    -- 按模块+操作类型筛选
    KEY `idx_oplogs_module` (`module`, `action`),
    -- 按目标记录查询变更历史
    KEY `idx_oplogs_target` (`module`, `target_id`),
    -- 按时间范围查询
    KEY `idx_oplogs_time`   (`created_at` DESC),

    CONSTRAINT `fk_oplogs_user`
        FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
        ON DELETE CASCADE
        ON UPDATE CASCADE
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci
  COMMENT='后台管理操作日志表';


-- +migrate Down

DROP TABLE IF EXISTS `operation_logs`;
