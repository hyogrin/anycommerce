<template>
  <Layout>
    <div class="content">

      <!-- Loading Indicator -->
      <div class="container mb-4" v-if="!orders">
        <i class="fas fa-spinner fa-spin fa-3x"></i>
      </div>

      <div class="container">
        <h1>주문 내역</h1>

        <table class="table" v-if="orders">
          <tr>
            <th>주문번호</th>
            <th>아이디</th>
            <th>세부정보</th>
            <th>합계</th>
          </tr>
          <tr v-for="order in orders" v-bind:key="order.id">
            <td>{{ order.id }}</td>
            <td>{{ order.username }}</td>
            <td>
              <div v-for="item in order.items" v-bind:key="item.product_id">
                {{ item.product_name }}: {{ item.quantity }}개
              </div>
            </td>
            <td>{{ formatPrice(order.total) }}</td>
          </tr>
        </table>

        <div class="alert alert-secondary no-orders" v-if="!orders || orders.length == 0">아직 아무 상품도 주문하지 않았습니다.</div>

      </div>
    </div>
  </Layout>
</template>

<script>
import { mapState } from 'vuex'

import { RepositoryFactory } from '@/repositories/RepositoryFactory'
import { formatPrice } from "@/util/formatPrice";

import Layout from '@/components/Layout/Layout'

const OrdersRepository = RepositoryFactory.get('orders')

export default {
  name: 'Orders',
  components: {
    Layout,
  },
  data() {
    return {
      errors: [],
      orders: null
    }
  },
  created() {
    this.getOrders()
  },
  methods: {
    async getOrders() {
      this.orders = null;

      const { data } = await OrdersRepository.getOrdersByUsername(this.user.username)

      this.orders = data
    },
    formatPrice
  },
  computed: {
    ...mapState({ user: state => state.user })
  },
  watch: {
    user() {
      this.getOrders()
    }
  }
}
</script>

<style scoped>
.no-orders {
  margin-bottom: 150px;
}
</style>