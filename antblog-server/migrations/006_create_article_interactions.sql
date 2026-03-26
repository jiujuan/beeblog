-- ============================================================
-- Migration: 006_create_article_interactions
-- Description: 创建文章点赞表和文章收藏表
-- Author: antblog
-- Date: 2026-03-05
-- Depends: 001_create_users, 004_create_articles
-- ============================================================

-- +migrate Up

-- 文章点赞表
CREATE TABLE IF NOT EXISTS `article_likes` (
    `id`         BIGINT   UNSIGNED NOT NULL AUTO_INCREMENT   COMMENT '点赞ID',
    `article_id` BIGINT   UNSIGNED NOT NULL                  COMMENT '文章ID',
    `user_id`    BIGINT   UNSIGNED NOT NULL                  COMMENT '用户ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '点赞时间',

    PRIMARY KEY (`id`),
    -- 保证同一用户对同一文章只能点赞一次
    UNIQUE KEY `uq_likes_article_user` (`article_id`, `user_id`),
    -- 查询某用户点赞的所有文章
    KEY `idx_likes_user`    (`user_id`),
    KEY `idx_likes_created` (`created_at`),

    CONSTRAINT `fk_likes_article`
        FOREIGN KEY (`article_id`) REFERENCES `articles` (`id`)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT `fk_likes_user`
        FOREIGN KEY (`user_id`)    REFERENCES `users`    (`id`)
        ON DELETE CASCADE
        ON UPDATE CASCADE
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci
  COMMENT='文章点赞表';


-- 文章收藏表
CREATE TABLE IF NOT EXISTS `article_bookmarks` (
    `id`         BIGINT   UNSIGNED NOT NULL AUTO_INCREMENT   COMMENT '收藏ID',
    `article_id` BIGINT   UNSIGNED NOT NULL                  COMMENT '文章ID',
    `user_id`    BIGINT   UNSIGNED NOT NULL                  COMMENT '用户ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '收藏时间',

    PRIMARY KEY (`id`),
    -- 保证同一用户对同一文章只能收藏一次
    UNIQUE KEY `uq_bookmarks_article_user` (`article_id`, `user_id`),
    -- 查询某用户的收藏列表
    KEY `idx_bookmarks_user`    (`user_id`, `created_at` DESC),
    KEY `idx_bookmarks_article` (`article_id`),

    CONSTRAINT `fk_bookmarks_article`
        FOREIGN KEY (`article_id`) REFERENCES `articles` (`id`)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT `fk_bookmarks_user`
        FOREIGN KEY (`user_id`)    REFERENCES `users`    (`id`)
        ON DELETE CASCADE
        ON UPDATE CASCADE
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci
  COMMENT='文章收藏表';


-- +migrate Down

DROP TABLE IF EXISTS `article_bookmarks`;
DROP TABLE IF EXISTS `article_likes`;
