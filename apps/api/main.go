package main

import (
	"context"
	"log"
	"net/http"

	"github.com/harusame0616/ijuku/apps/api/internal/db"
	"github.com/harusame0616/ijuku/apps/api/lib/env"
	"github.com/harusame0616/ijuku/apps/api/lib/txrunner"
	"github.com/harusame0616/ijuku/apps/api/routes/courses/commands"
	"github.com/harusame0616/ijuku/apps/api/routes/courses/queries"
	"github.com/harusame0616/ijuku/apps/api/routes/users/settings/apikeys"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, env.Require("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	q := db.New(pool)

	coursesHandler := queries.NewCoursesHandlers(q)
	enrollHandler := commands.NewHandler(commands.NewEnrollCourseUsecase(commands.NewSqrcCourseRepository(q), commands.NewSqrcUserTopicProgressRepository(q)))
	topicDetailHandler := queries.NewTopicDetailHandler(q)
	apikeysHandler := apikeys.NewGenerateApiKeyHandler(apikeys.NewGenerateApiKeyUsecase(apikeys.NewApiKeySqrcRepository(), txrunner.NewPgxTransactionRunner(pool)))

	http.HandleFunc("GET /v1/courses", coursesHandler.GetCoursesHandler)
	http.HandleFunc("POST /v1/courses/{courseId}/enrollment", enrollHandler.PostEnrollmentHandler)
	http.HandleFunc("GET /v1/courses/{courseId}/sections/{sectionId}/topics/{topicId}", topicDetailHandler.GetTopicDetailHandler)
	http.HandleFunc("POST /v1/users/{userID}/apikeys", apikeysHandler.GenerateApiKeyHandler)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
