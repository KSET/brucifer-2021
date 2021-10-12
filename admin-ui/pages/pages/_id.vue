<template>
  <nuxt-child />
</template>

<script>
  export default {
    async middleware({ store, error, params }) {
      const { status } = await store.dispatch("page/info", { id: params.id });

      if ("success" !== status) {
        return error({
          statusCode: 404,
        });
      }
    },
  };
</script>
