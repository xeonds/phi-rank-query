<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { loadHistory, removeHistory, saveHistory, sessionStore } from '@/sessions'
import type { PlayerResult } from '@/types'
import SongCard from '@/components/SongCard.vue'

const router = useRouter()
const entries = ref<PlayerResult[]>([])
const viewing = ref<PlayerResult | null>(null)

function refresh() {
  entries.value = loadHistory(sessionStore.selected)
}

function exportHistory() {
  const blob = new Blob([JSON.stringify({ data: entries.value })], { type: 'application/json' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = 'history.json'
  a.click()
  URL.revokeObjectURL(url)
}

function importHistory() {
  const input = document.createElement('input')
  input.type = 'file'
  input.accept = '.json'
  input.onchange = (e) => {
    const file = (e.target as HTMLInputElement).files?.[0]
    if (!file) return
    const reader = new FileReader()
    reader.onload = () => {
      try {
        const data = JSON.parse(reader.result as string)
        entries.value = data.data ?? []
        saveHistory(sessionStore.selected, entries.value)
        ElMessage.success('已导入')
      } catch {
        ElMessage.error('导入失败：文件格式错误')
      }
    }
    reader.readAsText(file)
  }
  input.click()
}

async function del(index: number) {
  await ElMessageBox.confirm('确定删除该记录吗？', '提示', { type: 'warning' })
  removeHistory(sessionStore.selected, index)
  refresh()
}

onMounted(() => {
  if (!sessionStore.selected) {
    router.replace('/session')
    return
  }
  refresh()
})
</script>

<template>
  <div>
    <div class="mb-4 flex flex-wrap items-center gap-2">
      <h1 class="text-2xl font-bold text-gray-200">查询历史</h1>
      <div class="ml-auto flex gap-2">
        <el-button @click="importHistory">导入历史</el-button>
        <el-button @click="exportHistory">导出历史</el-button>
      </div>
    </div>

    <template v-if="viewing">
      <el-button class="mb-4" @click="viewing = null">返回列表</el-button>
      <p class="mb-2 text-gray-200">玩家: <span class="font-semibold text-white">{{ viewing.player }}</span></p>
      <p class="mb-4 text-gray-200">Rks: <span class="font-semibold text-sky-300">{{ viewing.rks.toFixed(4) }}</span></p>
      <div class="grid grid-cols-1 gap-x-4 lg:grid-cols-2">
        <SongCard
          v-for="(song, i) in [...viewing.phi, ...viewing.records]"
          :key="i"
          :index="'#' + (i + 1)"
          :song="song"
        />
      </div>
    </template>

    <template v-else>
      <p class="mb-3 text-sm text-gray-400">共 {{ entries.length }} 条记录</p>
      <div class="frost p-2 sm:p-3">
        <el-table :data="entries" stripe>
          <el-table-column prop="date" label="日期" />
          <el-table-column label="Rks" width="140">
            <template #default="{ row }">{{ row.rks.toFixed(4) }}</template>
          </el-table-column>
          <el-table-column label="操作" width="200">
            <template #default="{ row, $index }">
              <el-button size="small" type="primary" @click="viewing = row">查看</el-button>
              <el-button size="small" type="danger" @click="del($index)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </template>
  </div>
</template>