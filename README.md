# Synth socket

An idiomatic, minimal-dependencies real-time chat app.

## Installing

**1. Prerequisites:**

- [Docker](https://docs.docker.com/get-started/get-docker/) installed with [docker-compose](https://docs.docker.com/get-started/get-docker/) pluggin.
- [Ngrok](https://dashboard.ngrok.com/get-started/setup/linux) (optional) if you wanna publish the app to the internet.

**2. Clone repository:**

  ```sh
    git clone https://github.com/anhtr13/synth-socket.git
    cd synth-socket
  ```

**3.Usage:**

- Run the application locally:

  ```sh
    docker compose up --build
  ```

- If you wanna publish to the internet to play with your friends:

    - Install Ngrok and add authtoken by following the instructions [here](https://dashboard.ngrok.com/get-started/setup/linux).

    -   Publish the application online:
        ```sh
            ngrok http 3000
        ```

## Dependencies

- [coder/websocket](https://github.com/coder/websocket)
- [jackc/pgx](https://github.com/jackc/pgx/v5)
- VueJS, Pinia, Zod, Tailwindcss, etc.

## Architecture (Oct 2025)

- **Overview**
    ![overview](./overview.png)

