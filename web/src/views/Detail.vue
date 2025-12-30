<template>
  <div v-if="loading" class="flex justify-center items-center min-h-screen">
    <span class="loading loading-spinner loading-lg text-primary"></span>
  </div>
  
  <div v-else-if="coin" class="max-w-7xl mx-auto px-4 lg:px-8 py-6">
    <div class="grid grid-cols-1 lg:grid-cols-12 gap-8 relative items-start">
        
        <!-- LEFT COLUMN: Sticky Gallery & Quick Stats (40%) -->
        <div class="lg:col-span-5 lg:sticky lg:top-24 space-y-6">
            <!-- Main Gallery Card -->
            <div class="card bg-base-100 shadow-2xl border border-gray-800 overflow-hidden group">
                <!-- Source Toggle (Floating) -->
                <div class="absolute top-4 left-0 right-0 z-20 flex justify-center opacity-0 group-hover:opacity-100 transition-opacity duration-300 pointer-events-none gap-2">
                     <div class="join shadow-lg pointer-events-auto scale-90">
                        <button class="btn btn-sm join-item" :class="{ 'btn-primary': activeImageSource === 'processed' }" @click="activeImageSource = 'processed'">
                            {{ $t('details.toggles.processed') }}
                        </button>
                        <button class="btn btn-sm join-item" :class="{ 'btn-primary': activeImageSource === 'original' }" @click="activeImageSource = 'original'">
                            {{ $t('details.toggles.original') }}
                        </button>
                    </div>
                     <!-- Sample Link (External) -->
                    <a v-if="coin.numista_details && (coin.numista_details.obverse_thumbnail || coin.numista_details.reverse_thumbnail)" 
                       :href="getNumistaUrl()" 
                       target="_blank"
                       class="btn btn-sm btn-accent shadow-lg pointer-events-auto scale-90 gap-1"
                       title="Ver ejemplo en Numista">
                        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-4 h-4"><path stroke-linecap="round" stroke-linejoin="round" d="M13.5 6H5.25A2.25 2.25 0 003 8.25v10.5A2.25 2.25 0 005.25 21h10.5A2.25 2.25 0 0018 18.75V10.5m-10.5 6L21 3m0 0h-5.25M21 3v5.25" /></svg>
                        Ejemplo
                    </a>
                </div>

                <!-- Combined Image Display (Flip-like or Grid) -->
                <div class="relative w-full bg-gray-900 flex items-start justify-center p-4 gap-8">
                    <!-- Front Column -->
                    <div class="flex flex-col items-center gap-1">
                        <figure class="relative group/img cursor-zoom-in transition-all duration-300 hover:scale-105" @click="openViewer('front')">
                            <img 
                                :src="getCurrentImageUrl('front')" 
                                class="w-32 h-32 sm:w-48 sm:h-48 object-contain drop-shadow-2xl"
                                :class="{ 'rounded-full': activeImageSource !== 'original', 'rounded-xl': activeImageSource === 'original' }"
                                :style="{ transform: `rotate(${rotations.front}deg)` }"
                                alt="Anverso"
                            />
                            <button @click.stop="openRotationEditor('front')" class="absolute bottom-0 right-0 btn btn-circle btn-sm btn-neutral opacity-0 group-hover/img:opacity-100 transition-opacity z-10 shadow-lg border border-gray-600" title="Rotar">
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" /></svg>
                            </button>
                        </figure>
                        <!-- Sample Link -->
                         <a v-if="coin.numista_details && coin.numista_details.obverse?.picture" 
                            :href="coin.numista_details.obverse.picture" 
                            target="_blank" 
                            class="link link-xs link-accent no-underline hover:underline flex items-center gap-1 mt-1 opacity-80 hover:opacity-100">
                            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-3 h-3"><path stroke-linecap="round" stroke-linejoin="round" d="M13.5 6H5.25A2.25 2.25 0 003 8.25v10.5A2.25 2.25 0 005.25 21h10.5A2.25 2.25 0 0018 18.75V10.5m-10.5 6L21 3m0 0h-5.25M21 3v5.25" /></svg>
                            Muestra
                        </a>
                    </div>

                    <!-- Back Column -->
                    <div class="flex flex-col items-center gap-1">
                         <figure class="relative group/img cursor-zoom-in transition-all duration-300 hover:scale-105" @click="openViewer('back')">
                            <img 
                                :src="getCurrentImageUrl('back')" 
                                class="w-32 h-32 sm:w-48 sm:h-48 object-contain drop-shadow-2xl"
                                :class="{ 'rounded-full': activeImageSource !== 'original', 'rounded-xl': activeImageSource === 'original' }"
                                :style="{ transform: `rotate(${rotations.back}deg)` }"
                                alt="Reverso"
                            />
                             <button @click.stop="openRotationEditor('back')" class="absolute bottom-0 right-0 btn btn-circle btn-sm btn-neutral opacity-0 group-hover/img:opacity-100 transition-opacity z-10 shadow-lg border border-gray-600" title="Rotar">
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" /></svg>
                            </button>
                        </figure>
                        <!-- Sample Link -->
                        <a v-if="coin.numista_details && coin.numista_details.reverse?.picture" 
                           :href="coin.numista_details.reverse.picture" 
                           target="_blank" 
                           class="link link-xs link-accent no-underline hover:underline flex items-center gap-1 mt-1 opacity-80 hover:opacity-100">
                           <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-3 h-3"><path stroke-linecap="round" stroke-linejoin="round" d="M13.5 6H5.25A2.25 2.25 0 003 8.25v10.5A2.25 2.25 0 005.25 21h10.5A2.25 2.25 0 0018 18.75V10.5m-10.5 6L21 3m0 0h-5.25M21 3v5.25" /></svg>
                           Muestra
                       </a>
                    </div>
                </div>
            </div>

            <!-- Dimensions Card (Compact) -->
            <div v-if="coin.diameter_mm > 0 || coin.weight_g > 0" class="card bg-base-100 shadow-xl border border-gray-800">
                <div class="card-body p-4">
                    <h3 class="text-xs font-bold text-gray-500 uppercase tracking-wider mb-4">{{ $t('details.sections.dimensions') }}</h3>
                    
                    <!-- SVG Diagram (Preserved) -->
                    <svg viewBox="0 0 400 320" class="w-full h-auto max-h-48" xmlns="http://www.w3.org/2000/svg">
                        <defs>
                            <linearGradient id="coinGradient" x1="0%" y1="0%" x2="100%" y2="100%">
                            <stop offset="0%" style="stop-color:#fbbf24;stop-opacity:0.3" />
                            <stop offset="50%" style="stop-color:#f59e0b;stop-opacity:0.5" />
                            <stop offset="100%" style="stop-color:#d97706;stop-opacity:0.3" />
                            </linearGradient>
                            <radialGradient id="coinShine" cx="30%" cy="30%">
                            <stop offset="0%" style="stop-color:#ffffff;stop-opacity:0.4" />
                            <stop offset="100%" style="stop-color:#ffffff;stop-opacity:0" />
                            </radialGradient>
                        </defs>
                        
                        <!-- Main coin circle -->
                        <circle cx="200" cy="160" r="90" fill="url(#coinGradient)" stroke="currentColor" stroke-width="2" class="text-primary" opacity="0.8"/>
                        <circle cx="200" cy="160" r="90" fill="url(#coinShine)" />
                        <circle cx="200" cy="160" r="85" fill="none" stroke="currentColor" stroke-width="1" class="text-primary" opacity="0.3" stroke-dasharray="5,5"/>
                        
                        <!-- Diameter measurement line -->
                        <g v-if="coin.diameter_mm > 0">
                            <line x1="110" y1="160" x2="290" y2="160" stroke="currentColor" stroke-width="2" class="text-secondary" marker-start="url(#arrowStart)" marker-end="url(#arrowEnd)"/>
                            <line x1="110" y1="155" x2="110" y2="165" stroke="currentColor" stroke-width="2" class="text-secondary"/>
                            <line x1="290" y1="155" x2="290" y2="165" stroke="currentColor" stroke-width="2" class="text-secondary"/>
                            
                            <!-- Diameter label -->
                            <rect x="170" y="140" width="60" height="24" rx="4" fill="currentColor" class="text-secondary" opacity="0.9"/>
                            <text x="200" y="156" text-anchor="middle" class="fill-white font-bold text-lg">Ø{{ coin.diameter_mm }}mm</text>
                        </g>
                        
                        <!-- Weight indicator -->
                        <g v-if="coin.weight_g > 0" transform="translate(60, 40)">
                            <path d="M20,30 L10,10 L30,10 Z M10,10 L30,10 L30,12 L10,12 Z" fill="currentColor" class="text-info" opacity="0.8"/>
                            <rect x="8" y="28" width="24" height="4" rx="2" fill="currentColor" class="text-info" opacity="0.8"/>
                            <rect x="0" y="35" width="40" height="20" rx="4" fill="currentColor" class="text-info" opacity="0.9"/>
                            <text x="20" y="49" text-anchor="middle" class="fill-white font-bold text-sm">{{ coin.weight_g }}g</text>
                        </g>
                        
                        <!-- Material badge -->
                        <g transform="translate(200, 270)" class="cursor-help">
                            <title>{{ coin.material }}</title>
                            <rect x="-80" y="0" width="160" height="32" rx="16" fill="currentColor" class="text-accent" opacity="0.9"/>
                            <text x="0" y="21" text-anchor="middle" class="fill-white font-bold text-lg">{{ coin.material ? coin.material.split('(')[0].trim() : 'N/A' }}</text>
                        </g>
                        
                        <!-- Thickness indicator (side view) -->
                        <g v-if="coin.thickness_mm > 0" transform="translate(320, 140)">
                            <rect x="0" y="0" width="60" height="40" rx="4" fill="currentColor" class="text-warning" opacity="0.2" stroke="currentColor" stroke-width="2"/>
                            <line x1="0" y1="0" x2="60" y2="0" stroke="currentColor" stroke-width="3" class="text-warning"/>
                            <line x1="0" y1="40" x2="60" y2="40" stroke="currentColor" stroke-width="3" class="text-warning"/>
                            <line x1="65" y1="0" x2="65" y2="40" stroke="currentColor" stroke-width="2" class="text-warning"/>
                            <line x1="63" y1="0" x2="67" y2="0" stroke="currentColor" stroke-width="2" class="text-warning"/>
                            <line x1="63" y1="40" x2="67" y2="40" stroke="currentColor" stroke-width="2" class="text-warning"/>
                            <text x="30" y="25" text-anchor="middle" class="fill-current text-base-content font-bold text-sm">{{ coin.thickness_mm }}mm</text>
                        </g>
                        
                        <defs>
                            <marker id="arrowStart" markerWidth="10" markerHeight="10" refX="5" refY="5" orient="auto"><polygon points="8,5 2,2 2,8" fill="currentColor" class="text-secondary"/></marker>
                            <marker id="arrowEnd" markerWidth="10" markerHeight="10" refX="5" refY="5" orient="auto"><polygon points="2,5 8,2 8,8" fill="currentColor" class="text-secondary"/></marker>
                        </defs>
                    </svg>
                </div>
            </div>
        </div>

        <!-- RIGHT COLUMN: Content (60%) -->
        <div class="lg:col-span-7 space-y-8">
            <!-- Header Section -->
            <div class="relative">
                <!-- Toolbar (Edit/Delete) -->
                <div class="flex justify-end gap-2 mb-2">
                     <router-link :to="`/edit/${coin.id}`" class="btn btn-ghost btn-xs text-info gap-1">
                        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-3 h-3"><path stroke-linecap="round" stroke-linejoin="round" d="M16.862 4.487l1.687-1.688a1.875 1.875 0 112.652 2.652L10.582 16.07a4.5 4.5 0 01-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 011.13-1.897l8.932-8.931zm0 0L19.5 7.125M18 14v4.75A2.25 2.25 0 0115.75 21H5.25A2.25 2.25 0 013 18.75V8.25A2.25 2.25 0 015.25 6H10" /></svg>
                        {{ $t('common.edit') }}
                    </router-link>
                    <button @click="deleteModalOpen = true" class="btn btn-ghost btn-xs text-error gap-1">
                        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-3 h-3"><path stroke-linecap="round" stroke-linejoin="round" d="M14.74 9l-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 01-2.244 2.077H8.084a2.25 2.25 0 01-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 00-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 013.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 00-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 00-7.5 0" /></svg>
                        {{ $t('common.delete') }}
                    </button>
                </div>

                <h2 v-if="coin.name" class="text-sm font-bold text-primary tracking-widest uppercase mb-1 opacity-80">{{ coin.country }}</h2>
                <h1 class="text-4xl lg:text-5xl font-extrabold tracking-tight mb-4 flex flex-wrap items-baseline gap-3">
                    {{ coin.face_value }} {{ coin.currency }}
                    <span class="text-2xl text-base-content/50 font-light" v-if="coin.year && coin.year !== 0"> '{{ coin.year }}</span>
                </h1>

                <!-- Tags Row -->
                <div class="flex flex-wrap gap-2 mb-6">
                    <div class="badge badge-lg badge-outline gap-1 pl-1 pr-3" v-if="coin.year && coin.year !== 0">
                       <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-3 h-3"><path stroke-linecap="round" stroke-linejoin="round" d="M6.75 3v2.25M17.25 3v2.25M3 18.75V7.5a2.25 2.25 0 012.25-2.25h13.5A2.25 2.25 0 0121 7.5v11.25m-18 0h18M5.25 12h13.5h-13.5zm0 5.25h13.5h-13.5z" /></svg>
                       {{ coin.year }}
                    </div>
                    <div class="badge badge-lg badge-outline gap-1 pl-1 pr-3">
                        {{ coin.currency }}
                    </div>
                    <div class="badge badge-lg badge-accent text-accent-content font-bold shadow-glow" v-if="coin.grade">
                        {{ coin.grade }}
                    </div>
                    <!-- Group Link -->
                    <router-link :to="`/list?group_id=${coin.group_id}`" class="badge badge-lg badge-primary badge-outline gap-1 hover:bg-primary hover:text-primary-content transition-colors cursor-pointer" v-if="coin.group_id && getGroupName(coin.group_id)">
                        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-3 h-3"><path stroke-linecap="round" stroke-linejoin="round" d="M2.25 12.75V12A2.25 2.25 0 014.5 9.75h15A2.25 2.25 0 0121.75 12v.75m-8.69-6.44l-2.12-2.12a1.5 1.5 0 00-1.061-.44H4.5A2.25 2.25 0 002.25 6v12a2.25 2.25 0 002.25 2.25h15A2.25 2.25 0 0021.75 18V9a2.25 2.25 0 00-2.25-2.25h-5.379a1.5 1.5 0 01-1.06-.44z" /></svg>
                        {{ getGroupName(coin.group_id) }}
                    </router-link>
                     <!-- KM Code -->
                    <div class="badge badge-lg badge-ghost gap-1 font-mono opacity-70" v-if="coin.km_code">
                        {{ coin.km_code }}
                    </div>
                </div>

                <!-- Appraisal Widget (Compact) -->
                <div v-if="coin.min_value > 0 || coin.max_value > 0" class="flex items-center p-4 bg-emerald-900 dark:bg-emerald-900 border border-emerald-500/30 rounded-xl relative overflow-hidden text-white shadow-lg">
                    <div class="absolute inset-0 bg-emerald-500/10 backdrop-blur-sm"></div>
                    <div class="relative flex items-center gap-4 w-full">
                         <div class="p-3 bg-emerald-500/20 rounded-full text-emerald-400">
                             <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6"><path stroke-linecap="round" stroke-linejoin="round" d="M12 6v12m-3-2.818l.879.659c1.171.879 3.07.879 4.242 0 1.172-.879 1.172-2.303 0-3.182C13.536 12.219 12.768 12 12 12c-.725 0-1.45-.22-2.003-.659-1.106-.879-1.106-2.303 0-3.182s2.9-.879 4.006 0l.415.33M21 12a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
                         </div>
                         <div class="flex-1">
                             <div class="text-xs font-bold text-emerald-400 uppercase tracking-wider mb-1">{{ $t('details.labels.est_value') }}</div>
                             <div class="text-2xl font-mono font-bold text-white tracking-tight">
                                 {{ getAppraisalText() }}
                             </div>
                         </div>
                    </div>
                </div>
            </div>

            <!-- Attribute Grid (3x2) -->
            <div class="grid grid-cols-2 sm:grid-cols-3 gap-4">
                 <!-- Material -->
                 <div class="p-3 bg-base-200/50 rounded-lg border border-base-300">
                    <div class="flex items-center gap-2 mb-1 opacity-60">
                         <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-4 h-4"><path stroke-linecap="round" stroke-linejoin="round" d="M21 7.5l-9-5.25L3 7.5m18 0l-9 5.25m9-5.25v9l-9 5.25M3 7.5l9 5.25M3 7.5v9l9 5.25m0-9v9" /></svg>
                         <span class="text-xs font-bold uppercase">{{ $t('details.labels.material') }}</span>
                    </div>
                    <div class="flex flex-col">
                         <div class="font-semibold text-sm leading-tight" :title="coin.material">
                             {{ getMaterialName(coin.material || coin.numista_details?.composition?.text) }}
                         </div>
                         <div v-if="getMaterialComposition(coin.material || coin.numista_details?.composition?.text)" class="text-xs opacity-70 leading-tight mt-0.5">
                             {{ getMaterialComposition(coin.material || coin.numista_details?.composition?.text) }}
                         </div>
                    </div>
                 </div>

                 <!-- Weight -->
                 <div class="p-3 bg-base-200/50 rounded-lg border border-base-300">
                    <div class="flex items-center gap-2 mb-1 opacity-60">
                         <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-4 h-4"><path stroke-linecap="round" stroke-linejoin="round" d="M12 3v18m9-9H3" /></svg> <!-- Placeholder icon -->
                         <span class="text-xs font-bold uppercase">{{ $t('details.labels.weight') }}</span>
                    </div>
                    <div class="font-semibold text-sm truncate">{{ coin.weight_g ? coin.weight_g + ' g' : (coin.numista_details?.weight ? coin.numista_details.weight + ' g' : 'N/A') }}</div>
                 </div>

                 <!-- Diameter -->
                 <div class="p-3 bg-base-200/50 rounded-lg border border-base-300">
                    <div class="flex items-center gap-2 mb-1 opacity-60">
                         <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-4 h-4"><path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-5.197-5.197m0 0A7.5 7.5 0 105.196 5.196a7.5 7.5 0 0010.607 10.607z" /></svg>
                         <span class="text-xs font-bold uppercase">{{ $t('details.labels.diameter') }}</span>
                    </div>
                    <div class="font-semibold text-sm truncate">{{ coin.diameter_mm ? coin.diameter_mm + ' mm' : (coin.numista_details?.size ? coin.numista_details.size + ' mm' : 'N/A') }}</div>
                 </div>

                 <!-- Mint -->
                 <div class="p-3 bg-base-200/50 rounded-lg border border-base-300">
                    <div class="flex items-center gap-2 mb-1 opacity-60">
                         <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-4 h-4"><path stroke-linecap="round" stroke-linejoin="round" d="M13.5 16.875h3.375m0 0h3.375m-3.375 0V13.5m0 3.375v3.375M6 10.5h2.25a2.25 2.25 0 002.25-2.25V6a2.25 2.25 0 00-2.25-2.25H6A2.25 2.25 0 003.75 6v2.25A2.25 2.25 0 006 10.5zm0 9.75h2.25A2.25 2.25 0 0010.5 18v-2.25a2.25 2.25 0 00-2.25-2.25H6a2.25 2.25 0 00-2.25 2.25V18A2.25 2.25 0 006 20.25zm9.75-9.75H18a2.25 2.25 0 002.25-2.25V6A2.25 2.25 0 0018 3.75h-2.25A2.25 2.25 0 0013.5 6v2.25a2.25 2.25 0 002.25 2.25z" /></svg>
                         <span class="text-xs font-bold uppercase">{{ $t('details.labels.mint') }}</span>
                    </div>
                    <div class="font-semibold text-sm break-words leading-tight" :title="coin.mint">{{ coin.mint || 'N/A' }}</div>
                 </div>

                  <!-- Mintage -->
                 <div class="p-3 bg-base-200/50 rounded-lg border border-base-300">
                    <div class="flex items-center gap-2 mb-1 opacity-60">
                         <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-4 h-4"><path stroke-linecap="round" stroke-linejoin="round" d="M3.75 3v11.25A2.25 2.25 0 006 16.5h2.25M3.75 3h-1.5m1.5 0h16.5m0 0h1.5m-1.5 0v11.25A2.25 2.25 0 0118 16.5h-2.25m-7.5 0h7.5m-7.5 0l-1 3m8.5-3l1 3m0 0l.5 1.5m-.5-1.5h-9.5m0 0l-.5 1.5M9 11.25v1.5M12 9v3.75m3-6v6" /></svg>
                         <span class="text-xs font-bold uppercase">{{ $t('details.labels.mintage') }}</span>
                    </div>
                    <div class="font-semibold text-sm">{{ coin.mintage ? formatMintage(coin.mintage) : 'N/A' }}</div>
                 </div>
            </div>

            <!-- Tabs Navigation -->
            <div class="overflow-x-auto -mx-4 px-4 lg:mx-0 lg:px-0 scrollbar-hide">
                <div role="tablist" class="tabs tabs-lifted tabs-md lg:tabs-lg min-w-max mx-auto lg:mx-0">
                    <a role="tab" class="tab" :class="{ 'tab-active font-bold': activeTab === 'overview' }" @click="activeTab = 'overview'">Resumen</a>
                    <a role="tab" class="tab" :class="{ 'tab-active font-bold': activeTab === 'technical' }" @click="activeTab = 'technical'">Numista</a>
                    <a role="tab" class="tab" :class="{ 'tab-active font-bold': activeTab === 'links' }" @click="activeTab = 'links'">{{ $t('details.links.title') || 'Enlaces' }}</a>
                    <a role="tab" class="tab" :class="{ 'tab-active font-bold': activeTab === 'gallery' }" @click="activeTab = 'gallery'">{{ $t('details.tabs.gallery') }}</a>
                    <a role="tab" class="tab" :class="{ 'tab-active font-bold': activeTab === 'stats' }" @click="activeTab = 'stats'">Estadísticas</a>
                    <a role="tab" class="tab" :class="{ 'tab-active font-bold': activeTab === 'notes' }" @click="activeTab = 'notes'">Notas</a>
                </div>
            </div>

            <!-- Tab Content Area -->
            <div class="bg-base-100 p-6 rounded-b-box rounded-tr-box border border-base-300 min-h-[300px]">
                
                <!-- TAB 1: OVERVIEW -->
                <div v-if="activeTab === 'overview'" class="space-y-6 animate-in fade-in duration-300">
                    <!-- Main Description -->
                    <div v-if="coin.description">
                         <h3 class="font-bold text-lg mb-2 flex items-center gap-2">
                            <span class="badge badge-neutral gap-1 text-xs"><svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="w-3 h-3"><path fill-rule="evenodd" d="M10 2a1 1 0 011 1v1.323l3.954 1.582 1.599-.8a1 1 0 01.894 1.79l-1.233.616 1.738 5.42a1 1 0 01-.285 1.05A3.989 3.989 0 0115 15a3.989 3.989 0 01-2.667-1.019 1 1 0 01-.633-.73l-3.111-9.64a1.184 1.184 0 01-.722-1.463A1.184 1.184 0 0110 2zM4.08 6.647a1 1 0 011.666-.086l1.248 1.94 1.258-.636a1 1 0 01.894 1.79l-1.233.616 1.738 5.42a1 1 0 01-.285 1.05A3.989 3.989 0 018 17a3.989 3.989 0 01-2.667-1.019 1 1 0 01-.633-.73L1.589 5.611a1 1 0 012.491-1.036z" clip-rule="evenodd" /></svg> Gemini</span>
                            {{ $t('details.sections.description') }}
                        </h3>
                        <p class="whitespace-pre-line text-sm leading-relaxed opacity-80">{{ coin.description }}</p>
                    </div>

                    <!-- Numista Links -->
                    <!-- AI Info & Actions (Moved from Notes) -->
                    <div v-if="coin.gemini_model" class="mt-6 pt-6 border-t border-base-200 text-xs text-gray-400">
                         <div class="flex flex-col sm:flex-row items-center justify-between gap-2">
                           <div class="flex gap-2 items-center">
                               <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-4 h-4">
                                  <path stroke-linecap="round" stroke-linejoin="round" d="M9.813 15.904L9 18.75l-.813-2.846a4.5 4.5 0 00-3.09-3.09L2.25 12l2.846-.813a4.5 4.5 0 003.09-3.09L9 5.25l.813 2.846a4.5 4.5 0 003.09 3.09L15.75 12l-2.846.813a4.5 4.5 0 00-3.09 3.09zM18.259 8.715L18 9.75l-.259-1.035a3.375 3.375 0 00-2.455-2.456L14.25 6l1.036-.259a3.375 3.375 0 002.455-2.456L18 2.25l.259 1.035a3.375 3.375 0 002.456 2.456L21.75 6l-1.035.259a3.375 3.375 0 00-2.456 2.456zM16.894 20.567L16.5 21.75l-.394-1.183a2.25 2.25 0 00-1.423-1.423L13.5 18.75l1.183-.394a2.25 2.25 0 001.423-1.423l.394-1.183.394 1.183a2.25 2.25 0 001.423 1.423l1.183.394-1.183.394a2.25 2.25 0 00-1.423 1.423z" />
                               </svg>
                               <span>AI Generated by {{ coin.gemini_model }} (Temp: {{ coin.gemini_temperature }})</span>
                           </div>
                           <button @click="openReprocessModal" class="btn btn-xs btn-outline btn-primary gap-1">
                               {{ $t('common.reprocess') || 'Reprocesar' }}
                           </button>
                       </div>
                    </div>
                </div>

                <!-- TAB: GALLERY -->
                <div v-if="activeTab === 'gallery'" class="space-y-6 animate-in fade-in duration-300">
                     <div class="flex justify-between items-center">
                        <h3 class="font-bold text-lg">{{ $t('details.gallery.title') }}</h3>
                        <div>
                            <input type="file" ref="galleryInput" class="hidden" @change="uploadGalleryImage" accept="image/*">
                            <button class="btn btn-primary btn-sm gap-2" @click="$refs.galleryInput.click()" :disabled="uploadingGallery">
                                <span v-if="uploadingGallery" class="loading loading-spinner loading-xs"></span>
                                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-4 h-4"><path stroke-linecap="round" stroke-linejoin="round" d="M3 16.5v2.25A2.25 2.25 0 005.25 21h13.5A2.25 2.25 0 0021 18.75V16.5m-13.5-9L12 3m0 0l4.5 4.5M12 3v13.5" /></svg>
                                {{ $t('details.gallery.upload') }}
                            </button>
                        </div>
                    </div>

                    <div v-if="(!coin.gallery_images || coin.gallery_images.length === 0) && groupImages.length === 0" class="text-center py-16 bg-base-200/50 rounded-xl border border-dashed border-base-300 flex flex-col items-center justify-center">
                         <div class="mb-4 text-base-content/20 w-16 h-16 flex items-center justify-center bg-base-200 rounded-full">
                            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-8 h-8"><path stroke-linecap="round" stroke-linejoin="round" d="M2.25 15.75l5.159-5.159a2.25 2.25 0 013.182 0l5.159 5.159m-1.5-1.5l1.409-1.409a2.25 2.25 0 013.182 0l2.909 2.909m-18 3.75h16.5a1.5 1.5 0 001.5-1.5V6a1.5 1.5 0 00-1.5-1.5H3.75A1.5 1.5 0 002.25 6v12a1.5 1.5 0 001.5 1.5zm10.5-11.25h.008v.008h-.008V8.25zm.375 0a.375.375 0 11-.75 0 .375.375 0 01.75 0z" /></svg>
                         </div>
                         <h3 class="font-bold text-lg opacity-60">{{ $t('details.gallery.empty_title') }}</h3>
                         <p class="text-sm opacity-50 max-w-xs mx-auto mb-6">{{ $t('details.gallery.empty_subtitle') }}</p>
                         <button class="btn btn-outline btn-sm gap-2" @click="$refs.galleryInput.click()" :disabled="uploadingGallery">
                            {{ $t('details.gallery.upload') }}
                         </button>
                    </div>

                    <div v-if="coin.gallery_images && coin.gallery_images.length > 0" class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-4">
                        <div v-for="img in coin.gallery_images" :key="img.id" class="relative group aspect-square bg-base-200 rounded-xl overflow-hidden shadow-sm border border-base-300">
                            <img :src="getImageUrl(img.path)" class="w-full h-full object-cover cursor-zoom-in transition-transform group-hover:scale-105" @click="openViewerForPath(img.path)">
                            <div class="absolute inset-0 bg-black/40 opacity-0 group-hover:opacity-100 transition-opacity flex items-start justify-end p-2 pointer-events-none">
                                <button @click.stop="deleteGalleryImage(img.id)" class="btn btn-xs btn-circle btn-error pointer-events-auto" :class="{'loading': deletingImageId === img.id}">
                                    <svg v-if="deletingImageId !== img.id" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-3 h-3"><path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" /></svg>
                                </button>
                            </div>
                        </div>
                    </div>

                    <!-- Group Images Section -->
                    <div v-if="groupImages.length > 0" class="pt-6 border-t border-base-200">
                        <h4 class="font-bold text-sm mb-3 opacity-70 flex items-center gap-2 uppercase tracking-wide">
                            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-4 h-4"><path stroke-linecap="round" stroke-linejoin="round" d="M2.25 12.75V12A2.25 2.25 0 014.5 9.75h15A2.25 2.25 0 0121.75 12v.75m-8.69-6.44l-2.12-2.12a1.5 1.5 0 00-1.061-.44H4.5A2.25 2.25 0 002.25 6v12a2.25 2.25 0 002.25 2.25h15A2.25 2.25 0 0021.75 18V9a2.25 2.25 0 00-2.25-2.25h-5.379a1.5 1.5 0 01-1.06-.44z" /></svg>
                            Imágenes del Grupo ({{ getGroupName(coin.group_id) }})
                        </h4>
                        <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-4">
                            <div v-for="img in groupImages" :key="img.id" class="relative group aspect-square bg-base-200 rounded-xl overflow-hidden shadow-sm border border-base-300">
                                <img :src="getImageUrl(img.path)" class="w-full h-full object-cover cursor-zoom-in transition-transform group-hover:scale-105" @click="openViewerForPath(img.path)">
                            </div>
                        </div>
                    </div>
                </div>

                <!-- TAB: STATISTICS -->
                <div v-if="activeTab === 'stats'" class="space-y-8 animate-in fade-in duration-300">
                    <div v-if="!coinStats" class="flex justify-center py-12">
                         <span class="loading loading-spinner text-primary"></span>
                    </div>
                    <div v-else class="space-y-8">
                         <!-- Radar Chart: Percentiles -->
                         <div class="card bg-base-100 shadow-sm border border-base-200">
                             <div class="card-body">
                                 <h3 class="card-title text-sm uppercase opacity-70 mb-4">Comparativa (Percentiles)</h3>
                                 <div class="h-64 relative flex justify-center">
                                     <Radar :data="percentileChartData" :options="radarOptions" />
                                 </div>
                                 <p class="text-xs text-center opacity-50 mt-2">Comparación con el resto de su colección.</p>
                             </div>
                         </div>

                         <!-- Year Distribution -->
                         <div class="card bg-base-100 shadow-sm border border-base-200">
                             <div class="card-body">
                                 <h3 class="card-title text-sm uppercase opacity-70 mb-4">Distribución por Año</h3>
                                 <div class="h-64 relative">
                                     <Bar :data="yearDistributionData" :options="barOptions" />
                                 </div>
                             </div>
                         </div>

                         <!-- Grade Distribution -->
                         <div class="card bg-base-100 shadow-sm border border-base-200">
                             <div class="card-body">
                                 <h3 class="card-title text-sm uppercase opacity-70 mb-4">Distribución por Estado</h3>
                                 <div class="h-64 relative">
                                     <Bar :data="gradeDistributionData" :options="barOptions" />
                                 </div>
                             </div>
                         </div>
                    </div>
                </div>

                <!-- TAB: LINKS (New) -->
                <div v-if="activeTab === 'links'" class="space-y-6 animate-in fade-in duration-300">
                    <div v-if="links.length === 0" class="text-center py-12 bg-base-200/50 rounded-xl border border-dashed border-base-300">
                        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-12 h-12 mx-auto text-base-content/30 mb-2"><path stroke-linecap="round" stroke-linejoin="round" d="M13.19 8.688a4.5 4.5 0 011.242 7.244l-4.5 4.5a4.5 4.5 0 01-6.364-6.364l1.757-1.757m13.35-.622l1.757-1.757a4.5 4.5 0 00-6.364-6.364l-4.5 4.5a4.5 4.5 0 001.242 7.244" /></svg>
                        <h3 class="font-bold text-lg opacity-60">{{ $t('details.links.empty') || 'No hay enlaces' }}</h3>
                        <p class="text-sm opacity-50 mb-4">{{ $t('details.links.empty_subtitle') }}</p>
                        <button class="btn btn-primary btn-sm gap-2" @click="openAddLinkModal">
                            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-4 h-4"><path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15" /></svg>
                            {{ $t('details.links.add') }}
                        </button>
                    </div>

                    <div v-else class="space-y-4">
                        <div class="flex justify-end">
                             <button class="btn btn-primary btn-xs gap-2" @click="openAddLinkModal">
                                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-3 h-3"><path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15" /></svg>
                                {{ $t('details.links.add') }}
                            </button>
                        </div>
                        
                        <div class="grid grid-cols-1 gap-4">
                            <a v-for="link in links" :key="link.id" :href="link.url" target="_blank" class="group relative block bg-base-100 hover:bg-base-200 border border-base-300 rounded-xl overflow-hidden shadow-sm hover:shadow-md transition-all">
                                <div class="flex h-full">
                                    <div v-if="link.og_image" class="w-32 min-w-[8rem] bg-gray-100 dark:bg-gray-800 bg-cover bg-center" :style="`background-image: url('${link.og_image}')`"></div>
                                    <div v-else class="w-32 min-w-[8rem] bg-gray-100 dark:bg-gray-800 flex items-center justify-center text-gray-400">
                                        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-10 h-10"><path stroke-linecap="round" stroke-linejoin="round" d="M13.19 8.688a4.5 4.5 0 011.242 7.244l-4.5 4.5a4.5 4.5 0 01-6.364-6.364l1.757-1.757m13.35-.622l1.757-1.757a4.5 4.5 0 00-6.364-6.364l-4.5 4.5a4.5 4.5 0 001.242 7.244" /></svg>
                                    </div>
                                    <div class="p-4 flex-1 min-w-0 flex flex-col justify-center">
                                        <div class="flex items-center gap-2 mb-1">
                                            <span class="badge badge-sm badge-accent badge-outline font-mono text-[10px] uppercase tracking-wider">
                                                {{ getDomain(link.url) }}
                                            </span>
                                        </div>
                                        <h4 class="font-bold text-base truncate pr-8 mb-1 leading-tight">{{ link.name || link.og_title || link.url }}</h4>
                                        <p v-if="link.og_description" class="text-xs opacity-70 line-clamp-2 leading-relaxed">{{ link.og_description }}</p>
                                        <div class="text-[10px] text-gray-400 mt-2 truncate font-mono opacity-40">{{ link.url }}</div>
                                    </div>
                                </div>
                                <div class="absolute top-2 right-2 flex gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
                                    <button @click.prevent="reloadLink(link)" class="btn btn-xs btn-circle btn-ghost text-info bg-base-100/80 backdrop-blur-sm" :class="{'loading': link.reloading}" title="Recargar metadatos">
                                        <svg v-if="!link.reloading" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-4 h-4"><path stroke-linecap="round" stroke-linejoin="round" d="M16.023 9.348h4.992v-.001M2.985 19.644v-4.992m0 0h4.992m-4.993 0l3.181 3.183a8.25 8.25 0 0013.803-3.7M4.031 9.865a8.25 8.25 0 0113.803-3.7l3.181 3.182m0-4.991v4.99" /></svg>
                                    </button>
                                    <button @click.prevent="openDeleteLinkModal(link)" class="btn btn-xs btn-circle btn-ghost text-error bg-base-100/80 backdrop-blur-sm" title="Eliminar enlace">
                                        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-4 h-4"><path stroke-linecap="round" stroke-linejoin="round" d="M14.74 9l-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 01-2.244 2.077H8.084a2.25 2.25 0 01-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 00-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 013.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 00-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 00-7.5 0" /></svg>
                                    </button>
                                </div>
                            </a>
                        </div>
                    </div>
                </div>

                <!-- TAB 2: TECHNICAL & NUMISTA -->

                <div v-if="activeTab === 'technical'" class="space-y-8 animate-in fade-in duration-300">
                    

                    

                    <div v-if="!coin.numista_number" class="text-center py-12 bg-base-200/50 rounded-xl border border-dashed border-base-300">
                         <div class="mb-4 text-primary bg-primary/10 w-16 h-16 rounded-full flex items-center justify-center mx-auto">
                            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-8 h-8"><path stroke-linecap="round" stroke-linejoin="round" d="M13.19 8.688a4.5 4.5 0 011.242 7.244l-4.5 4.5a4.5 4.5 0 01-6.364-6.364l1.757-1.757m13.35-.622l1.757-1.757a4.5 4.5 0 00-6.364-6.364l-4.5 4.5a4.5 4.5 0 001.242 7.244" /></svg>
                         </div>
                         <h3 class="font-bold text-lg opacity-80 mb-2">No vinculado con Numista</h3>
                         <p class="text-sm opacity-50 mb-6 max-w-xs mx-auto">Vincula esta moneda para obtener detalles técnicos, rareza y valores de catálogo.</p>
                         
                         <div class="flex justify-center gap-2">
                            <button @click="numistaModalOpen = true" class="btn btn-primary gap-2">
                                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-4 h-4"><path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-5.197-5.197m0 0A7.5 7.5 0 105.196 5.196a7.5 7.5 0 0010.607 10.607z" /></svg>
                                Buscar en Numista
                            </button>
                         </div>
                    </div>

                    <div v-else>
                        <!-- Numismatic Info (Ruler, Series, etc) -->
                        <div v-if="coin.ruler || coin.series || coin.commemorated_topic || coin.orientation || coin.edge" class="grid grid-cols-1 sm:grid-cols-2 gap-4">
                            <div v-if="coin.ruler">
                                <span class="font-bold block text-sm text-gray-500">{{ $t('details.labels.ruler') }}</span>
                                <span>{{ coin.ruler }}</span>
                            </div>
                            <div v-if="coin.series">
                                <span class="font-bold block text-sm text-gray-500">{{ $t('details.labels.series') }}</span>
                                <span>{{ coin.series }}</span>
                            </div>
                            <div v-if="coin.commemorated_topic">
                                <span class="font-bold block text-sm text-gray-500">{{ $t('details.labels.commemorated') }}</span>
                                <span>{{ coin.commemorated_topic }}</span>
                            </div>
                            <div v-if="coin.orientation">
                                <span class="font-bold block text-sm text-gray-500">{{ $t('details.labels.orientation') }}</span>
                                <span>{{ coin.orientation }}</span>
                            </div>
                            <div v-if="coin.edge">
                                <span class="font-bold block text-sm text-gray-500">{{ $t('details.labels.edge') }}</span>
                                <span>{{ coin.edge }}</span>
                            </div>
                        </div>

                        <div v-if="coin.numista_details">

                            
                            <!-- Obverse -->
                            <div v-if="coin.numista_details.obverse" class="mb-4">
                                <h3 class="font-bold border-b border-gray-200 dark:border-gray-700 pb-1 mb-2 text-primary">Anverso</h3>
                                <p v-if="coin.numista_details.obverse.lettering" class="mb-1 text-sm">
                                    <span class="font-semibold italic opacity-80">Leyenda:</span> 
                                    <span class="ml-1">{{ coin.numista_details.obverse.lettering }}</span>
                                </p>
                                <p v-if="coin.numista_details.obverse.description" class="text-sm">
                                    <span class="font-semibold italic opacity-80">Descripción:</span>
                                    <span class="ml-1">{{ coin.numista_details.obverse.description }}</span>
                                </p>
                            </div>

                            <!-- Reverse -->
                            <div v-if="coin.numista_details.reverse" class="mb-4">
                                <h3 class="font-bold border-b border-gray-200 dark:border-gray-700 pb-1 mb-2 text-primary">Reverso</h3>
                                <p v-if="coin.numista_details.reverse.lettering" class="mb-1 text-sm">
                                    <span class="font-semibold italic opacity-80">Leyenda:</span>
                                    <span class="ml-1">{{ coin.numista_details.reverse.lettering }}</span>
                                </p>
                                <p v-if="coin.numista_details.reverse.description" class="text-sm">
                                    <span class="font-semibold italic opacity-80">Descripción:</span>
                                    <span class="ml-1">{{ coin.numista_details.reverse.description }}</span>
                                </p>
                            </div>

                            <!-- Edge Detailed -->
                            <div v-if="coin.numista_details.edge && (coin.numista_details.edge.description || coin.numista_details.edge.lettering)" class="mb-4">
                                <h3 class="font-bold border-b border-gray-200 dark:border-gray-700 pb-1 mb-2 text-primary">Canto</h3>
                                <p v-if="coin.numista_details.edge.description" class="mb-1 text-sm">{{ coin.numista_details.edge.description }}</p>
                                <p v-if="coin.numista_details.edge.lettering" class="text-sm">
                                    <span class="font-semibold italic opacity-80">Leyenda:</span>
                                    <span class="ml-1">{{ coin.numista_details.edge.lettering }}</span>
                                </p>
                            </div>

                            <!-- Technique -->
                            <div v-if="coin.numista_details.technique && coin.numista_details.technique.text" class="mb-4">
                                <h3 class="font-bold border-b border-gray-200 dark:border-gray-700 pb-1 mb-2 text-primary">Técnica</h3>
                                <p class="text-sm">{{ coin.numista_details.technique.text }}</p>
                            </div>
                        </div>


                        <!-- Numista Actions (Moved to bottom) -->
                        <div class="flex flex-wrap items-center gap-4 mt-6 pt-6 border-t border-base-200">
                            <div class="flex gap-2">
                                <a v-if="coin.numista_number" :href="getNumistaUrl()" target="_blank" class="btn btn-sm btn-outline btn-info gap-2">
                                    Ver en Numista
                                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-3 h-3"><path stroke-linecap="round" stroke-linejoin="round" d="M13.5 6H5.25A2.25 2.25 0 003 8.25v10.5A2.25 2.25 0 005.25 21h10.5A2.25 2.25 0 0018 18.75V10.5m-10.5 6L21 3m0 0h-5.25M21 3v5.25" /></svg>
                                </a>
                                <button v-if="numistaCount >= 0" @click="numistaModalOpen = true" class="btn btn-sm btn-ghost gap-1 text-xs">
                                    {{ numistaCount }} resultados alternativos
                                </button>
                            </div>
                            <!-- Sync Button renamed and styled -->
                            <button @click="syncNumista" class="btn btn-xs btn-outline btn-primary gap-1 ml-auto" :disabled="syncing">
                                <span v-if="syncing" class="loading loading-spinner loading-xs"></span>
                                {{ $t('common.reprocess') || 'Reprocesar' }}
                            </button>
                        </div>
                    </div>
                </div>

                <!-- TAB 3: NOTES & HISTORY -->
                <div v-if="activeTab === 'notes'" class="space-y-6 animate-in fade-in duration-300">
                    
                    <!-- Transaction History -->
                     <div v-if="coin.acquired_at || coin.sold_at">
                        <h3 class="font-bold text-lg mb-3">{{ $t('details.sections.transaction_history') }}</h3>
                        <div class="space-y-3">
                            <div v-if="coin.acquired_at" class="flex items-center gap-3 p-3 bg-success/10 rounded-lg">
                                <div class="badge badge-success gap-2">
                                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-4 h-4"><path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15" /></svg>
                                    {{ $t('details.labels.acquired') }}
                                </div>
                                <div class="flex-1">
                                    <div class="font-bold">{{ new Date(coin.acquired_at).toLocaleDateString() }}</div>
                                    <div v-if="coin.price_paid > 0" class="text-sm opacity-70">{{ $t('details.labels.price_paid') }}: {{ formatCurrency(coin.price_paid) }}</div>
                                </div>
                            </div>
                            
                            <div v-if="coin.sold_at" class="flex items-center gap-3 p-3 bg-warning/10 rounded-lg">
                                <div class="badge badge-warning gap-2">
                                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-4 h-4"><path stroke-linecap="round" stroke-linejoin="round" d="M2.25 18.75a60.07 60.07 0 0115.797 2.101c.727.198 1.453-.342 1.453-1.096V18.75M3.75 4.5v.75A.75.75 0 013 6h-.75m0 0v-.375c0-.621.504-1.125 1.125-1.125H20.25M2.25 6v9m18-10.5v.75c0 .414.336.75.75.75h.75m-1.5-1.5h.375c.621 0 1.125.504 1.125 1.125v9.75c0 .621-.504 1.125-1.125 1.125h-.375m1.5-1.5H21a.75.75 0 00-.75.75v.75m0 0H3.75m0 0h-.375a1.125 1.125 0 01-1.125-1.125V15m1.5 1.5v-.75A.75.75 0 003 15h-.75M15 10.5a3 3 0 11-6 0 3 3 0 016 0zm3 0h.008v.008H18V10.5zm-12 0h.008v.008H6V10.5z" /></svg>
                                    {{ $t('details.labels.sold') }}
                                </div>
                                <div class="flex-1">
                                    <div class="font-bold">{{ new Date(coin.sold_at).toLocaleDateString() }}</div>
                                    <div v-if="coin.sold_price > 0" class="text-sm opacity-70">{{ $t('details.labels.sold_price') }}: {{ formatCurrency(coin.sold_price) }}</div>
                                    <div v-if="coin.sale_channel" class="text-sm opacity-70">{{ $t('details.labels.sale_channel') }}: {{ coin.sale_channel }}</div>
                                </div>
                            </div>

                             <div v-if="!coin.sold_at" class="alert alert-info py-2">
                                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="stroke-current shrink-0 w-6 h-6"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
                                <span>{{ $t('details.labels.in_collection') }}</span>
                            </div>
                        </div>
                    </div>

                    <!-- User Notes -->
                    <div v-if="coin.technical_notes || coin.personal_notes" class="space-y-4">
                        <div v-if="coin.technical_notes">
                            <h3 class="font-bold text-sm text-gray-500 mb-2 flex items-center gap-2">
                                <span class="badge badge-neutral gap-1 text-xs"><svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="w-3 h-3"><path fill-rule="evenodd" d="M10 2a1 1 0 011 1v1.323l3.954 1.582 1.599-.8a1 1 0 01.894 1.79l-1.233.616 1.738 5.42a1 1 0 01-.285 1.05A3.989 3.989 0 0115 15a3.989 3.989 0 01-2.667-1.019 1 1 0 01-.633-.73l-3.111-9.64a1.184 1.184 0 01-.722-1.463A1.184 1.184 0 0110 2zM4.08 6.647a1 1 0 011.666-.086l1.248 1.94 1.258-.636a1 1 0 01.894 1.79l-1.233.616 1.738 5.42a1 1 0 01-.285 1.05A3.989 3.989 0 018 17a3.989 3.989 0 01-2.667-1.019 1 1 0 01-.633-.73L1.589 5.611a1 1 0 012.491-1.036z" clip-rule="evenodd" /></svg> Gemini</span>
                                {{ $t('details.sections.technical_notes') }}
                            </h3>
                            <div class="text-sm whitespace-pre-line bg-base-200 p-3 rounded-lg">{{ coin.technical_notes }}</div>
                        </div>
                        
                        <div class="group relative">
                            <h3 class="font-bold text-sm text-gray-500 mb-2 flex items-center gap-2">
                                {{ $t('details.sections.personal_notes') }}
                                <button v-if="!editingNotes" @click="startEditingNotes" class="btn btn-ghost btn-xs btn-circle text-base-content/50 hover:text-primary transition-colors">
                                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-4 h-4"><path stroke-linecap="round" stroke-linejoin="round" d="M16.862 4.487l1.687-1.688a1.875 1.875 0 112.652 2.652L10.582 16.07a4.5 4.5 0 01-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 011.13-1.897l8.932-8.931zm0 0L19.5 7.125M18 14v4.75A2.25 2.25 0 0115.75 21H5.25A2.25 2.25 0 013 18.75V8.25A2.25 2.25 0 015.25 6H10" /></svg>
                                </button>
                            </h3>
                            
                            <div v-if="editingNotes" class="space-y-2">
                                <textarea 
                                    v-model="tempNotes" 
                                    class="textarea textarea-bordered w-full h-32 text-sm" 
                                    :placeholder="$t('form.placeholders.notes')"
                                ></textarea>
                                <div class="flex justify-end gap-2">
                                    <button @click="cancelNotes" class="btn btn-xs btn-ghost">{{ $t('common.cancel') }}</button>
                                    <button @click="saveNotes" class="btn btn-xs btn-primary" :disabled="savingNotes">
                                        <span v-if="savingNotes" class="loading loading-spinner loading-xs"></span>
                                        {{ $t('common.save') }}
                                    </button>
                                </div>
                            </div>
                            <div v-else>
                                <div v-if="coin.personal_notes" class="text-sm italic whitespace-pre-line bg-blue-50 dark:bg-blue-900/10 p-3 rounded-lg border border-blue-100 dark:border-blue-900/30">
                                    {{ coin.personal_notes }}
                                </div>
                                <div v-else @click="startEditingNotes" class="text-sm italic text-gray-500 bg-base-200/50 p-4 rounded-lg border border-dashed border-base-300 cursor-pointer hover:bg-base-200 transition-colors flex items-center justify-center gap-2">
                                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-4 h-4"><path stroke-linecap="round" stroke-linejoin="round" d="M16.862 4.487l1.687-1.688a1.875 1.875 0 112.652 2.652L10.582 16.07a4.5 4.5 0 01-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 011.13-1.897l8.932-8.931zm0 0L19.5 7.125M18 14v4.75A2.25 2.25 0 0115.75 21H5.25A2.25 2.25 0 013 18.75V8.25A2.25 2.25 0 015.25 6H10" /></svg>
                                    Añadir nota personal...
                                </div>
                            </div>
                        </div>
                    </div>



                </div>

            </div>
        </div>
    </div>

    <!-- Modals -->
    <!-- Delete Coin Modal (already exists presumably) -->
    
    <!-- Gallery Delete Modal -->
    <dialog class="modal" :class="{'modal-open': galleryDeleteId}">
        <div class="modal-box">
             <h3 class="font-bold text-lg">{{ $t('common.delete') }}</h3>
             <p class="py-4">{{ $t('common.delete_modal.confirm') }}?</p>
             <div class="modal-action">
                 <button class="btn" @click="galleryDeleteId = null">{{ $t('common.cancel') }}</button>
                 <button class="btn btn-error" @click="confirmDeleteGalleryImageAction">{{ $t('common.delete') }}</button>
             </div>
        </div>
    </dialog>
  </div>
  <div v-else class="flex justify-center p-20">
    <span class="loading loading-spinner loading-lg"></span>
  </div>

  <!-- Rotation Editor Modal -->
  <div v-if="editingSide" class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-90 backdrop-blur-sm transition-opacity" @click.self="cancelRotation">
    <div class="relative w-full max-w-2xl p-4 flex flex-col gap-4 animate-in zoom-in duration-200">
        <button @click="cancelRotation" class="absolute -top-10 right-0 btn btn-circle btn-ghost text-white z-50 hover:rotate-90 transition-transform">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" /></svg>
        </button>
        
        <!-- Editor Area -->
        <div class="relative flex justify-center items-center overflow-hidden bg-gray-900 rounded-2xl shadow-2xl h-[50vh] sm:h-[500px] border border-gray-700">
             <!-- Grid Overlay -->
             <div class="absolute inset-0 pointer-events-none z-10 opacity-20" 
                  style="background-image: linear-gradient(rgba(255,255,255,0.5) 1px, transparent 1px), linear-gradient(90deg, rgba(255,255,255,0.5) 1px, transparent 1px); background-size: 33.3% 33.3%;">
             </div>
             <!-- Center Crosshair -->
             <div class="absolute inset-0 pointer-events-none z-10 flex items-center justify-center opacity-30">
                <div class="w-8 h-8 border border-white rounded-full"></div>
                <div class="absolute w-4 h-[1px] bg-white"></div>
                <div class="absolute h-4 w-[1px] bg-white"></div>
             </div>
             
             <!-- Image -->
             <img :src="getCurrentImageUrl(editingSide)" 
                  class="max-h-full max-w-full transition-transform duration-75 ease-out object-contain select-none"
                  :style="{ transform: `rotate(${tempRotation}deg)` }" />
        </div>

        <!-- Controls -->
        <div class="bg-gray-800 rounded-2xl p-6 flex flex-col gap-6 shadow-xl border border-gray-700">
             <div class="flex justify-between text-white text-sm font-mono items-end">
                 <span class="opacity-50">-180°</span>
                 <div class="flex flex-col items-center">
                    <span class="text-xs uppercase text-gray-500 font-bold tracking-widest mb-1">{{ $t('common.rotation') }}</span>
                    <span class="font-bold text-3xl text-primary tracking-tighter">{{ tempRotation > 0 ? '+' : ''}}{{ tempRotation }}°</span>
                 </div>
                 <span class="opacity-50">+180°</span>
             </div>
             <input type="range" min="-180" max="180" step="0.5" v-model.number="tempRotation" class="range range-primary w-full" />
             <div class="w-full flex justify-between px-2 text-xs text-gray-500">
                <span>|</span><span>|</span><span>|</span><span>|</span><span>|</span>
             </div>
             
             <div class="flex justify-end gap-3 mt-2">
                 <button @click="cancelRotation" class="btn btn-ghost text-white hover:bg-white/10">{{ $t('common.cancel') }}</button>
                 <button @click="saveRotation" class="btn btn-primary px-8 gap-2">
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-5 h-5">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M4.5 12.75l6 6 9-13.5" />
                    </svg>
                    {{ $t('common.apply') }}
                 </button>
             </div>
        </div>
    </div>
  </div>

  <!-- Delete Modal -->
  <dialog v-if="deleteModalOpen" id="delete_modal" class="modal" :class="{ 'modal-open': deleteModalOpen }">
    <div class="modal-box">
      <h3 class="font-bold text-lg text-error">{{ $t('list.delete_modal.title') }}</h3>
      <p class="py-4">{{ $t('list.delete_modal.confirm') }} <span class="font-bold">{{ coin?.name || $t('common.unknown') }}</span>? {{ $t('list.delete_modal.warning') }}</p>
      <div class="modal-action">
        <button class="btn" @click="deleteModalOpen = false">{{ $t('common.cancel') }}</button>
        <button class="btn btn-error" @click="deleteCoin" :disabled="deleting">
          <span v-if="deleting" class="loading loading-spinner"></span>
          {{ $t('common.delete') }}
        </button>
      </div>
    </div>
  </dialog>

  <!-- Reprocess Modal -->
  <dialog v-if="reprocessModalOpen" id="reprocess_modal" class="modal" :class="{ 'modal-open': reprocessModalOpen }">
    <div class="modal-box">
      <h3 class="font-bold text-lg text-primary flex items-center gap-2">
          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6">
            <path stroke-linecap="round" stroke-linejoin="round" d="M16.023 9.348h4.992v-.001M2.985 19.644v-4.992m0 0h4.992m-4.993 0l3.181 3.183a8.25 8.25 0 0013.803-3.7M4.031 9.865a8.25 8.25 0 0113.803-3.7l3.181 3.182m0-4.991v4.99" />
          </svg>
          {{ $t('details.reprocess_modal.title') || 'Reprocesar con IA' }}
      </h3>
      
      <div class="py-4 space-y-4">
          <div class="alert alert-warning text-sm">
            <svg xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" /></svg>
            <span>{{ $t('details.reprocess_modal.warning') || '¡Advertencia! Esto sobreescribirá todos los detalles de la moneda con información nueva generada por la IA.' }}</span>
          </div>

          <GeminiConfig 
            v-model:isOpen="aiSettingsOpen"
            v-model:model="selectedModel"
            v-model:temperature="temperature"
            :available-models="availableModels"
            :title="$t('form.ai_settings') || 'Configuración del Modelo'"
          />
      </div>

      <div class="modal-action">
        <button class="btn" @click="reprocessModalOpen = false">{{ $t('common.cancel') }}</button>
        <button class="btn btn-primary" @click="reprocessCoin" :disabled="reprocessing">
          <span v-if="reprocessing" class="loading loading-spinner"></span>
          {{ $t('common.process') || 'Procesar' }}
        </button>
      </div>
    </div>
  </dialog>

  <!-- Add Link Modal -->
  <dialog v-if="addLinkModalOpen" id="add_link_modal" class="modal" :class="{ 'modal-open': addLinkModalOpen }">
    <div class="modal-box">
      <h3 class="font-bold text-lg text-primary flex items-center gap-2">
          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6"><path stroke-linecap="round" stroke-linejoin="round" d="M13.19 8.688a4.5 4.5 0 011.242 7.244l-4.5 4.5a4.5 4.5 0 01-6.364-6.364l1.757-1.757m13.35-.622l1.757-1.757a4.5 4.5 0 00-6.364-6.364l-4.5 4.5a4.5 4.5 0 001.242 7.244" /></svg>
          {{ $t('details.links.add') }}
      </h3>
      
      <div class="py-4 space-y-4">
          <div class="form-control w-full">
              <label class="label">
                  <span class="label-text">URL</span>
              </label>
              <input type="url" v-model="newLink.url" :placeholder="$t('details.links.url_placeholder')" class="input input-bordered w-full" :class="{ 'input-error': newLinkError }" @input="newLinkError = ''" />
              <label class="label" v-if="newLinkError">
                  <span class="label-text-alt text-error">{{ newLinkError }}</span>
              </label>
          </div>
          <!-- Optional Name? We'll rely on OG for now or add if user wants. Let's keep it simple first. -->
          <!-- Actually, let's allow optional name override -->
