import Vue from "vue";

export const state = () => ({
  pages: [],
  page: null,
});

export const getters = {
  pages(state) {
    return state.pages;
  },
  page(state) {
    return state.page;
  },
};

export const mutations = {
  SET_PAGES(state, pages) {
    Vue.set(state, "pages", pages);
  },
  SET_PAGE(state, page) {
    Vue.set(state, "page", page);
  },
};

export const actions = {
  async create(
    _,
    {
      name,
      markdown,
    },
  ) {
    const response = await this.$api.$post(
      "/page/",
      {
        name,
        markdown,
      },
    );

    return response;
  },

  async list(
    {
      commit,
    },
  ) {
    const response = await this.$api.$get(
      "/page/",
    );

    if ("success" === response?.status) {
      commit("SET_PAGES", response.data);
    } else {
      commit("SET_PAGES", []);
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
      `/page/${ id }`,
    );

    if ("success" === response?.status) {
      commit("SET_PAGE", response.data);
    } else {
      commit("SET_PAGE", null);
    }

    return response;
  },

  async update(
    {
      commit,
    },
    {
      id,
      name,
      markdown,
    },
  ) {
    const response = await this.$api.$patch(
      `/page/${ id }`,
      {
        name,
        markdown,
      },
    );

    if ("success" === response?.status) {
      commit("SET_PAGE", response.data);
    } else {
      commit("SET_PAGE", null);
    }

    return response;
  },

  async remove(
    {
      dispatch,
    },
    {
      id,
    },
  ) {
    const response = await this.$api.$delete(
      `/page/${ id }`,
    );

    if ("success" === response?.status) {
      await dispatch("list");
    }

    return response;
  },
};
