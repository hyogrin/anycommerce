<template>
  <div :class="{ 'get-started-container': true, mobile: isMobile }">
    <h1 class="heading mb-4 text-center">다른 사용자로 전환하기</h1>

    <div class="explanation mb-5">
      <p>
        애니커머스에서는 수천명의 가상의 사용자가 나이와 성별 등 알고리즘에 따라 구매를 하고 있습니다.
      </p>

      <p>지금 원하는 사용자를 선택해서 로그인없이 해당 가상 사용자로 전환하실 수 있습니다. </p>
      <p>아래 옵션 중 하나를 선택하세요.</p>
    </div>

    <div class="button-container mb-5 d-flex justify-content-center">
      <button type="button" class="auto-select btn btn-lg btn-outline-primary" @click="autoSelectShopper"
        data-toggle="tooltip" data-placement="bottom" title="무작위 사용자를 선택합니다." ref="autoSelectShopper">
        무작위로 선택하기
      </button>
      <button type="button" class="choose-shopper btn btn-lg btn-primary" @click="chooseAShopper" data-toggle="tooltip"
        data-placement="bottom" title="특정 연령과 성별을 직접 선택합니다." ref="chooseAShopper">
        직접 선택하기
      </button>
    </div>

    <hr class="mb-5" />

  </div>
</template>

<script>
import { mapState, mapActions } from 'vuex';
// import { AnalyticsHandler } from '@/analytics/AnalyticsHandler';
import { EventBus } from '@/event-bus';

import { RepositoryFactory } from '@/repositories/RepositoryFactory';

const UsersRepository = RepositoryFactory.get('users');

export default {
  name: 'GetStarted',
  computed: {
    ...mapState({ isMobile: (state) => state.modal.isMobile }),
  },
  mounted() {
    // eslint-disable-next-line no-undef
    $([this.$refs.autoSelectShopper, this.$refs.chooseAShopper]).tooltip();
  },
  beforeDestroy() {
    // eslint-disable-next-line no-undef
    $([this.$refs.autoSelectShopper, this.$refs.chooseAShopper]).tooltip('dispose');
  },
  methods: {
    ...mapActions(['setUser']),
    chooseAShopper() {
      this.$emit('chooseAShopper');
    },
    async autoSelectShopper() {
      const { data } = await UsersRepository.getRandomUser();

      this.$emit('autoSelectShopper', {
        assignedShopper: data[0],
      });
    },
    async useDefaultProfile() {
      const cognitoUser = await this.$Amplify.Auth.currentAuthenticatedUser();

      const { data: user } = await UsersRepository.getUserByUsername(cognitoUser.username);

      this.setUser(user);

      // AnalyticsHandler.identify(user);

      EventBus.$emit('authState', 'profileChanged');

      this.$emit('useDefaultProfile');
    },
  },
};
</script>

<style scoped>
.get-started-container {
  max-width: 800px;
  margin: auto;
}

.heading {
  margin-top: 15%;
  font-size: 1.75rem;
}

.explanation {
  font-size: 1.15rem;
}

.mobile .button-container {
  flex-direction: column;
  align-items: center;
}

.auto-select,
.choose-shopper {
  flex: 1;
}

.mobile .auto-select,
.mobile .choose-shopper {
  width: 100%;
  max-width: 350px;
}

.auto-select {
  margin-right: 16px;
  color: var(--blue-500);
}

.auto-select {
  border-color: var(--blue-500);
}

.mobile .auto-select {
  margin-right: 0px;
  margin-bottom: 16px;
}

.choose-shopper {
  background-color: var(--blue-500);
  border-color: var(--blue-500);
}

.auto-select:hover,
.auto-select:focus,
.choose-shopper:hover,
.choose-shopper:focus {
  background-color: var(--blue-600);
  border-color: var(--blue-600);
  color: var(--white);
}

.default-profile {
  color: var(--blue-600);
}
</style>