<!--           <div class="form-control w-full">
              <label class="label">
                  <span class="label-text">{{ $t('details.links.name_label') }}</span>
              </label>
              <input type="text" v-model="newLink.name" placeholder="Opcional" class="input input-bordered w-full" />
          </div> -->
      </div>

      <div class="modal-action">
        <button class="btn" @click="closeAddLinkModal">{{ $t('common.cancel') }}</button>
        <button class="btn btn-primary" @click="addLink" :disabled="addingLink || !newLink.url">
          <span v-if="addingLink" class="loading loading-spinner"></span>
          {{ $t('common.save') }}
        </button>
      </div>
    </div>
  </dialog>

  <!-- Delete Link Modal -->
  <dialog v-if="deleteLinkModalOpen" id="delete_link_modal" class="modal" :class="{ 'modal-open': deleteLinkModalOpen }">
    <div class="modal-box">
      <h3 class="font-bold text-lg text-error flex items-center gap-2">
          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6"><path stroke-linecap="round" stroke-linejoin="round" d="M14.74 9l-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 01-2.244 2.077H8.084a2.25 2.25 0 01-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 00-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 013.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 00-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 00-7.5 0" /></svg>
          {{ $t('details.links.delete_title') || 'Eliminar Enlace' }}
      </h3>
      <p class="py-4">
          {{ $t('details.links.delete_confirm') }}
          <br>
          <span class="font-bold opacity-70 block mt-2 text-sm">{{ linkToDelete?.name || linkToDelete?.og_title || linkToDelete?.url }}</span>
      </p>
      <div class="modal-action">
        <button class="btn" @click="deleteLinkModalOpen = false">{{ $t('common.cancel') }}</button>
        <button class="btn btn-error" @click="confirmDeleteLink" :disabled="deletingLink">
          <span v-if="deletingLink" class="loading loading-spinner"></span>
          {{ $t('common.delete') }}
        </button>
      </div>
    </div>
  </dialog>

  <!-- Numista Results Modal -->
  <dialog v-if="numistaModalOpen" id="numista_modal" class="modal" :class="{ 'modal-open': numistaModalOpen }">
    <div class="modal-box w-11/12 max-w-4xl">
      <h3 class="font-bold text-lg text-primary flex items-center gap-2 mb-4">
          Resultados de Numista ({{ numistaResults.length }})
      </h3>
      
      <div class="overflow-x-auto">
        <table class="table table-zebra w-full">
          <thead>
            <tr>
              <th>ID</th>
              <th>Título</th>
              <th>Año</th>
              <th>Emisor</th>
              <th>Acciones</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="result in numistaResults" :key="result.id">
              <td>
                  <a :href="`https://es.numista.com/catalogue/pieces${result.id}.html`" target="_blank" class="link link-primary font-mono text-xs">
                      {{ result.id }}
                  </a>
              </td>
              <td class="whitespace-normal">
                  <div class="font-bold">{{ result.title }}</div>
              </td>
              <td>{{ result.min_year }} - {{ result.max_year }}</td>
              <td>{{ result.issuer?.name }}</td>
              <td class="flex gap-2">
                  <!-- Redundant icon removed -->
                  <button @click="applyNumistaResult(result.id)" class="btn btn-xs btn-primary" :disabled="applyingNumista">
                      <span v-if="applyingNumista && selectedNumistaId === result.id" class="loading loading-spinner loading-xs"></span>
                      Aplicar
                  </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="divider">O importar manualmente</div>
      <div class="flex items-end gap-4 p-4 bg-base-200 rounded-lg">
          <div class="form-control w-full max-w-xs">
              <label class="label">
                  <span class="label-text">Numista ID</span>
              </label>
              <input type="number" v-model="manualNumistaId" placeholder="Ej: 12345" class="input input-bordered w-full max-w-xs" />
          </div>
          <button class="btn btn-primary" @click="applyManualNumista" :disabled="!manualNumistaId || applyingNumista">
              <span v-if="applyingNumista && selectedNumistaId === manualNumistaId" class="loading loading-spinner"></span>
              Enviar
          </button>
      </div>

      <div class="modal-action">
        <button class="btn" @click="numistaModalOpen = false">Cerrar</button>
      </div>
    </div>
  </dialog>

  <ImageViewer 
      :is-open="viewerOpen" 
      :image-url="viewerImage" 
      @close="viewerOpen = false" 
    />
