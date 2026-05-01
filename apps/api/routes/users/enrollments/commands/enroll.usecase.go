package commands

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type EnrollParams struct {
	userId   uuid.UUID
	courseId uuid.UUID
}

type EnrollResult struct {
	CourseId   string
	EnrolledAt time.Time
}

type EnrollmentCreator interface {
	findByUserAndCourse(ctx context.Context, userId, courseId uuid.UUID) (Enrollment, error)
	create(ctx context.Context, enrollment Enrollment) error
}

type EnrollUsecaseInterface interface {
	execute(ctx context.Context, params EnrollParams) (EnrollResult, error)
}

type EnrollUsecase struct {
	courseRepository     CourseRepository
	enrollmentRepository EnrollmentCreator
	now                  func() time.Time
}

func NewEnrollUsecase(courseRepo CourseRepository, enrollmentRepo EnrollmentCreator) EnrollUsecase {
	return EnrollUsecase{
		courseRepository:     courseRepo,
		enrollmentRepository: enrollmentRepo,
		now:                  time.Now,
	}
}

func (u EnrollUsecase) execute(ctx context.Context, params EnrollParams) (EnrollResult, error) {
	course, err := u.courseRepository.getCourseByCourseId(ctx, params.courseId)
	if err != nil {
		return EnrollResult{}, err
	}

	if err := course.checkEnrollable(params.userId); err != nil {
		return EnrollResult{}, err
	}

	if _, err := u.enrollmentRepository.findByUserAndCourse(ctx, params.userId, params.courseId); err == nil {
		return EnrollResult{}, ErrAlreadyEnrolled
	} else if err != ErrNotEnrolled {
		return EnrollResult{}, err
	}

	enrollment := NewEnrollment(params.userId, params.courseId, u.now())

	if err := u.enrollmentRepository.create(ctx, enrollment); err != nil {
		return EnrollResult{}, err
	}

	return EnrollResult{
		CourseId:   enrollment.CourseId().String(),
		EnrolledAt: enrollment.EnrolledAt(),
	}, nil
}
