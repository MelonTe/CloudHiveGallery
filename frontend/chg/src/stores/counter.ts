import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
/* 全局变量库，例如固定不变的头像 */
//一个状态就存储一类要共享的数据
export const useCounterStore = defineStore('counter', () => {
  //定义状态的初始值
  const count = ref(0)
  const doubleCount = computed(() => count.value * 2)
  //定义触发变量修改的函数
  function increment() {
    count.value++
  }

  return { count, doubleCount, increment }
})
