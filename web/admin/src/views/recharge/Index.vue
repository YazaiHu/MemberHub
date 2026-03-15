<template>
  <div class="recharge-container">
    <el-tabs v-model="activeTab">
      <!-- 充值活动管理 -->
      <el-tab-pane label="充值活动" name="promotions">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>充值活动列表</span>
              <el-button type="primary" @click="handleAddPromotion">
                <el-icon><Plus /></el-icon>
                新增活动
              </el-button>
            </div>
          </template>

          <el-table
            :data="promotionsData"
            border
            stripe
            v-loading="promotionsLoading"
            :empty-text="promotionsLoading ? '加载中...' : '暂无数据'"
            style="width: 100%"
          >
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="name" label="活动名称" width="200" />
            <el-table-column label="充值金额" width="120">
              <template #default="{ row }">
                <span class="money-text">¥{{ row.recharge_amount }}</span>
              </template>
            </el-table-column>
            <el-table-column label="赠送金额" width="120">
              <template #default="{ row }">
                <span class="bonus-text">¥{{ row.bonus_amount }}</span>
              </template>
            </el-table-column>
            <el-table-column label="赠送积分" width="120">
              <template #default="{ row }">
                <span class="points-text">{{ row.bonus_points || 0 }}</span>
              </template>
            </el-table-column>
            <el-table-column label="活动时间" width="300">
              <template #default="{ row }">
                {{ row.start_time }} ~ {{ row.end_time }}
              </template>
            </el-table-column>
            <el-table-column label="状态" width="100">
              <template #default="{ row }">
                <el-switch
                  v-model="row.status"
                  :active-value="1"
                  :inactive-value="0"
                  @change="handlePromotionStatusChange(row)"
                />
              </template>
            </el-table-column>
            <el-table-column label="操作" width="200" fixed="right">
              <template #default="{ row }">
                <el-button type="primary" size="small" @click="handleEditPromotion(row)">编辑</el-button>
                <el-button type="danger" size="small" @click="handleDeletePromotion(row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>

          <el-pagination
            v-if="promotionsPagination.total > 0"
            v-model:current-page="promotionsPagination.page"
            v-model:page-size="promotionsPagination.page_size"
            :total="promotionsPagination.total"
            layout="total, prev, pager, next"
            @current-change="loadPromotions"
            class="pagination"
          />
        </el-card>
      </el-tab-pane>

      <!-- 充值订单 -->
      <el-tab-pane label="充值订单" name="orders">
        <el-card>
          <el-form :inline="true" :model="ordersSearch" class="search-form">
            <el-form-item label="订单号">
              <el-input v-model="ordersSearch.order_no" placeholder="请输入订单号" clearable />
            </el-form-item>
            <el-form-item label="手机号">
              <el-input v-model="ordersSearch.phone" placeholder="请输入手机号" clearable />
            </el-form-item>
            <el-form-item label="订单状态">
              <el-select v-model="ordersSearch.status" placeholder="请选择状态" clearable>
                <el-option label="全部" value="" />
                <el-option label="待支付" value="pending" />
                <el-option label="已支付" value="paid" />
                <el-option label="已取消" value="cancelled" />
              </el-select>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleOrdersSearch">搜索</el-button>
              <el-button @click="handleOrdersReset">重置</el-button>
            </el-form-item>
          </el-form>

          <el-table
            :data="ordersData"
            border
            stripe
            v-loading="ordersLoading"
            :empty-text="ordersLoading ? '加载中...' : '暂无数据'"
            style="width: 100%"
          >
            <el-table-column prop="order_no" label="订单号" width="180" />
            <el-table-column prop="user_phone" label="会员手机" width="120" />
            <el-table-column prop="user_nickname" label="会员昵称" width="120" />
            <el-table-column label="充值金额" width="120">
              <template #default="{ row }">
                <span class="money-text">¥{{ row.amount }}</span>
              </template>
            </el-table-column>
            <el-table-column label="实付金额" width="120">
              <template #default="{ row }">
                <span class="money-text">¥{{ row.pay_amount }}</span>
              </template>
            </el-table-column>
            <el-table-column label="赠送金额" width="120">
              <template #default="{ row }">
                <span class="bonus-text">¥{{ row.bonus_amount || 0 }}</span>
              </template>
            </el-table-column>
            <el-table-column label="赠送积分" width="100">
              <template #default="{ row }">
                <span class="points-text">{{ row.bonus_points || 0 }}</span>
              </template>
            </el-table-column>
            <el-table-column label="订单状态" width="100">
              <template #default="{ row }">
                <el-tag :type="getOrderStatusType(row.status)">
                  {{ getOrderStatusText(row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="created_at" label="创建时间" width="160" />
            <el-table-column prop="paid_at" label="支付时间" width="160" />
          </el-table>

          <el-pagination
            v-if="ordersPagination.total > 0"
            v-model:current-page="ordersPagination.page"
            v-model:page-size="ordersPagination.page_size"
            :total="ordersPagination.total"
            layout="total, prev, pager, next"
            @current-change="loadOrders"
            class="pagination"
          />
        </el-card>
      </el-tab-pane>
    </el-tabs>

    <!-- 新增/编辑活动弹窗 -->
    <el-dialog v-model="promotionVisible" :title="promotionForm.id ? '编辑活动' : '新增活动'" width="600px">
      <el-form :model="promotionForm" :rules="promotionRules" ref="promotionFormRef" label-width="120px">
        <el-form-item label="活动名称" prop="name">
          <el-input v-model="promotionForm.name" placeholder="请输入活动名称" />
        </el-form-item>
        <el-form-item label="活动描述" prop="description">
          <el-input v-model="promotionForm.description" type="textarea" :rows="3" placeholder="请输入活动描述" />
        </el-form-item>
        <el-form-item label="充值金额" prop="recharge_amount">
          <el-input-number v-model="promotionForm.recharge_amount" :min="1" :max="100000" :precision="2" />
          <span style="margin-left: 10px;">元</span>
        </el-form-item>
        <el-form-item label="赠送金额" prop="bonus_amount">
          <el-input-number v-model="promotionForm.bonus_amount" :min="0" :max="100000" :precision="2" />
          <span style="margin-left: 10px;">元</span>
        </el-form-item>
        <el-form-item label="赠送积分" prop="bonus_points">
          <el-input-number v-model="promotionForm.bonus_points" :min="0" :max="100000" />
        </el-form-item>
        <el-form-item label="活动时间" prop="time_range">
          <el-date-picker
            v-model="promotionForm.time_range"
            type="datetimerange"
            range-separator="-"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="promotionForm.status">
            <el-radio :label="1">启用</el-radio>
            <el-radio :label="0">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="promotionVisible = false">取消</el-button>
        <el-button type="primary" @click="handlePromotionSubmit" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import {
  getRechargePromotions,
  createRechargePromotion,
  updateRechargePromotion,
  deleteRechargePromotion,
  getRechargeOrders
} from '@/api/recharge'

const activeTab = ref('promotions')
const promotionsLoading = ref(false)
const ordersLoading = ref(false)
const promotionsData = ref([])
const ordersData = ref([])
const promotionVisible = ref(false)
const submitting = ref(false)
const promotionFormRef = ref(null)

const promotionsPagination = reactive({
  page: 1,
  page_size: 10,
  total: 0
})

const ordersPagination = reactive({
  page: 1,
  page_size: 10,
  total: 0
})

const ordersSearch = reactive({
  order_no: '',
  phone: '',
  status: ''
})

const promotionForm = reactive({
  id: null,
  name: '',
  description: '',
  recharge_amount: 100,
  bonus_amount: 0,
  bonus_points: 0,
  time_range: [],
  status: 1
})

const promotionRules = {
  name: [{ required: true, message: '请输入活动名称', trigger: 'blur' }],
  recharge_amount: [{ required: true, message: '请输入充值金额', trigger: 'blur' }],
  time_range: [{ required: true, message: '请选择活动时间', trigger: 'change' }]
}

// 加载充值活动
const loadPromotions = async () => {
  promotionsLoading.value = true
  try {
    const data = await getRechargePromotions({
      page: promotionsPagination.page,
      page_size: promotionsPagination.page_size
    })
    promotionsData.value = data.list || []
    promotionsPagination.total = data.total || 0
  } catch (error) {
    ElMessage.error('加载充值活动失败')
  } finally {
    promotionsLoading.value = false
  }
}

// 加载充值订单
const loadOrders = async () => {
  ordersLoading.value = true
  try {
    const data = await getRechargeOrders({
      page: ordersPagination.page,
      page_size: ordersPagination.page_size,
      ...ordersSearch
    })
    ordersData.value = data.list || []
    ordersPagination.total = data.total || 0
  } catch (error) {
    ElMessage.error('加载充值订单失败')
  } finally {
    ordersLoading.value = false
  }
}

// 获取订单状态文本
const getOrderStatusText = (status) => {
  const map = {
    pending: '待支付',
    paid: '已支付',
    cancelled: '已取消'
  }
  return map[status] || status
}

// 获取订单状态类型
const getOrderStatusType = (status) => {
  const map = {
    pending: 'warning',
    paid: 'success',
    cancelled: 'info'
  }
  return map[status] || ''
}

// 新增活动
const handleAddPromotion = () => {
  Object.assign(promotionForm, {
    id: null,
    name: '',
    description: '',
    recharge_amount: 100,
    bonus_amount: 0,
    bonus_points: 0,
    time_range: [],
    status: 1
  })
  promotionVisible.value = true
}

// 编辑活动
const handleEditPromotion = (row) => {
  Object.assign(promotionForm, {
    ...row,
    time_range: [row.start_time, row.end_time]
  })
  promotionVisible.value = true
}

// 删除活动
const handleDeletePromotion = (row) => {
  ElMessageBox.confirm('确定要删除该活动吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await deleteRechargePromotion(row.id)
      ElMessage.success('删除成功')
      loadPromotions()
    } catch (error) {
      ElMessage.error(error.message || '删除失败')
    }
  })
}

// 状态切换
const handlePromotionStatusChange = async (row) => {
  try {
    await updateRechargePromotion(row.id, { status: row.status })
    ElMessage.success('状态更新成功')
  } catch (error) {
    ElMessage.error('状态更新失败')
    row.status = row.status === 1 ? 0 : 1
  }
}

// 提交活动
const handlePromotionSubmit = async () => {
  if (!promotionFormRef.value) return

  await promotionFormRef.value.validate(async (valid) => {
    if (valid) {
      submitting.value = true
      try {
        const submitData = {
          ...promotionForm,
          start_time: promotionForm.time_range[0],
          end_time: promotionForm.time_range[1]
        }
        delete submitData.time_range

        if (promotionForm.id) {
          await updateRechargePromotion(promotionForm.id, submitData)
        } else {
          await createRechargePromotion(submitData)
        }
        ElMessage.success(promotionForm.id ? '更新成功' : '新增成功')
        promotionVisible.value = false
        loadPromotions()
      } catch (error) {
        ElMessage.error(error.message || '操作失败')
      } finally {
        submitting.value = false
      }
    }
  })
}

// 搜索订单
const handleOrdersSearch = () => {
  ordersPagination.page = 1
  loadOrders()
}

// 重置搜索
const handleOrdersReset = () => {
  ordersSearch.order_no = ''
  ordersSearch.phone = ''
  ordersSearch.status = ''
  ordersPagination.page = 1
  loadOrders()
}

onMounted(() => {
  loadPromotions()
  loadOrders()
})
</script>

<style scoped>
.recharge-container {
  padding: 20px;
  min-height: 100%;
  box-sizing: border-box;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.search-form {
  margin-bottom: 20px;
}

.money-text {
  color: #67c23a;
  font-weight: bold;
}

.bonus-text {
  color: #e6a23c;
  font-weight: bold;
}

.points-text {
  color: #f56c6c;
  font-weight: bold;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>
