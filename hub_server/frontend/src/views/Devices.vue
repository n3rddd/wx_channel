<template>
  <div class="p-8">
    <!-- 添加新设备 -->
    <div class="bg-white shadow-card rounded-2xl p-6 mb-8 border border-slate-100">
      <div class="flex flex-col md:flex-row items-center justify-between gap-6">
        <div>
          <h2 class="text-xl font-bold text-slate-800 mb-2">添加新设备</h2>
          <p class="text-slate-500 text-sm">在您的客户端上运行以下命令以绑定此账号。验证码有效期为 5 分钟。</p>
        </div>
        <div class="flex flex-col items-end gap-3 w-full md:w-auto">
          <div v-if="bindToken" class="flex items-center gap-3 bg-bg shadow-neu-pressed px-4 py-3 rounded-xl">
            <span class="font-mono text-2xl font-bold text-primary tracking-widest">{{ bindToken }}</span>
            <button @click="copyToken" class="p-2 rounded-lg bg-bg shadow-neu hover:shadow-neu-sm text-slate-600 transition-all" title="复制">
              <component :is="Copy" class="w-5 h-5" />
            </button>
          </div>
          <button 
            v-else
            @click="generateToken" 
            class="px-6 py-3 rounded-xl bg-primary text-white font-bold hover:bg-primary/90 shadow-neu transition-all"
          >
            生成绑定码
          </button>
          <p v-if="bindToken" class="text-xs text-slate-400">命令: client bind {{ bindToken }}</p>
        </div>
      </div>
    </div>

    <!-- 统计卡片 -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
      <div class="bg-white shadow-card rounded-2xl p-6 border border-slate-100">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm text-slate-500 mb-1">总设备数</p>
            <p class="text-3xl font-bold text-slate-800">{{ devices.length }}</p>
          </div>
          <div class="w-12 h-12 bg-blue-100 rounded-xl flex items-center justify-center">
            <component :is="Monitor" class="w-6 h-6 text-blue-600" />
          </div>
        </div>
      </div>

      <div class="bg-white shadow-card rounded-2xl p-6 border border-slate-100">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm text-slate-500 mb-1">在线设备</p>
            <p class="text-3xl font-bold text-green-600">{{ onlineCount }}</p>
          </div>
          <div class="w-12 h-12 bg-green-100 rounded-xl flex items-center justify-center">
            <component :is="Wifi" class="w-6 h-6 text-green-600" />
          </div>
        </div>
      </div>

      <div class="bg-white shadow-card rounded-2xl p-6 border border-slate-100">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm text-slate-500 mb-1">离线设备</p>
            <p class="text-3xl font-bold text-slate-400">{{ offlineCount }}</p>
          </div>
          <div class="w-12 h-12 bg-slate-100 rounded-xl flex items-center justify-center">
            <component :is="WifiOff" class="w-6 h-6 text-slate-400" />
          </div>
        </div>
      </div>
    </div>

    <!-- 设备列表 -->
    <div class="bg-white shadow-card rounded-2xl p-6 border border-slate-100">
      <div class="flex items-center justify-between mb-6">
        <h2 class="text-xl font-bold text-slate-800">设备列表</h2>
        <button 
          @click="refreshDevices"
          class="px-4 py-2 bg-bg shadow-neu rounded-xl text-slate-600 hover:shadow-neu-sm transition-all flex items-center gap-2"
        >
          <component :is="RefreshCw" class="w-4 h-4" :class="{ 'animate-spin': loading }" />
          刷新
        </button>
      </div>

      <div v-if="loading && devices.length === 0" class="text-center py-12">
        <component :is="Loader2" class="w-8 h-8 text-slate-400 animate-spin mx-auto mb-4" />
        <p class="text-slate-500">加载中...</p>
      </div>

      <div v-else-if="devices.length === 0" class="text-center py-12">
        <component :is="Monitor" class="w-16 h-16 text-slate-300 mx-auto mb-4" />
        <p class="text-slate-500 mb-2">暂无设备</p>
        <p class="text-sm text-slate-400">请使用上方的绑定码绑定设备</p>
      </div>

      <div v-else class="space-y-4">
        <div 
          v-for="device in devices" 
          :key="device.id"
          class="bg-white shadow-card rounded-xl p-6 hover:shadow-lg transition-all"
        >
          <div class="flex items-start justify-between">
            <div class="flex items-start gap-4 flex-1">
              <!-- 状态指示器 -->
              <div class="mt-1">
                <div 
                  class="w-3 h-3 rounded-full"
                  :class="device.status === 'online' ? 'bg-green-500 animate-pulse' : 'bg-slate-300'"
                ></div>
              </div>

              <!-- 设备信息 -->
              <div class="flex-1">
                <div class="flex items-center gap-3 mb-2">
                  <h3 class="text-lg font-bold text-slate-800">{{ device.id }}</h3>
                  <span 
                    class="px-3 py-1 rounded-full text-xs font-medium"
                    :class="device.status === 'online' 
                      ? 'bg-green-100 text-green-700' 
                      : 'bg-slate-100 text-slate-500'"
                  >
                    {{ device.status === 'online' ? '在线' : '离线' }}
                  </span>
                </div>

                <div class="grid grid-cols-2 gap-4 text-sm">
                  <div>
                    <span class="text-slate-500">版本：</span>
                    <span class="text-slate-700 font-medium">{{ device.version || 'N/A' }}</span>
                  </div>
                  <div>
                    <span class="text-slate-500">主机名：</span>
                    <span class="text-slate-700 font-medium">{{ device.hostname || 'N/A' }}</span>
                  </div>
                  <div>
                    <span class="text-slate-500">最后在线：</span>
                    <span class="text-slate-700 font-medium">{{ formatTime(device.last_seen) }}</span>
                  </div>
                  <div>
                    <span class="text-slate-500">绑定时间：</span>
                    <span class="text-slate-700 font-medium">{{ formatTime(device.created_at) }}</span>
                  </div>
                </div>
              </div>
            </div>

            <!-- 操作按钮 -->
            <div class="flex gap-2 ml-4">
              <button
                @click="confirmUnbind(device)"
                class="px-4 py-2 bg-bg shadow-neu rounded-xl text-orange-600 hover:shadow-neu-sm transition-all flex items-center gap-2"
                title="解绑设备"
              >
                <component :is="Unlink" class="w-4 h-4" />
                解绑
              </button>
              <button
                @click="confirmDelete(device)"
                class="px-4 py-2 bg-bg shadow-neu rounded-xl text-red-600 hover:shadow-neu-sm transition-all flex items-center gap-2"
                title="删除设备"
              >
                <component :is="Trash2" class="w-4 h-4" />
                删除
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 确认对话框 -->
    <div 
      v-if="showConfirm"
      class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
      @click.self="showConfirm = false"
    >
      <div class="bg-white shadow-card rounded-2xl p-8 max-w-md w-full mx-4 border border-slate-100">
        <div class="flex items-center gap-4 mb-6">
          <div 
            class="w-12 h-12 rounded-xl flex items-center justify-center"
            :class="confirmAction === 'delete' ? 'bg-red-100' : 'bg-orange-100'"
          >
            <component 
              :is="confirmAction === 'delete' ? Trash2 : Unlink" 
              class="w-6 h-6"
              :class="confirmAction === 'delete' ? 'text-red-600' : 'text-orange-600'"
            />
          </div>
          <div>
            <h3 class="text-xl font-bold text-slate-800">
              {{ confirmAction === 'delete' ? '删除设备' : '解绑设备' }}
            </h3>
            <p class="text-sm text-slate-500">此操作需要确认</p>
          </div>
        </div>

        <div class="mb-6">
          <p class="text-slate-700 mb-4">
            {{ confirmAction === 'delete' 
              ? '确定要永久删除此设备吗？删除后无法恢复。' 
              : '确定要解绑此设备吗？解绑后设备将不再与您的账号关联。' 
            }}
          </p>
          <div class="bg-slate-50 rounded-xl p-4">
            <p class="text-sm text-slate-600">
              <span class="font-medium">设备 ID：</span>{{ selectedDevice?.id }}
            </p>
          </div>
        </div>

        <div class="flex gap-3">
          <button
            @click="showConfirm = false"
            class="flex-1 px-4 py-3 bg-bg shadow-neu rounded-xl text-slate-600 hover:shadow-neu-sm transition-all"
          >
            取消
          </button>
          <button
            @click="executeAction"
            :disabled="actionLoading"
            class="flex-1 px-4 py-3 rounded-xl text-white transition-all flex items-center justify-center gap-2"
            :class="confirmAction === 'delete' 
              ? 'bg-red-600 hover:bg-red-700' 
              : 'bg-orange-600 hover:bg-orange-700'"
          >
            <component 
              v-if="actionLoading"
              :is="Loader2" 
              class="w-4 h-4 animate-spin" 
            />
            <span>{{ actionLoading ? '处理中...' : '确认' }}</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { Monitor, Wifi, WifiOff, RefreshCw, Loader2, Unlink, Trash2, Copy } from 'lucide-vue-next'
