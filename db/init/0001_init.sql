CREATE USER gopher
    PASSWORD 'gopher';

CREATE DATABASE metrics_db
    OWNER 'gopher'
    ENCODING 'UTF8'
    LC_COLLATE = 'en_US.utf8'
    LC_CTYPE = 'en_US.utf8';

CREATE INDEX updated_at ON gauge_metrics (updated_at);
CREATE INDEX updated_at ON counter_metrics (updated_at);