# bits-and-bytes-bot
Discord bot for bits &amp; bytes book club

### Running locally
- Copy `.env.example` and rename it to `.env`
- Add `DISCORD_API_SECRET`, `GUILD_ID`, and `POSTGRES_URI`
- For more information on how to set up a discord api client please see: https://discord.com/developers/docs/getting-started
  - You will need to set up the `bot` scope and `Send Messages`, `Use Slash Commands`, and `Read Messages/View Channels` permissions under the Oauth2 settings
- TODO: SETUP DOCKER COMPOSE WITH PG DB
- run `ENV=local go run .` to start bot

# Docker
You can build a docker image by running the `task docker-build` command.
