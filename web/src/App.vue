<template>
  <div class="min-h-screen bg-base-200">
    <div class="navbar bg-base-100 shadow-lg">
      <div class="flex-1">
        <router-link to="/" class="btn btn-ghost normal-case text-xl">NumismaticApp</router-link>
      </div>
      <div class="flex-none">
        <ul class="menu menu-horizontal px-1">
          <li><router-link to="/">Home</router-link></li>
          <li><router-link to="/list">Collection</router-link></li>
          <li><router-link to="/add" class="btn btn-primary btn-sm text-white ml-2">Add Coin</router-link></li>
        </ul>
      </div>
    </div>
    <div class="container mx-auto p-4">
      <router-view></router-view>
    </div>
  </div>
</template>

<script setup>
</script>
