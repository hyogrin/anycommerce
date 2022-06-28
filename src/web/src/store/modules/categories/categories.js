import { RepositoryFactory } from '@/repositories/RepositoryFactory';
// import { capitalize } from '@/util/capitalize';
import { I18n } from '@/i18n';

const ProductsRepository = RepositoryFactory.get('products');

export const categories = {
  state: () => ({ categories: null }),
  getters: {
    formattedCategories: (state) => state.categories?.map(({ name }) => I18n['categories'][name]),
  },
  mutations: {
    setCategories: (state, newCategories) => (state.categories = newCategories),
  },
  actions: {
    getCategories: async ({ commit }) => {
      commit('setCategories', null);

      const { data: categories } = await ProductsRepository.getCategories();

      categories.sort(function (a, b) {
        // Compare the 2 dates
        if (a.id < b.id) return -1;
        if (a.id > b.id) return 1;
        return 0;
      });


      commit('setCategories', categories);
    },
  },
};
