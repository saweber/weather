# API

This service is for querying the postgres database for storm events.

This application is intentionally kept very simple; if there were significantly more endpoints, 
a more sophisticated project structure would be warranted.

## cURL examples
```curl "http://localhost:8080/stormReports?windowStart=2024-01-01T00%3A00%3A00Z&windowEnd=2024-02-01T00%3A00"```
Get storm reports for January 2024
windowStart and windowEnd can be used by themselves, or together.
Time must be RFC3339 formatted.

```curl "http://localhost:8080/stormReports?location=here"```
Get storm reports where the location contains the string 'here'

```curl "http://localhost:8080/stormReports"```
Get all storm reports

## Key Decisions
- Separate out database logic with a repository struct
- Use bun to build and execute queries
- In the interest of time, keep the endpoint very simple
- Error handling is very minimal
- No service layer, since there is no business logic to handle

## Shortcuts / Next Steps
- Add request level testing and a in-memory repository implementation for testing
- Database connection is hardcoded; and should be moved to environment variables
- Better time input
- Replace print statements with structured logger
- Improve query flexibility with more options and a query language for date ranges, partial and exact string matches, etc.
- OpenAPI documentation
- No security
- No observability
- No pagination
