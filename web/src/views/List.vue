<template>
  <div>
    <div class="flex flex-col md:flex-row justify-between items-center mb-6 gap-4">
      <h2 class="text-3xl font-bold">{{ $t('list.title') }}</h2>
      
      <div class="flex flex-wrap gap-2 w-full md:w-auto items-center">
        <!-- Active Filters Indicator -->
        <div v-if="hasActiveFilters" class="flex items-center gap-2">
          <div class="badge badge-primary gap-1">
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-3 h-3">
              <path stroke-linecap="round" stroke-linejoin="round" d="M12 3c2.755 0 5.455.232 8.083.678.533.09.917.556.917 1.096v1.044a2.25 2.25 0 01-.659 1.591l-5.432 5.432a2.25 2.25 0 00-.659 1.591v2.927a2.25 2.25 0 01-1.244 2.013L9.75 21v-6.568a2.25 2.25 0 00-.659-1.591L3.659 7.409A2.25 2.25 0 013 5.818V4.774c0-.54.384-1.006.917-1.096A48.32 48.32 0 0112 3z" />
            </svg>
            Filtros activos
          </div>
          <button @click="clearFilters" class="btn btn-ghost btn-xs gap-1">
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-3 h-3">
              <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
            </svg>
            Limpiar
          </button>
        </div>
        
        <!-- Search (always visible) -->
        <input type="text" v-model="filters.query" :placeholder="$t('list.search_placeholder')" class="input input-bordered w-full md:w-48 input-sm" :class="{ 'input-primary': filters.query }" />
        
        <!-- Filter Toggle Button (all screens) -->
        <button @click="showMobileFilters = !showMobileFilters" class="btn btn-sm btn-outline gap-2">
          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-4 h-4">
            <path stroke-linecap="round" stroke-linejoin="round" d="M10.5 6h9.75M10.5 6a1.5 1.5 0 11-3 0m3 0a1.5 1.5 0 10-3 0M3.75 6H7.5m3 12h9.75m-9.75 0a1.5 1.5 0 01-3 0m3 0a1.5 1.5 0 00-3 0m-3.75 0H7.5m9-6h3.75m-3.75 0a1.5 1.5 0 01-3 0m3 0a1.5 1.5 0 00-3 0m-9.75 0h9.75" />
          </svg>
          Filtros
          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-4 h-4 transition-transform" :class="{ 'rotate-180': showMobileFilters }">
            <path stroke-linecap="round" stroke-linejoin="round" d="M19.5 8.25l-7.5 7.5-7.5-7.5" />
          </svg>
        </button>
        
        <!-- Advanced Filters (collapsible on all screens) -->
        <div class="w-full md:w-auto contents" v-show="showMobileFilters">
        <!-- Group Filter -->
        <select v-model="filters.group_id" class="select select-bordered select-sm w-full md:w-40" :class="{ 'select-primary': filters.group_id }">
          <option value="">{{ $t('list.all_groups') }}</option>
          <option v-for="group in groups" :key="group.id" :value="group.id">{{ group.name }}</option>
        </select>

        <!-- Advanced Filters (Simple for now) -->
        <input type="number" v-model="filters.year" :placeholder="$t('list.filters.year')" class="input input-bordered input-sm w-24" :class="{ 'input-primary': filters.year }" />
        <input type="text" v-model="filters.country" :placeholder="$t('list.filters.country')" class="input input-bordered input-sm w-32" :class="{ 'input-primary': filters.country }" />
        
        <!-- Grade Filter -->
        <select v-model="filters.grade" class="select select-bordered select-sm w-32" :class="{ 'select-primary': filters.grade }">
          <option value="">{{ $t('list.filters.all_grades') }}</option>
          <option value="FDC">FDC</option>
          <option value="SC">SC</option>
          <option value="EBC">EBC</option>
          <option value="MBC">MBC</option>
          <option value="BC">BC</option>
          <option value="RC">RC</option>
        </select>
        
        <!-- Material Filter -->
        <select v-model="filters.material" class="select select-bordered select-sm w-36" :class="{ 'select-primary': filters.material }">
          <option value="">{{ $t('list.filters.all_materials') }}</option>
          <option v-for="material in uniqueMaterials" :key="material" :value="material">{{ material }}</option>
        </select>
        
        <input type="number" v-model="filters.min_price" :placeholder="$t('list.filters.min_price')" class="input input-bordered input-sm w-20" :class="{ 'input-primary': filters.min_price }" />
        <input type="number" v-model="filters.max_price" :placeholder="$t('list.filters.max_price')" class="input input-bordered input-sm w-20" :class="{ 'input-primary': filters.max_price }" />

        <!-- Sort -->
        <select v-model="filters.sort_by" class="select select-bordered select-sm w-full md:w-32">
          <option value="">{{ $t('list.sort.label') }}</option>
          <option value="year">{{ $t('list.sort.year') }}</option>
          <option value="min_value">{{ $t('list.sort.min_value') }}</option>
          <option value="max_value">{{ $t('list.sort.max_value') }}</option>
          <option value="created_at">{{ $t('list.sort.created_at') }}</option>
          <option value="country">{{ $t('list.sort.country') }}</option>
          <option value="name">{{ $t('list.sort.name') }}</option>
        </select>
        </div>
        
        <button class="btn btn-sm btn-square" @click="toggleOrder" :title="filters.order === 'asc' ? 'Ascending' : 'Descending'">
            <svg v-if="filters.order === 'asc'" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-4 h-4">
              <path stroke-linecap="round" stroke-linejoin="round" d="M3 4.5h14.25M3 9h9.75M3 13.5h9.75m4.5-4.5v12m0 0l-3.75-3.75M17.25 21L21 17.25" />
            </svg>
            <svg v-else xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-4 h-4">
              <path stroke-linecap="round" stroke-linejoin="round" d="M3 4.5h14.25M3 9h9.75M3 13.5h5.25m5.25-.75L17.25 9m0 0L21 12.75M17.25 9v12" />
            </svg>
        </button>

        <div class="join">
            <button class="join-item btn btn-sm" :class="{ 'btn-active': viewMode === 'grid' }" @click="viewMode = 'grid'">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z" /></svg>
            </button>
            <button class="join-item btn btn-sm" :class="{ 'btn-active': viewMode === 'table' }" @click="viewMode = 'table'">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 10h16M4 14h16M4 18h16" /></svg>
            </button>
        </div>
      </div>
    </div>
    
    <!-- Active Group Header -->
    <div v-if="filters.group_id && activeGroupName" class="mb-8 p-6 bg-primary/10 rounded-2xl border border-primary/20 flex flex-col md:flex-row items-center justify-between gap-4 animate-in fade-in slide-in-from-top-4 duration-500">
        <div class="flex items-center gap-4">
             <div class="p-3 bg-primary rounded-xl text-primary-content shadow-lg">
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-8 h-8"><path stroke-linecap="round" stroke-linejoin="round" d="M2.25 12.75V12A2.25 2.25 0 014.5 9.75h15A2.25 2.25 0 0121.75 12v.75m-8.69-6.44l-2.12-2.12a1.5 1.5 0 00-1.061-.44H4.5A2.25 2.25 0 002.25 6v12a2.25 2.25 0 002.25 2.25h15A2.25 2.25 0 0021.75 18V9a2.25 2.25 0 00-2.25-2.25h-5.379a1.5 1.5 0 01-1.06-.44z" /></svg>
             </div>
             <div>
                 <div class="text-sm font-bold uppercase tracking-wider opacity-60">Colección</div>
                 <h2 class="text-3xl font-black text-primary">{{ activeGroupName }}</h2>
                 <p class="opacity-70 text-sm mt-1" v-if="activeGroupDesc">{{ activeGroupDesc }}</p>
             </div>
        </div>
        <button @click="clearFilters" class="btn btn-outline btn-sm gap-2 hover:btn-error">
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-4 h-4"><path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" /></svg>
            Salir del grupo
        </button>
    </div>
    
    <div v-if="loading" class="flex justify-center p-10">
      <span class="loading loading-spinner loading-lg"></span>
    </div>

    <div v-else-if="coins.length === 0" class="flex flex-col items-center justify-center p-16 bg-base-100 rounded-box shadow-xl">
      <!-- Empty State Icon -->
      <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1" stroke="currentColor" class="w-24 h-24 text-base-content/20 mb-6">
        <path stroke-linecap="round" stroke-linejoin="round" d="M20.25 6.375c0 2.278-3.694 4.125-8.25 4.125S3.75 8.653 3.75 6.375m16.5 0c0-2.278-3.694-4.125-8.25-4.125S3.75 4.097 3.75 6.375m16.5 0v11.25c0 2.278-3.694 4.125-8.25 4.125s-8.25-1.847-8.25-4.125V6.375m16.5 0v3.75m-16.5-3.75v3.75m16.5 0v3.75C20.25 16.153 16.556 18 12 18s-8.25-1.847-8.25-4.125v-3.75m16.5 0c0 2.278-3.694 4.125-8.25 4.125s-8.25-1.847-8.25-4.125" />
      </svg>
      
      <!-- Empty State Message -->
      <h3 class="text-2xl font-bold mb-2">{{ $t('list.empty_state') }}</h3>
      <p class="text-base-content/60 mb-6 text-center max-w-md">
        {{ hasActiveFilters ? 'No se encontraron monedas con los filtros aplicados. Intenta ajustar los criterios de búsqueda.' : 'Tu colección está vacía. ¡Comienza añadiendo tu primera moneda!' }}
      </p>
      
      <!-- CTA Buttons -->
      <div class="flex gap-3">
        <button v-if="hasActiveFilters" @click="clearFilters" class="btn btn-outline gap-2">
          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-5 h-5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
          </svg>
          Limpiar Filtros
        </button>
        <router-link to="/add" class="btn btn-primary gap-2">
          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-5 h-5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15" />
          </svg>
          {{ $t('list.add_button') }}
        </router-link>
      </div>
    </div>

    <!-- Grid View -->
    <div v-else-if="viewMode === 'grid'" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <div v-for="coin in coins" :key="coin.id" class="card bg-base-100 shadow-xl hover:shadow-2xl transition-shadow cursor-pointer">
        <figure class="px-4 pt-4 relative group flex justify-center gap-2">
          <div class="relative group/img cursor-zoom-in" @click.stop="openViewer(coin, 'front')">
             <!-- Overlay for zoom hint -->
            <div class="absolute inset-0 bg-black bg-opacity-0 group-hover/img:bg-opacity-20 transition-all flex items-center justify-center z-10 rounded-full">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-white opacity-0 group-hover/img:opacity-100 transition-opacity" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0zM10 7v3m0 0v3m0-3h3m-3 0H7" />
                </svg>
            </div>
            <img :src="getThumbnailUrl(coin, 'front')" alt="Coin Front" class="rounded-full h-24 w-24 object-cover shadow-md hover:scale-110 transition-transform duration-300" />
          </div>
          <div class="relative group/img cursor-zoom-in" @click.stop="openViewer(coin, 'back')">
             <!-- Overlay for zoom hint -->
            <div class="absolute inset-0 bg-black bg-opacity-0 group-hover/img:bg-opacity-20 transition-all flex items-center justify-center z-10 rounded-full">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-white opacity-0 group-hover/img:opacity-100 transition-opacity" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0zM10 7v3m0 0v3m0-3h3m-3 0H7" />
                </svg>
            </div>
            <img :src="getThumbnailUrl(coin, 'back')" alt="Coin Back" class="rounded-full h-24 w-24 object-cover shadow-md hover:scale-110 transition-transform duration-300" />
          </div>
        </figure>
        <div class="card-body" @click="goToDetail(coin.id)">
          <h2 class="card-title flex-col items-start gap-1">
            <span v-if="coin.name" class="text-lg font-bold text-primary">{{ coin.name }}</span>
            <span class="text-base font-normal">
                {{ coin.country }} {{ coin.face_value }}
                <div class="badge badge-secondary ml-2" v-if="coin.year && coin.year !== 0">{{ coin.year }}</div>
            </span>
          </h2>
          <p class="text-sm text-gray-500">{{ coin.currency }}</p>
          <div v-if="coin.min_value > 0 || coin.max_value > 0" class="mt-1 font-bold text-success">
            {{ formatPriceRange(coin) }}
          </div>
          <div class="card-actions justify-between items-center mt-2">
            <div class="flex gap-2 flex-wrap">
              <div class="badge badge-outline" v-if="coin.grade">{{ coin.grade }}</div>
              <div class="tooltip" :data-tip="coin.material">
                <div class="badge badge-outline truncate max-w-[120px]">{{ coin.material?.split('(')[0].trim() || coin.material }}</div>
              </div>
            </div>
            <div class="flex gap-2">
                <button @click.stop="router.push(`/edit/${coin.id}`)" class="btn btn-square btn-sm btn-ghost text-info">
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-4 h-4">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M16.862 4.487l1.687-1.688a1.875 1.875 0 112.652 2.652L10.582 16.07a4.5 4.5 0 01-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 011.13-1.897l8.932-8.931zm0 0L19.5 7.125M18 14v4.75A2.25 2.25 0 0115.75 21H5.25A2.25 2.25 0 013 18.75V8.25A2.25 2.25 0 015.25 6H10" />
                    </svg>
                </button>
                <button @click.stop="confirmDelete(coin)" class="btn btn-square btn-sm btn-ghost text-error">
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-4 h-4">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M14.74 9l-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 01-2.244 2.077H8.084a2.25 2.25 0 01-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 00-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 013.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 00-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 00-7.5 0" />
                    </svg>
                </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Table View -->
    <div v-else class="overflow-x-auto bg-base-100 rounded-box shadow-xl">
        <table class="table table-zebra">
          <thead>
            <tr>
              <th>{{ $t('list.table.images') }}</th>
              <th class="cursor-pointer hover:bg-base-200" @click="toggleSort('name')">
                <div class="flex items-center gap-1">
                  {{ $t('list.table.name') }}
                  <svg v-if="filters.sort_by === 'name' && filters.order === 'asc'" xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7" />
                  </svg>
                  <svg v-else-if="filters.sort_by === 'name' && filters.order === 'desc'" xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                  </svg>
                </div>
              </th>
              <th class="hidden lg:table-cell">{{ $t('list.table.mint') }}</th>
              <th class="hidden lg:table-cell">{{ $t('list.table.mintage') }}</th>
              <th class="cursor-pointer hover:bg-base-200 hidden md:table-cell" @click="toggleSort('country')">
                <div class="flex items-center gap-1">
                  {{ $t('list.table.country') }}
                  <svg v-if="filters.sort_by === 'country' && filters.order === 'asc'" xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7" />
                  </svg>
                  <svg v-else-if="filters.sort_by === 'country' && filters.order === 'desc'" xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                  </svg>
                </div>
              </th>
              <th class="cursor-pointer hover:bg-base-200 hidden md:table-cell" @click="toggleSort('max_value')">
                <div class="flex items-center gap-1">
                  {{ $t('list.table.value') }}
                  <svg v-if="filters.sort_by === 'max_value' && filters.order === 'asc'" xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7" />
                  </svg>
                  <svg v-else-if="filters.sort_by === 'max_value' && filters.order === 'desc'" xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                  </svg>
                </div>
              </th>
              <th class="cursor-pointer hover:bg-base-200 hidden lg:table-cell" @click="toggleSort('year')">
                <div class="flex items-center gap-1">
                  {{ $t('list.table.year') }}
                  <svg v-if="filters.sort_by === 'year' && filters.order === 'asc'" xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7" />
                  </svg>
                  <svg v-else-if="filters.sort_by === 'year' && filters.order === 'desc'" xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                  </svg>
                </div>
              </th>
              <th class="hidden lg:table-cell">{{ $t('list.table.currency') }}</th>
              <th class="hidden lg:table-cell">{{ $t('list.table.grade') }}</th>
              <th class="hidden lg:table-cell">{{ $t('list.table.material') }}</th>
              <th>{{ $t('list.table.actions') }}</th>
            </tr>
          </thead>
        <tbody>
          <tr v-for="coin in coins" :key="coin.id" class="hover cursor-pointer" @click="goToDetail(coin.id)">
            <td>
                <div class="flex gap-2">
                    <div class="avatar cursor-zoom-in" @click.stop="openViewer(coin, 'front')">
                        <div class="w-12 h-12 rounded-full overflow-hidden">
                            <img :src="getThumbnailUrl(coin, 'front')" alt="Front" class="hover:scale-110 transition-transform duration-300" />
                        </div>
                    </div>
                    <div class="avatar cursor-zoom-in" @click.stop="openViewer(coin, 'back')">
                        <div class="w-12 h-12 rounded-full overflow-hidden">
                            <img :src="getThumbnailUrl(coin, 'back')" alt="Back" class="hover:scale-110 transition-transform duration-300" />
                        </div>
                    </div>
                </div>
            </td>
            <td>
                <div class="font-bold text-primary">{{ coin.name || '-' }}</div>
                <!-- Mobile-only info: show country and value below name on small screens -->
                <div class="text-xs opacity-70 md:hidden mt-1">
                    <div>{{ coin.country }}</div>
                    <div class="text-success font-semibold">{{ formatPriceRange(coin) }}</div>
                </div>
            </td>
            <td class="hidden lg:table-cell">{{ coin.mint || '-' }}</td>
            <td class="hidden lg:table-cell">{{ formatMintage(coin.mintage) }}</td>
            <td class="font-semibold hidden md:table-cell">{{ coin.country }}</td>
            <td class="hidden md:table-cell whitespace-nowrap font-bold text-success">
              {{ formatPriceRange(coin) }}
            </td>
            <td class="hidden lg:table-cell">{{ (coin.year && coin.year !== 0) ? coin.year : '-' }}</td>
            <td class="hidden lg:table-cell">{{ coin.currency }}</td>
            <td class="hidden lg:table-cell"><div class="badge badge-ghost" v-if="coin.grade">{{ coin.grade }}</div><span v-else>-</span></td>
            <td class="hidden lg:table-cell">
              <div class="tooltip" v-if="coin.material" :data-tip="coin.material">
                <div class="badge badge-ghost truncate max-w-[120px]">{{ coin.material?.split('(')[0].trim() || coin.material }}</div>
              </div>
              <span v-else>-</span>
            </td>
            <td>
                <div class="flex gap-1">
                    <button @click.stop="router.push(`/edit/${coin.id}`)" class="btn btn-square btn-sm btn-ghost text-info">
                        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-4 h-4">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M16.862 4.487l1.687-1.688a1.875 1.875 0 112.652 2.652L10.582 16.07a4.5 4.5 0 01-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 011.13-1.897l8.932-8.931zm0 0L19.5 7.125M18 14v4.75A2.25 2.25 0 0115.75 21H5.25A2.25 2.25 0 013 18.75V8.25A2.25 2.25 0 015.25 6H10" />
                        </svg>
                    </button>
                    <button @click.stop="confirmDelete(coin)" class="btn btn-square btn-sm btn-ghost text-error">
                        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-4 h-4">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M14.74 9l-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 01-2.244 2.077H8.084a2.25 2.25 0 01-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 00-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 013.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 00-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 00-7.5 0" />
                        </svg>
                    </button>
                </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Delete Modal -->
    <dialog id="delete_modal" class="modal" :class="{ 'modal-open': deleteModalOpen }">
      <div class="modal-box">
        <h3 class="font-bold text-lg text-error">{{ $t('list.delete_modal.title') }}</h3>
        <p class="py-4">{{ $t('list.delete_modal.confirm') }} <span class="font-bold">{{ coinToDelete?.name || $t('common.unknown') }}</span>? {{ $t('list.delete_modal.warning') }}</p>
        <div class="modal-action">
          <button class="btn" @click="deleteModalOpen = false">{{ $t('common.cancel') }}</button>
          <button class="btn btn-error" @click="deleteCoin" :disabled="deleting">
            <span v-if="deleting" class="loading loading-spinner"></span>
            {{ $t('common.delete') }}
          </button>
        </div>
      </div>
    </dialog>

    <ImageViewer 
      :is-open="viewerOpen" 
      :image-url="viewerImage" 
      @close="viewerOpen = false" 
    />
  </div>
