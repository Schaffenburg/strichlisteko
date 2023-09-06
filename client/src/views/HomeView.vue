<template>
  <div class="w-full">
    <transition name="fade">
      <Container class="items-center justify-center flex" v-if="pData.loading">
        <Spinner class="text-xl"/>
      </Container>
      <Container v-else-if="pData.products!==undefined" class="bg-bunker-500 max-w-5xl mx-auto">
        <div class="w-full grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-4">
          <div v-for="(product) in pData.products" v-bind:key="product.id" @click="addBasket(product)"
               class="relative rounded overflow-hidden hover:shadow-lg hover:shadow-bunker-400 transition-all active:scale-[0.98]">
            <img class="max-h-60 w-full h-full object-cover"
                 src="https://picsum.photos/600">
            <div class="absolute bottom-0 w-full backdrop-blur bg-opacity-40 bg-lgray-800 p-2">
              <div class="flex text-white select-none">
                <div class="text-lg">
                  {{product.name}}
                  <div class="text-xs text-gray-500">
                    {{product.note}}
                  </div>
                </div>
                <div class="ml-auto">
                  <div class="flex items-center justify-center" :class="(product.stock===0)?'text-red-500':'text-gray-400'">
                    <i class="ri-store-3-line mr-1"></i>
                    <div class="text-sm tracking-tight">
                      {{product.stock}}
                    </div>
                  </div>
                  <div class="pt-2 flex justify-end text-lg font-semibold">
                    {{product.getPrice().toFixed(2)}}€
                  </div>
                </div>

              </div>
            </div>
          </div>
        </div>
      </Container>
      <Container v-else>
        An error occurred.
      </Container>
    </transition>
  </div>
</template>

<script setup lang="ts">
import {defineComponent, onMounted, reactive} from 'vue';
import Container from "@/components/Container.vue";
import Spinner from "@/components/Spinner.vue";
import Product from "@/source/models/Product";
import {useStore} from "vuex";
import {useRouter} from "vue-router";

const pData: { products: Product[] | undefined, loading: boolean } = reactive({
  products: undefined,
  loading: true,
})

/**
 * Load product data on page mount
 */
onMounted(async () => {
  await loadProducts()
});

async function loadProducts() {
  pData.loading = true;
  try {
    pData.products = await Product.all();
  } catch (e) {
    console.log(e)
  }
  pData.loading = false;
}


const store = useStore();
const router = useRouter();
function addBasket(product:Product){
  store.commit("addBasket", product);
  setTimeout(()=>{
    router.push({name:"checkout"})
  },200)

}

</script>
