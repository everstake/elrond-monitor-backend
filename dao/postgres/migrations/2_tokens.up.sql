-- +migrate Up
create type tokenType as enum ('FungibleESDT', 'NonFungibleESDT', 'SemiFungibleESDT', 'MetaESDT');

create table tokens
(
    tkn_identity   varchar(255)              not null,
    tkn_name       varchar(255)              not null,
    tkn_type       tokenType               not null,
    tkn_owner      varchar(62)               not null,
    tkn_supply     numeric(52, 20) default 0 not null,
    tkn_decimals   integer                   not null,
    tkn_properties json                      not null,
    tkn_roles      json                      not null,
    constraint tokens_pk
        primary key (tkn_identity)
);
create index tokens_tkn_type_index
    on tokens (tkn_type);

create table nft_collections
(
    nfc_name       varchar(255) not null,
    nfc_identity   varchar(100) not null,
    nfc_owner      varchar(62)  not null,
    nfc_type       tokenType    not null,
    nfc_properties json         not null,
    nfc_created_at timestamp    not null
);
