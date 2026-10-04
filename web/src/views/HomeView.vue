<script setup lang="ts">
import { onMounted, ref } from 'vue'

const version = ref('...')
const updateTime = ref('...')

onMounted(async () => {
  try {
    const resp = await fetch('/config.json')
    const json = await resp.json()
    version.value = json.version
    updateTime.value = json.updateTime
  } catch {
    /* ignore */
  }
})

const features = [
  { title: '多存档查询', desc: 'SessionToken 保存在浏览器本地，随时切换要查询的存档。' },
  { title: '全成绩查询', desc: '一次性查询你的所有历史成绩。' },
  { title: '谱面定数查询', desc: '在线查看最新的谱面定数表。' },
  { title: '排行榜', desc: '看看大家的 Rks 排名。' },
]
</script>

<template>
  <div>
    <h1
      class="pb-6 bg-gradient-to-r from-sky-300 via-blue-200 to-indigo-300 bg-clip-text text-3xl font-bold text-transparent sm:text-4xl">
      欢迎使用 Phigros Rks 查询工具
    </h1>

    <h2 class="mb-4 text-2xl font-bold text-gray-300">Features</h2>
    <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
      <el-card v-for="f in features" :key="f.title" shadow="hover">
        <template #header><span class="font-semibold text-gray-100">{{ f.title }}</span></template>
        <p class="text-sm text-gray-300">{{ f.desc }}</p>
      </el-card>
    </div>

    <p class="mt-6 text-gray-300">
      用法：先在 Session 页面添加存档并选择，然后进入 B27 或 所有成绩 页面查询。
    </p>
    <p class="mt-4 text-gray-400">
      ps: 有任何意见和建议欢迎通过 <code class="rounded bg-white/10 px-1.5 py-0.5 text-gray-200">xeonds@stu.xidian.edu.cn</code> 联系。
    </p>
    <p class="mt-2 text-gray-400">ppps: 定数数据库已同步 {{ version }} 版本。</p>

    <footer class="mt-8 border-t border-white/10 pt-4 text-center text-gray-400">
      Last Update: {{ updateTime }} | Phi Rank Query
    </footer>
  </div>
</template>