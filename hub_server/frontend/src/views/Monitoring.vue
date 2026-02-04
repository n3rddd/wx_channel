<template>
  <div class="w-full space-y-8 p-8">
    <!-- Header -->
    <div class="flex flex-col md:flex-row justify-between items-start md:items-center gap-4">
      <div class="flex items-center gap-4">
        <button 
          @click="refreshData" 
          :disabled="loading"
          class="px-6 py-3 bg-primary text-white rounded-xl shadow-neu-btn hover:bg-primary-dark transition-all disabled:opacity-50 flex items-center gap-2"
        >
          <component :is="RefreshCw" class="w-5 h-5" :class="{ 'animate-spin': loading }" />
          {{ loading ? '刷新中...' : '刷新数据' }}
        </button>
        <select 
          v-model="timeRange" 
          @change="refreshData" 
          class="px-4 py-3 bg-white border border-slate-200 rounded-xl shadow-sm focus:outline-none focus:ring-2 focus:ring-primary text-slate-700"
        >
          <option value="5m">最近 5 分钟</option>
          <option value="15m">最近 15 分钟</option>
          <option value="1h">最近 1 小时</option>
          <option value="6h">最近 6 小时</option>
          <option value="24h">最近 24 小时</option>
        </select>
      </div>
    </div>

    <!-- 关键指标卡片 -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <!-- WebSocket 连接 -->
      <div class="bg-white rounded-3xl p-6 shadow-card border border-slate-100 flex items-center justify-between hover:shadow-lg transition-all">
        <div>
          <p class="text-slate-500 text-sm font-medium mb-1">WebSocket 连接</p>
          <h3 class="text-3xl font-bold text-slate-800">{{ metrics.connections }}</h3>
          <p class="text-sm font-medium mt-1" :class="getTrendClass(metrics.connectionsTrend)">
            {{ formatTrend(metrics.connectionsTrend) }}
          </p>
        </div>
        <div class="w-12 h-12 rounded-2xl bg-blue-50 text-blue-500 flex items-center justify-center">
          <component :is="Plug" class="w-6 h-6" />
        </div>
      </div>

      <!-- API 调用总数 -->
      <div class="bg-white rounded-3xl p-6 shadow-card border border-slate-100 flex items-center justify-between hover:shadow-lg transition-all">
        <div>
          <p class="text-slate-500 text-sm font-medium mb-1">API 调用总数</p>
          <h3 class="text-3xl font-bold text-slate-800">{{ formatNumber(metrics.apiCalls) }}</h3>
          <p class="text-sm font-medium mt-1" :class="getTrendClass(metrics.apiCallsTrend)">
            {{ formatTrend(metrics.apiCallsTrend) }}
          </p>
        </div>
        <div class="w-12 h-12 rounded-2xl bg-purple-50 text-purple-500 flex items-center justify-center">
          <component :is="Radio" class="w-6 h-6" />
        </div>
      </div>

      <!-- API 成功率 -->
      <div class="bg-white rounded-3xl p-6 shadow-card border border-slate-100 flex items-center justify-between hover:shadow-lg transition-all">
        <div>
          <p class="text-slate-500 text-sm font-medium mb-1">API 成功率</p>
          <h3 class="text-3xl font-bold text-slate-800">{{ metrics.successRate }}%</h3>
          <p class="text-sm font-medium mt-1" :class="getStatusClass(metrics.successRate)">
            {{ getStatusText(metrics.successRate) }}
          </p>
        </div>
        <div class="w-12 h-12 rounded-2xl bg-green-50 text-green-500 flex items-center justify-center">
          <component :is="CheckCircle" class="w-6 h-6" />
        </div>
      </div>

      <!-- 平均响应时间 -->
      <div class="bg-white rounded-3xl p-6 shadow-card border border-slate-100 flex items-center justify-between hover:shadow-lg transition-all">
        <div>
          <p class="text-slate-500 text-sm font-medium mb-1">平均响应时间</p>
          <h3 class="text-3xl font-bold text-slate-800">{{ metrics.avgResponseTime }}ms</h3>
          <p class="text-sm font-medium mt-1" :class="getTrendClass(-metrics.responseTimeTrend)">
            {{ formatTrend(metrics.responseTimeTrend) }}
          </p>
        </div>
        <div class="w-12 h-12 rounded-2xl bg-amber-50 text-amber-500 flex items-center justify-center">
          <component :is="Zap" class="w-6 h-6" />
        </div>
      </div>

      <!-- 心跳状态 -->
      <div class="bg-white rounded-3xl p-6 shadow-card border border-slate-100 flex items-center justify-between hover:shadow-lg transition-all">
        <div>
          <p class="text-slate-500 text-sm font-medium mb-1">心跳状态</p>
          <h3 class="text-3xl font-bold text-slate-800">{{ metrics.heartbeatsSent }}</h3>
          <p class="text-sm font-medium mt-1 text-green-600">
            失败: {{ metrics.heartbeatsFailed }}
          </p>
        </div>
        <div class="w-12 h-12 rounded-2xl bg-red-50 text-red-500 flex items-center justify-center">
          <component :is="Heart" class="w-6 h-6" />
        </div>
      </div>

      <!-- 压缩率 -->
      <div class="bg-white rounded-3xl p-6 shadow-card border border-slate-100 flex items-center justify-between hover:shadow-lg transition-all">
        <div>
          <p class="text-slate-500 text-sm font-medium mb-1">压缩率</p>
          <h3 class="text-3xl font-bold text-slate-800">{{ metrics.compressionRate.toFixed(2) }}%</h3>
          <p class="text-sm font-medium mt-1 text-green-600">
            节省 {{ formatBytes(metrics.bytesSaved) }}
          </p>
        </div>
        <div class="w-12 h-12 rounded-2xl bg-indigo-50 text-indigo-500 flex items-center justify-center">
          <component :is="Package" class="w-6 h-6" />
        </div>
      </div>
    </div>

    <!-- 图表区域 -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- 连接数趋势 -->
      <div class="bg-white rounded-3xl p-6 shadow-card border border-slate-100">
        <h3 class="text-lg font-bold text-slate-800 mb-4 font-serif">WebSocket 连接数趋势</h3>
        <div class="h-[300px] relative">
          <canvas ref="connectionsChart"></canvas>
        </div>
      </div>

      <!-- API 调用趋势 -->
      <div class="bg-white rounded-3xl p-6 shadow-card border border-slate-100">
        <h3 class="text-lg font-bold text-slate-800 mb-4 font-serif">API 调用趋势</h3>
        <div class="h-[300px] relative">
          <canvas ref="apiCallsChart"></canvas>
        </div>
      </div>

      <!-- 响应时间分布 -->
      <div class="bg-white rounded-3xl p-6 shadow-card border border-slate-100">
        <h3 class="text-lg font-bold text-slate-800 mb-4 font-serif">API 响应时间</h3>
        <div class="h-[300px] relative">
          <canvas ref="responseTimeChart"></canvas>
        </div>
      </div>

      <!-- 负载均衡分布 -->
      <div class="bg-white rounded-3xl p-6 shadow-card border border-slate-100">
        <h3 class="text-lg font-bold text-slate-800 mb-4 font-serif">负载均衡分布</h3>
        <div class="h-[300px] relative">
          <canvas ref="loadBalancerChart"></canvas>
        </div>
      </div>
    </div>

    <!-- 详细指标表格 -->
    <div class="bg-white rounded-3xl p-8 shadow-card border border-slate-100">
      <h3 class="text-xl font-bold text-slate-800 mb-6 font-serif">详细指标</h3>
      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse">
          <thead>
            <tr>
              <th class="p-4 border-b border-slate-100 text-slate-400 font-medium text-sm">指标名称</th>
              <th class="p-4 border-b border-slate-100 text-slate-400 font-medium text-sm">当前值</th>
              <th class="p-4 border-b border-slate-100 text-slate-400 font-medium text-sm">说明</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="metric in detailedMetrics" :key="metric.name" class="group hover:bg-slate-50 transition-colors">
              <td class="p-4 border-b border-slate-100 font-medium text-slate-700">{{ metric.name }}</td>
              <td class="p-4 border-b border-slate-100 font-mono font-bold text-slate-800">{{ metric.value }}</td>
              <td class="p-4 border-b border-slate-100 text-slate-400 text-sm">{{ metric.description }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { RefreshCw, Plug, Radio, CheckCircle, Zap, Heart, Package } from 'lucide-vue-next'
import Chart from 'chart.js/auto'

const loading = ref(false)
const timeRange = ref('15m')
const metrics = ref({
  connections: 0,
  connectionsTrend: 0,
  apiCalls: 0,
  apiCallsTrend: 0,
  successRate: 0,
  avgResponseTime: 0,
  responseTimeTrend: 0,
  heartbeatsSent: 0,
  heartbeatsFailed: 0,
  compressionRate: 0,
  bytesSaved: 0
})

const detailedMetrics = ref([])

// Chart 实例
const connectionsChart = ref(null)
const apiCallsChart = ref(null)
const responseTimeChart = ref(null)
const loadBalancerChart = ref(null)

let charts = {}
let refreshInterval = null

// 获取监控数据
async function fetchMetrics() {
  try {
    const token = localStorage.getItem('token')
    const response = await fetch('/api/metrics/summary', {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })
    const data = await response.json()
    
    metrics.value = {
      connections: data.connections || 0,
      connectionsTrend: data.connectionsTrend || 0,
      apiCalls: data.apiCalls || 0,
      apiCallsTrend: data.apiCallsTrend || 0,
      successRate: data.successRate || 0,
      avgResponseTime: data.avgResponseTime || 0,
      responseTimeTrend: data.responseTimeTrend || 0,
      heartbeatsSent: data.heartbeatsSent || 0,
      heartbeatsFailed: data.heartbeatsFailed || 0,
      compressionRate: data.compressionRate || 0,
      bytesSaved: data.bytesSaved || 0
    }

    detailedMetrics.value = data.detailedMetrics || []
    
    return data
  } catch (error) {
    console.error('获取监控数据失败:', error)
    return null
  }
}

// 获取时序数据
async function fetchTimeSeriesData() {
  try {
    const token = localStorage.getItem('token')
    const response = await fetch(`/api/metrics/timeseries?range=${timeRange.value}`, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })
    return await response.json()
  } catch (error) {
    console.error('获取时序数据失败:', error)
    return null
  }
}

