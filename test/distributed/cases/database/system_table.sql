SELECT table_catalog,table_schema,table_name,table_type from `information_schema`.`tables` where table_name = 'mo_tables';
SELECT * FROM `information_schema`.`character_sets` LIMIT 0,1000;
SELECT * FROM `information_schema`.`columns` where TABLE_NAME = 'mo_tables' order by ORDINAL_POSITION LIMIT 2;
SELECT * FROM `information_schema`.`key_column_usage` LIMIT 0,1000;
SELECT * FROM `information_schema`.`profiling` LIMIT 0,1000;
SELECT * FROM `information_schema`.`schemata` where schema_name = 'information_schema';
SELECT * FROM `information_schema`.`triggers` LIMIT 0,1000;
SELECT * FROM `information_schema`.`user_privileges` LIMIT 0,1000;
SELECT * FROM `mysql`.`columns_priv` LIMIT 0,1000;
SELECT * FROM `mysql`.`db` LIMIT 0,1000;
SELECT * FROM `mysql`.`procs_priv` LIMIT 0,1000;
SELECT * FROM `mysql`.`tables_priv` LIMIT 0,1000;
SELECT * FROM `mysql`.`user` LIMIT 0,1000;
use mysql;
show tables;
show columns from `user`;
show columns from `db`;
show columns from `procs_priv`;
show columns from `columns_priv`;
show columns from `tables_priv`;
use information_schema;
show tables;
show columns from `KEY_COLUMN_USAGE`;
show columns from `COLUMNS`;
show columns from `PROFILING`;
show columns from `USER_PRIVILEGES`;
-- @bvt:issue#12196
show columns from `SCHEMATA`;
-- @bvt:issue
show columns from `CHARACTER_SETS`;
show columns from `TRIGGERS`;
show columns from `TABLES`;
show columns from `PARTITIONS`;
drop database if exists test;
create database test;
use test;
drop table if exists t2;
create table t2(b int, a int);
desc t2;
drop table t2;
drop database test;