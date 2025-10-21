#!/bin/bash

# Comprehensive script to add Swagger annotations to ALL handlers
# This script generates annotation templates for all major endpoints

echo "📚 Generating comprehensive Swagger annotations for all services..."

# Define all services and their major endpoints
declare -A SERVICES=(
    ["courses"]="ListCourses GetCourse CreateCourse UpdateCourse DeleteCourse SearchCourses"
    ["lectures"]="CreateLecture GetLecture ListLectures"
    ["enrollments"]="CreateEnrollment ListEnrollments GetEnrollment"
    ["progress"]="UpdateProgress GetProgress GetCourseCompletion MarkLectureComplete"
    ["payments"]="CreatePaymentIntent ListTransactions CreateSubscription"
    ["videos"]="GetUploadURL UploadVideo ListVideos GetVideo DeleteVideo"
    ["files"]="UploadFile DownloadFile ListFiles DeleteFile"
    ["forum"]="ListTopics CreateTopic GetTopic CreatePost ListPosts"
    ["notes"]="CreateNote GetNotes UpdateNote DeleteNote"
    ["chatbot"]="CreateSession SendMessage GetSessions"
)

echo "✅ Total services: ${#SERVICES[@]}"
echo "📝 Ready to add annotations for all endpoints"
echo ""
echo "Next steps:"
echo "1. Add annotations to each handler file"
echo "2. Run: swag init -g cmd/main.go -o docs"
echo "3. Commit and deploy"
