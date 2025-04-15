<template>
    <a-button :type="type" :class="['green-button', `green-button-${type}`, { 'green-button-rounded': rounded }]"
        v-bind="$attrs">
        <slot></slot>
    </a-button>
</template>

<script setup lang="ts">
defineProps({
    type: {
        type: String,
        default: 'default',
        validator: (value: string) => ['default', 'primary', 'dashed', 'text', 'link'].includes(value)
    },
    rounded: {
        type: Boolean,
        default: false
    }
})
</script>

<style scoped>
.green-button {
    transition: all 0.3s ease;
    position: relative;
    overflow: hidden;
    border-radius: 6px;
}

/* 添加波纹效果 */
.green-button::after {
    content: '';
    position: absolute;
    top: 50%;
    left: 50%;
    width: 5px;
    height: 5px;
    background: rgba(255, 255, 255, 0.7);
    opacity: 0;
    border-radius: 100%;
    transform: scale(1, 1) translate(-50%, -50%);
    transform-origin: 50% 50%;
}

.green-button:active::after {
    animation: ripple 0.6s ease-out;
}

.green-button-rounded {
    border-radius: 50px;
}

.green-button-primary {
    background: #64d487;
    border-color: #64d487;
    color: white;
    box-shadow: 0 2px 0 rgba(0, 0, 0, 0.03);
}

.green-button-primary:hover {
    background: #43b16a;
    border-color: #43b16a;
    transform: translateY(-1px);
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.06);
}

.green-button-default {
    border-color: #e0e0e0;
    color: #333;
}

.green-button-default:hover {
    color: #64d487;
    border-color: #64d487;
    background-color: rgba(100, 212, 135, 0.02);
}

.green-button-dashed {
    border-color: #e0e0e0;
    color: #333;
}

.green-button-dashed:hover {
    color: #64d487;
    border-color: #64d487;
    background-color: rgba(100, 212, 135, 0.02);
}

.green-button-text:hover {
    color: #64d487;
    background-color: rgba(100, 212, 135, 0.05);
}

.green-button-link {
    color: #64d487;
}

.green-button-link:hover {
    color: #43b16a;
}

@keyframes ripple {
    0% {
        transform: scale(0, 0);
        opacity: 1;
    }

    20% {
        transform: scale(25, 25);
        opacity: 1;
    }

    100% {
        opacity: 0;
        transform: scale(40, 40);
    }
}
</style>