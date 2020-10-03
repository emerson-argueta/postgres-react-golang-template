--- Connect to TRUSTDONATIONS database

-- Creating schema for IDENTITY subdomain
DROP SCHEMA IF EXISTS IDENTITY CASCADE;
Create SCHEMA IDENTITY;

-- Creating json schema used by identity subdomain tables
CREATE EXTENSION is_jsonb_valid;
CREATE FUNCTION IDENTITY.SERVICES()
  RETURNS JSONB LANGUAGE SQL IMMUTABLE PARALLEL SAFE AS
  $$SELECT JSONB
  '{
    "type":"object",
    "patternProperties":{
      "^.+$":{
        "type":"object",
        "properties":{
          "role":{"type":"string"},
          "access":{"enum":["non-restricted","restricted"]}
        },
        "additionalProperties":false
      }
    },
    "additionalProperties":false
  }'
  $$
;

-- User
CREATE EXTENSION IF NOT EXISTS pgcrypto;
DROP TABLE IF EXISTS IDENTITY.USER;
CREATE TABLE IDENTITY.USER (
	UUID         UUID DEFAULT gen_random_uuid(),
	EMAIL        VARCHAR UNIQUE,
	PASSWORD     VARCHAR,
  SERVICES     JSONB CHECK (is_jsonb_valid(IDENTITY.SERVICES(),SERVICES))
);

-- Creating schema for CHURCH_FUND_MANAGING subdomain
DROP SCHEMA IF EXISTS CHURCH_FUND_MANAGING CASCADE;
CREATE SCHEMA CHURCH_FUND_MANAGING;

-- Create all json schemas here for CHURCH_FUND_MANAGING application
CREATE FUNCTION CHURCH_FUND_MANAGING.ADMINISTRATOR_CHURCHES()
  RETURNS JSONB LANGUAGE SQL IMMUTABLE PARALLEL SAFE AS
  $$SELECT JSONB
  '{
    "type":"object",
    "patternProperties":{
      "^[0-9]+$":{
        "type":"object",
        "properties":{
          "role":{"enum":["creator","support"]},
          "access":{"enum":["non-restricted","restricted"]}
        },
        "additionalProperties":false
      }
    },
    "additionalProperties":false
  }'
  $$
;
CREATE FUNCTION CHURCH_FUND_MANAGING.DONATOR_CHURCHES()
  RETURNS JSONB LANGUAGE SQL IMMUTABLE PARALLEL SAFE AS
  $$SELECT JSONB
  '{
    "type":"object",
    "patternProperties":{
      "^[0-9]+$":{
        "type":"object",
        "properties":{
          "donationcount":{"type":"integer"},
          "firstdonation":{"type":"string","format":"date"} 
        },
        "additionalProperties":true
      }
    },
    "additionalProperties":false
  }'
  $$
;
CREATE FUNCTION CHURCH_FUND_MANAGING.DONATORS()
  RETURNS JSONB LANGUAGE SQL IMMUTABLE PARALLEL SAFE AS
  $$SELECT JSONB
  '{
    "type":"object",
    "patternProperties":{
      "^[0-9]+$":{
        "type":"object",
        "properties":{
          "uuid":{"type":"string"}
        },
        "additionalProperties":true
      }
    },
    "additionalProperties":false
  }'
  $$
;
CREATE FUNCTION CHURCH_FUND_MANAGING.ADMINISTRATORS()
  RETURNS JSONB LANGUAGE SQL IMMUTABLE PARALLEL SAFE AS
  $$SELECT JSONB
  '{
    "type":"object",
    "patternProperties":{
      "^.+$":{
        "type":"object",
        "properties":{
          "role":{"enum":["creator","support"]},
          "access":{"enum":["non-restricted","restricted"]}
        },
        "additionalProperties":false
      }
    },
    "additionalProperties":false
  }'
  $$
;
CREATE FUNCTION CHURCH_FUND_MANAGING.ACCOUNTSTATEMENT()
  RETURNS JSONB LANGUAGE SQL IMMUTABLE PARALLEL SAFE AS
  $$SELECT JSONB
  '{
    "type":"object",
    "properties":{
      "closingbalance":{"type":"number"},
      "date":{"type":"string","format":"date"}
    },
    "additionalProperties":false
  }'
  $$
