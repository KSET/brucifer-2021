export default function(
  {
    $axios,
    $config,
  },
  inject,
) {
  const api = $axios.create({});

  api.setBaseURL(`${ $config.API_URL_BACKEND || $config.API_URL }/api`);
  api.setHeader("Content-Type", "application/json");

  api.interceptors.response.use(
    (response) => response,
    ({ response }) => response,
  );

  inject("api", api);
}
