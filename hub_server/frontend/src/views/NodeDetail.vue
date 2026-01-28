<template>
  <div class="view-container">
    <div v-if="loading && !node" class="loading-state">
      <div class="spinner"></div>
    </div>

    <div v-else-if="!node" class="empty-state">
      <p>节点不存在或已删除</p>
      <button class="btn btn-outline" @click="router.push('/dashboard')">返回概览</button>
    </div>

    <div v-else>
      <header class="header">
        <div class="flex items-center gap-4">
            <button class="btn-icon btn-outline" @click="router.back()">
                <ArrowLeft class="icon" />
            </button>
            <div>
                <h1>{{ node.hostname }}</h1>
                <p class="text-mono text-muted">{{ node.id }}</p>
            </div>
        </div>
        <div class="flex gap-2">
            <span class="status-tag" :class="node.status">{{ node.status }}</span>
        </div>
      </header>

      <div class="grid-stats">
         <div class="stat-card">
            <div class="label">IP 地址</div>
            <div class="value">{{ node.ip || 'Unknown' }}</div>
         </div>
         <div class="stat-card">
            <div class="label">客户端版本</div>
            <div class="value">{{ node.version }}</div>
         </div>
         <div class="stat-card">
            <div class="label">首次发现</div>
            <div class="value text-sm">{{ formatTime(node.created_at) }}</div>
         </div>
         <div class="stat-card">
            <div class="label">最近心跳</div>
            <div class="value text-sm">{{ formatTime(node.last_seen) }}</div>
         </div>
      </div>

      <div class="section-title">
        <h3>执行历史</h3>
      </div>

      <div class="table-container">
        <table class="data-table">
          <thead>
            <tr>
              <th>ID</th>
              <th>类型</th>
              <th>状态</th>
              <th>时间</th>
              <th>详情</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="task in tasks" :key="task.id">
              <td>#{{ task.id }}</td>
              <td><span class="type-badge">{{ task.type }}</span></td>
              <td><span class="status-badge" :class="task.status">{{ task.status }}</span></td>
              <td>{{ formatTime(task.created_at) }}</td>
              <td>
                <button class="btn-xs btn-outline" @click="showTaskDetail(task)">查看</button>
              </td>
            </tr>
            <tr v-if="tasks.length === 0">
                <td colspan="5" class="text-center text-muted p-4">暂无历史记录</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useClientStore } from '../store/client'
import { ArrowLeft } from 'lucide-vue-next'
import { formatTime } from '../utils/format'
import axios from 'axios'

const route = useRoute()
const router = useRouter()
const clientStore = useClientStore()

const node = ref(null)
const tasks = ref([])
const loading = ref(true)

onMounted(async () => {
    const id = route.params.id
    // Try to find in store first
    node.value = clientStore.getClientById(id)
    
    // If not found (offline node?), we might need a separate API to get single node detail 
    // For now we assume fetchClients has populated the list including offline nodes (if backend supports it)
    if (!node.value) {
        await clientStore.fetchClients()
        node.value = clientStore.getClientById(id)
    }

    if (node.value) {
        loadTasks(node.value.id)
    } else {
        loading.value = false
    }
})

const loadTasks = async (nodeId) => {
    try {
        const res = await axios.get(`/api/tasks?node_id=${nodeId}&limit=50`)
        tasks.value = res.data.list || []
    } catch (e) {
        console.error(e)
    } finally {
        loading.value = false
    }
}

const showTaskDetail = (task) => {
    alert(JSON.stringify(task, null, 2))
}
</script>

<style scoped>
.view-container { padding: 2rem 3rem; }
.header { margin-bottom: 2rem; display: flex; justify-content: space-between; align-items: center; }
.header h1 { font-family: 'Outfit'; font-size: 1.8rem; font-weight: 700; }
.text-mono { font-family: monospace; }
.text-center { text-align: center; }

.status-tag {
    padding: 4px 12px; border-radius: 6px; font-weight: 600; text-transform: uppercase; font-size: 0.8rem;
}
.status-tag.online { background: rgba(35, 165, 89, 0.15); color: var(--success); }
.status-tag.offline { background: rgba(255, 255, 255, 0.1); color: var(--text-muted); }

.grid-stats {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    gap: 1rem;
    margin-bottom: 2rem;
}

.stat-card {
    background: var(--bg-card);
    border: 1px solid var(--border);
    padding: 1.25rem;
    border-radius: 12px;
}
.stat-card .label { color: var(--text-muted); font-size: 0.8rem; text-transform: uppercase; margin-bottom: 0.5rem; }
.stat-card .value { font-size: 1.1rem; font-weight: 600; color: var(--text-main); }

.section-title { margin-bottom: 1rem; border-bottom: 1px solid var(--border); padding-bottom: 0.5rem; }

.table-container {
    background: var(--bg-card);
    border: 1px solid var(--border);
    border-radius: var(--radius-main);
    overflow: hidden;
}
.data-table { width: 100%; border-collapse: collapse; text-align: left; }
.data-table th { padding: 1rem; background: rgba(255,255,255,0.05); color: var(--text-muted); font-size: 0.85rem; }
.data-table td { padding: 1rem; border-top: 1px solid var(--border); color: var(--text-dim); }

.type-badge { background: rgba(88, 101, 242, 0.15); color: var(--primary); padding: 2px 8px; border-radius: 4px; font-size: 0.8rem; }
.status-badge { font-weight: 600; font-size: 0.8rem; text-transform: capitalize; }
.status-badge.success { color: var(--success); }
.status-badge.failed { color: var(--danger); }
</style>
