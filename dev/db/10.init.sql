------------------------------

CREATE SCHEMA IF NOT EXISTS lumos;

------------------------------

CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE plpgsql;


------------------------------

CREATE TABLE lumos.auth_logs (
    id BIGSERIAL PRIMARY KEY,
    login VARCHAR,
    ip VARCHAR,
    useragent VARCHAR,
    fingerprint VARCHAR,
    confidenceScore VARCHAR,
    created_at timestamptz NOT NULL DEFAULT NOW()
);

CREATE INDEX lumox_auth_logs_created_at ON
lumos.auth_logs (created_at DESC NULLS LAST);

CREATE INDEX lumos_auth_logs_login ON
lumos.auth_logs (login ASC NULLS LAST);

CREATE TABLE lumos.search_queries (
    id BIGSERIAL PRIMARY KEY,
    query VARCHAR,
    created_at timestamptz NOT NULL DEFAULT NOW()
);

CREATE INDEX lumox_search_query_query
ON lumos.search_queries (query ASC NULLS LAST);