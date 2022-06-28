<template>
  <button class="button btn btn-primary" @click="openDemoGuide">
    <div class="demo-guide">더 알아보기</div>
    <div class="learn-more">이 사이트에 대해서 더 알아보기</div>
  </button>
</template>

<script>
import { mapActions } from 'vuex';
import { Modals } from '@/partials/AppModal/config';

export default {
  name: 'DemoGuideButton',
  methods: {
    ...mapActions(['openModal']),
    openDemoGuide() {
      this.openModal(Modals.DemoGuide);
    },
  },
};
</script>

<style scoped>
.button {
  border: none;
  border-radius: 12px 12px 0 0;
  background: var(--blue-500);
}

.button:hover,
.button:focus {
  background: var(--blue-600);
}

.demo-guide {
  font-weight: bold;
  font-size: 1.4rem;
}

.learn-more {
  font-weight: bold;
  font-size: 1.2rem;
}
</style>
