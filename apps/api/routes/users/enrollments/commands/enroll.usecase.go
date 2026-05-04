package commands

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type EnrollParams struct {
	userId     uuid.UUID
	authorSlug string
	courseSlug string
}

type EnrollResult struct {
	CourseId   uuid.UUID
	EnrolledAt time.Time
}

type EnrollCourseRepository interface {
	getCourseBySlug(ctx context.Context, authorSlug, courseSlug string) (Course, error)
}

type EnrollmentCreator interface {
	findByUserAndCourse(ctx context.Context, userId, courseId uuid.UUID) (Enrollment, error)
	create(ctx context.Context, enrollment Enrollment) error
}

type EnrollUsecaseInterface interface {
	execute(ctx context.Context, params EnrollParams) (EnrollResult, error)
}

type EnrollUsecase struct {
	courseRepository     EnrollCourseRepository
	enrollmentRepository EnrollmentCreator
	now                  func() time.Time
}

func NewEnrollUsecase(courseRepo EnrollCourseRepository, enrollmentRepo EnrollmentCreator) EnrollUsecase {
	return EnrollUsecase{
		courseRepository:     courseRepo,
		enrollmentRepository: enrollmentRepo,
		now:                  time.Now,
	}
}

func (u EnrollUsecase) execute(ctx context.Context, params EnrollParams) (EnrollResult, error) {
	course, err := u.courseRepository.getCourseBySlug(ctx, params.authorSlug, params.courseSlug)
	if err != nil {
		return EnrollResult{}, err
	}

	if err := course.checkEnrollable(params.userId); err != nil {
		return EnrollResult{}, err
	}

	if _, err := u.enrollmentRepository.findByUserAndCourse(ctx, params.userId, course.courseId); err == nil {
		return EnrollResult{}, ErrAlreadyEnrolled
	} else if err != ErrNotEnrolled {
		return EnrollResult{}, err
	}

	enrollment := NewEnrollment(params.userId, course.courseId, u.now())

	if err := u.enrollmentRepository.create(ctx, enrollment); err != nil {
		return EnrollResult{}, err
	}

	return EnrollResult{
		CourseId:   enrollment.CourseId(),
		EnrolledAt: enrollment.EnrolledAt(),
	}, nil
}