import axios from 'axios'

const bindToken = ref('')
const devices = ref([])
const loading = ref(false)
const showConfirm = ref(false)
const confirmAction = ref('') // 'unbind' or 'delete'
const selectedDevice = ref(null)
const actionLoading = ref(false)

const onlineCount = computed(() => 
  devices.value.filter(d => d.status === 'online').length
)

const offlineCount = computed(() => 
  devices.value.filter(d => d.status !== 'online').length
)

const generateToken = async () => {
  try {
    const res = await axios.post('/api/device/bind_token')
    bindToken.value = res.data.token
  } catch (err) {
    console.error('Generate token failed', err)
    alert('生成绑定码失败')
  }
}

const copyToken = () => {
  navigator.clipboard.writeText(bindToken.value)
  // TODO: 添加复制成功提示
}

const formatTime = (time) => {
  if (!time) return 'N/A'
  const date = new Date(time)
  return date.toLocaleString('zh-CN')
}

const refreshDevices = async () => {
  loading.value = true
  try {
    const response = await axios.get('/api/device/list')
    devices.value = response.data || []
  } catch (error) {
    console.error('Failed to load devices:', error)
    alert('加载设备列表失败')
  } finally {
    loading.value = false
  }
}

const confirmUnbind = (device) => {
  selectedDevice.value = device
  confirmAction.value = 'unbind'
  showConfirm.value = true
}

const confirmDelete = (device) => {
  selectedDevice.value = device
  confirmAction.value = 'delete'
  showConfirm.value = true
}

const executeAction = async () => {
  if (!selectedDevice.value) return

  actionLoading.value = true
  try {
    const endpoint = confirmAction.value === 'delete' 
      ? '/api/device/delete' 
      : '/api/device/unbind'
    
    await axios.post(endpoint, {
      device_id: selectedDevice.value.id
    })

    // 刷新列表
    await refreshDevices()
    
    showConfirm.value = false
    selectedDevice.value = null
  } catch (error) {
    console.error('Action failed:', error)
    alert(error.response?.data || '操作失败')
  } finally {
    actionLoading.value = false
  }
}

onMounted(() => {
  refreshDevices()
})
</script>
