-- Re-enable Clockwork UI
UPDATE applications SET is_disabled = FALSE WHERE name = 'Clockwork';

-- Remove Turing
DELETE FROM applications WHERE name = 'Turing';