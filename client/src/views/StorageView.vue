<template>
  <div class="w-full px-6">
    <transition name="fade">
      <Container class="items-center justify-center flex border-spacing-2" v-if="pData.loading">
        <Spinner class="text-xl"/>
      </Container>
      <table v-else-if="pData.products!==undefined" class="max-w-6xl mx-auto text-white table-auto w-full">
        <thead>
        <tr class="text-left">
          <th>#</th>
          <th>Image</th>
          <th>Name</th>
          <th># available</th>
          <th>Size</th>
          <th>EAN</th>
          <th>Price</th>
          <th>Box-Size</th>
          <th>Actions</th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="(product) in pData.products" v-bind:key="product.id"
            class="relative rounded overflow-hidden hover:bg-bunker-400 transition-all">
          <td class="px-2">
            {{ product.id }}
          </td>
          <td class="px-2 relative">
            <div
                class="absolute h-full w-full left-0 right-0 bg-bunker transition-all opacity-0 hover:opacity-50 flex items-center justify-center"
                @click="pData.image.selProduct=product">
              <i class="ri-image-add-line text-white"></i>
            </div>
            <img class="max-h-20 w-full h-full object-cover"
                 src="https://picsum.photos/600">
          </td>
          <td class="px-2">
            <input type="text" class="min-w-0 w-full" v-model="product.name" @change="productChanged(product)">
          </td>
          <td class="px-2">
            <input type="text" class="min-w-0 w-full" v-model="product.stock">
          </td>
          <td class="px-2">
            <input type="text" class="min-w-0 w-full" v-model="product.amount">
          </td>
          <td class="px-2">
            <input type="text" class="min-w-0 w-full" v-model="product.ean">
          </td>
          <td class="px-2">
            <div class="flex items-center justify-center">
              <input type="text" class="min-w-0 w-full mr-1" v-model="product.price">
              <div>ct</div>
            </div>
          </td>
          <td class="px-2">
            <input type="text" class="min-w-0 w-full" v-model="product.boxSize">
          </td>
          <td class="px-2 flex gap-2 min-w-[200px]">
            <div class="grid grid-cols-3 gap-2">
              <Button class="bg-red-500 hover:bg-red-600 active:bg-red-700">
                -1
              </Button>
              <Button class="bg-primary-600 hover:bg-primary-700 active:bg-primary-800">
                +1
              </Button>
              <Button class="bg-primary-500 hover:bg-primary-600 active:bg-primary-700">
                +6
              </Button>
              <Button class="bg-primary-400 hover:bg-primary-500 active:bg-primary-600">
                +12
              </Button>
              <Button class="bg-primary-300 hover:bg-primary-400 active:bg-primary-500">
                +20
              </Button>
              <Button class="bg-primary-200 hover:bg-primary-300 active:bg-primary-400">
                +24
              </Button>
            </div>
            <div class="space-y-2">
              <Button class="bg-red-500 hover:bg-red-600 active:bg-red-700">
                <i class="ri-delete-bin-line"></i>
              </Button>
            </div>
          </td>
        </tr>

        <tr class="relative rounded overflow-hidden hover:bg-bunker-400 transition-all">
          <td class="px-2">
          </td>
          <td class="px-2">
            <input type="file" accept="image/jpeg, image/png" class="min-w-0 w-full">
          </td>
          <td class="px-2">
            <input type="text" class="min-w-0 w-full">
          </td>
          <td class="px-2">
            <input type="text" class="min-w-0 w-full">
          </td>
          <td class="px-2">
            <input type="number" step="1" class="min-w-0 w-full">
          </td>
          <td class="px-2">
            <input type="text" class="min-w-0 w-full">
          </td>
          <td class="px-2">
            <input type="number" step="1" class="min-w-0 w-full">
          </td>
          <td class="px-2">
            <input type="number" step="1" class="min-w-0 w-full">
          </td>
          <td class="px-2 flex gap-2">
            <Button class="bg-green-500 hover:bg-green-600 active:bg-green-700 my-3">
              Save
            </Button>
          </td>
        </tr>
        </tbody>
      </table>
      <Container v-else>
        An error occurred.
      </Container>
    </transition>
  </div>
  <div
      class="absolute top-0 left-0 bottom-0 right-0 backdrop-blur bg-bunker bg-opacity-50 flex items-center justify-center"
      v-if="pData.image.selProduct">
    <Container class="">
      <input type="file" accept="image/jpeg, image/png" class="min-w-0 w-full" @change="onImagePicked">
      <Button @click="saveImage(pData.image.selProduct, pData.image.imageData)"
              class="bg-green-500 hover:bg-green-600 active:bg-green-700">Save
      </Button>
    </Container>
  </div>
</template>

<script setup lang="ts">
import {defineComponent, onMounted, reactive} from 'vue';
import Container from "@/components/Container.vue";
import Spinner from "@/components/Spinner.vue";
import Product from "@/source/models/Product";
import {useStore} from "vuex";
import {useRouter} from "vue-router";
import Button from "@/components/Button.vue";

const pData: {
  products: Product[] | undefined,
  loading: boolean,
  image: {
    selProduct: Product | undefined,
    imageData: string
  }
} = reactive({
  products: undefined,
  loading: true,

  image: {
    selProduct: undefined,
    imageData: "",
  }
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


function productChanged(product: Product) {

}

async function saveImage(product: Product, imageData: string) {
  try {
    await product.saveImage(imageData);
  }catch (e){
    console.log(e)
  }
  pData.image.selProduct=undefined
  loadProducts();
}

function onImagePicked(e){
  const files = e.target.files
  let filename = files[0].name
  const fileReader = new FileReader()
  fileReader.addEventListener('load', () => {
    if(fileReader.result!==null)
      //@ts-ignore
      pData.image.imageData = fileReader.result
  })
  fileReader.readAsBinaryString(files[0])
}

</script>

<style scoped>
input {
  @apply rounded bg-lgray-800 border border-lgray focus:border-primary text-white px-2 py-1 my-3
}
</style>
