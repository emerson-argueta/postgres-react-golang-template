--- Connect to TRUSTDONATIONS database

-- Creating schema for IDENTITY subdomain
DROP SCHEMA IF EXISTS IDENTITY CASCADE;
Create SCHEMA IDENTITY;

-- Creating services json schema used by identity models
-- object structure {[service_id_number:number]:{role:string}}
-- ex: {1432:{role:administrator}}
--    in this example id 1432 could be service with name communit_goal_tracker
CREATE EXTENSION is_jsonb_valid;
CREATE FUNCTION IDENTITY.SERVICES()
  RETURNS JSONB LANGUAGE SQL IMMUTABLE PARALLEL SAFE AS
  $$SELECT JSONB
  '{
    "type":"object",
    "patternProperties":{
      "^[0-9]+$":{
        "type":"object",
        "properties":{
          "role":{"type":"string"},
        },
        "additionalProperties":false
      }
    },
    "additionalProperties":false
  }'
  $$
;

-- User model
CREATE EXTENSION IF NOT EXISTS pgcrypto;
DROP TABLE IF EXISTS IDENTITY.USER;
CREATE TABLE IDENTITY.USER (
	UUID         UUID DEFAULT gen_random_uuid(),
	EMAIL        VARCHAR UNIQUE,
	PASSWORD     VARCHAR,
  SERVICES     JSONB CHECK (is_jsonb_valid(IDENTITY.SERVICES(),SERVICES))
);

-- Service model
DROP TABLE IF EXISTS IDENTITY.SERVICE;
CREATE TABLE IDENTITY.SERVICE (
	ID           SERIAL,
	NAME         VARCHAR UNIQUE
);
-- Trigger to increment SERVICE ID
CREATE OR REPLACE FUNCTION SERVICE_TRIGGER_FUNCTION()
RETURNS TRIGGER AS $$
DECLARE LAST_VAL INTEGER;
DECLARE NEXT_ID INTEGER;
BEGIN
  SELECT last_value from IDENTITY.SERVICE_ID_SEQ
  INTO LAST_VAL;
  IF NOT NEW.ID = LAST_VAL THEN
      SELECT NEXTVAL(PG_GET_SERIAL_SEQUENCE('IDENTITY.SERVICE', 'id')) AS NEW_ID
      INTO NEXT_ID;
      NEW.ID = NEXT_ID;
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';
CREATE TRIGGER SERVICE_TRIGGER
BEFORE INSERT ON IDENTITY.SERVICE
FOR EACH ROW
EXECUTE PROCEDURE SERVICE_TRIGGER_FUNCTION();

-- Creating schema for COMMUNITY_GOAL_TRACKER subdomain
DROP SCHEMA IF EXISTS COMMUNITY_GOAL_TRACKER CASCADE;
CREATE SCHEMA COMMUNITY_GOAL_TRACKER;

-- Creating achievers json schema used by community_goal_tracker models
-- object structure {[achiever_id_number:number]:{progress:string,messages:{[message_number:number]:string}}}
-- ex: {1432:{progress:50,messages:{0:'layed down the concrete for the foundation'}}}
--     in this example achiever_id_numer 1432 is the achiever id 
CREATE FUNCTION COMMUNITY_GOAL_TRACKER.ACHIEVERS()
  RETURNS JSONB LANGUAGE SQL IMMUTABLE PARALLEL SAFE AS
  $$SELECT JSONB
  '{
    "type":"object",
    "patternProperties":{
      "^[0-9]+$":{
        "type":"object",
        "properties":{
          "state":{"enum":["abandoned","inprogress","completed"]},
          "progress":{"type":"integer"},
          "messages":{"type":"object","patternProperties":{"^[0-9]+$":{"type":"string"}}}
        },
        "additionalProperties":false
      }
    },
    "additionalProperties":false
  }'
  $$
;

-- Creating goals json schema used by community_goal_tracker models
-- object structure {[goals_id_number]}
-- ex: {[1432]}}
--     in this example goals_id_numer 1432 could refer to goal with name bulid-house
CREATE FUNCTION COMMUNITY_GOAL_TRACKER.GOALS()
  RETURNS JSONB LANGUAGE SQL IMMUTABLE PARALLEL SAFE AS
  $$SELECT JSONB
  '{
    "type":"array",
    "items":{
      "type":"number"
    }
  }'
  $$
;

-- Create all tables for COMMUNITY_GOAL_TRACKER domain
-- Achiever 
DROP TABLE IF EXISTS COMMUNITY_GOAL_TRACKER.ACHIEVER;
CREATE TABLE COMMUNITY_GOAL_TRACKER.ACHIEVER (
	UUID         VARCHAR UNIQUE,
	FIRSTNAME    VARCHAR,
	LASTNAME     VARCHAR,
  ADDRESS      VARCHAR,
  PHONE        VARCHAR,
  GOALS        JSONB CHECK (is_jsonb_valid(COMMUNITY_GOAL_TRACKER.GOALS(),GOALS))
);

-- Goal 
DROP TABLE IF EXISTS COMMUNITY_GOAL_TRACKER.GOAL;
CREATE TABLE COMMUNITY_GOAL_TRACKER.GOAL (
  ID                SERIAL,
  NAME              VARCHAR UNIQUE,
  ACHIEVERS    JSONB CHECK (is_jsonb_valid(COMMUNITY_GOAL_TRACKER.ACHIEVERS(),ACHIEVERS))
);
-- Trigger to increment CHURCH ID
CREATE OR REPLACE FUNCTION GOAL_TRIGGER_FUNCTION()
RETURNS TRIGGER AS $$
DECLARE LAST_VAL INTEGER;
DECLARE NEXT_ID INTEGER;
BEGIN
  SELECT last_value from COMMUNITY_GOAL_TRACKER.GOAL_ID_SEQ
  INTO LAST_VAL;
  IF NOT NEW.ID = LAST_VAL THEN
      SELECT NEXTVAL(PG_GET_SERIAL_SEQUENCE('COMMUNITY_GOAL_TRACKER.GOAL', 'id')) AS NEW_ID
      INTO NEXT_ID;
      NEW.ID = NEXT_ID;
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';
CREATE TRIGGER GOAL_TRIGGER
BEFORE INSERT ON COMMUNITY_GOAL_TRACKER.GOAL
FOR EACH ROW
EXECUTE PROCEDURE GOAL_TRIGGER_FUNCTION();