# API

This service is for querying the postgres database for storm events.

This application is intentionally kept very simple; if there were significantly more endpoints, 
a more sophisticated project structure would be warranted.

## cURL examples
```curl "http://localhost:8080/stormReports?windowStart=2024-01-01&windowEnd=2024-01-31"```
Get storm reports for January 2024
windowStart and windowEnd can be used by themselves, or together.

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
- Be consistent with `state` vs `state_code` in API vs storage
- Database connection is hardcoded; and should be moved to environment variables
- Better time input, right now only accepts days in 'YYYY-MM-DD' format
- Replace print statements with structured logger
- Improve query flexibility with more options and a query language for date ranges, partial and exact string matches, etc.
- API documentation is non-existent
- No API security
- No observability
- No pagination
