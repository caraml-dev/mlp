ALTER TABLE applications ADD config jsonb;

UPDATE applications
SET config = '{"sections": [{"name": "Routers", "href": "/routers"}, {"name": "Ensembling Jobs", "href": "/jobs"}]}'
WHERE name = 'Turing';