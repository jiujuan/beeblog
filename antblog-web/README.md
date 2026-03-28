# AntBlog Web

> Vue 3 + TypeScript + Tailwind CSS + shadcn-vue

## 技术栈

| 层级     | 技术                              |
| -------- | --------------------------------- |
| 框架     | Vue 3 (Composition API)           |
| 类型     | TypeScript 5                      |
| 状态     | Pinia 2                           |
| 路由     | Vue Router 4                      |
| HTTP     | Axios + 请求/响应拦截器           |
| 样式     | Tailwind CSS 3 + shadcn-vue       |
| Markdown | markdown-it + highlight.js        |
| 图标     | Lucide Vue Next                   |
| 工具     | @vueuse/core                      |

## 快速开始

```bash
# 安装依赖
npm install

# 开发服务器（代理 /api → localhost:8080）
npm run dev

# 类型检查
npm run type-check

# 生产构建
npm run build
```

## 目录结构

```
src/
├── api/           API 封装（axios）
├── assets/        全局样式 + 静态资源
├── components/
│   ├── ui/        shadcn-vue 基础组件
│   ├── common/    布局组件（Header / Footer / Sidebar）
│   ├── article/   文章组件（Card / Content / Like / Bookmark）
│   └── comment/   评论组件
├── composables/   组合式函数
├── router/        路由配置 + 守卫
├── stores/        Pinia Store
├── types/         TypeScript 类型定义
├── utils/         工具函数
└── views/
    ├── blog/      前台页面
    ├── auth/      登录注册
    └── admin/     后台管理
```

## 设计风格

Editorial / Magazine — Lora 衬线标题 + DM Sans 正文，暖石色调，支持深色模式切换。
