import { createRouter, createWebHistory } from 'vue-router'
import Home from '../views/Home.vue'
import AddCoin from '../views/AddCoin.vue'
import List from '../views/List.vue'
import Detail from '../views/Detail.vue'

const routes = [
    { path: '/', component: Home },
    { path: '/add', component: AddCoin },
    { path: '/list', component: List },
    { path: '/coin/:id', component: Detail },
]

const router = createRouter({
    history: createWebHistory(),
    routes,
})

export default router
