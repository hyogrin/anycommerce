import Vue from 'vue';
import Vuex from 'vuex';

import createPersistedState from 'vuex-persistedstate';
import { v4 as uuidv4 } from 'uuid';

import { welcomePageVisited } from './modules/welcomePageVisited/welcomePageVisited';
import { categories } from './modules/categories/categories';
import { cart } from './modules/cart/cart';
import { modal, manageResponsiveModalState } from './modules/modal/modal';
import { demoWalkthroughShown } from './modules/demoWalkthroughShown/demoWalkthroughShown';
import { lastVisitedPage } from './modules/lastVisitedPage/lastVisitedPage';
import { confirmationModal } from './modules/confirmationModal/confirmationModal';

Vue.use(Vuex);

const store = new Vuex.Store({
  modules: { welcomePageVisited, categories, cart, modal, demoWalkthroughShown, lastVisitedPage, confirmationModal },
  state: {
    user: null,
    volatileUser: null,
    provisionalUserID: uuidv4(),
    sessionEventsRecorded: 0,
  },
  mutations: {
    setLoggedOut(state) {
      state.user = null;
      state.volatileUser = null;
      state.provisionalUserID = uuidv4();
      state.sessionEventsRecorded = 0;
    },
    setUser(state, user) {
      if (user && Object.prototype.hasOwnProperty.call(user, 'storage')) {
        // Clear "user.storage" to prevent recursively nested user state
        // from being stored which eventually leads to exhausting local storage.
        user.storage = null;
      }
      state.user = user;
    },
    setVolatileUser(state, volatileUser) {
      state.volatileUser = volatileUser;
    },
    incrementSessionEventsRecorded(state) {
      state.sessionEventsRecorded += 1;
    },
  },
  getters: {
    username: (state) => state.user?.username ?? 'guest',
    personalizeUserID: (state) => {
      return state.user ? state.user.username : state.provisionalUserID;
    },
    personalizeRecommendationsForVisitor: (state) => {
      return state.user || (state.provisionalUserID && state.sessionEventsRecorded > 2);
    },
    volatileUser: (state) => {
      return state.volatileUser;
    },
  },
  actions: {
    setUser: ({ commit }, user) => {
      commit('setUser', user);
    },
    setVolatileUser: ({ commit }, volatileUser) => {
      commit('setVolatileUser', volatileUser);
    },
    logout: ({ commit, dispatch }) => {
      commit('setLoggedOut');
      dispatch('getNewCart');
    },
  },
  plugins: [
    createPersistedState({
      paths: [
        'user',
        'volatileUser',
        'provisionalUserID',
        'sessionEventsRecorded',
        'welcomePageVisited',
        'cart',
        'demoWalkthroughShown',
        'lastVisitedPage',
      ],
    }),
  ],
  strict: process.env.NODE_ENV !== 'production',
});

manageResponsiveModalState(store);

export default store;
