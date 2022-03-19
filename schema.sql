-- DROP TABLE IF EXISTS clients CASCADE;
-- DROP TABLE IF EXISTS notifications CASCADE;
-- DROP TABLE IF EXISTS wg_clients CASCADE;
-- DROP TABLE IF EXISTS clients_wireguard CASCADE;
-- DROP TABLE IF EXISTS wg_allowed_ips CASCADE;

CREATE TABLE clients (
    id          SERIAL      NOT NULL    UNIQUE  PRIMARY KEY,
    client_name VARCHAR     NOT NULL    UNIQUE,
    client_uuid UUID        NOT NULL    UNIQUE,
    last_seen   TIMESTAMP
);

CREATE TABLE notifications (
    id          SERIAL      NOT NULL    UNIQUE  PRIMARY KEY,
    client_id   SERIAL      NOT NULL            REFERENCES clients(id),
    title       VARCHAR     NOT NULL,
    content     VARCHAR     NOT NULL,
    time        TIMESTAMP   NOT NULL
);

CREATE TABLE wg_clients (
    id              SERIAL      NOT NULL    UNIQUE  PRIMARY KEY,
    public_key      VARCHAR     NOT NULL    UNIQUE,
    last_seen       TIMESTAMP               UNIQUE,
    endpoint        INET,
    endpoint_port   INTEGER
);

CREATE TABLE wg_allowed_ips (
    id          SERIAL  NOT NULL    UNIQUE  PRIMARY KEY,
    wg_id       SERIAL  NOT NULL            REFERENCES wg_clients(id),
    allowed_ip  INET    NOT NULL    UNIQUE
);

CREATE TABLE clients_wireguard (
    client_id   SERIAL  NOT NULL    REFERENCES clients(id),
    wg_id       SERIAL  NOT NULL    REFERENCES wg_clients(id),
    CONSTRAINT clients_wireguard_pkey PRIMARY KEY (client_id, wg_id)
);
