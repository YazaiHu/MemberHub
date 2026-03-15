<template>
  <div class="promotions-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>特价商品列表</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            新增特价商品
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="商品名称">
          <el-input v-model="searchForm.name" placeholder="请输入商品名称" clearable />
        </el-form-item>
        <el-form-item label="门店">
          <el-select v-model="searchForm.store_id" placeholder="请选择门店" clearable>
            <el-option label="全部" value="" />
            <el-option
              v-for="store in storesList"
              :key="store.id"
              :label="store.name"
              :value="store.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="请选择状态" clearable>
            <el-option label="全部" value="" />
            <el-option label="进行中" value="active" />
            <el-option label="未开始" value="pending" />
            <el-option label="已结束" value="ended" />
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
        <el-table-column label="商品图片" width="100">
          <template #default="{ row }">
            <el-image
              style="width: 60px; height: 60px"
              :src="row.image || '/default-product.png'"
              fit="cover"
            />
          </template>
        </el-table-column>
        <el-table-column prop="product_name" label="商品名称" width="200" />
        <el-table-column prop="store_name" label="所属门店" width="150" />
        <el-table-column label="价格信息" width="200">
          <template #default="{ row }">
            <div>
              <span style="text-decoration: line-through; color: #999;">
                原价：¥{{ row.original_price }}
              </span>
            </div>
            <div>
              <span style="color: #f56c6c; font-weight: bold; font-size: 16px;">
                特价：¥{{ row.promotion_price }}
              </span>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="库存" width="100">
          <template #default="{ row }">
            <el-tag :type="row.stock > 0 ? 'success' : 'danger'">
              {{ row.stock }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="促销时间" width="300">
          <template #default="{ row }">
            {{ row.start_time }} ~ {{ row.end_time }}
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button type="danger" size="small" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <el-pagination
        v-if="pagination.total > 0"
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.page_size"
        :page-sizes="[10, 20, 50]"
        :total="pagination.total"
        layout="total, sizes, prev, pager, next"
        @size-change="handleSizeChange"
        @current-change="handlePageChange"
        class="pagination"
      />
    </el-card>

    <!-- 新增/编辑弹窗 -->
    <el-dialog v-model="dialogVisible" :title="form.id ? '编辑特价商品' : '新增特价商品'" width="700px">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="120px">
        <el-form-item label="所属门店" prop="store_id">
          <el-select v-model="form.store_id" placeholder="请选择门店" style="width: 100%">
            <el-option
              v-for="store in storesList"
              :key="store.id"
              :label="store.name"
              :value="store.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="商品名称" prop="product_name">
          <el-input v-model="form.product_name" placeholder="请输入商品名称" />
        </el-form-item>
        <el-form-item label="商品描述" prop="description">
          <el-input v-model="form.description" type="textarea" :rows="3" placeholder="请输入商品描述" />
        </el-form-item>
        <el-form-item label="商品图片" prop="image">
          <el-upload
            class="avatar-uploader"
            action="#"
            :show-file-list="false"
            :auto-upload="false"
          >
            <img v-if="form.image" :src="form.image" class="avatar" />
            <el-icon v-else class="avatar-uploader-icon"><Plus /></el-icon>
          </el-upload>
          <div style="margin-top: 5px; color: #999; font-size: 12px;">支持jpg、png格式，建议尺寸800x800</div>
        </el-form-item>
        <el-form-item label="原价" prop="original_price">
          <el-input-number v-model="form.original_price" :min="0" :precision="2" :step="1" />
          <span style="margin-left: 10px;">元</span>
        </el-form-item>
        <el-form-item label="特价" prop="promotion_price">
          <el-input-number v-model="form.promotion_price" :min="0" :precision="2" :step="1" />
          <span style="margin-left: 10px;">元</span>
        </el-form-item>
        <el-form-item label="库存数量" prop="stock">
          <el-input-number v-model="form.stock" :min="0" :max="100000" />
        </el-form-item>
        <el-form-item label="促销时间" prop="time_range">
          <el-date-picker
            v-model="form.time_range"
            type="datetimerange"
            range-separator="-"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="排序" prop="sort_order">
          <el-input-number v-model="form.sort_order" :min="0" />
          <span style="margin-left: 10px;">数字越大越靠前</span>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Search, Refresh } from '@element-plus/icons-vue'