</template>

<script setup>
import { ref, onMounted, computed, reactive, watch } from 'vue'
import axios from 'axios'
import { useRoute, useRouter } from 'vue-router'
import { normalizeGrade, GRADE_ORDER } from '../utils/grades'
import ImageViewer from '../components/ImageViewer.vue'
import GeminiConfig from '../components/GeminiConfig.vue'
import { formatMintage } from '../utils/formatters'
import { useI18n } from 'vue-i18n'
import { Chart as ChartJS } from 'chart.js/auto'
import { Bar, Radar } from 'vue-chartjs'

const { t } = useI18n()

const route = useRoute()
const router = useRouter()
const coin = ref(null)
const loading = ref(true)
const API_URL = import.meta.env.VITE_API_URL || '/api/v1'
const STORAGE_URL = ''

const viewerOpen = ref(false)
const viewerImage = ref('')
const activeImageSource = ref('processed') // processed, original
const activeTab = ref('overview') // overview, technical, notes

const updatingLink = ref(null)

const coinStats = ref(null)

const fetchStats = async () => {
    if (!coin.value) return
    try {
        const res = await axios.get(`${API_URL}/coins/${coin.value.id}/stats`)
        coinStats.value = res.data
    } catch (e) {
        console.error("Failed to fetch coin stats", e)
    }
}

