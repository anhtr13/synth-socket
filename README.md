# Synth socket

An idiomatic, minimal-dependencies real-time chat app.

## Installation

**1. Prerequisites:**

- [Docker](https://docs.docker.com/get-started/get-docker/) installed with [docker-compose](https://docs.docker.com/get-started/get-docker/) pluggin.
- [Goose](https://github.com/pressly/goose) for database schemas migration.
- [Ngrok](https://dashboard.ngrok.com/get-started/setup/linux) (optional) if you wanna publish the app to the internet.
- [Sqlc](https://docs.sqlc.dev/en/latest) (optional) if you wanna make change in the sql schemas/queries.

**2. Clone repository:**

  ```sh
    git clone https://github.com/anhtr13/synth-socket.git
    cd synth-socket
  ```

**3. Start the app:**

- Run the application locally:

  ```sh
    docker compose up --build
  ```

- Migrate the SQL schemas:

  ```sh
    goose postgres postgres://xbro:123@localhost:5433/synth_socket up --dir sql/schemas
  ```

- If you wanna publish to the internet to play with your friends:

    - Install Ngrok and add authtoken by following the instructions [here](https://dashboard.ngrok.com/get-started/setup/linux).

    -   Publish the application online:
        ```sh
            ngrok http 3000
        ```

**4. How to use:**

- First create a room.
- Make some friends and invite them to your room.
- You can create multiple rooms but not with the same name.
- You can only invite friends to rooms that you created.

## Dependencies

- [coder/websocket](https://github.com/coder/websocket) to handle websocket connections.
- [jackc/pgx](https://github.com/jackc/pgx) for Postgres driver.
- Front-end: VueJS, Pinia, Zod, Tailwindcss, etc.

## Architecture (Oct 2025)

- **Overview**
    ![overview](./overview.png)

