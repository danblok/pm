CREATE TABLE IF NOT EXISTS accounts (
    "id" uuid DEFAULT gen_random_uuid(),
    "email" TEXT UNIQUE NOT NULL,
    "name" TEXT NOT NULL,
    "avatar" TEXT DEFAULT '',
    "deleted" BOOLEAN DEFAULT FALSE,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT now(),
    PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS projects (
    "id" uuid DEFAULT gen_random_uuid(),
    "name" TEXT NOT NULL,
    "description" TEXT NOT NULL DEFAULT '',
    "owner_id" uuid NOT NULL,
    "deleted" BOOLEAN DEFAULT FALSE,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT now(),
    PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS tasks (
    "id" uuid DEFAULT gen_random_uuid(),
    "name" TEXT NOT NULL,
    "start" TIMESTAMP(3) NOT NULL DEFAULT now(),
    "end" TIMESTAMP(3) NOT NULL DEFAULT now() + interval '1 month',
    "status_id" uuid NOT NULL,
    "project_id" uuid NOT NULL,
    "deleted" BOOLEAN DEFAULT FALSE,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT now(),
    PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS statuses (
    "id" uuid DEFAULT gen_random_uuid(),
    "name" TEXT NOT NULL,
    "project_id" uuid NOT NULL,
    "deleted" BOOLEAN DEFAULT FALSE,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT now(),
    PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS projects_to_accounts (
    "project_id" uuid NOT NULL,
    "account_id" uuid NOT NULL
);

CREATE UNIQUE INDEX projects_to_accounts_project_account_unique
ON projects_to_accounts(project_id, account_id);

CREATE UNIQUE INDEX status_name_project_id_unique
ON statuses(name, project_id);

ALTER TABLE projects
ADD CONSTRAINT fk_projects_accounts
FOREIGN KEY (owner_id) REFERENCES accounts(id)
ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE tasks
ADD CONSTRAINT fk_tasks_projects
FOREIGN KEY (project_id) REFERENCES projects(id)
ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE tasks
ADD CONSTRAINT fk_tasks_statuses
FOREIGN KEY (status_id) REFERENCES statuses(id)
ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE statuses
ADD CONSTRAINT fk_statuses_projects
FOREIGN KEY (project_id) REFERENCES projects(id)
ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE projects_to_accounts
ADD CONSTRAINT fk_projects_to_accounts_projects
FOREIGN KEY (project_id) REFERENCES projects(id)
ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE projects_to_accounts
ADD CONSTRAINT fk_projects_to_accounts_accounts
FOREIGN KEY (account_id) REFERENCES accounts(id)
ON DELETE CASCADE ON UPDATE CASCADE;