// Charts Data
const percentileChartData = computed(() => {
    if (!coinStats.value) return null
    return {
        labels: ['Valor', 'Rareza', 'Peso', 'Tamaño'],
        datasets: [{
            label: 'Percentil',
            data: [
                (coinStats.value.value_percentile || 0) * 100,
                (coinStats.value.rarity_percentile || 0) * 100,
                (coinStats.value.weight_percentile || 0) * 100,
                (coinStats.value.size_percentile || 0) * 100
            ],
            backgroundColor: 'rgba(234, 179, 8, 0.2)',
            borderColor: 'rgba(234, 179, 8, 1)',
            pointBackgroundColor: 'rgba(234, 179, 8, 1)',
            pointBorderColor: '#fff',
            pointHoverBackgroundColor: '#fff',
            pointHoverBorderColor: 'rgba(234, 179, 8, 1)'
        }]
    }
})

const radarOptions = {
    responsive: true,
    maintainAspectRatio: false,
    scales: {
        r: {
            angleLines: {
                color: 'rgba(128, 128, 128, 0.2)'
            },
            grid: {
                color: 'rgba(128, 128, 128, 0.2)'
            },
            pointLabels: {
                font: {
                    size: 12
                }
            },
            suggestedMin: 0,
            suggestedMax: 100
        }
    },
    plugins: {
        legend: { display: false }
    }
}

