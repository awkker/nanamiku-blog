import { defineConfig } from '@playwright/test'

const frontendPort = Number(process.env.SMOKE_FRONTEND_PORT || '4322')
const frontendURL = process.env.SMOKE_FRONTEND_URL || `http://127.0.0.1:${frontendPort}`
const browserChannel = process.env.PLAYWRIGHT_BROWSER_CHANNEL?.trim() || undefined
const cmsOrigin = process.env.SMOKE_CMS_API_ORIGIN || process.env.SMOKE_BACKEND_URL || 'http://127.0.0.1:8080'

export default defineConfig({
  testDir: './e2e',
  timeout: 45_000,
  expect: {
    timeout: 8_000,
  },
  fullyParallel: false,
  forbidOnly: Boolean(process.env.CI),
  retries: process.env.CI ? 1 : 0,
  reporter: process.env.CI ? [['github'], ['html', { open: 'never' }]] : [['list']],
  use: {
    baseURL: frontendURL,
    browserName: 'chromium',
    channel: browserChannel,
    headless: true,
    viewport: { width: 1440, height: 960 },
    trace: 'retain-on-failure',
    screenshot: 'only-on-failure',
  },
  webServer: {
    command: `CMS_API_ORIGIN=${cmsOrigin} npm run dev -- --host 127.0.0.1 --port ${frontendPort}`,
    cwd: '.',
    url: frontendURL,
    timeout: 120_000,
    reuseExistingServer: false,
  },
})
