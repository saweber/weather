# Collector

This service schedules daily tasks for retrieving CSV files from NOAA, parsing them,
and pushing each individual storm report into the raw-weather-reports kafka topic.

## Key Decisions
- Keep CSV parsing standardized so it can be reused
- Use standard lib for CSV parsing
- Normalize the tornado, hail, wind reports to have the message data be consistent
- No particular reason behind sarama, other than it seems to be the most popular based
on internet search and GitHub stars and does not require CGO.

## Shortcuts / Next Steps
- Retry logic, improved error handling, and metrics/alerting for edge cases, 
particularly for download failures, but also sending messages to Kafka.
- Replace print statements with structured logger
- URLs and Kafka connection info is hardcoded, and should be moved to environment variables
- Add metadata in the message for analytics, monitoring, troubleshooting
- No security
- No observability
- To reduce memory consumption for very large CSVs, finding a way to stream the repsonse 
rather than hold all of it in memory would be an improvement.
- More testing and hardening around timezones (ensure everything is as consistent as possible)
