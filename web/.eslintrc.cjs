/* eslint-env node */
module.exports = {
    root: true,
    extends: [
        'plugin:vue/essential'
    ],
    parserOptions: {
        ecmaVersion: 'latest'
    },
    env: {
        browser: true,
        node: true,
        es2021: true
    },
    rules: {
        'vue/multi-word-component-names': 'off',
        'no-unused-vars': 'warn',
        'no-undef': 'off' // Turning off temporarily if env doesn't catch all globals
    }
}
