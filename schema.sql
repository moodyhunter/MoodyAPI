DROP TABLE IF EXISTS clients CASCADE;
DROP TABLE IF EXISTS notifications CASCADE;
DROP TABLE IF EXISTS wg_clients CASCADE;
DROP TABLE IF EXISTS wg_allowed_ips CASCADE;
DROP TABLE IF EXISTS operation_logs CASCADE;

CREATE TABLE clients (
    id          SERIAL      NOT NULL    UNIQUE  PRIMARY KEY,
    client_name VARCHAR     NOT NULL,
    client_uuid VARCHAR     NOT NULL    UNIQUE,
    privileged  BOOLEAN     NOT NULL,
    last_seen   TIMESTAMP,
    enabled     BOOLEAN     NOT NULL
);

CREATE TABLE notifications (
    id          SERIAL      NOT NULL    UNIQUE  PRIMARY KEY,
    sender_id   INTEGER     NOT NULL            REFERENCES clients(id),
    channel_id  INTEGER     NOT NULL,
    title       VARCHAR     NOT NULL,
    content     VARCHAR     NOT NULL,
    time        TIMESTAMP   NOT NULL,
    icon        VARCHAR
);

CREATE TABLE wg_clients (
    id              SERIAL      NOT NULL    UNIQUE  PRIMARY KEY,
    public_key      VARCHAR     NOT NULL    UNIQUE,
    last_seen       TIMESTAMP               UNIQUE,
    endpoint        INET,
    endpoint_port   INTEGER,
    enabled         BOOLEAN     NOT NULL
);

CREATE TABLE wg_allowed_ips (
    id          SERIAL  NOT NULL    UNIQUE  PRIMARY KEY,
    wg_id       SERIAL  NOT NULL            REFERENCES wg_clients(id),
    allowed_ip  INET    NOT NULL    UNIQUE
);

CREATE TABLE operation_logs (
    id          SERIAL      NOT NULL    UNIQUE  PRIMARY KEY,
    client_id   NUMERIC     NOT NULL,
    client_name VARCHAR     NOT NULL,
    time        TIMESTAMP   NOT NULL,
    operation   VARCHAR     NOT NULL
);
