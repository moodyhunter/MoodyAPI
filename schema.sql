CREATE TABLE clients (
    id          UUID        NOT NULL    UNIQUE  PRIMARY KEY,
    client_name VARCHAR     NOT NULL    UNIQUE
);

CREATE TABLE notifications (
    id          SERIAL      NOT NULL    UNIQUE  PRIMARY KEY,
    client_id   UUID        NOT NULL            REFERENCES clients(id),
    title       VARCHAR     NOT NULL,
    content     VARCHAR     NOT NULL,
    time        TIMESTAMP   NOT NULL
);

CREATE TABLE wg_clients (
    id              SERIAL      NOT NULL    UNIQUE  PRIMARY KEY ,
    public_key      VARCHAR     NOT NULL    UNIQUE,
    last_seen       TIMESTAMP               UNIQUE,
    endpoint        INET,
    endpoint_port   INTEGER
);

CREATE TABLE clients_wireguard (
    client_id       UUID    NOT NULL    REFERENCES clients(id),
    wireguard_devid SERIAL  NOT NULL    REFERENCES wg_clients(id),
    CONSTRAINT clients_wireguard_pkey PRIMARY KEY (client_id, wireguard_devid)
);
