import {
  Context,
  Middleware,
} from "@nuxt/types";

const PAGE_NAME_LOGIN = "PageLogin";
const PAGE_NAME_HOME = "PageHome";

const auth: Middleware = (context: Context) => {
  if (!context.store.state.auth.user && PAGE_NAME_LOGIN !== context.route.name) {
    return context.redirect({
      name: PAGE_NAME_LOGIN,
    });
  }

  if (context.store.state.auth.user && PAGE_NAME_LOGIN === context.route.name) {
    return context.redirect({
      name: PAGE_NAME_HOME,
    });
  }
};

export default auth;
