import { createRouter, createWebHashHistory } from 'vue-router'

const routes = [
  { path: '/', name: 'home', component: () => import('@/views/HomeView.vue') },
  { path: '/info', name: 'rank', component: () => import('@/views/RankView.vue') },
  { path: '/b19', name: 'b19', component: () => import('@/views/B19View.vue') },
  { path: '/bn', name: 'bn', component: () => import('@/views/BNView.vue') },
  { path: '/calc', name: 'calc', component: () => import('@/views/CalcView.vue') },
  { path: '/history', name: 'history', component: () => import('@/views/HistoryView.vue') },
  { path: '/leaderboard', name: 'leaderboard', component: () => import('@/views/LeaderboardView.vue') },
  { path: '/session', name: 'session', component: () => import('@/views/SessionView.vue') },
]

export default createRouter({
  history: createWebHashHistory(),
  routes,
})