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
                    <Song v-if="index % 2 == 0" :index="`#${index}`" :song="rankTable[index]" />
                </template>
            </div>
            <div class="flex flex-col px-8">
                <template v-for="(_, index) in rankTable">
                    <Song v-if="index % 2 == 1" :index="`#${index}`" :song="rankTable[index]" />
                </template>
            </div>
        </div>
        <div v-else>
            <form>
                <label for="当前项目"></label>
                <input :model="curr" />
                <label for="输入项1"></label>
                <input :model="curr" />
                <label for="输入项2"></label>
                <input :model="curr" />
            </form>
        </div>
    </div>
</template>

<script lang="ts" setup>
import Song from '@/components/song.vue';
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