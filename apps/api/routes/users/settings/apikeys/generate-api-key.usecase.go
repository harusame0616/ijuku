package apikeys

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type generateApiKeyRepository interface {
	count(ctx context.Context, userId uuid.UUID) int
	save(ctx context.Context, apiKey hashedApiKey) error
}

type generateApiKeyUsecase struct {
	apiKeyRepository generateApiKeyRepository
}

func NewGenerateApiKeyUsecase(repository generateApiKeyRepository) generateApiKeyUsecase {
	return generateApiKeyUsecase{apiKeyRepository: repository}
}

type generateApiKeyExecuteResult struct {
	Apikey string `json:"apikey"`
}

func (usecase *generateApiKeyUsecase) Execute(ctx context.Context, userId uuid.UUID, expiredAt *time.Time) (generateApiKeyExecuteResult, error) {
	// 不変条件が上限チェックのみのため、Usecase 層でチェックする
	// 複数の不変条件が加わる場合は UserApiKeys ドメインコレクションへの昇格を検討する
	userApiKeyCount := usecase.apiKeyRepository.count(ctx, userId)
	if userApiKeyCount >= apiKeyMaxCount {
		return generateApiKeyExecuteResult{}, ErrApiKeyCountExceedsLimit
	}

	hashedApiKey, plainApiKey := NewHashedApiKey(NewHashedApiKeyParams{
		UserID:    userId,
		ExpiredAt: expiredAt,
	})

	if err := usecase.apiKeyRepository.save(ctx, hashedApiKey); err != nil {
		return generateApiKeyExecuteResult{}, err
	}

	return generateApiKeyExecuteResult{
		Apikey: plainApiKey,
	}, nil
}
