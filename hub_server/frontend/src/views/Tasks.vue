<template>
  <div class="view-container">
    <header class="header">
      <div>
        <h1>任务审计日志</h1>
        <p>所有远程操作的历史记录与执行结果</p>
      </div>
      <button class="btn btn-outline btn-sm" @click="taskStore.fetchTasks(taskStore.page)">
        <RefreshCw class="icon-sm" :class="{ 'spin': taskStore.loading }" />
        刷新
      </button>
    </header>

    <div class="content">
      <div v-if="taskStore.loading && !taskStore.tasks.length" class="loading-state">
        <div class="spinner"></div>
        <p>加载中...</p>
      </div>

      <div v-else class="table-container">
        <table class="data-table">
          <thead>
            <tr>
              <th>ID</th>
              <th>类型</th>
              <th>执行节点 (Client ID)</th>
              <th>状态</th>
              <th>创建时间</th>
              <th>详情</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="task in taskStore.tasks" :key="task.id">
              <td>#{{ task.id }}</td>
              <td>
                <span class="type-badge">{{ task.type }}</span>
              </td>
              <td class="font-mono text-sm">{{ task.node_id }}</td>
              <td>
                <span class="status-badge" :class="task.status">
                  {{ task.status }}
                </span>
              </td>
              <td>{{ formatTime(task.created_at) }}</td>
              <td>
                <button class="btn-xs btn-outline" @click="showDetail(task)">查看</button>
              </td>
            </tr>
          </tbody>
        </table>

        <div class="pagination">
          <button 
            :disabled="taskStore.page <= 1"
            class="btn-outline btn-sm"
            @click="taskStore.fetchTasks(taskStore.page - 1)"
          >
            上一页
          </button>
          <span class="page-info">第 {{ taskStore.page }} 页 / 共 {{ Math.ceil(taskStore.total / taskStore.pageSize) || 1 }} 页</span>
          <button 
            :disabled="taskStore.page * taskStore.pageSize >= taskStore.total"
            class="btn-outline btn-sm"
            @click="taskStore.fetchTasks(taskStore.page + 1)"
          >
            下一页
          </button>
        </div>
      </div>
    </div>

    <!-- Task Detail Modal -->
    <div v-if="selectedTask" class="modal-overlay" @click="selectedTask = null">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>任务详情 #{{ selectedTask.id }}</h3>
          <button class="close-btn" @click="selectedTask = null">×</button>
        </div>
        
        <div class="detail-row">
            <label>Payload:</label>
            <pre class="code-block">{{ formatJson(selectedTask.payload) }}</pre>
        </div>
        
        <div class="detail-row">
            <label>Result:</label>
            <pre class="code-block">{{ formatJson(selectedTask.result) }}</pre>
        </div>

        <div v-if="selectedTask.error" class="detail-row">
            <label>Error:</label>
            <div class="error-msg">{{ selectedTask.error }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { useTaskStore } from '../store/task'
import { RefreshCw } from 'lucide-vue-next'
import { formatTime } from '../utils/format'

const taskStore = useTaskStore()
const selectedTask = ref(null)

onMounted(() => {
  taskStore.fetchTasks()
})

const showDetail = (task) => {
    selectedTask.value = task
}

const formatJson = (str) => {
    if (!str) return '-'
    try {
        return JSON.stringify(JSON.parse(str), null, 2)
    } catch (e) {
        return str
    }
}
</script>

<style scoped>
.view-container { padding: 2rem 3rem; }
.header { margin-bottom: 2rem; display: flex; justify-content: space-between; align-items: flex-start; }
.header h1 { font-family: 'Outfit'; font-size: 1.8rem; font-weight: 700; margin-bottom: 0.5rem; }
.header p { color: var(--text-dim); font-size: 0.9rem; }

.table-container {
    background: var(--bg-card);
    border: 1px solid var(--border);
    border-radius: var(--radius-main);
    overflow: hidden;
}

.data-table {
    width: 100%;
    border-collapse: collapse;
    text-align: left;
}
.data-table th {
    padding: 1rem;
    background: rgba(255,255,255,0.05);
    color: var(--text-muted);
    font-weight: 600;
    font-size: 0.85rem;
}
.data-table td {
    padding: 1rem;
    border-top: 1px solid var(--border);
    color: var(--text-dim);
}

.type-badge {
    background: rgba(88, 101, 242, 0.15);
    color: var(--primary);
    padding: 2px 8px;
    border-radius: 4px;
    font-size: 0.8rem;
}

.status-badge {
    text-transform: capitalize;
    font-weight: 600;
    font-size: 0.8rem;
}
.status-badge.pending { color: var(--warning); }
.status-badge.success { color: var(--success); }
.status-badge.failed { color: var(--danger); }
.status-badge.timeout { color: var(--text-muted); }

.pagination {
    padding: 1rem;
    border-top: 1px solid var(--border);
    display: flex;
    justify-content: flex-end;
    gap: 1rem;
    align-items: center;
}
.page-info { font-size: 0.9rem; color: var(--text-muted); }

.modal-overlay {
    position: fixed; top: 0; left: 0; width: 100%; height: 100%;
    background: rgba(0,0,0,0.8); z-index: 100;
    display: flex; justify-content: center; align-items: center;
}
.modal-content {
    width: 600px;
    max-width: 90%;
    background: var(--bg-side);
    border: 1px solid var(--border);
    border-radius: 12px;
    padding: 1.5rem;
    max-height: 80vh;
    overflow-y: auto;
}
.modal-header { display: flex; justify-content: space-between; margin-bottom: 1.5rem; }
.close-btn { background: none; border: none; font-size: 1.5rem; color: var(--text-muted); cursor: pointer; }

.detail-row { margin-bottom: 1rem; }
.detail-row label { display: block; color: var(--text-muted); font-size: 0.8rem; margin-bottom: 0.5rem; }
.code-block {
    background: #000;
    padding: 0.8rem;
    border-radius: 6px;
    font-family: monospace;
    font-size: 0.85rem;
    color: #ccc;
    overflow-x: auto;
    white-space: pre-wrap;
    word-break: break-all;
}
.error-msg { color: var(--danger); }

.spin { animation: spin 1s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }
</style>
