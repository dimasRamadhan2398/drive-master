import { test } from '@playwright/test';

test('check page load and console errors', async ({ page }) => {
  page.on('console', msg => {
    console.log(`[CONSOLE ${msg.type().toUpperCase()}] ${msg.text()}`);
  });
  page.on('pageerror', err => {
    console.log('[PAGE ERROR]', err.stack || err.message);
  });

  console.log('--- NAVIGATING TO ROOT ---');
  await page.goto('http://localhost:3000/', { waitUntil: 'networkidle' });
  console.log('Final URL after load:', page.url());
  
  await page.waitForTimeout(2000);
  console.log('Final URL after 2s:', page.url());

  const mainExists = await page.evaluate(() => {
    const main = document.querySelector('main');
    return main ? { html: main.innerHTML.substring(0, 500), text: main.innerText.substring(0, 500) } : null;
  });
  console.log('Main Element Content:', JSON.stringify(mainExists, null, 2));
});
