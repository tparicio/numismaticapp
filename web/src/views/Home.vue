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
            <div v-for="coin in stats.top_valuable_coins" :key="coin.id" 
                 class="flex items-center gap-3 py-3 cursor-pointer hover:bg-base-200 transition-colors rounded-lg px-2 -mx-2"
                 @click="router.push(`/coin/${coin.id}`)">
              <div class="avatar">
                <div class="mask mask-squircle w-12 h-12 overflow-hidden">
                  <img :src="getThumbnail(coin)" class="hover:scale-110 transition-transform duration-300" />
                </div>
              </div>
              <div class="flex-1 min-w-0">
                <div class="font-bold truncate">{{ coin.name || 'Unknown Coin' }}</div>
                <div class="text-xs opacity-50">
                  <span v-if="coin.year && coin.year !== 0">{{ coin.year }} • </span>
                  {{ coin.country }}
                </div>
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
            <div v-for="coin in stats.recent_coins" :key="coin.id" 
                 class="flex items-center gap-3 py-3 cursor-pointer hover:bg-base-200 transition-colors rounded-lg px-2 -mx-2"
                 @click="router.push(`/coin/${coin.id}`)">
              <div class="avatar">
                <div class="mask mask-squircle w-12 h-12 overflow-hidden">
                  <img :src="getThumbnail(coin)" class="hover:scale-110 transition-transform duration-300" />
                </div>
              </div>
              <div class="flex-1 min-w-0">
                <div class="font-bold truncate">{{ coin.name || 'Unknown Coin' }}</div>
                <div class="text-xs opacity-50">{{ formatDate(coin.created_at) }}</div>
              </div>
              <div class="flex gap-1">
                <button @click.stop="router.push(`/edit/${coin.id}`)" class="btn btn-ghost btn-xs text-info">Edit</button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Advanced Analytics -->
    <div class="card bg-base-100 shadow-xl">
        <div class="card-body">
            <h2 class="card-title">Geographic Distribution</h2>
            <div class="h-64 relative">
                <Bar v-if="countryChartData" :data="countryChartData" :options="chartOptions" />
            </div>
        </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
        <!-- Timeline (Centuries) -->
        <div class="card bg-base-100 shadow-xl">
            <div class="card-body">
                <h2 class="card-title">Timeline (Centuries)</h2>
                <div class="h-64 relative">
                    <Bar v-if="timelineChartData" :data="timelineChartData" :options="chartOptions" />
                </div>
                <div v-if="stats.oldest_coin" class="alert alert-info mt-4">
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="stroke-current shrink-0 w-6 h-6"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
                    <div>
                        <h3 class="font-bold">Oldest Gem</h3>
                        <div class="text-xs">{{ stats.oldest_coin.year }}: {{ stats.oldest_coin.name }}</div>
                    </div>
                </div>
            </div>
        </div>

        <!-- Quality Matrix (Year vs Grade) -->
        <div class="card bg-base-100 shadow-xl">
            <div class="card-body">
                <h2 class="card-title">Quality Matrix (Year vs Grade)</h2>
                <div class="h-64 relative">
                    <Scatter v-if="qualityChartData" :data="qualityChartData" :options="scatterOptions" />
                </div>
            </div>
        </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <!-- Storage Distribution -->
        <div class="card bg-base-100 shadow-xl">
            <div class="card-body">
                <h2 class="card-title">Storage</h2>
                <div class="h-64 relative flex justify-center">
                    <Doughnut v-if="storageChartData" :data="storageChartData" :options="doughnutOptions" />
                </div>
            </div>
        </div>

        <!-- Rarest Coins -->
        <div class="card bg-base-100 shadow-xl">
            <div class="card-body">
                <h2 class="card-title text-sm uppercase text-base-content/70">Top Rarity (Lowest Mintage)</h2>
                <div class="divide-y divide-base-200">
                    <div v-for="coin in stats.rarest_coins" :key="coin.id" 
                         class="flex items-center gap-3 py-3 cursor-pointer hover:bg-base-200 transition-colors rounded-lg px-2 -mx-2"
                         @click="router.push(`/coin/${coin.id}`)">
                         <div class="flex-1">
                            <div class="font-bold">{{ coin.name }}</div>
                            <div class="text-xs opacity-50">{{ coin.year }} • {{ coin.country }}</div>
                         </div>
                         <div class="badge badge-accent">{{ formatMintage(coin.mintage) }}</div>
                    </div>
                </div>
            </div>
        </div>

        <!-- Trivia & Physical Stats -->
        <div class="card bg-base-100 shadow-xl">
            <div class="card-body">
                <h2 class="card-title text-sm uppercase text-base-content/70">Trivia & Physical Stats</h2>
                <div class="space-y-4">
                    <div class="stats stats-vertical shadow w-full">
                        <div class="stat">
                            <div class="stat-title">Silver Reserve</div>
                            <div class="stat-value text-primary text-2xl">{{ stats.total_silver_weight.toFixed(1) }}g</div>
                            <div class="stat-desc">Pure Silver (.999 approx)</div>
                        </div>
                        <div class="stat">
                            <div class="stat-title">Gold Reserve</div>
                            <div class="stat-value text-warning text-2xl">{{ stats.total_gold_weight.toFixed(1) }}g</div>
                            <div class="stat-desc">Pure Gold</div>
                        </div>
                    </div>
                    
                    <div v-if="stats.heaviest_coin" class="text-sm">
                        <span class="font-bold">Heaviest:</span> {{ stats.heaviest_coin.name }} ({{ stats.heaviest_coin.weight_g }}g)
                    </div>
                    <div v-if="stats.smallest_coin" class="text-sm">
                        <span class="font-bold">Smallest:</span> {{ stats.smallest_coin.name }} ({{ stats.smallest_coin.diameter_mm }}mm)
                    </div>
                </div>
            </div>
        </div>
    </div>

    <!-- Random Coin Feature -->
    <div v-if="stats.random_coin" class="card lg:card-side bg-base-100 shadow-xl">
        <figure class="lg:w-1/3 p-6">
            <img :src="getThumbnail(stats.random_coin)" class="rounded-xl shadow-lg hover:scale-105 transition-transform duration-300" />
        </figure>
        <div class="card-body">
            <h2 class="card-title">Featured Coin: {{ stats.random_coin.name }}</h2>
            <p>{{ stats.random_coin.description || 'No description available.' }}</p>
            <div class="card-actions justify-end">
                <button class="btn btn-primary" @click="router.push(`/coin/${stats.random_coin.id}`)">View Details</button>
            </div>
        </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import axios from 'axios'
