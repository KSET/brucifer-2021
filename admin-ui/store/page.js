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
      background,
      backgroundMobile,
    },
  ) {
    const data = new FormData();
    data.set("name", name);
    data.set("markdown", markdown);
    data.set("background", background);
    data.set("backgroundMobile", backgroundMobile);

    const response = await this.$api.$post(
      "/page/",
      data,
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
      background,
      backgroundId,
      backgroundMobile,
      backgroundMobileId,
    },
  ) {
    const data = new FormData();
    data.set("name", name);
    data.set("markdown", markdown);
    data.set("background", background);
    data.set("backgroundId", backgroundId);
    data.set("backgroundMobile", backgroundMobile);
    data.set("backgroundMobileId", backgroundMobileId);

    const response = await this.$api.$patch(
      `/page/${ id }`,
      data,
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
