package core

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

// statisticsManager implements the StatisticsManager interface.
type statisticsManager struct{}

// NewStatisticsManager creates a new instance of the StatisticsManager.
func NewStatisticsManager(client *mongo.Client) StatisticsManager {
	return &statisticsManager{}
}

func (s *statisticsManager) GetAverageMonthlyIncome(ctx context.Context, userID string) (float64, error) {
	panic("implement me")
}

func (s *statisticsManager) GetDebitPerCategory(ctx context.Context, params *ParamsStatsGetAmountDistribution) (
	map[string]float64, error,
) {
	panic("implement me")
}

func (s *statisticsManager) GetCreditPerCategory(ctx context.Context, params *ParamsStatsGetAmountDistribution) (
	map[string]float64, error,
) {
	panic("implement me")
}

func (s *statisticsManager) GetDebitPerTag(ctx context.Context, params *ParamsStatsGetAmountDistribution) (
	map[string]float64, error,
) {
	panic("implement me")
}

func (s *statisticsManager) GetCreditPerTag(ctx context.Context, params *ParamsStatsGetAmountDistribution) (
	map[string]float64, error,
) {
	panic("implement me")
}

func (s *statisticsManager) GetBalanceOverTime(ctx context.Context) (map[int64]float64, error) {
	panic("implement me")
}
