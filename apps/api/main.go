package main

import (
	"context"
	"log"
	"net/http"

	"github.com/harusame0616/ijuku/apps/api/internal/db"
	"github.com/harusame0616/ijuku/apps/api/lib/env"
	libauth "github.com/harusame0616/ijuku/apps/api/lib/auth"
	"github.com/harusame0616/ijuku/apps/api/lib/txrunner"
	"github.com/harusame0616/ijuku/apps/api/routes/courses/queries"
	enrollmentscommands "github.com/harusame0616/ijuku/apps/api/routes/users/enrollments/commands"
	enrollmentsqueries "github.com/harusame0616/ijuku/apps/api/routes/users/enrollments/queries"
	"github.com/harusame0616/ijuku/apps/api/routes/users/settings/apikeys"
	userscommands "github.com/harusame0616/ijuku/apps/api/routes/users/commands"
	usersqueries "github.com/harusame0616/ijuku/apps/api/routes/users/queries"
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
	enrollHandler := enrollmentscommands.NewHandler(enrollmentscommands.NewEnrollCourseUsecase(enrollmentscommands.NewSqrcCourseRepository(q), enrollmentscommands.NewSqrcUserTopicProgressRepository(q)))
	topicDetailHandler := queries.NewTopicDetailHandler(q)
	apikeysHandler := apikeys.NewGenerateApiKeyHandler(apikeys.NewGenerateApiKeyUsecase(apikeys.NewApiKeySqrcRepository(), txrunner.NewPgxTransactionRunner(pool)))
	verifier := libauth.NewVerifier(env.Require("SUPABASE_JWT_SECRET"), env.Require("SUPABASE_URL"))
	listApiKeysHandler := apikeys.NewListApiKeysHandler(q, verifier)

	getUserHandler := usersqueries.NewGetUserHandler(q)
	updateUserHandler := userscommands.NewUpdateUserHandler(
		userscommands.NewUpdateUserUsecase(userscommands.NewUserSqrcRepository(q)),
		verifier,
	)
	getEnrollmentsHandler := enrollmentsqueries.NewGetEnrollmentsHandler(q)

	http.HandleFunc("GET /v1/courses", coursesHandler.GetCoursesHandler)
	http.HandleFunc("POST /v1/courses/{courseId}/enrollment", enrollHandler.PostEnrollmentHandler)
	http.HandleFunc("GET /v1/courses/{courseId}/sections/{sectionId}/topics/{topicId}", topicDetailHandler.GetTopicDetailHandler)
	http.HandleFunc("POST /v1/users/{userID}/apikeys", apikeysHandler.GenerateApiKeyHandler)
	http.HandleFunc("GET /v1/users/{userID}/settings/apikeys", listApiKeysHandler.ListApiKeysHandler)
	http.HandleFunc("GET /v1/users/{userID}", getUserHandler.GetUserHandler)
	http.HandleFunc("PATCH /v1/users/{userID}", updateUserHandler.PatchUserHandler)
	http.HandleFunc("GET /v1/users/{userID}/enrollments", getEnrollmentsHandler.GetEnrollmentsHandler)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
