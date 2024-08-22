-- +migrate Up
CREATE TABLE tg_user
(
    tg_id            BIGINT NOT NULL,
    account_name     TEXT   NOT NULL,
    full_name        TEXT   NOT NULL,
    chosen_family_id UUID NULL,
    PRIMARY KEY (tg_id)
);

CREATE TABLE family
(
    id   UUID        NOT NULL,
    name VARCHAR(64) NOT NULL UNIQUE,
    PRIMARY KEY (id)
);

CREATE TABLE family_user
(
    user_id   BIGINT REFERENCES tg_user (tg_id) ON UPDATE CASCADE ON DELETE CASCADE,
    family_id UUID REFERENCES family (id) ON UPDATE CASCADE,
    CONSTRAINT family_user_pkey PRIMARY KEY (user_id, family_id) -- explicit pk
);

CREATE TABLE item
(
    id   UUID NOT NULL,
    name VARCHAR(64) NULL,
    link TEXT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE category
(
    id   UUID        NOT NULL,
    name VARCHAR(64) NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE category_item
(
    category_id UUID REFERENCES category (id) ON UPDATE CASCADE ON DELETE CASCADE,
    item_id     UUID REFERENCES item (id) ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT category_item_pkey PRIMARY KEY (category_id, item_id) -- explicit pk
);

CREATE TABLE category_family
(
    category_id UUID REFERENCES category (id) ON UPDATE CASCADE ON DELETE CASCADE,
    family_id   UUID REFERENCES family (id) ON UPDATE CASCADE,
    CONSTRAINT category_family_pkey PRIMARY KEY (category_id, family_id) -- explicit pk
);

-- +migrate Down

DROP TABLE IF EXISTS tg_user CASCADE;
DROP TABLE IF EXISTS family CASCADE;
DROP TABLE IF EXISTS family_user CASCADE;
DROP TABLE IF EXISTS item CASCADE;
DROP TABLE IF EXISTS category CASCADE;
DROP TABLE IF EXISTS category_item CASCADE;
DROP TABLE IF EXISTS category_family CASCADE;