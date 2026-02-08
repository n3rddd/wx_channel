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
    <div class="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
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

      <div class="bg-white shadow-card rounded-2xl p-6 border border-slate-100">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm text-slate-500 mb-1">锁定设备</p>
            <p class="text-3xl font-bold text-amber-600">{{ lockedCount }}</p>
          </div>
          <div class="w-12 h-12 bg-amber-100 rounded-xl flex items-center justify-center">
            <component :is="Lock" class="w-6 h-6 text-amber-600" />
          </div>
        </div>
      </div>
    </div>

    <!-- 筛选和搜索 -->
    <div class="bg-white shadow-card rounded-2xl p-6 mb-6 border border-slate-100">
      <div class="flex flex-col md:flex-row gap-4">
        <div class="flex-1">
          <input
            v-model="searchQuery"
            type="text"
            placeholder="搜索设备 ID、主机名或显示名称..."
            class="w-full px-4 py-2 bg-bg shadow-neu-pressed rounded-xl text-slate-700 placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-primary/20"
          />
        </div>
        <select
          v-model="filterStatus"
          class="px-4 py-2 bg-bg shadow-neu-pressed rounded-xl text-slate-700 focus:outline-none focus:ring-2 focus:ring-primary/20"
        >
          <option value="">全部状态</option>
          <option value="online">在线</option>
          <option value="offline">离线</option>
        </select>
        <select
          v-model="filterGroup"
          class="px-4 py-2 bg-bg shadow-neu-pressed rounded-xl text-slate-700 focus:outline-none focus:ring-2 focus:ring-primary/20"
        >
          <option value="">全部分组</option>
          <option v-for="group in deviceGroups" :key="group" :value="group">{{ group }}</option>
        </select>
      </div>
    </div>

    <!-- 设备列表 -->
    <div class="bg-white shadow-card rounded-2xl p-6 border border-slate-100">
      <div class="flex items-center justify-between mb-6">
        <h2 class="text-xl font-bold text-slate-800">设备列表 ({{ filteredDevices.length }})</h2>
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

      <div v-else-if="filteredDevices.length === 0" class="text-center py-12">
        <component :is="Monitor" class="w-16 h-16 text-slate-300 mx-auto mb-4" />
        <p class="text-slate-500 mb-2">{{ devices.length === 0 ? '暂无设备' : '没有符合条件的设备' }}</p>
        <p class="text-sm text-slate-400">{{ devices.length === 0 ? '请使用上方的绑定码绑定设备' : '尝试调整筛选条件' }}</p>
      </div>

      <div v-else class="space-y-4">
        <div 
          v-for="device in filteredDevices" 
          :key="device.id"
          class="bg-white shadow-card rounded-xl p-6 hover:shadow-lg transition-all border border-slate-100"
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
                <div class="flex items-center gap-3 mb-3">
                  <h3 class="text-lg font-bold text-slate-800">
                    {{ device.display_name || device.hostname || device.id }}
                  </h3>
                  <span 
                    class="px-3 py-1 rounded-full text-xs font-medium"
                    :class="device.status === 'online' 
                      ? 'bg-green-100 text-green-700' 
                      : 'bg-slate-100 text-slate-500'"
                  >
                    {{ device.status === 'online' ? '在线' : '离线' }}
                  </span>
                  <span 
                    v-if="device.is_locked"
                    class="px-3 py-1 rounded-full text-xs font-medium bg-amber-100 text-amber-700 flex items-center gap-1"
                  >
                    <component :is="Lock" class="w-3 h-3" />
                    已锁定
                  </span>
                  <span 
                    v-if="device.device_group"
                    class="px-3 py-1 rounded-full text-xs font-medium bg-blue-100 text-blue-700 flex items-center gap-1"
                  >
                    <component :is="FolderOpen" class="w-3 h-3" />
                    {{ device.device_group }}
                  </span>
                </div>

                <div class="grid grid-cols-2 md:grid-cols-3 gap-4 text-sm mb-3">
                  <div>
                    <span class="text-slate-500">设备 ID：</span>
                    <span class="text-slate-700 font-mono text-xs">{{ device.id }}</span>
                  </div>
                  <div>
                    <span class="text-slate-500">版本：</span>
                    <span class="text-slate-700 font-medium">{{ device.version || 'N/A' }}</span>
                  </div>
                  <div>
                    <span class="text-slate-500">主机名：</span>
                    <span class="text-slate-700 font-medium">{{ device.hostname || 'N/A' }}</span>
                  </div>
                  <div>
                    <span class="text-slate-500">IP 地址：</span>
                    <span class="text-slate-700 font-medium">{{ device.ip || 'N/A' }}</span>
                  </div>
                  <div>
                    <span class="text-slate-500">最后在线：</span>
                    <span class="text-slate-700 font-medium">{{ formatTime(device.last_seen) }}</span>
                  </div>
                  <div>
                    <span class="text-slate-500">首次连接：</span>
                    <span class="text-slate-700 font-medium">{{ formatTime(device.first_seen || device.created_at) }}</span>
                  </div>
                </div>

                <!-- 快捷操作 -->
                <div class="flex flex-wrap gap-2">
                  <button
                    @click="showRenameDialog(device)"
                    class="px-3 py-1.5 bg-bg shadow-neu rounded-lg text-slate-600 hover:shadow-neu-sm transition-all flex items-center gap-1.5 text-sm"
                  >
                    <component :is="Edit3" class="w-3.5 h-3.5" />
                    重命名
                  </button>
                  <button
                    @click="toggleLock(device)"
                    class="px-3 py-1.5 bg-bg shadow-neu rounded-lg hover:shadow-neu-sm transition-all flex items-center gap-1.5 text-sm"
                    :class="device.is_locked ? 'text-amber-600' : 'text-slate-600'"
                  >
                    <component :is="device.is_locked ? Unlock : Lock" class="w-3.5 h-3.5" />
                    {{ device.is_locked ? '解锁' : '锁定' }}
                  </button>
                  <button
                    @click="showGroupDialog(device)"
                    class="px-3 py-1.5 bg-bg shadow-neu rounded-lg text-slate-600 hover:shadow-neu-sm transition-all flex items-center gap-1.5 text-sm"
                  >
                    <component :is="FolderOpen" class="w-3.5 h-3.5" />
                    分组
                  </button>
                  <button
                    @click="showHardwareInfo(device)"
                    class="px-3 py-1.5 bg-bg shadow-neu rounded-lg text-slate-600 hover:shadow-neu-sm transition-all flex items-center gap-1.5 text-sm"
                  >
                    <component :is="Info" class="w-3.5 h-3.5" />
                    详情
                  </button>
                </div>
              </div>
            </div>

            <!-- 主要操作按钮 -->
            <div class="flex gap-2 ml-4">
              <button
                @click="showTransferDialog(device)"
                class="px-4 py-2 bg-bg shadow-neu rounded-xl text-blue-600 hover:shadow-neu-sm transition-all flex items-center gap-2"
                title="转移设备"
                :disabled="device.is_locked"
              >
                <component :is="ArrowRightLeft" class="w-4 h-4" />
                转移
              </button>
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
              <span class="font-medium">设备：</span>{{ selectedDevice?.display_name || selectedDevice?.hostname || selectedDevice?.id }}
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

    <!-- 重命名对话框 -->
    <div 
      v-if="showRename"
      class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
      @click.self="showRename = false"
    >
      <div class="bg-white shadow-card rounded-2xl p-8 max-w-md w-full mx-4 border border-slate-100">
        <div class="flex items-center gap-4 mb-6">
          <div class="w-12 h-12 bg-blue-100 rounded-xl flex items-center justify-center">
            <component :is="Edit3" class="w-6 h-6 text-blue-600" />
          </div>
          <div>
            <h3 class="text-xl font-bold text-slate-800">重命名设备</h3>
            <p class="text-sm text-slate-500">设置一个易于识别的名称</p>
          </div>
        </div>

        <div class="mb-6">
          <label class="block text-sm font-medium text-slate-700 mb-2">显示名称</label>
          <input
            v-model="renameInput"
            type="text"
            placeholder="例如：我的工作电脑"
            class="w-full px-4 py-3 bg-bg shadow-neu-pressed rounded-xl text-slate-700 placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-primary/20"
            @keyup.enter="executeRename"
          />
          <p class="text-xs text-slate-400 mt-2">留空则使用主机名</p>
        </div>

        <div class="flex gap-3">
          <button
            @click="showRename = false"
            class="flex-1 px-4 py-3 bg-bg shadow-neu rounded-xl text-slate-600 hover:shadow-neu-sm transition-all"
          >
            取消
          </button>
          <button
            @click="executeRename"
            :disabled="actionLoading"
            class="flex-1 px-4 py-3 bg-blue-600 hover:bg-blue-700 rounded-xl text-white transition-all flex items-center justify-center gap-2"
          >
            <component v-if="actionLoading" :is="Loader2" class="w-4 h-4 animate-spin" />
            <span>{{ actionLoading ? '保存中...' : '保存' }}</span>
          </button>
        </div>
      </div>
    </div>

    <!-- 分组对话框 -->
    <div 
      v-if="showGroup"
      class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
      @click.self="showGroup = false"
    >
      <div class="bg-white shadow-card rounded-2xl p-8 max-w-md w-full mx-4 border border-slate-100">
        <div class="flex items-center gap-4 mb-6">
          <div class="w-12 h-12 bg-purple-100 rounded-xl flex items-center justify-center">
            <component :is="FolderOpen" class="w-6 h-6 text-purple-600" />
          </div>
          <div>
            <h3 class="text-xl font-bold text-slate-800">设置分组</h3>
            <p class="text-sm text-slate-500">将设备分配到分组</p>
          </div>
        </div>

        <div class="mb-6">
          <label class="block text-sm font-medium text-slate-700 mb-2">设备分组</label>
          <input
            v-model="groupInput"
            type="text"
            placeholder="例如：办公室、开发环境"
            list="group-suggestions"
            class="w-full px-4 py-3 bg-bg shadow-neu-pressed rounded-xl text-slate-700 placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-primary/20"
            @keyup.enter="executeSetGroup"
          />
          <datalist id="group-suggestions">
            <option v-for="group in deviceGroups" :key="group" :value="group" />
          </datalist>
          <p class="text-xs text-slate-400 mt-2">留空则移除分组</p>
        </div>

        <div class="flex gap-3">
          <button
            @click="showGroup = false"
            class="flex-1 px-4 py-3 bg-bg shadow-neu rounded-xl text-slate-600 hover:shadow-neu-sm transition-all"
          >
            取消
          </button>
          <button
            @click="executeSetGroup"
            :disabled="actionLoading"
            class="flex-1 px-4 py-3 bg-purple-600 hover:bg-purple-700 rounded-xl text-white transition-all flex items-center justify-center gap-2"
          >
            <component v-if="actionLoading" :is="Loader2" class="w-4 h-4 animate-spin" />
            <span>{{ actionLoading ? '保存中...' : '保存' }}</span>
          </button>
        </div>
      </div>
    </div>

    <!-- 转移对话框 -->
    <div 
      v-if="showTransfer"
      class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
      @click.self="showTransfer = false"
    >
      <div class="bg-white shadow-card rounded-2xl p-8 max-w-md w-full mx-4 border border-slate-100">
        <div class="flex items-center gap-4 mb-6">
          <div class="w-12 h-12 bg-indigo-100 rounded-xl flex items-center justify-center">
            <component :is="ArrowRightLeft" class="w-6 h-6 text-indigo-600" />
          </div>
          <div>
            <h3 class="text-xl font-bold text-slate-800">转移设备</h3>
            <p class="text-sm text-slate-500">将设备转移给其他用户</p>
          </div>
        </div>

        <div class="mb-6">
          <label class="block text-sm font-medium text-slate-700 mb-2">目标用户 ID</label>
          <input
            v-model.number="transferUserId"
            type="number"
            placeholder="输入用户 ID"
            class="w-full px-4 py-3 bg-bg shadow-neu-pressed rounded-xl text-slate-700 placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-primary/20"
            @keyup.enter="executeTransfer"
          />
          <div class="mt-4 p-4 bg-amber-50 border border-amber-200 rounded-xl">
            <p class="text-sm text-amber-800">
              <strong>警告：</strong>转移后您将失去对此设备的控制权，且操作不可撤销。
            </p>
          </div>
        </div>

        <div class="flex gap-3">
          <button
            @click="showTransfer = false"
            class="flex-1 px-4 py-3 bg-bg shadow-neu rounded-xl text-slate-600 hover:shadow-neu-sm transition-all"
          >
            取消
          </button>
          <button
            @click="executeTransfer"
            :disabled="actionLoading || !transferUserId"
            class="flex-1 px-4 py-3 bg-indigo-600 hover:bg-indigo-700 rounded-xl text-white transition-all flex items-center justify-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <component v-if="actionLoading" :is="Loader2" class="w-4 h-4 animate-spin" />
            <span>{{ actionLoading ? '转移中...' : '确认转移' }}</span>
          </button>
        </div>
      </div>
    </div>

    <!-- 硬件信息对话框 -->
    <div 
      v-if="showHardware"
      class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
      @click.self="showHardware = false"
    >
      <div class="bg-white shadow-card rounded-2xl p-8 max-w-2xl w-full mx-4 border border-slate-100 max-h-[90vh] overflow-y-auto">
        <div class="flex items-center justify-between mb-6">
          <div class="flex items-center gap-4">
            <div class="w-12 h-12 bg-slate-100 rounded-xl flex items-center justify-center">
              <component :is="Info" class="w-6 h-6 text-slate-600" />
            </div>
            <div>
              <h3 class="text-xl font-bold text-slate-800">设备详情</h3>
              <p class="text-sm text-slate-500">{{ selectedDevice?.display_name || selectedDevice?.hostname || selectedDevice?.id }}</p>
            </div>
          </div>
          <button
            @click="showHardware = false"
            class="p-2 hover:bg-slate-100 rounded-lg transition-all"
          >
            <component :is="X" class="w-5 h-5 text-slate-400" />
          </button>
        </div>

        <div class="space-y-6">
          <!-- 基本信息 -->
          <div>
            <h4 class="text-sm font-bold text-slate-700 mb-3 flex items-center gap-2">
              <component :is="Monitor" class="w-4 h-4" />
              基本信息
            </h4>
            <div class="bg-slate-50 rounded-xl p-4 space-y-2">
              <div class="flex justify-between text-sm">
                <span class="text-slate-500">设备 ID</span>
                <span class="text-slate-700 font-mono">{{ selectedDevice?.id }}</span>
              </div>
              <div class="flex justify-between text-sm">
                <span class="text-slate-500">显示名称</span>
                <span class="text-slate-700">{{ selectedDevice?.display_name || '-' }}</span>
              </div>
              <div class="flex justify-between text-sm">
                <span class="text-slate-500">主机名</span>
                <span class="text-slate-700">{{ selectedDevice?.hostname || '-' }}</span>
              </div>
              <div class="flex justify-between text-sm">
                <span class="text-slate-500">版本</span>
                <span class="text-slate-700">{{ selectedDevice?.version || '-' }}</span>
              </div>
              <div class="flex justify-between text-sm">
                <span class="text-slate-500">IP 地址</span>
                <span class="text-slate-700">{{ selectedDevice?.ip || '-' }}</span>
              </div>
              <div class="flex justify-between text-sm">
                <span class="text-slate-500">状态</span>
                <span 
                  class="px-2 py-0.5 rounded-full text-xs font-medium"
                  :class="selectedDevice?.status === 'online' 
                    ? 'bg-green-100 text-green-700' 
                    : 'bg-slate-100 text-slate-500'"
                >
                  {{ selectedDevice?.status === 'online' ? '在线' : '离线' }}
                </span>
              </div>
              <div class="flex justify-between text-sm">
                <span class="text-slate-500">锁定状态</span>
                <span 
                  class="px-2 py-0.5 rounded-full text-xs font-medium"
                  :class="selectedDevice?.is_locked 
                    ? 'bg-amber-100 text-amber-700' 
                    : 'bg-slate-100 text-slate-500'"
                >
                  {{ selectedDevice?.is_locked ? '已锁定' : '未锁定' }}
                </span>
              </div>
              <div class="flex justify-between text-sm">
                <span class="text-slate-500">设备分组</span>
                <span class="text-slate-700">{{ selectedDevice?.device_group || '-' }}</span>
              </div>
              <div class="flex justify-between text-sm">
                <span class="text-slate-500">首次连接</span>
                <span class="text-slate-700">{{ formatTime(selectedDevice?.first_seen || selectedDevice?.created_at) }}</span>
              </div>
              <div class="flex justify-between text-sm">
                <span class="text-slate-500">最后在线</span>
                <span class="text-slate-700">{{ formatTime(selectedDevice?.last_seen) }}</span>
              </div>
            </div>
          </div>

          <!-- 硬件指纹 -->
          <div v-if="hardwareFingerprint">
            <h4 class="text-sm font-bold text-slate-700 mb-3 flex items-center gap-2">
              <component :is="Cpu" class="w-4 h-4" />
              硬件指纹
            </h4>
            <div class="bg-slate-50 rounded-xl p-4 space-y-2">
              <div v-if="hardwareFingerprint.mac_addresses?.length" class="text-sm">
                <span class="text-slate-500 block mb-1">MAC 地址</span>
                <div class="flex flex-wrap gap-2">
                  <span 
                    v-for="(mac, idx) in hardwareFingerprint.mac_addresses" 
                    :key="idx"
                    class="px-2 py-1 bg-white rounded-lg text-slate-700 font-mono text-xs"
                  >
                    {{ mac }}
                  </span>
                </div>
              </div>
              <div v-if="hardwareFingerprint.cpu_info" class="flex justify-between text-sm">
                <span class="text-slate-500">CPU 信息</span>
                <span class="text-slate-700 font-mono text-xs">{{ hardwareFingerprint.cpu_info }}</span>
              </div>
              <div v-if="hardwareFingerprint.motherboard_id" class="flex justify-between text-sm">
                <span class="text-slate-500">主板序列号</span>
                <span class="text-slate-700 font-mono text-xs">{{ hardwareFingerprint.motherboard_id }}</span>
              </div>
              <div v-if="hardwareFingerprint.disk_serial" class="flex justify-between text-sm">
                <span class="text-slate-500">硬盘序列号</span>
                <span class="text-slate-700 font-mono text-xs">{{ hardwareFingerprint.disk_serial }}</span>
              </div>
              <div v-if="hardwareFingerprint.hostname" class="flex justify-between text-sm">
                <span class="text-slate-500">系统主机名</span>
                <span class="text-slate-700">{{ hardwareFingerprint.hostname }}</span>
              </div>
              <div v-if="hardwareFingerprint.os" class="flex justify-between text-sm">
                <span class="text-slate-500">操作系统</span>
                <span class="text-slate-700">{{ hardwareFingerprint.os }}</span>
              </div>
            </div>
          </div>
          <div v-else class="text-center py-8">
            <p class="text-slate-400 text-sm">暂无硬件指纹信息</p>
            <p class="text-slate-400 text-xs mt-1">设备需要升级到最新版本</p>
          </div>
        </div>

        <div class="mt-6">
          <button
            @click="showHardware = false"
            class="w-full px-4 py-3 bg-bg shadow-neu rounded-xl text-slate-600 hover:shadow-neu-sm transition-all"
          >
            关闭
          </button>
        </div>
      </div>
    </div>

    <!-- Toast 通知 -->
    <div 
      v-if="toast.show"
      class="fixed bottom-8 right-8 z-50 animate-slide-up"
    >
      <div 
        class="bg-white shadow-card rounded-xl p-4 border flex items-center gap-3 min-w-[300px]"
        :class="{
          'border-green-200': toast.type === 'success',
          'border-red-200': toast.type === 'error',
          'border-blue-200': toast.type === 'info'
        }"
      >
        <div 
          class="w-10 h-10 rounded-lg flex items-center justify-center"
          :class="{
            'bg-green-100': toast.type === 'success',
            'bg-red-100': toast.type === 'error',
            'bg-blue-100': toast.type === 'info'
          }"
        >
          <component 
            :is="toast.type === 'success' ? CheckCircle : toast.type === 'error' ? XCircle : Info" 
            class="w-5 h-5"
            :class="{
              'text-green-600': toast.type === 'success',
              'text-red-600': toast.type === 'error',
              'text-blue-600': toast.type === 'info'
            }"
          />
        </div>
        <p class="text-slate-700 flex-1">{{ toast.message }}</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { 
  Monitor, Wifi, WifiOff, RefreshCw, Loader2, Unlink, Trash2, Copy, Lock, Unlock,
  Edit3, FolderOpen, ArrowRightLeft, Info, Cpu, X, CheckCircle, XCircle
} from 'lucide-vue-next'
import axios from 'axios'

