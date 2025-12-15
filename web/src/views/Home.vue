<template>
  <div v-if="loading" class="flex justify-center items-center min-h-screen">
    <span class="loading loading-spinner loading-lg text-primary"></span>
  </div>
  
  <div v-else class="space-y-8">
    <!-- Header -->
    <!-- Header -->
    <div class="flex flex-col md:flex-row justify-between items-center gap-4">
      <div class="text-center md:text-left">
        <h1 class="text-2xl md:text-3xl font-bold">{{ $t('dashboard.title') }}</h1>
        <p class="text-base-content/70">{{ $t('dashboard.subtitle') }}</p>
      </div>
      <router-link to="/add" class="btn btn-primary gap-2 w-full md:w-auto">
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6">
          <path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15" />
        </svg>
        {{ $t('nav.add_coin') }}
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
          <div class="stat-figure text-secondary">
             <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-8 h-8">
              <path stroke-linecap="round" stroke-linejoin="round" d="M3 13.125C3 12.504 3.504 12 4.125 12h2.25c.621 0 1.125.504 1.125 1.125v6.75C7.5 20.496 6.996 21 6.375 21h-2.25A1.125 1.125 0 013 19.875v-6.75zM9.75 8.625c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125v11.25c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 01-1.125-1.125V8.625zM16.5 4.125c0-.621.504-1.125 1.125-1.125h2.25C20.496 3 21 3.504 21 4.125v15.75c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 01-1.125-1.125V4.125z" />
            </svg>
          </div>
          <div class="stat-title">{{ $t('dashboard.stats.avg_value') }}</div>
          <div class="stat-value text-secondary">{{ formatCurrency(stats.average_value) }}</div>
          <div class="stat-desc">{{ $t('dashboard.stats.per_coin') }}</div>
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

      <div class="stats shadow">
        <div class="stat">
          <div class="stat-figure text-accent">
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="inline-block w-8 h-8 stroke-current"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 3v4M3 5h4M6 17v4m-2-2h4m5-16l2.286 6.857L21 12l-5.714 2.143L13 21l-2.286-6.857L5 12l5.714-2.143L13 3z" /></svg>
          </div>
          <div class="stat-title">{{ $t('dashboard.stats.top_rarity') }}</div>
          <div class="stat-value text-accent text-lg overflow-hidden truncate whitespace-nowrap max-w-[10rem]" :title="stats.rarest_coins?.[0]?.name">
            {{ stats.rarest_coins?.[0]?.name || 'N/A' }}
          </div>
          <div class="stat-desc" v-if="stats.rarest_coins?.[0]">
            {{ formatMintage(stats.rarest_coins[0].mintage) }} {{ $t('dashboard.stats.units') }}
          </div>
        </div>
      </div>
    </div>

    <!-- Charts Row 1: Distributions (Value, Grade, Material) -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
      <!-- Value Distribution -->
      <div class="card bg-base-100 shadow-xl">
        <div class="card-body">
          <h2 class="card-title">{{ $t('dashboard.charts.value_dist') }}</h2>
          <div class="h-64 relative">
            <Bar v-if="valueChartData" :data="valueChartData" :options="barOptions" />
          </div>
        </div>
      </div>

      <!-- Grade Distribution -->
      <div class="card bg-base-100 shadow-xl">
        <div class="card-body">
          <h2 class="card-title">{{ $t('dashboard.charts.grades') }}</h2>
          <div class="h-64 relative">
            <Bar v-if="gradeChartData" :data="gradeChartData" :options="barOptions" />
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
                    <a class="tab" :class="{ 'tab-active': activeTab === 'recent' }" @click="activeTab = 'recent'">{{ $t('dashboard.tabs.recent') || 'Recientes' }}</a>
                    <a class="tab" :class="{ 'tab-active': activeTab === 'valuable' }" @click="activeTab = 'valuable'">{{ $t('dashboard.tabs.valuable') || 'Mayor Valor' }}</a>
                    <a class="tab" :class="{ 'tab-active': activeTab === 'rare' }" @click="activeTab = 'rare'">{{ $t('dashboard.tabs.rare') || 'Top Rareza' }}</a>
                </div>

                <div class="divide-y divide-base-200" v-if="activeTab === 'recent'">
                    <div v-for="coin in stats.recent_coins" :key="coin.id" 
                         class="flex items-center gap-3 py-4 cursor-pointer hover:bg-base-200 transition-colors rounded-lg px-2 -mx-2"
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
                         class="flex items-center gap-3 py-4 cursor-pointer hover:bg-base-200 transition-colors rounded-lg px-2 -mx-2"
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

                <div class="divide-y divide-base-200" v-if="activeTab === 'rare'">
                    <div v-for="coin in stats.rarest_coins" :key="coin.id" 
                         class="flex items-center gap-3 py-4 cursor-pointer hover:bg-base-200 transition-colors rounded-lg px-2 -mx-2"
                         @click="router.push(`/coin/${coin.id}`)">
                         <div class="flex-1 min-w-0">
                            <div class="font-bold truncate">{{ coin.name }}</div>
                            <div class="text-xs opacity-50">{{ coin.year }} • {{ coin.country }}</div>
                         </div>
                         <div class="badge badge-accent badge-sm">{{ formatMintage(coin.mintage) }}</div>
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
                    <Bar v-if="timelineChartData" :data="timelineChartData" :options="barOptions" />
                </div>
                <div v-if="stats.oldest_coin" class="alert alert-info mt-4">
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="stroke-current shrink-0 w-6 h-6"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
                    <div>
                        <h3 class="font-bold">{{ $t('dashboard.stats.oldest') }}</h3>
                        <div class="text-xs">{{ stats.oldest_coin.year }}: {{ stats.oldest_coin.name }}</div>
                    </div>
                </div>
            </div>
        </div>

        <!-- Quality Matrix (Year vs Grade) -->
        <div class="card bg-base-100 shadow-xl">
            <div class="card-body">
                <h2 class="card-title">{{ $t('dashboard.charts.quality') }}</h2>
                <div class="h-64 relative">
                    <Scatter v-if="qualityChartData" :data="qualityChartData" :options="scatterOptions" />
                </div>
                <!-- Oldest High Grade Indicator -->
                <div v-if="stats.oldest_high_grade_coin" class="alert alert-success mt-4 bg-opacity-20">
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="stroke-current shrink-0 w-6 h-6"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
                    <div>
                        <h3 class="font-bold text-success-content">{{ $t('dashboard.stats.oldest_gem') || 'Joya Antigua (Calidad Alta)' }}</h3>
                        <div class="text-sm text-success-content opacity-90">{{ stats.oldest_high_grade_coin.year }}: {{ stats.oldest_high_grade_coin.name }} ({{ stats.oldest_high_grade_coin.grade }})</div>
                    </div>
                </div>
            </div>
        </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <!-- Combined Storage Distribution & Stats -->
        <div class="card bg-base-100 shadow-xl lg:col-span-3">
            <div class="card-body">
                <h2 class="card-title">{{ $t('dashboard.charts.storage') }}</h2>
                <div class="grid grid-cols-1 md:grid-cols-3 gap-8 items-center">
                    <!-- Chart Column (Left ~33%) -->
                    <div class="h-64 relative flex justify-center md:col-span-1">
                        <Doughnut v-if="storageChartData" :data="storageChartData" :options="storageChartOptions" />
                    </div>

                    <!-- Table Column (Right ~66%) -->
                    <div class="md:col-span-2 overflow-x-auto" v-if="stats.group_stats && stats.group_stats.length > 0">
                         <table class="table w-full">
                            <thead>
                                <tr>
                                    <th>{{ $t('common.name') || 'Nombre' }}</th>
                                    <th class="text-center">{{ $t('common.count') || 'Cantidad' }}</th>
                                    <th class="text-right">{{ $t('common.value_range') || 'Rango de Valor' }}</th>
                                    <th></th>
                                </tr>
                            </thead>
                            <tbody>
                                <tr v-for="(group, index) in stats.group_stats" :key="group.group_id" 
                                    class="cursor-pointer transition-colors duration-200"
                                    :class="{
                                        'bg-base-200': hoveredStorage === group.group_name,
                                        'opacity-50': hoveredStorage && hoveredStorage !== group.group_name
                                    }"
                                    @mouseenter="hoveredStorage = group.group_name"
                                    @mouseleave="hoveredStorage = null"
                                    @click="router.push({ path: '/list', query: { group_id: group.group_id } })">
                                    <td class="font-bold flex items-center gap-2">
                                        <!-- Color Indicator -->
                                        <div class="w-3 h-3 rounded-full" 
                                             :style="{ backgroundColor: storageChartData.datasets[0].backgroundColor[index] ? storageChartData.datasets[0].backgroundColor[index].replace('40', '') : '#ccc' }">
                                        </div>
                                        {{ group.group_name }}
                                    </td>
                                    <td class="text-center">
                                        <div class="badge badge-secondary badge-outline">{{ group.count }}</div>
                                    </td>
                                    <td class="text-right font-mono text-sm">
                                        <span class="text-success">{{ formatCurrency(group.min_value) }}</span>
                                        <span class="mx-1 text-base-content/50">-</span>
                                        <span class="text-primary font-bold">{{ formatCurrency(group.max_value) }}</span>
                                    </td>
                                    <td class="text-right">
                                        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-4 h-4 ml-auto opacity-50">
                                            <path stroke-linecap="round" stroke-linejoin="round" d="M8.25 4.5l7.5 7.5-7.5 7.5" />
                                        </svg>
                                    </td>
                                </tr>
                            </tbody>
                        </table>
                    </div>
                </div>
            </div>
        </div>


    </div>








    <!-- Bottom Row: Trivia & Random Coin -->
    <div class="grid grid-cols-1 lg:grid-cols-4 gap-8">
        <!-- Trivia & Physical Stats -->
        <div class="card bg-base-100 shadow-xl lg:col-span-1">
            <div class="card-body">
                <h2 class="card-title text-sm uppercase text-base-content/70">{{ $t('dashboard.lists.trivia') }}</h2>
                <div class="space-y-4">
                    <div class="stats stats-vertical shadow w-full">
                        <div class="stat">
                            <div class="stat-title">{{ $t('dashboard.stats.silver_reserve') }}</div>
                            <div class="stat-value text-primary text-2xl">{{ stats.total_silver_weight ? stats.total_silver_weight.toFixed(1) : 0 }}g</div>
                            <div class="stat-desc">{{ $t('dashboard.stats.pure_silver') }}</div>
                        </div>
                        <div class="stat">
                            <div class="stat-title">{{ $t('dashboard.stats.gold_reserve') }}</div>
                            <div class="stat-value text-warning text-2xl">{{ stats.total_gold_weight ? stats.total_gold_weight.toFixed(1) : 0 }}g</div>
                            <div class="stat-desc">{{ $t('dashboard.stats.pure_gold') }}</div>
                        </div>
                    </div>
                    
                    <div v-if="stats.heaviest_coin" class="text-sm">
                        <span class="font-bold">{{ $t('dashboard.stats.heaviest') }}</span> {{ stats.heaviest_coin.name }} ({{ stats.heaviest_coin.weight_g }}g)
                    </div>
                    <div v-if="stats.smallest_coin" class="text-sm">
                        <span class="font-bold">{{ $t('dashboard.stats.smallest') }}</span> {{ stats.smallest_coin.name }} ({{ stats.smallest_coin.diameter_mm }}mm)
                    </div>
                </div>
            </div>
        </div>

        <!-- Random Coin Feature -->
        <div v-if="stats.random_coin" class="card lg:card-side bg-base-100 shadow-xl overflow-hidden lg:col-span-3">
            <figure class="lg:w-1/3 bg-base-200 flex flex-col sm:flex-row lg:flex-col items-center justify-center gap-4 p-8">
                <div class="flex flex-col sm:flex-row lg:flex-col gap-4 items-center">
                    <div class="avatar cursor-pointer" @click="router.push(`/coin/${stats.random_coin.id}`)">
                        <div class="mask mask-circle w-40 h-40 flex items-center justify-center">
                            <img :src="getThumbnail(stats.random_coin, 'front')" class="object-contain w-full h-full" />
                        </div>
                    </div>
                    <div class="avatar cursor-pointer" @click="router.push(`/coin/${stats.random_coin.id}`)">
                        <div class="mask mask-circle w-40 h-40 flex items-center justify-center">
                            <img :src="getThumbnail(stats.random_coin, 'back')" class="object-contain w-full h-full" />
                        </div>
                    </div>
                </div>
            </figure>
            <div class="card-body">
                <div class="flex justify-between items-start">
                    <div>
                        <!-- Kicker -->
                        <div class="text-[0.65rem] uppercase tracking-widest text-base-content/50 font-bold mb-1">
                            {{ $t('dashboard.random.kicker') || 'MONEDA DESTACADA' }}
                        </div>
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
                        <div class="opacity-50">{{ $t('dashboard.random.face_value') }}</div>
                        <div class="font-bold">{{ stats.random_coin.face_value }} {{ stats.random_coin.currency }}</div>
                    </div>
                    <div>
                        <div class="opacity-50">{{ $t('dashboard.random.mint') }}</div>
                        <div class="font-bold">{{ stats.random_coin.mint || '-' }}</div>
                    </div>
                    <div>
                        <div class="opacity-50">{{ $t('dashboard.random.mintage') }}</div>
                        <div class="font-bold">{{ formatMintage(stats.random_coin.mintage) }}</div>
                    </div>
                    <div>
                        <div class="opacity-50">{{ $t('dashboard.random.km_code') }}</div>
                        <div class="font-bold">{{ stats.random_coin.km_code || '-' }}</div>
                    </div>
                    <div v-if="stats.random_coin.group_name">
                        <div class="opacity-50">{{ $t('dashboard.random.collection') }}</div>
                        <div class="badge badge-outline mt-1">{{ stats.random_coin.group_name }}</div>
                    </div>
                </div>

                <div class="card-actions justify-end mt-4">
                    <button class="btn btn-ghost border-primary/20 text-primary hover:bg-primary/10 gap-2" @click="router.push(`/coin/${stats.random_coin.id}`)">
                        {{ $t('common.view_details') }}
                        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-4 h-4">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M13.5 4.5L21 12m0 0l-7.5 7.5M21 12H3" />
                        </svg>
                    </button>
                </div>
            </div>
        </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import axios from 'axios'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
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
import { GRADE_ORDER, getGradeColor, getGradeValue } from '../utils/grades'

