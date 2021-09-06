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
VALUES (1, 'elrond', 0);

create type block_status as ENUM ('on-chain');

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

create table daily_stats
(
    das_title      varchar(36)     not null,
    das_value      numeric(36, 18) not null,
    das_created_at timestamp       not null
);
create index daily_stats_das_created_at_index
    on daily_stats (das_created_at);
create index daily_stats_das_title_index
    on daily_stats (das_title);


create type stake_event_type as ENUM ('claimRewards', 'delegate', 'unDelegate', 'reDelegateRewards', 'withdraw', 'stake', 'unStake', 'reStakeRewards', 'unBond');
create table stake_events
(
    ste_tx_hash    varchar(64)      not null
        constraint stake_events_pk
            primary key,
    ste_type       stake_event_type not null,
    ste_validator  varchar(64)      not null,
    ste_delegator  varchar(64)      not null,
    ste_epoch      integer          not null,
    ste_amount     numeric(36, 18)  not null,
    ste_created_at timestamp        not null
);
create index stake_events_ste_created_at_index
    on stake_events (ste_created_at);
create index stake_events_ste_delegator_index
    on stake_events (ste_delegator);


create table storage
(
    stg_key   varchar(50)           not null
        constraint storage_pk
            primary key,
    stg_value text default ''::text not null
);
INSERT INTO storage (stg_key) VALUES ('stats'), ('validator_stats'), ('staking_providers'), ('nodes'), ('validators_map'), ('validators'), ('ranking');



