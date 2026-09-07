package service

import (
	"context"

	"github.com/gosharplite/aixbdd-go/internal/domain"
)

type PipelineService struct {
	oppRepo domain.OpportunityRepository
}

func NewPipelineService(oppRepo domain.OpportunityRepository) *PipelineService {
	return &PipelineService{oppRepo: oppRepo}
}

func (s *PipelineService) GetPipelineSummary(ctx context.Context, actor *domain.CurrentActor) (*domain.PipelineSummary, error) {
	if actor == nil || !actor.IsAuthenticated {
		return nil, domain.ErrUnauthorized
	}

	var targetOwner *int64
	if actor.User.Role != domain.RoleSalesManager {
		targetOwner = &actor.User.ID
	}

	allOpps, err := s.oppRepo.List(ctx, targetOwner, nil)
	if err != nil {
		return nil, err
	}

	stagesOrder := []domain.SalesStage{
		domain.StageProspecting,
		domain.StageContacted,
		domain.StageProposal,
		domain.StageClosedWon,
		domain.StageClosedLost,
	}

	stageBuckets := make(map[domain.SalesStage]*domain.StageSummary)
	for _, st := range stagesOrder {
		stageBuckets[st] = &domain.StageSummary{
			Stage:         st,
			Count:         0,
			TotalAmount:   0.0,
			Opportunities: []domain.Opportunity{},
		}
	}

	var activeTotal float64
	for _, opp := range allOpps {
		bucket, ok := stageBuckets[opp.Stage]
		if !ok {
			bucket = &domain.StageSummary{
				Stage:         opp.Stage,
				Count:         0,
				TotalAmount:   0.0,
				Opportunities: []domain.Opportunity{},
			}
			stageBuckets[opp.Stage] = bucket
		}

		bucket.Count++
		bucket.TotalAmount += opp.Amount
		bucket.Opportunities = append(bucket.Opportunities, opp)

		// Active stages are Prospecting, Contacted, Proposal
		if opp.Stage == domain.StageProspecting || opp.Stage == domain.StageContacted || opp.Stage == domain.StageProposal {
			activeTotal += opp.Amount
		}
	}

	stagesList := make([]domain.StageSummary, 0, len(stagesOrder))
	for _, st := range stagesOrder {
		if bucket, exists := stageBuckets[st]; exists {
			stagesList = append(stagesList, *bucket)
		}
	}

	summary := &domain.PipelineSummary{
		ViewerRole:            actor.User.Role,
		Stages:                stagesList,
		ActiveTotalAmount:     activeTotal,
		TotalOpportunityCount: len(allOpps),
	}

	return summary, nil
}
