<template>
    <div class="container mx-auto py-4">
        <h1 class="text-2xl font-bold mb-4 text-gray-400">Rank Query Page</h1>
        <div v-if="isLoading" class="text-center">
            <p>咕咕咕...正在加载</p>
        </div>
        <div v-else>
            <table class="table border-collapse table-auto w-full">
                <thead>
                    <tr>
                        <th class="px-4 border-b py-2 text-gray-400">名称</th>
                        <th class="px-4 border-b py-2 text-gray-400">难度</th>
                        <th class="px-4 border-b py-2 text-gray-400">定数</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="(item, index) in rankTable" :key="index">
                        <td class="border-b px-4 py-2 text-gray-400">{{ item.title }}</td>
                        <td class="border-b px-4 py-2 text-gray-400">{{ item.difficulty }}</td>
                        <td class="border-b px-4 py-2 text-gray-400">{{ item.rank }}</td>
                    </tr>
                </tbody>
            </table>
        </div>
    </div>
</template>

<script lang="ts" setup>
import { ref, onMounted, Ref } from 'vue';

interface RankTableItem {
    title: string;
    difficulty: string;
    rank: number;
}

const rankTable: Ref<RankTableItem[]> = ref([]);
const isLoading = ref(true);

onMounted(async () => {
    const response = await fetch('/api/v1/rank_table');
    const data = await response.json();
    rankTable.value = Object.values(data)
        .map((item: any) => ([
            { title: item.title, difficulty: 'EZ', rank: parseFloat(item.EZ) },
            { title: item.title, difficulty: 'HD', rank: parseFloat(item.HD) },
            { title: item.title, difficulty: 'IN', rank: parseFloat(item.IN) },
            { title: item.title, difficulty: 'AT', rank: parseFloat(item.AT) || NaN }
        ]))
        .flat()
        .filter((item: any) => !Number.isNaN(item.rank))
        .concat()
        .sort((a: any, b: any) => b.rank - a.rank);
    isLoading.value = false;
});
</script>

<style>
/* Add any custom styles here */
</style>