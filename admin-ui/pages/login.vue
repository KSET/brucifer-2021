<template>
  <div :class="$style.container">
    <v-form
      ref="form"
      v-model="form.valid"
      :disabled="form.loading"
      method="POST"
      @submit.prevent="doLogin"
    >
      <v-card>
        <v-card-text>
          <v-row>
            <v-col>
              <h1>Login</h1>
            </v-col>
          </v-row>

          <v-row>
            <v-col
              cols="12"
            >
              <v-text-field
                v-model="form.inputs.username"
                :rules="form.rules.username"
                label="Username"
                required
              />
            </v-col>

            <v-col
              cols="12"
            >
              <v-text-field
                v-model="form.inputs.password"
                :rules="form.rules.password"
                label="Password"
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
            Login
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-form>
  </div>
</template>

<router>
name: PageLogin
</router>

<script>
  import {
    mapActions,
  } from "vuex";

  export default {
    data: () => ({
      form: {
        valid: false,
        inputs: {
          username: "",
          password: "",
        },
        rules: {
          username: [
            (v) => !!v || "Username required",
          ],
          password: [
            (v) => !!v || "Password required",
          ],
        },
        loading: false,
      },
    }),

    methods: {
      ...mapActions("auth", {
        login: "login",
      }),

      async doLogin() {
        this.form.loading = true;
        try {
          const res = await this.login(this.form.inputs);

          if ("success" === res.status) {
            return this.$router.push({
              name: "PageHome",
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
