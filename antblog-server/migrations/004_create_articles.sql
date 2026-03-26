-- ============================================================
-- Migration: 004_create_articles
-- Description: 创建文章主表及文章-标签关联表
-- Author: antblog
-- Date: 2026-03-05
-- Depends: 001_create_users, 002_create_categories, 003_create_tags
-- ============================================================

-- +migrate Up

CREATE TABLE IF NOT EXISTS `articles` (
    `id`             BIGINT       UNSIGNED NOT NULL AUTO_INCREMENT        COMMENT '文章ID',
    `uuid`           CHAR(36)     NOT NULL                                COMMENT '全局唯一标识 UUID v4',
    `author_id`      BIGINT       UNSIGNED NOT NULL                       COMMENT '作者用户ID',
    `category_id`    BIGINT       UNSIGNED                                COMMENT '所属分类ID（NULL=未分类）',
    `title`          VARCHAR(256) NOT NULL                                COMMENT '文章标题',
    `slug`           VARCHAR(300) NOT NULL                                COMMENT 'URL slug（全局唯一，用于前台路由）',
    `summary`        VARCHAR(512) NOT NULL DEFAULT ''                     COMMENT '文章摘要（手动填写或自动截取正文前200字）',
    `content`        LONGTEXT     NOT NULL                                COMMENT 'Markdown 原始内容',
    `content_html`   LONGTEXT     NOT NULL                                COMMENT '服务端渲染后的 HTML 内容（可用于 RSS / 缓存）',
    `cover`          VARCHAR(512) NOT NULL DEFAULT ''                     COMMENT '封面图 URL',
    `status`         TINYINT      NOT NULL DEFAULT 1                      COMMENT '发布状态: 1=草稿 2=已发布 3=已归档',
    `is_top`         TINYINT      NOT NULL DEFAULT 0                      COMMENT '是否置顶: 0=否 1=是',
    `is_featured`    TINYINT      NOT NULL DEFAULT 0                      COMMENT '是否精选: 0=否 1=是',
    `allow_comment`  TINYINT      NOT NULL DEFAULT 1                      COMMENT '是否允许评论: 0=关闭 1=开启',
    `view_count`     INT          NOT NULL DEFAULT 0                      COMMENT '阅读次数',
    `like_count`     INT          NOT NULL DEFAULT 0                      COMMENT '点赞数（与 article_likes 表冗余）',
    `comment_count`  INT          NOT NULL DEFAULT 0                      COMMENT '评论数（与 comments 表冗余）',
    `bookmark_count` INT          NOT NULL DEFAULT 0                      COMMENT '收藏数（与 article_bookmarks 表冗余）',
    `word_count`     INT          NOT NULL DEFAULT 0                      COMMENT '正文字数（Markdown 去除标记后的字符数）',
    `published_at`   DATETIME                                             COMMENT '发布时间（status=2 时由系统自动设置）',
    `created_at`     DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP      COMMENT '创建时间',
    `updated_at`     DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP
                                  ON UPDATE CURRENT_TIMESTAMP             COMMENT '更新时间',
    `deleted_at`     DATETIME                                             COMMENT '软删除时间（NULL=未删除）',

    PRIMARY KEY (`id`),
    UNIQUE KEY `uq_articles_uuid`      (`uuid`),
    UNIQUE KEY `uq_articles_slug`      (`slug`),
    KEY        `idx_articles_author`   (`author_id`),
    KEY        `idx_articles_category` (`category_id`),
    KEY        `idx_articles_status`   (`status`),
    -- 前台列表常用查询：已发布 + 按时间降序
    KEY        `idx_articles_pub_time` (`status`, `published_at` DESC),
    -- 置顶排序复合索引
    KEY        `idx_articles_top`      (`is_top` DESC, `published_at` DESC),
    -- 精选文章
    KEY        `idx_articles_featured` (`is_featured`, `status`),
    KEY        `idx_articles_deleted`  (`deleted_at`),

    CONSTRAINT `fk_articles_author`
        FOREIGN KEY (`author_id`)   REFERENCES `users`      (`id`)
        ON DELETE RESTRICT
        ON UPDATE CASCADE,

    CONSTRAINT `fk_articles_category`
        FOREIGN KEY (`category_id`) REFERENCES `categories` (`id`)
        ON DELETE SET NULL
        ON UPDATE CASCADE
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci
  COMMENT='文章主表';


-- 文章-标签 多对多关联表
CREATE TABLE IF NOT EXISTS `article_tags` (
    `article_id` BIGINT   UNSIGNED NOT NULL                  COMMENT '文章ID',
    `tag_id`     BIGINT   UNSIGNED NOT NULL                  COMMENT '标签ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '关联创建时间',

    PRIMARY KEY (`article_id`, `tag_id`),
    -- 反向查询（某标签下所有文章）
    KEY `idx_article_tags_tag` (`tag_id`),

    CONSTRAINT `fk_at_article`
        FOREIGN KEY (`article_id`) REFERENCES `articles` (`id`)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT `fk_at_tag`
        FOREIGN KEY (`tag_id`)     REFERENCES `tags`     (`id`)
        ON DELETE CASCADE
        ON UPDATE CASCADE
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci
  COMMENT='文章标签关联表（多对多）';


-- +migrate Down

DROP TABLE IF EXISTS `article_tags`;
DROP TABLE IF EXISTS `articles`;