const yearDistributionData = computed(() => {
    if (!coinStats.value || !coinStats.value.year_distribution) return null
    
    // Convert map to sorted arrays
    const years = Object.keys(coinStats.value.year_distribution).map(Number).sort((a,b) => a-b)
    // Filter to a reasonable range around the current coin if too many?
    // For now show all, or maybe aggregate by decade if too many.
    // Let's just show all for now.
    
    const relevantYears = years.filter(y => y > 0)
    const labels = relevantYears.map(String)
    const data = relevantYears.map(y => coinStats.value.year_distribution[y])
    
    const backgroundColors = relevantYears.map(y => y === coin.value.year ? 'rgba(234, 179, 8, 0.8)' : 'rgba(128, 128, 128, 0.2)')

    return {
        labels,
        datasets: [{
            label: 'Monedas',
            data,
            backgroundColor: backgroundColors,
            borderRadius: 4
        }]
    }
})

const gradeDistributionData = computed(() => {
    if (!coinStats.value || !coinStats.value.grade_distribution) return null

    // Use predefined grade order for sorting
    const labels = GRADE_ORDER.filter(g => coinStats.value.grade_distribution[g] !== undefined)
    const data = labels.map(g => coinStats.value.grade_distribution[g])
    
    const currentGrade = coin.value.grade_code || coin.value.grade // Assuming grade_code or grade matches map keys
    // Need to ensure coin.grade matches keys in GRADE_ORDER ("S/C" vs "SC" etc).
    // The backend uses normalized keys probably.
    
    const backgroundColors = labels.map(g => g === currentGrade ? 'rgba(234, 179, 8, 0.8)' : 'rgba(128, 128, 128, 0.2)')

    return {
        labels,
        datasets: [{
            label: 'Monedas',
            data,
            backgroundColor: backgroundColors,
            borderRadius: 4
        }]
    }
})

