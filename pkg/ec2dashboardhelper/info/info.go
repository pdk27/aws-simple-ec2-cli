package info

import (
	"fmt"
	"simple-ec2/pkg/table"
	"strconv"
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
	CpuUtilPercentage Metric
	NetworkIn Metric
	NetworkOut Metric
	EbsReadBytesPerSec Metric
	EbsWriteBytesPerSec Metric
	CpuCreditUsagePercentage Metric // CPUCreditUsage / CPUCreditBalance
	MaxCPUSurplusCreditsCharged string
}

type Metric struct {
	Avg string
	Max string
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
	instancesInfo := make([]InstanceInfo, 0, len(result))
	for _,v := range result {
		instancesInfo = append(instancesInfo, v)
	}
	//fmt.Printf("%+v", instancesInfo)

	var data [][]string
	//var indexedOptions []string

	for _, i := range instancesInfo {
		cef := fmt.Sprintf("%0.4f", getCostEffectivenessFactor(i))
		//ebs := fmt.Sprintf("%s", i.EbsAttached)
		c := fmt.Sprintf("%+v", i.AvgCostPerPeriod)
		fi := fmt.Sprintf("%+v", i.Recommendations.Finding)
		rec := fmt.Sprintf("%+v", i.Recommendations.RecommendedInstanceTypesWithRank)

		cu := fmt.Sprintf("%+v", i.UsageInfo.CpuUtilPercentage)
		ni := fmt.Sprintf("%+v", i.UsageInfo.NetworkIn)
		no := fmt.Sprintf("%+v", i.UsageInfo.NetworkOut)
		cp := fmt.Sprintf("%+v", i.UsageInfo.CpuCreditUsagePercentage)
		cs := fmt.Sprintf("%+v", i.UsageInfo.MaxCPUSurplusCreditsCharged)

		row := []string{i.InstanceId, i.InstanceType, i.CapacityType, i.Region, fi, rec, cef, c, cu, ni, no, cp, cs}
		//indexedOptions = append(indexedOptions, "Capacity Type")

		data = append(data, row)
	}

	header := []string{"Instance Id", "Type", "Capacity Type", "Region", "Finding", "Recommended Instance Types", "Cost Effectiveness Factor", "Avg Cost", "Avg Cpu Utilization %", "Avg Network In",
		"Avg Network Out", "Avg Cpu Credit Usage %", "Max CPU Surplus Credits Charged"}

	if data != nil {
		table := table.BuildTable(data, header)
		fmt.Print(table)
		fmt.Println("\n\n")
	}
}

func getCostEffectivenessFactor(info InstanceInfo) float64 {
	// todo: impl, use recommendation util metrics? learn more.

	// 80% cpu = avg cost
	// max cpu util = ?

	cost, _ := strconv.ParseFloat(info.AvgCostPerPeriod.Blended, 64)
	cpuUtil, _ := strconv.ParseFloat(info.UsageInfo.CpuUtilPercentage.Max, 64)

	a := cost / 80
	res :=  a * cpuUtil * 100

	return res
}

func Merge(src map[string]InstanceInfo, dest map[string]InstanceInfo) map[string]InstanceInfo{
	//fmt.Printf("Merging %d records to result\n\n", len(src))
	for id, srcInfo := range src {
		destInfo, found := dest[id]
		if !found {
			dest[id] = srcInfo
		} else {
			info := destInfo.merge(srcInfo)
			dest[id] = info
		}
	}

	return dest
}

func (insInf InstanceInfo) merge(sourceInsInfo InstanceInfo) InstanceInfo {
	//fmt.Println("SRC: %+v", sourceInsInfo)
	//fmt.Println("before merge: %+v", insInf)
	result := insInf
	if result.InstanceId == "" {
		//fmt.Println("merging id")
		result.InstanceId = sourceInsInfo.InstanceId
	}

	if result.InstanceType == "" {
		//fmt.Println("merging type")
		result.InstanceType = sourceInsInfo.InstanceType
	}

	if result.CapacityType == "" {
		//fmt.Println("merging CapacityType")
		result.CapacityType = sourceInsInfo.CapacityType
	}

	if result.Region == "" {
		//fmt.Println("merging region")
		result.Region = sourceInsInfo.Region
	}

	if result.AvgCostPerPeriod == (CostPerPeriod{}) {
		//fmt.Println("merging AvgCostPerPeriod")
		result.AvgCostPerPeriod = sourceInsInfo.AvgCostPerPeriod
	}

	if result.UsageInfo == (InstanceUsageInfo{}) {
		//fmt.Println("merging UsageInfo")
		result.UsageInfo = sourceInsInfo.UsageInfo
	}

	if result.Recommendations.InstanceId == "" {
		//fmt.Println("merging Recommendations")
		result.Recommendations = sourceInsInfo.Recommendations
	}

	if result.CostEffectivenessFactor == 0 {
		//fmt.Println("merging CostEffectivenessFactor")
		result.CostEffectivenessFactor = sourceInsInfo.CostEffectivenessFactor
	}

	//fmt.Println("after merge: %+v", result)
	return result
}

