# 会员管理系统 - 管理后台前端

基于 Vue 3 + Element Plus + Vite 构建的会员管理系统管理后台。

## 技术栈

- **框架**: Vue 3.4.21 (Composition API)
- **UI 组件库**: Element Plus 2.6.3
- **状态管理**: Pinia 2.1.7
- **路由**: Vue Router 4.3.0
- **HTTP 客户端**: Axios 1.6.8
- **构建工具**: Vite 5.2.0

## 功能模块

- ✅ **登录认证** - JWT Token 认证，路由守卫
- ✅ **仪表板** - 统计数据展示、快捷操作、系统信息
- ✅ **会员管理** - 会员列表、详情查询、积分/余额调整
- ✅ **积分管理** - 兑换规则配置、兑换记录、核销功能
- ✅ **充值管理** - 充值活动配置、订单查询
- ✅ **优惠券管理** - 模板管理、批量推送、推送任务
- ✅ **门店管理** - 门店 CRUD、位置信息管理
- ✅ **促销管理** - 特价商品管理

## 项目结构

```
web/admin/
├── public/              # 静态资源
├── src/
│   ├── api/            # API 接口模块
│   │   ├── auth.js     # 认证相关
│   │   ├── member.js   # 会员相关
│   │   ├── coupon.js   # 优惠券相关
│   │   └── store.js    # 门店相关
│   ├── assets/         # 资源文件（图片、样式等）
│   ├── components/     # 公共组件
│   ├── layouts/        # 布局组件
│   │   └── MainLayout.vue  # 主布局
│   ├── router/         # 路由配置
│   │   └── index.js
│   ├── stores/         # Pinia 状态管理
│   │   └── user.js     # 用户状态
│   ├── utils/          # 工具函数
│   │   └── request.js  # Axios 封装
│   ├── views/          # 页面组件
│   │   ├── auth/       # 登录页
│   │   ├── dashboard/  # 仪表板
│   │   ├── members/    # 会员管理
│   │   ├── points/     # 积分管理
│   │   ├── recharge/   # 充值管理
│   │   ├── coupons/    # 优惠券管理
│   │   ├── stores/     # 门店管理
│   │   └── promotions/ # 促销管理
│   ├── App.vue         # 根组件
│   └── main.js         # 入口文件
├── index.html
├── package.json
├── vite.config.js
└── README.md
```

## 快速开始

### 1. 安装依赖

```bash
cd web/admin
npm install
```

### 2. 启动开发服务器

```bash
npm run dev
```

开发服务器将运行在 `http://localhost:3000`

### 3. 构建生产版本

```bash
npm run build
```

构建产物将生成在 `dist/` 目录。

### 4. 预览生产版本

```bash
npm run preview
```

## 配置说明

### 后端 API 地址

开发环境下，Vite 已配置代理将 `/api` 请求转发到 `http://localhost:8080`。

如需修改，请编辑 `vite.config.js`：

```javascript
server: {
  port: 3000,
  proxy: {
    '/api': {
      target: 'http://localhost:8080',  // 修改为你的后端地址
      changeOrigin: true
    }
  }
}
```

### 默认登录账号

- **用户名**: admin
- **密码**: admin123

## 主要功能说明

### 会员管理

- 会员列表展示（头像、昵称、手机号、等级、积分、余额、优惠券）
- 搜索过滤（手机号、昵称、状态）
- 会员详情查看
- 手动调整积分/余额

### 积分管理

- **兑换规则**: 配置各门店的积分兑换规则，设置库存、每日限兑
- **兑换记录**: 查看所有兑换记录，支持核销操作

### 充值管理

- **充值活动**: 配置充值送活动（如充200送50），设置赠送金额和积分
- **充值订单**: 查看所有充值订单及支付状态

### 优惠券管理

- **优惠券模板**: 创建满减券/折扣券/代金券模板
- **批量推送**: 支持全部会员、指定等级、指定用户推送
- **推送任务**: 查看推送任务执行进度和结果
- **微信通知**: 推送时可选择发送微信模板消息

### 门店管理

- 门店信息管理（名称、地址、电话、经纬度）
- 营业时间配置
- 门店状态控制

### 促销管理

- 特价商品配置（原价、特价、库存）
- 促销时间范围设置
- 按门店筛选

## 开发指南

### 添加新页面

1. 在 `src/views/` 下创建页面组件
2. 在 `src/router/index.js` 中添加路由配置
3. 在 `src/layouts/MainLayout.vue` 中添加菜单项

### 添加新 API

1. 在 `src/api/` 下创建或编辑模块文件
2. 导出 API 函数，使用 `request` 工具发起请求

示例：

```javascript
import request from '@/utils/request'

export function getData(params) {
  return request({
    url: '/admin/data',
    method: 'get',
    params
  })
}
```

### 状态管理

使用 Pinia 进行状态管理，在 `src/stores/` 下创建 store：

```javascript
import { defineStore } from 'pinia'

export const useMyStore = defineStore('my', {
  state: () => ({
    data: null
  }),
  actions: {
    async fetchData() {
      // ...
    }
  }
})
```

## 注意事项

1. **认证**: 所有需要认证的请求会自动在请求头中添加 `Authorization: Bearer <token>`
2. **错误处理**: Axios 拦截器已配置全局错误处理，401 会自动跳转登录页
3. **响应格式**: 后端 API 统一返回格式为：
   ```json
   {
     "code": 0,
     "message": "success",
     "data": {}
   }
   ```
4. **TODO 标记**: 代码中有 `// TODO: 调用API` 的地方需要在后端 API 完成后替换为实际调用

## 部署

### Nginx 配置示例

```nginx
server {
    listen 80;
    server_name your-domain.com;
    root /var/www/admin/dist;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }

    location /api {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

## 浏览器支持

- Chrome (推荐)
- Firefox
- Safari
- Edge

不支持 IE11 及以下版本。

## License

MIT
