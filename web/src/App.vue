<template>
    <div class="relative min-h-screen">
        <div class="absolute inset-0 bg-cover bg-center bg-fixed z-0 blur-md"
            style="background-image: url('/assets/Star1.png');"></div>
        <div class="absolute inset-0 bg-cover bg-center bg-fixed z-0 bg-gray-900 bg-opacity-50"></div>
        <div class="container mx-auto relative z-10">
            <nav class="flex justify-between items-center py-4 relative">
                <div>
                    <h1 class="text-gray-200 text-3xl font-bold">Phi Rank Query</h1>
                </div>
                <div class="space-x-4">
                    <a href="/#/" class="link link-hover text-gray-300">主页</a>
                    <a href="/#/info" class="link link-hover text-gray-300">定数查询</a>
                    <a href="/#/b19" class="link link-hover text-gray-300">B19查询</a>
                    <a href="/#/bn" class="link link-hover text-gray-300">所有成绩查询</a>
                    <a href="/#/calc" class="link link-hover text-gray-300">Rks手动计算</a>
                    <a href="/#/history" class="link link-hover text-gray-300">查询历史</a>
                    <a href="/#/leaderboard" class="link link-hover text-gray-300">排行榜</a>
                    <a href="/#/session" class="link link-hover text-gray-300">Sessions</a>
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