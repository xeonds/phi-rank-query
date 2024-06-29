<template>
    <div>
        <h1 class="text-2xl font-bold mb-4 text-gray-400">排行榜</h1>
        <table class="table table-bordered">
            <thead>
                <tr>
                    <th class="border-b px-4 py-4 text-gray-400">排行</th>
                    <th class="border-b px-4 py-4 text-gray-400">User ID</th>
                    <th class="border-b px-4 py-4 text-gray-400">Rks</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="(player, index) in leaderboardData" :key="player.id">
                    <td class="border-b px-4 py-4 text-gray-400">{{ index + 1 }}</td>
                    <td class="border-b px-4 py-4 text-gray-400">{{ player.username }}</td>
                    <td class="border-b px-4 py-4 text-gray-400">{{ player.rks.toFixed(4) }}</td>
                </tr>
            </tbody>
        </table>
    </div>
</template>

<script lang="ts" setup>
import { ref, onMounted, Ref } from 'vue';
const leaderboardData: Ref<{ id: number, username: string, rks: number }[]> = ref([]);

const fetchData = async () => {
    const data: Ref<any> = ref(null);
    const err: Ref<any> = ref(null);
    await fetch('/api/v1/leaderboard', {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json'
        }
    })
        .then(response => response.json())
        .then(json => data.value = json)
        .catch(error => err.value = error);

    return { data, err };
};

onMounted(async () => {
    const { data, err } = await fetchData();
    if (err.value) { alert("查询榜单失败" + err.value); }
    else { leaderboardData.value = data.value; }
});
</script>