;
CREATE FUNCTION CHURCH_FUND_MANAGING.DONATIONCATEGORIES()
  RETURNS JSONB LANGUAGE SQL IMMUTABLE PARALLEL SAFE AS
  $$SELECT JSONB
  '{
    "type":"object",
    "patternProperties":{
      "^.+$":{"type":"string"}
    }
  }'
  $$
;
CREATE FUNCTION CHURCH_FUND_MANAGING.DONATION()
  RETURNS JSONB LANGUAGE SQL IMMUTABLE PARALLEL SAFE AS
  $$SELECT JSONB
  '{
    "type":"object",
    "properties":{
      "donatorid":{"type":"number"},
      "churchid":{"type":"number"},
      "amount":{"type":"number"},
      "type":{"type":"string"},
      "currency":{"type":"string"},
      "account":{"enum":["donator","church"]},
      "category":{"type":"string"},
      "details":{"type":"string"},
      "date":{"type":"string","format":"date"}
    },
    "additionalProperties":false
  }'
  $$
;
CREATE FUNCTION CHURCH_FUND_MANAGING.SUBSCRIPTION()
  RETURNS JSONB LANGUAGE SQL IMMUTABLE PARALLEL SAFE AS
  $$SELECT JSONB
  '{
    "type":"object",
    "properties":{
      "freeusagelimitcount":{"type":"integer"},
      "customeremail":{"type":"string"},
      "type":{"type":"string"},
      "paymentgateway":{"type":"object"}
    },
    "additionalProperties":false
  }'
  $$
;
-- Create all tables for CHURCH_FUND_MANAGING application
-- Admin 
DROP TABLE IF EXISTS CHURCH_FUND_MANAGING.ADMINISTRATOR;
CREATE TABLE CHURCH_FUND_MANAGING.ADMINISTRATOR (
  --ID           SERIAL,
	UUID         VARCHAR UNIQUE,
	FIRSTNAME    VARCHAR,
	LASTNAME     VARCHAR,
  ADDRESS      VARCHAR,
  PHONE        VARCHAR,
  CHURCHES     JSONB CHECK (is_jsonb_valid(CHURCH_FUND_MANAGING.ADMINISTRATOR_CHURCHES(),CHURCHES)),
  SUBSCRIPTION JSONB CHECK (is_jsonb_valid(CHURCH_FUND_MANAGING.SUBSCRIPTION(),SUBSCRIPTION))
);
-- Trigger to increment ADMIN ID
-- CREATE OR REPLACE FUNCTION ADMINISTRATOR_TRIGGER_FUNCTION()
-- RETURNS TRIGGER AS $$
-- DECLARE LAST_VAL INTEGER;
-- DECLARE NEXT_ID INTEGER;
-- BEGIN
--   SELECT last_value from CHURCH_FUND_MANAGING.ADMINISTRATOR_ID_SEQ
--   INTO LAST_VAL;
--   IF NOT NEW.ID = LAST_VAL THEN
--       SELECT NEXTVAL(PG_GET_SERIAL_SEQUENCE('CHURCH_FUND_MANAGING.ADMINISTRATOR', 'id')) AS NEW_ID
--       INTO NEXT_ID;
--       NEW.ID = NEXT_ID;
--   END IF;
--   RETURN NEW;
-- END;
-- $$ LANGUAGE 'plpgsql';
-- CREATE TRIGGER ADMINISTRATOR_TRIGGER
-- BEFORE INSERT ON CHURCH_FUND_MANAGING.ADMINISTRATOR
-- FOR EACH ROW
-- EXECUTE PROCEDURE ADMINISTRATOR_TRIGGER_FUNCTION();

