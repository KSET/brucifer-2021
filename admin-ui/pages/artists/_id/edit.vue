<template>
  <v-form
    ref="form"
    v-model="form.valid"
    :disabled="form.loading"
    method="POST"
    @submit.prevent="doSubmit"
  >
    <v-card>
      <v-card-text>
        <v-row>
          <v-col>
            <h1>Edit artist</h1>
          </v-col>
        </v-row>

        <v-row>
          <v-col
            cols="12"
          >
            <v-text-field
              v-model="form.inputs.name"
              :rules="form.rules.name"
              label="Name"
              required
            />
          </v-col>

          <v-col
            cols="12"
          >
            <v-file-input
              v-model="form.inputs.logo"
              :loading="logo.loading"
              :rules="form.rules.logo"
              accept="image/*"
              label="Logo"
              required
              show-size
            />
          </v-col>
        </v-row>

        <v-expand-transition>
          <v-row v-if="logo.src">
            <v-col cols="1" md="4" lg="3">
              <v-card>
                <v-img
                  :src="logo.src"
                  aspect-ratio="1"
                  contain
                />
              </v-card>
            </v-col>
          </v-row>
        </v-expand-transition>
      </v-card-text>

      <v-card-actions>
        <v-spacer />

        <v-btn
          :disabled="!form.valid"
          :loading="form.loading"
          color="success"
          large
          text
          type="submit"
        >
          Update
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-form>
</template>

<router>
name: PageArtistsEdit
</router>

<script>
  import {
    mapActions,
    mapGetters,
  } from "vuex";

  export default {
    data() {
      const artist = this.$store.getters["artist/artist"];

      return ({
        logo: {
          src: artist.logo[artist.logo.length - 1].url,
          loading: false,
        },
        form: {
          valid: false,
          inputs: {
            name: String(artist.name),
            logo: null,
          },
          rules: {
            username: [
              (v) => !!v || "Name required",
            ],
            logo: [],
          },
          loading: false,
        },
      });
    },

    computed: {
      ...mapGetters("artist", [
        "artist",
      ]),

      defaultUrl() {
        const { logo } = this.artist;

        return logo[logo.length - 1].url;
      },
    },

    watch: {
      "form.inputs.logo"(logo) {
        this.logo.src = this.defaultUrl;

        if (!logo) {
          return;
        }

        this.logo.loading = true;
        const reader = new FileReader();
        reader.readAsDataURL(logo);
        reader.onload = () => {
          this.logo.src = reader.result;
          this.logo.loading = false;
        };
        reader.onabort = () => {
          this.logo.loading = false;
        };
        reader.onerror = reader.onabort;
      },
    },

    methods: {
      ...mapActions("artist", [
        "update",
      ]),

      async doSubmit() {
        const { inputs } = this.form;

        this.form.loading = true;
        try {
          const res = await this.update({
            id: this.artist.id,
            ...inputs,
          });

          if ("success" === res.status) {
            return this.$router.push({
              name: "PageArtistsList",
            });
          } else {
            alert(res.message);
          }
        } catch (e) {
          console.error(e);
          alert("Something went wrong");
        } finally {
          this.form.loading = false;
        }
      },
    },
  };
</script>
