<template>
  <div class="h-full w-full relative z-0">
    <div ref="chartContainer" class="w-full h-full"></div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch, onBeforeUnmount } from 'vue'
import Highcharts from 'highcharts'
import HighchartsMap from 'highcharts/modules/map'

// Initialize the map module
if (typeof HighchartsMap === 'function') {
    HighchartsMap(Highcharts)
} else if (HighchartsMap && typeof HighchartsMap.default === 'function') {
    HighchartsMap.default(Highcharts)
}

const props = defineProps({
  data: {
    type: Object,
    default: () => ({})
  }
})

const chartContainer = ref(null)
let chart = null

// Name Mapping: Spanish -> Highcharts Map Name (English usually)
// We need to match the "name" property in Highcharts World Map
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



const transformData = () => {
    // Convert the props.data object {'España': 10} to Highcharts array [{'code3': 'ESP', 'value': 10}]
    // Since we don't have ISO codes easily available without a big list, we will try to join by 'name'
    // The Highcharts map data has 'name' property.
    
    const data = []
    
    // 1. Iterate over our data
    for (const [key, count] of Object.entries(props.data)) {
        if (count > 0) {
             let mapName = key
             // Try to map from Spanish to English if needed
             if (countryMapping[key]) {
                 mapName = countryMapping[key]
             }
             
             data.push({
                 name: mapName,
                 value: count,
                 // Add color explicitly for neon effect if handled by data, 
                 // but easier to handle via series color.
             })
        }
    }
    return data
}

const initChart = async () => {
    // Load Map Data dynamically
    const topology = await fetch(
        'https://code.highcharts.com/mapdata/custom/world.topo.json'
    ).then(response => response.json());

    const data = transformData()

    chart = Highcharts.mapChart(chartContainer.value, {
        chart: {
            backgroundColor: 'transparent',
            map: topology
        },
        title: {
            text: ''
        },
        mapNavigation: {
            enabled: true,
            buttonOptions: {
                verticalAlign: 'bottom',
                align: 'right'
            }
        },
        legend: {
            enabled: false
        },
        credits: { enabled: false },
        plotOptions: {
            map: {
                allAreas: true,
                joinBy: ['name', 'name'], // Join by name property
                nullColor: '#1a1a24',     // Design Req: Dark background for empty
                borderColor: '#404040',   // Design Req: Subtle Border
                borderWidth: 0.5,
                states: {
                    hover: {
                        color: null,      // Inherit from point
                        brightness: 0.1,  // Design Req: Lighten on hover
                        borderColor: '#ffffff',
                        borderWidth: 1
                    }
                }
            }
        },
        series: [{
            name: 'Coins',
            data: data,
            color: '#00FFFF', // Design Req: Neon Cyan for active countries
            states: {
                hover: {
                    color: '#80ffff' // Slightly brighter cyan
                }
            },
            tooltip: {
                headerFormat: '',
                pointFormat: '<b>{point.name}</b>: {point.value} coins'
            }
        }]
    });
}

onMounted(() => {
    initChart()
})

watch(() => props.data, () => {
    if (chart) {
        chart.series[0].setData(transformData())
    }
}, { deep: true })

onBeforeUnmount(() => {
    if (chart) {
        chart.destroy()
    }
})
</script>

<style scoped>
/* Ensure tooltip z-index is high enough if needed */
</style>
