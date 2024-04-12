<template>
    <div class="container mx-auto">
        <h1 class="text-2xl font-bold mb-4">Rank Query Page</h1>
        <div v-if="isLoading" class="text-center">
            <p>Loading...</p>
        </div>
        <div v-else>
            <table class="table table-auto w-full">
                <thead>
                    <tr>
                        <th class="px-4 py-2">ID</th>
                        <th class="px-4 py-2">EZ</th>
                        <th class="px-4 py-2">HD</th>
                        <th class="px-4 py-2">IN</th>
                        <th class="px-4 py-2">AT</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="(item, index) in rankTable" :key="index">
                        <td class="border px-4 py-2">{{ item.title }}</td>
                        <td class="border px-4 py-2">{{ item.EZ }}</td>
                        <td class="border px-4 py-2">{{ item.HD }}</td>
                        <td class="border px-4 py-2">{{ item.IN }}</td>
                        <td class="border px-4 py-2">{{ item.AT || '-' }}</td>
                    </tr>
                </tbody>
            </table>
        </div>
    </div>
</template>

<script lang="ts" setup>
import { ref, onMounted, Ref } from 'vue';

const rankTable: Ref<{ title: string, EZ: string, HD: string, IN: string, AT: string }[]> = ref([]);
const isLoading = ref(true);

onMounted(async () => {
    try {
        const response = await fetch('/api/v1/rank_table');
        const data = await response.json();
        rankTable.value = data;
        isLoading.value = false;
    } catch (error) {
        console.error(error);
        isLoading.value = false;
    }
});
</script>

<style>
/* Add any custom styles here */
</style>