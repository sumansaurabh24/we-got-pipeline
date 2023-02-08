CREATE TABLE files (
    id SERIAL PRIMARY KEY ,
    batch_id VARCHAR(50) NOT NULL ,
    filename VARCHAR(100) NOT NULL ,
    file_path TEXT NOT NULL ,
    status VARCHAR(100) NOT NULL ,
    stage VARCHAR(100) NOT NULL ,
    last_entry JSONB NULL ,
    timestamp TIMESTAMP NOT NULL
);

CREATE TABLE time_series_data(
    file_path TEXT NOT NULL ,
    flat_no VARCHAR(10) NOT NULL ,
    metadata JSONB NOT NULL,
    total_consumption INT NOT NULL ,
    timestamp TIMESTAMPTZ NOT NULL
);

SELECT create_hypertable('time_series_data', 'timestamp');

CREATE INDEX ix_tsd_timestamp ON time_series_data (flat_no, timestamp DESC );

CREATE INDEX idx_gin_tsd_metadata ON time_series_data USING GIN (metadata);