<template>
  <div class="coupons-container">
    <el-tabs v-model="activeTab">
      <!-- 优惠券模板 -->
      <el-tab-pane label="优惠券模板" name="templates">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>模板列表</span>
              <el-button type="primary" @click="handleAddTemplate">
                <el-icon><Plus /></el-icon>
                新增模板
              </el-button>
            </div>
          </template>

          <el-table
            :data="templatesData"
            border
            stripe
            v-loading="templatesLoading"
            :empty-text="templatesLoading ? '加载中...' : '暂无数据'"
            style="width: 100%"
          >
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="name" label="优惠券名称" width="200" />
            <el-table-column label="类型" width="120">
              <template #default="{ row }">
                <el-tag :type="getTypeColor(row.type)">{{ getTypeName(row.type) }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column label="优惠内容" width="200">
              <template #default="{ row }">
                <span v-if="row.type === 'full_discount'">
                  满 ¥{{ row.min_amount }} 减 ¥{{ row.discount_amount }}
                </span>
                <span v-else-if="row.type === 'discount'">
                  {{ row.discount_rate }}折
                </span>
                <span v-else>
                  抵用 ¥{{ row.discount_amount }}
                </span>
              </template>
            </el-table-column>
            <el-table-column label="发放数量" width="120">
              <template #default="{ row }">
                <span>{{ row.issued_quantity || 0 }} / {{ row.total_quantity }}</span>
              </template>
            </el-table-column>
            <el-table-column label="使用限制" width="150">
              <template #default="{ row }">
                <div>每人限领：{{ row.per_user_limit || '不限' }}</div>
                <div>有效期：{{ row.valid_days }}天</div>
              </template>
            </el-table-column>
            <el-table-column label="状态" width="100">
              <template #default="{ row }">
                <el-switch
                  v-model="row.status"
                  :active-value="1"
                  :inactive-value="0"
                  @change="handleTemplateStatusChange(row)"
                />
              </template>
            </el-table-column>
            <el-table-column prop="created_at" label="创建时间" width="160" />
            <el-table-column label="操作" width="250" fixed="right">
              <template #default="{ row }">
                <el-button type="success" size="small" @click="handlePush(row)">
                  <el-icon><Promotion /></el-icon>
                  推送
                </el-button>
                <el-button type="primary" size="small" @click="handleEditTemplate(row)">编辑</el-button>
                <el-button type="danger" size="small" @click="handleDeleteTemplate(row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>

          <el-pagination
            v-if="templatesPagination.total > 0"
            v-model:current-page="templatesPagination.page"
            v-model:page-size="templatesPagination.page_size"
            :total="templatesPagination.total"
            layout="total, prev, pager, next"
            @current-change="loadTemplates"
            class="pagination"
          />
        </el-card>
      </el-tab-pane>

      <!-- 推送任务 -->
      <el-tab-pane label="推送任务" name="tasks">
        <el-card>
          <el-table
            :data="tasksData"
            border
            stripe
            v-loading="tasksLoading"
            :empty-text="tasksLoading ? '加载中...' : '暂无数据'"
            style="width: 100%"
          >
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="template_name" label="优惠券模板" width="200" />
            <el-table-column label="目标用户" width="150">
              <template #default="{ row }">
                <span v-if="row.target_type === 'all'">全部会员</span>
                <span v-else-if="row.target_type === 'level'">{{ getLevelName(row.target_value) }}</span>
                <span v-else>指定用户</span>
              </template>
            </el-table-column>
            <el-table-column label="执行进度" width="200">
              <template #default="{ row }">
                <el-progress
                  :percentage="getProgress(row)"
                  :status="row.status === 'completed' ? 'success' : ''"
                />
              </template>
            </el-table-column>
            <el-table-column label="发放数量" width="150">
              <template #default="{ row }">
                <span>{{ row.success_count }} / {{ row.total_count }}</span>
                <span v-if="row.failed_count > 0" style="color: #f56c6c; margin-left: 10px;">
                  (失败{{ row.failed_count }})
                </span>
              </template>
            </el-table-column>
            <el-table-column label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="getTaskStatusType(row.status)">
                  {{ getTaskStatusText(row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="created_at" label="创建时间" width="160" />
            <el-table-column prop="completed_at" label="完成时间" width="160" />
          </el-table>

          <el-pagination
            v-if="tasksPagination.total > 0"
            v-model:current-page="tasksPagination.page"
            v-model:page-size="tasksPagination.page_size"
            :total="tasksPagination.total"
            layout="total, prev, pager, next"
            @current-change="loadTasks"
            class="pagination"
          />
        </el-card>
      </el-tab-pane>
    </el-tabs>

    <!-- 新增/编辑模板弹窗 -->
    <el-dialog v-model="templateVisible" :title="templateForm.id ? '编辑模板' : '新增模板'" width="700px">
      <el-form :model="templateForm" :rules="templateRules" ref="templateFormRef" label-width="120px">
        <el-form-item label="优惠券名称" prop="name">
          <el-input v-model="templateForm.name" placeholder="请输入优惠券名称" />
        </el-form-item>
        <el-form-item label="优惠券类型" prop="type">
          <el-radio-group v-model="templateForm.type">
            <el-radio label="full_discount">满减券</el-radio>
            <el-radio label="discount">折扣券</el-radio>
            <el-radio label="voucher">代金券</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item v-if="templateForm.type === 'full_discount'" label="满减金额" prop="min_amount">
          <el-input-number v-model="templateForm.min_amount" :min="0" :precision="2" />
          <span style="margin: 0 10px;">满</span>
          <el-input-number v-model="templateForm.discount_amount" :min="0" :precision="2" />
          <span style="margin-left: 10px;">减</span>
        </el-form-item>
        <el-form-item v-if="templateForm.type === 'discount'" label="折扣率" prop="discount_rate">
          <el-input-number v-model="templateForm.discount_rate" :min="1" :max="9.9" :precision="1" :step="0.1" />
          <span style="margin-left: 10px;">折</span>
        </el-form-item>
        <el-form-item v-if="templateForm.type === 'voucher'" label="抵用金额" prop="discount_amount">
          <el-input-number v-model="templateForm.discount_amount" :min="0" :precision="2" />
          <span style="margin-left: 10px;">元</span>
        </el-form-item>
        <el-form-item label="发放总量" prop="total_quantity">
          <el-input-number v-model="templateForm.total_quantity" :min="1" :max="1000000" />
        </el-form-item>
        <el-form-item label="每人限领" prop="per_user_limit">
          <el-input-number v-model="templateForm.per_user_limit" :min="0" :max="100" />
          <span style="margin-left: 10px;">0表示不限制</span>
        </el-form-item>
        <el-form-item label="有效期" prop="valid_days">
          <el-input-number v-model="templateForm.valid_days" :min="1" :max="365" />
          <span style="margin-left: 10px;">天</span>
        </el-form-item>
        <el-form-item label="适用门店" prop="applicable_stores">
          <el-select v-model="templateForm.applicable_stores" multiple placeholder="请选择适用门店" style="width: 100%">
            <el-option
              v-for="store in storesList"
              :key="store.id"
              :label="store.name"
              :value="store.id"
            />
          </el-select>
          <div style="margin-top: 5px; color: #999; font-size: 12px;">不选择则适用全部门店</div>
        </el-form-item>
        <el-form-item label="使用说明" prop="description">
          <el-input v-model="templateForm.description" type="textarea" :rows="3" placeholder="请输入使用说明" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="templateForm.status">
            <el-radio :label="1">启用</el-radio>
            <el-radio :label="0">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="templateVisible = false">取消</el-button>
        <el-button type="primary" @click="handleTemplateSubmit" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>

    <!-- 推送优惠券弹窗 -->
    <el-dialog v-model="pushVisible" title="推送优惠券" width="600px">
      <el-form :model="pushForm" :rules="pushRules" ref="pushFormRef" label-width="120px">
        <el-form-item label="优惠券模板">
          <el-input :value="currentTemplate.name" disabled />
        </el-form-item>
        <el-form-item label="目标用户" prop="target_type">
          <el-radio-group v-model="pushForm.target_type">
            <el-radio label="all">全部会员</el-radio>
            <el-radio label="level">指定等级</el-radio>
            <el-radio label="user_ids">指定用户</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item v-if="pushForm.target_type === 'level'" label="会员等级" prop="target_value">
          <el-select v-model="pushForm.target_value" placeholder="请选择会员等级" style="width: 100%">
            <el-option label="普通会员" :value="1" />
            <el-option label="白银会员" :value="2" />
            <el-option label="黄金会员" :value="3" />
            <el-option label="VIP会员" :value="4" />
          </el-select>
        </el-form-item>
        <el-form-item v-if="pushForm.target_type === 'user_ids'" label="用户ID" prop="user_ids">
          <el-input
            v-model="pushForm.user_ids"
            type="textarea"
            :rows="3"
            placeholder="请输入用户ID，多个用逗号分隔"
          />
        </el-form-item>
        <el-form-item label="推送消息">
          <el-checkbox v-model="pushForm.send_wechat_msg">发送微信模板消息</el-checkbox>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="pushVisible = false">取消</el-button>
        <el-button type="primary" @click="handlePushSubmit" :loading="submitting">确定推送</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Promotion } from '@element-plus/icons-vue'
import { getStoreList } from '@/api/store'
import {
  getCouponTemplates,
  createCouponTemplate,
  updateCouponTemplate,
  deleteCouponTemplate,
  createPushTask,
  getPushTasks
} from '@/api/coupon'

const activeTab = ref('templates')
const templatesLoading = ref(false)
const tasksLoading = ref(false)
const templatesData = ref([])
const tasksData = ref([])
const templateVisible = ref(false)
const pushVisible = ref(false)
const submitting = ref(false)
const templateFormRef = ref(null)
const pushFormRef = ref(null)
const storesList = ref([])
const currentTemplate = ref({})

const templatesPagination = reactive({
  page: 1,
  page_size: 10,
  total: 0
})

const tasksPagination = reactive({
  page: 1,
  page_size: 10,
  total: 0
})

const templateForm = reactive({
  id: null,
  name: '',
  type: 'full_discount',
  min_amount: 100,
  discount_amount: 10,
  discount_rate: 8.8,
  total_quantity: 1000,
  per_user_limit: 1,
  valid_days: 30,
  applicable_stores: [],
  description: '',
  status: 1
})

const pushForm = reactive({
  target_type: 'all',
  target_value: null,
  user_ids: '',
  send_wechat_msg: true
})

const templateRules = {
  name: [{ required: true, message: '请输入优惠券名称', trigger: 'blur' }],
  type: [{ required: true, message: '请选择优惠券类型', trigger: 'change' }]
}

const pushRules = {
  target_type: [{ required: true, message: '请选择目标用户', trigger: 'change' }]
}

// 获取类型名称
const getTypeName = (type) => {
  const map = {
    full_discount: '满减券',
    discount: '折扣券',
    voucher: '代金券'
  }
  return map[type] || type
}

// 获取类型颜色
const getTypeColor = (type) => {
  const map = {
    full_discount: 'success',
    discount: 'warning',
    voucher: 'danger'
  }
  return map[type] || ''
}

// 获取等级名称
const getLevelName = (level) => {
  const levels = { 1: '普通会员', 2: '白银会员', 3: '黄金会员', 4: 'VIP会员' }
  return levels[level] || '普通会员'
}

// 获取任务进度
const getProgress = (row) => {
  if (row.total_count === 0) return 0
  return Math.round((row.success_count / row.total_count) * 100)
}

// 获取任务状态文本
const getTaskStatusText = (status) => {
  const map = {
    pending: '待执行',
    processing: '执行中',
    completed: '已完成',
    failed: '失败'
  }
  return map[status] || status
}

// 获取任务状态类型
const getTaskStatusType = (status) => {
  const map = {
    pending: 'warning',
    processing: 'primary',
    completed: 'success',
    failed: 'danger'
  }
  return map[status] || ''
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

// 加载优惠券模板
const loadTemplates = async () => {
  templatesLoading.value = true
  try {
    const data = await getCouponTemplates({
      page: templatesPagination.page,
      page_size: templatesPagination.page_size
    })
    templatesData.value = data.list || []
    templatesPagination.total = data.total || 0
  } catch (error) {
    ElMessage.error('加载优惠券模板失败')
  } finally {
    templatesLoading.value = false
  }
}

// 加载推送任务
const loadTasks = async () => {
  tasksLoading.value = true
  try {
    const data = await getPushTasks({
      page: tasksPagination.page,
      page_size: tasksPagination.page_size
    })
    tasksData.value = data.list || []
    tasksPagination.total = data.total || 0
  } catch (error) {
    ElMessage.error('加载推送任务失败')
  } finally {
    tasksLoading.value = false
  }
}

// 新增模板
const handleAddTemplate = () => {
  Object.assign(templateForm, {
    id: null,
    name: '',
    type: 'full_discount',
    min_amount: 100,
    discount_amount: 10,
    discount_rate: 8.8,
    total_quantity: 1000,
    per_user_limit: 1,
    valid_days: 30,
    applicable_stores: [],
    description: '',
    status: 1
  })
  templateVisible.value = true
}

// 编辑模板
const handleEditTemplate = (row) => {
  Object.assign(templateForm, row)
  templateVisible.value = true
}

// 删除模板
const handleDeleteTemplate = (row) => {
  ElMessageBox.confirm('确定要删除该模板吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await deleteCouponTemplate(row.id)
      ElMessage.success('删除成功')
      loadTemplates()
    } catch (error) {
      ElMessage.error(error.message || '删除失败')
    }
  })
}

// 状态切换
const handleTemplateStatusChange = async (row) => {
  try {
    await updateCouponTemplate(row.id, { status: row.status })
    ElMessage.success('状态更新成功')
  } catch (error) {
    ElMessage.error('状态更新失败')
    row.status = row.status === 1 ? 0 : 1
  }
}

// 提交模板
const handleTemplateSubmit = async () => {
  if (!templateFormRef.value) return

  await templateFormRef.value.validate(async (valid) => {
    if (valid) {
      submitting.value = true
      try {
        if (templateForm.id) {
          await updateCouponTemplate(templateForm.id, templateForm)
        } else {
          await createCouponTemplate(templateForm)
        }
        ElMessage.success(templateForm.id ? '更新成功' : '新增成功')
        templateVisible.value = false
        loadTemplates()
      } catch (error) {
        ElMessage.error(error.message || '操作失败')
      } finally {
        submitting.value = false
      }
    }
  })
}

// 推送优惠券
const handlePush = (row) => {
  currentTemplate.value = row
  pushForm.target_type = 'all'
  pushForm.target_value = null
  pushForm.user_ids = ''
  pushForm.send_wechat_msg = true
  pushVisible.value = true
}

// 提交推送
const handlePushSubmit = async () => {
  if (!pushFormRef.value) return

  await pushFormRef.value.validate(async (valid) => {
    if (valid) {
      submitting.value = true
      try {
        const data = {
          template_id: currentTemplate.value.id,
          target_type: pushForm.target_type,
          send_wechat_msg: pushForm.send_wechat_msg
        }

        if (pushForm.target_type === 'level') {
          data.target_value = pushForm.target_value
        } else if (pushForm.target_type === 'user_ids') {
          data.user_ids = pushForm.user_ids.split(',').map(id => parseInt(id.trim()))
        }

        await createPushTask(data)
        ElMessage.success('推送任务创建成功，正在后台处理')
        pushVisible.value = false
        activeTab.value = 'tasks'
        loadTasks()
      } catch (error) {
        ElMessage.error(error.message || '推送失败')
      } finally {
        submitting.value = false
      }
    }
  })
}

onMounted(() => {
  loadStores()
  loadTemplates()
  loadTasks()
})
</script>

<style scoped>
.coupons-container {
  padding: 20px;
  min-height: 100%;
  box-sizing: border-box;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>
