-- Add foreign key constraints with proper cascade behavior as recommended in architecture review

-- Add missing foreign key constraint for videos.upload_user_id
ALTER TABLE videos
ADD CONSTRAINT fk_videos_uploader
FOREIGN KEY (upload_user_id) REFERENCES users(id) ON DELETE RESTRICT;

-- Add foreign key constraints with proper cascade behavior for courses
ALTER TABLE courses
ADD CONSTRAINT fk_courses_creator
FOREIGN KEY (creator_id) REFERENCES users(id) ON DELETE SET NULL;

-- Add foreign key constraints with proper cascade behavior for lectures
ALTER TABLE lectures
ADD CONSTRAINT fk_lectures_course
FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE;

-- Add foreign key constraints with proper cascade behavior for videos
ALTER TABLE videos
ADD CONSTRAINT fk_videos_lecture
FOREIGN KEY (lecture_id) REFERENCES lectures(id) ON DELETE CASCADE;

-- Add foreign key constraints for enrollment with transaction reference
-- First add transaction_id column if it doesn't exist
ALTER TABLE enrollment
ADD COLUMN IF NOT EXISTS transaction_id UUID;

ALTER TABLE enrollment
ADD CONSTRAINT fk_enrollment_transaction
FOREIGN KEY (transaction_id) REFERENCES transactions(id) ON DELETE RESTRICT;