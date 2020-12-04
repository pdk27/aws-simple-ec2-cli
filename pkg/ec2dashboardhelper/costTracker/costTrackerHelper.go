package costTracker

import (
	"fmt"
	"simple-ec2/pkg/table"
)

type CostTracker struct {
	ByInstanceType []Item
	ByRegion []Item
	ByCapacityType []Item // spot or on-demand
	CostUnit string
	Period int
}

type Item struct {
	Key string // e.g. instance type / region / capacity type
	Value string // total cost
}

type Ec2SpotTracker struct {
	SpotRunningPercentage int // cost explorer API: spot/ total ec2 running
}

//func PopulateCostsByCategories(cfg config.Config, aggregatedRes map[string]info.InstanceInfo) {
//	reqInput := getReqInputWithCommonParams(cfg)
//
//	reqInput.GroupBy = getGroupBys("CapacityType")
//	output := getCostAndUsageOutput(cfg, reqInput)
//	ct := buildCostTracker(output)
//	printCostTrackerTables("CostsByCapacityType", ct.ByCapacityType)
//
//	reqInput.GroupBy = getGroupBys("InstanceType")
//	output = getCostAndUsageOutput(cfg, reqInput)
//	ct = buildCostTracker(output)
//	printCostTrackerTables("CostsByInstanceType", ct.ByInstanceType)
//
//	if cfg.Region == "" {
//		reqInput.GroupBy = getGroupBys("Region")
//		output = getCostAndUsageOutput(cfg, reqInput)
//		ct = buildCostTracker(output)
//		printCostTrackerTables("CostsByRegion", ct.ByRegion)
//	}
//}
//
//func getGroupBys(category string) []*costexplorer.GroupDefinition {
//	var groupBys []*costexplorer.GroupDefinition
//
//	// group by resource id by default
//	groupBys = append(groupBys, &costexplorer.GroupDefinition{
//		Type: aws.String(costexplorer.GroupDefinitionTypeDimension), Key: aws.String(costexplorer.DimensionResourceId)})
//
//	// append ByInstanceType
//	if category == "InstanceType" {
//		groupBys = append(groupBys, &costexplorer.GroupDefinition{
//			Type: aws.String(costexplorer.GroupDefinitionTypeDimension), Key: aws.String(costexplorer.DimensionInstanceType)})
//	}
//
//	//append ByRegion
//	if category == "Region" {
//		groupBys = append(groupBys, &costexplorer.GroupDefinition{
//			Type: aws.String(costexplorer.GroupDefinitionTypeDimension), Key: aws.String(costexplorer.DimensionRegion)})
//	}
//
//	//append ByCapacityType
//	if category == "CapacityType" {
//		groupBys = append(groupBys, &costexplorer.GroupDefinition{
//			Type: aws.String(costexplorer.GroupDefinitionTypeDimension), Key: aws.String(costexplorer.DimensionPurchaseType)})
//	}
//	return groupBys
//}

//func buildCostTracker(output costexplorer.GetCostAndUsageWithResourcesOutput) CostTracker{
//	fmt.Println(output)
//
//
//	return CostTracker{}
//}

func printCostTrackerTables(category string, ctItems []Item) {
	var data [][]string
	indexedOptions := []string{}

	for _, cti := range ctItems {
		row := []string{category, cti.Key, cti.Value}
		indexedOptions = append(indexedOptions, category)
		data = append(data, row)
	}

	header := []string{"Category", "Key", "Value"} // InstanceType // Region // CapacityType
	table := table.BuildTable(data, header)

	fmt.Println(category)
	fmt.Print(table)
}
