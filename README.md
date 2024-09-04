# weather

## Setup
- `docker compose up`
- Connect to database (see docker compose file, default credentials) and run script/command in database/migrations

## See other READMEs for each service
Each service has its own readme, please also read those.

## Database
I chose Postgres for familiarity and simplicity.
Id is an increment counter for simplicity, a UUID or a custom id format likely makes more sense here.
A relational database also makes sense here, as it provides great query flexibility and scalability.
No DB migration tooling has been set up, but is also needed for CI/CD purposes.

## Docker Compose
This is the minimum to spin up a development environment, and not intended for a public deployment.
Security is nonexistent.
