<template>
    <div class="gradient-background" :class="type">
        <slot></slot>
    </div>
</template>

<script setup lang="ts">
defineProps({
    type: {
        type: String,
        default: 'primary',
        validator: (value: string) => ['primary', 'secondary', 'light', 'dark'].includes(value)
    }
})
</script>

<style scoped>
.gradient-background {
    position: relative;
    width: 100%;
    height: 100%;
    overflow: hidden;
}

.gradient-background::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    z-index: -1;
}

.primary::before {
    background: linear-gradient(135deg, #f9f9f9, #ffffff);
}

.secondary::before {
    background: linear-gradient(145deg, #f5fbf7, #e8f5e9);
}

.light::before {
    background: linear-gradient(to right, #ffffff, #fafafa);
}

.dark::before {
    background: linear-gradient(145deg, #64d487, #43b16a);
}

/* 添加微妙的装饰效果 */
.gradient-background::after {
    content: '';
    position: absolute;
    top: 10px;
    right: 10px;
    width: 40px;
    height: 40px;
    border-radius: 50%;
    background: linear-gradient(135deg, rgba(100, 212, 135, 0.3), rgba(100, 212, 135, 0.1));
    z-index: -1;
    box-shadow: 
        80px 40px 100px rgba(100, 212, 135, 0.1),
        -60px 150px 80px rgba(100, 212, 135, 0.05);
}
</style>