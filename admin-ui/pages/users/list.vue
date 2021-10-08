<template>
  <v-row>
    <v-col cols="12">
      <v-row>
        <v-col cols="12">
          <v-card
            :loading="form.loading"
          >
            <v-card-title>Create invitation</v-card-title>

            <v-card-text>
              <v-form
                ref="createForm"
                v-model="form.valid"
                :disabled="form.loading"
                method="POST"
                @submit.prevent="doSubmitCreate()"
              >
                <v-row>
                  <v-col cols="12" md="8">
                    <v-text-field
                      v-model="form.inputs.comment"
                      label="Comment for invitation (only visible here)"
                    />
                  </v-col>
                  <v-col cols="12" md="4">
                    <v-btn
                      :loading="form.loading"
                      color="success"
                      outlined
                      type="submit"
                    >
                      Create
                    </v-btn>
                  </v-col>
                </v-row>
              </v-form>

              <client-only>
                <v-expand-transition>
                  <v-row v-if="token">
                    <v-col cols="12">
                      <v-text-field
                        :value="registerUrlFor(token)"
                        label="Register URL"
                        readonly
                      />
                    </v-col>
                  </v-row>
                </v-expand-transition>
              </client-only>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>

      <v-row>
        <v-col cols="12">
          <v-card
            :loading="loading"
          >
            <v-card-title>Invitations</v-card-title>

            <v-card-text>
              <v-simple-table>
                <template #default>
                  <thead>
                    <tr>
                      <th class="text-left">
                        Creator
                      </th>

                      <th class="text-left">
                        Comment
                      </th>

                      <th class="text-left">
                        Link
                      </th>

                      <th class="text-left">
                        Used by
                      </th>

                      <th />
                    </tr>
                  </thead>

                  <tbody>
                    <tr
                      v-for="token in tokens"
                      :key="token._id"
                    >
                      <td v-text="`${token.creator.username} (${token.creator.email})`" />
                      <td v-text="token.comment" />
                      <td>
                        <v-text-field
                          readonly
                          :disabled="token.usedUpBy !== null"
                          :value="registerUrlFor(token)"
                        />
                      </td>
                      <td
                        v-if="token.usedUpBy"
                        v-text="`${token.usedUpBy.username} (${token.usedUpBy.email})`"
                      />
                      <td
                        v-else
                      />
                      <td>
                        <v-btn
                          v-if="!token.usedUpBy"
                          :loading="loading"
                          color="error"
                          icon
                          @click.prevent="doDeleteToken(token)"
                        >
                          <v-icon>mdi-delete</v-icon>
                        </v-btn>
                      </td>
                    </tr>
                  </tbody>
                </template>
              </v-simple-table>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>

      <v-row>
        <v-col cols="12">
          <v-card
            :loading="loading"
          >
            <v-card-title>Users</v-card-title>

            <v-card-text>
              <v-simple-table>
                <template #default>
                  <thead>
                    <tr>
                      <th class="text-left">
                        Username
                      </th>

                      <th class="text-left">
                        Email
                      </th>
                    </tr>
                  </thead>

                  <tbody>
                    <tr
                      v-for="user in users"
                      :key="user._id"
                    >
                      <td v-text="user.username" />
                      <td v-text="user.email" />
                    </tr>
                  </tbody>
                </template>
              </v-simple-table>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
    </v-col>
  </v-row>
</template>

<router>
name: PageUsersList
</router>

<script>
  import {
    mapActions,
    mapGetters,
  } from "vuex";

  export default {
    data: () => ({
      loading: false,
      form: {
        inputs: {
          comment: "",
        },
      },
    }),

    async fetch({ store }) {
      await Promise.all([
        store.dispatch("user/list"),
        store.dispatch("userInvitations/list"),
      ]);
    },

    computed: {
      ...mapGetters("user", [
        "users",
      ]),

      ...mapGetters("userInvitations", [
        "tokens",
        "token",
      ]),
    },

    methods: {
      ...mapActions("userInvitations", [
        "create",
        "list",
        "remove",
      ]),

      async doSubmitCreate() {
        try {
          this.form.loading = true;
          const res = await this.create(this.form.inputs);

          if ("success" === res.status) {
            await this.list();
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

      async doDeleteToken(token) {
        if (!window.confirm(`Delete token with comment \`${ token.comment }\` and code \`${ token.code }\``)) {
          return;
        }

        this.loading = true;
        try {
          const res = await this.remove({
            id: token._id,
          });

          if ("success" === res.status) {
            await this.list();
          } else {
            alert(res.message);
          }
        } catch (e) {
          console.error(e);
          alert("Something went wrong");
        } finally {
          this.loading = false;
        }
      },

      registerUrlFor(token) {
        const { href } = this.$router.resolve({
          name: "PageRegister",
          query: {
            code: token.code,
          },
        });

        const origin =
          this.$isServer
            ? ""
            : window.location.origin
        ;

        return `${ origin }${ href }`;
      },
    },
  };
</script>
