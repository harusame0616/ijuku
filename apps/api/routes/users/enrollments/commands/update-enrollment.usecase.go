package commands

import (
	"context"

	"github.com/google/uuid"
)

type UpdateEnrollmentParams struct {
	userId   uuid.UUID
	courseId uuid.UUID
	topicId  uuid.UUID
	status   ProgressStatus
}

type UpdateEnrollmentResult struct {
	TopicId string
	Status  string
}

type CourseRepository interface {
	getCourseByCourseId(ctx context.Context, courseId uuid.UUID) (Course, error)
}

type EnrollmentRepository interface {
	findByUserAndCourse(ctx context.Context, userId, courseId uuid.UUID) (Enrollment, error)
	upsertTopicProgress(ctx context.Context, userId uuid.UUID, progress TopicProgress) error
}

type UpdateEnrollmentUsecaseInterface interface {
	execute(ctx context.Context, params UpdateEnrollmentParams) (UpdateEnrollmentResult, error)
}

type UpdateEnrollmentUsecase struct {
	courseRepository     CourseRepository
	enrollmentRepository EnrollmentRepository
}

func NewUpdateEnrollmentUsecase(courseRepo CourseRepository, enrollmentRepo EnrollmentRepository) UpdateEnrollmentUsecase {
	return UpdateEnrollmentUsecase{
		courseRepository:     courseRepo,
		enrollmentRepository: enrollmentRepo,
	}
}

func (usecase UpdateEnrollmentUsecase) execute(ctx context.Context, params UpdateEnrollmentParams) (UpdateEnrollmentResult, error) {
	course, err := usecase.courseRepository.getCourseByCourseId(ctx, params.courseId)
	if err != nil {
		return UpdateEnrollmentResult{}, err
	}

	if _, err := course.findTopicById(params.topicId); err != nil {
		return UpdateEnrollmentResult{}, err
	}

	enrollment, err := usecase.enrollmentRepository.findByUserAndCourse(ctx, params.userId, params.courseId)
	if err != nil {
		return UpdateEnrollmentResult{}, err
	}

	progress, err := enrollment.UpdateTopicProgress(params.topicId, params.status)
	if err != nil {
		return UpdateEnrollmentResult{}, err
	}

	if err := usecase.enrollmentRepository.upsertTopicProgress(ctx, enrollment.UserId(), progress); err != nil {
		return UpdateEnrollmentResult{}, err
	}

	return UpdateEnrollmentResult{
		TopicId: progress.TopicId().String(),
		Status:  string(progress.Status()),
	}, nil
}
