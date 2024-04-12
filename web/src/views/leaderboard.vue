<template>
    <div>
        <h1 class="text-2xl font-bold mb-4 text-gray-400">Leaderboard</h1>
        <table class="table table-bordered">
            <thead>
                <tr>
                    <th>Ranking</th>
                    <th>Username</th>
                    <th>Rks</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="(player, index) in leaderboardData" :key="player.id">
                    <td>{{ index + 1 }}</td>
                    <td>{{ player.username }}</td>
                    <td>{{ player.rks.toFixed(4) }}</td>
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
