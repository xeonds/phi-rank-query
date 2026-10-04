<template>
  <div class="relative min-h-screen">
    <div class="fixed inset-0 z-0 bg-cover bg-center blur-md"
      style="background-image: url('/assets/Star1.png')"></div>
    <div class="fixed inset-0 z-0 bg-gray-900/60"></div>

    <div class="relative z-10">
      <header
        class="sticky top-0 z-30 flex items-center justify-between gap-2 border-b border-white/10 bg-black/40 px-4 py-3 backdrop-blur-md sm:px-6">
        <h1
          class="truncate bg-gradient-to-r from-sky-300 via-blue-200 to-indigo-300 bg-clip-text text-xl font-bold text-transparent sm:text-2xl">
          Phi Rank Query
        </h1>

        <el-menu class="app-menu hidden md:flex" mode="horizontal" router :default-active="$route.path"
          :ellipsis="false">
          <el-menu-item v-for="item in nav" :key="item.path" :index="item.path">{{ item.label }}</el-menu-item>
        </el-menu>

        <el-button class="md:hidden" text @click="drawer = true">☰ 菜单</el-button>
      </header>

      <el-drawer v-model="drawer" direction="rtl" size="70%" :with-header="false">
        <el-menu mode="vertical" router :default-active="$route.path" @select="drawer = false">
          <el-menu-item v-for="item in nav" :key="item.path" :index="item.path">{{ item.label }}</el-menu-item>
        </el-menu>
      </el-drawer>

      <main class="mx-auto max-w-6xl px-3 py-4 sm:px-4 sm:py-6">
        <router-view />
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRoute } from 'vue-router'

const $route = useRoute()
const drawer = ref(false)

const nav = [
  { path: '/', label: '主页' },
  { path: '/info', label: '定数查询' },
  { path: '/b19', label: 'B27 查询' },
  { path: '/bn', label: '所有成绩' },
  { path: '/calc', label: 'Rks 计算' },
  { path: '/history', label: '查询历史' },
  { path: '/leaderboard', label: '排行榜' },
  { path: '/session', label: 'Session' },
]
</script>

<style>
body {
  margin: 0;
}
.app-menu {
  --el-menu-bg-color: transparent;
  --el-menu-text-color: #cbd5e1;
  --el-menu-active-color: #60a5fa;
  --el-menu-hover-bg-color: rgba(255, 255, 255, 0.08);
  border-bottom: none !important;
}
</style>