package cloudwatch


import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"simple-ec2/pkg/ec2dashboardhelper/config"
	"simple-ec2/pkg/ec2dashboardhelper/info"
	"strconv"
	"time"
)

const (
	timeLayout = "2006-01-02T15:04:05Z"
)

func PopulateMetrics(instanceIds []string, cfg config.Config) map[string]info.InstanceInfo {

	client := cloudwatch.New(cfg.AWSSession, &aws.Config{
		Region: aws.String(cfg.Region),
	})
	instancesInfo := make(map[string]info.InstanceInfo)

	for _, id := range instanceIds {
		instancesInfo[id] = getMetricsForId(client, id, cfg.EvaluationPeriodInDays)
	}

	return instancesInfo
}

func getMetricsForId(client *cloudwatch.CloudWatch, id string, evalPeriodDays int) info.InstanceInfo {
	end := time.Now()
	start := end.Add(time.Duration(-evalPeriodDays) * 24 * time.Hour)
	period, _ := strconv.ParseInt("604800", 10, 64) // 5 mins

	idDimension := cloudwatch.Dimension{
		Name: aws.String("InstanceId"),
		Value: aws.String(id),
	}
	input1 := cloudwatch.GetMetricStatisticsInput{
		Namespace: aws.String("AWS/EC2"),
		MetricName: aws.String("CPUUtilization"),
		//"NetworkIn", "NetworkOut", "CPUCreditUsage", "CPUSurplusCreditsCharged"),
		Statistics: []*string{aws.String("Average"), aws.String("Maximum")},
		StartTime: &start,
		EndTime: &end,
		Period: &period,
		Dimensions: []*cloudwatch.Dimension{&idDimension},
	}

	//avg := aws.String("Average")
	//max := aws.String("Maximum")

	//cpuAvg := cloudwatch.MetricDataQuery{
	//	Id:         aws.String("m1"),
	//	MetricStat: &cloudwatch.MetricStat{
	//		Metric: &cloudwatch.Metric{
	//			//Dimensions: nil,
	//			MetricName: aws.String("CPUUtilization"),
	//			Namespace:  ec2Namespace,
	//		},
	//		Period: &period,
	//		Stat:   avg,
	//	},
	//	Expression: nil,
	//	Label:      nil,
	//	ReturnData: aws.Bool(false),
	//}

	//cpuAvg := getMetricQuery("cAvg", "CPUUtilization", avg, &idDimension)
	//cpuMax := getMetricQuery("cMax", "CPUUtilization", max, &idDimension)
	//ninAvg := getMetricQuery("inAvg", "NetworkIn", avg, &idDimension)
	//ninMax := getMetricQuery("inMax", "NetworkIn", max, &idDimension)
	//noutAvg := getMetricQuery("outAvg", "NetworkOut", avg, &idDimension)
	//noutMax := getMetricQuery("outMax", "NetworkOut", max, &idDimension)

	//input := cloudwatch.GetMetricDataInput{
	//	MetricDataQueries: []*cloudwatch.MetricDataQuery{cpuAvg, cpuMax},// cpuMax, ninAvg, ninMax, noutAvg, noutMax},
	//	StartTime:         aws.Time(start),
	//	EndTime:           aws.Time(end),
	//}
	//result, err := client.GetMetricData(&input)

	result, err := client.GetMetricStatistics(&input1)
	//fmt.Println("result", result)
	if err != nil {
		panic(err)
	}

	usageInfo :=  info.InstanceUsageInfo{}
	//for _, m := range result.MetricDataResults {
	//	switch *m.Id {
	//	case "cAvg": usageInfo.CpuUtilPercentage.Avg = fmt.Sprintf("%.4f", *m.Values[0])
	//	case "cMax": usageInfo.CpuUtilPercentage.Max = fmt.Sprintf("%.4f", *m.Values[0])
	//	}
	//}
	usageInfo.CpuUtilPercentage.Avg = fmt.Sprintf("%.4f", *result.Datapoints[0].Average)
	usageInfo.CpuUtilPercentage.Max = fmt.Sprintf("%.4f", *result.Datapoints[0].Maximum)

	return info.InstanceInfo{
		InstanceId: id,
		UsageInfo:  usageInfo,
	}
}

func getMetricQuery(id string, metricName string, stat *string, idDimension *cloudwatch.Dimension) *cloudwatch.MetricDataQuery {
	period, _ := strconv.ParseInt("604800", 10, 64) // 5 mins
	return &cloudwatch.MetricDataQuery{
		Id:         aws.String(id),
		MetricStat: &cloudwatch.MetricStat{
			Metric: &cloudwatch.Metric{
				Dimensions: []*cloudwatch.Dimension{idDimension},
				MetricName: aws.String(metricName),
				Namespace:  aws.String("AWS/EC2"),
			},
			Period: &period,
			Stat:   stat,
		},
		ReturnData: aws.Bool(false),
	}

}
