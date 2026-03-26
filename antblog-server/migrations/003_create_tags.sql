-- ============================================================
-- Migration: 003_create_tags
-- Description: 创建文章标签表
-- Author: antblog
-- Date: 2026-03-05
-- Depends: 001_create_users
-- ============================================================

-- +migrate Up

CREATE TABLE IF NOT EXISTS `tags` (
    `id`            BIGINT      UNSIGNED NOT NULL AUTO_INCREMENT        COMMENT '标签ID',
    `name`          VARCHAR(64) NOT NULL                                COMMENT '标签名称（全局唯一）',
    `slug`          VARCHAR(128) NOT NULL                               COMMENT 'URL slug（全局唯一）',
    `color`         VARCHAR(16) NOT NULL DEFAULT '#6B7280'              COMMENT '标签展示颜色（Hex 色值，如 #00ADD8）',
    `article_count` INT         NOT NULL DEFAULT 0                      COMMENT '关联文章数量（冗余字段）',
    `created_at`    DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP      COMMENT '创建时间',
    `updated_at`    DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP
                                ON UPDATE CURRENT_TIMESTAMP             COMMENT '更新时间',
    `deleted_at`    DATETIME                                            COMMENT '软删除时间（NULL=未删除）',

    PRIMARY KEY (`id`),
    UNIQUE KEY `uq_tags_name`   (`name`),
    UNIQUE KEY `uq_tags_slug`   (`slug`),
    KEY        `idx_tags_count` (`article_count` DESC),
    KEY        `idx_tags_deleted`(`deleted_at`)
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci
  COMMENT='文章标签表';


-- +migrate Down

DROP TABLE IF EXISTS `tags`;
