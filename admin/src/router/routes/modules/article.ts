import { DEFAULT_LAYOUT } from '../base';
import { AppRouteRecordRaw } from '../types';

export default {
  path: '/article',
  component: DEFAULT_LAYOUT,
  meta: {
    title: 'Articles',
    requiresAuth: true,
    icon: 'icon-apps',
    order: 4,
    hideChildrenInMenu: true,
  },
  children: [
    {
      path: '',
      name: 'Article',
      component: () => import('@/views/pages/article/index.vue'),
      meta: {
        title: 'Articles',
        requiresAuth: true,
        roles: ['*'],
      },
    },
  ],
} as AppRouteRecordRaw;
