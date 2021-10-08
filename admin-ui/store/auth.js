import Vue from "vue";

export const state = () => ({
  user: null,
});

export const getters = {
  loggedIn(state) {
    return null !== state.user;
  },
  user(state) {
    return state.user;
  },
};

export const mutations = {
  SET_USER(state, user) {
    Vue.set(state, "user", user);
  },
};

export const actions = {
  async login(
    {
      commit,
    },
    {
      username,
      password,
    },
  ) {
    try {
      const response = await this.$api.$post(
        "/auth/login",
        {
          identity: username,
          password,
        },
      );

      const { data = {} } = response || {};

      await commit("SET_USER", data);

      return response;
    } catch (e) {
      console.log(e);

      return null;
    }
  },

  async logout(
    {
      commit,
    },
  ) {
    try {
      await this.$api.$post(
        "/auth/logout",
      );
    } catch {
    }

    await commit("SET_USER", null);
  },

  async nuxtServerInit(
    {
      commit,
    },
  ) {
    const response = await this.$api.$get(
      "/auth/user",
    );

    const { data = null, status } = response || {};

    if ("success" === status) {
      await commit("SET_USER", data);
    } else {
      await commit("SET_USER", null);
    }
  },
};
