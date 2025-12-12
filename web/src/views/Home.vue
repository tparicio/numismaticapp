<template>
  <div class="container mx-auto p-4 space-y-8">
    <!-- Header -->
    <div class="flex justify-between items-center">
      <div>
        <h1 class="text-3xl font-bold">{{ $t('dashboard.title') }}</h1>
        <p class="text-base-content/70">{{ $t('dashboard.subtitle') }}</p>
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
          <div class="stat-title">{{ $t('dashboard.stats.total_value') }}</div>
          <div class="stat-value text-primary">{{ formatCurrency(stats.total_value) }}</div>
          <div class="stat-desc">{{ $t('dashboard.stats.market_value') }}</div>
        </div>
      </div>

      <div class="stats shadow">
        <div class="stat">
          <div class="stat-figure text-accent">
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="inline-block w-8 h-8 stroke-current"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 3v4M3 5h4M6 17v4m-2-2h4m5-16l2.286 6.857L21 12l-5.714 2.143L13 21l-2.286-6.857L5 12l5.714-2.143L13 3z" /></svg>
          </div>
          <div class="stat-title">Top Rarity</div>
          <div class="stat-value text-accent text-lg overflow-hidden truncate whitespace-nowrap max-w-[10rem]" :title="stats.rarest_coins?.[0]?.name">
            {{ stats.rarest_coins?.[0]?.name || 'N/A' }}
          </div>
          <div class="stat-desc" v-if="stats.rarest_coins?.[0]">
            {{ formatMintage(stats.rarest_coins[0].mintage) }} units
          </div>
        </div>
      </div>
      
      <div class="stats shadow">
        <div class="stat">
          <div class="stat-figure text-secondary">
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="inline-block w-8 h-8 stroke-current"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"></path></svg>
          </div>
          <div class="stat-title">{{ $t('dashboard.stats.total_coins') }}</div>
          <div class="stat-value text-secondary">{{ stats.total_coins }}</div>
          <div class="stat-desc">{{ $t('dashboard.stats.in_collection') }}</div>
        </div>
      </div>
    </div>

    <!-- Charts Row 1 -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
      <!-- Value Distribution -->
      <div class="card bg-base-100 shadow-xl">
        <div class="card-body">
          <h2 class="card-title">{{ $t('dashboard.charts.value_dist') }}</h2>
          <div class="h-64 relative">
            <Bar v-if="valueChartData" :data="valueChartData" :options="chartOptions" />
          </div>
        </div>
      </div>

      <!-- Material Distribution -->
      <div class="card bg-base-100 shadow-xl">
        <div class="card-body">
          <h2 class="card-title">{{ $t('dashboard.charts.materials') }}</h2>
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
          <h2 class="card-title">{{ $t('dashboard.charts.grades') }}</h2>
          <div class="h-64 relative">
            <Bar v-if="gradeChartData" :data="gradeChartData" :options="chartOptions" />
          </div>
        </div>
      </div>



    <!-- Main Content: Map & Lists -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <!-- Map (2/3 width) -->
        <div class="card bg-base-100 shadow-xl lg:col-span-2">
            <div class="card-body">
                <h2 class="card-title">{{ $t('dashboard.charts.geo') }}</h2>
                <div class="h-96 relative w-full rounded-box overflow-hidden border border-base-200">
                    <MapChart v-if="stats.country_distribution" :data="stats.country_distribution" />
                </div>
            </div>
        </div>

        <!-- Labs / Lists (1/3 width) -->
        <div class="card bg-base-100 shadow-xl">
            <div class="card-body p-4">
                <div class="tabs tabs-boxed bg-base-200 mb-4">
                    <a class="tab" :class="{ 'tab-active': activeTab === 'recent' }" @click="activeTab = 'recent'">{{ $t('dashboard.tabs.recent') }}</a>
                    <a class="tab" :class="{ 'tab-active': activeTab === 'valuable' }" @click="activeTab = 'valuable'">{{ $t('dashboard.tabs.valuable') }}</a>
                </div>

                <div class="divide-y divide-base-200" v-if="activeTab === 'recent'">
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
                        <button @click.stop="router.push(`/edit/${coin.id}`)" class="btn btn-ghost btn-xs text-info">{{ $t('common.edit') }}</button>
                      </div>
                    </div>
                </div>

                <div class="divide-y divide-base-200" v-if="activeTab === 'valuable'">
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
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
        <!-- Timeline (Centuries) -->
        <div class="card bg-base-100 shadow-xl">
            <div class="card-body">
                <h2 class="card-title">{{ $t('dashboard.charts.timeline') }}</h2>
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
    <div v-if="stats.random_coin" class="card lg:card-side bg-base-100 shadow-xl overflow-hidden">
        <figure class="lg:w-auto lg:shrink-0 bg-base-200 p-6 flex flex-col gap-4 justify-center items-center">
            <div class="flex gap-4">
                <div class="avatar cursor-pointer" @click="router.push(`/coin/${stats.random_coin.id}`)">
                    <div class="w-40 sm:w-48 rounded-xl shadow-xl hover:scale-105 transition-transform duration-300">
                        <img :src="getThumbnail(stats.random_coin, 'front')" class="object-contain" />
                    </div>
                </div>
                <div class="avatar cursor-pointer" @click="router.push(`/coin/${stats.random_coin.id}`)">
                    <div class="w-40 sm:w-48 rounded-xl shadow-xl hover:scale-105 transition-transform duration-300">
                        <img :src="getThumbnail(stats.random_coin, 'back')" class="object-contain" />
                    </div>
                </div>
            </div>
        </figure>
        <div class="card-body">
            <div class="flex justify-between items-start">
                <div>
                    <h2 class="card-title text-3xl font-bold mb-1">{{ stats.random_coin.name }}</h2>
                    <div class="text-lg opacity-70 flex items-center gap-2">
                        <span class="font-semibold">{{ stats.random_coin.country }}</span>
                        <span>•</span>
                        <span>{{ stats.random_coin.year }}</span>
                    </div>
                </div>
                <div class="badge badge-lg badge-primary font-bold">{{ formatCurrency(stats.random_coin.max_value) }}</div>
            </div>

            <div class="divider my-2"></div>

            <div class="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
                <div>
                    <div class="opacity-50">Face Value</div>
                    <div class="font-bold">{{ stats.random_coin.face_value }} {{ stats.random_coin.currency }}</div>
                </div>
                <div>
                    <div class="opacity-50">Mint</div>
                    <div class="font-bold">{{ stats.random_coin.mint || '-' }}</div>
                </div>
                <div>
                    <div class="opacity-50">Mintage</div>
                    <div class="font-bold">{{ formatMintage(stats.random_coin.mintage) }}</div>
                </div>
                <div>
                    <div class="opacity-50">KM Code</div>
                    <div class="font-bold">{{ stats.random_coin.km_code || '-' }}</div>
                </div>
                <div v-if="stats.random_coin.group_name">
                    <div class="opacity-50">Collection</div>
                    <div class="badge badge-outline mt-1">{{ stats.random_coin.group_name }}</div>
                </div>
            </div>

            <p class="mt-4 text-base-content/80 italic border-l-4 border-primary pl-4 py-2 bg-base-200/50 rounded-r">
                "{{ stats.random_coin.description || 'No description available.' }}"
            </p>

            <div class="card-actions justify-end mt-4">
                <button class="btn btn-primary gap-2" @click="router.push(`/coin/${stats.random_coin.id}`)">
                    View Full Details
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-4 h-4">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M13.5 4.5L21 12m0 0l-7.5 7.5M21 12H3" />
                    </svg>
                </button>
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
import MapChart from '../components/MapChart.vue'

