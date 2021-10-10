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
            <h1>Create sponsor</h1>
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
            <v-text-field
              v-model="form.inputs.link"
              :rules="form.rules.link"
              label="URL"
              required
              type="url"
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
            <v-col
              cols="12"
              lg="3"
              md="4"
              sm="6"
            >
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
          Create
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-form>
</template>

<router>
name: PageSponsorsCreate
</router>

<script>
  import {
    mapActions,
  } from "vuex";

  export default {
    data: () => ({
      logo: {
        src: null,
        loading: false,
      },
      form: {
        valid: false,
        inputs: {
          name: "",
          link: "",
          logo: null,
        },
        rules: {
          username: [
            (v) => !!v || "Name required",
          ],
          link: [
            (v) => !!v || "Link required",
          ],
          logo: [
            (v) => !!v || "Logo required",
          ],
        },
        loading: false,
      },
    }),

    watch: {
      "form.inputs.logo"(logo) {
        this.logo.src = null;

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
      ...mapActions("sponsor", [
        "create",
      ]),

      async doSubmit() {
        const { inputs } = this.form;

        this.form.loading = true;
        try {
          const res = await this.create(inputs);

          if ("success" === res.status) {
            return this.$router.push({
              name: "PageSponsorsList",
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