const bindToken = ref('')
const devices = ref([])
const loading = ref(false)
const showConfirm = ref(false)
const confirmAction = ref('') // 'unbind' or 'delete'
const selectedDevice = ref(null)
const actionLoading = ref(false)

// 新增状态
const showRename = ref(false)
const showGroup = ref(false)
const showTransfer = ref(false)
const showHardware = ref(false)
const renameInput = ref('')
const groupInput = ref('')
const transferUserId = ref(null)
const hardwareFingerprint = ref(null)

// 筛选状态
const searchQuery = ref('')
const filterStatus = ref('')
const filterGroup = ref('')

// Toast 通知
const toast = ref({
  show: false,
  type: 'success', // 'success', 'error', 'info'
  message: ''
})

const onlineCount = computed(() => 
  devices.value.filter(d => d.status === 'online').length
)

const offlineCount = computed(() => 
  devices.value.filter(d => d.status !== 'online').length
)

const lockedCount = computed(() =>
  devices.value.filter(d => d.is_locked).length
)

const deviceGroups = computed(() => {
  const groups = new Set()
  devices.value.forEach(d => {
    if (d.device_group) groups.add(d.device_group)
  })
  return Array.from(groups).sort()
})

const filteredDevices = computed(() => {
  return devices.value.filter(device => {
    // 搜索过滤
    if (searchQuery.value) {
      const query = searchQuery.value.toLowerCase()
      const matchId = device.id.toLowerCase().includes(query)
      const matchHostname = device.hostname?.toLowerCase().includes(query)
      const matchDisplayName = device.display_name?.toLowerCase().includes(query)
      if (!matchId && !matchHostname && !matchDisplayName) return false
    }

    // 状态过滤
    if (filterStatus.value && device.status !== filterStatus.value) {
      return false
    }

    // 分组过滤
    if (filterGroup.value && device.device_group !== filterGroup.value) {
      return false
    }

    return true
  })
})

