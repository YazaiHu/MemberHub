<template>
  <div class="points-container">
    <el-tabs v-model="activeTab">
      <!-- 兑换规则管理 -->
      <el-tab-pane label="兑换规则管理" name="rules">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>兑换规则列表</span>
              <el-button type="primary" @click="handleAddRule">
                <el-icon><Plus /></el-icon>
                新增规则
              </el-button>
            </div>
          </template>

          <el-table
            :data="rulesData"
            border
            stripe
            v-loading="rulesLoading"
            :empty-text="rulesLoading ? '加载中...' : '暂无数据'"
            style="width: 100%"
          >
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="store_name" label="门店" width="150" />
            <el-table-column prop="product_name" label="商品名称" width="200" />
            <el-table-column label="兑换积分" width="120">
              <template #default="{ row }">
                <span class="points-text">{{ row.points_required }}</span>
              </template>
            </el-table-column>
            <el-table-column label="库存" width="100">
              <template #default="{ row }">
                <el-tag :type="row.stock > 0 ? 'success' : 'danger'">
                  {{ row.stock }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="daily_limit" label="每日限兑" width="100" />
            <el-table-column label="状态" width="100">
              <template #default="{ row }">
                <el-switch
                  v-model="row.status"
                  :active-value="1"
                  :inactive-value="0"
                  @change="handleStatusChange(row)"
                />
              </template>
            </el-table-column>
            <el-table-column prop="created_at" label="创建时间" width="160" />
            <el-table-column label="操作" width="200" fixed="right">
              <template #default="{ row }">
                <el-button type="primary" size="small" @click="handleEditRule(row)">编辑</el-button>
                <el-button type="danger" size="small" @click="handleDeleteRule(row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>

          <el-pagination
            v-if="rulesPagination.total > 0"
            v-model:current-page="rulesPagination.page"
            v-model:page-size="rulesPagination.page_size"
            :total="rulesPagination.total"
            layout="total, prev, pager, next"
            @current-change="loadRules"
            class="pagination"
          />
        </el-card>
      </el-tab-pane>

      <!-- 兑换记录 -->
      <el-tab-pane label="兑换记录" name="records">
        <el-card>
          <el-form :inline="true" :model="recordsSearch" class="search-form">
            <el-form-item label="手机号">
              <el-input v-model="recordsSearch.phone" placeholder="请输入手机号" clearable />
            </el-form-item>
            <el-form-item label="状态">
              <el-select v-model="recordsSearch.status" placeholder="请选择状态" clearable>
                <el-option label="全部" value="" />
                <el-option label="待核销" value="pending" />
                <el-option label="已核销" value="used" />
                <el-option label="已过期" value="expired" />
              </el-select>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleRecordsSearch">搜索</el-button>
              <el-button @click="handleRecordsReset">重置</el-button>
            </el-form-item>
          </el-form>

          <el-table
            :data="recordsData"
            border
            stripe
            v-loading="recordsLoading"
            :empty-text="recordsLoading ? '加载中...' : '暂无数据'"
            style="width: 100%"
          >
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="user_phone" label="会员手机" width="120" />
            <el-table-column prop="user_nickname" label="会员昵称" width="120" />
            <el-table-column prop="product_name" label="兑换商品" width="200" />
            <el-table-column label="消耗积分" width="120">
              <template #default="{ row }">
                <span class="points-text">{{ row.points_used }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="exchange_code" label="兑换码" width="150" />
            <el-table-column label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="getRecordStatusType(row.status)">
                  {{ getRecordStatusText(row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="exchange_time" label="兑换时间" width="160" />
            <el-table-column prop="use_time" label="核销时间" width="160" />
            <el-table-column label="操作" width="120" fixed="right">
              <template #default="{ row }">
                <el-button
                  v-if="row.status === 'pending'"
                  type="success"
                  size="small"
                  @click="handleVerify(row)"
                >
                  核销
                </el-button>
              </template>
            </el-table-column>
          </el-table>

          <el-pagination
            v-if="recordsPagination.total > 0"
            v-model:current-page="recordsPagination.page"
            v-model:page-size="recordsPagination.page_size"
            :total="recordsPagination.total"
            layout="total, prev, pager, next"
            @current-change="loadRecords"
            class="pagination"
          />
        </el-card>
      </el-tab-pane>
    </el-tabs>

    <!-- 新增/编辑规则弹窗 -->
    <el-dialog v-model="ruleVisible" :title="ruleForm.id ? '编辑规则' : '新增规则'" width="600px">
      <el-form :model="ruleForm" :rules="ruleRules" ref="ruleFormRef" label-width="120px">
        <el-form-item label="门店" prop="store_id">
          <el-select v-model="ruleForm.store_id" placeholder="请选择门店" style="width: 100%">
            <el-option
              v-for="store in storesList"
              :key="store.id"
              :label="store.name"
              :value="store.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="商品名称" prop="product_name">
          <el-input v-model="ruleForm.product_name" placeholder="请输入商品名称" />
        </el-form-item>
        <el-form-item label="商品描述" prop="description">
          <el-input v-model="ruleForm.description" type="textarea" :rows="3" placeholder="请输入商品描述" />
        </el-form-item>
        <el-form-item label="兑换积分" prop="points_required">
          <el-input-number v-model="ruleForm.points_required" :min="1" :max="100000" />
        </el-form-item>
        <el-form-item label="库存数量" prop="stock">
          <el-input-number v-model="ruleForm.stock" :min="0" :max="100000" />
        </el-form-item>
        <el-form-item label="每日限兑" prop="daily_limit">
          <el-input-number v-model="ruleForm.daily_limit" :min="0" :max="1000" />
          <span style="margin-left: 10px; color: #999;">0表示不限制</span>
        </el-form-item>
        <el-form-item label="有效期(天)" prop="valid_days">
          <el-input-number v-model="ruleForm.valid_days" :min="1" :max="365" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="ruleForm.status">
            <el-radio :label="1">启用</el-radio>
            <el-radio :label="0">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="ruleVisible = false">取消</el-button>
        <el-button type="primary" @click="handleRuleSubmit" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { getStoreList } from '@/api/store'
import {
  getExchangeRules,
  createExchangeRule,
  updateExchangeRule,
  getExchangeRecords,
  verifyExchangeCode
} from '@/api/points'

const activeTab = ref('rules')
const rulesLoading = ref(false)
const recordsLoading = ref(false)
const rulesData = ref([])
const recordsData = ref([])
const ruleVisible = ref(false)
const submitting = ref(false)
const ruleFormRef = ref(null)
const storesList = ref([])

const rulesPagination = reactive({
  page: 1,
  page_size: 10,
  total: 0
})

const recordsPagination = reactive({
  page: 1,
  page_size: 10,
  total: 0
})

const recordsSearch = reactive({
  phone: '',
  status: ''
})

const ruleForm = reactive({
  id: null,
  store_id: null,
  product_name: '',
  description: '',
  points_required: 100,
  stock: 100,
  daily_limit: 0,
  valid_days: 30,
  status: 1
})

const ruleRules = {
  store_id: [{ required: true, message: '请选择门店', trigger: 'change' }],
  product_name: [{ required: true, message: '请输入商品名称', trigger: 'blur' }],
  points_required: [{ required: true, message: '请输入兑换积分', trigger: 'blur' }],
  stock: [{ required: true, message: '请输入库存数量', trigger: 'blur' }]
}

// 加载门店列表
const loadStores = async () => {
  try {
    const data = await getStoreList({ page: 1, page_size: 100 })
    storesList.value = data.list || []
  } catch (error) {
    console.error('加载门店列表失败', error)
  }
}

// 加载兑换规则
const loadRules = async () => {
  rulesLoading.value = true
  try {
    const data = await getExchangeRules({
      page: rulesPagination.page,
      page_size: rulesPagination.page_size
    })
    rulesData.value = data.list || []
    rulesPagination.total = data.total || 0
  } catch (error) {
    ElMessage.error('加载兑换规则失败')
  } finally {
    rulesLoading.value = false
  }
}

// 加载兑换记录
const loadRecords = async () => {
  recordsLoading.value = true
  try {
    const data = await getExchangeRecords({
      page: recordsPagination.page,
      page_size: recordsPagination.page_size,
      ...recordsSearch
    })
    recordsData.value = data.list || []
    recordsPagination.total = data.total || 0
  } catch (error) {
    ElMessage.error('加载兑换记录失败')
  } finally {
    recordsLoading.value = false
  }
}

// 获取记录状态文本
const getRecordStatusText = (status) => {
  const map = {
    pending: '待核销',
    used: '已核销',
    expired: '已过期'
  }
  return map[status] || status
}

// 获取记录状态类型
const getRecordStatusType = (status) => {
  const map = {
    pending: 'warning',
    used: 'success',
    expired: 'info'
  }
  return map[status] || ''
}

// 新增规则
const handleAddRule = () => {
  Object.assign(ruleForm, {
    id: null,
    store_id: null,
    product_name: '',
    description: '',
    points_required: 100,
    stock: 100,
    daily_limit: 0,
    valid_days: 30,
    status: 1
  })
  ruleVisible.value = true
}

// 编辑规则
const handleEditRule = (row) => {
  Object.assign(ruleForm, row)
  ruleVisible.value = true
}

// 删除规则
const handleDeleteRule = (row) => {
  ElMessageBox.confirm('确定要删除该规则吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      // TODO: 调用API
      ElMessage.success('删除成功')
      loadRules()
    } catch (error) {
      ElMessage.error('删除失败')
    }
  })
}

// 状态切换
const handleStatusChange = async (row) => {
  try {
    // TODO: 调用API
    ElMessage.success('状态更新成功')
  } catch (error) {
    ElMessage.error('状态更新失败')
    row.status = row.status === 1 ? 0 : 1
  }
}

// 提交规则
const handleRuleSubmit = async () => {
  if (!ruleFormRef.value) return

  await ruleFormRef.value.validate(async (valid) => {
    if (valid) {
      submitting.value = true
      try {
        // TODO: 调用API
        ElMessage.success(ruleForm.id ? '更新成功' : '新增成功')
        ruleVisible.value = false
        loadRules()
      } catch (error) {
        ElMessage.error(error.message || '操作失败')
      } finally {
        submitting.value = false
      }
    }
  })
}

// 搜索记录
const handleRecordsSearch = () => {
  recordsPagination.page = 1
  loadRecords()
}

// 重置搜索
const handleRecordsReset = () => {
  recordsSearch.phone = ''
  recordsSearch.status = ''
  recordsPagination.page = 1
  loadRecords()
}

// 核销
const handleVerify = (row) => {
  ElMessageBox.prompt('请输入核销码进行核销', '核销', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    inputPattern: /^EX\d{11}$/,
    inputErrorMessage: '核销码格式不正确'
  }).then(async ({ value }) => {
    if (value === row.exchange_code) {
      try {
        // TODO: 调用API
        ElMessage.success('核销成功')
        loadRecords()
      } catch (error) {
        ElMessage.error('核销失败')
      }
    } else {
      ElMessage.error('核销码不匹配')
    }
  })
}

onMounted(() => {
  loadStores()
  loadRules()
  loadRecords()
})
</script>

<style scoped>
.points-container {
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
