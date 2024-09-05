# weather

## Setup
- Run `docker compose up` in a terminal in this folder.
- `docker exec -it weather-kafka-1 bash -c "/bin/kafka-topics --create --topic raw-weather-reports --bootstrap-server localhost:9092"`
- `docker exec -it weather-kafka-1 bash -c "/bin/kafka-topics --create --topic transformed-weather-data --bootstrap-server localhost:9092"`
- Connect to database (see docker compose file, default credentials) and run script/command in database/migrations
- May need to restart the collector container to get the initial processing (or otherwise wait 24h).

## See other READMEs for each service
Each service has its own readme, please also read those.
I would recommend looking at collector, then etl, then api.

## Overall Structure
I went with a monorepo structure here, which each service isolated in its own folder.
diagram.excalidraw provides a visual of the overall system.

## Overarching Issues
- There are opportunities to move some structs and funcs to common libraries and packages rather than
duplicate code, but I am choosing to defer that work.
- I have not done nearly as much manual or automated testing as I would like here - the code is much more 'proof of concept'
than production ready.
- There are a lot of hardcoded values that need to be moved to environment variables
- There are opportunities to do better data transformations that ensure data cleanliness and consistency.

## Database
- I chose Postgres for familiarity and simplicity.
- I'm making an enormous assumption about what defines a unique report to deal with reprocessing and more than once delivery.
- A relational database also makes sense here, as it provides great query flexibility and scalability.
- No DB migration tooling has been set up, but is needed.
- Based on business needs, I would further refine the data types for columns in the table. Depends on
how the data is consumed and what is relevant. Examples: normalizing 'size' to a standard measurement
in mm, Latitude/Longitude to a specific format.

## Docker Compose
This is the minimum to spin up a development environment, and not intended for a public deployment.
Security is nonexistent.
