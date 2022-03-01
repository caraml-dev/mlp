INSERT INTO applications
(
    name,
    href,
    description,
    icon,
    use_projects,
    is_in_beta,
    is_disabled
)
VALUES
    ('Pipeline', '/pipeline', 'Platform for orchestrating and managing ML pipeline', 'pipelineApp', TRUE, TRUE, FALSE);

-- Previously, Excalibur uses pipelineApp icon. It should be Pipeline app above that uses pipelineApp icon.
UPDATE applications SET icon = 'fleetApp' WHERE name = 'Excalibur';
