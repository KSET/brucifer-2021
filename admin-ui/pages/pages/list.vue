<template>
  <v-row>
    <v-col
      v-for="page in pages"
      :key="page.id"
      cols="12"
      lg="3"
      md="4"
      sm="6"
    >
      <v-card>
        <v-card-title v-text="page.name" />

        <v-card-actions>
          <v-spacer />

          <v-btn
            :loading="loading"
            color="error"
            icon
            @click="doDelete(page)"
          >
            <v-icon>mdi-delete</v-icon>
          </v-btn>

          <v-btn
            :loading="loading"
            :to="{ name: 'PagePagesEdit', params: { id: page.id } }"
            color="warning"
            icon
            nuxt
          >
            <v-icon>mdi-pencil</v-icon>
          </v-btn>

          <v-spacer />
        </v-card-actions>
      </v-card>
    </v-col>
  </v-row>
</template>

<router>
name: PagePagesList
</router>

<script>
  import {
    mapActions,
    mapGetters,
  } from "vuex";

  export default {
    data: () => ({
      loading: false,
    }),

    async fetch({ store }) {
      await store.dispatch("page/list");
    },

    computed: {
      ...mapGetters("page", [
        "pages",
      ]),
    },

    methods: {
      ...mapActions("page", [
        "remove",
      ]),

      async doDelete(sponsor) {
        if (!window.confirm(`Delete ${ sponsor.name }?`)) {
          return;
        }

        this.loading = true;
        try {
          const { status, message } = await this.remove({
            id: sponsor.id,
          });

          if ("success" !== status) {
            alert(message);
          }
        } finally {
          this.loading = false;
        }
      },
    },
  };
</script>