// 刷新数据
async function refreshData() {
  loading.value = true
  try {
    await fetchMetrics()
    const timeSeriesData = await fetchTimeSeriesData()
    if (timeSeriesData) {
      updateCharts(timeSeriesData)
    }
  } finally {
    loading.value = false
  }
}

// 初始化图表
function initCharts() {
  // 连接数趋势图
  if (connectionsChart.value) {
    charts.connections = new Chart(connectionsChart.value, {
      type: 'line',
      data: {
        labels: [],
        datasets: [{
          label: '连接数',
          data: [],
          borderColor: '#3b82f6',
          backgroundColor: 'rgba(59, 130, 246, 0.1)',
          tension: 0.4,
          fill: true
        }]
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
          legend: { display: false }
        },
        scales: {
          y: { beginAtZero: true }
        }
      }
    })
  }

  // API 调用趋势图
  if (apiCallsChart.value) {
    charts.apiCalls = new Chart(apiCallsChart.value, {
      type: 'line',
      data: {
        labels: [],
        datasets: [
          {
            label: '成功',
            data: [],
            borderColor: '#10b981',
            backgroundColor: 'rgba(16, 185, 129, 0.1)',
            tension: 0.4,
            fill: true
          },
          {
            label: '失败',
            data: [],
            borderColor: '#ef4444',
            backgroundColor: 'rgba(239, 68, 68, 0.1)',
            tension: 0.4,
            fill: true
          }
        ]
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        scales: {
          y: { beginAtZero: true }
        }
      }
    })
  }

  // 响应时间图
  if (responseTimeChart.value) {
    charts.responseTime = new Chart(responseTimeChart.value, {
      type: 'line',
      data: {
        labels: [],
        datasets: [
          {
            label: 'P50',
            data: [],
            borderColor: '#3b82f6',
            tension: 0.4
          },
          {
            label: 'P95',
            data: [],
            borderColor: '#f59e0b',
            tension: 0.4
          },
          {
            label: 'P99',
            data: [],
            borderColor: '#ef4444',
            tension: 0.4
          }
        ]
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        scales: {
          y: { beginAtZero: true }
        }
      }
    })
  }

  // 负载均衡分布图
  if (loadBalancerChart.value) {
    charts.loadBalancer = new Chart(loadBalancerChart.value, {
      type: 'bar',
      data: {
        labels: [],
        datasets: [{
          label: '请求数',
          data: [],
          backgroundColor: [
            'rgba(59, 130, 246, 0.8)',
            'rgba(16, 185, 129, 0.8)',
            'rgba(245, 158, 11, 0.8)',
            'rgba(139, 92, 246, 0.8)'
          ]
        }]
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        scales: {
          y: { beginAtZero: true }
        }
      }
    })
  }
}

