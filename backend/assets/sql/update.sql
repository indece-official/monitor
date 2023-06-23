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

--- Update before rev. 3 are not supported
DO $$
BEGIN
    IF (SELECT "value" FROM "mo_dbinfo" WHERE "name" = 'revision') != '1' THEN
        RETURN;
    END IF;

    RAISE EXCEPTION 'Automatic db updates before db revision 3 are not supported';
END $$;

--- Update before rev. 3 are not supported
DO $$
BEGIN
    IF (SELECT "value" FROM "mo_dbinfo" WHERE "name" = 'revision') != '2' THEN
        RETURN;
    END IF;

    RAISE EXCEPTION 'Automatic db updates before db revision 3 are not supported';
END $$;

--- Update from rev. 3 to rev. 4
DO $$
DECLARE 
    check_row mo_check%rowtype;
    old_checker_row mo_checker%rowtype;
    agent_row mo_agent%rowtype;
    updated_checker_row mo_checker%rowtype;
    new_checker_uid varchar;
BEGIN
    IF (SELECT "value" FROM "mo_dbinfo" WHERE "name" = 'revision') != '3' THEN
        RETURN;
    END IF;

    -- add column agent_uid to checker (nullable)
    -- make column agent_type in checker nullable
    -- for each check: get checker with agent_type, get agent by agent_type and check.host_uid, add new checker, update check 
    -- delete all checkers with agent_uid null
    -- remove column checker.agent_type
    -- set column checker.agent_uid not nullable
    DROP VIEW mo_checker_v1;
    DROP VIEW mo_check_v1;

    ALTER TABLE "mo_checker"
        ADD COLUMN "agent_uid" VARCHAR(36) NULL REFERENCES mo_agent(uid);

    ALTER TABLE "mo_checker"
        ALTER COLUMN "agent_type" DROP NOT NULL;

    FOR check_row in SELECT * FROM mo_check LOOP
        SELECT
        INTO old_checker_row *
        FROM mo_checker
        WHERE
            mo_checker.uid = check_row.checker_uid
        LIMIT 1;

        SELECT
        INTO agent_row *
        FROM mo_agent
        WHERE
            mo_agent.type = old_checker_row.agent_type AND
            mo_agent.host_uid = check_row.host_uid
        LIMIT 1;

        SELECT
        INTO updated_checker_row *
        FROM mo_checker
        WHERE
            mo_checker.type = old_checker_row.type AND
            mo_checker.agent_uid = agent_row.uid
        LIMIT 1;

        IF updated_checker_row IS NOT NULL THEN
            UPDATE mo_check SET
                checker_uid = updated_checker_row.uid,
                datetime_updated = NOW()
            WHERE
                uid = check_row.uid;
        ELSE 
            new_checker_uid := gen_random_uuid();

            INSERT INTO mo_checker
                (uid, type, agent_uid, version, name, custom_checks, capabilities)
            VALUES 
                (new_checker_uid, old_checker_row.type, agent_row.uid, old_checker_row.version, old_checker_row.name, old_checker_row.custom_checks, old_checker_row.capabilities);

            UPDATE mo_check SET
                checker_uid = new_checker_uid,
                datetime_updated = NOW()
            WHERE
                uid = check_row.uid;
        END IF;
    END LOOP;

    DELETE FROM mo_checker WHERE agent_uid IS NULL;

    ALTER TABLE "mo_checker"
        DROP COLUMN "agent_type";

    ALTER TABLE "mo_checker"
        ALTER COLUMN "agent_uid" SET NOT NULL;

    ALTER TABLE "mo_check"
        DROP COLUMN "host_uid";

    CREATE VIEW mo_checker_v1 AS SELECT * FROM mo_checker;
    CREATE VIEW mo_check_v1 AS SELECT * FROM mo_check;

    UPDATE "mo_dbinfo" SET "value" = '4' WHERE "name" = 'revision';
END $$;