const showToast = (message, type = 'success') => {
  toast.value = { show: true, message, type }
  setTimeout(() => {
    toast.value.show = false
  }, 3000)
}

const generateToken = async () => {
  try {
    const res = await axios.post('/api/device/bind_token')
    bindToken.value = res.data.token
    showToast('绑定码生成成功', 'success')
  } catch (err) {
    console.error('Generate token failed', err)
    showToast('生成绑定码失败', 'error')
  }
}

const copyToken = () => {
  navigator.clipboard.writeText(bindToken.value)
  showToast('绑定码已复制到剪贴板', 'success')
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
    // 适配新的响应格式 {code: 0, devices: [...]}
    if (response.data.code === 0) {
      devices.value = response.data.devices || []
    } else {
      devices.value = response.data || []
    }
  } catch (error) {
    console.error('Failed to load devices:', error)
    showToast('加载设备列表失败', 'error')
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

    showToast(
      confirmAction.value === 'delete' ? '设备已删除' : '设备已解绑',
      'success'
    )

    // 刷新列表
    await refreshDevices()
    
    showConfirm.value = false
    selectedDevice.value = null
  } catch (error) {
    console.error('Action failed:', error)
    showToast(error.response?.data || '操作失败', 'error')
  } finally {
    actionLoading.value = false
  }
}