// 更新图表
function updateCharts(data) {
  if (charts.connections && data.connections) {
    charts.connections.data.labels = data.connections.labels
    charts.connections.data.datasets[0].data = data.connections.values
    charts.connections.update()
  }

  if (charts.apiCalls && data.apiCalls) {
    charts.apiCalls.data.labels = data.apiCalls.labels
    charts.apiCalls.data.datasets[0].data = data.apiCalls.success
    charts.apiCalls.data.datasets[1].data = data.apiCalls.failed
    charts.apiCalls.update()
  }

  if (charts.responseTime && data.responseTime) {
    charts.responseTime.data.labels = data.responseTime.labels
    charts.responseTime.data.datasets[0].data = data.responseTime.p50
    charts.responseTime.data.datasets[1].data = data.responseTime.p95
    charts.responseTime.data.datasets[2].data = data.responseTime.p99
    charts.responseTime.update()
  }

  if (charts.loadBalancer && data.loadBalancer) {
    charts.loadBalancer.data.labels = data.loadBalancer.labels
    charts.loadBalancer.data.datasets[0].data = data.loadBalancer.values
    charts.loadBalancer.update()
  }
}

// 格式化函数
function formatNumber(num) {
  if (num >= 1000000) return (num / 1000000).toFixed(1) + 'M'
  if (num >= 1000) return (num / 1000).toFixed(1) + 'K'
  return num.toString()
}

