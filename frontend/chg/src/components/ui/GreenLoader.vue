<template>
    <div class="loader-container" :class="{ 'centered': centered }">
        <div class="loader" :class="size">
            <svg class="loader-circle" viewBox="0 0 50 50">
                <circle class="loader-path" cx="25" cy="25" r="20" fill="none" stroke-width="4"></circle>
            </svg>
        </div>
        <div v-if="text" class="loader-text">{{ text }}</div>
    </div>
</template>

<script setup lang="ts">
defineProps({
    size: {
        type: String,
        default: 'medium',
        validator: (value: string) => ['small', 'medium', 'large'].includes(value)
    },
    text: {
        type: String,
        default: ''
    },
    centered: {
        type: Boolean,
        default: false
    }
})
</script>

<style scoped>
.loader-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
}

.centered {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(255, 255, 255, 0.85);
    backdrop-filter: blur(6px);
    -webkit-backdrop-filter: blur(6px);
    z-index: 9999;
}

.loader {
    position: relative;
}

.small {
    width: 24px;
    height: 24px;
}

.medium {
    width: 36px;
    height: 36px;
}

.large {
    width: 48px;
    height: 48px;
}

.loader-circle {
    animation: rotate 2s linear infinite;
    height: 100%;
    transform-origin: center center;
    width: 100%;
}

.loader-path {
    stroke-dasharray: 150;
    stroke-dashoffset: 0;
    stroke: #64d487;
    stroke-linecap: round;
    animation: dash 1.5s ease-in-out infinite;
}

.loader-text {
    margin-top: 12px;
    color: #666;
    font-size: 14px;
    font-weight: 500;
    letter-spacing: 0.2px;
}

@keyframes rotate {
    100% {
        transform: rotate(360deg);
    }
}

@keyframes dash {
    0% {
        stroke-dasharray: 1, 150;
        stroke-dashoffset: 0;
    }
    50% {
        stroke-dasharray: 90, 150;
        stroke-dashoffset: -35;
    }
    100% {
        stroke-dasharray: 90, 150;
        stroke-dashoffset: -124;
    }
}
</style>