/**
 * plugins/index.js
 *
 * Automatically included in `./src/main.js`
 */

// Plugins
import vuetify from './vuetify'
import pinia from '@/stores'
import router from '@/router'
import V3ScrollLock from 'v3-scroll-lock';
import VBodyScrollLock from 'v-body-scroll-lock'

export function registerPlugins (app) {
  app
    .use(vuetify)
    .use(router)
    .use(pinia)
    .use(V3ScrollLock)
    .use(VBodyScrollLock)
}
