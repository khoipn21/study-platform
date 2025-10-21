#!/bin/bash

# Script to add Swagger annotations to handler files
# This is a helper script - you'll still need to review and customize the annotations

echo "🔧 Adding Swagger annotations to handlers..."

# Add Course Listendpoints annotations
cat > /tmp/course_annotations.txt << 'EOF'
// ListCourses godoc
// @Summary      List all courses
// @Description  Get a paginated list of all available courses
// @Tags         Courses
// @Accept       json
// @Produce      json
// @Param        page query int false "Page number" default(1)
// @Param        page_size query int false "Items per page" default(10)
// @Param        q query string false "Search query"
// @Success      200 {object} APIResponse "List of courses"
// @Router       /courses [get]
EOF

# Add more endpoint annotations here...

echo "✅ Annotation templates created in /tmp/"
echo "📝 Review and manually add to handler files"
