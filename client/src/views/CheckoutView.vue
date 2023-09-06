<template>
  <div class="flex">
    <Container class="bg-lgray-800 mt-20 mx-auto max-w-xl" v-if="cData.product">
      <div class="mb-2 text-xl text-white">
        Basket
      </div>
      <div class="flex gap-4">
        <img src="https://picsum.photos/600" class="aspect-square max-h-16 rounded">
        <div class="text-xl text-white">
          {{ cData.product.name }}
          <div class="text-md text-gray-500">
            {{ cData.product.note }}
          </div>
        </div>

        <div class="text-white">
          <div class="pt-2 flex justify-end text-2xl font-semibold">
            {{ cData.product.price.toFixed(2) }}€
          </div>
        </div>
      </div>
    </Container>
  </div>
  <div class="w-full px-6">
    <Container class="bg-lgray-800 mx-auto max-w-5xl mt-4">
      <transition appear name="fade" mode="out-in">
        <div v-if="cData.loading" class="flex items-center justify-center w-full">
          <Spinner/>
        </div>
        <transition appear name="fade" mode="out-in" v-else>
          <div class="text-white" v-if="cData.flowStep===0">
            <div class="text-center text-2xl font-semibold">
              Scan Card
            </div>
            <div class="text-center text-6xl">
              <i class="ri-base-station-line"></i>
            </div>

            <div class="flex w-full items-center">
              <div class="bg-lgray-400 h-0.5 flex-grow"></div>
              <div class="mx-4 mb-1 text-gray-500">
                or
              </div>
              <div class="bg-lgray-400 h-0.5 flex-grow"></div>
            </div>
            <div class="text-center text-2xl font-semibold mb-4">
              Select User
            </div>
            <div class="grid grid-cols-2 sm:grid-cols-4 md:grid-cols-5 lg:grid-cols-6">
              <div v-for="(user) in cData.users" v-bind:key="user.id" @click="selectUser(user)"
                   class="relative rounded overflow-hidden hover:shadow-lg hover:shadow-bunker-400 transition-all active:scale-[0.98] p-0">
                <img class="max-h-60 w-full h-full object-cover"
                     src="https://picsum.photos/600">
                <div class="absolute bottom-0 w-full backdrop-blur bg-opacity-40 bg-lgray-800 p-2">
                  <div class="flex text-white select-none">
                    <div class="text-lg">
                      {{ user.username }}
                    </div>
                    <div class="ml-auto">
                      <div class="pt-2 flex justify-end text-lg font-semibold"
                           :class="(user.getBalance()<0)?'text-red-500':'text-green-400'">
                        {{ user.getBalance().toFixed(2) }}€
                      </div>
                    </div>

                  </div>
                </div>
              </div>
            </div>
          </div>
          <div class="text-white" v-else-if="cData.flowStep===1">
            <div class="text-white text-2xl mb-4">
              Select payment method:
            </div>
            <Container
                class="bg-lgray-700 hover:cursor-pointer hover:bg-lgray-600 active:bg-lgray-500 active:scale-[0.98] transition-all"
                @click="checkout(true)">
              <div class="text-center text-2xl font-semibold">
                User Balance
              </div>
              <div class="text-center text-2xl font-semibold text-gray-500">
                (<span
                  :class="(cData.selectedUser.getBalance()<0)?'text-red-500':'text-green-400'">{{ cData.selectedUser.getBalance().toFixed(2) }}</span>)
              </div>
              <div class="text-center text-6xl">
                <i class="ri-user-3-line"></i>
              </div>
            </Container>
            <div class="flex w-full items-center">
              <div class="bg-lgray-400 h-0.5 flex-grow"></div>
              <div class="mx-4 mb-1 text-gray-500">
                or
              </div>
              <div class="bg-lgray-400 h-0.5 flex-grow"></div>
            </div>
            <Container
                class="bg-lgray-700 hover:cursor-pointer hover:bg-lgray-600 active:bg-lgray-500 active:scale-[0.98] transition-all"
                @click="checkout(false)">
              <div class="text-center text-2xl font-semibold mb-4">
                Cash
              </div>
              <div class="text-center text-6xl">
                <i class="ri-currency-line"></i>
              </div>
            </Container>
          </div>
          <div class="" v-else-if="cData.flowStep===2">
            <div class="text-center text-2xl font-semibold text-green-500">
              Purchase successful
            </div>
            <div class="flex w-full justify-center my-6 text-green-500">
              <i class="ri-check-double-line text-8xl"></i>
            </div>
            <div class="text-center text-xl font-semibold text-gray-300">
              Thank you for your purchase!
            </div>
            <div class="text-center text-sm text-gray-500">
              Click anywhere to return home
            </div>
          </div>
        </transition>
      </transition>
    </Container>
  </div>
</template>

<script setup lang="ts">

import {useStore} from "vuex";
import Container from "@/components/Container.vue";
import {onMounted, onUnmounted, reactive} from "vue";
import Product from "@/source/models/Product";
import User from "@/source/models/User";
import Spinner from "@/components/Spinner.vue";
import {useRouter} from "vue-router";

const router = useRouter();
const store = useStore()
const cData: {
  product: Product | undefined,
  flowStep: number,
  users: User[] | undefined,
  loading: boolean,
  selectedUser: User | undefined,
} = reactive({
  flowStep: 0,
  product: undefined,
  users: undefined,
  loading: true,
  selectedUser: undefined,
})
if (store.state.basket.length > 0) {
  cData.product = store.state.basket[0];
}


onMounted(() => {
  loadUsers();
  document.getElementsByTagName("body")[0].addEventListener("click", homeRedirect);
});

onUnmounted(() => {
  document.getElementsByTagName("body")[0].removeEventListener("click", homeRedirect);
})

async function loadUsers() {
  cData.loading = true;
  try {
    cData.users = await User.all();
  } catch (e) {
    console.log(e)
  }
  cData.loading = false;

}

function selectUser(user: User) {
  cData.selectedUser = user;
  store.commit("setUser", user);
  cData.flowStep = 1;
}

async function checkout(account: boolean) {
  //TODO send stuff to api
  if (account) {
    //a
  } else {
//a
  }
  setTimeout(() => {
    cData.flowStep = 2;
  }, 200)

}

function homeRedirect() {
  if (cData.flowStep === 2) {
    router.push({name: 'home'})
  }
}

</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: all 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
