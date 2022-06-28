<template>
  <div>
    <div class="banner mt-2 mb-3 py-2">
      <div class="container d-flex justify-content-center align-items-center">
        <a href="/" aria-label="home"><img src="/logo-footer.png" width="320" /></a>
        <!-- <PoweredByAWS class="powered-by-logo"></PoweredByAWS> -->
      </div>
    </div>

    <slot></slot>

    <footer class="mt-5 container">
      <!-- <TermsAndConditions></TermsAndConditions> -->
    </footer>
    <AppModal></AppModal>

    <ConfirmationModal></ConfirmationModal>
  </div>

</template>

<script>
// import PoweredByAWS from '@/components/PoweredByAWS/PoweredByAWS';
// import TermsAndConditions from '@/partials/TermsAndConditions/TermsAndConditions.vue';
import AppModal from '@/partials/AppModal/AppModal';
// import DemoGuideButton from '@/partials/DemoGuideButton/DemoGuideButton';
import ConfirmationModal from '@/partials/ConfirmationModal/ConfirmationModal';

export default {
  name: 'SecondaryLayout',
  components: {
    AppModal,
    // DemoGuideButton,
    ConfirmationModal
  },
};
</script>

<style scoped>
.banner {
  background: var(white);
}

.logo {
  width: 120px;
}

.powered-by-logo {
  /* fine-tune spacing to make logo appear more centered */
  margin-top: 15px;
  height: 45px;
  fill: var(--white);
}
</style>
