# weather

## Setup
- Run `docker compose up` in a terminal in this folder.
- `docker exec -it weather-kafka-1 bash -c "/bin/kafka-topics --create --topic raw-weather-reports --bootstrap-server localhost:9092"`
- Connect to database (see docker compose file, default credentials) and run script/command in database/migrations

## See other READMEs for each service
Each service has its own readme, please also read those.

## Overall Structure
I went with a monorepo structure here, which each service isolated in its own folder.

## Database
I chose Postgres for familiarity and simplicity.
Id is an increment counter for simplicity, a UUID or a custom id format likely makes more sense here.
A relational database also makes sense here, as it provides great query flexibility and scalability.
No DB migration tooling has been set up, but is also needed for CI/CD purposes.
Based on business needs, I would further refine the data types for columns in the table. Depends on
how the data is consumed and what is relevant. Examples: normalizing 'size' to a standard measurement
in mm, Latitude/Longitude to a specific format.

## Docker Compose
This is the minimum to spin up a development environment, and not intended for a public deployment.
Security is nonexistent.
