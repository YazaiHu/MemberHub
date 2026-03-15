<template>
  <div class="layout-container" ref="layoutRef">
    <el-container>
      <!-- 侧边栏 -->
      <el-aside :width="isCollapse ? '64px' : '200px'" class="sidebar">
        <div class="logo">
          <h2 v-if="!isCollapse">会员管理系统</h2>
          <h2 v-else class="logo-mini">会员</h2>
        </div>
        <el-menu
          :default-active="activeMenu"
          :collapse="isCollapse"
          :collapse-transition="false"
          router
          background-color="#001529"
          text-color="#fff"
          active-text-color="#1890ff"
        >
          <el-menu-item index="/dashboard">
            <el-icon><DataAnalysis /></el-icon>
            <template #title>仪表板</template>
          </el-menu-item>
          <el-menu-item index="/members">
            <el-icon><User /></el-icon>
            <template #title>会员管理</template>
          </el-menu-item>
          <el-menu-item index="/points">
            <el-icon><Medal /></el-icon>
            <template #title>积分管理</template>
          </el-menu-item>
          <el-menu-item index="/recharge">
            <el-icon><Wallet /></el-icon>
            <template #title>充值管理</template>
          </el-menu-item>
          <el-menu-item index="/coupons">
            <el-icon><Ticket /></el-icon>
            <template #title>优惠券管理</template>
          </el-menu-item>
          <el-menu-item index="/stores">
            <el-icon><Shop /></el-icon>
            <template #title>门店管理</template>
          </el-menu-item>
          <el-menu-item index="/promotions">
            <el-icon><PriceTag /></el-icon>
            <template #title>促销管理</template>
          </el-menu-item>
        </el-menu>
      </el-aside>

      <el-container>
        <!-- 顶部导航 -->
        <el-header class="header">
          <div class="header-left">
            <!-- 折叠按钮 -->
            <el-icon class="collapse-icon" @click="toggleCollapse">
              <Fold v-if="!isCollapse" />
              <Expand v-else />
            </el-icon>
            <span class="welcome">欢迎，{{ userInfo.username || '管理员' }}</span>
          </div>
          <div class="header-right">
            <!-- 全屏按钮 -->
            <el-tooltip :content="isFullscreen ? '退出全屏' : '全屏显示'" placement="bottom">
              <el-button
                class="action-btn"
                :icon="isFullscreen ? Close : FullScreen"
                @click="toggleFullscreen"
                circle
                text
              />
            </el-tooltip>

            <!-- 用户下拉菜单 -->
            <el-dropdown @command="handleCommand">
              <span class="user-dropdown">
                <el-icon><User /></el-icon>
                {{ userInfo.username || '管理员' }}
                <el-icon class="el-icon--right"><ArrowDown /></el-icon>
              </span>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="logout">退出登录</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
        </el-header>

        <!-- 主内容区 -->
        <el-main class="main-content">
          <router-view />
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessageBox, ElMessage } from 'element-plus'
import {
  DataAnalysis,
  User,
  Medal,
  Wallet,
  Ticket,
  Shop,
  PriceTag,
  ArrowDown,
  Fold,
  Expand,
  FullScreen,
  Close
} from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const layoutRef = ref(null)
const isCollapse = ref(false)
const isFullscreen = ref(false)

const activeMenu = computed(() => route.path)
const userInfo = computed(() => userStore.userInfo)

// 切换侧边栏折叠状态
const toggleCollapse = () => {
  isCollapse.value = !isCollapse.value
  // 保存到 localStorage
  localStorage.setItem('sidebar-collapse', isCollapse.value)
}

// 切换全屏
const toggleFullscreen = () => {
  if (!document.fullscreenElement) {
    // 进入全屏
    if (layoutRef.value) {
      layoutRef.value.requestFullscreen().then(() => {
        isFullscreen.value = true
      }).catch(err => {
        ElMessage.error('全屏失败: ' + err.message)
      })
    }
  } else {
    // 退出全屏
    document.exitFullscreen().then(() => {
      isFullscreen.value = false
    })
  }
}

// 监听全屏变化事件
const handleFullscreenChange = () => {
  isFullscreen.value = !!document.fullscreenElement
}

// 处理用户下拉菜单命令
const handleCommand = (command) => {
  if (command === 'logout') {
    ElMessageBox.confirm('确定要退出登录吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }).then(() => {
      userStore.logout()
      router.push('/login')
      ElMessage.success('已退出登录')
    })
  }
}

// 组件挂载时读取折叠状态
onMounted(() => {
  const savedCollapse = localStorage.getItem('sidebar-collapse')
  if (savedCollapse !== null) {
    isCollapse.value = savedCollapse === 'true'
  }

  // 监听全屏变化
  document.addEventListener('fullscreenchange', handleFullscreenChange)
  document.addEventListener('webkitfullscreenchange', handleFullscreenChange)
  document.addEventListener('mozfullscreenchange', handleFullscreenChange)
  document.addEventListener('MSFullscreenChange', handleFullscreenChange)
})

// 组件卸载时移除监听
onBeforeUnmount(() => {
  document.removeEventListener('fullscreenchange', handleFullscreenChange)
  document.removeEventListener('webkitfullscreenchange', handleFullscreenChange)
  document.removeEventListener('mozfullscreenchange', handleFullscreenChange)
  document.removeEventListener('MSFullscreenChange', handleFullscreenChange)
})
</script>

<style scoped>
.layout-container {
  height: 100vh;
  width: 100vw;
  overflow: hidden;
}

.sidebar {
  background-color: #001529;
  overflow-x: hidden;
  overflow-y: auto;
  transition: width 0.3s ease;
}

.sidebar::-webkit-scrollbar {
  width: 6px;
}

.sidebar::-webkit-scrollbar-thumb {
  background-color: rgba(255, 255, 255, 0.2);
  border-radius: 3px;
}

.logo {
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #002140;
  overflow: hidden;
  transition: all 0.3s ease;
}

.logo h2 {
  color: #fff;
  font-size: 18px;
  margin: 0;
  font-weight: bold;
  white-space: nowrap;
}

.logo-mini {
  font-size: 16px;
}

.header {
  background: #fff;
  border-bottom: 1px solid #f0f0f0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 20px;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 20px;
}

.collapse-icon {
  font-size: 20px;
  cursor: pointer;
  color: #666;
  transition: color 0.3s;
}

.collapse-icon:hover {
  color: #1890ff;
}

.action-btn {
  font-size: 20px;
}

.action-btn:hover {
  color: #1890ff;
}

.welcome {
  font-size: 14px;
  color: #666;
}

.user-dropdown {
  display: flex;
  align-items: center;
  cursor: pointer;
  padding: 0 12px;
  transition: background-color 0.3s;
  border-radius: 4px;
}

.user-dropdown:hover {
  background-color: #f5f5f5;
}

.main-content {
  background-color: #f0f2f5;
  padding: 0;
  overflow-y: auto;
  height: calc(100vh - 60px);
}

.el-menu {
  border-right: none;
}

/* 折叠状态下的菜单项样式 */
.el-menu--collapse {
  width: 64px;
}

.el-menu--collapse .el-menu-item {
  padding: 0 20px;
}

/* 全屏模式样式优化 */
.layout-container:fullscreen {
  background-color: #f0f2f5;
}

.layout-container:-webkit-full-screen {
  background-color: #f0f2f5;
}

.layout-container:-moz-full-screen {
  background-color: #f0f2f5;
}

.layout-container:-ms-fullscreen {
  background-color: #f0f2f5;
}
</style>
