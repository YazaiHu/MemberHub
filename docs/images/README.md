# MemberHub 图片资源

本目录包含 MemberHub 项目的所有图片资源。

## 目录结构

- `screenshots/` - UI 界面截图（登录页、仪表板、管理页面等）
- `architecture/` - 架构图、系统设计图、ER 图
- `demo/` - 功能演示 GIF 动画

## 命名规范

### 截图命名
- `login-page.png` - 登录页面
- `dashboard.png` - 仪表板
- `member-list.png` - 会员列表
- `store-management.png` - 门店管理
- `points-exchange.png` - 积分兑换
- `coupon-management.png` - 优惠券管理

### 架构图命名
- `system-architecture.png` - 系统架构图
- `ddd-layers.png` - DDD 分层架构
- `database-er.png` - 数据库 ER 图
- `api-flow.png` - API 请求流程

### 演示动画命名
- `login-demo.gif` - 登录演示
- `member-crud-demo.gif` - 会员 CRUD 演示
- `points-exchange-demo.gif` - 积分兑换演示

## 图片规范

### 格式
- **截图**: PNG 格式（支持透明背景，质量好）
- **照片/复杂图**: JPG 格式（文件更小）
- **动画**: GIF 格式（< 5MB）或使用 GitHub 支持的 MP4

### 尺寸
- **全屏截图**: 宽度 1920px 或 1440px
- **局部截图**: 宽度 800-1200px
- **Logo/图标**: 宽度 200-400px
- **缩略图**: 宽度 300px

### 优化
使用工具压缩图片以减小仓库体积：
- PNG: TinyPNG (https://tinypng.com/)
- JPG: JPEGmini
- GIF: ezgif.com 或使用 MP4 替代

## 在 Markdown 中引用

```markdown
# 相对路径引用（推荐）
![登录页面](docs/images/screenshots/login-page.png)

# 带链接的图片
[![仪表板](docs/images/screenshots/dashboard.png)](docs/images/screenshots/dashboard.png)

# 指定大小（HTML 方式）
<img src="docs/images/screenshots/member-list.png" width="600" alt="会员列表">

# 并排显示多张图片
<p align="center">
  <img src="docs/images/screenshots/login-page.png" width="45%" alt="登录">
  <img src="docs/images/screenshots/dashboard.png" width="45%" alt="仪表板">
</p>
```
