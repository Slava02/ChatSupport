package outbox

import (
	"context"
	"time"

	"github.com/Slava02/ChatSupport/internal/types"
)

func (s *Service) Put(ctx context.Context, name, payload string, availableAt time.Time) (types.JobID, error) {
	return s.jobsRepo.CreateJob(ctx, name, payload, availableAt)
}