ChartJS.register(CategoryScale, LinearScale, BarElement, Title, Tooltip, Legend, ArcElement, PointElement)



const router = useRouter()
const { t } = useI18n()

const API_URL = import.meta.env.VITE_API_URL || '/api/v1'
const activeTab = ref('recent')
const loading = ref(true)

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
  group_stats: [],
  all_coins: []
})

const barOptions = computed(() => ({
  responsive: true,
  maintainAspectRatio: false,
  scales: {
    y: {
      beginAtZero: true,
      ticks: {
        precision: 0  // Force integer values only
      }
    }
  },
  plugins: {
    legend: { display: false }
  },
  onClick: (event, elements, chart) => {
    if (elements.length > 0) {
      const index = elements[0].index
      const label = chart.data.labels[index]
      
      // Detect chart type by checking dataset label or data structure
      const datasetLabel = chart.data.datasets[0]?.label || ''
      
      // Timeline chart - labels are decade numbers like "1990", "2000"
      if (label && !isNaN(parseInt(label)) && parseInt(label) >= 1800) {
        const decade = parseInt(label)
        router.push({ path: '/list', query: { min_year: decade, max_year: decade + 9 } })
      }
      // Value distribution - labels contain "€"
      else if (label && label.includes('€')) {
        const cleanLabel = label.replace(/€/g, '').replace(/\s/g, '').trim()
        if (cleanLabel.includes('+')) {
          const min = cleanLabel.replace('+', '')
          router.push({ path: '/list', query: { min_price: min } })
        } else if (cleanLabel.includes('-')) {
          const [min, max] = cleanLabel.split('-')
          router.push({ path: '/list', query: { min_price: min, max_price: max } })
        }
      }
      // Grade distribution - labels are grade codes (FDC, SC, EBC, etc.)
      else if (label && ['FDC', 'SC', 'EBC', 'MBC', 'BC', 'RC', 'MC'].includes(label)) {
        router.push({ path: '/list', query: { grade: label } })
      }
    }
  }
}))
const chartOptions = computed(() => ({
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: { display: false },
    tooltip: {
       callbacks: {
           label: (context) => {
               let label = context.dataset.label || '';
               if (label) {
                   label += ': ';
               }
               if (context.parsed.y !== null) {
                   // Check if this chart displays currency values (Value Distribution)
                   // We guess based on the variable name if possible, or we make a specific option for Value Distribution.
                   // Since this options object is shared, we might need separate options for the Value Chart.
                   // For now, let's just format as number.
                   // WAIT, user asked for Currency in charts.
                   // Only Value Chart needs currency. Grade/Timeline do not.
                   // I should split chartOptions.
                   label += context.parsed.y;
               }
               return label;
           }
       }
    }
  },
  onClick: (event, elements, chart) => {
     if (elements.length > 0) {
        const index = elements[0].index
        // Timeline Chart (Decades) logic
        // Need to check if this is indeed the timeline chart.
        // The options are shared, which is tricky.
        // However, this `chartOptions` is only used for Timeline currently.
        // Timeline labels are "1990s", "2000s"
        const label = chart.data.labels[index]
        if (label.endsWith('s')) {
            const decade = parseInt(label.replace('s', ''))
            if (!isNaN(decade)) {
                router.push({ path: '/list', query: { min_year: decade, max_year: decade + 9 } })
            }
        }
     }
  }
}))

