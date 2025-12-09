<template>
  <div class="hero min-h-[50vh] bg-base-100 rounded-box shadow-xl">
    <div class="hero-content text-center">
      <div class="max-w-md">
        <h1 class="text-5xl font-bold">Your Collection</h1>
        <p class="py-6">Manage, analyze, and value your numismatic treasures with AI.</p>
        <div class="stats shadow w-full">
          <div class="stat">
            <div class="stat-title">Total Coins</div>
            <div class="stat-value text-primary">{{ totalCoins }}</div>
            <div class="stat-desc">In your collection</div>
          </div>
        </div>
        <div class="mt-6">
            <router-link to="/add" class="btn btn-primary">Add New Coin</router-link>
            <router-link to="/list" class="btn btn-ghost ml-2">View Gallery</router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'

const totalCoins = ref(0)
const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1'

onMounted(async () => {
  try {
    // Ideally we have a stats endpoint, for now we list and count (inefficient but works for MVP)
    // Or we can add a count endpoint.
    // Let's assume we fetch list with limit 1 just to see if it works, or implement count later.
    // For now, let's just show "Welcome"
    // Actually, let's try to fetch list
    const res = await axios.get(`${API_URL}/coins?limit=1`)
    // We don't have total count in response yet.
    // Let's just leave it as 0 or "..."
    totalCoins.value = "..."
  } catch (e) {
    console.error(e)
  }
})
</script>
