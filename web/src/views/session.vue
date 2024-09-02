<template>
    <div class="p-4">
        <h1 class="text-2xl font-bold mb-4 text-gray-400">Sessions</h1>
        <div class="flex items-center mb-4">
            <input type="text" v-model="newSession" placeholder="输入一个Session..."
                class="input input-ghost input-bordered w-full max-w-xs mr-2 p-2 border border-gray-800 rounded" />
            <input type="text" v-model="newAlias" placeholder="输入别名..."
                class="input input-ghost input-bordered w-full max-w-xs mr-2 p-2 border border-gray-800 rounded" />
            <button @click="addSession" class="px-4 py-2 text-white btn btn-outline btn-primary">
                添加 Session
            </button>
        </div>
        <div v-if="sessions.length === 0" class="text-gray-500">
            咕，还没有添加任何 Session 哦
        </div>
        <table class="table border-collapse text-gray-400" v-else>
            <thead>
                <tr>
                    <th class="py-2 px-4 border-b text-gray-400">Session</th>
                    <th class="py-2 px-4 border-b text-gray-400">别名</th>
                    <th class="py-2 px-4 border-b text-gray-400">操作</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="(session, index) in sessions" :key="index" >
                    <td class="py-2 px-4 border-b">
                        <input type="radio" :id="'session-' + index" :value="session" v-model="selectedSession"
                            class="mr-2 radio" />
                        <label :for="'session-' + index" class="mr-2">
                            {{ session.split('').slice(0, 3).join('') }}***
                            {{ session.split('').slice(-3).join('') }}
                        </label>
                    </td>
                    <td class="py-2 px-4 border-b">
                        <label :for="'alias-' + index" class="mr-2">{{ aliases[index] }}</label>
                    </td>
                    <td class="py-2 px-4 border-b">
                        <button @click="copySession(index)" class="px-2 mr-2 py-1 btn btn-outline btn-info">
                            复 制 Session
                        </button>
                        <button @click="deleteSession(index)" class="px-2 py-1 btn btn-outline btn-warning">
                            删 除 Session
                        </button>
                    </td>
                </tr>
            </tbody>
        </table>
    </div>
</template>

<script lang="ts" setup>
import { onMounted, ref, watch } from 'vue';

const sessions = ref<string[]>([]);
const aliases = ref<string[]>([]);
const selectedSession = ref<string | null>(null);
const newSession = ref('');
const newAlias = ref('');

const loadSessions = () => {
    sessions.value = JSON.parse(localStorage.getItem('sessions') || '[]');
    aliases.value = JSON.parse(localStorage.getItem('aliases') || '[]');
};
const loadSelectedSession = () => {
    const savedSelectedSession = localStorage.getItem('selectedSession');
    if (savedSelectedSession) selectedSession.value = savedSelectedSession;
};
const addSession = () => {
    if (newSession.value.trim() !== '') {
        sessions.value.push(newSession.value);
        aliases.value.push(newAlias.value || newSession.value);
        newSession.value = '';
        newAlias.value = '';
    }
};
const copySession = (index: number) => {
    navigator.clipboard.writeText(sessions.value[index]);
    alert('Session 已复制到剪贴板');
};
const deleteSession = (index: number) => {
    if (!confirm('确定要删除这个 Session 吗？')) return;
    sessions.value.splice(index, 1);
    aliases.value.splice(index, 1);
};

watch(
    sessions,
    () => {
        if (sessions.value.length > 0)
            localStorage.setItem('sessions', JSON.stringify(sessions.value));
        if (aliases.value.length > 0)
            localStorage.setItem('aliases', JSON.stringify(aliases.value));
    },
    { immediate: true, deep: true }
);
watch(
    selectedSession,
    () => {
        if (selectedSession.value)
            localStorage.setItem('selectedSession', selectedSession.value);
    },
    { immediate: true }
);
onMounted(() => {
    loadSessions();
    loadSelectedSession();
});
</script>
