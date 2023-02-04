-- DROP TABLE IF EXISTS clients CASCADE;
-- DROP TABLE IF EXISTS notifications CASCADE;
-- DROP TABLE IF EXISTS notification_channels CASCADE;
-- DROP TABLE IF EXISTS dns CASCADE;
-- DROP TABLE IF EXISTS operation_logs CASCADE;

CREATE TABLE clients (
    id          SERIAL      NOT NULL    UNIQUE  PRIMARY KEY,
    client_name VARCHAR     NOT NULL,
    client_uuid VARCHAR     NOT NULL    UNIQUE,
    privileged  BOOLEAN     NOT NULL,
    last_seen   TIMESTAMP,
    enabled     BOOLEAN     NOT NULL
);

CREATE TABLE notification_channels (
    id          SERIAL      NOT NULL    UNIQUE  PRIMARY KEY,
    name        VARCHAR     NOT NULL    UNIQUE
);

CREATE TABLE notifications (
    id          SERIAL      NOT NULL    UNIQUE  PRIMARY KEY,
    sender_id   INTEGER     NOT NULL            REFERENCES clients(id),
    channel_id  INTEGER     NOT NULL            REFERENCES notification_channels(id),
    title       VARCHAR     NOT NULL,
    content     VARCHAR     NOT NULL,
    time        TIMESTAMP   NOT NULL,
    urgency     INTEGER     NOT NULL    DEFAULT 0,
    private     BOOLEAN     NOT NULL    DEFAULT FALSE,
    icon        VARCHAR
);

CREATE TABLE dns (
    hostname    VARCHAR     NOT NULL,
    type        VARCHAR     NOT NULL,
    ip          VARCHAR     NOT NULL,
    PRIMARY KEY (hostname, type)
);

CREATE TABLE operation_logs (
    id          SERIAL      NOT NULL    UNIQUE  PRIMARY KEY,
    client_id   NUMERIC     NOT NULL,
    client_name VARCHAR     NOT NULL,
    time        TIMESTAMP   NOT NULL,
    operation   VARCHAR     NOT NULL
);
