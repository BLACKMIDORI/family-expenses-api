CREATE TABLE app_user
(
    id text not null unique,
    creation_unix bigint not null,
    PRIMARY KEY (id)
);

CREATE TABLE app_user_login
(
    id                text not null unique,
    creation_unix bigint not null,
    identity_provider text not null,
    key               text not null,
    fk_app_user_id    text not null,
    PRIMARY KEY (id),
    FOREIGN KEY (fk_app_user_id) REFERENCES app_user (id),
    CONSTRAINT unique_identity_provider_and_key UNIQUE (identity_provider, key)
);

CREATE TABLE persistent_grant
(
    id text not null unique,
    creation_unix bigint not null,
    key_digest text unique,
    client_id text not null,
    fk_app_user_id  text not null,
    session_id  text not null,
    expiration_unix bigint not null,
    consumed_unix bigint,
    PRIMARY KEY (id),
    FOREIGN KEY (fk_app_user_id) REFERENCES app_user (id)
);


CREATE TABLE workspace
(
    id             text not null unique,
    creation_unix bigint not null,
    name           text not null,
    fk_app_user_id text not null,
    PRIMARY KEY (id),
    FOREIGN KEY (fk_app_user_id) REFERENCES app_user (id)
);
