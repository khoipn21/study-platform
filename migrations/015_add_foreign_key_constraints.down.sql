-- Remove foreign key constraints

ALTER TABLE videos DROP CONSTRAINT IF EXISTS fk_videos_uploader;
ALTER TABLE courses DROP CONSTRAINT IF EXISTS fk_courses_creator;
ALTER TABLE lectures DROP CONSTRAINT IF EXISTS fk_lectures_course;
ALTER TABLE videos DROP CONSTRAINT IF EXISTS fk_videos_lecture;
ALTER TABLE enrollment DROP CONSTRAINT IF EXISTS fk_enrollment_transaction;

-- Remove transaction_id column if it was added
ALTER TABLE enrollment DROP COLUMN IF EXISTS transaction_id;