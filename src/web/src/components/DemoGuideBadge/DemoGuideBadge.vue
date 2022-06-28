<template>
  <button type="button" @click="onClick" :aria-label="copy" class="demo-guide-badge align-items-center text-left">
    <!-- <div class="logo mr-1"><img :src="serviceLogo" alt="" class="img-fluid" /></div> -->
    <div :class="{ text: true, 'hide-text-on-small-screens': hideTextOnSmallScreens }">
      <div>{{ copy }}</div>

    </div>
  </button>
</template>

<script>
import { mapActions } from 'vuex';

import { Articles } from '@/partials/AppModal/DemoGuide/config';

const Services = {
  Pinpoint: 'Pinpoint',
  Personalize: 'Personalize',
};

export default {
  name: 'DemoGuideBadge',
  props: {
    article: { type: String, required: true },
    hideTextOnSmallScreens: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    service() {
      switch (this.article) {
        case Articles.SMS_MESSAGING:
        case Articles.PERSONALIZED_EMAILS:
          return Services.Pinpoint;
        case Articles.USER_PERSONALIZATION:
        case Articles.PERSONALIZED_RANKING:
        case Articles.SIMS_RECOMMENDATIONS:
        case Articles.SIMILAR_ITEMS_RECOMMENDATIONS:
        case Articles.ECOMM_CUSTOMERS_WHO_VIEWED_X:
        case Articles.ECOMM_FBT:
        case Articles.ECOMM_POPULAR_BY_PURCHASES:
        case Articles.ECOMM_POPULAR_BY_VIEWS:
        case Articles.ECOMM_RFY:
          return Services.Personalize;
      }

      throw new Error('Invalid article passed to DemoGuideBadge');
    },
    serviceLogo() {
      switch (this.service) {
        case Services.Pinpoint:
          return '/pinpoint.svg';
        case Services.Personalize:
          return '/personalize.svg';
      }

      throw new Error('Invalid article passed to DemoGuideBadge');
    },
    copy() {
      switch (this.article) {
        case Articles.SMS_MESSAGING:
          return 'Learn more about personalized product recommendations via SMS';
        case Articles.USER_PERSONALIZATION:
          return '유저가 구매한 상품 목록을 통해 자동으로 추천된 상품들입니다.';
        case Articles.PERSONALIZED_RANKING:
          return '현재 유저가 관심있어할 상품이 가장 앞에 우선적으로 보여집니다.';
        case Articles.SIMS_RECOMMENDATIONS:
          return '현재 상품과 비슷한 상품이 보여집니다.';
        case Articles.SIMILAR_ITEMS_RECOMMENDATIONS:
          return '현재 상품과 비슷한 상품이 보여집니다.';
      }

      throw new Error('Invalid article passed to DemoGuideBadge');
    },
    poweredByService() {
      switch (this.service) {
        case Services.Pinpoint:
          return 'Amazon Pinpoint';
        case Services.Personalize:
          return 'Amazon Personalize';
      }

      throw new Error('Invalid article passed to DemoGuideBadge');
    },
  },
  methods: {
    ...mapActions(['demoGuideBadgeClicked']),
    onClick() {
      // this.demoGuideBadgeClicked(this.article);
    },
  },
};
</script>

<style scoped>
.demo-guide-badge {
  border: none;
  background: none;
  display: inline-flex;
}

.logo {
  width: 40px;
}

.text {
  flex: 1;
  font-size: 0.8rem;
  line-height: 1rem;
}

.powered-by {
  font-size: 0.7rem;
}

.service {
  color: var(--blue-600);
}

@media (max-width: 768px) {
  .hide-text-on-small-screens {
    display: none;
  }
}
</style>
