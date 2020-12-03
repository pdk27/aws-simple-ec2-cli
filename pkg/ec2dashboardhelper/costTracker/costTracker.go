package costTracker

import (
	"fmt"
	"simple-ec2/pkg/ec2dashboardhelper/info"
	"strconv"
	"strings"
	"time"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/costexplorer"

	"simple-ec2/pkg/ec2dashboardhelper/config"
)
//
//type CostTracker struct {
//	//ByInstanceType []Item
//	//ByRegion []Item
//	ByCapacityType []Item // spot or on-demand
//	ByResource []Item
//	CostUnit string
//	Period int
//}
//
//type Item struct {
//	Key string
//	Value string
//}
//
//type Ec2SpotTracker struct {
//	SpotRunningPercentage int // cost explorer API: spot/ total ec2 running
//}

//https://aws.amazon.com/blogs/aws-cost-management/understanding-your-aws-cost-datasets-a-cheat-sheet/
type CostType int
const (
	AmortizedCost = iota
	BlendedCost
	NetAmortizedCost
	NetUnblendedCost
	UnblendedCost
)
func (c CostType) String() string {
	return [...]string{"AmortizedCost", "BlendedCost", "NetAmortizedCost", "NetUnblendedCost", "UnblendedCost"}[c]
}

const (
    layout = "2006-01-02"
)

func PopulateCosts(cfg config.Config, instancesInfo []info.InstanceInfo) []info.InstanceInfo {
	client := costexplorer.New(cfg.AWSSession)

	var groupBys []*costexplorer.GroupDefinition

	// group by resource is a must
	groupBys = append(groupBys, &costexplorer.GroupDefinition{
		Type: aws.String(costexplorer.GroupDefinitionTypeDimension), Key: aws.String(costexplorer.DimensionResourceId)})

	// append ByInstanceType
	groupBys = append(groupBys, &costexplorer.GroupDefinition{
		Type: aws.String(costexplorer.GroupDefinitionTypeDimension), Key: aws.String(costexplorer.DimensionInstanceType)})

	// append ByRegion
	//groupBys = append(groupBys, &costexplorer.GroupDefinition{
	//	Type: aws.String(costexplorer.GroupDefinitionTypeDimension), Key: aws.String(costexplorer.DimensionRegion)})

	// append ByCapacityType
	//groupBys = append(groupBys, &costexplorer.GroupDefinition{
	//	Type: aws.String(costexplorer.GroupDefinitionTypeDimension), Key: aws.String(costexplorer.DimensionPurchaseType)})

	return getCostsWithGroupBy(client, cfg, groupBys, instancesInfo)
}

