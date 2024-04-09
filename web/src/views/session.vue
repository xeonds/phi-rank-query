<template>
    <div class="p-4">
        <h1 class="text-2xl font-bold mb-4">Sessions</h1>
        <div class="flex items-center mb-4">
            <input type="text" v-model="newSession" placeholder="Enter a session"
                class="input input-ghost input-bordered w-full max-w-xs mr-2 p-2 border border-gray-800 rounded" />
            <input type="text" v-model="newAlias" placeholder="Enter alias"
                class="input input-ghost input-bordered w-full max-w-xs mr-2 p-2 border border-gray-800 rounded" />
            <button @click="addSession" class="px-4 py-2 text-white btn btn-outline btn-primary">
                Add Session
            </button>
        </div>
        <div v-if="sessions.length === 0" class="text-gray-500">
            No sessions added yet.
        </div>
        <table class="table border-collapse" v-else>
            <thead>
                <tr>
                    <th>Session</th>
                    <th>Alias</th>
                    <th>Action</th>
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
                        <button @click="deleteSession(index)" class="px-2 py-1 btn btn-outline btn-warning">
                            Delete
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
    const savedSelectedSession = sessionStorage.getItem('selectedSession');
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
const deleteSession = (index: number) => {
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
            sessionStorage.setItem('selectedSession', selectedSession.value);
    },
    { immediate: true }
);
onMounted(() => {
    loadSessions();
    loadSelectedSession();
});
</script>