</template>

<script setup>
import { ref, onMounted, watch, computed } from 'vue'
import axios from 'axios'
import { useRouter, useRoute } from 'vue-router'
import ImageViewer from '../components/ImageViewer.vue'
import { formatMintage } from '../utils/formatters'

const coins = ref([])
const groups = ref([])
const loading = ref(true)
const router = useRouter()
const showMobileFilters = ref(false)

// Detect mobile viewport
const isMobile = computed(() => {
  if (typeof window !== 'undefined') {
    return window.innerWidth < 768
  }
  return false
})

// View Mode with localStorage persistence
const viewMode = ref(localStorage.getItem('coinListViewMode') || 'table')

// Watch for changes and save to localStorage
watch(viewMode, (newMode) => {
  localStorage.setItem('coinListViewMode', newMode)
})
const API_URL = import.meta.env.VITE_API_URL || '/api/v1'
const STORAGE_URL = '' // Base URL for static files (relative)

const viewerOpen = ref(false)
const viewerImage = ref('')

// Helper to clean filter values from URL (remove €, +, spaces, etc.)
const cleanFilterValue = (value) => {
  if (!value) return ''
  return String(value).replace(/[€+\s]/g, '').trim()
}

import { useSettingsStore } from '../stores/settings'
import { storeToRefs } from 'pinia'

