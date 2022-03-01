-- Remove Pipeline
DELETE FROM applications WHERE name = 'Pipeline';

-- Should we give pipelineApp icon to Excalibur back?
UPDATE applications SET icon = 'pipelineApp' WHERE name = 'Excalibur';
