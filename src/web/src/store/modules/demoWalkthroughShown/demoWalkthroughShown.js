export const demoWalkthroughShown = {
  state: () => ({ shown: true }),
  mutations: {
    setDemoWalkthroughShown: (state, newVisited) => (state.shown = newVisited),
  },
  actions: {
    markDemoWalkthroughAsShown: ({ commit }) => commit('setDemoWalkthroughShown', true),
    firstTimeSignInDetected: ({ commit }) => commit('setDemoWalkthroughShown', true),
  },
};
