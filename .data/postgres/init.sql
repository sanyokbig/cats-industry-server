create table characters
(
  id      serial  not null
    constraint characters_pkey
    primary key,
  name    text    not null,
  is_main boolean not null
);

create table users
(
  id serial not null
    constraint users_pkey
    primary key
);

create table groups
(
  id   serial not null
    constraint groups_pkey
    primary key,
  name text   not null
);

create table roles
(
  id   serial not null
    constraint roles_pkey
    primary key,
  name text   not null
);

create table tokens
(
  id            serial  not null
    constraint tokens_pkey
    primary key,
  character_id  integer
    constraint tokens_character_id_fkey
    references characters
    on update cascade,
  expires_at    integer not null,
  scopes        text    not null,
  access_token  text    not null,
  refresh_token text    not null,
  constraint tokens_character_id_scopes_key
  unique (character_id, scopes)
);

create table users_characters
(
  user_id      integer not null
    constraint users_characters_user_id_fkey
    references users
    on update cascade,
  character_id integer not null
    constraint users_characters_character_id_key
    unique
    constraint users_characters_character_id_fkey
    references characters
    on update cascade,
  constraint users_characters_pkey
  primary key (user_id, character_id)
);

create table users_groups
(
  user_id  integer not null
    constraint users_groups_user_id_fkey
    references users
    on update cascade,
  group_id integer not null
    constraint users_groups_group_id_fkey
    references groups
    on update cascade,
  constraint users_groups_pkey
  primary key (user_id, group_id)
);

create table groups_roles
(
  group_id integer not null
    constraint groups_roles_group_id_fkey
    references groups
    on update cascade,
  role_id  integer not null
    constraint groups_roles_role_id_fkey
    references roles
    on update cascade,
  constraint groups_roles_pkey
  primary key (group_id, role_id)
);

create table sde_ram_activities
(
  id          bigint not null
    constraint ram_activities_pkey
    primary key,
  name        text,
  description text,
  icon        varchar(10)
);

create unique index ram_activities_id_uindex
  on sde_ram_activities (id);

create table sde_product_types
(
  id   bigint not null
    constraint table_name_pkey
    primary key,
  name varchar(160)
);

create table jobs
(
  id                     serial           not null
    constraint jobs_pkey
    primary key,
  eve_id                 integer          not null
    constraint jobs_eve_id_key
    unique,
  installer_id           integer          not null,
  facility_id            bigint           not null,
  station_id             bigint           not null,
  activity_id            integer          not null
    constraint jobs_ram_activities_id_fk
    references sde_ram_activities,
  blueprint_id           bigint           not null,
  blueprint_type_id      integer          not null,
  blueprint_location_id  bigint           not null,
  output_location_id     bigint           not null,
  runs                   smallint         not null,
  cost                   double precision not null,
  licensed_runs          smallint         not null,
  probability            double precision not null,
  product_type_id        integer          not null
    constraint jobs_product_types_id_fk
    references sde_product_types,
  status                 text             not null,
  duration               integer          not null,
  start_date             bigint           not null,
  end_date               bigint           not null,
  pause_date             bigint           not null,
  completed_date         bigint           not null,
  completed_character_id integer          not null,
  successful_runs        smallint         not null
);

create unique index table_name_id_uindex
  on sde_product_types (id);

