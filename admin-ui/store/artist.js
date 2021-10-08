import Vue from "vue";

export const state = () => ({
  artists: [],
  artist: null,
});

export const getters = {
  artists(state) {
    return state.artists;
  },
  artist(state) {
    return state.artist;
  },
};

export const mutations = {
  SET_ARTISTS(state, artists) {
    Vue.set(state, "artists", artists);
  },
  SET_ARTIST(state, artist) {
    Vue.set(state, "artist", artist);
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
      "/artist",
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
      data.set("logo", logo);
    }

    const response = await this.$api.$patch(
      `/artist/${ id }`,
      data,
    );

    return response;
  },

  async remove(
    {
      dispatch,
    },
    artist,
  ) {
    const response = await this.$api.$delete(
      `/artist/${ artist.id }`,
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
      "/artist",
    );

    const { data = [] } = response || {};

    commit("SET_ARTISTS", data);

    return data;
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
      `/artist/${ id }`,
    );

    const { data = null } = response || {};

    commit("SET_ARTIST", data);

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
      `/artist/swap/${ a.id }`,
      {
        with: b.id,
      },
    );

    if ("success" === response?.status) {
      commit("SET_ARTISTS", response.data);
    }

    return response;
  },
};