const formatPriceRange = (coin) => {
    const settingsStore = useSettingsStore()
    if (settingsStore.privacyMode) return '***'
    
    if (!coin.min_value && !coin.max_value) return ''
    if (coin.min_value === 0 && coin.max_value > 0) return `< ${coin.max_value}€`
    if (coin.min_value > 0 && coin.max_value > 0) return `${coin.min_value} - ${coin.max_value}€`
    if (coin.max_value > 0) return `${coin.max_value}€`
    if (coin.min_value > 0) return `> ${coin.min_value}€`
    return ''
}

// Filters
const route = useRoute()
const filters = ref({
    query: route.query.q || '',
    group_id: route.query.group_id || '',
    year: cleanFilterValue(route.query.year) || '',
    country: route.query.country || '',
    min_price: cleanFilterValue(route.query.min_price) || '',
    max_price: cleanFilterValue(route.query.max_price) || '',
    grade: route.query.grade || '',
    material: route.query.material || '',
    min_year: cleanFilterValue(route.query.min_year) || '',
    max_year: cleanFilterValue(route.query.max_year) || '',
    sort_by: route.query.sort_by || '',
    order: route.query.order || 'asc'
})

// Toggle sort function
const toggleSort = (column) => {
  if (filters.value.sort_by === column) {
    // Toggle order if same column
    filters.value.order = filters.value.order === 'asc' ? 'desc' : 'asc'
  } else {
    // New column, default to asc
    filters.value.sort_by = column
    filters.value.order = 'asc'
  }
}

