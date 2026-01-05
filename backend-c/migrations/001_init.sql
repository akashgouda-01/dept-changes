-- Core tables for EduVault

CREATE TABLE IF NOT EXISTS students (
    id              SERIAL PRIMARY KEY,
    register_number TEXT UNIQUE NOT NULL,
    name            TEXT NOT NULL,
    email           TEXT,
    section         TEXT NOT NULL,
    semester        INT,
    is_present      BOOLEAN NOT NULL DEFAULT TRUE,
    faculty_email   TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS certificates (
    id               SERIAL PRIMARY KEY,
    student_reg_no   TEXT NOT NULL REFERENCES students(register_number) ON DELETE CASCADE,
    section          TEXT NOT NULL,
    faculty_email    TEXT NOT NULL,
    drive_link       TEXT NOT NULL,
    ml_status        TEXT NOT NULL DEFAULT 'pending',       -- pending | verified | duplicate (for future ML)
    faculty_status   TEXT NOT NULL DEFAULT 'pending',       -- pending | legit | not_legit
    ml_score         NUMERIC,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_students_faculty_email ON students(faculty_email);
CREATE INDEX IF NOT EXISTS idx_students_section ON students(section);
CREATE INDEX IF NOT EXISTS idx_certificates_faculty_email ON certificates(faculty_email);
CREATE INDEX IF NOT EXISTS idx_certificates_section ON certificates(section);
CREATE INDEX IF NOT EXISTS idx_certificates_status ON certificates(faculty_status, ml_status);



