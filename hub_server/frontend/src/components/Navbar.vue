<template>
  <nav class="w-full bg-white shadow-neu-sm z-50 sticky top-0 px-6 md:px-12 py-4 flex items-center justify-between shrink-0">
    <!-- Brand -->
    <div class="flex items-center gap-3">
      <div class="w-10 h-10 rounded-xl bg-primary/10 flex items-center justify-center text-primary shrink-0">
        <svg viewBox="0 0 24 24" width="22" height="22" fill="currentColor">
          <path d="M12 2L4.5 20.29l.71.71L12 18l6.79 3 .71-.71z" />
        </svg>
      </div>
      <span class="font-serif text-xl font-bold text-slate-800 tracking-tight">Hub Control</span>
    </div>

    <!-- Navigation Links -->
    <div class="flex items-center gap-1 md:gap-4">
      <router-link to="/dashboard" active-class="!text-primary !bg-primary/5" class="px-4 py-2 rounded-lg text-slate-600 font-medium text-sm transition-all hover:text-primary hover:bg-slate-50">
        概览
      </router-link>
      <router-link to="/search" active-class="!text-primary !bg-primary/5" class="px-4 py-2 rounded-lg text-slate-600 font-medium text-sm transition-all hover:text-primary hover:bg-slate-50">
        搜索
      </router-link>
      <router-link to="/tasks" active-class="!text-primary !bg-primary/5" class="px-4 py-2 rounded-lg text-slate-600 font-medium text-sm transition-all hover:text-primary hover:bg-slate-50">
        任务
      </router-link>
      <router-link to="/settings" active-class="!text-primary !bg-primary/5" class="px-4 py-2 rounded-lg text-slate-600 font-medium text-sm transition-all hover:text-primary hover:bg-slate-50">
        设置
      </router-link>
    </div>

    <!-- Status Pill -->
    <div class="hidden md:flex items-center gap-2 px-3 py-1.5 rounded-full bg-green-50 border border-green-100">
      <span class="w-2 h-2 rounded-full bg-green-500 animate-pulse"></span>
      <span class="text-green-700 font-bold text-xs uppercase tracking-wide">系统在线</span>
    </div>
  </nav>
</template>

<script setup>
import { LayoutDashboard, Globe, ListTodo, Settings } from 'lucide-vue-next'
</script>

<style scoped>
/* Hide scrollbar for cleaner look */
.no-scrollbar::-webkit-scrollbar {
  display: none;
}
.no-scrollbar {
  -ms-overflow-style: none;
  scrollbar-width: none;
}
</style>
