import Vue from "vue";

export const state = () => ({
  sponsors: [],
});

export const getters = {
  sponsors(state) {
    return state.sponsors;
  },
};

export const mutations = {
  SET_SPONSORS(state, sponsors) {
    Vue.set(state, "sponsors", sponsors);
  },
};

export const actions = {
  async create(
    _,
    {
      name,
      logo,
    },
  ) {
    const data = new FormData();
    data.set("name", name);
    data.set("logo", logo);

    const response = await this.$api.$post(
      "/sponsor",
      data,
    );

    return response;
  },

  async remove(
    {
      dispatch,
    },
    sponsor,
  ) {
    const response = await this.$api.$delete(
      `/sponsor/${ sponsor.id }`,
    );

    if ("success" === response?.status) {
      await dispatch("listAll");
    }

    return response;
  },

  async listAll(
    {
      commit,
    },
  ) {
    const response = await this.$api.$get(
      "/sponsor",
    );

    const { data = [] } = response || {};

    commit("SET_SPONSORS", data);

    return data;
  },

  async swap(
    {
      commit,
    },
    {
      a,
      b,
    },
  ) {
    const response = await this.$api.$patch(
      `/sponsor/swap/${ a.id }`,
      {
        with: b.id,
      },
    );

    if ("success" === response?.status) {
      commit("SET_SPONSORS", response.data);
    }

    return response;
  },
};
