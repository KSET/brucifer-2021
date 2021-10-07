FROM keymetrics/pm2:16-alpine

WORKDIR /app

COPY yarn.lock package.json ./

ENV NPM_CONFIG_LOGLEVEL warn

RUN yarn

ENV NUXT_HOST 0.0.0.0
ENV NUXT_PORT 3000
ENV NODE_ENV production

COPY . .

RUN yarn build

CMD [ "sh", "-c", "pm2 start ecosystem.config.js && pm2 logs" ]