const barOptions = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
        legend: { display: false }
    },
    scales: {
        y: {
            beginAtZero: true,
            ticks: { precision: 0 }
        },
        x: {
            grid: { display: false }
        }
    }
}

// Group State
const groups = ref([])
const groupImages = ref([])

const fetchGroupImages = async () => {
    if (!coin.value || !coin.value.group_id) return
    try {
        const res = await axios.get(`${API_URL}/groups/${coin.value.group_id}/images`)
        groupImages.value = res.data
    } catch (e) {
        console.error("Failed to fetch group images", e)
    }
}

const fetchGroups = async () => {
    try {
        const res = await axios.get(`${API_URL}/groups`)
        groups.value = res.data
    } catch (e) {
        console.error("Failed to fetch groups", e)
    }
}

const getGroupName = (groupId) => {
    if (!groupId) return null
    const group = groups.value.find(g => g.id === groupId)
    return group ? group.name : null
}

// Delete Modal State
const deleteModalOpen = ref(false)
const deleting = ref(false)

// Reprocess Modal State
const reprocessModalOpen = ref(false)
const reprocessing = ref(false)
const selectedModel = ref('gemini-2.5-flash')
const temperature = ref(0.1)
const availableModels = ref([])

const aiSettingsOpen = ref(true)

// Notes Editing State
const editingNotes = ref(false)
const tempNotes = ref('')
const savingNotes = ref(false)

