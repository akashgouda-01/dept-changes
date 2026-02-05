ALTER TABLE faculty_certificates ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW();

UPDATE faculty_certificates SET updated_at = uploaded_at WHERE updated_at IS NULL;
