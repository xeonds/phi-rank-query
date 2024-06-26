import { createRouter, createWebHashHistory } from "vue-router";

const routes = [
  {
    path: "/",
    name: "home",
    children: [
      { path: "/", name: "home", component: () => import("@/views/home.vue") },
      { path: "info", name: "info", component: () => import("@/views/info.vue") },
      { path: "b19", name: "b19", component: () => import("@/views/b19.vue") },
      { path: "bn", name: "bn", component: () => import("@/views/bn.vue") },
      { path: "calc", name: "calc", component: () => import("@/views/calc.vue") },
      {
        path: "leaderboard",
        name: "leaderboard",
        component: () => import("@/views/leaderboard.vue"),
      },
      {
        path: "session",
        name: "session",
        component: () => import("@/views/session.vue"),
      },
      {
        path: "history",
        name: "history",
        component: () => import("@/views/history.vue"),
      },
    ],
  },
];

const router = createRouter({
  history: createWebHashHistory(),
  routes,
});

export default router;
