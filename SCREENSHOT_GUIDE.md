# 📸 截图指南

## ✅ 服务已启动

### 后端 API
- **地址**: http://localhost:8080
- **状态**: ✅ 运行中
- **日志**: `logs/api_server.log`

### 前端管理后台
- **地址**: http://localhost:3000
- **状态**: ✅ 运行中
- **日志**: `web/admin/frontend.log`

---

## 🔐 登录信息

**管理员账号**:
- 用户名: `admin`
- 密码: `admin123`

---

## 📸 截图步骤

### 1. 打开浏览器访问管理后台

```
http://localhost:3000
```

### 2. 按顺序截图以下页面

#### ① 登录页面 (login-page.png)
- 访问 http://localhost:3000
- **不要登录**，直接截图登录页面
- 保存为: `docs/images/screenshots/login-page.png`

#### ② 登录并进入仪表板 (dashboard.png)
- 输入用户名: `admin`
- 输入密码: `admin123`
- 点击登录
- 截图仪表板页面
- 保存为: `docs/images/screenshots/dashboard.png`

#### ③ 会员管理 (member-list.png)
- 点击左侧菜单 "会员管理"
- 等待列表加载
- 截图会员列表页面
- 保存为: `docs/images/screenshots/member-list.png`

#### ④ 门店管理 (store-management.png)
- 点击左侧菜单 "门店管理"
- 等待列表加载
- 截图门店管理页面
- 保存为: `docs/images/screenshots/store-management.png`

#### ⑤ 积分管理 (points-exchange.png)
- 点击左侧菜单 "积分管理"
- 等待列表加载
- 截图积分管理页面
- 保存为: `docs/images/screenshots/points-exchange.png`

#### ⑥ 充值管理 (recharge-management.png)
- 点击左侧菜单 "充值管理"
- 等待列表加载
- 截图充值管理页面
- 保存为: `docs/images/screenshots/recharge-management.png`

#### ⑦ 优惠券管理 (coupon-management.png)
- 点击左侧菜单 "优惠券管理"
- 等待列表加载
- 截图优惠券管理页面
- 保存为: `docs/images/screenshots/coupon-management.png`

#### ⑧ 促销管理 (promotion-management.png)
- 点击左侧菜单 "促销管理"
- 等待列表加载
- 截图促销管理页面
- 保存为: `docs/images/screenshots/promotion-management.png`

---

## 🎨 截图技巧

### 推荐设置
- **浏览器窗口宽度**: 1440px 或 1920px（全屏）
- **截图格式**: PNG
- **截图工具**:
  - macOS: `Shift + Cmd + 4` (选区截图)
  - CleanShot X (专业工具)

### 注意事项
1. ✅ 确保浏览器地址栏和书签栏**隐藏**或**不截入**
2. ✅ 关闭浏览器扩展图标（AdBlock 等）
3. ✅ 使用**统一的浏览器**（推荐 Chrome）
4. ✅ 截图时确保页面内容**完整显示**
5. ✅ 数据是测试数据，无敏感信息

---

## 🗜️ 图片优化

截图完成后，压缩图片：

### 方法一：在线压缩（推荐）
1. 访问 https://tinypng.com/
2. 上传所有 PNG 截图
3. 下载压缩后的文件
4. 替换原文件

### 方法二：命令行压缩
```bash
cd docs/images/screenshots

# 使用 ImageMagick (如果已安装)
for img in *.png; do
  convert "$img" -quality 85 -resize 1440x\> "opt-$img"
  mv "opt-$img" "$img"
done
```

**目标**: 每张图片 < 500KB

---

## ✅ 完成后

### 1. 检查图片
```bash
cd /Users/niumc/Downloads/code/gostudy/memberhub
ls -lh docs/images/screenshots/
```

应该看到 8 个 PNG 文件：
- login-page.png
- dashboard.png
- member-list.png
- store-management.png
- points-exchange.png
- recharge-management.png
- coupon-management.png
- promotion-management.png

### 2. 提交到 Git
```bash
git add docs/images/screenshots/
git add README.md
git commit -m "docs: add UI screenshots of management backend

- Added 8 high-quality screenshots
- Login page, dashboard, and all management modules
- All images optimized to < 500KB
- Updated README.md with visual preview section
"
git push origin main
```

### 3. 查看效果
推送后访问 GitHub 仓库查看 README.md 的界面预览区块：
https://github.com/YazaiHu/MemberHub

---

## 🛑 停止服务

截图完成后，如需停止服务：

```bash
# 查看运行中的服务
ps aux | grep -E "(cmd/api/main.go|vite)" | grep -v grep

# 停止服务（会自动查找并停止）
pkill -f "cmd/api/main.go"
pkill -f "vite"
```

或者使用 Ctrl+C 在对应终端停止。

---

## 💡 额外建议

### 可选：录制 GIF 演示
如果想添加动画演示：

1. 下载录屏工具
   - macOS: [Kap](https://getkap.co/) (免费)
   - 或者: [Gifox](https://gifox.io/)

2. 录制关键流程
   - 登录流程 (5-10秒)
   - 添加会员流程
   - 积分兑换流程

3. 保存到
   - `docs/images/demo/login-demo.gif`
   - `docs/images/demo/member-add-demo.gif`
   - `docs/images/demo/points-exchange-demo.gif`

4. 压缩 GIF
   - 使用 https://ezgif.com/optimize
   - 目标: < 3MB

---

**准备好了吗？现在可以开始截图了！** 🎉

浏览器访问: http://localhost:3000
