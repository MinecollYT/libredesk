import { createRouter, createWebHistory } from 'vue-router'
import DashboardView from '../views/DashboardView.vue'
import ConversationsView from '../views/ConversationView.vue'
import UserLoginView from '../views/UserLoginView.vue'
import AccountView from '@/views/AccountView.vue'
import AdminView from '@/views/AdminView.vue'

const routes = [
  {
    path: '/',
    name: 'login',
    component: UserLoginView
  },
  {
    path: '/dashboard',
    name: 'dashboard',
    component: DashboardView
  },
  {
    path: '/conversations/:uuid?',
    name: 'conversations',
    component: ConversationsView,
    props: true
  },
  {
    path: '/account/:page?',
    name: 'account',
    component: AccountView,
    props: true,
    beforeEnter: (to, from, next) => {
      if (!to.params.page) {
        next({ ...to, params: { ...to.params, page: 'profile' } })
      } else {
        next()
      }
    }
  },
  {
    path: '/admin',
    name: 'admin',
    component: AdminView,
    children: [
      {
        path: 'inboxes',
        component: () => import('@/components/admin/inbox/Inbox.vue')
      },
      {
        path: 'inboxes/new',
        component: () => import('@/components/admin/inbox/NewInbox.vue')
      },
      {
        path: 'inboxes/:id/edit',
        props: true,
        component: () => import('@/components/admin/inbox/EditInbox.vue')
      },
      {
        path: 'notification',
        component: () => import('@/components/admin/notification/NotificationSetting.vue')
      },
      {
        path: 'teams',
        component: () => import('@/components/admin/team/Team.vue')
      },
      {
        path: 'teams/users',
        component: () => import('@/components/admin/team/users/UsersCard.vue')
      },
      {
        path: 'teams/users/new',
        component: () => import('@/components/admin/team/users/AddUserForm.vue')
      },
      {
        path: 'teams/users/:id/edit',
        props: true,
        component: () => import('@/components/admin/team/users/EditUserForm.vue')
      },
      {
        path: 'teams/teams',
        component: () => import('@/components/admin/team/teams/Teams.vue')
      },
      {
        path: 'teams/teams/new',
        component: () => import('@/components/admin/team/teams/AddTeamForm.vue')
      },
      {
        path: 'teams/teams/:id/edit',
        props: true,
        component: () => import('@/components/admin/team/teams/EditTeamForm.vue')
      },
      {
        path: 'teams/roles',
        component: () => import('@/components/admin/team/roles/Roles.vue')
      },
      {
        path: 'teams/roles/new',
        component: () => import('@/components/admin/team/roles/NewRole.vue')
      },
      {
        path: 'teams/roles/:id/edit',
        props: true,
        component: () => import('@/components/admin/team/roles/EditRole.vue')
      },
      {
        path: 'automations',
        component: () => import('@/components/admin/automation/Automation.vue')
      },
      {
        path: 'automations/new',
        props: true,
        component: () => import('@/components/admin/automation/CreateOrEditRule.vue')
      },
      {
        path: 'automations/:id/edit',
        props: true,
        component: () => import('@/components/admin/automation/CreateOrEditRule.vue')
      },
      {
        path: 'general',
        component: () => import('@/components/admin/general/General.vue')
      },
      {
        path: 'templates',
        component: () => import('@/components/admin/templates/Templates.vue')
      },
      {
        path: 'templates/:id/edit',
        props: true,
        component: () => import('@/components/admin/templates/AddEditTemplate.vue')
      },
      {
        path: 'templates/new',
        component: () => import('@/components/admin/templates/AddEditTemplate.vue')
      },
      {
        path: 'oidc',
        component: () => import('@/components/admin/oidc/OIDC.vue')
      },
      {
        path: 'oidc/:id/edit',
        props: true,
        component: () => import('@/components/admin/oidc/AddEditOIDC.vue')
      },
      {
        path: 'oidc/new',
        component: () => import('@/components/admin/oidc/AddEditOIDC.vue')
      },
      {
        path: 'conversations',
        component: () => import('@/components/admin/conversation/Conversation.vue')
      },
      {
        path: 'conversations/tags',
        component: () => import('@/components/admin/conversation/tags/Tags.vue')
      },
      {
        path: 'conversations/statuses',
        component: () => import('@/components/admin/conversation/status/Status.vue')
      },
      {
        path: 'conversations/canned-responses',
        component: () => import('@/components/admin/conversation/canned_responses/CannedResponses.vue')
      }
    ]
  },
  // Fallback to dashboard.
  {
    path: '/:pathMatch(.*)*',
    redirect: (to) => {
      alert(`Redirecting to dashboard from: ${to.fullPath}`)
      return '/dashboard'
    }
  }
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: routes
})

export default router
