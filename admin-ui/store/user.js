import Vue from "vue";

export const state = () => ({
  users: [],
});

export const getters = {
  users(state) {
    return state.users;
  },
};

export const mutations = {
  SET_USERS(state, users) {
    Vue.set(state, "users", users);
  },
};

export const actions = {
  async list(
    {
      commit,
    },
  ) {
    const response = await this.$api.$get(
      "/user/",
    );

    if ("success" === response?.status) {
      commit("SET_USERS", response.data);
    }

    return response;
  },

  async changePassword(
    {
      dispatch,
    },
    {
      oldPassword,
      newPassword,
    },
  ) {
    const response = await this.$api.$patch(
      "/user/",
      {
        password: {
          old: oldPassword,
          new: newPassword,
        },
      },
    );

    await dispatch(
      "auth/fetchUser",
      null,
      {
        root: true,
      },
    );

    return response;
  },
};
