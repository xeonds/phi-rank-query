<template>
    <div class="container mx-auto">
        <nav class="flex justify-between items-center py-4">
            <div>
                <h1 class="text-3xl font-bold">Phi Rank Query</h1>
            </div>
            <div class="space-x-4">
                <a href="/#/" class="text-lg font-medium text-gray-700 hover:text-gray-100">Home</a>
                <a href="/#/b19" class="text-lg font-medium text-gray-700 hover:text-gray-100">Queries</a>
                <a href="/#/session" class="text-lg font-medium text-gray-700 hover:text-gray-100">Sessions</a>
                <a href="/#/about" class="text-lg font-medium text-gray-700 hover:text-gray-100">About</a>
            </div>
        </nav>
        <router-view></router-view>
    </div>
</template>

<style>
body{
  margin: 0;
}
</style>