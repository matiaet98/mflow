CREATE SEQUENCE  "SIPER"."SEQ_TASKS_MASTER"  MINVALUE 1 MAXVALUE 9999999999999999999999999999 INCREMENT BY 1 START WITH 1 CACHE 20 ORDER  NOCYCLE  NOKEEP  NOSCALE  GLOBAL;

CREATE SEQUENCE  "SIPER"."SEQ_TASKS"  MINVALUE 1 MAXVALUE 9999999999999999999999999999 INCREMENT BY 1 START WITH 1 CACHE 20 ORDER  NOCYCLE  NOKEEP  NOSCALE  GLOBAL;

CREATE TABLE "SIPER"."TASKS_MASTER"(
    "ID" NUMBER,
    "START_DATE" DATE,
	"END_DATE" DATE,
	"STATUS" VARCHAR2(20)
);
COMMENT ON COLUMN "SIPER"."TASKS_MASTER"."STATUS" IS 'STARTED / ENDED';
CREATE UNIQUE INDEX "SIPER"."ID_TASKS_MASTER_PK1" ON "SIPER"."TASKS_MASTER" ("ID");

CREATE TABLE "SIPER"."TASKS" 
   ("ID_MASTER" NUMBER,
	"ID_TASK" NUMBER, 
	"START_DATE" DATE, 
	"END_DATE" DATE, 
	"STATUS" VARCHAR2(20)
   ) ;
COMMENT ON COLUMN "SIPER"."TASKS"."STATUS" IS 'SUCCESS / FAILED / RUNNING / NONE';
CREATE UNIQUE INDEX "SIPER"."TASKS_PK1" ON "SIPER"."TASKS" ("ID_MASTER","ID_TASK");


ALTER TABLE "SIPER"."TASKS_MASTER" MODIFY ("ID" NOT NULL ENABLE);
ALTER TABLE "SIPER"."TASKS_MASTER" MODIFY ("START_DATE" NOT NULL ENABLE);
ALTER TABLE "SIPER"."TASKS_MASTER" MODIFY ("STATUS" NOT NULL ENABLE);

ALTER TABLE "SIPER"."TASKS" MODIFY ("ID_MASTER" NOT NULL ENABLE);
ALTER TABLE "SIPER"."TASKS" MODIFY ("ID_TASK" NOT NULL ENABLE);
ALTER TABLE "SIPER"."TASKS" MODIFY ("START_DATE" NOT NULL ENABLE);
ALTER TABLE "SIPER"."TASKS" MODIFY ("STATUS" NOT NULL ENABLE);
/

create or replace PACKAGE BODY "PKG_TASKMAN" AS

PROCEDURE CREATE_MASTER(I_ID OUT TASKS_MASTER.ID%TYPE) AS
BEGIN

I_ID := SEQ_TASKS_MASTER.NEXTVAL;

INSERT INTO TASKS_MASTER(ID,START_DATE,END_DATE,STATUS)
VALUES(I_ID,SYSDATE,NULL,'STARTED');

END CREATE_MASTER;

FUNCTION GET_MASTER_STATUS(I_ID IN TASKS_MASTER.ID%TYPE) RETURN TASKS_MASTER.STATUS%TYPE IS
L_STATUS TASKS_MASTER.STATUS%TYPE;
BEGIN

SELECT STATUS
INTO L_STATUS
FROM TASKS_MASTER
WHERE ID = I_ID;

RETURN L_STATUS;

END GET_MASTER_STATUS;

PROCEDURE END_MASTER(I_ID IN TASKS_MASTER.ID%TYPE) IS

BEGIN

UPDATE TASKS_MASTER
SET STATUS = 'ENDED',END_DATE = SYSDATE
WHERE ID = I_ID;

END END_MASTER;


PROCEDURE START_TASK(I_ID_MASTER IN TASKS.ID_MASTER%TYPE,I_ID_TASK IN TASKS.ID_TASK%TYPE) AS
BEGIN
    INSERT INTO TASKS(ID_MASTER,ID_TASK,START_DATE,END_DATE,STATUS)
    VALUES (I_ID_MASTER,I_ID_TASK,SYSDATE,NULL,'RUNNING');
END START_TASK;

PROCEDURE UPDATE_TASK(I_ID_MASTER IN TASKS.ID_MASTER%TYPE,I_ID_TASK IN TASKS.ID_TASK%TYPE,I_STATUS IN TASKS.STATUS%TYPE) AS
BEGIN
    UPDATE TASKS
    SET STATUS = I_STATUS,END_DATE = SYSDATE
    WHERE ID_MASTER = I_ID_MASTER 
    AND ID_TASK = I_ID_TASK
    AND END_DATE IS NULL;
END UPDATE_TASK;

PROCEDURE GET_STATUS(I_ID_MASTER IN TASKS.ID_MASTER%TYPE,I_ID_TASK IN TASKS.ID_TASK%TYPE,O_STATUS OUT TASKS.STATUS%TYPE, O_START_DATE OUT TASKS.START_DATE%TYPE) IS
BEGIN
    SELECT STATUS,START_DATE
    INTO O_STATUS,O_START_DATE
    FROM TASKS
    WHERE ID_MASTER = I_ID_MASTER
    AND ID_TASK = I_ID_TASK;

EXCEPTION WHEN OTHERS THEN
    O_STATUS := 'NONE';
    O_START_DATE := NULL;
END GET_STATUS;

END PKG_TASKMAN;
/