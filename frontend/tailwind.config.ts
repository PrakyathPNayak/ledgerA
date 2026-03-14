import type { Config } from 'tailwindcss'

const config: Config = {
    darkMode: ['class'],
    content: ['./index.html', './src/**/*.{ts,tsx}'],
    theme: {
        extend: {
            colors: {
                brand: {
                    50: '#f8fafc',
                    900: '#0f172a',
                },
            },
        },
    },
    plugins: [],
}

export default config
