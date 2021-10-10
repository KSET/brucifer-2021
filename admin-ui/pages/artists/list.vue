<template>
  <v-row>
    <v-col
      v-for="(artist, i) in artists"
      :key="artist.id"
      cols="12"
      sm="6"
      lg="3"
      md="4"
    >
      <v-card>
        <v-img
          :lazy-src="artist.logo[0].url"
          :src="artist.logo[artist.logo.length - 1].url"
          aspect-ratio="1"
          contain
        />

        <v-card-title v-text="artist.name" />

        <v-card-actions>
          <v-btn
            :disabled="i === 0"
            :loading="loading"
            icon
            @click.prevent="moveLeft(artist)"
          >
            <v-icon>mdi-arrow-left</v-icon>
          </v-btn>

          <v-spacer />

          <v-btn
            :loading="loading"
            color="error"
            icon
            @click="doDelete(artist)"
          >
            <v-icon>mdi-delete</v-icon>
          </v-btn>

          <v-btn
            :loading="loading"
            :to="{ name: 'PageArtistsEdit', params: { id: artist.id } }"
            color="warning"
            icon
          >
            <v-icon>mdi-pencil</v-icon>
          </v-btn>

          <v-spacer />

          <v-btn
            :disabled="i === artists.length - 1"
            :loading="loading"
            icon
            @click.prevent="moveRight(artist)"
          >
            <v-icon>mdi-arrow-right</v-icon>
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-col>
  </v-row>
</template>

<router>
name: PageArtistsList
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
      await store.dispatch("artist/listAll");
    },

    computed: {
      ...mapGetters("artist", [
        "artists",
      ]),
    },

    methods: {
      ...mapActions("artist", [
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

      async doDelete(artist) {
        if (!window.confirm(`Delete ${ artist.name }?`)) {
          return;
        }

        this.loading = true;
        try {
          const { status, message } = await this.remove(artist);

          if ("success" !== status) {
            alert(message);
          }
        } finally {
          this.loading = false;
        }
      },

      async moveRight(artist) {
        const { artists } = this;
        const i = artists.findIndex((s) => s === artist);

        await this.doSwap(artists[i], artists[i + 1]);
      },

      async moveLeft(artist) {
        const { artists } = this;
        const i = artists.findIndex((s) => s === artist);

        await this.doSwap(artists[i], artists[i - 1]);
      },
    },
  };
</script>
