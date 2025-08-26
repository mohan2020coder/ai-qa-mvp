CREATE TABLE agents (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    type VARCHAR(20) NOT NULL, -- "LLM", "API", "Script", etc.
    config JSONB NOT NULL      -- config for this agent
);


--Example rows:

--name	type	config
--Product	LLM	{"system":"You are a product manager..."}
--Deploy	Script	{"command":"./deploy.sh"}
--Notify	API	{"url":"https://hooks.slack.com/...", "method":"POST"}
--
--



-- Workflows table
CREATE TABLE IF NOT EXISTS workflows (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Nodes inside workflows
CREATE TABLE IF NOT EXISTS workflow_nodes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workflow_id UUID NOT NULL,
    type TEXT NOT NULL,             -- e.g. "llm", "condition", "parallel"
    config JSONB NOT NULL,          -- agent config (model, max tokens, prompt, etc.)
    position_x INTEGER,             -- optional, for editor UI
    position_y INTEGER,             -- optional, for editor UI
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (workflow_id) REFERENCES workflows(id) ON DELETE CASCADE
);

-- Edges connecting nodes
CREATE TABLE IF NOT EXISTS workflow_edges (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workflow_id UUID NOT NULL,
    source_node_id UUID NOT NULL,
    target_node_id UUID NOT NULL,
    condition TEXT,                 -- optional (e.g. "success", "failure", "branch A")
    FOREIGN KEY (workflow_id) REFERENCES workflows(id) ON DELETE CASCADE,
    FOREIGN KEY (source_node_id) REFERENCES workflow_nodes(id) ON DELETE CASCADE,
    FOREIGN KEY (target_node_id) REFERENCES workflow_nodes(id) ON DELETE CASCADE
);


ALTER TABLE workflow_nodes ADD COLUMN name TEXT;
ALTER TABLE workflow_nodes ADD COLUMN next_ids JSONB;
ALTER TABLE workflow_nodes ADD COLUMN condition TEXT;

-- Insert a workflow
-- Insert a workflow (ID auto-generated)
INSERT INTO workflows (name) 
VALUES ('Shopping App Builder')
RETURNING id;


-- Insert nodes
INSERT INTO workflow_nodes (workflow_id, type, config) VALUES
('c941286a-f580-4f08-9142-b1b9f64c75a2', 'llm',
 '{"name":"ProductAgent","max_tokens":512,"prompt":"Generate product requirements"}'::jsonb),
('c941286a-f580-4f08-9142-b1b9f64c75a2', 'llm',
 '{"name":"CodeAgent","max_tokens":1024,"prompt":"Generate Go backend code"}'::jsonb),
('c941286a-f580-4f08-9142-b1b9f64c75a2', 'llm',
 '{"name":"DesignAgent","max_tokens":512,"prompt":"Suggest UI/UX design"}'::jsonb);


 select * from workflow_nodes

INSERT INTO workflow_edges (workflow_id, source_node_id, target_node_id, condition)
VALUES
('c941286a-f580-4f08-9142-b1b9f64c75a2', 'f433d8f4-e90c-445d-bbb5-0cc545883c8f', '22695631-c485-46c0-831b-0f7c1d83d190', NULL),
('c941286a-f580-4f08-9142-b1b9f64c75a2', 'f433d8f4-e90c-445d-bbb5-0cc545883c8f', '841219b1-d443-444b-995e-217d2a64a5da', NULL);

