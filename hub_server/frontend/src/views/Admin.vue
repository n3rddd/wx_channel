<template>
  <div class="w-full space-y-8 p-8">
    <!-- Stats Cards -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
       <!-- Users -->
       <div class="bg-white rounded-3xl p-6 shadow-neu border border-slate-100 flex items-center justify-between">
           <div>
               <p class="text-slate-500 text-sm font-medium mb-1">总用户数</p>
               <h3 class="text-3xl font-bold text-slate-800">{{ stats.users || 0 }}</h3>
           </div>
           <div class="w-12 h-12 rounded-2xl bg-blue-50 text-blue-500 flex items-center justify-center">
               <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" /></svg>
           </div>
       </div>
       
       <!-- Devices -->
       <div class="bg-white rounded-3xl p-6 shadow-neu border border-slate-100 flex items-center justify-between">
           <div>
               <p class="text-slate-500 text-sm font-medium mb-1">活跃设备</p>
               <h3 class="text-3xl font-bold text-slate-800">{{ stats.devices || 0 }}</h3>
           </div>
           <div class="w-12 h-12 rounded-2xl bg-purple-50 text-purple-500 flex items-center justify-center">
               <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 18h.01M8 21h8a2 2 0 002-2V5a2 2 0 00-2-2H8a2 2 0 00-2 2v14a2 2 0 002 2z" /></svg>
           </div>
       </div>

       <!-- Transactions -->
       <div class="bg-white rounded-3xl p-6 shadow-neu border border-slate-100 flex items-center justify-between">
           <div>
               <p class="text-slate-500 text-sm font-medium mb-1">交易记录</p>
               <h3 class="text-3xl font-bold text-slate-800">{{ stats.transactions || 0 }}</h3>
           </div>
           <div class="w-12 h-12 rounded-2xl bg-green-50 text-green-500 flex items-center justify-center">
               <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 7h6m0 10v-3m-3 3h.01M9 17h.01M9 14h.01M12 14h.01M15 11h.01M12 11h.01M9 11h.01M7 21h10a2 2 0 002-2V5a2 2 0 00-2-2H7a2 2 0 00-2 2v14a2 2 0 002 2z" /></svg>
           </div>
       </div>

       <!-- Total Credits -->
       <div class="bg-white rounded-3xl p-6 shadow-neu border border-slate-100 flex items-center justify-between">
           <div>
               <p class="text-slate-500 text-sm font-medium mb-1">积分流通量</p>
               <h3 class="text-3xl font-bold text-amber-500">{{ stats.total_credits || 0 }}</h3>
           </div>
           <div class="w-12 h-12 rounded-2xl bg-amber-50 text-amber-500 flex items-center justify-center">
               <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
           </div>
       </div>
    </div>

    <!-- User Table -->
    <div class="bg-white rounded-3xl p-8 shadow-neu border border-slate-100">
      <h2 class="text-xl font-bold text-slate-800 mb-6">用户列表</h2>
      
      <div v-if="loading" class="text-center py-10 text-slate-400">
          加载中...
      </div>

      <div v-else class="overflow-x-auto">
          <table class="w-full text-left border-collapse">
              <thead>
                  <tr>
                      <th class="p-4 border-b border-slate-100 text-slate-400 font-medium text-sm">ID</th>
                      <th class="p-4 border-b border-slate-100 text-slate-400 font-medium text-sm">用户邮箱</th>
                      <th class="p-4 border-b border-slate-100 text-slate-400 font-medium text-sm">角色</th>
                      <th class="p-4 border-b border-slate-100 text-slate-400 font-medium text-sm">当前积分</th>
                      <th class="p-4 border-b border-slate-100 text-slate-400 font-medium text-sm">注册时间</th>
                  </tr>
              </thead>
              <tbody>
                  <tr v-for="user in users" :key="user.id" class="group hover:bg-slate-50 transition-colors">
                      <td class="p-4 border-b border-slate-100 text-slate-500 font-mono text-xs">#{{ user.id }}</td>
                      <td class="p-4 border-b border-slate-100 font-medium text-slate-700">{{ user.email }}</td>
                      <td class="p-4 border-b border-slate-100">
                          <span :class="user.role === 'admin' ? 'bg-purple-100 text-purple-700' : 'bg-slate-100 text-slate-600'" class="px-2 py-1 rounded-md text-xs font-bold uppercase">{{ user.role }}</span>
                      </td>
                      <td class="p-4 border-b border-slate-100 font-mono font-bold text-amber-600">{{ user.credits }}</td>
                      <td class="p-4 border-b border-slate-100 text-slate-400 text-sm">{{ formatDate(user.created_at) }}</td>
                  </tr>
              </tbody>
          </table>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'
import { useRouter } from 'vue-router'

const stats = ref({})
const users = ref([])
const loading = ref(true)
const router = useRouter()

const fetchData = async () => {
    loading.value = true
    try {
        const [statsRes, usersRes] = await Promise.all([
            axios.get('/api/admin/stats'),
            axios.get('/api/admin/users')
        ])
        stats.value = statsRes.data
        users.value = usersRes.data.list
    } catch (err) {
        if (err.response && err.response.status === 403) {
            alert("需要管理员权限")
            router.push('/dashboard')
        }
    } finally {
        loading.value = false
    }
}

const formatDate = (dateStr) => {
    return new Date(dateStr).toLocaleDateString()
}

onMounted(() => {
    fetchData()
})
</script>
