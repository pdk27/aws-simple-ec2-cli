package computeOptimizer

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/computeoptimizer"
	"simple-ec2/pkg/ec2dashboardhelper/config"
	"simple-ec2/pkg/ec2dashboardhelper/info"
	"strings"
)


func PopulateRecommendations(cfg config.Config) map[string]info.InstanceInfo {

	client := computeoptimizer.New(cfg.AWSSession, &aws.Config{
		Region: aws.String(cfg.Region),
	})
	instancesInfo := make(map[string]info.InstanceInfo)

	// Hard-code accountID for now
	accountId := aws.String("710616526111")
	params := computeoptimizer.GetEC2InstanceRecommendationsInput{
		AccountIds: aws.StringSlice([]string{*accountId}),
	}
	result, err := client.GetEC2InstanceRecommendations(&params)
	//fmt.Println(result)
	if err != nil {
		panic(err)
	}

	for _, res := range result.InstanceRecommendations {
		instanceId := strings.Split(*res.InstanceArn, "/")[1]
		region := strings.Split(*res.InstanceArn, ":")[3]

		insInf := info.InstanceInfo{
			InstanceId: instanceId,
			Region:     region,
		}

		// add recommendations
		reco := info.Recommendation{
			InstanceId:   instanceId,
			Finding: *res.Finding,
		}
		var rs []info.RecommendedType // collect recommended types and ranks
		for _, g := range res.RecommendationOptions {
			rs = append(rs, info.RecommendedType{
				InstanceType: *g.InstanceType,
				Rank: fmt.Sprintf("%d", *g.Rank)})
		}
		reco.RecommendedInstanceTypesWithRank = rs
		insInf.Recommendations = reco

		instancesInfo[instanceId] = insInf
	}
	return instancesInfo
}
