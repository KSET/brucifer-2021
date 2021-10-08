<template>
  <div :class="$style.container">
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
              <h1>Register</h1>
            </v-col>
          </v-row>

          <v-row>
            <v-col cols="12">
              <v-text-field
                :value="token.by"
                label="Invited by"
                readonly
              />
            </v-col>
          </v-row>

          <v-row>
            <v-col
              cols="12"
            >
              <v-text-field
                v-model="form.inputs.username"
                :rules="rules.username"
                label="Username"
                required
              />
            </v-col>

            <v-col
              cols="12"
            >
              <v-text-field
                v-model="form.inputs.email"
                :rules="rules.email"
                label="Email"
                required
                type="email"
              />
            </v-col>

            <v-col
              cols="12"
            >
              <v-text-field
                v-model="form.inputs.password"
                :rules="rules.password"
                label="Password"
                required
                type="password"
              />
            </v-col>

            <v-col
              cols="12"
            >
              <v-text-field
                v-model="form.inputs.passwordRepeat"
                :rules="rules.passwordRepeat"
                label="Repeat password"
                required
                type="password"
              />
            </v-col>
          </v-row>
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
            Register
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-form>
  </div>
</template>

<router>
name: PageRegister
</router>

<script>
  import {
    mapActions,
    mapGetters,
  } from "vuex";

  export default {
    async middleware({ store, query, error }) {
      if (!query.code) {
        return error({
          statusCode: 403,
        });
      }

      await store.dispatch("userInvitations/info", { token: query.code });

      const tokenValid = store.getters["userInvitations/tokenValid"];
      if (!tokenValid) {
        return error({
          statusCode: 403,
        });
      }

      return true;
    },

    data: () => ({
      form: {
        valid: false,
        inputs: {
          username: "",
          email: "",
          password: "",
          passwordRepeat: "",
        },
        loading: false,
      },
    }),

    computed: {
      ...mapGetters("userInvitations", [
        "token",
      ]),

      rules() {
        return {
          username: [
            (v) => !!v || "Username required",
          ],
          email: [
            (v) => !!v || "Email required",
          ],
          password: [
            (v) => !!v || "Password required",
          ],
          passwordRepeat: [
            (v) => !!v || "Password required",
            (v) => (!!v && v === this.form.inputs.password) || "Passwords must match",
          ],
        };
      },
    },

    methods: {
      ...mapActions("auth", [
        "register",
      ]),

      async doSubmit() {
        this.form.loading = true;
        try {
          const res = await this.register({
            token: this.token.code,
            ...this.form.inputs,
          });

          if ("success" === res.status) {
            this.$store.commit("userInvitations/SET_TOKEN", null);
            return this.$router.push({
              name: "PageLogin",
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

      resetValidation() {
        this.$refs.form.resetValidation();
      },
    },
  };
</script>

<style lang="scss" module>
  .container {
    max-width: 400px;
    margin: 0 auto;
  }
</style>
