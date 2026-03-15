<template>
  <div class="dashboard">
    <h2 class="page-title">仪表板</h2>

    <!-- 统计卡片 -->
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon user">
              <el-icon><User /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">1,234</div>
              <div class="stat-label">总会员数</div>
            </div>
          </div>
        </el-card>
      </el-col>

      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon points">
              <el-icon><Medal /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">56,789</div>
              <div class="stat-label">今日积分发放</div>
            </div>
          </div>
        </el-card>
      </el-col>

      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon money">
              <el-icon><Wallet /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">¥12,345</div>
              <div class="stat-label">今日充值金额</div>
            </div>
          </div>
        </el-card>
      </el-col>

      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon coupon">
              <el-icon><Ticket /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">890</div>
              <div class="stat-label">今日优惠券发放</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 快捷操作 -->
    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="24">
        <el-card>
          <template #header>
            <span>快捷操作</span>
          </template>
          <div class="quick-actions">
            <el-button type="primary" @click="$router.push('/members')">
              <el-icon><User /></el-icon>
              会员管理
            </el-button>
            <el-button type="success" @click="$router.push('/coupons')">
              <el-icon><Ticket /></el-icon>
              发放优惠券
            </el-button>
            <el-button type="warning" @click="$router.push('/points')">
              <el-icon><Medal /></el-icon>
              积分兑换规则
            </el-button>
            <el-button type="info" @click="$router.push('/stores')">
              <el-icon><Shop /></el-icon>
              门店管理
            </el-button>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 系统信息 -->
    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>系统信息</span>
          </template>
          <el-descriptions :column="1" border>
            <el-descriptions-item label="系统版本">v1.0.0</el-descriptions-item>
            <el-descriptions-item label="技术栈">
              Go + Gin + MySQL + Redis + RabbitMQ
            </el-descriptions-item>
            <el-descriptions-item label="前端框架">
              Vue 3 + Element Plus
            </el-descriptions-item>
            <el-descriptions-item label="架构模式">DDD 领域驱动设计</el-descriptions-item>
          </el-descriptions>
        </el-card>
      </el-col>

      <el-col :span="12">
        <el-card>
          <template #header>
            <span>功能模块</span>
          </template>
          <ul class="feature-list">
            <li>✅ 会员管理 - 用户信息、资产查询</li>
            <li>✅ 积分系统 - 积分兑换、高并发控制</li>
            <li>✅ 充值系统 - 充值活动、微信支付</li>
            <li>✅ 优惠券系统 - 模板管理、批量推送</li>
            <li>✅ 门店管理 - CRUD、附近门店查询</li>
            <li>✅ 促销管理 - 特价商品管理</li>
            <li>✅ Worker 系统 - 异步任务处理</li>
          </ul>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { User, Medal, Wallet, Ticket, Shop } from '@element-plus/icons-vue'
</script>

<style scoped>
.dashboard {
  padding: 20px;
}

.page-title {
  margin: 0 0 20px 0;
  font-size: 24px;
  font-weight: 600;
  color: #333;
}

.stat-card {
  cursor: pointer;
  transition: all 0.3s;
}

.stat-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.stat-content {
  display: flex;
  align-items: center;
}

.stat-icon {
  width: 60px;
  height: 60px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 28px;
  color: #fff;
  margin-right: 20px;
}

.stat-icon.user {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.stat-icon.points {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.stat-icon.money {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
}

.stat-icon.coupon {
  background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
}

.stat-value {
  font-size: 28px;
  font-weight: bold;
  color: #333;
}

.stat-label {
  font-size: 14px;
  color: #999;
  margin-top: 5px;
}

.quick-actions {
  display: flex;
  gap: 10px;
}

.feature-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.feature-list li {
  padding: 8px 0;
  color: #666;
  line-height: 1.8;
}
</style>
