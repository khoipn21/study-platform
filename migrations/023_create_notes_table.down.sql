-- Drop notes table and related functions/triggers
DROP TRIGGER IF EXISTS trigger_notes_updated_at ON notes;
DROP FUNCTION IF EXISTS update_notes_updated_at();
DROP TABLE IF EXISTS notes;