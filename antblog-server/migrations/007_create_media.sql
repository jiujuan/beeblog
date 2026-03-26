-- ============================================================
-- Migration: 007_create_media
-- Description: 创建媒体资源表（图片上传管理）
-- Author: antblog
-- Date: 2026-03-05
-- Depends: 001_create_users, 004_create_articles
-- ============================================================

-- +migrate Up

CREATE TABLE IF NOT EXISTS `media` (
    `id`            BIGINT        UNSIGNED NOT NULL AUTO_INCREMENT        COMMENT '媒体ID',
    `uploader_id`   BIGINT        UNSIGNED NOT NULL                       COMMENT '上传者用户ID',
    `article_id`    BIGINT        UNSIGNED                                COMMENT '关联文章ID（NULL=独立资源，未绑定文章）',
    `original_name` VARCHAR(256)  NOT NULL DEFAULT ''                     COMMENT '客户端上传时的原始文件名',
    `storage_path`  VARCHAR(512)  NOT NULL                                COMMENT '服务端存储相对路径（如 uploads/2026/03/xxx.jpg）',
    `url`           VARCHAR(512)  NOT NULL                                COMMENT '可访问的公开 URL',
    `mime_type`     VARCHAR(128)  NOT NULL DEFAULT ''                     COMMENT 'MIME 类型，如 image/jpeg image/png image/webp',
    `file_size`     BIGINT        NOT NULL DEFAULT 0                      COMMENT '文件大小（字节）',
    `width`         INT           NOT NULL DEFAULT 0                      COMMENT '图片宽度（px），非图片类型为 0',
    `height`        INT           NOT NULL DEFAULT 0                      COMMENT '图片高度（px），非图片类型为 0',
    `hash`          CHAR(64)      NOT NULL DEFAULT ''                     COMMENT 'SHA256 文件内容哈希（用于去重检测）',
    `created_at`    DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP      COMMENT '上传时间',
    `updated_at`    DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP
                                  ON UPDATE CURRENT_TIMESTAMP             COMMENT '更新时间',
    `deleted_at`    DATETIME                                              COMMENT '软删除时间',

    PRIMARY KEY (`id`),
    -- 按上传者查询资源库
    KEY `idx_media_uploader`  (`uploader_id`, `created_at` DESC),
    -- 按文章查询已绑定的图片
    KEY `idx_media_article`   (`article_id`),
    -- 哈希去重查询
    KEY `idx_media_hash`      (`hash`),
    -- 按文件类型筛选
    KEY `idx_media_mime`      (`mime_type`),
    KEY `idx_media_deleted`   (`deleted_at`),

    CONSTRAINT `fk_media_uploader`
        FOREIGN KEY (`uploader_id`) REFERENCES `users`    (`id`)
        ON DELETE RESTRICT
        ON UPDATE CASCADE,

    CONSTRAINT `fk_media_article`
        FOREIGN KEY (`article_id`)  REFERENCES `articles` (`id`)
        ON DELETE SET NULL
        ON UPDATE CASCADE
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci
  COMMENT='媒体资源表（图片/文件上传管理）';


-- +migrate Down

DROP TABLE IF EXISTS `media`;
