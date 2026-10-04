<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { fetchRankTable } from '@/api'
import type { RankItem } from '@/types'

interface Row {
  title: string
  level: string
  rank: number
}

const items = ref<RankItem[]>([])
const keyword = ref('')
const loading = ref(true)

const rows = computed<Row[]>(() =>
  items.value
    .flatMap((it) =>
      (['EZ', 'HD', 'IN', 'AT'] as const).map((level) => ({
        title: it.title,
        level,
        rank: parseFloat(it[level]),
      })),
    )
    .filter((r) => !Number.isNaN(r.rank))
    .sort((a, b) => b.rank - a.rank),
)

const filtered = computed(() =>
  rows.value.filter((r) => r.title.toLowerCase().includes(keyword.value.trim().toLowerCase())),
)

onMounted(async () => {
  try {
    items.value = await fetchRankTable()
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div>
    <div class="mb-4 flex flex-wrap items-center justify-between gap-4">
      <h1 class="text-2xl font-bold text-gray-200">谱面定数表</h1>
      <el-input v-model="keyword" placeholder="搜索曲名" clearable class="w-full sm:w-64" />
    </div>

    <div class="frost p-2 sm:p-3">
      <el-table v-loading="loading" :data="filtered" height="70vh" stripe>
        <el-table-column prop="title" label="名称" min-width="220" />
        <el-table-column prop="level" label="难度" width="100" />
        <el-table-column prop="rank" label="定数" width="100" />
      </el-table>
    </div>
  </div>
</template>