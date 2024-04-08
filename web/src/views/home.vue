<template>
    <div>
        <header class="py-4">
            <div class="container mx-auto px-4">
                <h1 class="text-gray-200 text-white text-4xl font-bold">欢迎使用Phigros Rks查询工具</h1>
                <p class="text-white mt-2">👽</p>
            </div>
        </header>
        <main class="container mx-auto px-4 py-8">
            <h2 class="text-2xl font-bold mb-4">Features</h2>
            <ul class="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <li class="bg-gray-900 bg-opacity-60 p-4 rounded shadow">
                    <h3 class="text-lg font-bold">定数查询</h3>
                    <p>反正能用（）功能上比较纯粹，后端用Go重构，性能相对更高</p>
                </li>
                <li class="bg-gray-900 bg-opacity-60 p-4 rounded shadow">
                    <h3 class="text-lg font-bold">多存档查询</h3>
                    <p>SessionToken全部保存在浏览器存储中，随时切换要查询的存档</p>
                </li>
            </ul>
            <br>
            <p> 用法：先在Session页面添加存档，然后点击存档左侧来选择要查询的存档，最后切换到B19页面等待一会即可查询到结果 </p>
            <br>
            <p> 为减轻鸽游服务器压力，站点内置1分钟查询cd，期间显示结果为最近一次查询结果 </p>
        </main>
        <footer class="py-4">
            <div class="container mx-auto px-4">
                <p class="text-center text-gray-400"> 2024 &copy; Phi Rank Query</p>
            </div>
        </footer>
    </div>
</template>

