export const formatMintage = (value) => {
    if (!value) return 'N/A'

    if (value >= 1000000) {
        return (value / 1000000).toFixed(1).replace(/\.0$/, '') + 'M'
    }
    if (value >= 1000) {
        return (value / 1000).toFixed(1).replace(/\.0$/, '') + 'k'
    }

    return value.toLocaleString()
}
