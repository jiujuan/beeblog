-- ============================================================
-- Migration: 009_create_views
-- Description: 创建数据库视图（文章归档时间线、分类统计）
-- Author: antblog
-- Date: 2026-03-05
-- Depends: 002_create_categories, 003_create_tags, 004_create_articles
-- ============================================================

-- +migrate Up

-- 文章归档视图：按年月聚合已发布文章数，驱动前台时间线归档页
CREATE OR REPLACE VIEW `article_archive_view` AS
SELECT
    YEAR(`published_at`)  AS `year`,
    MONTH(`published_at`) AS `month`,
    COUNT(*)              AS `article_count`
FROM `articles`
WHERE `status`     = 2        -- 只统计已发布
  AND `deleted_at` IS NULL    -- 排除软删除
  AND `published_at` IS NOT NULL
GROUP BY
    YEAR(`published_at`),
    MONTH(`published_at`)
ORDER BY
    `year`  DESC,
    `month` DESC;


-- 分类文章统计视图：实时统计各分类下已发布文章数（可作为冗余字段同步依据）
CREATE OR REPLACE VIEW `category_article_count_view` AS
SELECT
    c.`id`                                   AS `category_id`,
    c.`name`                                 AS `category_name`,
    c.`slug`                                 AS `category_slug`,
    COUNT(a.`id`)                            AS `article_count`,
    MAX(a.`published_at`)                    AS `last_published_at`
FROM `categories` c
LEFT JOIN `articles` a
    ON  a.`category_id` = c.`id`
    AND a.`status`       = 2
    AND a.`deleted_at`  IS NULL
WHERE c.`deleted_at` IS NULL
GROUP BY c.`id`, c.`name`, c.`slug`;


-- 标签文章统计视图：实时统计各标签下已发布文章数
CREATE OR REPLACE VIEW `tag_article_count_view` AS
SELECT
    t.`id`                AS `tag_id`,
    t.`name`              AS `tag_name`,
    t.`slug`              AS `tag_slug`,
    t.`color`             AS `tag_color`,
    COUNT(at.`article_id`) AS `article_count`
FROM `tags` t
LEFT JOIN `article_tags` at ON at.`tag_id` = t.`id`
LEFT JOIN `articles`      a
    ON  a.`id`         = at.`article_id`
    AND a.`status`      = 2
    AND a.`deleted_at` IS NULL
WHERE t.`deleted_at` IS NULL
GROUP BY t.`id`, t.`name`, t.`slug`, t.`color`;


-- +migrate Down

DROP VIEW IF EXISTS `tag_article_count_view`;
DROP VIEW IF EXISTS `category_article_count_view`;
DROP VIEW IF EXISTS `article_archive_view`;