const valueChartOptions = computed(() => ({
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: { display: false },
    tooltip: {
       callbacks: {
           label: (context) => {
               let label = context.dataset.label || '';
               if (label) {
                   label += ': ';
               }
               if (context.parsed.y !== null) {
                   // This IS the value chart but the Y axis is COUNT (coins), not Value.
                   // The X axis is the bins (0-10, 10-50).
                   // User said: "reference to value should show € in tooltips and axis labels"
                   // For Value Distribution: X axis are ranges of currency. Y axis is Count.
                   // So the TOOLTIP for the BAR represents COUNT.
                   // But the LABELS on X axis represent CURRENCY.
                   // The labels are strings "0-10", "10-50". I can format those in the computed data if needed, or just leave them.
                   // But wait, user said "reference to value should show €".
                   // Maybe they mean charts where the VALUE is plotted?
                   // The scatter plot is Grade vs Year.
                   // The tabs list shows values.
                   // The Value Distribution X labels are "0-10", "10-50". These imply €.
                   // Let's add € to these labels in the data computation.
                   return label + context.parsed.y;
               }
               return label;
           }
       }
    }
  },
  onClick: (event, elements, chart) => {
     if (elements.length > 0) {
        const index = elements[0].index
        const label = chart.data.labels[index]
        // Label format is "0-10 €" or "500+ €"
        // We need to parse min_price and max_price
        // Common regex or split
        const cleanLabel = label.replace(' €', '').trim()
        if (cleanLabel.includes('+')) {
            const min = cleanLabel.replace('+', '')
            router.push({ path: '/list', query: { min_price: min } })
        } else if (cleanLabel.includes('-')) {
            const parts = cleanLabel.split('-')
            if (parts.length === 2) {
                router.push({ path: '/list', query: { min_price: parts[0], max_price: parts[1] } })
            }
        }
     }
  }
}))

const doughnutOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: { position: 'right' }
  },
  onClick: (event, elements, chart) => {
     if (elements.length > 0) {
        const index = elements[0].index
        const label = chart.data.labels[index]
        router.push({ path: '/list', query: { material: label } })
     }
  }
}

onMounted(async () => {
    try {
        const response = await axios.get(`${API_URL}/dashboard`)
        stats.value = response.data
    } catch (error) {
        console.error('Error fetching stats:', error)
    } finally {
        loading.value = false
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

const scatterOptions = computed(() => ({
    responsive: true,
    maintainAspectRatio: false,
    scales: {
        x: {
            type: 'linear',
            position: 'bottom',
            title: { display: true, text: t('dashboard.charts.labels.year') },
            ticks: {
                callback: function(value) {
                    return value.toString().replace(',', '');
                }
            }
        },
        y: {
            type: 'linear',
            title: { display: true, text: t('dashboard.charts.labels.grade') },
            min: 0,
            max: 75,
            ticks: {
                stepSize: 10,
                callback: function(value) {
                    // Map numeric values to grade labels
                    const gradeMap = {
                        70: 'FDC',
                        60: 'SC',
                        50: 'EBC',
                        40: 'MBC',
                        30: 'BC',
                        20: 'RC',
                        10: 'MC'
                    }
                    return gradeMap[value] || ''
                }
            }
        }
    },
    plugins: {
        legend: { display: false },
        tooltip: {
            callbacks: {
                label: (context) => {
                    const gradeMap = {
                        70: 'FDC',
                        60: 'SC',
                        50: 'EBC',
                        40: 'MBC',
                        30: 'BC',
                        20: 'RC',
                        10: 'MC'
                    }
                    const grade = gradeMap[context.parsed.y] || context.parsed.y
                    return `${context.raw.name} (${context.parsed.x}) - ${grade}`
                }
            }
        }
    },
    onClick: (event, elements, chart) => {
        if (elements.length > 0) {
            const index = elements[0].index
            const point = chart.data.datasets[0].data[index]
            // This assumes we have the coin ID or can access it.
            // In qualityChartData mapping we only added x, y, name, grade. 
            // We need to add ID to the data point.
            if (point.id) {
                router.push(`/coin/${point.id}`)
            }
        }
    }
}))

// Chart Data Computed Properties

const valueChartData = computed(() => {
  const dist = stats.value.value_distribution
  if (!dist) return null
  
  // Get all coin values to determine dynamic ranges
  const allValues = stats.value.all_coins?.map(c => c.max_value).filter(v => v > 0) || []
  if (allValues.length === 0) return null
  
  const maxValue = Math.max(...allValues)
  
  // Create dynamic ranges based on max value
  let rawLabels = []
  if (maxValue <= 10) {
    rawLabels = ['0-2', '2-5', '5-10']
  } else if (maxValue <= 50) {
    rawLabels = ['0-10', '10-20', '20-50']
  } else if (maxValue <= 100) {
    rawLabels = ['0-10', '10-50', '50-100']
  } else if (maxValue <= 500) {
    rawLabels = ['0-50', '50-100', '100-500']
  } else {
    rawLabels = ['0-100', '100-500', '500+']
  }
  
  // Count coins in each range
  const distribution = {}
  rawLabels.forEach(r => distribution[r] = 0)
  
  allValues.forEach(value => {
    for (const range of rawLabels) {
      const parts = range.split('-')
      const min = parseFloat(parts[0])
      const max = parts[1] === '+' ? Infinity : parseFloat(parts[1])
      if (value >= min && value < max) {
        distribution[range]++
        break
      } else if (range.endsWith('+') && value >= min) {
        distribution[range]++
        break
      }
    }
  })
  
  // Format labels with €
  const labels = rawLabels.map(l => {
    if (l.includes('+')) {
      return `${l.replace('+', '')}€+`
    } else {
      const [min, max] = l.split('-')
      return `${min}€ - ${max}€`
    }
  })
  
  const data = rawLabels.map(l => distribution[l] || 0)

  return {
    labels,
    datasets: [{
      label: t('dashboard.charts.labels.coins'),
      backgroundColor: '#7c3aed',
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
  
  // Use all base grades to ensure strict ordering
  const labels = [...GRADE_ORDER]
  
  const data = labels.map(l => dist[l] || 0)
  
  const colors = labels.map(label => getGradeColor(label))

  return {
    labels,
    datasets: [{
      label: t('dashboard.charts.labels.coins'),
      backgroundColor: colors,
      data
    }]
  }
})

const gradeChartOptions = computed(() => ({
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: { display: false },
    tooltip: {
        callbacks: {
            title: (context) => {
                const label = context[0].label
                return `${label} - ${t('grades.' + label + '.name')}`
            },
            label: (context) => {
                 return `${t('dashboard.charts.labels.coins')}: ${context.parsed.y}`
            },
            afterBody: (context) => {
                const label = context[0].label
                return t('grades.' + label + '.desc')
            }
        }
    }
  },
  onClick: (event, elements, chart) => {
     if (elements.length > 0) {
        const index = elements[0].index
        const label = chart.data.labels[index] // Grade e.g. "EBC"
        router.push({ path: '/list', query: { grade: label } })
     }
  }
}))

const timelineChartData = computed(() => {
    const dist = stats.value.decade_distribution
    if (!dist) return null

    // Sort decades and find range
    const existingDecades = Object.keys(dist).map(Number).sort((a, b) => a - b)
    if (existingDecades.length === 0) return null
    
    const minDecade = existingDecades[0]
    const maxDecade = existingDecades[existingDecades.length - 1]
    
    // Create all decades in range
    const allDecades = []
    for (let decade = minDecade; decade <= maxDecade; decade += 10) {
        allDecades.push(decade)
    }
    
    // Map data, filling 0 for missing decades
    const labels = allDecades.map(d => d.toString())
    const data = allDecades.map(d => dist[d.toString()] || 0)

    return {
        labels,
        datasets: [{
            label: t('dashboard.charts.labels.coins'),
            backgroundColor: '#3b82f6',
            data
        }]
    }
})

const qualityChartData = computed(() => {
    const coins = stats.value.all_coins
    
    // Debug logging
    console.log('Quality chart - all_coins:', coins)
    console.log('Quality chart - all_coins length:', coins ? coins.length : 0)
    
    if (!coins || coins.length === 0) {
        console.log('Quality chart - no coins data')
        return null
    }

    const data = coins
        .filter(c => c.year > 0 && c.grade)
        .map(c => {
            const gradeValue = getGradeValue(c.grade)
            console.log(`Coin ${c.name}: grade=${c.grade}, gradeValue=${gradeValue}, year=${c.year}`)
            return {
                x: c.year,
                y: gradeValue || 0,
                name: c.name,
                grade: c.grade,
                id: c.id
            }
        })
        .filter(point => point.y > 0)
    
    console.log('Quality chart - final data points:', data.length)
    console.log('Quality chart - data:', data)
    
    if (data.length === 0) {
        console.log('Quality chart - no valid data points after filtering')
        return null
    }

    return {
        datasets: [{
            label: t('dashboard.charts.labels.quality_vs_year'),
            backgroundColor: '#f59e0b',
            data
        }]
    }
})

const hoveredStorage = ref(null)

const storageChartData = computed(() => {
    const dist = stats.value.group_distribution
    if (!dist) return null

    const labels = Object.keys(dist)
    const data = Object.values(dist)
    // Base colors
    const baseColors = [
        '#3b82f6', '#10b981', '#f59e0b', '#ef4444', '#8b5cf6', '#ec4899', '#6366f1'
    ]
    
    // Create background colors based on hover state
    // If hoveredStorage is set (Group Name), dim others
    const bgColors = labels.map((label, index) => {
        const color = baseColors[index % baseColors.length]
        if (hoveredStorage.value && hoveredStorage.value !== label) {
            return color + '40' // Add transparency (hex alpha) to dim
        }
        return color
    })

    const borderColors = labels.map((label) => {
         // Highlight border if active
        if (hoveredStorage.value === label) return '#ffffff'
        return 'transparent'
    })

    const borderWidths = labels.map((label) => {
        if (hoveredStorage.value === label) return 2
        return 0
    })

    return {
        labels,
        datasets: [{
            backgroundColor: bgColors,
            borderColor: borderColors,
            borderWidth: borderWidths,
            data,
            hoverOffset: 10
        }]
    }
})

const storageChartOptions = computed(() => ({
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
        legend: { display: false }, // We use the custom table as legend
        tooltip: {
            enabled: false // Disable default tooltip to rely on table? Or keep it? User didn't say disable. Keep for now.
             // Actually user said "Left... Donut... Right... Table... acting as explicit legend".
             // Let's keep tooltip for precise numbers on chart.
        }
    },
    onHover: (event, elements) => {
        if (elements && elements.length > 0) {
            // Chart.js 3/4 structure
            const index = elements[0].index
            const dist = stats.value.group_distribution
            if (dist) {
                const labels = Object.keys(dist)
                hoveredStorage.value = labels[index]
            }
        } else {
            hoveredStorage.value = null
        }
    }
}))


</script>
