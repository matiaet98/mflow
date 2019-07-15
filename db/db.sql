--------------------------------------------------------
--  File created - Monday-July-15-2019   
--------------------------------------------------------
--------------------------------------------------------
--  DDL for Sequence SEQ_TASKS_MASTER
--------------------------------------------------------

   CREATE SEQUENCE  "SIPER"."SEQ_TASKS_MASTER"  MINVALUE 1 MAXVALUE 9999999999999999999999999999 INCREMENT BY 1 START WITH 141 CACHE 20 ORDER  NOCYCLE  NOKEEP  NOSCALE  GLOBAL ;
--------------------------------------------------------
--  DDL for Table TASKS_MASTER
--------------------------------------------------------

  CREATE TABLE "SIPER"."TASKS_MASTER" 
   (	"ID" NUMBER, 
	"TASK_ID" NUMBER, 
	"START_DATE" DATE, 
	"END_DATE" DATE, 
	"STATUS" VARCHAR2(20)
   ) ;

   COMMENT ON COLUMN "SIPER"."TASKS_MASTER"."STATUS" IS 'SUCCESS / FAILED / RUNNING / NONE';
--------------------------------------------------------
--  DDL for Index LOG_CORRIDAS_PK
--------------------------------------------------------

  CREATE UNIQUE INDEX "SIPER"."LOG_CORRIDAS_PK" ON "SIPER"."TASKS_MASTER" ("ID") 
  ;
--------------------------------------------------------
--  DDL for Trigger TRIGGER1
--------------------------------------------------------

  CREATE OR REPLACE EDITIONABLE TRIGGER "SIPER"."TRIGGER1" 
BEFORE INSERT ON TASKS_MASTER FOR EACH ROW
BEGIN
  :NEW.ID := SEQ_TASKS_MASTER.NEXTVAL;
END;
/
ALTER TRIGGER "SIPER"."TRIGGER1" ENABLE;
--------------------------------------------------------
--  DDL for Package PKG_TASKMAN
--------------------------------------------------------

  CREATE OR REPLACE EDITIONABLE PACKAGE "SIPER"."PKG_TASKMAN" AS 

  PROCEDURE START_TASK(I_TASK_ID TASKS_MASTER.TASK_ID%TYPE);
  PROCEDURE UPDATE_TASK(I_TASK_ID  IN TASKS_MASTER.TASK_ID%TYPE,I_STATUS IN TASKS_MASTER.STATUS%TYPE);
  PROCEDURE GET_STATUS(I_TASK_ID  IN TASKS_MASTER.TASK_ID%TYPE,O_STATUS OUT TASKS_MASTER.STATUS%TYPE, O_START_DATE OUT TASKS_MASTER.START_DATE%TYPE);

END PKG_TASKMAN;

/
--------------------------------------------------------
--  DDL for Package Body PKG_TASKMAN
--------------------------------------------------------

  CREATE OR REPLACE EDITIONABLE PACKAGE BODY "SIPER"."PKG_TASKMAN" AS

  PROCEDURE START_TASK(I_TASK_ID TASKS_MASTER.TASK_ID%TYPE) AS
  BEGIN
    INSERT INTO TASKS_MASTER(TASK_ID,START_DATE,END_DATE,STATUS)
    VALUES(I_TASK_ID,SYSDATE,NULL,'RUNNING');
  END START_TASK;

  PROCEDURE UPDATE_TASK(I_TASK_ID  IN TASKS_MASTER.TASK_ID%TYPE,I_STATUS IN TASKS_MASTER.STATUS%TYPE) AS
  BEGIN
    UPDATE TASKS_MASTER
    SET STATUS = I_STATUS,END_DATE = SYSDATE
    WHERE TASK_ID = I_TASK_ID
    AND END_DATE IS NULL;
  END UPDATE_TASK;

  PROCEDURE GET_STATUS(I_TASK_ID  IN TASKS_MASTER.TASK_ID%TYPE,O_STATUS OUT TASKS_MASTER.STATUS%TYPE, O_START_DATE OUT TASKS_MASTER.START_DATE%TYPE) IS
  BEGIN
        SELECT STATUS,START_DATE
        INTO O_STATUS,O_START_DATE
        FROM TASKS_MASTER
        WHERE TASK_ID = I_TASK_ID
        AND ID = (SELECT MAX(ID) FROM TASKS_MASTER WHERE TASK_ID = I_TASK_ID);
        
  EXCEPTION WHEN OTHERS THEN
    O_STATUS := 'NONE';
    O_START_DATE := NULL;
  END GET_STATUS;

END PKG_TASKMAN;

/
--------------------------------------------------------
--  DDL for Function FNC_TEST
--------------------------------------------------------

  CREATE OR REPLACE EDITIONABLE FUNCTION "SIPER"."FNC_TEST" 
(
  I_P1 IN NUMBER 
) RETURN VARCHAR2 AS 
BEGIN
  RETURN 'Hello Fucking idiot';
END FNC_TEST;

/
--------------------------------------------------------
--  Constraints for Table TASKS_MASTER
--------------------------------------------------------

  ALTER TABLE "SIPER"."TASKS_MASTER" MODIFY ("ID" NOT NULL ENABLE);
  ALTER TABLE "SIPER"."TASKS_MASTER" MODIFY ("TASK_ID" NOT NULL ENABLE);
  ALTER TABLE "SIPER"."TASKS_MASTER" MODIFY ("START_DATE" NOT NULL ENABLE);
  ALTER TABLE "SIPER"."TASKS_MASTER" MODIFY ("STATUS" NOT NULL ENABLE);
  ALTER TABLE "SIPER"."TASKS_MASTER" ADD CONSTRAINT "LOG_CORRIDAS_PK" PRIMARY KEY ("ID")
  USING INDEX  ENABLE;
