<template>
    <div class="p-4">
        <h1 class="text-2xl font-bold mb-4 text-gray-400">Rks手动计算</h1>
        <div class="flex items-center mb-4">
            <el-button class="px-4 py-2 text-white btn btn-outline btn-primary"
                @click="editVisible = true">编辑b19</el-button>
            <span class="mx-4"></span>
            <button class="px-4 py-2 text-white btn btn-outline btn-primary" @click="editVisible = false">开算！</button>
        </div>
        <div v-if="!editVisible" class="flex flex-row justify-around">
            <div class="flex flex-col px-8">
                <template v-for="(_, index) in rankTable">
                    <SongItem v-if="index % 2 == 0" :index="`#${index}`" :song="rankTable[index]" />
                </template>
            </div>
            <div class="flex flex-col px-8">
                <template v-for="(_, index) in rankTable">
                    <SongItem v-if="index % 2 == 1" :index="`#${index}`" :song="rankTable[index]" />
                </template>
            </div>
        </div>
        <div v-else class="flex items-center mb-4">
            <input type="text" v-model="curr" placeholder="选择id"
                class="input input-ghost input-bordered w-full max-w-xs mr-2 p-2 border border-gray-800 rounded" />
            <input type="text" v-model="rankTable[curr].difficulty" placeholder="输入歌曲定数"
                class="input input-ghost input-bordered w-full max-w-xs mr-2 p-2 border border-gray-800 rounded" />
            <input type="text" v-model="rankTable[curr].acc" placeholder="输入歌曲准度"
                class="input input-ghost input-bordered w-full max-w-xs mr-2 p-2 border border-gray-800 rounded" />
        </div>
    </div>
</template>

<script lang="ts" setup>
import SongItem from '@/components/song.vue';
import { Ref, onMounted, ref } from 'vue';

interface RankTableItem {
    num: number;
    song: string;
    illustration: string;
    rank: string;
    difficulty: string;
    rks: string;
    Rating: string;
    score: string;
    acc: string;
    suggest: string;
}
const editVisible = ref(false);
const curr: Ref<number> = ref(0);
const rankTable: Ref<RankTableItem[]> = ref([]);

onMounted(() => {
    for (let i = 0; i < 20; i++) {
        rankTable.value.push({ song: 'test', difficulty: '', rank: '0' } as RankTableItem);
    }
})
</script>