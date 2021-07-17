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

create type miniblock_type as ENUM ('TxBlock', 'SmartContractResultBlock', 'InvalidBlock', 'RewardsBlock');
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

create index miniblocks_mlk_receiver_block_hash_index
    on miniblocks (mlk_receiver_block_hash);

create index miniblocks_mlk_sender_block_hash_index
    on miniblocks (mlk_sender_block_hash);

create type tx_status as ENUM ('success', 'fail', 'invalid');
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

create index transactions_mlk_mini_block_hash_index
    on transactions (mlk_mini_block_hash);

create table sc_results
(
    scr_hash    varchar(64) not null
        constraint sc_results_pk
            primary key,
    trn_hash    varchar(64) not null
        constraint sc_results_transactions_trn_hash_fk
            references transactions,
    scr_from    varchar(62) not null,
    scr_to      varchar(62) not null,
    scr_value   numeric(36) not null,
    scr_data    text        not null,
    scr_message text        not null
);

create
    index sc_results_trn_hash_index
    on sc_results (trn_hash);


create table accounts
(
    acc_address    varchar(255) not null
        primary key,
    acc_created_at timestamp(0) not null
);

create table rewards
(
    rwd_tx_hash          varchar(64)     not null
        constraint rewards_pk
            primary key,
    rwd_hyperblock_id    bigint          not null,
    rwd_receiver_address varchar(64)     not null,
    rwd_amount           numeric(36, 18) not null,
    rwd_created_at       timestamp       not null
);

create index rewards_rwd_receiver_address_index
    on rewards (rwd_receiver_address);

create table delegations
(
    dlg_tx_hash    varchar(64)     not null
        constraint delegations_pk
            primary key,
    dlg_delegator  varchar(64)     not null,
    dlg_validator  varchar(64)     not null,
    dlg_amount     numeric(36, 18) not null,
    dlg_created_at timestamp       not null
);

create table stakes
(
    stk_tx_hash    varchar(64)     not null
        constraint stakes_pk
            primary key,
    stk_validator  varchar(64)     not null,
    stk_amount     numeric(36, 18) not null,
    stk_created_at timestamp       not null
);

