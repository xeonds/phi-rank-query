<template>
    <div class="p-4">
        <h1 class="text-2xl font-bold mb-4 text-gray-400">Rks手动计算</h1>
        <div class="flex items-center mb-4">
            <el-button class="px-4 py-2 text-white btn btn-outline btn-primary"
                @click="editVisible = true">编辑b19</el-button>
            <span class="mx-4"></span>
            <button class="px-4 py-2 text-white btn btn-outline btn-primary" @click="calc">开算！</button>
        </div>
        <div v-if="!editVisible" class="flex flex-row justify-around">
            <div class="flex flex-col px-8">
                <SongItem v-if="true" :index="`#phi`" :song="b0" />
                <template v-for="(_, index) in b19">
                    <SongItem v-if="index % 2 == 0" :index="`#${index}`" :song="rankTable[index]" />
                </template>
            </div>
            <div class="flex flex-col px-8">
                <template v-for="(_, index) in b19">
                    <SongItem v-if="index % 2 == 1" :index="`#${index}`" :song="rankTable[index]" />
                </template>
            </div>
        </div>
        <div v-else class="flex items-center mb-4 *:mx-2">
            <button @click="curr < 19 ? curr++ : curr" class="btn btn-outline btn-primary"> + </button>
            <button @click="curr > 0 ? curr-- : curr" class="btn btn-outline btn-primary"> - </button>
            <span class="mx-4">{{ curr }}</span>
            <select class="select select-ghost select-bordered w-full max-w-xs" v-model="b19[curr].song"
                placeholder="选择歌曲id">
                <option v-for="(item, index) in rankTable" :value="index">{{ item.title }}</option>
            </select>
            <select class="select select-ghost select-bordered w-full max-w-xs" v-model="b19[curr].difficulty"
                placeholder="选择难度">
                <option v-for="item in songDifficulties" :value="item">{{ item }}</option>
            </select>
            <input type="text" v-model="b19[curr].acc" placeholder="输入准度" class="input input-ghost input-bordered" />
        </div>
    </div>
</template>

<script lang="ts" setup>
import SongItem from '@/components/song.vue';
import { Ref, onMounted, ref, watch } from 'vue';
import { Song } from '@/common';

const editVisible = ref(true);
const curr: Ref<number> = ref(0);
const rankTable: Ref<any[]> = ref([]);
const b0: Ref<Song> = ref({} as Song);
const b19: Ref<Song[]> = ref(Array.from({ length: 20 }, () => ({} as Song)));
const songDifficulties: Ref<any> = ref(null);

watch(() => b19,
    () => {
        // fix: 996
        songDifficulties.value = Object.values(rankTable.value[b19.value[curr.value].song as any])
            .filter((item: any) => item != '' && !isNaN(parseFloat(item)))
            .map((item: any) => parseFloat(item))
            .sort((a: any, b: any) => b - a);
    },
    { deep: true }
)

const calc = async () => {
    const response = await fetch('/api/v1/calc', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            b19: b19.value
        })
    });
    const data = response.json();
    console.log(data);
}

onMounted(async () => {
    const response = await fetch('/api/v1/rank_table');
    const data = await response.json();
    rankTable.value = data;
})
</script>