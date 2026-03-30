package main

import (
	"context"
	"log"
	"net/http"

	"github.com/harusame0616/ijuku/apps/api/internal/db"
	"github.com/harusame0616/ijuku/apps/api/routes/courses/commands"
	"github.com/harusame0616/ijuku/apps/api/routes/courses/queries"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, "postgresql://postgres:password@localhost:5432/postgres")
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	q := db.New(pool)

	coursesHandler := queries.NewCoursesHandlers(queries.NewSqrcCourseQueryService(q))
	enrollHandler := commands.NewHandler(commands.NewEnrollCourseUsecase(commands.NewSqrcCourseRepository(q), commands.NewSqrcUserTopicProgressRepository(q)))

	http.HandleFunc("GET /v1/courses", coursesHandler.GetCoursesHandler)
	http.HandleFunc("POST /v1/courses/{courseId}/enrollment", enrollHandler.PostEnrollmentHandler)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
