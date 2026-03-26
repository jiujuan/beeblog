# AntBlog — 数据库迁移说明

## 文件列表与执行顺序

| 文件 | 说明 | 依赖 |
|------|------|------|
| `001_create_users.sql`              | 用户表 + 刷新令牌表           | —                          |
| `002_create_categories.sql`         | 文章分类表                    | 001                        |
| `003_create_tags.sql`               | 文章标签表                    | 001                        |
| `004_create_articles.sql`           | 文章主表 + 文章标签关联表     | 001, 002, 003              |
| `005_create_comments.sql`           | 评论表（支持二级嵌套）        | 001, 004                   |
| `006_create_article_interactions.sql` | 点赞表 + 收藏表             | 001, 004                   |
| `007_create_media.sql`              | 媒体资源表                    | 001, 004                   |
| `008_create_operation_logs.sql`     | 后台操作日志表                | 001                        |
| `009_create_views.sql`              | 归档视图 + 统计视图           | 002, 003, 004              |
| `010_seed_data.sql`                 | 初始种子数据                  | 001, 002, 003              |

## 执行方式

### 方式一：直接 MySQL 执行（按顺序）

```bash
mysql -u root -p antblog < migrations/001_create_users.sql
mysql -u root -p antblog < migrations/002_create_categories.sql
mysql -u root -p antblog < migrations/003_create_tags.sql
mysql -u root -p antblog < migrations/004_create_articles.sql
mysql -u root -p antblog < migrations/005_create_comments.sql
mysql -u root -p antblog < migrations/006_create_article_interactions.sql
mysql -u root -p antblog < migrations/007_create_media.sql
mysql -u root -p antblog < migrations/008_create_operation_logs.sql
mysql -u root -p antblog < migrations/009_create_views.sql
mysql -u root -p antblog < migrations/010_seed_data.sql
```

### 方式二：Shell 一键执行

```bash
for f in migrations/0*.sql; do
  echo ">>> Running $f"
  mysql -u root -p"your_password" antblog < "$f"
done
```

### 方式三：golang-migrate（推荐生产环境）

```bash
# 安装
go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# 执行 Up（升级到最新）
migrate -path ./migrations -database "mysql://root:password@tcp(127.0.0.1:3306)/antblog" up

# 回滚一步
migrate -path ./migrations -database "mysql://root:password@tcp(127.0.0.1:3306)/antblog" down 1

# 查看当前版本
migrate -path ./migrations -database "mysql://..." version
```

## 迁移文件格式

每个文件遵循以下结构：

```sql
-- +migrate Up
-- （正向迁移 SQL，建表、加字段等）

-- +migrate Down
-- （反向回滚 SQL，删表、删字段等）
```

## 注意事项

1. **执行顺序**：严格按文件编号顺序执行，外键依赖不能乱序
2. **幂等性**：所有建表语句均使用 `CREATE TABLE IF NOT EXISTS`，重复执行安全
3. **软删除**：所有业务表均含 `deleted_at` 字段，物理删除仅用于 `article_tags`（级联）
4. **种子数据**：`010_seed_data.sql` 中管理员初始密码为 `Admin@2026`，**生产环境请立即修改**
5. **字符集**：全库统一使用 `utf8mb4 + utf8mb4_unicode_ci`，确保完整 Unicode 支持（含 Emoji）
