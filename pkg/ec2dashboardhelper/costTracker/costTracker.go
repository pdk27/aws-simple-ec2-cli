package costTracker

import (
	"fmt"
	"simple-ec2/pkg/ec2dashboardhelper/info"
	"strconv"
	"strings"
	"regexp"
	"time"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/costexplorer"

	"simple-ec2/pkg/ec2dashboardhelper/config"
)

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
    regionPattern = "(us(-gov)?|ap|ca|cn|eu|sa)-(central|(north|south)?(east|west)?)-\\d"
    instanceTypePattern = "[a-zA-Z][0-9][a-zA-Z]*\\.[0-9]*[a-zA-Z]+" // i-09bdb27b3fcb0cf92 t2.small
)

var	client *costexplorer.CostExplorer

func PopulateCostsAndType(cfg config.Config) map[string]info.InstanceInfo {
	var groupBys []*costexplorer.GroupDefinition

	client = costexplorer.New(cfg.AWSSession, &aws.Config{
		Region: aws.String(cfg.Region),
	})

	// group by resource is a must
	groupBys = append(groupBys, &costexplorer.GroupDefinition{
		Type: aws.String(costexplorer.GroupDefinitionTypeDimension), Key: aws.String(costexplorer.DimensionResourceId)})

	// append ByInstanceType
	groupBys = append(groupBys, &costexplorer.GroupDefinition{
		Type: aws.String(costexplorer.GroupDefinitionTypeDimension), Key: aws.String(costexplorer.DimensionInstanceType)})

	return getCostsWithGroupBy(cfg, groupBys)
}

func PopulateCapacityType(cfg config.Config) map[string]info.InstanceInfo {
	var groupBys []*costexplorer.GroupDefinition

	// group by resource is a must
	groupBys = append(groupBys, &costexplorer.GroupDefinition{
		Type: aws.String(costexplorer.GroupDefinitionTypeDimension), Key: aws.String(costexplorer.DimensionResourceId)})

	// append ByCapacityType
	groupBys = append(groupBys, &costexplorer.GroupDefinition{
		Type: aws.String(costexplorer.GroupDefinitionTypeDimension), Key: aws.String(costexplorer.DimensionPurchaseType)})

	return getCostsWithGroupBy(cfg, groupBys)
}

func PopulateRegion(cfg config.Config) map[string]info.InstanceInfo {
	var groupBys []*costexplorer.GroupDefinition

	// group by resource is a must
	groupBys = append(groupBys, &costexplorer.GroupDefinition{
		Type: aws.String(costexplorer.GroupDefinitionTypeDimension), Key: aws.String(costexplorer.DimensionResourceId)})

	// append ByRegion
	groupBys = append(groupBys, &costexplorer.GroupDefinition{
		Type: aws.String(costexplorer.GroupDefinitionTypeDimension), Key: aws.String(costexplorer.DimensionRegion)})
	return getCostsWithGroupBy(cfg, groupBys)
}

func getCostsWithGroupBy(
	cfg config.Config,
	groupBy []*costexplorer.GroupDefinition) map[string]info.InstanceInfo {

	// todo: validations

	reqInput := getReqInputWithCommonParams(cfg, groupBy)
	result := getCostAndUsageOutput(reqInput)
	//fmt.Println(result)

	instancesInfo, blendedCostsById, amortizedCostsById := populateCommonInstanceInfo(cfg.Region, result.ResultsByTime)
	//fmt.Println(blendedCostsById)
	//fmt.Println(amortizedCostsById)
	instancesInfo = populateCosts(instancesInfo, blendedCostsById, amortizedCostsById)

	return instancesInfo
}

func getReqInputWithCommonParams(cfg config.Config, groupBy []*costexplorer.GroupDefinition) costexplorer.GetCostAndUsageWithResourcesInput {

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
	//https://docs.aws.amazon.com/awsaccountbilling/latest/aboutv2/ce-filtering.html
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

	if cfg.Region != "" {
		regionFilter := &costexplorer.Expression{
			Dimensions: &costexplorer.DimensionValues{
				Key: aws.String("REGION"),
				Values: []*string{aws.String(cfg.Region)},
			},
		}
		andFilters = append(andFilters, regionFilter)
	}
	filters := &costexplorer.Expression{And: andFilters}

	//fmt.Println(cfg)
	return costexplorer.GetCostAndUsageWithResourcesInput{
		Granularity: aws.String(cfg.Granularity),
		Metrics:     costType,
		GroupBy: groupBy,
		TimePeriod:  timePeriod,
		Filter:		 filters,
	}
}

