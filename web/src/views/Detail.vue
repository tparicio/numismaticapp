<template>
  <div v-if="coin" class="grid grid-cols-1 lg:grid-cols-2 gap-8">
    <!-- Images Section -->
    <div class="card bg-base-100 shadow-xl">
      <div class="card-body">
        
        <!-- Image Source Toggles -->
        <div class="flex justify-center mb-6">
            <div class="join">
                <button 
                    class="btn join-item" 
                    :class="{ 'btn-primary': activeImageSource === 'processed' }"
                    @click="activeImageSource = 'processed'"
                >
                    {{ $t('details.toggles.processed') }}
                </button>
                <button 
                    class="btn join-item" 
                    :class="{ 'btn-primary': activeImageSource === 'original' }"
                    @click="activeImageSource = 'original'"
                >
                    {{ $t('details.toggles.original') }}
                </button>
            </div>
        </div>

        <div class="flex flex-col sm:flex-row justify-center gap-8 items-center">
            <!-- Front -->
            <div class="text-center relative group">
                <figure class="cursor-zoom-in relative inline-block" @click="openViewer('front')">
                    <div class="absolute inset-0 bg-black bg-opacity-0 group-hover:bg-opacity-10 transition-all flex items-center justify-center z-10"
                         :class="{ 'rounded-full': activeImageSource !== 'original', 'rounded-xl': activeImageSource === 'original' }">
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-10 w-10 text-gray-600 opacity-0 group-hover:opacity-100 transition-opacity" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0zM10 7v3m0 0v3m0-3h3m-3 0H7" />
                        </svg>
                    </div>
                    <img :src="getCurrentImageUrl('front')" 
                         class="shadow-lg max-w-xs hover:scale-105 transition-transform duration-300" 
                         :class="{ 'rounded-full': activeImageSource !== 'original', 'rounded-xl': activeImageSource === 'original' }"
                         :style="{ transform: `rotate(${rotations.front}deg)` }"
                         alt="Front" @error="handleImageError" />
                </figure>
                <button @click.stop="openRotationEditor('front')" class="absolute top-2 right-10 btn btn-circle btn-sm btn-neutral bg-opacity-70 border-none hover:bg-opacity-100 opacity-0 group-hover:opacity-100 transition-opacity z-20" title="Corregir Rotación">
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-4 h-4">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M16.023 9.348h4.992v-.001M2.985 19.644v-4.992m0 0h4.992m-4.993 0l3.181 3.183a8.25 8.25 0 0013.803-3.7M4.031 9.865a8.25 8.25 0 0113.803-3.7l3.181 3.182m0-4.991v4.99" />
                    </svg>
                </button>
                <div class="mt-4 font-bold text-lg flex flex-col items-center gap-2">
                    {{ $t('details.obverse') }}
                    <a v-if="coin?.numista_details?.obverse?.picture" 
                       :href="coin.numista_details.obverse.picture" 
                       target="_blank" 
                       class="btn btn-xs btn-outline gap-1 font-normal text-xs">
                        {{ $t('details.view_sample') }}
                        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-3 h-3">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M13.5 6H5.25A2.25 2.25 0 003 8.25v10.5A2.25 2.25 0 005.25 21h10.5A2.25 2.25 0 0018 18.75V10.5m-10.5 6L21 3m0 0h-5.25M21 3v5.25" />
                        </svg>
                    </a>
                </div>
            </div>

            <!-- Back -->
            <div class="text-center relative group">
                <figure class="cursor-zoom-in relative inline-block" @click="openViewer('back')">
                    <div class="absolute inset-0 bg-black bg-opacity-0 group-hover:bg-opacity-10 transition-all flex items-center justify-center z-10"
                         :class="{ 'rounded-full': activeImageSource !== 'original', 'rounded-xl': activeImageSource === 'original' }">
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-10 w-10 text-gray-600 opacity-0 group-hover:opacity-100 transition-opacity" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0zM10 7v3m0 0v3m0-3h3m-3 0H7" />
                        </svg>
                    </div>
                    <img :src="getCurrentImageUrl('back')" 
                         class="shadow-lg max-w-xs hover:scale-105 transition-transform duration-300" 
                         :class="{ 'rounded-full': activeImageSource !== 'original', 'rounded-xl': activeImageSource === 'original' }"
                         :style="{ transform: `rotate(${rotations.back}deg)` }"
                         alt="Back" @error="handleImageError" />
                </figure>
                <button @click.stop="openRotationEditor('back')" class="absolute top-2 right-10 btn btn-circle btn-sm btn-neutral bg-opacity-70 border-none hover:bg-opacity-100 opacity-0 group-hover:opacity-100 transition-opacity z-20" title="Corregir Rotación">
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-4 h-4">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M16.023 9.348h4.992v-.001M2.985 19.644v-4.992m0 0h4.992m-4.993 0l3.181 3.183a8.25 8.25 0 0013.803-3.7M4.031 9.865a8.25 8.25 0 0113.803-3.7l3.181 3.182m0-4.991v4.99" />
                    </svg>
                </button>
                <div class="mt-4 font-bold text-lg flex flex-col items-center gap-2">
                    {{ $t('details.reverse') }}
                    <a v-if="coin?.numista_details?.reverse?.picture" 
                       :href="coin.numista_details.reverse.picture" 
                       target="_blank" 
                       class="btn btn-xs btn-outline gap-1 font-normal text-xs">
                        {{ $t('details.view_sample') }}
                        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-3 h-3">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M13.5 6H5.25A2.25 2.25 0 003 8.25v10.5A2.25 2.25 0 005.25 21h10.5A2.25 2.25 0 0018 18.75V10.5m-10.5 6L21 3m0 0h-5.25M21 3v5.25" />
                        </svg>
                    </a>
                </div>
            </div>
        </div>
      </div>
    </div>

    <!-- Details Section -->
    <div class="card bg-base-100 shadow-xl h-fit">
      <div class="card-body">
        <h2 v-if="coin.name" class="text-2xl font-bold text-primary mb-1">{{ coin.name }}</h2>
        <h1 class="card-title text-4xl mb-2">{{ coin.country }} {{ coin.face_value }} {{ coin.currency }}</h1>
        <div class="flex gap-2 mb-6">
            <div class="badge badge-lg badge-primary" v-if="coin.year && coin.year !== 0">{{ coin.year }}</div>
            <div class="badge badge-lg badge-secondary">{{ coin.currency }}</div>
            <div class="tooltip" :data-tip="getGradeDescription(coin.grade)" v-if="coin.grade">
                <div class="badge badge-lg badge-accent cursor-help">{{ coin.grade }}</div>
            </div>
        </div>

        <div class="divider">{{ $t('details.sections.details') }}</div>

        <div class="grid grid-cols-2 gap-4">
            <div>
                <span class="font-bold block text-sm text-gray-500">{{ $t('details.labels.material') }}</span>
                <span>{{ coin.material }}</span>
            </div>
            <div>
                <span class="font-bold block text-sm text-gray-500">{{ $t('details.labels.km') }}</span>
                <span class="flex items-center gap-2">
                    {{ coin.km_code || 'N/A' }}
                    <a v-if="getNumistaUrl()" 
                       :href="getNumistaUrl()" 
                       target="_blank" 
                       class="btn btn-xs btn-outline btn-info gap-1"
                       title="Ver en Numista"
                    >
                        N
                        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-3 h-3">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M13.5 6H5.25A2.25 2.25 0 003 8.25v10.5A2.25 2.25 0 005.25 21h10.5A2.25 2.25 0 0018 18.75V10.5m-10.5 6L21 3m0 0h-5.25M21 3v5.25" />
                        </svg>
                    </a>
                    <button v-if="numistaCount > 0" 
                            @click="numistaModalOpen = true" 
                            class="btn btn-xs btn-outline btn-accent gap-1">
                        {{ numistaCount }} resultados numista
                    </button>
                </span>
            </div>
            <div v-if="coin.mint">
                <span class="font-bold block text-sm text-gray-500">{{ $t('details.labels.mint') }}</span>
                <span>{{ coin.mint }}</span>
            </div>
            <div v-if="coin.mintage && coin.mintage > 0">
                <span class="font-bold block text-sm text-gray-500">{{ $t('details.labels.mintage') }}</span>
                <span>{{ formatMintage(coin.mintage) }}</span>
            </div>
            <div v-if="coin.min_value > 0 || coin.max_value > 0">
                <span class="font-bold block text-sm text-gray-500">{{ $t('details.labels.est_value') }}</span>
                <span>{{ coin.min_value }} - {{ coin.max_value }}</span>
            </div>
             <div>
                <span class="font-bold block text-sm text-gray-500">{{ $t('details.labels.added_on') }}</span>
                <span>{{ new Date(coin.created_at).toLocaleDateString() }}</span>
            </div>
        </div>

        <div class="divider">
            {{ $t('details.sections.description') }}
            <div class="badge badge-neutral gap-1 text-xs">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="w-3 h-3">
                    <path fill-rule="evenodd" d="M10 2a1 1 0 011 1v1.323l3.954 1.582 1.599-.8a1 1 0 01.894 1.79l-1.233.616 1.738 5.42a1 1 0 01-.285 1.05A3.989 3.989 0 0115 15a3.989 3.989 0 01-2.667-1.019 1 1 0 01-.633-.73l-3.111-9.64a1.184 1.184 0 01-.722-1.463A1.184 1.184 0 0110 2zM4.08 6.647a1 1 0 011.666-.086l1.248 1.94 1.258-.636a1 1 0 01.894 1.79l-1.233.616 1.738 5.42a1 1 0 01-.285 1.05A3.989 3.989 0 018 17a3.989 3.989 0 01-2.667-1.019 1 1 0 01-.633-.73L1.589 5.611a1 1 0 012.491-1.036z" clip-rule="evenodd" />
                </svg>
                Gemini
            </div>
        </div>
        <p class="whitespace-pre-line">{{ coin.description }}</p>

        <!-- Numista Extended Details -->
        <template v-if="coin.numista_details">
             <div class="divider mt-6">
                {{ $t('details.sections.numista_details') || 'Detalles Numista' }}
                <div class="badge badge-neutral gap-1 text-xs">
                    Numista
                </div>
            </div>
            <!-- Obverse -->
            <div v-if="coin.numista_details.obverse">
                <h3 class="font-bold border-b border-gray-200 dark:border-gray-700 pb-1 mb-2 text-primary">Anverso</h3>
                <p v-if="coin.numista_details.obverse.lettering" class="mb-1">
                    <span class="font-semibold italic opacity-80">Leyenda:</span> 
                    <span class="ml-1">{{ coin.numista_details.obverse.lettering }}</span>
                </p>
                <p v-if="coin.numista_details.obverse.description">
                    <span class="font-semibold italic opacity-80">Descripción:</span>
                    <span class="ml-1">{{ coin.numista_details.obverse.description }}</span>
                </p>
            </div>

            <!-- Reverse -->
            <div v-if="coin.numista_details.reverse">
                <h3 class="font-bold border-b border-gray-200 dark:border-gray-700 pb-1 mb-2 text-primary">Reverso</h3>
                <p v-if="coin.numista_details.reverse.lettering" class="mb-1">
                    <span class="font-semibold italic opacity-80">Leyenda:</span>
                    <span class="ml-1">{{ coin.numista_details.reverse.lettering }}</span>
                </p>
                <p v-if="coin.numista_details.reverse.description">
                    <span class="font-semibold italic opacity-80">Descripción:</span>
                    <span class="ml-1">{{ coin.numista_details.reverse.description }}</span>
                </p>
            </div>

            <!-- Edge -->
            <div v-if="coin.numista_details.edge && (coin.numista_details.edge.description || coin.numista_details.edge.lettering)">
                <h3 class="font-bold border-b border-gray-200 dark:border-gray-700 pb-1 mb-2 text-primary">Edge</h3>
                <p v-if="coin.numista_details.edge.description" class="mb-1">{{ coin.numista_details.edge.description }}</p>
                <p v-if="coin.numista_details.edge.lettering">
                    <span class="font-semibold italic opacity-80">Leyenda:</span>
                    <span class="ml-1">{{ coin.numista_details.edge.lettering }}</span>
                </p>
            </div>

            <!-- Technique -->
            <div v-if="coin.numista_details.technique && coin.numista_details.technique.text">
                 <h3 class="font-bold border-b border-gray-200 dark:border-gray-700 pb-1 mb-2 text-primary">Técnica</h3>
                 <p>{{ coin.numista_details.technique.text }}</p>
            </div>
        </template>


        <div v-if="coin.gemini_model" class="mt-4 flex flex-col text-xs text-gray-400">
           <div class="flex items-center justify-between">
               <div class="flex gap-2 items-center">
                   <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-4 h-4">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M9.813 15.904L9 18.75l-.813-2.846a4.5 4.5 0 00-3.09-3.09L2.25 12l2.846-.813a4.5 4.5 0 003.09-3.09L9 5.25l.813 2.846a4.5 4.5 0 003.09 3.09L15.75 12l-2.846.813a4.5 4.5 0 00-3.09 3.09zM18.259 8.715L18 9.75l-.259-1.035a3.375 3.375 0 00-2.455-2.456L14.25 6l1.036-.259a3.375 3.375 0 002.455-2.456L18 2.25l.259 1.035a3.375 3.375 0 002.456 2.456L21.75 6l-1.035.259a3.375 3.375 0 00-2.456 2.456zM16.894 20.567L16.5 21.75l-.394-1.183a2.25 2.25 0 00-1.423-1.423L13.5 18.75l1.183-.394a2.25 2.25 0 001.423-1.423l.394-1.183.394 1.183a2.25 2.25 0 001.423 1.423l1.183.394-1.183.394a2.25 2.25 0 00-1.423 1.423z" />
                   </svg>
                   <span>AI Generated by {{ coin.gemini_model }} (Temp: {{ coin.gemini_temperature }})</span>
               </div>
               <div class="tooltip" :data-tip="$t('details.reprocess_tooltip') || 'Reprocesar'">
                   <button @click="openReprocessModal" class="btn btn-xs btn-outline btn-primary flex flex-row items-center gap-1 flex-nowrap">
                       <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-4 h-4 flex-shrink-0">
                           <path stroke-linecap="round" stroke-linejoin="round" d="M16.023 9.348h4.992v-.001M2.985 19.644v-4.992m0 0h4.992m-4.993 0l3.181 3.183a8.25 8.25 0 0013.803-3.7M4.031 9.865a8.25 8.25 0 0113.803-3.7l3.181 3.182m0-4.991v4.99" />
                       </svg>
                       <span class="whitespace-nowrap">{{ $t('common.reprocess') || 'Reprocesar' }}</span>
                   </button>
               </div>
                <!-- Numista Reprocess Action -->
                <div class="tooltip" :data-tip="$t('details.reprocess_numista_tooltip') || 'Reprocesar Numista'">
                    <button @click="syncNumista" class="btn btn-xs btn-outline btn-primary flex flex-row items-center gap-1 flex-nowrap" :disabled="syncing">
                        <span v-if="syncing" class="loading loading-spinner loading-xs"></span>
                        <svg v-else xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-4 h-4 flex-shrink-0">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M16.023 9.348h4.992v-.001M2.985 19.644v-4.992m0 0h4.992m-4.993 0l3.181 3.183a8.25 8.25 0 0013.803-3.7M4.031 9.865a8.25 8.25 0 0113.803-3.7l3.181 3.182m0-4.991v4.99" />
                        </svg>
                        <span class="whitespace-nowrap">{{ $t('common.reprocess') }}</span>
                    </button>
                </div>



                <!-- Numista Search Results Action -->
                <div class="tooltip" :data-tip="$t('details.show_results_tooltip') || 'Ver resultados de búsqueda'" v-if="numistaResults.length > 0">
                    <button @click="numistaModalOpen = true" class="btn btn-xs btn-outline btn-accent flex flex-row items-center gap-1 flex-nowrap">
                        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-4 h-4 flex-shrink-0">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M8.25 6.75h12M8.25 12h12m-12 5.25h12M3.75 6.75h.007v.008H3.75V6.75zm.375 0a.375.375 0 11-.75 0 .375.375 0 01.75 0zM3.75 12h.007v.008H3.75V12zm.375 0a.375.375 0 11-.75 0 .375.375 0 01.75 0zM3.75 17.25h.007v.008H3.75v-.008zm.375 0a.375.375 0 11-.75 0 .375.375 0 01.75 0z" />
                        </svg>
                        <span class="whitespace-nowrap">{{ numistaResults.length }} Resultados</span>
                    </button>
                </div>
           </div>
           <div v-if="coin.gemini_details && coin.gemini_details.error" class="mt-1 text-error">
               Warning: {{ coin.gemini_details.error }}
           </div>
        </div>

        <template v-if="coin.notes">
            <div class="divider">{{ $t('details.sections.notes') }}</div>
            <p class="text-sm italic">{{ coin.notes }}</p>
        </template>
        
        <div class="card-actions justify-end mt-8 gap-2">
            <router-link to="/list" class="btn btn-ghost">{{ $t('details.back_gallery') }}</router-link>
            <router-link :to="`/edit/${coin.id}`" class="btn btn-info">
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-5 h-5 mr-1">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M16.862 4.487l1.687-1.688a1.875 1.875 0 112.652 2.652L10.582 16.07a4.5 4.5 0 01-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 011.13-1.897l8.932-8.931zm0 0L19.5 7.125M18 14v4.75A2.25 2.25 0 0115.75 21H5.25A2.25 2.25 0 013 18.75V8.25A2.25 2.25 0 015.25 6H10" />
                </svg>
                {{ $t('common.edit') }}
            </router-link>
            <button @click="deleteModalOpen = true" class="btn btn-error">
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-5 h-5 mr-1">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M14.74 9l-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 01-2.244 2.077H8.084a2.25 2.25 0 01-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 00-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 013.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 00-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 00-7.5 0" />
                </svg>
                {{ $t('common.delete') }}
            </button>
        </div>
      </div>
    </div>
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
        <div class="relative flex justify-center items-center overflow-hidden bg-gray-900 rounded-2xl shadow-2xl h-[400px] sm:h-[500px] border border-gray-700">
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
  <dialog id="delete_modal" class="modal" :class="{ 'modal-open': deleteModalOpen }">
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
  <dialog id="reprocess_modal" class="modal" :class="{ 'modal-open': reprocessModalOpen }">
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

  <!-- Numista Results Modal -->
  <dialog id="numista_modal" class="modal" :class="{ 'modal-open': numistaModalOpen }">
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
import { ref, onMounted, computed, reactive } from 'vue'
import axios from 'axios'
import { useRoute, useRouter } from 'vue-router'
import { normalizeGrade } from '../utils/grades'
import ImageViewer from '../components/ImageViewer.vue'
import GeminiConfig from '../components/GeminiConfig.vue'
import { formatMintage } from '../utils/formatters'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const route = useRoute()
const router = useRouter()
const coin = ref(null)
const API_URL = import.meta.env.VITE_API_URL || '/api/v1'
const STORAGE_URL = ''

const viewerOpen = ref(false)
const viewerImage = ref('')
const activeImageSource = ref('processed') // processed, original

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



onMounted(async () => {
  try {
    const res = await axios.get(`${API_URL}/coins/${route.params.id}`)
    coin.value = res.data
  } catch (e) {
    console.error(e)
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
</script>
