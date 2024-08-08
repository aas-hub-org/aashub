/**
 * router/index.ts
 *
 * Automatic routes for `./src/pages/*.vue`
 */

// Composables
import { createRouter, createWebHistory } from 'vue-router';
// Views
import HelloWorld from '@/pages/HelloWorld.vue';

const routes = [{ path: '/', component: HelloWorld }];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;
