<template>
  <h1>
    Hello, {{ user.username }}!
  </h1>
</template>

<router>
name: PageHome
</router>

<script>
  import {
    mapGetters,
  } from "vuex";

  export default {
    computed: {
      ...mapGetters("auth", [
        "user",
      ]),
    },
  };
</script>
