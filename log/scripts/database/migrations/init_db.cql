CREATE KEYSPACE db
WITH replication = {
    'class': 'SimpleStrategy',
    'replication_factor': 1
};

CREATE TABLE db.migrations (
    cluster text,
    timestamp text,
    script text,
    PRIMARY KEY ((cluster), timestamp, script)
)
WITH CLUSTERING ORDER BY (
    timestamp DESC
);

CREATE KEYSPACE log
WITH replication = {
    'class': 'SimpleStrategy',
    'replication_factor': 1
};

CREATE TABLE log.source_logs (
    source_type text,
    timestamp timestamp,
    transaction_id uuid,
    source_id uuid,
    level tinyint,
    location text,
    data text,
    PRIMARY KEY ((source_type), timestamp, transaction_id, source_id)
)
WITH CLUSTERING ORDER BY (
    timestamp DESC
);

CREATE TABLE log.transaction_logs (
    transaction_id uuid,
    transaction_timestamp timestamp,
    timestamp timestamp,
    source_id uuid,
    source_type text,
    level tinyint,
    location text,
    data text,
    PRIMARY KEY ((transaction_id), transaction_timestamp, timestamp, source_id)
)
WITH CLUSTERING ORDER BY (
    transaction_timestamp DESC,
    timestamp DESC
);