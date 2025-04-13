CREATE TABLE projects (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE suite_runs (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL REFERENCES projects(id),
    git_branch TEXT,
    git_commit TEXT,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP,
    config JSONB,
    summary JSONB,
    created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE spec_runs (
    id TEXT PRIMARY KEY,
    suite_run_id TEXT NOT NULL REFERENCES suite_runs(id),
    name TEXT NOT NULL,
    file_name TEXT NOT NULL,
    status TEXT NOT NULL,
    duration_ms INT,
    line_number INT,
    created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE attachments (
    id TEXT PRIMARY KEY,
    spec_run_id TEXT REFERENCES spec_runs(id),
    name TEXT,
    type TEXT,
    content TEXT
);

CREATE TABLE steps (
    id TEXT PRIMARY KEY,
    spec_run_id TEXT REFERENCES spec_runs(id),
    description TEXT,
    status TEXT,
    duration_ms INT,
    created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE custom_reports (
    id TEXT PRIMARY KEY,
    suite_run_id TEXT REFERENCES suite_runs(id),
    name TEXT,
    content JSONB,
    created_at TIMESTAMP DEFAULT now()
);
