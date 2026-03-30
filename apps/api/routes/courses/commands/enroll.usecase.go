package commands

import (
	"context"

	"golang.org/x/sync/errgroup"
)

type EnrollCourseUsecaseParams struct {
	userId        string
	courseId      string
	sectionNumber *int
	topicNumber   *int
}

type EnrollCourseUsecase struct {
	courseRepository   CourseRepository
	progressRepository ProgressRepository
}

func NewEnrollCourseUsecase(courseRepo CourseRepository, progressRepo ProgressRepository) EnrollCourseUsecase {
	return EnrollCourseUsecase{
		courseRepository:   courseRepo,
		progressRepository: progressRepo,
	}
}

type CourseRepository interface {
	getCourseByCourseId(ctx context.Context, courseId string) (Course, error)
}

type ProgressRepository interface {
	getProgress(ctx context.Context, userId string, courseId string) (Progress, error)
	save(ctx context.Context, progress *Progress) error
}

type EnrollCourseUsecaseInterface interface {
	execute(ctx context.Context, params EnrollCourseUsecaseParams) (string, error)
}

func (usecase EnrollCourseUsecase) execute(ctx context.Context, params EnrollCourseUsecaseParams) (string, error) {
	var course Course
	var progress Progress

	g, gctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		var err error
		course, err = usecase.courseRepository.getCourseByCourseId(gctx, params.courseId)
		return err
	})
	g.Go(func() error {
		var err error
		progress, err = usecase.progressRepository.getProgress(gctx, params.userId, params.courseId)
		return err
	})

	if err := g.Wait(); err != nil {
		return "", err
	}

	if err := course.checkEnrollable(params.userId); err != nil {
		return "", err
	}

	section, topic, err := course.findTopicToStart(params.sectionNumber, params.topicNumber, progress.lastTopicPosition())
	if err != nil {
		return "", err
	}

	topicId := progress.start(section, topic)

	return topicId, usecase.progressRepository.save(ctx, &progress)
}
