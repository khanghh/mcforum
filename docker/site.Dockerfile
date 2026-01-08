# ---- site builder ----
FROM node:20-alpine AS site-builder

WORKDIR /workdir

COPY ./site ./
# RUN npm install -g pnpm --registry=https://registry.npmmirror.com
# RUN pnpm install --registry=https://registry.npmmirror.com
RUN npm install -g pnpm
RUN pnpm install
RUN pnpm build:docker

# ---- Final site image ----
FROM node:20-alpine

ENV NODE_ENV=production
ENV PORT=8080
ENV BBSGO_BACKEND_URL=http://localhost:3000/api/

WORKDIR /app

COPY --from=site-builder /workdir/.output ./.output
COPY --from=site-builder /workdir/node_modules ./node_modules

EXPOSE 8080

ENTRYPOINT ["node", "/app/.output/server/index.mjs"]
