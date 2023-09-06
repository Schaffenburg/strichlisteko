import { createStore } from 'vuex'
import User from "@/source/models/User";
import Product from "@/source/models/Product";

export default createStore({
  state: {
    user:new User(-1,"Select user"),
    basket: []
  },
  getters: {
  },
  mutations: {
    addBasket(state,product:Product){
      // @ts-ignore
      state.basket.push(product);
    },
    setUser(state,user:User){
      state.user = user;
    }
  },
  actions: {
  },
  modules: {
  }
})