-- Church 
DROP TABLE IF EXISTS CHURCH_FUND_MANAGING.CHURCH;
CREATE TABLE CHURCH_FUND_MANAGING.CHURCH (
  ID                SERIAL,
	TYPE              VARCHAR,
  NAME              VARCHAR,
  ADDRESS           VARCHAR,
  PHONE             VARCHAR,
  EMAIL             VARCHAR UNIQUE,
	PASSWORD          VARCHAR,
  ADMINISTRATORS    JSONB CHECK (is_jsonb_valid(CHURCH_FUND_MANAGING.ADMINISTRATORS(),ADMINISTRATORS)),
  DONATORS          JSONB CHECK (is_jsonb_valid(CHURCH_FUND_MANAGING.DONATORS(),DONATORS)),
  ACCOUNTSTATEMENT  JSONB CHECK (is_jsonb_valid(CHURCH_FUND_MANAGING.ACCOUNTSTATEMENT(),ACCOUNTSTATEMENT)),
  DONATIONCATEGORIES JSONB CHECK (is_jsonb_valid(CHURCH_FUND_MANAGING.DONATIONCATEGORIES(),DONATIONCATEGORIES))
);
-- Trigger to increment CHURCH ID
CREATE OR REPLACE FUNCTION CHURCH_TRIGGER_FUNCTION()
RETURNS TRIGGER AS $$
DECLARE LAST_VAL INTEGER;
DECLARE NEXT_ID INTEGER;
BEGIN
  SELECT last_value from CHURCH_FUND_MANAGING.CHURCH_ID_SEQ
  INTO LAST_VAL;
  IF NOT NEW.ID = LAST_VAL THEN
      SELECT NEXTVAL(PG_GET_SERIAL_SEQUENCE('CHURCH_FUND_MANAGING.CHURCH', 'id')) AS NEW_ID
      INTO NEXT_ID;
      NEW.ID = NEXT_ID;
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';
CREATE TRIGGER CHURCH_TRIGGER
BEFORE INSERT ON CHURCH_FUND_MANAGING.CHURCH
FOR EACH ROW
EXECUTE PROCEDURE CHURCH_TRIGGER_FUNCTION();

-- Donator 
DROP TABLE IF EXISTS CHURCH_FUND_MANAGING.DONATOR;
CREATE TABLE CHURCH_FUND_MANAGING.DONATOR (
  ID                SERIAL,
	UUID              VARCHAR,
	FIRSTNAME         VARCHAR,
	LASTNAME          VARCHAR,
  EMAIL             VARCHAR UNIQUE,
  ADDRESS           VARCHAR,
  PHONE             VARCHAR,
  CHURCHES          JSONB CHECK (is_jsonb_valid(CHURCH_FUND_MANAGING.DONATOR_CHURCHES(),CHURCHES)),
  ACCOUNTSTATEMENT  JSONB CHECK (is_jsonb_valid(CHURCH_FUND_MANAGING.ACCOUNTSTATEMENT(),ACCOUNTSTATEMENT))

);
-- Trigger to increment DONATOR ID
CREATE OR REPLACE FUNCTION DONATOR_TRIGGER_FUNCTION()
RETURNS TRIGGER AS $$
DECLARE LAST_VAL INTEGER;
DECLARE NEXT_ID INTEGER;
BEGIN
  SELECT last_value from CHURCH_FUND_MANAGING.DONATOR_ID_SEQ
  INTO LAST_VAL;
  IF NOT NEW.ID = LAST_VAL THEN
      SELECT NEXTVAL(PG_GET_SERIAL_SEQUENCE('CHURCH_FUND_MANAGING.DONATOR', 'id')) AS NEW_ID
      INTO NEXT_ID;
      NEW.ID = NEXT_ID;
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';
CREATE TRIGGER DONATOR_TRIGGER
BEFORE INSERT ON CHURCH_FUND_MANAGING.DONATOR
FOR EACH ROW
EXECUTE PROCEDURE DONATOR_TRIGGER_FUNCTION();

-- TRANSACTION 
--     Partition by CHURCHID
CREATE TYPE TRANSACTIONTYPE AS ENUM ('credit', 'debit');
CREATE TABLE CHURCH_FUND_MANAGING.TRANSACTION(
	DONATORID   INT,
  CHURCHID    INT,                     
	AMOUNT      NUMERIC(15,6),
  TYPE        TRANSACTIONTYPE, 
	Donation    JSONB CHECK (is_jsonb_valid(CHURCH_FUND_MANAGING.DONATION(),Donation)),
	CREATEDAT   DATE,
  UPDATEDAT   DATE
);