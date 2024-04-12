<template>
    <div class="p-4">
        <h1 class="text-2xl font-bold mb-4 text-gray-400">History</h1>
        <div class="flex items-center mb-4">
            <button @click="exportHistory" class="px-4 py-2 text-white btn btn-outline btn-primary">
                Export History
            </button>
        </div>
        <table class="table w-full border-collapse text-gray-400">
            <thead>
                <tr>
                    <th class="py-2 px-4 border-b text-gray-400">Score</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="score in historyScores" :key="score">
                    <td class="py-2 px-4 border-b text-gray-400">{{ score }}</td>
                </tr>
            </tbody>
        </table>
    </div>
</template>

<script lang="ts" setup>
import { Ref, onMounted, ref } from 'vue';

const historyScores: Ref<any[]> = ref([])
const sessionToken = ref<string>('');

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
const askForSessionToken = () => {
  var token = sessionStorage.getItem('selectedSession');
  if (!token) {
    return false;
  } else {
    sessionToken.value = token!;
    return true;
  }
}

onMounted(() => {
    if (!askForSessionToken()) {
      alert('Please select a session first');
      window.location.href = '/#/session';
    }
    historyScores.value = JSON.parse(localStorage.getItem('history-'+sessionToken.value)||'[]').data;
})
</script>

<style>
@import 'tailwindcss/base';
@import 'tailwindcss/components';
@import 'tailwindcss/utilities';

@import 'daisyui/dist/full.css';
</style>
