CREATE TABLE storm_reports
(
    report_date DATE         NOT NULL,
    storm_type  VARCHAR(32)  NOT NULL,
    latitude    VARCHAR(32),
    longitude   VARCHAR(32),
    location    VARCHAR(128) NOT NULL,
    county      VARCHAR(64)  NOT NULL,
    state       VARCHAR(2)   NOT NULL,
    comments    VARCHAR,
    speed       INTEGER,
    size        INTEGER,
    f_scale     INTEGER,
    time        TIMESTAMP WITH TIME ZONE,
    PRIMARY KEY (report_date, storm_type, location, county, state)
)
