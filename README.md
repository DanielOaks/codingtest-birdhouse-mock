# birdhouse-mock

This is a companion to the [birdhouse-admin frontend](https://github.com/DanielOaks/codingtest-birdhouse-admin), and creates a mock server to test against with fake data.

## Environment variables

You can configure the server through these environment variables:

- `BH_REGISTRATIONS`: How many registrations to create. Default: 20
- `BH_EMPTY_REGISTRATIONS`: How many registrations should be empty and not contain a birdhouse, percentage-wise. Default: 0.1 (10%)
- `BH_OCCUPANCY_WEEKS`: How many weeks to generate occupancy history for. Default: 25
- `BH_UPDATES_PER_WEEK`: How many updates to generate per week. Default: 14
- `BH_SERVE_PORT`: Which port to expose the API endpoints on. Default: 7000

## Docker quick start

You can get the mock server up and running quickly by using Docker!

### Docker

Run the prebuilt docker image with this command:

```bash
docker run -it -e -p 7000:7000 ghcr.io/danieloaks/codingtest-birdhouse-mock:release
```

### Docker compose

The default [docker compose file](./compose.yaml) sets up the mock server and the admin panel:

```bash
docker compose up
```

If you're running the compose file on a server, you'll want to change the `NUXT_PUBLIC_API_BASE` variable to point towards the server's public address rather than `localhost`.

The dashboard will be exposed on port 3000. If running the command locally, at http://localhost:3000
