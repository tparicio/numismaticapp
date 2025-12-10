<template>
  <div class="container mx-auto p-4 space-y-8">
    <!-- Header -->
    <div class="flex justify-between items-center">
      <div>
        <h1 class="text-3xl font-bold">Dashboard</h1>
        <p class="text-base-content/70">Overview of your numismatic collection</p>
      </div>
      <router-link to="/add" class="btn btn-primary gap-2">
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6">
          <path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15" />
        </svg>
        Add Coin
      </router-link>
    </div>

    <!-- Stats Cards -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
      <div class="stats shadow">
        <div class="stat">
          <div class="stat-figure text-primary">
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="inline-block w-8 h-8 stroke-current"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
          </div>
          <div class="stat-title">Total Value</div>
          <div class="stat-value text-primary">{{ formatCurrency(stats.total_value) }}</div>
          <div class="stat-desc">Estimated market value</div>
        </div>
      </div>
      
      <div class="stats shadow">
        <div class="stat">
          <div class="stat-figure text-secondary">
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="inline-block w-8 h-8 stroke-current"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"></path></svg>
          </div>
          <div class="stat-title">Total Coins</div>
          <div class="stat-value text-secondary">{{ stats.total_coins }}</div>
          <div class="stat-desc">In collection</div>
        </div>
      </div>
    </div>

    <!-- Charts Row 1 -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
      <!-- Value Distribution -->
      <div class="card bg-base-100 shadow-xl">
        <div class="card-body">
          <h2 class="card-title">Value Distribution</h2>
          <div class="h-64 relative">
            <Bar v-if="valueChartData" :data="valueChartData" :options="chartOptions" />
          </div>
        </div>
      </div>

      <!-- Material Distribution -->
      <div class="card bg-base-100 shadow-xl">
        <div class="card-body">
          <h2 class="card-title">Materials</h2>
          <div class="h-64 relative flex justify-center">
            <Doughnut v-if="materialChartData" :data="materialChartData" :options="doughnutOptions" />
          </div>
        </div>
      </div>
    </div>

    <!-- Charts Row 2 & Lists -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
      <!-- Grade Distribution -->
      <div class="card bg-base-100 shadow-xl lg:col-span-1">
        <div class="card-body">
          <h2 class="card-title">Grades</h2>
          <div class="h-64 relative">
            <Bar v-if="gradeChartData" :data="gradeChartData" :options="chartOptions" />
          </div>
        </div>
      </div>

      <!-- Top Valuable -->
      <div class="card bg-base-100 shadow-xl lg:col-span-1">
        <div class="card-body">
          <h2 class="card-title text-sm uppercase text-base-content/70">Most Valuable</h2>
          <div class="divide-y divide-base-200">
            <div v-for="coin in stats.top_valuable_coins" :key="coin.id" class="flex items-center gap-3 py-3">
              <div class="avatar">
                <div class="mask mask-squircle w-12 h-12">
                  <img :src="getThumbnail(coin)" />
                </div>
              </div>
              <div class="flex-1 min-w-0">
                <div class="font-bold truncate">{{ coin.name || 'Unknown Coin' }}</div>
                <div class="text-xs opacity-50">{{ coin.year }} • {{ coin.country }}</div>
              </div>
              <div class="font-bold text-primary">{{ formatCurrency(coin.max_value) }}</div>
            </div>
          </div>
        </div>
      </div>

      <!-- Recent -->
      <div class="card bg-base-100 shadow-xl lg:col-span-1">
        <div class="card-body">
          <h2 class="card-title text-sm uppercase text-base-content/70">Recently Added</h2>
          <div class="divide-y divide-base-200">
            <div v-for="coin in stats.recent_coins" :key="coin.id" class="flex items-center gap-3 py-3">
              <div class="avatar">
                <div class="mask mask-squircle w-12 h-12">
                  <img :src="getThumbnail(coin)" />
                </div>
              </div>
              <div class="flex-1 min-w-0">
                <div class="font-bold truncate">{{ coin.name || 'Unknown Coin' }}</div>
                <div class="text-xs opacity-50">{{ formatDate(coin.created_at) }}</div>
              </div>
              <router-link :to="`/coin/${coin.id}`" class="btn btn-ghost btn-xs">View</router-link>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import axios from 'axios'
import {
  Chart as ChartJS,
  Title,
  Tooltip,
  Legend,
  BarElement,
  CategoryScale,
  LinearScale,
  ArcElement
} from 'chart.js'
import { Bar, Doughnut } from 'vue-chartjs'

ChartJS.register(CategoryScale, LinearScale, BarElement, Title, Tooltip, Legend, ArcElement)

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1'

const stats = ref({
  total_coins: 0,
  total_value: 0,
  top_valuable_coins: [],
  recent_coins: [],
  value_distribution: {},
  material_distribution: {},
  grade_distribution: {}
})

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: { display: false }
  }
}

const doughnutOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: { position: 'right' }
  }
}

onMounted(async () => {
  try {
    const res = await axios.get(`${API_URL}/dashboard`)
    stats.value = res.data
  } catch (e) {
    console.error("Failed to load dashboard stats", e)
  }
})

const formatCurrency = (val) => {
  return new Intl.NumberFormat('es-ES', { style: 'currency', currency: 'EUR' }).format(val || 0)
}

const formatDate = (dateStr) => {
  return new Date(dateStr).toLocaleDateString()
}

const getThumbnail = (coin) => {
  if (coin.images && coin.images.length > 0) {
    const thumb = coin.images.find(i => i.image_type === 'thumbnail' && i.side === 'front')
    if (thumb) return `${API_URL.replace('/api/v1', '')}/${thumb.path}`
    const crop = coin.images.find(i => i.image_type === 'crop' && i.side === 'front')
    if (crop) return `${API_URL.replace('/api/v1', '')}/${crop.path}`
  }
  return 'https://placehold.co/100x100?text=No+Image'
}

// Chart Data Computed Properties

const valueChartData = computed(() => {
  const dist = stats.value.value_distribution
  if (!dist) return null
  
  // Ensure order
  const labels = ["0-10", "10-50", "50-100", "100-500", "500+"]
  const data = labels.map(l => dist[l] || 0)

  return {
    labels,
    datasets: [{
      label: 'Coins',
      backgroundColor: '#3b82f6',
      data
    }]
  }
})

const materialChartData = computed(() => {
  const dist = stats.value.material_distribution
  if (!dist) return null
  
  const labels = Object.keys(dist)
  const data = Object.values(dist)
  const colors = [
    '#FF6384', '#36A2EB', '#FFCE56', '#4BC0C0', '#9966FF', '#FF9F40',
    '#C9CBCF', '#FF6384', '#36A2EB', '#FFCE56'
  ]

  return {
    labels,
    datasets: [{
      backgroundColor: colors,
      data
    }]
  }
})

const gradeChartData = computed(() => {
  const dist = stats.value.grade_distribution
  if (!dist) return null
  
  const labels = Object.keys(dist)
  const data = Object.values(dist)

  return {
    labels,
    datasets: [{
      label: 'Coins',
      backgroundColor: '#10b981',
      data
    }]
  }
})
</script>
