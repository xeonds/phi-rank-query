<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { fetchRankTable } from '@/api'
import { songRks } from '@/utils'
import type { RankItem, SongRecord } from '@/types'
import SongCard from '@/components/SongCard.vue'

interface Entry {
  id: string
  level: string
  acc: number
  score: number
}

const rankTable = ref<RankItem[]>([])
const entries = ref<Entry[]>([])
const songIndex = computed(() => new Map(rankTable.value.map((r) => [r.id, r])))

const levelsOf = (id: string): string[] => {
  const item = songIndex.value.get(id)
  if (!item) return []
  return (['EZ', 'HD', 'IN', 'AT', 'Legacy'] as const).filter((l) => item[l])
}

const rankOf = (id: string, level: string): number => {
  const item = songIndex.value.get(id)
  return item ? parseFloat(item[level as keyof RankItem]) : NaN
}

const computed_ = computed(() => {
  const records: SongRecord[] = entries.value
    .filter((e) => e.id && e.level)
    .map((e) => {
      const item = songIndex.value.get(e.id)
      return {
        id: e.id,
        song: item?.title ?? e.id,
        level: e.level,
        difficulty: String(rankOf(e.id, e.level)),
        rks: songRks(e.acc, rankOf(e.id, e.level)),
        score: e.score,
        acc: e.acc,
        fullCombo: false,
        illustration: `/assets/illustrations/${e.id}.png`,
      }
    })
    .sort((a, b) => b.rks - a.rks)

  const b0 = records.find((r) => r.acc === 100)
  let rks = records.slice(0, 19).reduce((sum, r) => sum + r.rks / 20, 0)
  if (b0) rks += b0.rks / 20
  return { records, rks }
})

function addEntry() {
  entries.value.push({ id: '', level: '', acc: 100, score: 1000000 })
}

function exportData() {
  const blob = new Blob([JSON.stringify({ data: entries.value })], { type: 'application/json' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = 'calc.json'
  a.click()
  URL.revokeObjectURL(url)
}

function importData() {
  const input = document.createElement('input')
  input.type = 'file'
  input.accept = '.json'
  input.onchange = (e) => {
    const file = (e.target as HTMLInputElement).files?.[0]
    if (!file) return
    const reader = new FileReader()
    reader.onload = () => {
      try {
        entries.value = JSON.parse(reader.result as string).data ?? []
      } catch {
        /* ignore */
      }
    }
    reader.readAsText(file)
  }
  input.click()
}

onMounted(async () => {
  rankTable.value = await fetchRankTable()
})
</script>

<template>
  <div>
    <div class="mb-4 flex flex-wrap items-center gap-2">
      <h1 class="text-2xl font-bold text-gray-200">Rks 手动计算</h1>
      <div class="ml-auto flex gap-2">
        <el-button @click="addEntry">添加一行</el-button>
        <el-button @click="importData">导入</el-button>
        <el-button @click="exportData">导出</el-button>
      </div>
    </div>

    <div class="frost mb-4 flex flex-wrap items-center gap-x-6 gap-y-1 px-4 py-3 text-gray-100">
      <span>Player: <span class="font-semibold text-white">OFFLINE</span></span>
      <span>RankingScore: <span class="font-semibold text-sky-300">{{ computed_.rks.toFixed(4) }}</span></span>
    </div>

    <div class="frost mb-4 space-y-3 p-3 sm:p-4">
      <p v-if="entries.length === 0" class="text-sm text-gray-400">还没有条目，点击「添加一行」开始。</p>
      <div v-for="(entry, i) in entries" :key="i" class="flex flex-wrap items-center gap-2">
        <el-select v-model="entry.id" filterable placeholder="选择歌曲" class="w-full sm:w-64">
          <el-option v-for="r in rankTable" :key="r.id" :label="r.title" :value="r.id" />
        </el-select>
        <el-select v-model="entry.level" placeholder="难度" class="w-full sm:w-28">
          <el-option v-for="l in levelsOf(entry.id)" :key="l" :label="l" :value="l" />
        </el-select>
        <el-input-number v-model="entry.acc" :min="0" :max="100" :step="0.01" :precision="2" />
        <el-input-number v-model="entry.score" :min="0" :max="1000000" :step="1" />
        <el-button type="danger" @click="entries.splice(i, 1)">删除</el-button>
      </div>
    </div>

    <div class="grid grid-cols-1 gap-x-4 lg:grid-cols-2">
      <SongCard
        v-for="(song, i) in computed_.records"
        :key="i"
        :index="'#' + (i + 1)"
        :song="song"
      />
    </div>
  </div>
</template>