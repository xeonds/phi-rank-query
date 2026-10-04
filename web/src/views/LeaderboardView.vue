<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { fetchLeaderboard } from '@/api'
import type { LeaderboardItem } from '@/types'

const data = ref<LeaderboardItem[]>([])
const loading = ref(true)

onMounted(async () => {
  try {
    data.value = await fetchLeaderboard()
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div>
    <h1 class="mb-4 text-2xl font-bold text-gray-200">排行榜</h1>
    <div class="frost p-2 sm:p-3">
      <el-table v-loading="loading" :data="data" stripe>
        <el-table-column type="index" label="排行" width="80" />
        <el-table-column prop="username" label="玩家" />
        <el-table-column label="Rks" width="140">
          <template #default="{ row }">{{ row.rks.toFixed(4) }}</template>
        </el-table-column>
      </el-table>
    </div>
  </div>
</template>