func getCostsWithGroupBy(
	client *costexplorer.CostExplorer,
	cfg config.Config,
	groupBy []*costexplorer.GroupDefinition,
	instancesInfoInput []info.InstanceInfo) []info.InstanceInfo {

	// todo: validations

	var costType []*string
	for _, ct := range strings.Split(cfg.CostType, ",") {
		costType = append(costType, aws.String(ct))
	}

	end := time.Now()
	start := end.Add(time.Duration(-cfg.EvaluationPeriodInDays) * 24 * time.Hour)
	timePeriod := &costexplorer.DateInterval{
		Start: aws.String(start.Format(layout)),
		End:   aws.String(end.Format(layout)),
	}

	// always use the filters
	var andFilters []*costexplorer.Expression
	svcFilter := &costexplorer.Expression{
		Dimensions: &costexplorer.DimensionValues{
			Key: aws.String("SERVICE"),
			Values: []*string{aws.String("Amazon Elastic Compute Cloud - Compute")},
		},
	}
	andFilters = append(andFilters, svcFilter)

	usageTypeGroupFilter := &costexplorer.Expression{
		Dimensions: &costexplorer.DimensionValues{
			Key: aws.String("USAGE_TYPE_GROUP"),
			Values: []*string{aws.String("EC2: Running Hours")},
		},
	}
	andFilters = append(andFilters, usageTypeGroupFilter)

	if cfg.Region == "" {
		regionFilter := &costexplorer.Expression{
			Dimensions: &costexplorer.DimensionValues{
				Key: aws.String("REGION"),
				Values: []*string{aws.String(cfg.Region)},
			},
		}
		andFilters = append(andFilters, regionFilter)
	}
	filters := &costexplorer.Expression{And: andFilters}

	//https://docs.aws.amazon.com/awsaccountbilling/latest/aboutv2/ce-filtering.html
	reqInput := costexplorer.GetCostAndUsageWithResourcesInput{
		Granularity: aws.String(cfg.Granularity),
		Metrics:     costType,
		GroupBy:     groupBy,
		TimePeriod:  timePeriod,
		Filter:		 filters,
	}

	result, err := client.GetCostAndUsageWithResources(&reqInput)

	if err != nil {
		panic(err)
	}

	//fmt.Printf("%+v\n\n\n", result)

	instanceTypeById := make(map[string]string) // {instanceId: X, type: Y}
	blendedCostsPerInstance := make(map[string][]string) // {instanceId: X, blendedCost: Y}
	unblendedCostsPerInstance := make(map[string][]string)// {instanceId: X, unblendedCost: Y}
	for _, res := range result.ResultsByTime {
		var id, iType string
		for _, g := range res.Groups {
			keys := g.Keys
			if strings.Contains(*keys[0], "i-") {
				id = *keys[0]
				iType = *keys[1]
			} else {
				id = *keys[1]
				iType = *keys[0]
			}
			instanceTypeById[id] = iType

			// populate the map to calculate avg later
			bCosts := blendedCostsPerInstance[id]
			bCosts = append(bCosts, *g.Metrics["BlendedCost"].Amount)
			blendedCostsPerInstance[id] = bCosts

			ubCosts := unblendedCostsPerInstance[id]
			ubCosts = append(ubCosts, *g.Metrics["UnblendedCost"].Amount)
			unblendedCostsPerInstance[id] = ubCosts
		}
	}

	var instancesInfo []info.InstanceInfo
	//fmt.Println(len(instancesInfoInput), instancesInfoInput)
	if len(instancesInfoInput) != 0 {
		for _, i := range instancesInfoInput {
			var blendedSum, unblendedSum float64
			for _, c := range blendedCostsPerInstance[i.InstanceId] {
				cost, err := strconv.ParseFloat(c, 64)
				if err == nil {
					blendedSum += cost
				}
			}
			blendedAvg := (blendedSum) / float64(len(blendedCostsPerInstance[i.InstanceId]))

			for _, c := range unblendedCostsPerInstance[i.InstanceId] {
				cost, err := strconv.ParseFloat(c, 64)
				if err == nil {
					unblendedSum += cost
				}
			}
			unblendedAvg := (unblendedSum) / float64(len(unblendedCostsPerInstance[i.InstanceId]))

			i.AvgCostPerPeriod = info.AvgCostPerPeriod{
				Blended: fmt.Sprintf("%.4f", blendedAvg),
				Unblended: fmt.Sprintf("%.4f", unblendedAvg),
			}
		}
		instancesInfo = instancesInfoInput
	} else {
		// temporary build logic
		for id, iType := range instanceTypeById {
			var blendedSum, unblendedSum float64
			for _, c := range blendedCostsPerInstance[id] {
				cost, err := strconv.ParseFloat(c, 64)
				if err == nil {
					blendedSum += cost
				}
			}
			blendedAvg := (blendedSum) / float64(len(blendedCostsPerInstance[id]))

			for _, c := range unblendedCostsPerInstance[id] {
				cost, err := strconv.ParseFloat(c, 64)
				if err == nil {
					unblendedSum += cost
				}
			}
			unblendedAvg := (unblendedSum) / float64(len(unblendedCostsPerInstance[id]))

			instancesInfo = append(instancesInfo, info.InstanceInfo{
				InstanceId: id,
				InstanceType: iType,
				Region: cfg.Region,
				AvgCostPerPeriod: info.AvgCostPerPeriod{
					Blended: fmt.Sprintf("%.4f", blendedAvg),
					Unblended: fmt.Sprintf("%.4f", unblendedAvg),
				},
			})
		}
	}

	return instancesInfo
}