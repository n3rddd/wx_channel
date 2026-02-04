<template>
  <div class="w-full space-y-8 p-8">
    <h1 class="text-3xl font-serif font-bold text-slate-800">设备管理</h1>
    
    <!-- Bind New Device -->
    <div class="bg-white rounded-3xl p-8 shadow-neu border border-slate-100">
      <div class="flex flex-col md:flex-row items-center justify-between gap-6">
        <div>
          <h2 class="text-xl font-bold text-slate-800 mb-2">添加新设备</h2>
          <p class="text-slate-500 text-sm">在您的客户端上运行以下命令以绑定此账号。验证码有效期为 5 分钟。</p>
        </div>
        <div class="flex flex-col items-end gap-3 w-full md:w-auto">
             <div v-if="bindToken" class="flex items-center gap-3 bg-slate-50 border border-slate-200 px-4 py-2 rounded-xl">
                 <span class="font-mono text-2xl font-bold text-primary tracking-widest">{{ bindToken }}</span>
                 <button @click="copyToken" class="p-2 rounded-lg hover:bg-slate-200 text-slate-500 hover:text-slate-700 transition-colors" title="复制">
                     <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" /></svg>
                 </button>
             </div>
             <button 
                v-else
                @click="generateToken" 
                class="px-6 py-3 rounded-xl bg-slate-900 text-white font-bold hover:bg-slate-800 shadow-lg shadow-slate-200 hover:-translate-y-0.5 transition-all"
             >
                生成绑定码
             </button>
             <p v-if="bindToken" class="text-xs text-slate-400">命令: client bind {{ bindToken }}</p>
        </div>
      </div>
    </div>

    <!-- Bound Devices List -->
    <div class="bg-white rounded-3xl p-8 shadow-neu border border-slate-100">
      <h2 class="text-xl font-bold text-slate-800 mb-6">已绑定设备</h2>
      
      <div v-if="loading" class="text-center py-10 text-slate-400">
          加载中...
      </div>
      
      <div v-else-if="devices.length === 0" class="text-center py-10 text-slate-400 bg-slate-50 rounded-2xl border border-dashed border-slate-200">
          暂无绑定设备
      </div>

      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div v-for="device in devices" :key="device.id" class="bg-white border border-slate-200 rounded-2xl p-5 hover:shadow-lg transition-all hover:border-primary/50 group">
          <div class="flex justify-between items-start mb-4">
              <div class="flex items-center gap-3">
                  <div class="w-10 h-10 rounded-xl bg-slate-100 flex items-center justify-center text-slate-500 group-hover:bg-primary/10 group-hover:text-primary transition-colors">
                      <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" /></svg>
                  </div>
                  <div>
                      <h3 class="font-bold text-slate-800 text-sm truncate max-w-[120px]" :title="device.hostname">{{ device.hostname || 'Unknown Host' }}</h3>
                      <p class="text-xs text-slate-400 font-mono">{{ device.id.substring(0, 8) }}...</p>
                  </div>
              </div>
              <span :class="device.status === 'online' ? 'bg-green-100 text-green-700' : 'bg-slate-100 text-slate-500'" class="px-2 py-1 rounded-md text-xs font-bold uppercase tracking-wide">
                  {{ device.status }}
              </span>
          </div>
          
          <div class="space-y-2 mb-4">
              <div class="flex justify-between text-xs">
                  <span class="text-slate-400">版本</span>
                  <span class="text-slate-700 font-mono">{{ device.version }}</span>
              </div>
              <div class="flex justify-between text-xs">
                  <span class="text-slate-400">最后在线</span>
                  <span class="text-slate-700">{{ formatDate(device.last_seen) }}</span>
              </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'

const bindToken = ref('')
const devices = ref([])
const loading = ref(false)

const generateToken = async () => {
    try {
        const res = await axios.post('/api/device/bind_token')
        bindToken.value = res.data.token
    } catch (err) {
        console.error("Generate token failed", err)
    }
}

const copyToken = () => {
    navigator.clipboard.writeText(bindToken.value)
    // Optional: Toast notification
}

const fetchDevices = async () => {
    loading.value = true
    try {
        const res = await axios.get('/api/device/list')
        devices.value = res.data || []
    } catch (err) {
        console.error("Fetch devices failed", err)
    } finally {
        loading.value = false
    }
}

const formatDate = (dateStr) => {
    if (!dateStr) return '-'
    return new Date(dateStr).toLocaleString()
}

onMounted(() => {
    fetchDevices()
})
</script>
