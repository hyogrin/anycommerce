import Vue from 'vue'
import App from './App.vue'
import router from './router';
import { Interactions } from 'aws-amplify';
import { components } from 'aws-amplify-vue';
import store from '@/store/store';
import Amplitude from 'amplitude-js'
import mParticle from '@mparticle/web-sdk';
// import AmplifyStore from '@/store/store';
import VueGtag from "vue-gtag";

import './styles/tokens.css'

// Base configuration for Amplify
const amplifyConfig = {}

// if (AmplifyStore.state.user?.id) {
//   amplifyConfig.Analytics.AWSPinpoint.endpoint = {
//     userId: AmplifyStore.state.user.id
//   }
// }

// Only add Personalize event tracking if configured.
if (process.env.VUE_APP_PERSONALIZE_TRACKING_ID && process.env.VUE_APP_PERSONALIZE_TRACKING_ID != 'NONE') {
  // Amazon Personalize event tracker.

  amplifyConfig.Analytics.AmazonPersonalize = {
    trackingId: process.env.VUE_APP_PERSONALIZE_TRACKING_ID,
    region: process.env.VUE_APP_AWS_REGION,
    // OPTIONAL - The number of events to be deleted from the buffer when flushed.
    flushSize: 5,
    // OPTIONAL - The interval in milliseconds to perform a buffer check and flush if necessary.
    flushInterval: 2000, // 2s
  }
}

// Initialize Amplitude if a valid API key is specified.
if (process.env.VUE_APP_AMPLITUDE_API_KEY && process.env.VUE_APP_AMPLITUDE_API_KEY != 'NONE') {
  Amplitude.getInstance().init(process.env.VUE_APP_AMPLITUDE_API_KEY)
}

// Initialize mParticle if a valid API key is specified.
if (process.env.VUE_APP_MPARTICLE_API_KEY && process.env.VUE_APP_MPARTICLE_API_KEY != 'NONE') {
  const mParticleConfig = {
    isDevelopmentMode: true,
    logLevel: "verbose"
  };
  mParticle.init(process.env.VUE_APP_MPARTICLE_API_KEY, mParticleConfig);
}

if (process.env.VUE_APP_GOOGLE_ANALYTICS_ID && process.env.VUE_APP_GOOGLE_ANALYTICS_ID != 'NONE') {
  Vue.use(VueGtag, {
    config: {
      id: process.env.VUE_APP_GOOGLE_ANALYTICS_ID,
      params: {
        send_page_view: false
      }
    }
  }, router);
}
else {
  Vue.use(VueGtag, {
    enabled: false,
    disableScriptLoad: true
  });
}

// Set the configuration
Interactions.configure(amplifyConfig);

require('dotenv').config()

Vue.config.productionTip = false

new Vue({
  el: '#app',
  router: router,
  template: '<App/>',
  store,
  components: {
    App,
    ...components
  },
  render: h => h(App)
}).$mount('#app')