# README.md 中添加图片的示例

## 📸 如何在 README.md 中添加 UI 截图

### 1. 单张大图展示

```markdown
## 🎨 界面预览

### 登录页面
![登录页面](docs/images/screenshots/login-page.png)

### 管理后台仪表板
![仪表板](docs/images/screenshots/dashboard.png)
```

---

### 2. 多张图片并排显示（推荐）

```markdown
## 🎨 界面预览

<div align="center">

### 登录与仪表板
<p>
  <img src="docs/images/screenshots/login-page.png" width="48%" alt="登录页面">
  <img src="docs/images/screenshots/dashboard.png" width="48%" alt="仪表板">
</p>

### 会员管理与门店管理
<p>
  <img src="docs/images/screenshots/member-list.png" width="48%" alt="会员管理">
  <img src="docs/images/screenshots/store-management.png" width="48%" alt="门店管理">
</p>

### 积分管理与优惠券管理
<p>
  <img src="docs/images/screenshots/points-exchange.png" width="48%" alt="积分管理">
  <img src="docs/images/screenshots/coupon-management.png" width="48%" alt="优惠券管理">
</p>

</div>
```

---

### 3. 可点击放大的图片

```markdown
## 🎨 界面预览

点击图片查看大图：

[![登录页面](docs/images/screenshots/login-page.png)](docs/images/screenshots/login-page.png)
[![仪表板](docs/images/screenshots/dashboard.png)](docs/images/screenshots/dashboard.png)
```

---

### 4. 图片画廊风格（最专业）

```markdown
## 🎨 界面预览

<table>
  <tr>
    <td width="50%">
      <h3 align="center">登录页面</h3>
      <img src="docs/images/screenshots/login-page.png" alt="登录页面">
    </td>
    <td width="50%">
      <h3 align="center">管理后台仪表板</h3>
      <img src="docs/images/screenshots/dashboard.png" alt="仪表板">
    </td>
  </tr>
  <tr>
    <td width="50%">
      <h3 align="center">会员管理</h3>
      <img src="docs/images/screenshots/member-list.png" alt="会员管理">
    </td>
    <td width="50%">
      <h3 align="center">门店管理</h3>
      <img src="docs/images/screenshots/store-management.png" alt="门店管理">
    </td>
  </tr>
  <tr>
    <td width="50%">
      <h3 align="center">积分兑换</h3>
      <img src="docs/images/screenshots/points-exchange.png" alt="积分兑换">
    </td>
    <td width="50%">
      <h3 align="center">优惠券管理</h3>
      <img src="docs/images/screenshots/coupon-management.png" alt="优惠券管理">
    </td>
  </tr>
</table>
```

---

### 5. GIF 动画演示

```markdown
## 🎬 功能演示

### 会员管理演示
![会员管理演示](docs/images/demo/member-crud-demo.gif)

### 积分兑换流程
![积分兑换演示](docs/images/demo/points-exchange-demo.gif)
```

---

### 6. 架构图展示

```markdown
## 🏗️ 系统架构

### DDD 分层架构
![DDD 分层架构](docs/images/architecture/ddd-layers.png)

### 数据库 ER 图
![数据库设计](docs/images/architecture/database-er.png)

### API 请求流程
![API 流程](docs/images/architecture/api-flow.png)
```

---

## 📋 推荐的 README.md 图片区块结构

建议在 README.md 的以下位置添加图片：

```markdown
# MemberHub - 连锁会员管理系统

[徽章区域]

**一个会员，通行全店** | One Member, All Stores

[功能特性] • [快速开始] • [技术栈] • [AI 开发说明]

---

## 🤖 AI 辅助开发说明
[现有内容...]

---

## 🎨 界面预览  👈 在这里添加 UI 截图

<div align="center">

### 管理后台界面

<table>
  <tr>
    <td width="50%">
      <h4 align="center">登录页面</h4>
      <img src="docs/images/screenshots/login-page.png" alt="登录">
    </td>
    <td width="50%">
      <h4 align="center">仪表板</h4>
      <img src="docs/images/screenshots/dashboard.png" alt="仪表板">
    </td>
  </tr>
  <tr>
    <td width="50%">
      <h4 align="center">会员管理</h4>
      <img src="docs/images/screenshots/member-list.png" alt="会员管理">
    </td>
    <td width="50%">
      <h4 align="center">门店管理</h4>
      <img src="docs/images/screenshots/store-management.png" alt="门店管理">
    </td>
  </tr>
</table>

</div>

---

## 🎬 功能演示  👈 在这里添加演示动画

### 会员管理流程
![会员管理演示](docs/images/demo/member-crud-demo.gif)

---

## 📖 项目简介
[现有内容...]
```

---

## 💡 图片优化技巧

### 1. 截图工具推荐
- **macOS**: Shift + Cmd + 4（选区截图）
- **专业工具**:
  - [CleanShot X](https://cleanshot.com/) - macOS 最佳截图工具
  - [Snipaste](https://www.snipaste.com/) - 跨平台免费

### 2. GIF 录制工具
- **macOS**:
  - [Kap](https://getkap.co/) - 开源免费
  - [Gifox](https://gifox.io/) - 专业工具
- **Windows**:
  - [ScreenToGif](https://www.screentogif.com/) - 免费开源

### 3. 图片压缩（重要！）
```bash
# 使用 ImageMagick 批量压缩 PNG
cd docs/images/screenshots
for img in *.png; do
  convert "$img" -quality 85 -resize 1440x\> "optimized-$img"
done

# 使用 gifsicle 压缩 GIF
gifsicle -O3 --colors 256 input.gif -o output.gif
```

**在线压缩工具**：
- PNG: https://tinypng.com/
- JPG: https://compressor.io/
- GIF: https://ezgif.com/optimize

### 4. 图片尺寸建议
```markdown
# 全屏截图（管理后台）
宽度: 1440px 或 1920px
格式: PNG
压缩后大小: < 500KB

# 局部截图（某个模块）
宽度: 800-1000px
格式: PNG
压缩后大小: < 200KB

# GIF 动画
宽度: 800-1000px
帧率: 10-15 fps
时长: 5-10 秒
压缩后大小: < 3MB
```

---

## ⚠️ 注意事项

1. **隐私保护**: 截图前确保没有敏感信息（真实手机号、邮箱、API keys）
2. **使用测试数据**: 显示的会员数据应该是虚拟的测试数据
3. **统一风格**: 所有截图使用相同的浏览器、相同的窗口尺寸
4. **高清显示**: 使用 Retina 屏幕截图，然后缩放到合适尺寸
5. **文件命名**: 使用小写字母和连字符，不要用空格和中文
6. **版本控制**: 大量图片会增加仓库体积，考虑使用 Git LFS

---

## 🚀 快速开始

### 1. 拍摄截图
```bash
# 启动前端
cd web/admin
npm run dev

# 在浏览器中访问 http://localhost:3000
# 使用 admin/admin123 登录
# 依次访问各个页面并截图
```

### 2. 保存截图
```bash
# 将截图保存到对应目录
cp ~/Desktop/login.png docs/images/screenshots/login-page.png
cp ~/Desktop/dashboard.png docs/images/screenshots/dashboard.png
# ... 其他截图
```

### 3. 优化图片
```bash
# 使用 TinyPNG 或其他工具压缩
# 确保每张图片 < 500KB
```

### 4. 更新 README.md
```bash
# 在 README.md 的合适位置添加图片展示区块
# 使用上面提供的模板
```

### 5. 提交到 Git
```bash
git add docs/images/
git add README.md
git commit -m "docs: add UI screenshots and demo GIFs"
git push origin main
```