const startEditingNotes = () => {
    tempNotes.value = coin.value.personal_notes || ''
    editingNotes.value = true
}

const cancelNotes = () => {
    editingNotes.value = false
    tempNotes.value = ''
}

const saveNotes = async () => {
    if (!coin.value) return
    savingNotes.value = true
    try {
        const updatedCoin = { ...coin.value, personal_notes: tempNotes.value } 
        // We need to send the full object or backend compatible update
        // Assuming PUT /coins/:id expects the full body or partial. 
        // Based on EditCoin.vue, usually we send what we edit. 
        // Let's try sending the updated fields.
        
        await axios.put(`${API_URL}/coins/${coin.value.id}`, updatedCoin)
        
        coin.value.personal_notes = tempNotes.value
        editingNotes.value = false
    } catch (e) {
        console.error("Failed to save notes", e)
        alert('Failed to save notes: ' + (e.response?.data?.error || e.message))
    } finally {
        savingNotes.value = false
    }
}

// Numista Sync State
const syncing = ref(false)

const syncNumista = async () => {
    if (!coin.value) return
    syncing.value = true
    try {
        await axios.post(`${API_URL}/coins/${coin.value.id}/reprocess-numista`)
        // Refresh coin data
        const res = await axios.get(`${API_URL}/coins/${coin.value.id}`)
        coin.value = res.data
    } catch (e) {
        console.error("Failed to sync with Numista", e)
        alert('Failed to sync: ' + (e.response?.data?.error || e.message))
    } finally {
        syncing.value = false
    }
}

const openReprocessModal = async () => {
    reprocessModalOpen.value = true
    // Initialize with current values if available
    if (coin.value.gemini_model) selectedModel.value = coin.value.gemini_model
    if (coin.value.gemini_temperature) temperature.value = coin.value.gemini_temperature
    
    // Fetch models if empty
    if (availableModels.value.length === 0) {
        try {
            const res = await axios.get(`${API_URL}/gemini/models`)
            availableModels.value = res.data
        } catch (e) {
            console.error("Failed to fetch models", e)
        }
    }
}

const reprocessCoin = async () => {
    if (!coin.value) return
    reprocessing.value = true
    try {
        const res = await axios.post(`${API_URL}/coins/${coin.value.id}/analyze`, {
            model_name: selectedModel.value,
            temperature: temperature.value
        })
        coin.value = res.data
        reprocessModalOpen.value = false
        // Optional: Show success toast
    } catch (e) {
        console.error(e)
        alert('Failed to reprocess coin: ' + (e.response?.data?.error || e.message))
    } finally {
        reprocessing.value = false
    }
}

