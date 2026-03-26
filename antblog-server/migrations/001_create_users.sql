-- ============================================================
-- Migration: 001_create_users
-- Description: 创建用户表和用户刷新令牌表
-- Author: antblog
-- Date: 2026-03-05
-- ============================================================

-- +migrate Up

CREATE TABLE IF NOT EXISTS `users` (
    `id`         BIGINT       UNSIGNED NOT NULL AUTO_INCREMENT        COMMENT '用户ID',
    `uuid`       CHAR(36)     NOT NULL                                COMMENT '全局唯一标识 UUID v4',
    `username`   VARCHAR(32)  NOT NULL                                COMMENT '用户名（3-32位字母数字下划线）',
    `email`      VARCHAR(128) NOT NULL                                COMMENT '邮箱',
    `password`   VARCHAR(255) NOT NULL                                COMMENT 'bcrypt 哈希密码',
    `nickname`   VARCHAR(64)  NOT NULL DEFAULT ''                     COMMENT '昵称，默认同用户名',
    `avatar`     VARCHAR(512) NOT NULL DEFAULT ''                     COMMENT '头像 URL',
    `bio`        VARCHAR(512) NOT NULL DEFAULT ''                     COMMENT '个人简介',
    `role`       TINYINT      NOT NULL DEFAULT 1                      COMMENT '角色: 1=普通用户 2=管理员',
    `status`     TINYINT      NOT NULL DEFAULT 1                      COMMENT '状态: 1=正常 2=禁用',
    `last_login` DATETIME                                             COMMENT '最后登录时间',
    `created_at` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP      COMMENT '创建时间',
    `updated_at` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP
                              ON UPDATE CURRENT_TIMESTAMP             COMMENT '更新时间',
    `deleted_at` DATETIME                                             COMMENT '软删除时间（NULL=未删除）',

    PRIMARY KEY (`id`),
    UNIQUE KEY `uq_users_uuid`     (`uuid`),
    UNIQUE KEY `uq_users_email`    (`email`),
    UNIQUE KEY `uq_users_username` (`username`),
    KEY        `idx_users_role`    (`role`),
    KEY        `idx_users_status`  (`status`),
    KEY        `idx_users_deleted` (`deleted_at`)
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci
  COMMENT='用户表';


CREATE TABLE IF NOT EXISTS `user_tokens` (
    `id`            BIGINT       UNSIGNED NOT NULL AUTO_INCREMENT   COMMENT '令牌ID',
    `user_id`       BIGINT       UNSIGNED NOT NULL                  COMMENT '关联用户ID',
    `refresh_token` VARCHAR(512) NOT NULL                           COMMENT 'JWT Refresh Token 字符串',
    `user_agent`    VARCHAR(512) NOT NULL DEFAULT ''                COMMENT '客户端 User-Agent',
    `client_ip`     VARCHAR(64)  NOT NULL DEFAULT ''                COMMENT '登录客户端 IP',
    `expires_at`    DATETIME     NOT NULL                           COMMENT 'Token 过期时间',
    `created_at`    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',

    PRIMARY KEY (`id`),
    KEY `idx_tokens_user_id`      (`user_id`),
    KEY `idx_tokens_refresh_token`(`refresh_token`(64)),
    KEY `idx_tokens_expires`      (`expires_at`),
    CONSTRAINT `fk_tokens_user`
        FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
        ON DELETE CASCADE
        ON UPDATE CASCADE
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci
  COMMENT='用户刷新令牌表（支持多设备登录）';


-- +migrate Down

DROP TABLE IF EXISTS `user_tokens`;
DROP TABLE IF EXISTS `users`;
