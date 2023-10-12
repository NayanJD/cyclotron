CREATE TABLE users
 ( 
     id                  uuid PRIMARY KEY, 
     username            text  NOT NULL UNIQUE, 
     first_name          text  NOT NULL, 
     last_name           text  NOT NULL, 
     hashed_password     text NOT NULL, 
     dob                 timestamptz NOT NULL,  
    
     created_at          timestamptz NOT NULL, 
     updated_at          timestamptz NOT NULL, 
     deleted_at          timestamptz 
);

CREATE TABLE tokens
(
    id                  serial PRIMARY KEY,
    refresh_token       text NOT NULL UNIQUE,
    access_token        text NOT NULL UNIQUE,
    user_id             text  NOT NULL,
    
    created_at          timestamptz NOT NULL,
    valid_till          timestamptz NOT NULL
)

/* SELECT sum(pg_column_size(t.*)) as filesize, count(*) as filerow FROM users as t where id = (select id from users limit 1); */
