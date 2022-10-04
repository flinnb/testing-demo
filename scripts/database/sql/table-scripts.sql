CREATE TABLE user_profiles (
    id SERIAL PRIMARY KEY,
    first_name character varying(255),
    last_name character varying(255),
    username character varying(255) UNIQUE,
    city character varying(255),
    zip_code character varying(10),
    created_on timestamp NOT NULL,
    updated_on timestamp NOT NULL
);

CREATE TABLE user_password_history (
    id SERIAL PRIMARY KEY,
    user_id integer REFERENCES user_profiles ON DELETE CASCADE,
    password_hash character varying(255),
    created_on timestamp NOT NULL,
    is_active boolean
);
