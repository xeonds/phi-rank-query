<script setup lang="ts">
import { computed } from 'vue'
import type { SongRecord } from '@/types'
import { getRating } from '@/utils'

const props = defineProps<{ index: string; song: SongRecord }>()

const rating = computed(() => getRating(props.song.fullCombo, props.song.score))
const acc = computed(() => (props.song.acc + 0.005).toFixed(2))
const rankClass = computed(() => ({
  EZ: 'bg-green-600',
  HD: 'bg-sky-600',
  IN: 'bg-red-600',
  AT: 'bg-zinc-600',
  Legacy: 'bg-purple-600',
}[props.song.level] ?? 'bg-zinc-600'))
</script>

<template>
  <div
    class="my-2 flex w-full items-stretch overflow-hidden rounded-lg bg-black/45 shadow-lg ring-1 ring-white/10 backdrop-blur-sm transition hover:ring-white/25">
    <div class="relative w-28 shrink-0 sm:w-40">
      <img loading="lazy" :src="song.illustration" :alt="song.song" class="h-20 w-28 object-cover sm:h-24 sm:w-40" />
      <span
        class="absolute left-1 top-1 rounded bg-black/70 px-1.5 text-xs font-semibold text-gray-100 backdrop-blur">{{ index }}</span>
      <span :class="['absolute bottom-0 left-1 rounded-tr px-1.5 py-0.5 text-[11px] font-semibold text-white sm:px-2 sm:text-xs', rankClass]">
        {{ song.level }} {{ Number(song.difficulty).toFixed(1) }}
      </span>
      <span class="absolute bottom-0 right-1 text-xs font-bold text-white drop-shadow sm:text-sm">
        {{ song.rks.toFixed(2) }}
      </span>
    </div>
    <div class="flex min-w-0 flex-1 items-center justify-between gap-2 px-2 sm:px-3">
      <p class="truncate text-xs text-gray-100 sm:text-sm" :title="song.song">{{ song.song }}</p>
      <div class="flex shrink-0 items-center gap-1.5 sm:gap-2">
        <img :src="`/assets/${rating}.png`" :alt="rating" class="h-6 w-6 sm:h-7 sm:w-7" />
        <div class="text-right leading-tight">
          <p class="text-base font-semibold text-white sm:text-lg">{{ song.score }}</p>
          <p class="text-[11px] text-gray-300 sm:text-xs">{{ acc }}%</p>
        </div>
      </div>
    </div>
  </div>
</template>