// Computed property to check if any filters are active
const hasActiveFilters = computed(() => {
    return filters.value.query || 
           filters.value.group_id || 
           filters.value.year || 
           filters.value.country || 
           filters.value.min_price || 
           filters.value.max_price || 
           filters.value.grade || 
           filters.value.material || 
           filters.value.min_year || 
           filters.value.max_year
})

// Clear all filters
const clearFilters = () => {
    filters.value = {
        query: '',
        group_id: '',
        year: '',
        country: '',
        min_price: '',
        max_price: '',
        grade: '',
        material: '',
        min_year: '',
        max_year: '',
        sort_by: filters.value.sort_by, // Keep sort
        order: filters.value.order // Keep order
    }
}

const toggleOrder = () => {
    filters.value.order = filters.value.order === 'asc' ? 'desc' : 'asc'
}

// Debounce helper
let timeout = null
const debouncedFetch = () => {
    if (timeout) clearTimeout(timeout)
    timeout = setTimeout(() => {
        fetchCoins()
    }, 300)
}

watch(filters, () => {
    debouncedFetch()
}, { deep: true })

const fetchGroups = async () => {
    try {
        const res = await axios.get(`${API_URL}/groups`)
        groups.value = res.data
    } catch (e) {
        console.error("Failed to fetch groups", e)
    }
}

