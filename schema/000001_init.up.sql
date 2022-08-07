BEGIN;

SET statement_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = ON;
SET check_function_bodies = FALSE;
SET client_min_messages = WARNING;
SET search_path = public, extensions;
SET default_tablespace = '';
SET default_with_oids = FALSE;

CREATE TABLE person (
    id VARCHAR(100) PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    age INT,
    created_at    TIMESTAMPTZ,
    updated_at    TIMESTAMPTZ
);

CREATE TABLE device (
    id VARCHAR(100) PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    age INT,
    category TEXT,
    created_at    TIMESTAMPTZ,
    updated_at    TIMESTAMPTZ
);

CREATE TABLE person_device (
    id VARCHAR(100) PRIMARY KEY DEFAULT gen_random_uuid(),
    device_id VARCHAR(100) NOT NULL,
    person_id VARCHAR(100) NOT NULL,
    created_at    TIMESTAMPTZ,
    updated_at    TIMESTAMPTZ,

    CONSTRAINT device_fk FOREIGN KEY (device_id)  REFERENCES device (id),
    CONSTRAINT person_fk  FOREIGN KEY  (person_id)  REFERENCES person (id),
    CONSTRAINT device_person_unique UNIQUE (device_id, person_id)
);

INSERT INTO person (id, name, age)
VALUES ('7c9628d7-3151-49fa-b5e3-4', 'Oleg', '50');
INSERT INTO person (id, name, age)
VALUES ('2aa33c62-8c48-4c84-afee-', 'Nika', '30');
INSERT INTO person (id, name, age)
VALUES ('fef50e36-1e5e-4e23-9d2e', 'Alex', '25');

INSERT INTO device (id, name, age, category)
VALUES ('b730d758-b6f2-415f-a13b', 'car', '5', 'for road');
INSERT INTO device (id, name, age, category)
VALUES ('9e5b590d-0347-4dda-a44f-6c76', 'train', '20', 'for railways');
INSERT INTO device (id, name, age, category)
VALUES ('90fc6524-789c-4c1f-bd59-f9', 'airplane', '10', 'for air');

INSERT INTO person_device (person_id, device_id)
VALUES ('7c9628d7-3151-49fa-b5e3-4', '90fc6524-789c-4c1f-bd59-f9');
INSERT INTO person_device (person_id, device_id)
VALUES ('2aa33c62-8c48-4c84-afee-', 'b730d758-b6f2-415f-a13b');
INSERT INTO person_device (person_id, device_id)
VALUES ('fef50e36-1e5e-4e23-9d2e', '9e5b590d-0347-4dda-a44f-6c76');