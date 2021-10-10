<template>
  <nuxt-child />
</template>

<script>
  export default {
    async middleware({ store, error, params }) {
      const sponsor = await store.dispatch("sponsor/info", { id: params.id });

      if (null === sponsor) {
        return error({
          statusCode: 404,
        });
      }
    },
  };
</script>
