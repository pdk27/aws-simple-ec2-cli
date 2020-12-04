package info

import (
	"fmt"
	"simple-ec2/pkg/table"
	"time"
)

// impl info
// tabulate info
// populate info from Cost tracker
// on init of dsahboard, fetch from AWS if refresh - from CE, CW, EC2
// formulate Cost-effectiveness factor
// impl CW API integration

type InstanceInfo struct {
	// info tracked per instance
	InstanceId        string
	//InstanceName      string
	InstanceType      string
	CapacityType		string
	Region 			  string
	//EbsAttached       string
	//CPUCreditsBalance float64
	AvgCostPerPeriod  CostPerPeriod
	TotalCostPerPeriod	CostPerPeriod

	UsageInfo 		InstanceUsageInfo
	Recommendations Recommendation

	// derived values
	CostEffectivenessFactor	int  // percentage, better name?

	// API parameters
	apiConfig ApiConfig

	// misc
	LastFetchedAt			time.Time
}

type InstanceUsageInfo struct {
	//AvgUsageHours	int
	AvgCpuUtilPercentage int
	AvgNetworkIn int
	AvgNetworkOut int
	AvgEbsReadBytesPerSec int
	AvgEbsWriteBytesPerSec int
	AvgCpuCreditUsagePercentage int // CPUCreditUsage / CPUCreditBalance
	MaxCPUSurplusCreditsCharged int
}

type Recommendation struct {
	InstanceId        string
	Finding string
	RecommendedInstanceTypesWithRank []RecommendedType
	SavingsPercentage int
}

type RecommendedType struct {
	InstanceType string
	Rank string
}

type CostPerPeriod struct {
	Amortized string
	Blended 	string
	//NetAmortized	string
	//NetUnblended	string
	//Unblended	string
}

type ApiConfig struct {
	// API parameters
	Granularity				string
	CostType			    []string
	EvaluationPeriodInDays  int
	MetricStatistic			string
	StartDate		        time.Time
	EndDate		 	  		time.Time
}

type GenericInfo struct {
	// derived values
	SpotRunningPercentage   int
	//CostEffectivenessFactor	int  // percentage, better name?

	// misc
	LastFetchedAt			time.Time
}

func PrintTable(result map[string]InstanceInfo) {
	fmt.Println(len(result))
	instancesInfo := make([]InstanceInfo, 0, len(result))
	for _,v := range result {
		instancesInfo = append(instancesInfo, v)
	}

	var data [][]string
	//var indexedOptions []string

	for _, i := range instancesInfo {
		cef := fmt.Sprintf("%d", getCostEffectivenessFactor(i))
		//ebs := fmt.Sprintf("%s", i.EbsAttached)
		c := fmt.Sprintf("%+v", i.AvgCostPerPeriod)
		fi := fmt.Sprintf("%+v", i.Recommendations.Finding)
		rec := fmt.Sprintf("%+v", i.Recommendations.RecommendedInstanceTypesWithRank)

		cu := fmt.Sprintf("%d", i.UsageInfo.AvgCpuCreditUsagePercentage)
		ni := fmt.Sprintf("%d", i.UsageInfo.AvgNetworkIn)
		no := fmt.Sprintf("%d", i.UsageInfo.AvgNetworkOut)
		cp := fmt.Sprintf("%d", i.UsageInfo.AvgCpuCreditUsagePercentage)
		cs := fmt.Sprintf("%d", i.UsageInfo.MaxCPUSurplusCreditsCharged)

		row := []string{i.InstanceId, i.InstanceType, i.CapacityType, i.Region, fi, rec, cef, c, cu, ni, no, cp, cs}
		//indexedOptions = append(indexedOptions, "Capacity Type")

		data = append(data, row)
	}

	header := []string{"Instance_Id", "Type", "CapacityType", "Region", "Finding", "Recommended Instance Types", "Cost_Effectiveness_Factor", "Avg_Cost", "Avg_Cpu_Utilization %", "Avg_Network_In",
		"Avg_Network_Out", "Avg_Cpu_Credit_Usage %", "Max_CPU_Surplus_Credits_Charged"}

	if data != nil {
		table := table.BuildTable(data, header)
		fmt.Print(table)
		fmt.Println("\n\n")
	}
}

func getCostEffectivenessFactor(info InstanceInfo) int {
	// todo: impl, use recommendation util metrics? learn more.
	return 0
}

func Merge(src map[string]InstanceInfo, dest map[string]InstanceInfo) map[string]InstanceInfo{
	for k, v := range src {
		fmt.Println(k)
		fmt.Printf("%+v", v)
		v.merge(dest[k])
	}

	return src
}

func (insInf InstanceInfo) merge(insInfoToMerge InstanceInfo) {
	if insInf.InstanceId == "" {
		insInf.InstanceId = insInfoToMerge.InstanceId
	}

	if insInf.InstanceType == "" {
		insInf.InstanceType = insInfoToMerge.InstanceType
	}

	if insInf.CapacityType == "" {
		insInf.CapacityType = insInfoToMerge.CapacityType
	}

	if insInf.Region == "" {
		insInf.Region = insInfoToMerge.Region
	}

	if insInf.AvgCostPerPeriod == (CostPerPeriod{}) {
		insInf.AvgCostPerPeriod = insInfoToMerge.AvgCostPerPeriod
	}

	if insInf.UsageInfo == (InstanceUsageInfo{}) {
		insInf.UsageInfo = insInfoToMerge.UsageInfo
	}

	if insInf.Recommendations.InstanceId == "" {
		insInf.Recommendations = insInfoToMerge.Recommendations
	}

	if insInf.CostEffectivenessFactor == 0 {
		insInf.CostEffectivenessFactor = insInfoToMerge.CostEffectivenessFactor
	}
}

