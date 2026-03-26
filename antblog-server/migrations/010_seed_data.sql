-- ============================================================
-- Migration: 010_seed_data
-- Description: 初始化种子数据（管理员账号、默认分类、默认标签）
-- Author: antblog
-- Date: 2026-03-05
-- Depends: 001_create_users, 002_create_categories, 003_create_tags
-- Note: 仅用于开发/首次部署，生产环境请在部署后立即修改管理员密码
-- ============================================================

-- +migrate Up

-- ── 默认管理员账号 ────────────────────────────────────────────────────────────
-- 初始密码：Admin@2026
-- bcrypt hash (cost=10): $2a$10$GfxMXS0M2Vzrw7I.ehMpHuoufOp1tC1UfnI0QYBYd3B8vx0bsPFQu
-- ⚠️ 请在生产环境首次登录后立即修改密码
INSERT INTO `users`
    (`uuid`, `username`, `email`, `password`, `nickname`, `role`, `status`)
VALUES
    (
        'a0000000-0000-0000-0000-000000000001',
        'admin',
        'admin@antblog.dev',
        '$2a$10$GfxMXS0M2Vzrw7I.ehMpHuoufOp1tC1UfnI0QYBYd3B8vx0bsPFQu',
        'AntBlog Admin',
        2,   -- role=2 管理员
        1    -- status=1 正常
    )
ON DUPLICATE KEY UPDATE `id` = `id`;  -- 幂等：已存在时跳过


-- ── 默认文章分类 ─────────────────────────────────────────────────────────────
INSERT INTO `categories`
    (`name`, `slug`, `description`, `sort_order`)
VALUES
    ('技术',    'tech',          '编程、架构、工具等技术相关文章',     100),
    ('随笔',    'essay',         '生活随笔、读书笔记与个人感悟',        90),
    ('工具',    'tools',         '开发工具、效率软件与资源推荐',         80),
    ('开源',    'open-source',   '开源项目介绍与参与经验分享',           70),
    ('未分类',  'uncategorized', '暂未归入具体分类的文章',                0)
ON DUPLICATE KEY UPDATE `id` = `id`;


-- ── 默认文章标签 ─────────────────────────────────────────────────────────────
INSERT INTO `tags`
    (`name`, `slug`, `color`)
VALUES
    ('Go',          'go',           '#00ADD8'),
    ('Vue',         'vue',          '#42B883'),
    ('TypeScript',  'typescript',   '#3178C6'),
    ('MySQL',       'mysql',        '#4479A1'),
    ('Redis',       'redis',        '#DC382D'),
    ('Docker',      'docker',       '#2496ED'),
    ('Kubernetes',  'kubernetes',   '#326CE5'),
    ('架构设计',    'architecture', '#7C3AED'),
    ('性能优化',    'performance',  '#F59E0B'),
    ('开源',        'open-source',  '#10B981'),
    ('Git',         'git',          '#F05032'),
    ('Linux',       'linux',        '#FCC624')
ON DUPLICATE KEY UPDATE `id` = `id`;


-- +migrate Down

-- 仅删除种子数据，保留表结构
DELETE FROM `tags`       WHERE `slug` IN ('go','vue','typescript','mysql','redis','docker','kubernetes','architecture','performance','open-source','git','linux');
DELETE FROM `categories` WHERE `slug` IN ('tech','essay','tools','open-source','uncategorized');
DELETE FROM `users`      WHERE `uuid` = 'a0000000-0000-0000-0000-000000000001';
