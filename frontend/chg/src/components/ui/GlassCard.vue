<template>
    <div class="glass-card" :class="{ 'hover-effect': hover, [className]: className }">
        <div class="glass-card-inner">
            <slot></slot>
        </div>
    </div>
</template>

<script setup lang="ts">
defineProps({
    hover: {
        type: Boolean,
        default: true
    },
    className: {
        type: String,
        default: ''
    }
})
</script>

<style scoped>
.glass-card {
    background: rgba(255, 255, 255, 0.7);
    backdrop-filter: blur(10px);
    -webkit-backdrop-filter: blur(10px);
    border-radius: 12px;
    border: 1px solid rgba(255, 255, 255, 0.5);
    box-shadow: 0 4px 15px rgba(0, 0, 0, 0.04);
    transition: all 0.3s ease;
    overflow: hidden;
    position: relative;
}

.glass-card::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 3px;
    background: linear-gradient(90deg, rgba(100, 212, 135, 0.7), rgba(150, 230, 170, 0.3));
    border-radius: 3px 3px 0 0;
}

.glass-card-inner {
    padding: 20px;
}

.hover-effect:hover {
    transform: translateY(-3px);
    box-shadow: 0 6px 20px rgba(0, 0, 0, 0.08);
    border-color: rgba(100, 212, 135, 0.3);
}
</style>