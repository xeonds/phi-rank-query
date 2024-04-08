<template>
    <div class="relative min-h-screen">
        <div class="absolute inset-0 bg-cover bg-center bg-fixed z-0 blur-md" style="background-image: url('/assets/Star1.png');"></div>
        <div class="absolute inset-0 bg-cover bg-center bg-fixed z-0 bg-gray-900 bg-opacity-50"></div>
        <div class="container mx-auto relative z-10">
            <nav class="flex justify-between items-center py-4 relative">
                <div>
                    <h1 class="text-gray-400 text-3xl font-bold text-white">Phi Rank Query</h1>
                </div>
                <div class="space-x-4">
                    <a href="/#/" class="link link-hover">主页</a>
                    <a href="/#/b19" class="link link-hover">B19查询</a>
                    <a href="/#/history" class="link link-hover">查询历史</a>
                    <a href="/#/session" class="link link-hover">Sessions</a>
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