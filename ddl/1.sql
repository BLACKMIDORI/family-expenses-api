CREATE TABLE app_user
(
    id            text   not null unique,
    creation_unix bigint not null,
    PRIMARY KEY (id)
);

CREATE TABLE app_user_login
(
    id                text   not null unique,
    creation_unix     bigint not null,
    identity_provider text   not null,
    key               text   not null,
    fk_app_user_id    text   not null,
    PRIMARY KEY (id),
    FOREIGN KEY (fk_app_user_id) REFERENCES app_user (id),
    CONSTRAINT unique_identity_provider_and_key UNIQUE (identity_provider, key)
);

CREATE TABLE persistent_grant
(
    id              text   not null unique,
    creation_unix   bigint not null,
    key_digest      text unique,
    client_id       text   not null,
    fk_app_user_id  text   not null,
    session_id      text   not null,
    expiration_unix bigint not null,
    consumed_unix   bigint,
    PRIMARY KEY (id),
    FOREIGN KEY (fk_app_user_id) REFERENCES app_user (id)
);


CREATE TABLE workspace
(
    id             text   not null unique,
    creation_unix  bigint not null,
    name           text   not null,
    fk_app_user_id text   not null,
    PRIMARY KEY (id),
    FOREIGN KEY (fk_app_user_id) REFERENCES app_user (id)
);

CREATE TABLE expense
(
    id              text   not null unique,
    creation_unix   bigint not null,
    name            text   not null,
    fk_workspace_id text   not null,
    PRIMARY KEY (id),
    FOREIGN KEY (fk_workspace_id) REFERENCES workspace (id)
);

CREATE TABLE payer
(
    id              text   not null unique,
    creation_unix   bigint not null,
    name            text   not null,
    fk_workspace_id text   not null,
    PRIMARY KEY (id),
    FOREIGN KEY (fk_workspace_id) REFERENCES workspace (id)
);

CREATE TABLE charges_model
(
    id              text   not null unique,
    creation_unix   bigint not null,
    name            text   not null,
    fk_workspace_id text   not null,
    PRIMARY KEY (id),
    FOREIGN KEY (fk_workspace_id) REFERENCES workspace (id)
);

CREATE TABLE charge_association
(
    id                  text   not null unique,
    creation_unix       bigint not null,
    name                text   not null,
    fk_expense_id       text   not null,
    fk_payer_id         text   not null,
    fk_charges_model_id text   not null,
    PRIMARY KEY (id),
    FOREIGN KEY (fk_expense_id) REFERENCES expense (id),
    FOREIGN KEY (fk_payer_id) REFERENCES payer (id),
    FOREIGN KEY (fk_charges_model_id) REFERENCES charges_model (id)
);

CREATE TABLE payer_payment_weight
(
    id                       text   not null unique,
    creation_unix            bigint not null,
    weight                   real   not null,
    fk_payer_id              text   not null,
    fk_charge_association_id text   not null,
    PRIMARY KEY (id),
    FOREIGN KEY (fk_payer_id) REFERENCES payer (id),
    FOREIGN KEY (fk_charge_association_id) REFERENCES charge_association (id)
);
