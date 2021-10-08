import Vue from "vue";

export const state = () => ({
  token: null,
  tokens: [],
});

export const getters = {
  token(state) {
    return state.token;
  },

  tokens(state) {
    return state.tokens;
  },

  tokenValid(state) {
    return null !== state.token;
  },
};

export const mutations = {
  SET_TOKEN(state, token) {
    Vue.set(state, "token", token);
  },
  SET_TOKENS(state, tokens) {
    Vue.set(state, "tokens", tokens);
  },
};

export const actions = {
  async list(
    {
      commit,
    },
  ) {
    const response = await this.$api.$get(
      "/user-invitation/",
    );

    const { data = null, status } = response || {};

    if ("success" === status) {
      await commit("SET_TOKENS", data);
    } else {
      await commit("SET_TOKENS", []);
    }
  },

  async create(
    {
      commit,
    },
    {
      comment,
    },
  ) {
    const response = await this.$api.$post(
      "/user-invitation/",
      {
        comment,
      },
    );

    const { data = null, status } = response || {};

    if ("success" === status) {
      await commit("SET_TOKEN", data);
    } else {
      await commit("SET_TOKEN", null);
    }

    return response || {};
  },

  async remove(
    _,
    {
      id,
    },
  ) {
    const response = await this.$api.$delete(
      `/user-invitation/${ id }`,
    );

    return response || {};
  },

  async info(
    {
      commit,
    },
    {
      token,
    },
  ) {
    const response = await this.$api.$get(
      `/user-invitation/info/${ token }`,
    );

    const { data = null, status } = response || {};

    if ("success" === status) {
      await commit("SET_TOKEN", data);
    } else {
      await commit("SET_TOKEN", null);
    }
  },
};
