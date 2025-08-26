-- Insert a workflow
INSERT INTO workflows (id, name) VALUES ('wf1', 'Shopping App Builder');

-- Insert nodes
INSERT INTO workflow_nodes (id, workflow_id, type, config) VALUES
('n1', 'wf1', 'llm', '{"name":"ProductAgent","max_tokens":512,"prompt":"Generate product requirements"}'),
('n2', 'wf1', 'llm', '{"name":"CodeAgent","max_tokens":1024,"prompt":"Generate Go backend code"}'),
('n3', 'wf1', 'llm', '{"name":"DesignAgent","max_tokens":512,"prompt":"Suggest UI/UX design"}');

-- Insert edges
INSERT INTO workflow_edges (id, workflow_id, source_node_id, target_node_id, condition) VALUES
('e1', 'wf1', 'n1', 'n2', NULL),
('e2', 'wf1', 'n1', 'n3', NULL);
