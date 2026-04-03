-- 001_init_schema.down.sql
-- 政策雷达数据库回滚脚本

DROP TABLE IF EXISTS ${prefix}_crawler_log;
DROP TABLE IF EXISTS ${prefix}_policy_recommend;
DROP TABLE IF EXISTS ${prefix}_local_standard;
DROP TABLE IF EXISTS ${prefix}_industry_standard;
DROP TABLE IF EXISTS ${prefix}_national_standard;
DROP TABLE IF EXISTS ${prefix}_policy_category;