// Compute unique materials (text before parenthesis)
const uniqueMaterials = computed(() => {
  const materials = new Set()
  coins.value.forEach(coin => {
    if (coin.material) {
      // Extract text before parenthesis
      const cleanMaterial = coin.material.split('(')[0].trim()
      if (cleanMaterial) materials.add(cleanMaterial)
    }
  })
  return Array.from(materials).sort()
})

const activeGroupName = computed(() => {
    if (!filters.value.group_id) return null
    const group = groups.value.find(g => g.id === parseInt(filters.value.group_id)) || groups.value.find(g => g.id === filters.value.group_id) // Handle string/int types
    return group ? group.name : null
})

const activeGroupDesc = computed(() => {
    if (!filters.value.group_id) return null
    const group = groups.value.find(g => g.id === parseInt(filters.value.group_id)) || groups.value.find(g => g.id === filters.value.group_id)
    return group ? group.description : null
})

// Fetch coins
const fetchCoins = async () => {
    loading.value = true
    try {
        const params = {}
        if (filters.value.query) params.q = filters.value.query
        if (filters.value.group_id) params.group_id = filters.value.group_id
        if (filters.value.year) params.year = filters.value.year
        if (filters.value.country) params.country = filters.value.country
        if (filters.value.min_price) params.min_price = filters.value.min_price
        if (filters.value.max_price) params.max_price = filters.value.max_price
        if (filters.value.grade) params.grade = filters.value.grade
        if (filters.value.material) params.material = filters.value.material
        if (filters.value.min_year) params.min_year = filters.value.min_year
        if (filters.value.max_year) params.max_year = filters.value.max_year
        if (filters.value.sort_by) params.sort_by = filters.value.sort_by
        if (filters.value.order) params.sort_order = filters.value.order // Changed from 'order' to 'sort_order'

        const res = await axios.get(`${API_URL}/coins`, { params })
        coins.value = res.data || []
    } catch (e) {
        console.error("Failed to fetch coins", e)
        coins.value = [] // Ensure coins is an empty array on error
    } finally {
        loading.value = false
    }
}

