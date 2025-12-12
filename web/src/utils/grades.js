export const GRADES = {
    'PROOF': { level: 80, key: 'grades.PROOF' },
    'FDC': { level: 70, key: 'grades.FDC' },
    'SC': { level: 60, key: 'grades.SC' },
    'EBC': { level: 50, key: 'grades.EBC' },
    'MBC': { level: 40, key: 'grades.MBC' },
    'BC': { level: 30, key: 'grades.BC' },
    'RC': { level: 20, key: 'grades.RC' },
    'MC': { level: 10, key: 'grades.MC' }
}

export const GRADE_ORDER = ['MC', 'RC', 'BC', 'MBC', 'EBC', 'SC', 'FDC', 'PROOF']

export const BASE_GRADES = ['MC', 'RC', 'BC', 'MBC', 'EBC', 'SC', 'FDC', 'PROOF']

export const normalizeGrade = (grade) => {
    if (!grade) return null
    // Remove + or - or ++ or --
    return grade.replace(/[+\-]+/g, '').trim()
}

export const getGradeValue = (grade) => {
    const base = normalizeGrade(grade)
    if (!base || !GRADES[base]) return 0
    let val = GRADES[base].level

    if (grade.includes('++')) val += 2
    else if (grade.includes('+')) val += 1
    else if (grade.includes('--')) val -= 2
    else if (grade.includes('-')) val -= 1

    return val
}

export const getGradeColor = (grade) => {
    const base = normalizeGrade(grade)
    switch (base) {
        case 'MC':
        case 'RC': return '#EF4444' // Red
        case 'BC': return '#F97316' // Orange
        case 'MBC': return '#EAB308' // Yellow
        case 'EBC': return '#84CC16' // Lime
        case 'SC': return '#22C55E' // Green
        case 'FDC': return '#14B8A6' // Teal
        case 'PROOF': return '#3B82F6' // Blue
        default: return '#9CA3AF' // Gray
    }
}
