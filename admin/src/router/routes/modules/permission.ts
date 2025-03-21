import { DEFAULT_LAYOUT } from '../base';
import { AppRouteRecordRaw } from '../types';

export default {
  path: '/permission',
  component: DEFAULT_LAYOUT,
  meta: {
    title: 'Permission',
    requiresAuth: true,
    icon: 'icon-apps',
    order: 8,
  },
  children: [
    {
      path: 'role',
      name: 'Role',
      component: () => import('@/views/pages/system/role/index.vue'),
      meta: {
        title: 'Roles',
        requiresAuth: true,
        roles: ['*'],
      },
    },
    {
      path: 'menu',
      name: 'Menu',
      component: () => import('@/views/pages/system/menu/index.vue'),
      meta: {
        title: 'Menu',
        requiresAuth: true,
        roles: ['*'],
      },
    },
    {
      path: 'index',
      name: 'Permission',
      component: () => import('@/views/pages/system/permission/index.vue'),
      meta: {
        title: 'Permission',
        requiresAuth: true,
        roles: ['*'],
      },
    },
  ],
} as AppRouteRecordRaw;
