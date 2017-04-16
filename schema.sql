create table post (
  id bigint unsigned not null primary key,
  created_at datetime not null,
  updated_at datetime not null,
  url_name text not null,
  title text not null,
  body text not null
);

create table last_id (
  post_last_id bigint not null default 0
);

create table label (
  post_id bigint not null,
  name text
);