// 重命名功能
const showRenameDialog = (device) => {
  selectedDevice.value = device
  renameInput.value = device.display_name || device.hostname || ''
  showRename.value = true
}

const executeRename = async () => {
  if (!selectedDevice.value) return

  actionLoading.value = true
  try {
    await axios.post('/api/device/rename', {
      device_id: selectedDevice.value.id,
      display_name: renameInput.value.trim()
    })

    showToast('设备重命名成功', 'success')
    await refreshDevices()
    showRename.value = false
  } catch (error) {
    console.error('Rename failed:', error)
    showToast(error.response?.data || '重命名失败', 'error')
  } finally {
    actionLoading.value = false
  }
}

// 锁定/解锁功能
const toggleLock = async (device) => {
  const action = device.is_locked ? '解锁' : '锁定'
  
  try {
    await axios.post('/api/device/lock', {
      device_id: device.id,
      is_locked: !device.is_locked
    })

    showToast(`设备${action}成功`, 'success')
    await refreshDevices()
  } catch (error) {
    console.error('Lock toggle failed:', error)
    showToast(error.response?.data || `${action}失败`, 'error')
  }
}

// 分组功能
const showGroupDialog = (device) => {
  selectedDevice.value = device
  groupInput.value = device.device_group || ''
  showGroup.value = true
}

