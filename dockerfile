# BUILDER --------------------------------
FROM node:24 as builder
WORKDIR /app
COPY ./front-end/package*.json ./
RUN npm install
COPY ./front-end/ ./
RUN npm run build

# RUNNER --------------------------------
FROM nginx as runner
RUN mkdir -p /var/www/synth_socket/static
COPY --from=builder /app/dist /var/www/synth_socket/static
COPY ./nginx.conf /etc/nginx/nginx.conf
