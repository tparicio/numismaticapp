<template>
  <div class="h-full w-full relative z-0 bg-base-100">
    <l-map 
      ref="map" 
      v-model:zoom="zoom" 
      :center="[20, 0]" 
      :use-global-leaflet="false"
      :options="{ 
        maxBounds: [[-90, -180], [90, 180]], 
        maxBoundsViscosity: 1.0, 
        minZoom: 1.5,
        worldCopyJump: false
      }"
    >
      <l-tile-layer
        url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
        layer-type="base"
        name="OpenStreetMap"
        :no-wrap="true"
      ></l-tile-layer>
      <l-geo-json 
        v-if="geojson" 
        :geojson="geojson" 
        :options="geoJsonOptions"
      ></l-geo-json>
    </l-map>
    
    <!-- Legend overlay -->
    <div class="absolute bottom-4 left-4 bg-white/90 p-2 rounded shadow text-xs z-[1000] backdrop-blur-sm border border-base-300">
      <div class="font-bold mb-1 text-base-content">Coins</div>
      <div class="flex flex-col gap-1">
        <div class="flex items-center gap-2"><span class="w-3 h-3 bg-[#08306b]"></span> 50+</div>
        <div class="flex items-center gap-2"><span class="w-3 h-3 bg-[#2171b5]"></span> 20-49</div>
        <div class="flex items-center gap-2"><span class="w-3 h-3 bg-[#6baed6]"></span> 10-19</div>
        <div class="flex items-center gap-2"><span class="w-3 h-3 bg-[#bdd7e7]"></span> 1-9</div>
        <div class="flex items-center gap-2"><span class="w-3 h-3 bg-base-200 border border-base-300"></span> 0</div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, watch } from 'vue'
import "leaflet/dist/leaflet.css"
import { LMap, LTileLayer, LGeoJson } from "@vue-leaflet/vue-leaflet"

const props = defineProps({
  data: {
    type: Object,
    default: () => ({})
  }
})

const zoom = ref(2)
const geojson = ref(null)

// Normalize name helper
const normalize = (str) => str ? str.toLowerCase().trim() : ''

// Name Mapping: Spanish -> English/ISO Name in GeoJSON
const countryMapping = {
  "España": "Spain",
  "Reino Unido": "United Kingdom",
  "Estados Unidos": "United States of America",
  "Alemania": "Germany",
  "Francia": "France",
  "Italia": "Italy",
  "Portugal": "Portugal",
  "Grecia": "Greece",
  "México": "Mexico",
  "Canadá": "Canada",
  "Australia": "Australia",
  "Japón": "Japan",
  "China": "China",
  "Rusia": "Russia",
  "Brasil": "Brazil",
  "Argentina": "Argentina",
  "Chile": "Chile",
  "Perú": "Peru",
  "Colombia": "Colombia",
  "Irlanda": "Ireland",
  "Bélgica": "Belgium",
  "Países Bajos": "Netherlands",
  "Austria": "Austria",
  "Suiza": "Switzerland",
  "Suecia": "Sweden",
  "Noruega": "Norway",
  "Dinamarca": "Denmark",
  "Finlandia": "Finland",
  "Polonia": "Poland"
}

// Helper to get count for a feature
const getCount = (feature) => {
    const featureName = normalize(feature.properties.name)
    const featureId = normalize(feature.id) // Sometimes ISO code is in ID

    // 1. Check direct match (Spanish key matching Feature Name)
    // iterate props.data keys
    for (const [key, val] of Object.entries(props.data)) {
        if (normalize(key) === featureName) return val
    }

    // 2. Check via mapping
    // Iterate mapping to find if we have data for a mapped country
    for (const [spanishKey, englishVal] of Object.entries(countryMapping)) {
        if (normalize(englishVal) === featureName || normalize(englishVal) === featureId) {
            if (props.data[spanishKey]) return props.data[spanishKey]
        }
    }
    
    return 0
}

// Color scale
const getColor = (d) => {
    return d > 50 ? '#08306b' :
           d > 20  ? '#2171b5' :
           d > 10  ? '#6baed6' :
           d > 0   ? '#bdd7e7' :
                     'transparent'; 
}

const style = (feature) => {
    const count = getCount(feature)
    const hasData = count > 0

    return {
        fillColor: hasData ? getColor(count) : '#f8f9fa', // Light gray for empty
        weight: 1,
        opacity: 1,
        color: '#dee2e6', // Subtle border
        dashArray: '3',
        fillOpacity: hasData ? 0.7 : 0.3
    }
}

const onEachFeature = (feature, layer) => {
    const count = getCount(feature)
    const countryName = feature.properties.name
    
    // Find our internal name for display
    let displayName = countryName
    for (const [spanishKey, englishVal] of Object.entries(countryMapping)) {
        if (normalize(englishVal) === normalize(countryName)) {
            displayName = spanishKey
            break
        }
    }

    if (count > 0) {
        layer.bindTooltip(`<strong>${displayName}</strong><br/>Coins: ${count}`, { sticky: true })
        // Highlight effect
        layer.on({
            mouseover: (e) => {
                const layer = e.target;
                layer.setStyle({
                    weight: 2,
                    color: '#666',
                    dashArray: '',
                    fillOpacity: 0.9
                });
                layer.bringToFront();
            },
            mouseout: (e) => {
                const layer = e.target;
                // create a mock feature to reset style
                layer.setStyle(style(feature)); 
            }
        });
    } else {
       layer.bindTooltip(`${countryName}`, { sticky: true }) 
    }
}

const geoJsonOptions = computed(() => ({
    style: style,
    onEachFeature: onEachFeature
}))

onMounted(async () => {
    try {
        const response = await fetch('https://raw.githubusercontent.com/johan/world.geo.json/master/countries.geo.json')
        geojson.value = await response.json()
    } catch (e) {
        console.error("Failed to load GeoJSON", e)
    }
})
</script>

<style>
/* Leaflet generic fix for z-index issues if any */
.leaflet-pane { z-index: 10 !important; }
.leaflet-bottom { z-index: 20 !important; }
</style>