const executeSetGroup = async () => {
  if (!selectedDevice.value) return

  actionLoading.value = true
  try {
    await axios.post('/api/device/group', {
      device_id: selectedDevice.value.id,
      device_group: groupInput.value.trim()
    })

    showToast('设备分组设置成功', 'success')
    await refreshDevices()
    showGroup.value = false
  } catch (error) {
    console.error('Set group failed:', error)
    showToast(error.response?.data || '设置分组失败', 'error')
  } finally {
    actionLoading.value = false
  }
}

// 转移功能
const showTransferDialog = (device) => {
  if (device.is_locked) {
    showToast('设备已锁定，无法转移', 'error')
    return
  }
  selectedDevice.value = device
  transferUserId.value = null
  showTransfer.value = true
}

const executeTransfer = async () => {
  if (!selectedDevice.value || !transferUserId.value) return

  actionLoading.value = true
  try {
    await axios.post('/api/device/transfer', {
      device_id: selectedDevice.value.id,
      target_user_id: transferUserId.value
    })

    showToast('设备转移成功', 'success')
    await refreshDevices()
    showTransfer.value = false
  } catch (error) {
    console.error('Transfer failed:', error)
    showToast(error.response?.data || '转移失败', 'error')
  } finally {
    actionLoading.value = false
  }
}

// 硬件信息功能
const showHardwareInfo = (device) => {
  selectedDevice.value = device
  
  // 解析硬件指纹
  if (device.hardware_fingerprint) {
    try {
      hardwareFingerprint.value = JSON.parse(device.hardware_fingerprint)
    } catch (e) {
      hardwareFingerprint.value = null
    }
  } else {
    hardwareFingerprint.value = null
  }
  
  showHardware.value = true
}

onMounted(() => {
  refreshDevices()
})
</script>

<style scoped>
@keyframes slide-up {
  from {
    transform: translateY(100%);
    opacity: 0;
  }
  to {
    transform: translateY(0);
    opacity: 1;
  }
}

.animate-slide-up {
  animation: slide-up 0.3s ease-out;
}
</style>
