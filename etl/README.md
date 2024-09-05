# ETL

This service consumes from the raw-weather-reports kafka topic, and inserts transformed records 
into the postgres database for the API. Additionally, the transformed records are placed into a
new kafka topic - transformed-weather-reports for consumption by other future services.

## Key Decisions
- No particular reason behind sarama, other than it seems to be the most popular based 
on internet search and GitHub stars and does not require CGO.
- The code to handle consumption was lifted from sarama's documentation and simplified/modified
for use here.

## Shortcuts / Next Steps
- Needs more cleanup and refactoring - I am not fully satisfied on the boundary between the ETL and consumer structs.
- Retry logic, improved error handling, and metrics/alerting for edge cases,
  particularly for connection to Kafka failures
- Replace print statements with structured logger
- URLs and Kafka and Postgres connection info is hardcoded, and should be moved to environment variables
- Add metadata in the message for analytics, monitoring, troubleshooting
- No security
- No observability
