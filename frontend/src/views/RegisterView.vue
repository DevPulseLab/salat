<script setup>
import Register from '@/components/Register.vue'
</script>

<template>
    <Register />
    <router-link :to="{name: 'home'}">Login</router-link>
</template>