import { useRouter } from 'vue-router'
import {
  Chart as ChartJS,
  Title,
  Tooltip,
  Legend,
  BarElement,
  CategoryScale,
  LinearScale,
  ArcElement,
  PointElement
} from 'chart.js'
import { Bar, Doughnut, Scatter } from 'vue-chartjs'

ChartJS.register(CategoryScale, LinearScale, BarElement, Title, Tooltip, Legend, ArcElement, PointElement)

const router = useRouter()
const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1'

const stats = ref({
  total_coins: 0,
  total_value: 0,
  top_valuable_coins: [],
  recent_coins: [],
  value_distribution: {},
  material_distribution: {},
  grade_distribution: {},
  country_distribution: {},
  century_distribution: {},
  oldest_coin: null,
  rarest_coins: [],
  group_distribution: {},
  total_silver_weight: 0,
  total_gold_weight: 0,
  heaviest_coin: null,
  smallest_coin: null,
  random_coin: null,
  all_coins: []
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

const formatMintage = (value) => {
    if (value === 0) return 'Unknown'
    if (value >= 1000000) return (value / 1000000).toFixed(1) + 'M'
    if (value >= 1000) return (value / 1000).toFixed(0) + 'k'
    return value.toString()
}

const scatterOptions = {
    responsive: true,
    maintainAspectRatio: false,
    scales: {
        x: {
            type: 'linear',
            position: 'bottom',
            title: { display: true, text: 'Year' }
        },
        y: {
            type: 'linear',
            title: { display: true, text: 'Grade (0-70)' },
            min: 0,
            max: 75
        }
    },
    plugins: {
        legend: { display: false },
        tooltip: {
            callbacks: {
                label: (context) => {
                    const point = context.raw
                    return `${point.name} (${point.grade}): ${point.x}`
                }
            }
        }
    }
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

const timelineChartData = computed(() => {
    const dist = stats.value.century_distribution
    if (!dist) return null

    // Sort centuries
    const labels = Object.keys(dist).sort()
    const data = labels.map(l => dist[l])

    return {
        labels,
        datasets: [{
            label: 'Coins per Century',
            backgroundColor: '#8b5cf6',
            data
        }]
    }
})

const qualityChartData = computed(() => {
    const coins = stats.value.all_coins
    if (!coins || coins.length === 0) return null

    const gradeMap = {
        'MC': 10, 'RC': 20, 'BC': 30, 'MBC': 40, 'EBC': 50, 'SC': 60, 'FDC': 65, 'PROOF': 70
    }

    const data = coins
        .filter(c => c.year > 0 && c.grade && gradeMap[c.grade])
        .map(c => ({
            x: c.year,
            y: gradeMap[c.grade],
            name: c.name,
            grade: c.grade
        }))

    return {
        datasets: [{
            label: 'Quality vs Year',
            backgroundColor: '#f59e0b',
            data
        }]
    }
})

const storageChartData = computed(() => {
    const dist = stats.value.group_distribution
    if (!dist) return null

    const labels = Object.keys(dist)
    const data = Object.values(dist)
    const colors = [
        '#3b82f6', '#10b981', '#f59e0b', '#ef4444', '#8b5cf6', '#ec4899', '#6366f1'
    ]

    return {
        labels,
        datasets: [{
            backgroundColor: colors,
            data
        }]
    }
})

const countryChartData = computed(() => {
    const dist = stats.value.country_distribution
    if (!dist) return null

    // Sort by count desc
    const sorted = Object.entries(dist).sort((a, b) => b[1] - a[1])
    const labels = sorted.map(e => e[0])
    const data = sorted.map(e => e[1])

    return {
        labels,
        datasets: [{
            label: 'Coins by Country',
            backgroundColor: '#0ea5e9',
            data
        }]
    }
})
</script>
