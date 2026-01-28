<template>
  <div class="min-h-screen w-full bg-bg flex flex-col font-sans">
    <Navbar />
    <main class="flex-1 relative w-full overflow-hidden">
      <!-- Adaptive container: full width on mobile, centered max-width on larger screens if needed, or just fluid -->
      <div class="h-full w-full overflow-y-auto px-4 md:px-8 py-6">
        <div class="max-w-7xl mx-auto h-full"> 
          <router-view v-slot="{ Component }">
            <transition name="fade" mode="out-in">
              <component :is="Component" />
            </transition>
          </router-view>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup>
import Navbar from '../components/Navbar.vue'
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
