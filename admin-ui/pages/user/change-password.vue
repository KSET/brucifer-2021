<template>
  <v-row>
    <v-col cols="12">
      <v-card>
        <v-card-title>
          Change password
        </v-card-title>

        <v-card-text>
          <v-form
            v-model="form.valid"
            :disabled="form.loading"
            method="POST"
            @submit.prevent="doSubmit"
          >
            <v-row>
              <v-col cols="12">
                <v-text-field
                  v-model="form.inputs.oldPassword"
                  :rules="rules.oldPassword"
                  label="Old password"
                  type="password"
                />
              </v-col>
              <v-col cols="12">
                <v-text-field
                  v-model="form.inputs.password"
                  :rules="rules.password"
                  label="New password"
                  type="password"
                />
              </v-col>
              <v-col cols="12">
                <v-text-field
                  v-model="form.inputs.passwordRepeat"
                  :rules="rules.passwordRepeat"
                  label="Repeat new password"
                  type="password"
                />
              </v-col>
            </v-row>

            <v-row>
              <v-col cols="12">
                <v-btn
                  :disabled="!form.valid"
                  :loading="form.loading"
                  color="success"
                  large
                  type="submit"
                >
                  Change
                </v-btn>
              </v-col>
            </v-row>
          </v-form>
        </v-card-text>
      </v-card>
    </v-col>
  </v-row>
</template>

<router>
name: PageUserChangePassword
</router>

<script>
  import {
    mapActions,
  } from "vuex";

  export default {
    data: () => ({
      form: {
        loading: false,
        valid: false,
        inputs: {
          oldPassword: "",
          password: "",
          passwordRepeat: "",
        },
      },
    }),

    computed: {
      rules() {
        return {
          oldPassword: [
            (v) => !!v || "Password required",
          ],
          password: [
            (v) => !!v || "Password required",
            (v) => (!!v && v !== this.form.inputs.oldPassword) || "Password must be different from old password",
          ],
          passwordRepeat: [
            (v) => !!v || "Password required",
            (v) => (!!v && v === this.form.inputs.password) || "Passwords must match",
          ],
        };
      },
    },

    methods: {
      ...mapActions("user", [
        "changePassword",
      ]),

      async doSubmit() {
        this.form.loading = true;
        try {
          const res = await this.changePassword({
            oldPassword: this.form.inputs.oldPassword,
            newPassword: this.form.inputs.password,
          });

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
    },
  };
</script>