func getCostAndUsageOutput(
	input costexplorer.GetCostAndUsageWithResourcesInput) costexplorer.GetCostAndUsageWithResourcesOutput {
	output, err := client.GetCostAndUsageWithResources(&input)

	if err != nil {
		panic(err)
	}

	//fmt.Println(output)
	return *output
}

func populateCommonInstanceInfo(regionFilter string, results []*costexplorer.ResultByTime) (map[string]info.InstanceInfo, map[string][]string, map[string][]string) {
	//instanceTypeById := make(map[string]string) // {instanceId: X, type: Y}
	blendedCostsPerInstance := make(map[string][]string) // {instanceId: X, blendedCost: Y}
	//unblendedCostsPerInstance := make(map[string][]string)// {instanceId: X, unblendedCost: Y}
	amortizedCostsPerInstance := make(map[string][]string)// {instanceId: X, amortizedCost: Y}

	instancesInfo := make(map[string]info.InstanceInfo)
	for _, res := range results {
		for _, g := range res.Groups {
			var info info.InstanceInfo
			var id string

			// there can be a max of 2 keys
			for _, k := range g.Keys {
				regionMatched, _ := regexp.Match(regionPattern, []byte(*k))
				typeMatched, _ := regexp.Match(instanceTypePattern, []byte(*k))
				fmt.Println(id, *k)

				// instance id
				if strings.Contains(*k, "i-") {
					id = *k
					info.InstanceId = *k
				} else if strings.Contains(*k, "On Demand") || strings.Contains(*k, "Spot")  {
					info.CapacityType = *k
				}  else if regionMatched {
					info.Region = *k
				} else if typeMatched {
					// instance type
					info.InstanceType = *k
				} else {
					// nothing
				}
			}
			if regionFilter != "" {
				info.Region = regionFilter
			}
			instancesInfo[id] = info

			// populate the map to calculate avg later
			bCosts := blendedCostsPerInstance[id]
			bCosts = append(bCosts, *g.Metrics["BlendedCost"].Amount)
			blendedCostsPerInstance[id] = bCosts

			//ubCosts := unblendedCostsPerInstance[id]
			//ubCosts = append(ubCosts, *g.Metrics["UnblendedCost"].Amount)
			//unblendedCostsPerInstance[id] = ubCosts

			amCosts := amortizedCostsPerInstance[id]
			amCosts = append(amCosts, *g.Metrics["AmortizedCost"].Amount)
			amortizedCostsPerInstance[id] = amCosts
		}
	}

	return instancesInfo, blendedCostsPerInstance, amortizedCostsPerInstance
}

func populateCosts(instancesInfo map[string]info.InstanceInfo, blendedCostsById map[string][]string, amortizedCostsById map[string][]string) map[string]info.InstanceInfo {
	var result = make(map[string]info.InstanceInfo)
	for id, inf := range instancesInfo {
		// add costs
		var blendedSum, amortizedSum float64
		for _, c := range blendedCostsById[id] {
			cost, err := strconv.ParseFloat(c, 64)
			if err == nil {
				blendedSum += cost
			}
		}
		blendedAvg := (blendedSum) / float64(len(blendedCostsById[id]))

		//for _, c := range unblendedCostsPerInstance[i.InstanceId] {
		//	cost, err := strconv.ParseFloat(c, 64)
		//	if err == nil {
		//		unblendedSum += cost
		//	}
		//}
		//unblendedAvg := (unblendedSum) / float64(len(unblendedCostsPerInstance[i.InstanceId]))

		for _, c := range amortizedCostsById[id] {
			cost, err := strconv.ParseFloat(c, 64)
			if err == nil {
				amortizedSum += cost
			}
		}
		amortizedAvg := (amortizedSum) / float64(len(amortizedCostsById[id]))

		inf.AvgCostPerPeriod = info.CostPerPeriod{
			Blended: fmt.Sprintf("%.4f", blendedAvg),
			Amortized: fmt.Sprintf("%.4f", amortizedAvg),
		}
		inf.TotalCostPerPeriod = info.CostPerPeriod{
			Blended: fmt.Sprintf("%.4f", blendedSum),
			Amortized: fmt.Sprintf("%.4f", amortizedSum),
		}
		result[id] = inf
	}
	return result
}
