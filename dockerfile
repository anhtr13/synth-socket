FROM node:latest as build-stage
WORKDIR /app
COPY ./front-end/package*.json ./
RUN npm install
COPY ./front-end/ ./
RUN npm run build

FROM nginx as production-stage
RUN mkdir /static
COPY --from=build-stage /app/dist /static
COPY ./nginx.conf /etc/nginx/nginx.conf
