<template>
  <div :class="{ mobile: isMobile }">
    <h1 class="confirm-shopper mb-4 text-center">사용자 바꿔보기</h1>

    <div class="flex-container mb-5 d-flex">
      <div v-if="selection" class="your-selections px-4 py-3">
        <h2 class="your-selections-heading mb-4 text-center">
          다음 기준에 따라 사용자를 찾았습니다.
        </h2>
        <dl class="selections">
          <dt class="key">나이:</dt>
          <dd>{{ selection.ageRange }}대</dd>

          <dt class="key">성별:</dt>
          <dd>
            {{
                profile.gender
            }}
          </dd>
        </dl>
      </div>

      <div class="assigned-shopper px-4 py-3">
        <h2 class="assigned-shopper-heading mb-4 text-center">
          발견한 사용자
        </h2>

        <div class="mb-3">{{ profile.name }}. {{ profile.gender }}. {{ profile.age }}살</div>

        <dl>
          <div class="mb-2 d-flex flex-wrap">
            <dt class="mr-1">아이디:</dt>
            <dd class="mb-0">{{ profile.username }}</dd>
          </div>

        </dl>
      </div>
    </div>

    <div class="button-container d-flex">
      <button class="try-again btn btn-outline-primary" @click="tryAgain">다시 찾기</button>
      <button class="confirm btn btn-primary" @click="confirmShopper">선택하기</button>
    </div>
  </div>
</template>

<script>
import { mapState, mapActions } from 'vuex';
import { EventBus } from '@/event-bus';

import { RepositoryFactory } from '@/repositories/RepositoryFactory';
// import { AnalyticsHandler } from '@/analytics/AnalyticsHandler';

const UsersRepository = RepositoryFactory.get('users');

export default {
  name: 'ConfirmShopper',
  props: {
    selection: { type: Object, required: false },
    assignedShopper: { type: Object, required: true },
  },
  computed: {
    ...mapState({
      isMobile: (state) => state.modal.isMobile,
    }),
    profile() {
      return {
        age: this.assignedShopper.age,
        gender: this.assignedShopper.gender === 'M' ? '남성' : '여성',
        name: `${this.assignedShopper.last_name}${this.assignedShopper.first_name}`,
        segment: this.assignedShopper.segment,
        email: this.assignedShopper.email,
        username: this.assignedShopper.username,
      };
    },
  },
  mounted() {
    // eslint-disable-next-line no-undef
    $(this.$refs.learnMore).tooltip();
  },
  beforeDestroy() {
    // eslint-disable-next-line no-undef
    $(this.$refs.learnMore).tooltip('dispose');
  },
  methods: {
    ...mapActions(['setUser']),
    tryAgain() {
      this.$emit('tryAgain');
    },
    async confirmShopper() {
      await UsersRepository.claimUser(this.assignedShopper.id);

      this.setUser(this.assignedShopper);

      // AnalyticsHandler.identify(this.user);

      EventBus.$emit('authState', 'profileChanged');

      this.$emit('confirm');
    },
  },
};
</script>

<style scoped>
.confirm-shopper {
  font-size: 1.75rem;
}

.mobile .flex-container {
  flex-direction: column;
}

.your-selections,
.assigned-shopper {
  flex: 1;
  border: 1px solid var(--grey-400);
  border-radius: 4px;
}

.your-selections {
  margin-right: 24px;
  font-size: 1.35rem;
}

.your-selections-heading {
  font-size: 1.75rem;
}

.mobile .your-selections {
  margin-right: 0;
  margin-bottom: 24px;
}

.selections {
  display: grid;
  grid-template-columns: 1fr 1fr;
  grid-row-gap: 8px;
  grid-column-gap: 16px;
}

.key {
  justify-self: end;
  font-weight: normal;
  color: var(--grey-600);
}

.assigned-shopper {
  font-size: 1.35rem;
}

.assigned-shopper-heading {
  font-size: 1.75rem;
}

.learn-more {
  font-style: italic;
  font-size: 0.9rem;
  color: var(--blue-600);
}

.button-container {
  justify-content: flex-end;
}

.mobile .button-container {
  flex-direction: column;
  align-items: center;
}

.try-again,
.confirm {
  width: 200px;
  font-size: 1.25rem;
}

.try-again:hover,
.try-again:focus,
.confirm:hover,
.confirm:focus {
  border-color: var(--blue-600);
  background: var(--blue-600);
  color: var(--white);
}

.mobile .try-again,
.mobile .confirm {
  width: 100%;
  max-width: 350px;
}

.try-again {
  margin-right: 8px;
  border-color: var(--blue-500);
  color: var(--blue-500);
}

.mobile .try-again {
  margin-right: 0;
  margin-bottom: 8px;
}

.confirm {
  border-color: var(--blue-500);
  background: var(--blue-500);
}
</style>
