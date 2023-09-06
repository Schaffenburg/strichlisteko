<template>
  <div class="absolute w-full backdrop-blur bg-lgray-800 bg-opacity-50 p-4 flex px-6 md:px-12">
    <router-link :to="{name:'home'}">
      <NavButton>
        <i class="ri-home-line text-white"></i>
        <div class="text-white ml-1">Home</div>
      </NavButton>
    </router-link>
    <div class="flex-grow"></div>
    <div class="" v-if="store.state.user">
      <div class="flex items-center gap-2">
        <div class="text-lg text-white font-semibold">
          {{store.state.user.username}}
        </div>
        <img src="https://picsum.photos/200" class="aspect-square w-8 rounded-full" alt="profile picture">
      </div>
    </div>
  </div>
  <router-view class="pt-16"/>
</template>

<script setup lang="ts">

import NavButton from "@/components/NavButton.vue";
import {useStore} from "vuex";

const store = useStore();

</script>

<style>
#app{
  @apply min-w-full min-h-screen bg-bunker-600
}
</style>
