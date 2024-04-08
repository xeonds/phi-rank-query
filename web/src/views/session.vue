<template>
    <div class="p-4">
        <h1 class="text-2xl font-bold mb-4">Sessions</h1>
        <div class="flex items-center mb-4">
            <input type="text" v-model="newSession" placeholder="Enter a session"
                class="mr-2 p-2 border border-gray-800 rounded" />
            <button @click="addSession" class="px-4 py-2 bg-blue-500 text-white rounded">Add Session</button>
        </div>
        <div v-if="sessions.length === 0" class="text-gray-500">No sessions added yet.</div>
        <ul v-else>
            <li v-for="(session, index) in sessions" :key="index" class="flex items-center mb-2">
                <input type="radio" :id="'session-' + index" :value="session" v-model="selectedSession" class="mr-2" />
                <label :for="'session-' + index" class="mr-2">{{ session }}</label>
                <input type="text" v-model="aliases[index]" placeholder="Enter an alias"
                    class="mr-2 p-2 border border-gray-300 rounded" />
                <button @click="deleteSession(index)" class="px-2 py-1 bg-red-500 text-white rounded">Delete</button>
            </li>
        </ul>
    </div>
</template>

<script lang="ts" setup>
import { onMounted, ref, watch } from 'vue';

const sessions = ref<string[]>([]);
const aliases = ref<string[]>([]);
const selectedSession = ref<string | null>(null);
const newSession = ref('');

// Load sessions from localStorage on component mount
const loadSessions = () => {
    const savedSessions = localStorage.getItem('sessions');
    if (savedSessions) {
        sessions.value = JSON.parse(savedSessions);
    }
};

// Save sessions to localStorage whenever it changes
watch(sessions, () => {
    localStorage.setItem('sessions', JSON.stringify(sessions.value));
});

// Load selected session from sessionStorage on component mount
const loadSelectedSession = () => {
    const savedSelectedSession = sessionStorage.getItem('selectedSession');
    if (savedSelectedSession) {
        selectedSession.value = savedSelectedSession;
    }
};

// Save selected session to sessionStorage whenever it changes
watch(selectedSession, () => {
    sessionStorage.setItem('selectedSession', selectedSession.value || '');
});

const addSession = () => {
    if (newSession.value.trim() !== '') {
        sessions.value.push(newSession.value);
        aliases.value.push('');
        newSession.value = '';
    }
};

const deleteSession = (index: number) => {
    sessions.value.splice(index, 1);
    aliases.value.splice(index, 1);
};

// Load sessions and selected session on component mount
onMounted(() => {
    loadSessions();
    loadSelectedSession();
});
</script>
