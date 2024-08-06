<template>
    <div>
        <header class="py-4">
            <div class="container mx-auto px-4">
                <h1 class="text-gray-200 text-4xl font-bold pb-8">欢迎使用Phigros Rks查询工具</h1>
                <p class="text-white mt-2 pb-16">👽</p>
            </div>
        </header>
        <main class="container mx-auto px-4">
            <h2 class="text-2xl font-bold mb-4 text-gray-400">Features</h2>
            <ul class="grid grid-cols-1 sm:grid-cols-2 gap-4 pb-16">
                <li class="bg-gray-900 bg-opacity-60 p-4 rounded shadow">
                    <h3 class="text-lg font-bold text-gray-400">多存档查询</h3>
                    <p>SessionToken全部保存在浏览器存储中，随时切换要查询的存档</p>
                </li>
                <li class="bg-gray-900 bg-opacity-60 p-4 rounded shadow">
                    <h3 class="text-lg font-bold text-gray-400">全成绩查询</h3>
                    <p>一次性查询你的所有历史成绩！和B19共享cd哦</p>
                </li>
                <li class="bg-gray-900 bg-opacity-60 p-4 rounded shadow">
                    <h3 class="text-lg font-bold text-gray-400">谱面定数查询</h3>
                    <p>最新（？）最热（？）定数表在线查看</p>
                </li>
                <li class="bg-gray-900 bg-opacity-60 p-4 rounded shadow">
                    <h3 class="text-lg font-bold text-gray-400">排行榜</h3>
                    <p>你们👽在地球住着还习惯吗</p>
                </li>
            </ul>
            <br>
            <p class="text-gray-300"> 用法：先在Session页面添加存档，然后点击存档左侧来选择要查询的存档，最后切换到B19页面等待一会即可查询到结果 </p>
            <br>
            <p class="text-gray-300"> 为减轻鸽游服务器压力，站点内置1分钟查询cd，期间显示结果为最近一次查询结果 </p>
            <br>
            <p class="text-gray-300"> ps: 有任何意见和建议欢迎通过<code>xeonds@stu.xidian.edu.cn</code>给我发邮件哦</p>
            <br>
            <p class="text-gray-300"> pps: 查询历史做好了</p>
            <br>
            <p class="text-gray-300"> ppps: 定数数据库同步3.8.0版本了（</p>
        </main>
        <footer class="py-4">
            <div class="container mx-auto px-4">
                <p class="text-center text-gray-400"> 2024 &copy; Phi Rank Query</p>
            </div>
        </footer>
    </div>
</template>
