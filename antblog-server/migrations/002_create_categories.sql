-- ============================================================
-- Migration: 002_create_categories
-- Description: 创建文章分类表
-- Author: antblog
-- Date: 2026-03-05
-- Depends: 001_create_users
-- ============================================================

-- +migrate Up

CREATE TABLE IF NOT EXISTS `categories` (
    `id`            BIGINT       UNSIGNED NOT NULL AUTO_INCREMENT        COMMENT '分类ID',
    `name`          VARCHAR(64)  NOT NULL                                COMMENT '分类名称（全局唯一）',
    `slug`          VARCHAR(128) NOT NULL                                COMMENT 'URL slug（全局唯一，用于前台路由）',
    `description`   VARCHAR(512) NOT NULL DEFAULT ''                     COMMENT '分类描述',
    `cover`         VARCHAR(512) NOT NULL DEFAULT ''                     COMMENT '分类封面图 URL',
    `sort_order`    INT          NOT NULL DEFAULT 0                      COMMENT '排序权重，数值越大越靠前',
    `article_count` INT          NOT NULL DEFAULT 0                      COMMENT '关联文章数量（冗余字段，用于加速列表查询）',
    `created_at`    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP      COMMENT '创建时间',
    `updated_at`    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP
                                 ON UPDATE CURRENT_TIMESTAMP             COMMENT '更新时间',
    `deleted_at`    DATETIME                                             COMMENT '软删除时间（NULL=未删除）',

    PRIMARY KEY (`id`),
    UNIQUE KEY `uq_categories_name`    (`name`),
    UNIQUE KEY `uq_categories_slug`    (`slug`),
    KEY        `idx_categories_sort`   (`sort_order` DESC),
    KEY        `idx_categories_deleted`(`deleted_at`)
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci
  COMMENT='文章分类表';


-- +migrate Down

DROP TABLE IF EXISTS `categories`;
