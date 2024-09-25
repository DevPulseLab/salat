<script setup>
import Login from '@/components/Login.vue'
</script>

<template>
    <login />
    <router-link :to="{name: 'register'}">Register</router-link>
</template>