const deleteCoin = async () => {
    if (!coin.value) return
    deleting.value = true
    try {
        await axios.delete(`${API_URL}/coins/${coin.value.id}`)
        router.push('/list')
    } catch (e) {
        console.error("Failed to delete coin", e)
        alert("Failed to delete coin")
    } finally {
        deleting.value = false
        deleteModalOpen.value = false
    }
}

// Links Logic
const links = ref([])
const addLinkModalOpen = ref(false)
const addingLink = ref(false)
const newLink = reactive({ url: '', name: '' })
const newLinkError = ref('')

const deleteLinkModalOpen = ref(false)
const linkToDelete = ref(null)
const deletingLink = ref(false)

const getDomain = (url) => {
    try {
        const hostname = new URL(url).hostname
        return hostname.replace(/^www\./, '')
    } catch (e) {
        return url
    }
}

const fetchLinks = async () => {
    if (!coin.value) return
    try {
        const res = await axios.get(`${API_URL}/coins/${coin.value.id}/links`)
        links.value = res.data
    } catch (e) {
        console.error("Failed to fetch links", e)
    }
}

const openAddLinkModal = () => {
    newLink.url = ''
    newLink.name = ''
    newLinkError.value = ''
    addLinkModalOpen.value = true
}

const closeAddLinkModal = () => {
    addLinkModalOpen.value = false
}

const isValidUrl = (string) => {
  try {
    new URL(string);
    return true;
  } catch (_) {
    return false;  
  }
}

const addLink = async () => {
    if (!newLink.url) return
    if (!isValidUrl(newLink.url)) {
        newLinkError.value = 'URL inválida'
        return
    }

    addingLink.value = true
    try {
        const res = await axios.post(`${API_URL}/coins/${coin.value.id}/links`, {
            url: newLink.url
        })
        links.value.unshift(res.data) // Add to top
        closeAddLinkModal()
    } catch (e) {
        console.error("Failed to add link", e)
        newLinkError.value = e.response?.data?.error || e.message
    } finally {
        addingLink.value = false
    }
}

const openDeleteLinkModal = (link) => {
    linkToDelete.value = link
    deleteLinkModalOpen.value = true
}

const confirmDeleteLink = async () => {
    if (!linkToDelete.value) return
    deletingLink.value = true
    try {
        await axios.delete(`${API_URL}/coins/${coin.value.id}/links/${linkToDelete.value.id}`)
        links.value = links.value.filter(l => l.id !== linkToDelete.value.id)
        deleteLinkModalOpen.value = false
        linkToDelete.value = null
    } catch (e) {
         console.error("Failed to delete link", e)
         alert('Failed to delete link')
    } finally {
        deletingLink.value = false
    }
}

const getImageUrl = (path) => {
    if (!path) return '/broken_coin.png'
    if (path.includes('storage/')) {
        return `${STORAGE_URL}/storage/${path.split('storage/')[1]}`
    }
    return path
}

const getNumistaUrl = () => {
    if (coin.value && coin.value.numista_number) {
        return `https://es.numista.com/${coin.value.numista_number}`
    }
    return null
}

// Rotation Logic
const rotations = reactive({ front: 0, back: 0 })
const editingSide = ref(null)
const tempRotation = ref(0)

const imageTimestamp = ref(Date.now())

const openRotationEditor = (side) => {
    editingSide.value = side
    tempRotation.value = rotations[side]
}

const saveRotation = () => {
    if (editingSide.value && coin.value) {
        const side = editingSide.value
        const angle = tempRotation.value
        const id = coin.value.id
        
        // Optimistic update:
        // Update local state to reflect rotation visually (via CSS)
        rotations[side] = angle
        editingSide.value = null

        // Fire and forget (log error if needed)
        axios.post(`${API_URL}/coins/${id}/rotate`, {
            side: side,
            angle: angle
        }).catch(e => {
            console.error("Failed to save rotation background", e)
            // Optional: Revert via toast? For now silent fail or alert.
        })
    }
}

const cancelRotation = () => {
    editingSide.value = null
}

const getCurrentImageUrl = (side) => {
    if (!coin.value) return '/broken_coin.png'
    
    // Standard Logic (Processed/Original)
    if (coin.value.images && coin.value.images.length > 0) {
        let typeToFind = 'crop'
        if (activeImageSource.value === 'original') {
            typeToFind = 'original'
        }
        
        const img = coin.value.images.find(img => img.image_type === typeToFind && img.side === side)
        if (img) {
            return `${getImageUrl(img.path)}?t=${imageTimestamp.value}`
        }
    }
    
    // Fallback logic
    return '/broken_coin.png'
}

const handleImageError = (e) => {
    e.target.src = '/broken_coin.png'
}

const openViewer = (side) => {
    viewerImage.value = getCurrentImageUrl(side)
    viewerOpen.value = true
}

const getGradeDescription = (code) => {
    if (!code) return t('common.unknown')
    const base = normalizeGrade(code)
    return t(`grades.${base}.desc`)
}

// Gallery Logic
const galleryInput = ref(null)
const uploadingGallery = ref(false)
const deletingImageId = ref(null)

const uploadGalleryImage = async (event) => {
    const file = event.target.files[0]
    if (!file) return

    uploadingGallery.value = true
    const formData = new FormData()
    formData.append('image', file)

    try {
        await axios.post(`${API_URL}/coins/${coin.value.id}/gallery`, formData, {
            headers: { 'Content-Type': 'multipart/form-data' }
        })
        // Refresh coin
        const res = await axios.get(`${API_URL}/coins/${coin.value.id}`)
        coin.value = res.data
    } catch (e) {
        console.error("Failed to upload gallery image", e)
        alert(t('form.errors.upload_failed'))
    } finally {
        uploadingGallery.value = false
        // Reset input
        event.target.value = ''
    }
}

const galleryDeleteId = ref(null)

const confirmDeleteGalleryImage = (id) => {
    galleryDeleteId.value = id
}

const confirmDeleteGalleryImageAction = async () => {
    if (!galleryDeleteId.value) return
    const imageId = galleryDeleteId.value
    galleryDeleteId.value = null
    
    deletingImageId.value = imageId
    try {
        await axios.delete(`${API_URL}/coins/${coin.value.id}/gallery/${imageId}`)
        // Remove locally
        if (coin.value.gallery_images) {
            coin.value.gallery_images = coin.value.gallery_images.filter(img => img.id !== imageId)
        }
    } catch (e) {
        console.error("Failed to delete gallery image", e)
        alert(t('form.errors.delete_failed'))
    } finally {
        deletingImageId.value = null
    }
}

const openViewerForPath = (path) => {
    viewerImage.value = getImageUrl(path)
    viewerOpen.value = true
} 




onMounted(async () => {
    fetchGroups() // Fetch groups
    try {
        const response = await axios.get(`${API_URL}/coins/${route.params.id}`)
        coin.value = response.data
    } catch (error) {
        console.error('Error fetching coin:', error)
        router.push('/list')
    } finally {
        loading.value = false
        // Fetch Numista sync status (or just rely on coin data)
        // checkNumistaCount() // Defined elsewhere? Not in variables I saw.
        // Wait, checkNumistaCount is not defined in the code I inserted or saw?
        // Ah, numistaCount is computed.
        // I will just add fetchLinks() here.
        fetchLinks()
        fetchGroupImages()
    }
})

watch(activeTab, (newTab) => {
    if (newTab === 'links' && links.value.length === 0) {
        fetchLinks()
    }
    if (newTab === 'stats' && !coinStats.value) {
        fetchStats()
    }
})
// Numista Manual Selection
const numistaModalOpen = ref(false)
const applyingNumista = ref(false)
const selectedNumistaId = ref(null)
const manualNumistaId = ref(null)

const parsedNumistaSearch = computed(() => {
    if (!coin.value || !coin.value.numista_search) return null
    // If it's already an object, return it directly
    if (typeof coin.value.numista_search === 'object') {
        return coin.value.numista_search
    }
    try {
        return JSON.parse(coin.value.numista_search)
    } catch (e) {
        console.error("Failed to parse numista_search", e)
        return null
    }
})

const numistaResults = computed(() => parsedNumistaSearch.value?.types || [])
const numistaCount = computed(() => parsedNumistaSearch.value?.count || 0)

const applyNumistaResult = async (numistaId) => {
    if (!coin.value) return
    applyingNumista.value = true
    selectedNumistaId.value = numistaId
    try {
        await axios.post(`${API_URL}/coins/${coin.value.id}/apply-numista/${numistaId}`)
        // Refresh
        const res = await axios.get(`${API_URL}/coins/${coin.value.id}`)
        coin.value = res.data
        numistaModalOpen.value = false
        // Optional success toast
    } catch (e) {
        console.error("Failed to apply numista result", e)
        alert('Failed: ' + (e.response?.data?.error || e.message))
    } finally {
        applyingNumista.value = false
        selectedNumistaId.value = null
    }
}

const applyManualNumista = () => {
    if (manualNumistaId.value) {
        applyNumistaResult(manualNumistaId.value)
    }
}

import { useSettingsStore } from '../stores/settings'
import { storeToRefs } from 'pinia'

const formatCurrency = (val) => {
  const settingsStore = useSettingsStore()
  if (settingsStore.privacyMode) return '***'
  return new Intl.NumberFormat('es-ES', { style: 'currency', currency: 'EUR' }).format(val || 0)
}

const getMaterialName = (material) => {
    if (!material) return 'N/A'
    return material.split('(')[0].trim()
}

const getMaterialComposition = (material) => {
    if (!material || !material.includes('(')) return null
    const parts = material.split('(')
    if (parts.length < 2) return null
    return '(' + parts.slice(1).join('(').trim()
}

const getAppraisalText = () => {
    const settingsStore = useSettingsStore()
    if (settingsStore.privacyMode) return '***'
    if (!coin.value) return ''
    
    if (coin.value.min_value > 0 && coin.value.max_value > 0) return `${coin.value.min_value}€ - ${coin.value.max_value}€`
    if (coin.value.min_value > 0) return `> ${coin.value.min_value}€`
    if (coin.value.max_value > 0) return `< ${coin.value.max_value}€`
    return ''
}


</script>
