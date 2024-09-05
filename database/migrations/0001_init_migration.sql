CREATE TABLE storm_reports (
    id serial PRIMARY KEY,
    report_date DATE NOT NULL,
    storm_type VARCHAR(32) NOT NULL,
    latitude VARCHAR(32),
    longitude VARCHAR(32),
    location VARCHAR(128),
    county VARCHAR(64),
    state VARCHAR(2),
    comments VARCHAR,
    speed INTEGER,
    size INTEGER,
    FScale INTEGER
)
