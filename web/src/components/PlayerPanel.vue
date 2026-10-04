<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { fetchPlayer } from '@/api'
import { appendHistory, sessionStore } from '@/sessions'
import { challengeMode, challengeRank, formatDate } from '@/utils'
import type { PlayerResult } from '@/types'
import SongCard from '@/components/SongCard.vue'

const props = defineProps<{ all?: boolean }>()

const router = useRouter()
const loading = ref(true)
const error = ref('')
const player = ref<PlayerResult | null>(null)

const displayRecords = computed(() => {
  if (!player.value) return []
  const records = props.all ? player.value.records : player.value.records.slice(0, 27)
  return [...player.value.phi, ...records]
})

onMounted(async () => {
  if (!sessionStore.selected) {
    router.replace('/session')
    return
  }
  try {
    const result = await fetchPlayer(sessionStore.selected)
    player.value = result
    appendHistory(sessionStore.selected, result)
  } catch (e) {
    error.value = e instanceof Error ? e.message : String(e)
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div v-loading="loading" class="min-h-[200px]">
    <el-result v-if="error" icon="error" :title="error" sub-title="请确认 Session 有效" />
    <template v-else-if="player">
      <div class="frost mb-4 p-4 text-gray-200">
        <div class="flex flex-wrap items-center justify-between gap-4">
          <div class="leading-7">
            <p>玩家: <span class="font-semibold text-white">{{ player.player }}</span></p>
            <p>Rks: <span class="font-semibold text-sky-300">{{ player.rks.toFixed(4) }}</span></p>
            <p>同步时间: {{ formatDate(player.date) }}</p>
          </div>
          <div class="flex items-center gap-2">
            <span>课题模式:</span>
            <img v-if="challengeMode(player.challenge)" :src="`/assets/${challengeMode(player.challenge)}.png`" class="h-5" />
            <span>{{ challengeRank(player.challenge) }}</span>
          </div>
        </div>
      </div>
      <div class="grid grid-cols-1 gap-x-4 lg:grid-cols-2">
        <SongCard v-for="(song, i) in displayRecords" :key="i" :index="'#' + (i + 1)" :song="song" />
      </div>
    </template>
  </div>
</template>