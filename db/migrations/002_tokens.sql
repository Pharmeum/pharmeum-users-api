-- +migrate Up

CREATE TABLE tokens(
  email varchar(256) not null,
  token varchar(128) PRIMARY KEY ,
  last_sent_at timestamp without time zone
);

-- +migrate Down

DROP TABLE tokens;