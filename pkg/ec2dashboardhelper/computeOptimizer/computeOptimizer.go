package computeOptimizer

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/computeoptimizer"
	"simple-ec2/pkg/ec2dashboardhelper/config"
	"simple-ec2/pkg/ec2dashboardhelper/info"
	"strings"
)


func GetRecommendation(cfg config.Config, instancesInfo []info.InstanceInfo) []info.InstanceInfo {
	client := computeoptimizer.New(cfg.AWSSession)

	// Hard-code accountID for now
	accountId := aws.String("xxxxxx")
	params := computeoptimizer.GetEC2InstanceRecommendationsInput{
		AccountIds: aws.StringSlice([]string{*accountId}),
	}
	result, err := client.GetEC2InstanceRecommendations(&params)
	if err != nil {
		panic(err)
	}

	for _, res := range result.InstanceRecommendations {
		instanceId := strings.Split(*res.InstanceArn, "/")[1]
		for _, g := range res.RecommendationOptions {
			instancesInfo = append(instancesInfo, info.InstanceInfo{
				InstanceId:   instanceId,
				Recommendation: *g.InstanceType,
			})
		}
	}
	return instancesInfo
}
