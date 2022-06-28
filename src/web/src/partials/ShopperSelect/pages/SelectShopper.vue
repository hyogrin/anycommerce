<template>
  <div>
    <h1 class="heading">사용자 선택하기</h1>
    <p>
      사용자를 선택해 실제 취향이 어떻게 반영되어 추천되고 있는지 살펴보세요.
    </p>

    <div class="form-container p-4">
      <div class="mb-4">아래 도구를 사용해 연령대와 성별을 선택해 주세요.</div>

      <form @submit.prevent="onSubmit" class="form">
        <div class="form-group d-flex align-items-center">
          <label for="age-range" class="label mr-3 mb-0">1</label>
          <select class="form-control" id="age-range" placeholder="연령대를 선택해 주세요" v-model="ageRange">
            <option value="">연령대를 선택해 주세요</option>
            <option value="10">10대</option>
            <option value="20">20대</option>
            <option value="30">30대</option>
            <option value="40">40대</option>
            <option value="50">50대</option>
          </select>
        </div>

        <div class="form-group d-flex flex-column">
          <div class="d-flex align-items-center">
            <label for="gender" class="label mr-3 mb-0">2</label>
            <select class="form-control" id="gender" placeholder="성별을 선택해 주세요" v-model="gender">
              <option value="">성별을 선택해 주세요</option>
              <option value="F">여성</option>
              <option value="M">남성</option>
            </select>
          </div>

        </div>

        <div class="text-center">
          <button class="submit btn btn-primary btn-lg" type="submit" :disabled="!(ageRange && gender)">
            선택하기
          </button>
        </div>

        <div v-if="shopperNotFound" class="alert alert-warning mt-4" role="alert">
          그런 사용자를 찾지 못했습니다. 다른 연령대와 성별을 선택해 주세요.
        </div>
      </form>
    </div>
  </div>
</template>

<script>
import { RepositoryFactory } from '@/repositories/RepositoryFactory';
import { mapGetters, mapState } from 'vuex';

const UsersRepository = RepositoryFactory.get('users');

export default {
  name: 'SelectShopper',
  data() {
    return {
      ageRange: '',
      gender: '',
      shopperNotFound: false,
    };
  },
  mounted() {
    // eslint-disable-next-line no-undef
    $(this.$refs.learnMore).tooltip();
  },
  beforeDestroy() {
    // eslint-disable-next-line no-undef
    $(this.$refs.learnMore).tooltip('dispose');
  },
  computed: {
    ...mapState({ categories: (state) => state.categories.categories }),
    ...mapGetters(['formattedCategories']),
  },
  methods: {
    async onSubmit() {
      const { gender, ageRange } = this;
      const segment = gender + ageRange;
      const { data } = await UsersRepository.getUnclaimedUser({ segment });

      if (!data) {
        this.shopperNotFound = true;
      } else {
        this.$emit('shopperSelected', {
          selection: { gender, ageRange },
          assignedShopper: data[0],
        });
      }
    },
    resetShopperNotFound() {
      if (this.shopperNotFound) this.shopperNotFound = false;
    },
  },
  watch: {
    gender() {
      this.resetShopperNotFound();
    },
    ageRange() {
      this.resetShopperNotFound();
    },
  },
};
</script>

<style scoped>
.heading {
  font-size: 1.75rem;
}

.form-container {
  border: 1px solid var(--grey-400);
  border-radius: 4px;
}

.form {
  margin: auto;
  max-width: 350px;
}

.label {
  font-size: 2.5rem;
  color: var(--grey-600);
  width: 25px;
}

.learn-more {
  font-style: italic;
  font-size: 0.9rem;
  color: var(--blue-600);
  align-self: flex-end;
}

.submit {
  width: 200px;
  background: var(--blue-500);
  border-color: var(--blue-500);
}

.submit:hover:not([disabled]),
.submit:focus:not([disabled]) {
  background: var(--blue-600);
  border-color: var(--blue-600);
}
</style>
