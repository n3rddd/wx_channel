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

    <!-- WebSocket 连接详情 -->
    <div class="bg-white rounded-3xl p-8 shadow-card border border-slate-100">
      <div class="flex justify-between items-center mb-6">
        <h3 class="text-xl font-bold text-slate-800 font-serif">WebSocket 连接详情</h3>
        <div class="flex items-center gap-4">
          <span class="text-sm text-slate-500">
            总连接: <span class="font-bold text-slate-800">{{ wsStats.total_connections }}</span>
          </span>
          <span class="text-sm text-slate-500">
            Ping 成功率: <span class="font-bold" :class="getPingSuccessRateClass(wsStats.ping_success_rate)">
              {{ wsStats.ping_success_rate }}%
            </span>
          </span>
        </div>
      </div>

      <!-- WebSocket 统计卡片 -->
      <div class="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
        <div class="bg-gradient-to-br from-blue-50 to-blue-100 rounded-2xl p-4 border border-blue-200">
          <p class="text-blue-600 text-xs font-medium mb-1">总 Ping 次数</p>
          <p class="text-2xl font-bold text-blue-900">{{ formatNumber(wsStats.total_pings) }}</p>
        </div>
        <div class="bg-gradient-to-br from-green-50 to-green-100 rounded-2xl p-4 border border-green-200">
          <p class="text-green-600 text-xs font-medium mb-1">总 Pong 次数</p>
          <p class="text-2xl font-bold text-green-900">{{ formatNumber(wsStats.total_pongs) }}</p>
        </div>
        <div class="bg-gradient-to-br from-purple-50 to-purple-100 rounded-2xl p-4 border border-purple-200">
          <p class="text-purple-600 text-xs font-medium mb-1">总消息数</p>
          <p class="text-2xl font-bold text-purple-900">{{ formatNumber(wsStats.total_messages) }}</p>
        </div>
        <div class="bg-gradient-to-br from-amber-50 to-amber-100 rounded-2xl p-4 border border-amber-200">
          <p class="text-amber-600 text-xs font-medium mb-1">平均延迟</p>
          <p class="text-2xl font-bold text-amber-900">{{ wsStats.avg_latency }}</p>
        </div>
      </div>

      <!-- 连接列表 -->
      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse">
          <thead>
            <tr>
              <th class="p-4 border-b border-slate-100 text-slate-400 font-medium text-sm">客户端 ID</th>
              <th class="p-4 border-b border-slate-100 text-slate-400 font-medium text-sm">主机名</th>
              <th class="p-4 border-b border-slate-100 text-slate-400 font-medium text-sm">IP 地址</th>
              <th class="p-4 border-b border-slate-100 text-slate-400 font-medium text-sm">运行时长</th>
              <th class="p-4 border-b border-slate-100 text-slate-400 font-medium text-sm">Ping/Pong</th>
              <th class="p-4 border-b border-slate-100 text-slate-400 font-medium text-sm">平均延迟</th>
              <th class="p-4 border-b border-slate-100 text-slate-400 font-medium text-sm">消息数</th>
              <th class="p-4 border-b border-slate-100 text-slate-400 font-medium text-sm">状态</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="wsStats.clients && wsStats.clients.length === 0">
              <td colspan="8" class="p-8 text-center text-slate-400">
                暂无连接
              </td>
            </tr>
            <tr v-for="client in wsStats.clients" :key="client.id" class="group hover:bg-slate-50 transition-colors">
              <td class="p-4 border-b border-slate-100 font-mono text-sm text-slate-700">
                {{ client.id.substring(0, 12) }}...
              </td>
              <td class="p-4 border-b border-slate-100 font-medium text-slate-700">
                {{ client.hostname || '-' }}
              </td>
              <td class="p-4 border-b border-slate-100 font-mono text-sm text-slate-600">
                {{ client.ip }}
              </td>
              <td class="p-4 border-b border-slate-100 text-slate-700">
                {{ client.uptime }}
              </td>
              <td class="p-4 border-b border-slate-100 font-mono text-sm">
                <span class="text-green-600">{{ client.ping_count }}</span> / 
                <span class="text-blue-600">{{ client.pong_count }}</span>
              </td>
              <td class="p-4 border-b border-slate-100 font-mono text-sm" :class="getLatencyClass(client.avg_latency)">
                {{ client.avg_latency }}
              </td>
              <td class="p-4 border-b border-slate-100 font-mono text-sm text-slate-700">
                <span class="text-purple-600">↑{{ client.messages_sent }}</span> / 
                <span class="text-blue-600">↓{{ client.messages_recv }}</span>
              </td>
              <td class="p-4 border-b border-slate-100">
                <span v-if="client.failure_count === 0" class="px-3 py-1 bg-green-100 text-green-700 rounded-full text-xs font-medium">
                  正常
                </span>
                <span v-else class="px-3 py-1 bg-red-100 text-red-700 rounded-full text-xs font-medium">
                  失败 {{ client.failure_count }}
                </span>
              </td>
            </tr>
          </tbody>
        </table>
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
const wsStats = ref({
  total_connections: 0,
  total_pings: 0,
  total_pongs: 0,
  total_messages: 0,
  ping_success_rate: 0,
  avg_latency: '-',
  clients: []
})

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

// 获取 WebSocket 统计
async function fetchWSStats() {
  try {
    const token = localStorage.getItem('token')
    const response = await fetch('/api/ws/stats', {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })
    const result = await response.json()
    
    if (result.code === 0 && result.data) {
      const data = result.data
      
      // 计算 Ping 成功率
      const pingSuccessRate = data.total_pings > 0 
        ? ((data.total_pongs / data.total_pings) * 100).toFixed(2)
        : 0
      
      // 计算平均延迟
      let avgLatency = '-'
      if (data.clients && data.clients.length > 0) {
        const latencies = data.clients
          .map(c => c.avg_latency)
          .filter(l => l && l !== '-')
        
        if (latencies.length > 0) {
          avgLatency = latencies[0] // 使用第一个客户端的延迟作为示例
        }
      }
      
      wsStats.value = {
        total_connections: data.total_connections || 0,
        total_pings: data.total_pings || 0,
        total_pongs: data.total_pongs || 0,
        total_messages: data.total_messages || 0,
        ping_success_rate: pingSuccessRate,
        avg_latency: avgLatency,
        clients: data.clients || []
      }
    }
  } catch (error) {
    console.error('获取 WebSocket 统计失败:', error)
  }
}

// 刷新数据
async function refreshData() {
  loading.value = true
  try {
    await fetchMetrics()
    await fetchWSStats()
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

function getPingSuccessRateClass(rate) {
  const numRate = parseFloat(rate)
  if (numRate >= 95) return 'text-green-600'
  if (numRate >= 90) return 'text-amber-600'
  return 'text-red-600'
}

function getLatencyClass(latency) {
  if (!latency || latency === '-') return 'text-slate-500'
  const ms = parseInt(latency)
  if (ms < 100) return 'text-green-600'
  if (ms < 500) return 'text-amber-600'
  return 'text-red-600'
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
