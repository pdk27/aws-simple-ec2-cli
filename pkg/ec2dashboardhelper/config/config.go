package config

import "github.com/aws/aws-sdk-go/aws/session"

type Config struct {
	AWSSession  *session.Session
	Region string
	Granularity string
	CostType string // list?
	EvaluationPeriodInDays int
	ShowCostsByCategories bool
}