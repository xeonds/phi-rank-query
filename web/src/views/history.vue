<template>
    <div class="p-4">
        <h1 class="text-2xl font-bold mb-4">History</h1>
        <table class="table w-full border-collapse">
            <thead>
                <tr>
                    <th class="py-2 px-4 border-b">Score</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="score in historyScores" :key="score">
                    <td class="py-2 px-4 border-b">{{ score }}</td>
                </tr>
            </tbody>
        </table>
    </div>
</template>

<script lang="ts" setup>
import { onMounted, ref } from 'vue';

const historyScores = ref([])
onMounted(() => {
    historyScores.value = JSON.parse(localStorage.getItem('history')||'[]');
    console.log(historyScores.value)
})
</script>

<style>
@import 'tailwindcss/base';
@import 'tailwindcss/components';
@import 'tailwindcss/utilities';

@import 'daisyui/dist/full.css';
</style>
