package main

import (
	"context"
	"log"
	"net/http"

	"github.com/harusame0616/ijuku/apps/api/internal/db"
	libauth "github.com/harusame0616/ijuku/apps/api/lib/auth"
	"github.com/harusame0616/ijuku/apps/api/lib/env"
	"github.com/harusame0616/ijuku/apps/api/lib/txrunner"
	"github.com/harusame0616/ijuku/apps/api/routes/contacts"
	"github.com/harusame0616/ijuku/apps/api/routes/courses/queries"
	userscommands "github.com/harusame0616/ijuku/apps/api/routes/users/commands"
	enrollmentscommands "github.com/harusame0616/ijuku/apps/api/routes/users/enrollments/commands"
	enrollmentsqueries "github.com/harusame0616/ijuku/apps/api/routes/users/enrollments/queries"
	usersqueries "github.com/harusame0616/ijuku/apps/api/routes/users/queries"
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
	courseDetailHandler := queries.NewGetCourseDetailHandler(q)
	courseRepo := enrollmentscommands.NewSqrcCourseRepository(q)
	enrollmentRepo := enrollmentscommands.NewSqrcEnrollmentRepository(q)
	updateEnrollmentHandler := enrollmentscommands.NewUpdateEnrollmentHandler(enrollmentscommands.NewUpdateEnrollmentUsecase(courseRepo, enrollmentRepo))
	enrollHandler := enrollmentscommands.NewEnrollHandler(enrollmentscommands.NewEnrollUsecase(courseRepo, enrollmentRepo))
	topicDetailHandler := queries.NewTopicDetailHandler(q)
	verifier := libauth.NewVerifier(env.Require("SUPABASE_JWT_SECRET"), env.Require("SUPABASE_URL"))
	apikeysHandler := apikeys.NewGenerateApiKeyHandler(apikeys.NewGenerateApiKeyUsecase(apikeys.NewApiKeySqrcRepository(), txrunner.NewPgxTransactionRunner(pool)))
	listApiKeysHandler := apikeys.NewListApiKeysHandler(q)

	getUserHandler := usersqueries.NewGetUserHandler(q)
	updateUserHandler := userscommands.NewUpdateUserHandler(
		userscommands.NewUpdateUserUsecase(userscommands.NewUserSqrcRepository(q)),
	)
	getEnrollmentsHandler := enrollmentsqueries.NewGetEnrollmentsHandler(q)
	getEnrollmentHandler := enrollmentsqueries.NewGetEnrollmentHandler(q)

	postContactHandler := contacts.NewPostContactHandler(q)

	authMiddleware := libauth.Middleware(verifier, q)
	optionalAuthMiddleware := libauth.OptionalMiddleware(verifier, q)

	http.HandleFunc("POST /v1/contacts", postContactHandler.PostContactHandler)

	http.HandleFunc("GET /v1/courses", coursesHandler.GetCoursesHandler)
	http.HandleFunc("GET /v1/courses/{courseId}/sections/{sectionId}/topics/{topicId}", topicDetailHandler.GetTopicDetailHandler)
	http.Handle("GET /v1/courses/{authorSlug}/{courseSlug}", optionalAuthMiddleware(http.HandlerFunc(courseDetailHandler.GetCourseDetailHandler)))

	http.Handle("GET /v1/me", authMiddleware(http.HandlerFunc(getUserHandler.GetUserHandler)))
	http.Handle("PATCH /v1/me", authMiddleware(http.HandlerFunc(updateUserHandler.PatchUserHandler)))
	http.Handle("POST /v1/me/apikeys", authMiddleware(http.HandlerFunc(apikeysHandler.GenerateApiKeyHandler)))
	http.Handle("GET /v1/me/settings/apikeys", authMiddleware(http.HandlerFunc(listApiKeysHandler.ListApiKeysHandler)))
	http.Handle("GET /v1/me/enrollments", authMiddleware(http.HandlerFunc(getEnrollmentsHandler.GetEnrollmentsHandler)))
	http.Handle("GET /v1/me/enrollments/{courseId}", authMiddleware(http.HandlerFunc(getEnrollmentHandler.GetEnrollmentHandler)))
	http.Handle("POST /v1/me/enrollments", authMiddleware(http.HandlerFunc(enrollHandler.PostEnrollmentHandler)))
	http.Handle("PATCH /v1/me/enrollments/{courseId}", authMiddleware(http.HandlerFunc(updateEnrollmentHandler.PatchEnrollmentHandler)))

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
