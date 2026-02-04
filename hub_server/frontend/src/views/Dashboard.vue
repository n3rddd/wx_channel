<template>
  <div class="w-full min-h-screen bg-bg px-6 py-8 md:px-12 md:py-10">
    <!-- Page Header -->
    <div class="w-full mb-12 flex flex-col md:flex-row md:items-end justify-between gap-6">
      <button 
        class="group px-6 py-3 rounded-full bg-primary text-white font-bold text-sm tracking-wide shadow-neu-btn flex items-center gap-2 transition-all hover:scale-105 hover:shadow-lg active:scale-95" 
        @click="clientStore.fetchClients"
      >
        <RefreshCw class="w-4 h-4 transition-transform group-hover:rotate-180" :class="{ 'animate-spin': clientStore.loading }" />
        刷新状态
      </button>
    </div>

    <!-- Main Content Area -->
    <div class="w-full">
      
      <!-- Loading State -->
      <div v-if="clientStore.loading && !clientStore.clients.length" class="flex flex-col items-center justify-center py-40">
        <div class="w-16 h-16 border-4 border-primary/20 border-t-primary rounded-full animate-spin mb-6"></div>
        <p class="text-slate-400 font-medium text-lg">正在连接 Hub...</p>
      </div>

      <!-- Empty State -->
      <div v-else-if="clientStore.clients.length === 0" class="flex flex-col items-center justify-center py-40 bg-white rounded-3xl shadow-neu border border-slate-100">
        <div class="w-24 h-24 rounded-full bg-slate-50 flex items-center justify-center mb-6">
          <svg class="w-10 h-10 text-slate-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19.428 15.428a2 2 0 00-1.022-.547l-2.384-.477a6 6 0 00-3.86.517l-.318.158a6 6 0 01-3.86.517L6.05 15.21a2 2 0 00-1.806.547M8 4h8l-1 1v5.172a2 2 0 00.586 1.414l5 5c1.26 1.26.367 3.414-1.415 3.414H4.828c-1.782 0-2.674-2.154-1.414-3.414l5-5A2 2 0 009 10.172V5L8 4z"/>
          </svg>
        </div>
        <h3 class="text-2xl font-serif font-bold text-slate-800 mb-2">暂无在线终端</h3>
        <p class="text-slate-500 max-w-md text-center">请并在目标机器上启动客户端应用程序并配置 Hub URL。</p>
      </div>

      <!-- Terminal Grid -->
      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-3 gap-8">
        <div 
          v-for="client in clientStore.clients" 
          :key="client.id" 
          class="group bg-white rounded-[2rem] p-8 shadow-neu border border-slate-50 transition-all duration-500 hover:shadow-xl hover:-translate-y-1"
        >
          <!-- Card Header -->
          <div class="flex items-start justify-between mb-8">
            <div class="flex items-center gap-4">
              <!-- Icon Container -->
              <div class="w-14 h-14 rounded-2xl flex items-center justify-center transition-colors duration-300"
                   :class="client.status === 'online' 
                     ? 'bg-green-50 text-primary group-hover:bg-primary group-hover:text-white' 
                     : 'bg-gray-100 text-gray-400'">
                <svg class="w-7 h-7" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <rect x="2" y="3" width="20" height="14" rx="2" ry="2" stroke-width="2"/>
                  <line x1="8" y1="21" x2="16" y2="21" stroke-width="2"/>
                  <line x1="12" y1="17" x2="12" y2="21" stroke-width="2"/>
                </svg>
              </div>
              <div>
                <h3 class="font-bold text-xl text-slate-800 leading-tight mb-1">{{ client.hostname }}</h3>
                <div class="flex items-center gap-2">
                  <span class="relative flex h-2.5 w-2.5">
                    <span v-if="client.status === 'online'" class="animate-ping absolute inline-flex h-full w-full rounded-full bg-green-400 opacity-75"></span>
                    <span class="relative inline-flex rounded-full h-2.5 w-2.5" :class="client.status === 'online' ? 'bg-green-500' : 'bg-gray-400'"></span>
                  </span>
                  <span class="text-xs font-bold uppercase tracking-widest" :class="client.status === 'online' ? 'text-green-600' : 'text-gray-500'">
                    {{ client.status === 'online' ? '在线' : '离线' }}
                  </span>
                </div>
              </div>
            </div>
          </div>

          <!-- Card Metrics -->
          <div class="space-y-4 mb-8">
            <div class="flex items-center justify-between p-3 rounded-xl bg-slate-50 group-hover:bg-slate-100/80 transition-colors">
              <span class="text-xs font-bold text-slate-400 uppercase tracking-wider">版本 (Version)</span>
              <span class="text-sm font-semibold text-slate-700">v{{ client.version || '1.0.0' }}</span>
            </div>
            
            <div class="flex items-center justify-between p-3 rounded-xl bg-slate-50 group-hover:bg-slate-100/80 transition-colors">
              <span class="text-xs font-bold text-slate-400 uppercase tracking-wider">客户端 ID</span>
              <span class="text-xs font-mono font-medium text-slate-600 truncate max-w-[120px]" :title="client.id">
                {{ client.id }}
              </span>
            </div>

            <div class="flex items-center justify-between p-3 rounded-xl bg-slate-50 group-hover:bg-slate-100/80 transition-colors">
              <span class="text-xs font-bold text-slate-400 uppercase tracking-wider">最近心跳</span>
              <span class="text-sm font-semibold text-slate-700">{{ timeAgo(client.last_seen) }}</span>
            </div>
          </div>

          <!-- Actions -->
          <div class="grid grid-cols-2 gap-4">
            <button 
              class="py-3 rounded-xl border border-slate-200 text-slate-600 font-bold text-sm hover:bg-slate-50 hover:text-slate-900 transition-colors"
              @click="router.push('/nodes/' + client.id)"
            >
              详情
            </button>
            <button 
              class="py-3 rounded-xl font-bold text-sm transition-all"
              :class="client.status === 'online' 
                ? 'bg-slate-800 text-white hover:bg-slate-900 shadow-lg shadow-slate-200 hover:-translate-y-0.5 cursor-pointer' 
                : 'bg-gray-200 text-gray-400 cursor-not-allowed'"
              :disabled="client.status !== 'online'"
              @click="client.status === 'online' && enterConsole(client)"
            >
              控制台
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useClientStore } from '../store/client'
import { RefreshCw } from 'lucide-vue-next'
import { timeAgo } from '../utils/format'

const clientStore = useClientStore()
const router = useRouter()
let timer = null

onMounted(() => {
  clientStore.fetchClients()
  timer = setInterval(() => {
    clientStore.fetchClients()
  }, 5000)
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
})

const enterConsole = (client) => {
  clientStore.setCurrentClient(client.id)
  router.push('/search')
}
</script>
