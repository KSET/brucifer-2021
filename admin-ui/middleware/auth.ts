import {
  Context,
  Middleware,
} from "@nuxt/types";

const PAGE_NAME_LOGIN = "PageLogin";
const PAGE_NAME_REGISTER = "PageRegister";
const PAGE_NAME_HOME = "PageHome";

const auth: Middleware = (context: Context) => {
  const isAuthenticated =
    context.store.getters["auth/loggedIn"]
  ;
  const isOnAuthPage =
    context.route.name === PAGE_NAME_LOGIN ||
    context.route.name === PAGE_NAME_REGISTER
  ;

  if (!isAuthenticated && !isOnAuthPage) {
    return context.redirect({
      name: PAGE_NAME_LOGIN,
    });
  }

  if (isAuthenticated && isOnAuthPage) {
    return context.redirect({
      name: PAGE_NAME_HOME,
    });
  }
};

export default auth;
