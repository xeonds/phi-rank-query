<template>
    <div class="relative min-h-screen">
        <div class="absolute inset-0 bg-cover bg-center bg-fixed z-0 blur-md" style="background-image: url('/assets/Star1.png');"></div>
        <div class="absolute inset-0 bg-cover bg-center bg-fixed z-0 bg-gray-900 bg-opacity-50"></div>
        <div class="container mx-auto relative z-10">
            <nav class="flex justify-between items-center py-4 relative">
                <div>
                    <h1 class="text-3xl font-bold text-white">Phi Rank Query</h1>
                </div>
                <div class="space-x-4">
                    <a href="/#/" class="text-lg font-medium text-gray-300 hover:text-gray-100">Home</a>
                    <a href="/#/b19" class="text-lg font-medium text-gray-300 hover:text-gray-100">Queries</a>
                    <a href="/#/session" class="text-lg font-medium text-gray-300 hover:text-gray-100">Sessions</a>
                </div>
            </nav>
            <router-view></router-view>
        </div>
    </div>
</template>

<style>
body {
    margin: 0;

    p {
        color: white;
        margin: 0px;
        margin-block: 0;
        margin-inline: 0;
    }
}
</style>