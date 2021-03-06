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
  async fetchUser(
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

      const { status = "error", data = {} } = response || {};

      if ("success" === status) {
        await commit("SET_USER", data);
      }

      return response;
    } catch (e) {
      console.log(e);

      return null;
    }
  },

  async register(
    _,
    {
      email,
      username,
      password,
      token,
    },
  ) {
    const response = await this.$api.$post(
      "/auth/register",
      {
        email,
        identity: username,
        password,
        code: token,
      },
    );

    return response;
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
      dispatch,
    },
  ) {
    await dispatch("fetchUser");
  },
};
