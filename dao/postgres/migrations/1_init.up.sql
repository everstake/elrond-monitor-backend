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
        primary key,
    blk_nonce            int             not null,
    blk_round            int             not null,
    blk_shard            bigint          not null,
    blk_num_txs          int             not null,
    blk_epoch            int             not null,
    blk_status           block_status    not null,
    blk_prev_block_hash  varchar(64)     not null,
    blk_created_at       timestamp       not null,
    blk_developer_fees   decimal(36, 18) not null,
    blk_accumulated_fees decimal(36, 18) not null
);

create type miniblock_type as ENUM ('TxBlock', 'SmartContractResultBlock');

create table miniblocks
(
    mlk_hash                varchar(64)    not null,
    mlk_receiver_block_hash varchar(64) null,
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
        primary key,
    trn_status          tx_status       not null,
    mlk_mini_block_hash varchar(64)     not null,
    trn_value           decimal(36, 18) not null,
    trn_fee             decimal(36, 18) not null,
    trn_sender          varchar(62)     not null,
    trn_sender_shard    bigint          not null,
    trn_receiver        varchar(62)     not null,
    trn_receiver_shard  bigint          not null,
    trn_gas_price       numeric         not null,
    trn_gas_used        numeric         not null,
    trn_nonce           int             not null,
    trn_data            text            not null,
    trn_created_at      timestamp       not null
);

create table sc_results
(
    trn_hash  varchar(64)     not null,
    scr_from  varchar(62)     not null,
    scr_to    varchar(62)     not null,
    scr_value decimal(36, 0) not null,
    scr_data  text            not null
);

create table accounts
(
    acc_address    varchar(255) not null
        primary key,
    acc_created_at timestamp(0) not null
);



