<template>
  <v-app dark>
    <v-navigation-drawer
      v-model="drawer"
      :mini-variant="miniVariant"
      app
      clipped
      fixed
    >
      <v-list>
        <v-list-item
          :to="{ name: 'PageHome' }"
          exact
          router
        >
          <v-list-item-action>
            <v-icon>mdi-home</v-icon>
          </v-list-item-action>
          <v-list-item-content>
            <v-list-item-title>
              Home
            </v-list-item-title>
          </v-list-item-content>
        </v-list-item>

        <v-list-item
          v-if="!loggedIn"
          :to="{ name: 'PageLogin' }"
          exact
          router
        >
          <v-list-item-icon>
            <v-icon>mdi-account</v-icon>
          </v-list-item-icon>

          <v-list-item-title>
            Login
          </v-list-item-title>
        </v-list-item>
        <transition-group
          v-else
        >
          <v-list-group
            key="account"
            no-action
            prepend-icon="mdi-account-circle"
          >
            <template #activator>
              <v-list-item-title v-text="user.username" />
            </template>

            <v-list-item
              exact
              href="#"
              link
              @click.prevent="doLogout"
            >
              <v-list-item-icon>
                <v-icon>mdi-account-off</v-icon>
              </v-list-item-icon>

              <v-list-item-title>
                Logout
              </v-list-item-title>
            </v-list-item>
          </v-list-group>

          <v-list-group
            key="sponsors"
            no-action
            prepend-icon="mdi-domain"
          >
            <template #activator>
              <v-list-item-title>Sponsors</v-list-item-title>
            </template>

            <v-list-item
              :to="{ name: 'PageSponsorsList' }"
              exact
              nuxt
            >
              <v-list-item-icon>
                <v-icon>mdi-view-list</v-icon>
              </v-list-item-icon>

              <v-list-item-title>
                List
              </v-list-item-title>
            </v-list-item>

            <v-list-item
              :to="{ name: 'PageSponsorsCreate' }"
              exact
              nuxt
            >
              <v-list-item-icon>
                <v-icon>mdi-playlist-plus</v-icon>
              </v-list-item-icon>

              <v-list-item-title>
                Create
              </v-list-item-title>
            </v-list-item>
          </v-list-group>
        </transition-group>
      </v-list>
    </v-navigation-drawer>

    <v-app-bar
      app
      clipped-left
      fixed
    >
      <v-app-bar-nav-icon @click.stop="drawer = !drawer" />

      <v-btn
        icon
        @click.stop="miniVariant = !miniVariant"
      >
        <v-icon>mdi-{{ `chevron-${ miniVariant ? "right" : "left" }` }}</v-icon>
      </v-btn>

      <v-toolbar-title>
        <a href="/">Brucošijada</a>
      </v-toolbar-title>
    </v-app-bar>

    <v-main>
      <v-container>
        <nuxt />
      </v-container>
    </v-main>
  </v-app>
</template>

<script>
  import {
    mapActions,
    mapGetters,
  } from "vuex";

  export default {
    middleware: [
      "auth",
    ],

    data: () => ({
      drawer: false,
      miniVariant: false,
    }),

    computed: {
      ...mapGetters("auth", [
        "loggedIn",
        "user",
      ]),
    },

    methods: {
      ...mapActions("auth", [
        "logout",
      ]),

      async doLogout() {
        await this.logout();
        await this.$router.push({
          name: "PageLogin",
        });
      },
    },
  };
</script>