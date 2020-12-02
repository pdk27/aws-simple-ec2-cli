package types

type InstanceInfo struct {
	Region string
	InstanceId string
	InstanceName string
	InstanceType string
	CpuCredits string
	EbsAttached bool
	Recommendations []Recommendation
}

type Recommendation struct {
	Finding string
	RecommendedInstanceType string
}

type IdleInstanceTracker struct {
	AvgCpuUtilPercentage int
	AvgMemoryPercentage int
	AvgNetworkIn int
	AvgNetworkOut int
	AvgEbsReadBytesPerSec int
	AvgEbsWriteBytesPerSec int
	AvgCpuCreditUsagePercentage int // CPUCreditUsage / CPUCreditBalance
	MaxCPUSurplusCreditsCharged int
	Period int
}

type CostTracker struct {
	ByInstanceType []CostTrackerItem
	ByRegion []CostTrackerItem
	ByCapacityType []CostTrackerItem // spot or on-demand
	ByResource []CostTrackerItem
	CostUnit string
	Period int
}

type CostTrackerItem struct {
	Key string
	Value string
}

type Ec2SpotTracker struct {
	spotRunningPercentage int // cost explorer API: spot/ total ec2 running
}