// Delete Modal State
const deleteModalOpen = ref(false)
const coinToDelete = ref(null)
const deleting = ref(false)

const confirmDelete = (coin) => {
    coinToDelete.value = coin
    deleteModalOpen.value = true
}

const deleteCoin = async () => {
    if (!coinToDelete.value) return
    deleting.value = true
    try {
        await axios.delete(`${API_URL}/coins/${coinToDelete.value.id}`)
        coins.value = coins.value.filter(c => c.id !== coinToDelete.value.id)
        deleteModalOpen.value = false
        coinToDelete.value = null
    } catch (e) {
        console.error("Failed to delete coin", e)
        alert("Failed to delete coin")
    } finally {
        deleting.value = false
    }
}

const getImageUrl = (path) => {
    if (!path) return 'https://via.placeholder.com/150'
    // Hacky fix for absolute paths from docker
    // If path contains "storage", take everything after
    if (path.includes('storage/')) {
        return `${STORAGE_URL}/storage/${path.split('storage/')[1]}`
    }
    return path
}

const getThumbnailUrl = (coin, side = 'front') => {
    // Try to find thumbnail in images array
    if (coin.images && coin.images.length > 0) {
        const thumb = coin.images.find(img => img.image_type === 'thumbnail' && img.side === side)
        if (thumb) {
            return getImageUrl(thumb.path)
        }
    }
    // Fallback to sample image
    return side === 'front' ? getImageUrl(coin.sample_image_url_front) : getImageUrl(coin.sample_image_url_back)
}

const getFullResUrl = (coin, side = 'front') => {
    // Try to find processed crop first, then original
    if (coin.images && coin.images.length > 0) {
        const processed = coin.images.find(img => img.image_type === 'crop' && img.side === side)
        if (processed) return getImageUrl(processed.path)
        
        const original = coin.images.find(img => img.image_type === 'original' && img.side === side)
        if (original) return getImageUrl(original.path)
    }
    return side === 'front' ? getImageUrl(coin.sample_image_url_front) : getImageUrl(coin.sample_image_url_back)
}

const openViewer = (coin, side) => {
    viewerImage.value = getFullResUrl(coin, side)
    viewerOpen.value = true
}


const goToDetail = (id) => {
  router.push(`/coin/${id}`)
}

onMounted(async () => {
    fetchGroups()
    fetchCoins()
})
</script>