ChartJS.register(CategoryScale, LinearScale, BarElement, Title, Tooltip, Legend, ArcElement, PointElement)

const router = useRouter()

const API_URL = import.meta.env.VITE_API_URL || '/api/v1'
const activeTab = ref('recent')

const stats = ref({
  total_coins: 0,
  total_value: 0,
  average_value: 0,
  top_valuable_coins: [],
  recent_coins: [],
  value_distribution: {},
  material_distribution: {},
  grade_distribution: {},
  country_distribution: {},
  decade_distribution: {},
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

const getThumbnail = (coin, side = 'front') => {
  if (coin.images && coin.images.length > 0) {
    const thumb = coin.images.find(i => i.image_type === 'thumbnail' && i.side === side)
    if (thumb) return `${API_URL.replace('/api/v1', '')}/${thumb.path}`
    const crop = coin.images.find(i => i.image_type === 'crop' && i.side === side)
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
      backgroundColor: '#7c3aed', // Violeta (Money)
      data
    }]
  }
})

const materialChartData = computed(() => {
  const dist = stats.value.material_distribution
  if (!dist) return null
  
  const labels = Object.keys(dist)
  const data = Object.values(dist)
  const colors = labels.map(label => {
      const l = label.toLowerCase()
      if (l.includes('gold') || l.includes('oro')) return '#D97706'
      if (l.includes('silver') || l.includes('plata')) return '#94A3B8'
      if (l.includes('copper') || l.includes('bronze') || l.includes('cobre') || l.includes('laton')) return '#C2410C'
      return '#CBD5E1' // Nickel/Steel/Other
  })

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
  
  // Custom sort order for grades if possible, otherwise rely on backend or simple sort
  const gradeOrder = ['RC', 'BC', 'MBC', 'EBC', 'SC', 'FDC', 'PROOF']
  const labels = Object.keys(dist).sort((a, b) => {
      return gradeOrder.indexOf(a) - gradeOrder.indexOf(b)
  })
  
  const data = labels.map(l => dist[l])
  
  const colors = labels.map(label => {
      if (['RC', 'BC'].includes(label)) return '#BEF264' // Low
      if (['MBC', 'EBC'].includes(label)) return '#22C55E' // Mid
      if (['SC', 'FDC', 'PROOF'].includes(label)) return '#15803D' // High
      return '#22C55E' // Default
  })

  return {
    labels,
    datasets: [{
      label: 'Coins',
      backgroundColor: colors,
      data
    }]
  }
})

const timelineChartData = computed(() => {
    const dist = stats.value.decade_distribution
    if (!dist) return null

    // Sort decades
    const labels = Object.keys(dist).sort()
    const data = labels.map(l => dist[l])

    return {
        labels,
        datasets: [{
            label: 'Coins per Decade',
            backgroundColor: '#4F46E5', // Indigo
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
