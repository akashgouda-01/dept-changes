-- Create enums
DO $$ BEGIN
    CREATE TYPE ml_status_enum AS ENUM ('PENDING', 'VERIFIED', 'DUPLICATE');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

DO $$ BEGIN
    CREATE TYPE faculty_status_enum AS ENUM ('PENDING', 'LEGIT', 'NOT_LEGIT');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

-- Create faculty_certificates table
CREATE TABLE IF NOT EXISTS faculty_certificates (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    drive_link TEXT NOT NULL,
    reg_no TEXT NOT NULL,
    section TEXT NOT NULL,
    student_name TEXT NOT NULL,
    faculty_id TEXT NOT NULL,
    uploaded_at TIMESTAMP WITH TIME ZONE NOT NULL,
    ml_status ml_status_enum DEFAULT 'PENDING' NOT NULL,
    faculty_status faculty_status_enum DEFAULT 'PENDING' NOT NULL,
    is_legit BOOLEAN,
    ml_score DOUBLE PRECISION,
    archived BOOLEAN DEFAULT FALSE NOT NULL
);

-- Create student_statistics table
CREATE TABLE IF NOT EXISTS student_statistics (
    reg_no TEXT PRIMARY KEY,
    student_name TEXT,
    section TEXT,
    total_uploaded INTEGER DEFAULT 0 NOT NULL,
    legit_count INTEGER DEFAULT 0 NOT NULL,
    not_legit_count INTEGER DEFAULT 0 NOT NULL,
    pending_certificates INTEGER DEFAULT 0 NOT NULL,
    ml_verified_certificates INTEGER DEFAULT 0 NOT NULL,
    last_updated TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);

-- Create section_statistics table
CREATE TABLE IF NOT EXISTS section_statistics (
    section TEXT PRIMARY KEY,
    total_uploaded INTEGER DEFAULT 0 NOT NULL,
    legit_count INTEGER DEFAULT 0 NOT NULL,
    not_legit_count INTEGER DEFAULT 0 NOT NULL,
    last_updated TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);
