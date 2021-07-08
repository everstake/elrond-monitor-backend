-- +migrate Up
create sequence parsers_seq;

create table parsers
(
    par_id     int default nextval('parsers_seq')
        primary key,
    par_title  varchar(255)  not null,
    par_height int default 0 not null,
    constraint parsers_par_title_uindex
        unique (par_title)
);

insert into parsers (par_id, par_title, par_height)
VALUES (1, 'elrond', 4910084);

create type block_status as ENUM ('on-chain');

create table blocks
(
    blk_hash             varchar(64)     not null
        constraint blocks_pkey
            primary key,
    blk_nonce            integer         not null,
    blk_round            integer         not null,
    blk_shard            bigint          not null,
    blk_num_txs          integer         not null,
    blk_epoch            integer         not null,
    blk_status           block_status    not null,
    blk_prev_block_hash  varchar(64)     not null,
    blk_created_at       timestamp       not null,
    blk_developer_fees   numeric(36, 18) not null,
    blk_accumulated_fees numeric(36, 18) not null
);

create type miniblock_type as ENUM ('TxBlock', 'SmartContractResultBlock', 'InvalidBlock');

create table miniblocks
(
    mlk_hash                varchar(64)    not null
        constraint miniblocks_pk
            primary key,
    mlk_receiver_block_hash varchar(64),
    mlk_receiver_shard      bigint         not null,
    mlk_sender_block_hash   varchar(64)    not null,
    mlk_sender_shard        bigint         not null,
    mlk_type                miniblock_type not null,
    mlk_created_at          timestamp(0)   not null
);
create type tx_status as ENUM ('success', 'fail');

create table transactions
(
    trn_hash            varchar(64)     not null
        constraint transactions_pkey
            primary key,
    trn_status          tx_status       not null,
    mlk_mini_block_hash varchar(64)     not null,
    trn_value           numeric(36, 18) not null,
    trn_fee             numeric(36, 18) not null,
    trn_sender          varchar(62)     not null,
    trn_sender_shard    bigint          not null,
    trn_receiver        varchar(62)     not null,
    trn_receiver_shard  bigint          not null,
    trn_gas_price       numeric         not null,
    trn_gas_used        numeric         not null,
    trn_nonce           integer         not null,
    trn_data            text            not null,
    trn_created_at      timestamp       not null
);

create table sc_results
(
    scr_hash  varchar(64) not null
        constraint sc_results_pk
            primary key,
    trn_hash  varchar(64) not null,
    scr_from  varchar(62) not null,
    scr_to    varchar(62) not null,
    scr_value numeric(36) not null,
    scr_data  text        not null
);


create table accounts
(
    acc_address    varchar(255) not null
        primary key,
    acc_created_at timestamp(0) not null
);



