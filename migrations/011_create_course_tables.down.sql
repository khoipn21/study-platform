-- Drop tables in reverse order
DROP TABLE IF EXISTS enrollments;
DROP TABLE IF EXISTS lectures;
DROP TABLE IF EXISTS courses;

-- Drop enums
DROP TYPE IF EXISTS enrollment_status;
DROP TYPE IF EXISTS lecture_status;
DROP TYPE IF EXISTS course_level;
DROP TYPE IF EXISTS course_status;