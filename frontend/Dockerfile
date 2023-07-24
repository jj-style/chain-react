FROM node:20-alpine AS build

WORKDIR /build

COPY package.json package.json
COPY package-lock.json package-lock.json
RUN npm install --force

COPY public/ public
COPY src/ src
COPY .env.production .
RUN npm run build:prod

FROM nginx:alpine
COPY --from=build /build/build/ /usr/share/nginx/html