import { getStoreList } from '@/api/store'
import {
  getPromotionProducts,
  createPromotionProduct,
  updatePromotionProduct,
  deletePromotionProduct
} from '@/api/promotion'

const loading = ref(false)
const tableData = ref([])
const dialogVisible = ref(false)
const submitting = ref(false)
const formRef = ref(null)
const storesList = ref([])

const searchForm = reactive({
  name: '',
  store_id: '',
  status: ''
})

const pagination = reactive({
  page: 1,
  page_size: 10,
  total: 0
})

const form = reactive({
  id: null,
  store_id: null,
  product_name: '',
  description: '',
  image: '',
  original_price: 0,
  promotion_price: 0,
  stock: 100,
  time_range: [],
  sort_order: 0
})

const rules = {
  store_id: [{ required: true, message: '请选择门店', trigger: 'change' }],
  product_name: [{ required: true, message: '请输入商品名称', trigger: 'blur' }],
  original_price: [{ required: true, message: '请输入原价', trigger: 'blur' }],
  promotion_price: [{ required: true, message: '请输入特价', trigger: 'blur' }],
  stock: [{ required: true, message: '请输入库存数量', trigger: 'blur' }],
  time_range: [{ required: true, message: '请选择促销时间', trigger: 'change' }]
}

// 获取状态文本
const getStatusText = (status) => {
  const map = {
    pending: '未开始',
    active: '进行中',
    ended: '已结束'
  }
  return map[status] || status
}

// 获取状态类型
const getStatusType = (status) => {
  const map = {
    pending: 'info',
    active: 'success',
    ended: 'danger'
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

// 加载特价商品列表
const loadPromotions = async () => {
  loading.value = true
  try {
    const data = await getPromotionProducts({
      page: pagination.page,
      page_size: pagination.page_size,
      ...searchForm
    })
    tableData.value = data.list || []
    pagination.total = data.total || 0
  } catch (error) {
    ElMessage.error('加载特价商品列表失败')
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  loadPromotions()
}

// 重置
const handleReset = () => {
  searchForm.name = ''
  searchForm.store_id = ''
  searchForm.status = ''
  pagination.page = 1
  loadPromotions()
}

// 分页大小改变
const handleSizeChange = (size) => {
  pagination.page_size = size
  loadPromotions()
}

// 页码改变
const handlePageChange = (page) => {
  pagination.page = page
  loadPromotions()
}

// 新增特价商品
const handleAdd = () => {
  Object.assign(form, {
    id: null,
    store_id: null,
    product_name: '',
    description: '',
    image: '',
    original_price: 0,
    promotion_price: 0,
    stock: 100,
    time_range: [],
    sort_order: 0
  })
  dialogVisible.value = true
}

// 编辑特价商品
const handleEdit = (row) => {
  Object.assign(form, {
    ...row,
    time_range: [row.start_time, row.end_time]
  })
  dialogVisible.value = true
}

// 删除特价商品
const handleDelete = (row) => {
  ElMessageBox.confirm('确定要删除该特价商品吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await deletePromotionProduct(row.id)
      ElMessage.success('删除成功')
      loadPromotions()
    } catch (error) {
      ElMessage.error(error.message || '删除失败')
    }
  })
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (valid) {
      submitting.value = true
      try {
        const submitData = {
          ...form,
          start_time: form.time_range[0],
          end_time: form.time_range[1]
        }
        delete submitData.time_range

        if (form.id) {
          await updatePromotionProduct(submitData)
        } else {
          await createPromotionProduct(submitData)
        }
        ElMessage.success(form.id ? '更新成功' : '新增成功')
        dialogVisible.value = false
        loadPromotions()
      } catch (error) {
        ElMessage.error(error.message || '操作失败')
      } finally {
        submitting.value = false
      }
    }
  })
}

onMounted(() => {
  loadStores()
  loadPromotions()
})
</script>

<style scoped>
.promotions-container {
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

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.avatar-uploader {
  border: 1px dashed #d9d9d9;
  border-radius: 6px;
  cursor: pointer;
  position: relative;
  overflow: hidden;
  transition: var(--el-transition-duration-fast);
}

.avatar-uploader:hover {
  border-color: var(--el-color-primary);
}

.avatar-uploader-icon {
  font-size: 28px;
  color: #8c939d;
  width: 148px;
  height: 148px;
  text-align: center;
  line-height: 148px;
}

.avatar {
  width: 148px;
  height: 148px;
  display: block;
}
</style>
