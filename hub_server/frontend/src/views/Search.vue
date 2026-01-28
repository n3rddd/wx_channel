<template>
  <div class="min-h-screen bg-bg p-8 lg:p-12 font-sans text-text">
    <header class="flex justify-between items-center mb-12">
      <div>
        <h1 class="font-serif font-bold text-3xl lg:text-4xl mb-2 text-text">远程穿透搜索</h1>
        <p class="text-text-muted">通过选中终端执行实时视频号搜索与解析</p>
      </div>
      <div v-if="client" class="px-4 py-2 rounded-xl bg-bg shadow-neu-sm border border-white/50 text-primary font-medium flex items-center gap-2">
        <span class="text-xs uppercase tracking-wider text-text-muted">Connected to</span>
        <strong>{{ client.hostname }}</strong>
      </div>
    </header>

    <!-- Client Selector if none selected -->
    <div v-if="!client" class="p-12 text-center bg-bg rounded-[2rem] shadow-neu">
      <p class="text-text-muted mb-4">请先选择一个操作目标</p>
      <router-link to="/dashboard" class="inline-block px-6 py-3 rounded-full bg-bg shadow-neu-btn text-primary font-semibold hover:text-primary-dark transition-all active:shadow-neu-btn-active">
          前往在线终端
      </router-link>
    </div>

    <div v-else>
      <!-- Search Box -->
      <div class="flex gap-4 mb-12 p-4 bg-bg rounded-[2rem] shadow-neu items-center">
        <input 
          v-model="keyword" 
          type="text" 
          class="flex-1 bg-transparent border-none outline-none text-lg px-4 text-text placeholder-text-muted/50" 
          placeholder="输入视频号名称、博主ID或关键词..."
          @keyup.enter="handleSearch"
        >
        <button 
            class="px-8 py-3 rounded-full bg-primary text-white font-semibold shadow-lg shadow-primary/30 hover:bg-primary-dark transition-all transform hover:-translate-y-0.5 disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
            @click="handleSearch" 
            :disabled="searching"
        >
          <Search v-if="!searching" class="w-5 h-5" />
          <div v-else class="w-5 h-5 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
          <span>{{ searching ? '搜索中...' : '开始搜索' }}</span>
        </button>
      </div>

      <!-- Results Grid -->
      <div v-if="results.length > 0" class="grid grid-cols-[repeat(auto-fill,minmax(200px,1fr))] gap-8">
        <div 
          v-for="item in results" 
          :key="item.username" 
          class="bg-bg rounded-3xl p-6 text-center shadow-neu border border-white/40 cursor-pointer transition-all hover:-translate-y-1 hover:shadow-neu-sm hover:border-primary/20 group"
          @click="openProfile(item)"
        >
          <div class="w-20 h-20 mx-auto rounded-full bg-bg shadow-neu-sm p-1 mb-4 group-hover:shadow-neu-pressed transition-shadow">
             <img :src="item.headUrl || placeholderImg" class="w-full h-full rounded-full object-cover" @error="onImgError">
          </div>
          <div class="font-bold text-lg mb-2 text-text group-hover:text-primary transition-colors">{{ item.nickname }}</div>
          <div class="text-sm text-text-muted line-clamp-2 px-2">{{ item.signature || '暂无签名' }}</div>
        </div>
      </div>
      
      <div v-if="hasMoreSearch" class="text-center mt-12 mb-8">
          <button class="px-8 py-3 rounded-full bg-bg shadow-neu-btn text-text-muted font-medium hover:text-primary transition-all active:shadow-neu-btn-active disabled:opacity-50" @click="handleSearch(true)" :disabled="searching">
              {{ searching ? '加载中...' : '加载更多账号' }}
          </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useClientStore } from '../store/client'
import { useRouter } from 'vue-router'
import { Search } from 'lucide-vue-next'

const clientStore = useClientStore()
const router = useRouter()
const client = computed(() => clientStore.currentClient)

const keyword = ref('')
const searching = ref(false)
const results = ref([])
const placeholderImg = 'https://via.placeholder.com/100'

const lastSearchBuffer = ref('')
const hasMoreSearch = ref(false)

const handleSearch = async (loadMore = false) => {
  if (!keyword.value || !client.value) return
  searching.value = true
  if (!loadMore) {
      results.value = []
      lastSearchBuffer.value = ''
      hasMoreSearch.value = false
  }
  
  try {
    const res = await clientStore.remoteCall('api_call', {
      key: 'key:channels:contact_list',
      body: { 
          keyword: keyword.value,
          lastBuffer: loadMore ? lastSearchBuffer.value : ''
      }
    })
    
    // Config adapter
    const findList = (obj) => {
        if (!obj) return null
        if (Array.isArray(obj.infoList)) return obj.infoList
        if (Array.isArray(obj.objectList)) return obj.objectList
        if (Array.isArray(obj.list)) return obj.list
        return null
    }

    let list = null
    // 1. Try res.data (Hub payload -> data)
    if (res.data) {
        list = findList(res.data)
        // Check for pagination in data
        if (res.data.continueFlag !== undefined) {
             hasMoreSearch.value = !!res.data.continueFlag
             lastSearchBuffer.value = res.data.lastBuff || ''
        }

        // 2. Try res.data.data (Hub payload -> data -> business payload)
        if (!list && res.data.data) {
            list = findList(res.data.data)
             // Check for nested pagination
            if (res.data.data.continueFlag !== undefined) {
                 hasMoreSearch.value = !!res.data.data.continueFlag
                 lastSearchBuffer.value = res.data.data.lastBuff || ''
            }
        }
    }
    // 3. Try root just in case
    if (!list) {
        list = findList(res)
        if (res.continueFlag !== undefined) {
            hasMoreSearch.value = !!res.continueFlag
            lastSearchBuffer.value = res.lastBuff || ''
        }
    }

    if (!list) {
         console.warn("No list found in response", res)
         list = []
    }
    
    const newItems = list.map(item => item.contact || item)
    if (loadMore) {
        results.value = [...results.value, ...newItems]
    } else {
        results.value = newItems
    }

  } catch (err) {
    alert('搜索失败: ' + err.message)
  } finally {
    searching.value = false
  }
}

const openProfile = (user) => {
    router.push({
        path: '/profile',
        query: {
            username: user.username,
            nickname: user.nickname,
            headUrl: user.headUrl,
            signature: user.signature
        }
    })
}

const onImgError = (e) => {
  e.target.src = placeholderImg
}
</script>
