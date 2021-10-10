import Vue from "vue";

export const state = () => ({
  sponsors: [],
  sponsor: null,
});

export const getters = {
  sponsors(state) {
    return state.sponsors;
  },
  sponsor(state) {
    return state.sponsor;
  },
};

export const mutations = {
  SET_SPONSORS(state, sponsors) {
    Vue.set(state, "sponsors", sponsors);
  },
  SET_SPONSOR(state, sponsor) {
    Vue.set(state, "sponsor", sponsor);
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

  async update(
    _,
    {
      id,
      name,
      logo,
    },
  ) {
    const data = new FormData();
    data.set("name", name);
    if (logo) {
      data.set("image", logo);
    }

    const response = await this.$api.$patch(
      `/sponsor/${ id }`,
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

  async info(
    {
      commit,
    },
    {
      id,
    },
  ) {
    const response = await this.$api.$get(
      `/sponsor/${ id }`,
    );

    const { status = "error", data = null } = response || {};

    if ("success" === status) {
      commit("SET_SPONSOR", data);
    } else {
      commit("SET_SPONSOR", null);
    }

    return data;
  },
};
