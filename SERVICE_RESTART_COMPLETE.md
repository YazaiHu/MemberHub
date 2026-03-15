# ✅ MemberHub 服务重启完成

**重启时间**: 2026-03-15 21:13

---

## 🚀 服务状态

### ✅ 后端 API 服务
- **状态**: 正常运行
- **地址**: http://localhost:8080
- **进程ID**: 56904
- **日志文件**: `logs/api_server.log`
- **健康检查**: ✅ 通过

### ✅ 前端管理后台
- **状态**: 正常运行
- **地址**: http://localhost:3000
- **进程ID**: 59645
- **日志文件**: `logs/frontend.log`
- **访问测试**: ✅ 正常

---

## 🔐 登录信息

**管理员账号**:
- 用户名: `admin`
- 密码: `admin123`

**访问地址**: http://localhost:3000

---

## 🛠️ 解决的问题

### 问题1: 端口冲突
- **原因**: 8080 端口被旧的后端进程占用
- **解决**: 停止了进程 21170

### 问题2: 前端服务异常
- **原因**: 多个 vite 实例运行，端口切换到 3001
- **解决**: 清理所有旧进程，重新启动使用 3000 端口

---

## 📊 运行中的服务

```
进程ID  服务类型         命令
------  ------------    ---------------------------
56904   后端 API        go run cmd/api/main.go
59645   前端 Vite       npm run dev
```

---

## 📝 日志查看

### 后端日志
```bash
# 实时查看后端日志
tail -f logs/api_server.log

# 查看最近 50 行
tail -50 logs/api_server.log
```

### 前端日志
```bash
# 实时查看前端日志
tail -f logs/frontend.log

# 查看最近 50 行
tail -50 logs/frontend.log
```

---

## 🛑 停止服务

如需停止服务：

```bash
# 方法1: 使用进程ID
kill 56904 59645

# 方法2: 使用 pkill
pkill -f "cmd/api/main.go"
pkill -f "vite"

# 方法3: 停止所有相关服务
ps aux | grep -E "(cmd/api/main.go|vite)" | grep -v grep | awk '{print $2}' | xargs kill
```

---

## 🔄 重启服务

如需重新启动：

```bash
# 1. 停止旧服务
pkill -f "cmd/api/main.go"
pkill -f "vite"

# 2. 启动后端
cd /Users/niumc/Downloads/code/gostudy/MemberHub
CONFIG_PATH=configs/config.dev.yaml nohup go run cmd/api/main.go > logs/api_server.log 2>&1 &

# 3. 启动前端
cd web/admin
nohup npm run dev > ../../logs/frontend.log 2>&1 &

# 4. 等待5秒
sleep 5

# 5. 验证服务
curl http://localhost:8080/health
curl -I http://localhost:3000
```

---

## ✅ 验证清单

- [x] 后端 API 服务正常运行 (http://localhost:8080)
- [x] 前端界面服务正常运行 (http://localhost:3000)
- [x] 健康检查接口返回正常
- [x] 前端页面可以访问
- [x] 日志文件正常写入

---

## 🌐 现在可以访问

**打开浏览器访问**: http://localhost:3000

**使用以下凭证登录**:
- 用户名: `admin`
- 密码: `admin123`

---

**服务已成功重启！** 🎉
