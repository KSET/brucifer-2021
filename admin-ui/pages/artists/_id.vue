<template>
  <nuxt-child />
</template>

<script>
  export default {
    async middleware({ store, error, params }) {
      const artist = await store.dispatch("artist/info", { id: params.id });

      if (null === artist) {
        return error({
          statusCode: 404,
        });
      }
    },
  };
</script>
