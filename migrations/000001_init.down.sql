BEGIN;
ALTER TABLE projects_to_accounts DROP CONSTRAINT fk_projects_to_accounts_accounts;
ALTER TABLE projects_to_accounts DROP CONSTRAINT fk_projects_to_accounts_projects;
ALTER TABLE comments DROP CONSTRAINT fk_comments_accounts_sender;
ALTER TABLE comments DROP CONSTRAINT fk_comments_tasks;
ALTER TABLE statuses DROP CONSTRAINT fk_statuses_projects;
ALTER TABLE tasks DROP CONSTRAINT fk_tasks_statuses;
ALTER TABLE tasks DROP CONSTRAINT fk_tasks_projects;
ALTER TABLE projects DROP CONSTRAINT fk_projects_accounts;

DROP TABLE IF EXISTS projects_to_accounts;
DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS statuses;
DROP TABLE IF EXISTS tasks;
DROP TABLE IF EXISTS projects;
DROP TABLE IF EXISTS accounts;
COMMIT;
