<template>
  <v-row>
    <v-col
      v-for="(sponsor, i) in sponsors"
      :key="sponsor.id"
      cols="1"
      lg="3"
      md="4"
    >
      <v-card>
        <v-img
          :lazy-src="sponsor.logo[0].url"
          :src="sponsor.logo[sponsor.logo.length - 1].url"
          aspect-ratio="1"
          contain
        />

        <v-card-title v-text="sponsor.name" />

        <v-card-actions>
          <v-btn
            :disabled="i === 0"
            :loading="loading"
            icon
            @click.prevent="moveLeft(sponsor)"
          >
            &larr;
          </v-btn>

          <v-spacer />

          <v-btn
            :loading="loading"
            color="error"
            icon
            @click="doDelete(sponsor)"
          >
            <v-icon>mdi-delete</v-icon>
          </v-btn>

          <v-btn
            :loading="loading"
            color="warning"
            icon
            nuxt
            :to="{ name: 'PageSponsorsEdit', params: { id: sponsor.id } }"
          >
            <v-icon>mdi-pencil</v-icon>
          </v-btn>

          <v-spacer />

          <v-btn
            :disabled="i === sponsors.length - 1"
            :loading="loading"
            icon
            @click.prevent="moveRight(sponsor)"
          >
            &rarr;
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-col>
  </v-row>
</template>

<router>
name: PageSponsorsList
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
      await store.dispatch("sponsor/listAll");
    },

    computed: {
      ...mapGetters("sponsor", [
        "sponsors",
      ]),
    },

    methods: {
      ...mapActions("sponsor", [
        "swap",
        "remove",
      ]),

      async doSwap(a, b) {
        this.loading = true;
        try {
          const { status, message } = await this.swap({
            a,
            b,
          });

          if ("success" !== status) {
            alert(message);
          }
        } finally {
          this.loading = false;
        }
      },

      async doDelete(sponsor) {
        if (!window.confirm(`Delete ${ sponsor.name }?`)) {
          return;
        }

        this.loading = true;
        try {
          const { status, message } = await this.remove(sponsor);

          if ("success" !== status) {
            alert(message);
          }
        } finally {
          this.loading = false;
        }
      },

      async moveRight(sponsor) {
        const { sponsors } = this;
        const i = sponsors.findIndex((s) => s === sponsor);

        await this.doSwap(sponsors[i], sponsors[i + 1]);
      },

      async moveLeft(sponsor) {
        const { sponsors } = this;
        const i = sponsors.findIndex((s) => s === sponsor);

        await this.doSwap(sponsors[i], sponsors[i - 1]);
      },
    },
  };
</script>
