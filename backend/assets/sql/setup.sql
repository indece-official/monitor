--- indece Monitor
--- Copyright (C) 2023 indece UG (haftungsbeschränkt)
---
--- This program is free software: you can redistribute it and/or modify
--- it under the terms of the GNU General Public License as published by
--- the Free Software Foundation, either version 3 of the License or any
--- later version.
---
--- This program is distributed in the hope that it will be useful,
--- but WITHOUT ANY WARRANTY; without even the implied warranty of
--- MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
--- GNU General Public License for more details.
---
--- You should have received a copy of the GNU General Public License
--- along with this program. If not, see <https:--www.gnu.org/licenses/>.

--- Default table for storing database informations (revision, ...)
CREATE TABLE IF NOT EXISTS mo_dbinfo (
    name VARCHAR (64) UNIQUE NOT NULL,
    value INT NOT NULL
);

INSERT INTO mo_dbinfo(name, value)
    VALUES
        ('revision', 4)
    ON CONFLICT DO NOTHING;

--- Table for storing config properties
CREATE TABLE IF NOT EXISTS mo_configproperty (
    key VARCHAR(256) NOT NULL,
    value TEXT NOT NULL,
    datetime_created TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    datetime_updated TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS mo_configproperty_key
    ON mo_configproperty(key);

DO $$
BEGIN
   IF NOT EXISTS (
        SELECT *
        FROM information_schema.tables
        WHERE table_type='VIEW' AND
        table_schema='public'
        AND table_catalog=current_database()
        AND table_name='mo_configproperty_v1') THEN
        CREATE VIEW mo_configproperty_v1 AS SELECT * FROM mo_configproperty;
    END IF;
END $$;

--- Table for storing users
CREATE TABLE IF NOT EXISTS mo_user (
    uid VARCHAR(36) NOT NULL,
    source VARCHAR(32) NOT NULL,
    username VARCHAR(256) NOT NULL,
    name VARCHAR(256) NULL,
    email VARCHAR(256) NULL,
    local_roles VARCHAR(256) ARRAY NOT NULL,
    password_hash VARCHAR(256) NULL,
    datetime_created TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    datetime_updated TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    datetime_locked TIMESTAMP WITH TIME ZONE NULL DEFAULT NULL,
    datetime_deleted TIMESTAMP WITH TIME ZONE NULL DEFAULT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS mo_user_uid
    ON mo_user(uid);

CREATE UNIQUE INDEX IF NOT EXISTS mo_user_username
    ON mo_user(username)
    WHERE
        mo_user.datetime_deleted IS NULL;

DO $$
BEGIN
   IF NOT EXISTS (
        SELECT *
        FROM information_schema.tables
        WHERE table_type='VIEW' AND
        table_schema='public'
        AND table_catalog=current_database()
        AND table_name='mo_user_v1') THEN
        CREATE VIEW mo_user_v1 AS SELECT * FROM mo_user;
    END IF;
END $$;

--- Table for storing host tags
CREATE TABLE IF NOT EXISTS mo_tag (
    uid VARCHAR(36) NOT NULL,
    name VARCHAR(256) NOT NULL,
    color VARCHAR(32) NOT NULL,
    datetime_created TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    datetime_updated TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    datetime_deleted TIMESTAMP WITH TIME ZONE NULL DEFAULT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS mo_tag_uid
    ON mo_tag(uid);

DO $$
BEGIN
   IF NOT EXISTS (
        SELECT *
        FROM information_schema.tables
        WHERE table_type='VIEW' AND
        table_schema='public'
        AND table_catalog=current_database()
        AND table_name='mo_tag_v1') THEN
        CREATE VIEW mo_tag_v1 AS SELECT * FROM mo_tag;
    END IF;
END $$;

--- Table for storing hosts
CREATE TABLE IF NOT EXISTS mo_host (
    uid VARCHAR(36) NOT NULL,
    name VARCHAR(256) NOT NULL,
    tag_uids VARCHAR(36) ARRAY NOT NULL,
    datetime_created TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    datetime_updated TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    datetime_deleted TIMESTAMP WITH TIME ZONE NULL DEFAULT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS mo_host_uid
    ON mo_host(uid);

DO $$
BEGIN
   IF NOT EXISTS (
        SELECT *
        FROM information_schema.tables
        WHERE table_type='VIEW' AND
        table_schema='public'
        AND table_catalog=current_database()
        AND table_name='mo_host_v1') THEN
        CREATE VIEW mo_host_v1 AS SELECT * FROM mo_host;
    END IF;
END $$;

--- Table for storing agents
CREATE TABLE IF NOT EXISTS mo_agent (
    uid VARCHAR(36) NOT NULL,
    host_uid VARCHAR(36) NOT NULL REFERENCES mo_host(uid),
    type VARCHAR(256) NULL,
    version VARCHAR(32) NULL,
    certs JSON NOT NULL,
    datetime_created TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    datetime_updated TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    datetime_registered TIMESTAMP WITH TIME ZONE NULL DEFAULT NULL,
    datetime_deleted TIMESTAMP WITH TIME ZONE NULL DEFAULT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS mo_agent_uid
    ON mo_agent(uid);

DO $$
BEGIN
   IF NOT EXISTS (
        SELECT *
        FROM information_schema.tables
        WHERE table_type='VIEW' AND
        table_schema='public'
        AND table_catalog=current_database()
        AND table_name='mo_agent_v1') THEN
        CREATE VIEW mo_agent_v1 AS SELECT * FROM mo_agent;
    END IF;
END $$;

--- Table for storing checkers
CREATE TABLE IF NOT EXISTS mo_checker (
    uid VARCHAR(36) NOT NULL,
    type VARCHAR(256) NOT NULL,
    agent_uid VARCHAR(36) NOT NULL REFERENCES mo_agent(uid),
    version VARCHAR(32) NULL,
    name VARCHAR(256) NOT NULL,
    custom_checks BOOLEAN NOT NULL,
    capabilities JSON NOT NULL,
    datetime_created TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    datetime_updated TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    datetime_deleted TIMESTAMP WITH TIME ZONE NULL DEFAULT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS mo_checker_uid
    ON mo_checker(uid);

DO $$
BEGIN
   IF NOT EXISTS (
        SELECT *
        FROM information_schema.tables
        WHERE table_type='VIEW' AND
        table_schema='public'
        AND table_catalog=current_database()
        AND table_name='mo_checker_v1') THEN
        CREATE VIEW mo_checker_v1 AS SELECT * FROM mo_checker;
    END IF;
END $$;

--- Table for storing checks
CREATE TABLE IF NOT EXISTS mo_check (
    uid VARCHAR(36) NOT NULL,
    checker_uid VARCHAR(36) NOT NULL REFERENCES mo_checker(uid),
    name VARCHAR(256) NOT NULL,
    type VARCHAR(256) NULL,
    schedule VARCHAR(64) NULL,
    config JSON NOT NULL,
    custom BOOLEAN NOT NULL,
    datetime_created TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    datetime_updated TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    datetime_disabled TIMESTAMP WITH TIME ZONE NULL DEFAULT NULL,
    datetime_deleted TIMESTAMP WITH TIME ZONE NULL DEFAULT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS mo_check_uid
    ON mo_check(uid);

DO $$
BEGIN
   IF NOT EXISTS (
        SELECT *
        FROM information_schema.tables
        WHERE table_type='VIEW' AND
        table_schema='public'
        AND table_catalog=current_database()
        AND table_name='mo_check_v1') THEN
        CREATE VIEW mo_check_v1 AS SELECT * FROM mo_check;
    END IF;
END $$;

--- Table for storing checkstatuses
CREATE TABLE IF NOT EXISTS mo_checkstatus (
    uid VARCHAR(36) NOT NULL,
    check_uid VARCHAR(36) NOT NULL REFERENCES mo_check(uid),
    status VARCHAR(16) NOT NULL,
    message TEXT NOT NULL,
    data JSON NOT NULL,
    datetime_created TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    datetime_updated TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    datetime_deleted TIMESTAMP WITH TIME ZONE NULL DEFAULT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS mo_checkstatus_uid
    ON mo_checkstatus(uid);

CREATE UNIQUE INDEX IF NOT EXISTS mo_checkstatus_check_uid_datetime_created
    ON mo_checkstatus(check_uid, datetime_created DESC);

DO $$
BEGIN
   IF NOT EXISTS (
        SELECT *
        FROM information_schema.tables
        WHERE table_type='VIEW' AND
        table_schema='public'
        AND table_catalog=current_database()
        AND table_name='mo_checkstatus_v1') THEN
        CREATE VIEW mo_checkstatus_v1 AS SELECT * FROM mo_checkstatus;
    END IF;
END $$;

--- Table for storing notifiers
CREATE TABLE IF NOT EXISTS mo_notifier (
    uid VARCHAR(36) NOT NULL,
    name VARCHAR(256) NOT NULL,
    type VARCHAR(64) NOT NULL,
    config JSON NOT NULL,
    datetime_created TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    datetime_updated TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    datetime_deleted TIMESTAMP WITH TIME ZONE NULL DEFAULT NULL,
    datetime_disabled TIMESTAMP WITH TIME ZONE NULL DEFAULT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS mo_notifier_uid
    ON mo_notifier(uid);

DO $$
BEGIN
   IF NOT EXISTS (
        SELECT *
        FROM information_schema.tables
        WHERE table_type='VIEW' AND
        table_schema='public'
        AND table_catalog=current_database()
        AND table_name='mo_notifier_v1') THEN
        CREATE VIEW mo_notifier_v1 AS SELECT * FROM mo_notifier;
    END IF;
END $$;

--- Table for storing maintenance
CREATE TABLE IF NOT EXISTS mo_maintenance (
    uid VARCHAR(36) NOT NULL,
    title VARCHAR(256) NOT NULL,
    message TEXT NOT NULL,
    details JSON NOT NULL,
    datetime_created TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    datetime_updated TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    datetime_start TIMESTAMP WITH TIME ZONE NOT NULL,
    datetime_finish TIMESTAMP WITH TIME ZONE NULL DEFAULT NULL,
    datetime_deleted TIMESTAMP WITH TIME ZONE NULL DEFAULT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS mo_maintenance_uid
    ON mo_maintenance(uid);

DO $$
BEGIN
   IF NOT EXISTS (
        SELECT *
        FROM information_schema.tables
        WHERE table_type='VIEW' AND
        table_schema='public'
        AND table_catalog=current_database()
        AND table_name='mo_maintenance_v1') THEN
        CREATE VIEW mo_maintenance_v1 AS SELECT * FROM mo_maintenance;
    END IF;
END $$;
