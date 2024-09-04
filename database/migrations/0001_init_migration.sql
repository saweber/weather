CREATE TABLE storm_reports (
                               id serial PRIMARY KEY,
                               report_date DATE NOT NULL,
                               storm_type VARCHAR(64) NOT NULL,
                               latitude VARCHAR(32) NOT NULL,
                               longitude VARCHAR(32) NOT NULL,
                               location VARCHAR(128),
                               county VARCHAR,
                               state_code VARCHAR(2),
                               comments VARCHAR,
                               description TEXT
)
