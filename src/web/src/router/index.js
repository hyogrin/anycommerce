import Vue from 'vue';
import Router from 'vue-router';
import Main from '@/public/Main.vue'
import ProductDetail from '@/public/ProductDetail.vue'
import CategoryDetail from '@/public/CategoryDetail.vue'
import Help from '@/public/Help.vue'
import Cart from '@/public/Cart.vue'
import AuthScreen from '@/public/Auth.vue'
import RegisterScreen from '@/public/Register.vue'
import Checkout from '@/public/Checkout.vue'
import Orders from '@/authenticated/Orders.vue'
import Admin from '@/authenticated/Admin.vue'
import ShopperSelectPage from '@/authenticated/ShopperSelectPage'

import Location from "@/public/Location";
import Collections from "@/public/Collections";

import { EventBus } from '@/event-bus';
import AmplifyStore from '@/store/store';

import { RepositoryFactory } from '@/repositories/RepositoryFactory'
import { AnalyticsHandler } from '@/analytics/AnalyticsHandler'

const UsersRepository = RepositoryFactory.get('users')

Vue.use(Router);
// Explicitly add only components needed to keep deployment as small as possible

// Load User
// eslint-disable-next-line
getUser().then((_user, error) => {
  if (error) {
    // eslint-disable-next-line
    console.log(error)
  }
})

function getVolatileUser() {
  return AmplifyStore.getters.volatileUser
}

// Event Handles for Authentication
EventBus.$on('authState', async (state) => {
  if (state === 'signedOut') {
    AmplifyStore.dispatch('logout');
    AnalyticsHandler.clearUser()

    if (router.currentRoute.path !== '/') router.push({ path: '/' })
  }
  else if (state === 'signedIn') {
    const volatileUser = await getVolatileUser()

    let storeUser = null

    const hasAssignedShopperProfile = !!volatileUser.attributes?.['custom:profile_user_id'];

    if (hasAssignedShopperProfile) {
      const { data } = await UsersRepository.getUserByID(volatileUser.attributes['custom:profile_user_id'])
      storeUser = data
    }
    else {
      // Perhaps our auth user is one without an associated "profile" - so there may be no profile_user_id on the
      // volatile record - so we see if we've created a user in the user service (see below) for this non-profile user
      const { data } = await UsersRepository.getUserByUsername(volatileUser.username)
      storeUser = data
    }

    const credentials = null;

    if (!storeUser.id) {
      // Store user does not exist. Create one on the fly.
      // This takes the personalize User ID which was a UUID4 for the current session and turns it into a user user ID.
      console.log('store user does not exist for volatile user... creating on the fly')
      // let identityId = credentials ? credentials.identityId : null;
      // let provisionalUserId = AmplifyStore.getters.personalizeUserID;
      const { data } = await UsersRepository.getUserByUsername(volatileUser.username)
      storeUser = data
    }

    console.log('Syncing store user state to volatile user custom attributes')

    // Sync identityId with user to support reverse lookup.
    if (credentials && storeUser.identity_id != credentials.identityId) {
      console.log('Syncing credentials identity_id with store user profile')
      storeUser.identity_id = credentials.identityId
    }

    // Update last sign in and sign up dates on user.
    let newSignUp = false

    const now = new Date()
    storeUser.last_sign_in_date = now.toISOString()

    if (!storeUser.sign_up_date) {
      storeUser.sign_up_date = now.toISOString()
      newSignUp = true
    }

    AmplifyStore.commit('setUser', storeUser);

    if (newSignUp && !hasAssignedShopperProfile) {
      AmplifyStore.dispatch('firstTimeSignInDetected');

      router.push({ path: '/shopper-select' });
    } else {
      router.push({ path: '/' });
    }
  }
  else if (state === 'profileChanged') {
    const volatileUser = await getVolatileUser()
    const storeUser = AmplifyStore.state.user

    if (volatileUser && storeUser) {
      // Store user exists. Use this as opportunity to sync store user
      // attributes to volatile custom attributes.
      // Vue.prototype.$Amplify.Auth.updateUserAttributes(volatileUser, {
      //   'custom:profile_user_id': storeUser.id.toString(),
      //   'custom:profile_email': storeUser.email,
      //   'custom:profile_first_name': storeUser.first_name,
      //   'custom:profile_last_name': storeUser.last_name,
      //   'custom:profile_gender': storeUser.gender,
      //   'custom:profile_age': storeUser.age.toString(),
      //   'custom:profile_segment': storeUser.segment
      // })
    } else if (!volatileUser && storeUser) {
      AmplifyStore.commit('setVolatileUser', storeUser)
      router.push({ path: '/' });
    }

    // Sync identityId with user to support reverse lookup.
    // const credentials = await Credentials.get();
    const credentials = null;
    if (credentials && storeUser.identity_id != credentials.identityId) {
      console.log('Syncing credentials identity_id with store user profile')
      // storeUser.identity_id = credentials.identityId
      // UsersRepository.updateUser(storeUser)
    }
  }
});

// Get store user from local storage, making sure session is authenticated
async function getUser() {
  const volatileUser = await getVolatileUser()
  if (!volatileUser) {
    AmplifyStore.commit('setUser', null);
  }

  return AmplifyStore.state.user;
}

// Routes
const router = new Router({
  // mode: 'history',
  routes: [
    {
      path: '/',
      name: 'Main',
      component: Main,
      meta: { requiresAuth: false }
    },
    {
      path: '/product/:id',
      name: 'ProductDetail',
      component: ProductDetail,
      props: route => ({ discount: route.query.di === "true" || route.query.di === true }),
      meta: { requiresAuth: false }
    },
    {
      path: '/category/:id',
      name: 'CategoryDetail',
      component: CategoryDetail,
      meta: { requiresAuth: false }
    },
    {
      path: '/help',
      name: 'Help',
      component: Help,
      meta: { requiresAuth: false }
    },
    {
      path: '/orders',
      name: 'Orders',
      component: Orders,
      meta: { requiresAuth: true }
    },
    {
      path: '/cart',
      name: 'Cart',
      component: Cart,
      meta: { requiresAuth: false }
    },
    {
      path: '/checkout',
      name: 'Checkout',
      component: Checkout,
      meta: { requiresAuth: false }
    },
    {
      path: '/admin',
      name: 'Admin',
      component: Admin,
      meta: { requiresAuth: true }
    },
    {
      path: '/auth',
      name: 'Authenticator',
      component: AuthScreen,
    },
    {
      path: '/register',
      name: 'Register',
      component: RegisterScreen,
    },
    {
      path: '/shopper-select',
      name: 'ShopperSelect',
      component: ShopperSelectPage,
      meta: { requiresAuth: true },
    },
    {
      path: '/location',
      name: 'Location',
      component: Location,
      meta: { requiresAuth: true }
    },
    {
      path: '/collections',
      name: 'Collections',
      component: Collections,
      meta: { requiresAuth: true }
    }
  ],
  scrollBehavior(_to, _from, savedPosition) {
    if (savedPosition) {
      return savedPosition
    } else {
      return { x: 0, y: 0 }
    }
  }
});

// Check if we need to redirect to welcome page - if redirection has never taken place and user is not authenticated
// Check For Authentication
router.beforeResolve(async (to, from, next) => {
  AmplifyStore.dispatch('pageVisited', from.fullPath);

  if (to.matched.some(record => record.meta.requiresAuth)) {
    const user = await getUser();
    if (!user) {
      return next({
        path: '/auth',
        query: {
          redirect: to.fullPath,
        }
      });
    }
    return next()
  }
  return next()
})

export default router