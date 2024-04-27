<template>
    <h1 class="text-2xl font-bold mb-4 text-gray-400">Rks手动计算</h1>
    <div class="flex items-center mb-4">
        <el-button class="px-4 py-2 text-white btn btn-outline btn-primary"
            @click="editVisible = true">编辑b19</el-button>
        <span class="mx-4"></span>
        <button class="px-4 py-2 text-white btn btn-outline btn-primary" @click="editVisible = false">开算！</button>
    </div>
    <div v-if="!editVisible" class="flex flex-col justify-around w-full">
        <p>Player: OFFLINE</p>
        <p>RankingScore: {{ calc(b19).rks.toFixed(4) }}</p>
        <div class="flex flex-row flex-wrap">
            <SongItem :index="'#phi'" v-if="calc(b19).b0" :song="calc(b19).b0" class="mx-4" />
            <SongItem v-for="song, index in calc(b19).b19" :key="index" :index="'#' + (index + 1).toString()"
                :song="song" class="mx-4" />
        </div>
    </div>
    <div v-else class="flex flex-col items-center pb-8 *:mx-2">
        <div class="*:mx-2 *:my-4">
            <button @click="b19.push({} as Song)" class="btn btn-outline btn-primary"> + </button>
            <span>{{ b19.length }}</span>
            <button @click="b19.pop()" class="btn btn-outline btn-primary"> - </button>
        </div>
        <div v-for="item in b19" class="flex flex-row *:mx-2 *:my-2 w-full">
            <select class="select select-ghost select-bordered w-full max-w-xs" v-model="item.id" placeholder="选择歌曲id">
                <option v-for="song in rankTable" :value="song.id">{{ song.title }}</option>
            </select>
            <select class="select select-ghost select-bordered w-full max-w-xs" v-model="item.difficulty"
                placeholder="选择难度">
                <template
                    v-for="difficulty, rank in item.id ? rankTable.filter(e => e.id == item.id!)[0].difficulty : []">
                    <option v-if="difficulty" :value="rank">{{ rank }}</option>
                </template>
            </select>
            <input type="number" v-model="item.acc" placeholder="输入准度" class="input input-ghost input-bordered" />
            <input type="number" v-model="item.score" placeholder="输入成绩" class="input input-ghost input-bordered" />
        </div>
    </div>
</template>

<script lang="ts" setup>
import SongItem from '@/components/song.vue';
import { Ref, onMounted, ref } from 'vue';
import { Song, getRating } from '@/common';

interface RankTableItem {
    id: string;
    title: string;
    difficulty: {
        'EZ': number;
        'HD': number;
        'IN': number;
        'AT': number;
    };
}

const editVisible = ref(true);
const rankTable: Ref<RankTableItem[]> = ref([]);
const b19: Ref<Song[]> = ref([]);

const calc = (b19: Song[]) => {
    if (b19.length == 0) { return { b19: [], b0: {}, rks: 0 } }
    const getTitle = (id?: string) => { return id ? (rankTable.value.filter(e => e.id == id))[0].title : '' }
    const getDifficulty = (id?: string) => { return id ? (rankTable.value.filter(e => e.id == id))[0].difficulty : [] as any }
    b19 = b19.map(item => ({
        id: item.id,
        song: getTitle(item.id),
        difficulty: getDifficulty(item.id)[item.difficulty],
        acc: item.acc,
        score: item.score,
        rank: item.difficulty,
        rks: calcSongRks(parseFloat(item.acc), parseFloat(getDifficulty(item.id)[item.difficulty])).toString(),
        Rating: getRating(false, parseInt(item.score)),
        illustration: `/assets/illustrations/${item.id}.png`,
    } as Song))
    b19.sort((a: any, b: any) => parseFloat(b.rks) - parseFloat(a.rks))
    const b0 = b19.filter((item) => parseFloat(item.acc) == 100)[0] ?? null;
    let rks = b19.slice(0, 19)
        .map(item => parseFloat(item.rks) / 20)
        .reduce((sum: number, curr: number) => sum + curr)
    if (b0) { rks += parseFloat(b0.rks) / 20 }
    return { b19, b0, rks }
}

const calcSongRks = (acc: number, rank: number) => {
    if (acc == 100) { return rank }
    else if (acc < 70) { return 0 }
    else { return rank * (((acc - 55) / 45) * ((acc - 55) / 45)) }
}

onMounted(async () => {
    const response = await fetch('/api/v1/rank_table');
    const data = await response.json();
    rankTable.value = Object.entries(data)
        .map(([key, item]: [any, any]) => ({
            id: key,
            title: item.title,
            difficulty: {
                'EZ': parseFloat(item.EZ),
                'HD': parseFloat(item.HD),
                'IN': parseFloat(item.IN),
                'AT': parseFloat(item.AT) || NaN
            },
        }))
})
</script>