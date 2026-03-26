-- ============================================================
-- Migration: 005_create_comments
-- Description: 创建评论表（支持二级嵌套、游客评论）
-- Author: antblog
-- Date: 2026-03-05
-- Depends: 001_create_users, 004_create_articles
-- ============================================================

-- +migrate Up

CREATE TABLE IF NOT EXISTS `comments` (
    `id`          BIGINT       UNSIGNED NOT NULL AUTO_INCREMENT        COMMENT '评论ID',
    `article_id`  BIGINT       UNSIGNED NOT NULL                       COMMENT '所属文章ID',
    `user_id`     BIGINT       UNSIGNED                                COMMENT '评论用户ID（NULL=游客）',
    `parent_id`   BIGINT       UNSIGNED                                COMMENT '父评论ID（NULL=顶级评论）',
    `root_id`     BIGINT       UNSIGNED                                COMMENT '根评论ID，用于聚合同一楼层的所有子评论',
    `reply_to_id` BIGINT       UNSIGNED                                COMMENT '回复的目标评论ID（区分直接回复和楼层内回复）',
    `nickname`    VARCHAR(64)  NOT NULL DEFAULT ''                     COMMENT '游客昵称（登录用户置空）',
    `email`       VARCHAR(128) NOT NULL DEFAULT ''                     COMMENT '游客邮箱（用于 Gravatar 头像，登录用户置空）',
    `content`     TEXT         NOT NULL                                COMMENT '评论正文（纯文本，前端做 XSS 过滤）',
    `ip`          VARCHAR(64)  NOT NULL DEFAULT ''                     COMMENT '评论者 IP（用于反垃圾）',
    `user_agent`  VARCHAR(512) NOT NULL DEFAULT ''                     COMMENT '客户端 User-Agent',
    `status`      TINYINT      NOT NULL DEFAULT 1                      COMMENT '审核状态: 1=待审核 2=已通过 3=已拒绝 4=垃圾评论',
    `like_count`  INT          NOT NULL DEFAULT 0                      COMMENT '评论点赞数',
    `created_at`  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP      COMMENT '创建时间',
    `updated_at`  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP
                               ON UPDATE CURRENT_TIMESTAMP             COMMENT '更新时间',
    `deleted_at`  DATETIME                                             COMMENT '软删除时间',

    PRIMARY KEY (`id`),
    -- 前台：按文章+审核状态查评论列表
    KEY `idx_comments_article`  (`article_id`, `status`),
    -- 按根评论聚合子评论
    KEY `idx_comments_root`     (`root_id`, `status`),
    -- 按父评论查子评论
    KEY `idx_comments_parent`   (`parent_id`),
    -- 后台：按审核状态筛选
    KEY `idx_comments_status`   (`status`, `created_at` DESC),
    -- 关联用户
    KEY `idx_comments_user`     (`user_id`),
    KEY `idx_comments_deleted`  (`deleted_at`),

    CONSTRAINT `fk_comments_article`
        FOREIGN KEY (`article_id`) REFERENCES `articles` (`id`)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT `fk_comments_user`
        FOREIGN KEY (`user_id`)    REFERENCES `users`    (`id`)
        ON DELETE SET NULL
        ON UPDATE CASCADE
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci
  COMMENT='文章评论表（支持二级嵌套，支持游客评论）';


-- +migrate Down

DROP TABLE IF EXISTS `comments`;
