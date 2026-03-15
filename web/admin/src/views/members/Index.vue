<template>
  <div class="members-container">
    <el-card>
      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="手机号">
          <el-input v-model="searchForm.phone" placeholder="请输入手机号" clearable />
        </el-form-item>
        <el-form-item label="昵称">
          <el-input v-model="searchForm.nickname" placeholder="请输入昵称" clearable />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="请选择状态" clearable>
            <el-option label="全部" value="" />
            <el-option label="正常" :value="1" />
            <el-option label="禁用" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">
            <el-icon><Search /></el-icon>
            搜索
          </el-button>
          <el-button @click="handleReset">
            <el-icon><Refresh /></el-icon>
            重置
          </el-button>
        </el-form-item>
      </el-form>

      <!-- 数据表格 -->
      <el-table
        :data="tableData"
        border
        stripe
        v-loading="loading"
        :empty-text="loading ? '加载中...' : '暂无数据'"
        style="width: 100%"
      >
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="头像" width="80">
          <template #default="{ row }">
            <el-avatar :src="row.avatar || '/default-avatar.png'" />
          </template>
        </el-table-column>
        <el-table-column prop="nickname" label="昵称" width="120" />
        <el-table-column prop="phone" label="手机号" width="120" />
        <el-table-column label="会员等级" width="100">
          <template #default="{ row }">
            <el-tag :type="getLevelType(row.level)">{{ getLevelName(row.level) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="积分" width="100">
          <template #default="{ row }">
            <span class="points-text">{{ row.available_points || 0 }}</span>
          </template>
        </el-table-column>
        <el-table-column label="余额" width="100">
          <template #default="{ row }">
            <span class="balance-text">¥{{ (row.balance || 0).toFixed(2) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="优惠券" width="100">
          <template #default="{ row }">
            {{ row.coupon_count || 0 }} 张
          </template>
        </el-table-column>
        <el-table-column label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'">
              {{ row.status === 1 ? '正常' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="注册时间" width="160" />
        <el-table-column label="操作" width="280" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" size="small" @click="handleDetail(row)">
              <el-icon><View /></el-icon>
              详情
            </el-button>
            <el-button type="warning" size="small" @click="handleAdjustPoints(row)">
              <el-icon><Medal /></el-icon>
              调整积分
            </el-button>
            <el-button type="success" size="small" @click="handleAdjustBalance(row)">
              <el-icon><Wallet /></el-icon>
              调整余额
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <el-pagination
        v-if="pagination.total > 0"
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.page_size"
        :page-sizes="[10, 20, 50, 100]"
        :total="pagination.total"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handlePageChange"
        class="pagination"
      />
    </el-card>

    <!-- 会员详情弹窗 -->
    <el-dialog v-model="detailVisible" title="会员详情" width="800px">
      <el-descriptions :column="2" border v-if="currentMember">
        <el-descriptions-item label="用户ID">{{ currentMember.id }}</el-descriptions-item>
        <el-descriptions-item label="OpenID">{{ currentMember.openid }}</el-descriptions-item>
        <el-descriptions-item label="昵称">{{ currentMember.nickname }}</el-descriptions-item>
        <el-descriptions-item label="手机号">{{ currentMember.phone || '未绑定' }}</el-descriptions-item>
        <el-descriptions-item label="会员等级">{{ getLevelName(currentMember.level) }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="currentMember.status === 1 ? 'success' : 'danger'">
            {{ currentMember.status === 1 ? '正常' : '禁用' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="可用积分">{{ currentMember.available_points || 0 }}</el-descriptions-item>
        <el-descriptions-item label="累计积分">{{ currentMember.total_points || 0 }}</el-descriptions-item>
        <el-descriptions-item label="账户余额">¥{{ (currentMember.balance || 0).toFixed(2) }}</el-descriptions-item>
        <el-descriptions-item label="累计充值">¥{{ (currentMember.total_recharge || 0).toFixed(2) }}</el-descriptions-item>
        <el-descriptions-item label="优惠券数量">{{ currentMember.coupon_count || 0 }} 张</el-descriptions-item>
        <el-descriptions-item label="注册时间">{{ currentMember.created_at }}</el-descriptions-item>
      </el-descriptions>
    </el-dialog>

    <!-- 调整积分弹窗 -->
    <el-dialog v-model="pointsVisible" title="调整积分" width="500px">
      <el-form :model="pointsForm" :rules="pointsRules" ref="pointsFormRef" label-width="100px">
        <el-form-item label="会员昵称">
          <el-input v-model="currentMember.nickname" disabled />
        </el-form-item>
        <el-form-item label="当前积分">
          <el-input v-model="currentMember.available_points" disabled />
        </el-form-item>
        <el-form-item label="调整类型" prop="type">
          <el-radio-group v-model="pointsForm.type">
            <el-radio :label="1">增加</el-radio>
            <el-radio :label="-1">减少</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="调整数量" prop="amount">
          <el-input-number v-model="pointsForm.amount" :min="1" :max="100000" />
        </el-form-item>
        <el-form-item label="备注" prop="remark">
          <el-input v-model="pointsForm.remark" type="textarea" :rows="3" placeholder="请输入调整原因" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="pointsVisible = false">取消</el-button>
        <el-button type="primary" @click="handlePointsSubmit" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>

    <!-- 调整余额弹窗 -->
    <el-dialog v-model="balanceVisible" title="调整余额" width="500px">
      <el-form :model="balanceForm" :rules="balanceRules" ref="balanceFormRef" label-width="100px">
        <el-form-item label="会员昵称">
          <el-input v-model="currentMember.nickname" disabled />
        </el-form-item>
        <el-form-item label="当前余额">
          <el-input :value="'¥' + (currentMember.balance || 0).toFixed(2)" disabled />
        </el-form-item>
        <el-form-item label="调整类型" prop="type">
          <el-radio-group v-model="balanceForm.type">
            <el-radio :label="1">增加</el-radio>
            <el-radio :label="-1">减少</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="调整金额" prop="amount">
          <el-input-number v-model="balanceForm.amount" :min="0.01" :max="100000" :precision="2" :step="10" />
        </el-form-item>
        <el-form-item label="备注" prop="remark">
          <el-input v-model="balanceForm.remark" type="textarea" :rows="3" placeholder="请输入调整原因" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="balanceVisible = false">取消</el-button>
        <el-button type="primary" @click="handleBalanceSubmit" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Refresh, View, Medal, Wallet } from '@element-plus/icons-vue'
import { getMemberList, getMemberDetail, adjustPoints, adjustBalance } from '@/api/member'

const loading = ref(false)
const tableData = ref([])
const detailVisible = ref(false)
const pointsVisible = ref(false)
const balanceVisible = ref(false)
const submitting = ref(false)
const currentMember = ref({})
const pointsFormRef = ref(null)
const balanceFormRef = ref(null)

const searchForm = reactive({
  phone: '',
  nickname: '',
  status: ''
})

const pagination = reactive({
  page: 1,
  page_size: 10,
  total: 0
})

const pointsForm = reactive({
  type: 1,
  amount: 0,
  remark: ''
})

const balanceForm = reactive({
  type: 1,
  amount: 0,
  remark: ''
})

const pointsRules = {
  amount: [{ required: true, message: '请输入调整数量', trigger: 'blur' }],
  remark: [{ required: true, message: '请输入调整原因', trigger: 'blur' }]
}

const balanceRules = {
  amount: [{ required: true, message: '请输入调整金额', trigger: 'blur' }],
  remark: [{ required: true, message: '请输入调整原因', trigger: 'blur' }]
}

// 获取会员等级名称
const getLevelName = (level) => {
  const levels = { 1: '普通会员', 2: '白银会员', 3: '黄金会员', 4: 'VIP会员' }
  return levels[level] || '普通会员'
}

// 获取会员等级标签类型
const getLevelType = (level) => {
  const types = { 1: '', 2: 'info', 3: 'warning', 4: 'danger' }
  return types[level] || ''
}

// 加载会员列表
const loadMembers = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.page_size,
      ...searchForm
    }
    const data = await getMemberList(params)
    // 处理返回的数据，合并资产信息到会员列表
    tableData.value = (data.list || []).map(member => ({
      ...member,
      available_points: member.points?.available_points || 0,
      balance: member.balance?.balance || 0,
      coupon_count: member.coupons_count || 0
    }))
    pagination.total = data.total || 0
  } catch (error) {
    ElMessage.error('加载会员列表失败')
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  loadMembers()
}

// 重置
const handleReset = () => {
  searchForm.phone = ''
  searchForm.nickname = ''
  searchForm.status = ''
  pagination.page = 1
  loadMembers()
}

// 分页大小改变
const handleSizeChange = (size) => {
  pagination.page_size = size
  loadMembers()
}

// 页码改变
const handlePageChange = (page) => {
  pagination.page = page
  loadMembers()
}

// 查看详情
const handleDetail = async (row) => {
  try {
    const data = await getMemberDetail(row.id)
    currentMember.value = data
    detailVisible.value = true
  } catch (error) {
    ElMessage.error('获取会员详情失败')
  }
}

// 调整积分
const handleAdjustPoints = (row) => {
  currentMember.value = row
  pointsForm.type = 1
  pointsForm.amount = 0
  pointsForm.remark = ''
  pointsVisible.value = true
}

// 提交积分调整
const handlePointsSubmit = async () => {
  if (!pointsFormRef.value) return

  await pointsFormRef.value.validate(async (valid) => {
    if (valid) {
      submitting.value = true
      try {
        const points = pointsForm.type * pointsForm.amount
        await adjustPoints(currentMember.value.id, {
          points,
          remark: pointsForm.remark
        })
        ElMessage.success('积分调整成功')
        pointsVisible.value = false
        loadMembers()
      } catch (error) {
        ElMessage.error(error.message || '积分调整失败')
      } finally {
        submitting.value = false
      }
    }
  })
}

// 调整余额
const handleAdjustBalance = (row) => {
  currentMember.value = row
  balanceForm.type = 1
  balanceForm.amount = 0
  balanceForm.remark = ''
  balanceVisible.value = true
}

// 提交余额调整
const handleBalanceSubmit = async () => {
  if (!balanceFormRef.value) return

  await balanceFormRef.value.validate(async (valid) => {
    if (valid) {
      submitting.value = true
      try {
        const amount = balanceForm.type * balanceForm.amount
        await adjustBalance(currentMember.value.id, {
          amount,
          remark: balanceForm.remark
        })
        ElMessage.success('余额调整成功')
        balanceVisible.value = false
        loadMembers()
      } catch (error) {
        ElMessage.error(error.message || '余额调整失败')
      } finally {
        submitting.value = false
      }
    }
  })
}

onMounted(() => {
  loadMembers()
})
</script>

<style scoped>
.members-container {
  padding: 20px;
  min-height: 100%;
  box-sizing: border-box;
}

.search-form {
  margin-bottom: 20px;
}

.points-text {
  color: #f56c6c;
  font-weight: bold;
}

.balance-text {
  color: #67c23a;
  font-weight: bold;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>