function formatBytes(bytes) {
  if (bytes >= 1073741824) return (bytes / 1073741824).toFixed(2) + ' GB'
  if (bytes >= 1048576) return (bytes / 1048576).toFixed(2) + ' MB'
  if (bytes >= 1024) return (bytes / 1024).toFixed(2) + ' KB'
  return bytes + ' B'
}

function formatTrend(trend) {
  if (trend > 0) return `↑ ${trend.toFixed(1)}%`
  if (trend < 0) return `↓ ${Math.abs(trend).toFixed(1)}%`
  return '→ 0%'
}

function getTrendClass(trend) {
  if (trend > 0) return 'text-green-600'
  if (trend < 0) return 'text-red-600'
  return 'text-slate-500'
}

function getStatusClass(rate) {
  if (rate >= 95) return 'text-green-600'
  if (rate >= 90) return 'text-amber-600'
  return 'text-red-600'
}

function getStatusText(rate) {
  if (rate >= 95) return '优秀'
  if (rate >= 90) return '良好'
  return '需关注'
}

onMounted(async () => {
  await refreshData()
  initCharts()
  
  // 每 10 秒自动刷新
  refreshInterval = setInterval(refreshData, 10000)
})

onUnmounted(() => {
  if (refreshInterval) {
    clearInterval(refreshInterval)
  }
  
  // 销毁图表
  Object.values(charts).forEach(chart => {
    if (chart) chart.destroy()
  })
})
</script>
