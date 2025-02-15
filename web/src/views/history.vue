<template>
    <h1 class="text-2xl font-bold mb-4 text-gray-400">查询历史</h1>
    <div class="flex items-center mb-4">
        <button @click="importHistory" class="mx-2 px-4 py-2 text-white btn btn-outline btn-primary">
            导入历史
        </button>
        <button @click="exportHistory" class="mx-2 px-4 py-2 text-white btn btn-outline btn-primary">
            导出历史
        </button>
        <button v-if="!isSelectHistory" @click="isSelectHistory = true"
            class="mx-2 px-4 py-2 text-white btn btn-outline btn-primary">
            选择历史
        </button>
    </div>
    <div v-if="isSelectHistory">
        <h3>查询历史记录，共{{ historyScores.length }}条</h3>
        <p>注意：3.11.0更新了历史格式，因此旧版本记录格式将无法查看</p>
        <table class="table w-full border-collapse text-gray-400">
            <thead>
                <tr>
                    <th class="py-2 px-4 border-b text-gray-400"></th>
                    <th class="py-2 px-4 border-b text-gray-400">日期</th>
                    <th class="py-2 px-4 border-b text-gray-400">Rks</th>
                    <th class="py-2 px-4 border-b text-gray-400">操作</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="(score, index) in historyScores" :key="score">
                    <td class="py-2 px-4 border-b text-gray-400">
                        <button class="px-2 py-1 btn btn-outline btn-info" @click="checkScore(score)">查 看</button>
                    </td>
                    <td class="py-2 px-4 border-b text-gray-400">{{ score.date }}</td>
                    <td class="py-2 px-4 border-b text-gray-400">{{ score.rks }}</td>
                    <td class="py-2 px-4 border-b">
                        <button @click="deleteHistory(index)" class="px-2 py-1 btn btn-outline btn-warning">
                            删除该记录
                        </button>
                    </td>
                </tr>
            </tbody>
        </table>
    </div>
    <div v-else>
        <p>Player: {{ selectedScore.player }}</p>
        <p>RankingScore: {{ selectedScore.rks.toFixed(4) }}</p>
        <div class="flex flex-row flex-wrap">

            <SongItem v-for="(song, index) in parseSongList([...selectedScore.phi, ...(selectedScore.b19 || selectedScore.b27)])" :key="index"
                :index="'#' + (index + 1).toString()" :song="song" class="mx-4"/>
        </div>
    </div>
</template>

<script lang="ts" setup>
import { Ref, onMounted, ref } from 'vue';
import SongItem from '@/components/song.vue';
import { parseSongList } from '@/common.ts';

const historyScores: Ref<any[]> = ref([])
const sessionToken = ref<string>('');
const isSelectHistory = ref(true);
const selectedScore = ref<any>({});

const exportHistory = () => {
    const data = JSON.stringify({ data: historyScores.value });
    const blob = new Blob([data], { type: 'application/json' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'history.json';
    a.click();
    URL.revokeObjectURL(url);
}
const importHistory = () => {
    const input = document.createElement('input');
    input.type = 'file';
    input.accept = '.json';
    input.onchange = (e) => {
        const file = (e.target as HTMLInputElement).files?.[0];
        if (!file) {
            return;
        }
        const reader = new FileReader();
        reader.onload = (e) => {
            const data = JSON.parse(e.target?.result as string);
            historyScores.value = data.data;
            localStorage.setItem('history-' + sessionToken.value, JSON.stringify(data));
        }
        reader.readAsText(file);
    }
    input.click();
}
const askForSessionToken = () => {
    var token = localStorage.getItem('selectedSession');
    if (!token) {
        return false;
    } else {
        sessionToken.value = token!;
        return true;
    }
}
const deleteHistory = (index: number) => {
    historyScores.value.splice(index, 1);
    localStorage.setItem('history-' + sessionToken.value, JSON.stringify({ data: historyScores.value }));
}
const checkScore = (score: any) => {
    selectedScore.value = score;
    isSelectHistory.value = false;
}

onMounted(() => {
    if (!askForSessionToken()) {
        alert('Please select a session first');
        window.location.href = '/#/session';
    }
    historyScores.value = JSON.parse(localStorage.getItem('history-' + sessionToken.value) || '[]').data;
})
</script>

<style>
@import 'tailwindcss/base';
@import 'tailwindcss/components';
@import 'tailwindcss/utilities';

@import 'daisyui/dist/full.css';
</style>
