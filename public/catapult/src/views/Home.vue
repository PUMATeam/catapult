<template>
    <h1>Home</h1>
</template>

<script lang="ts">
    import Home from '../views/Home';

    const router = new VueRouter({
    routes: [
        {
        path: '/',
        name: 'Home',
        component: Home
        }
    ]
    })

</script>
