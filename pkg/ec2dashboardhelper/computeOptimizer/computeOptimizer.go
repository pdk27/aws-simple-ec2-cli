package computeOptimizer

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/computeoptimizer"
	"simple-ec2/pkg/ec2dashboardhelper/config"
	"simple-ec2/pkg/ec2dashboardhelper/info"
	"strings"
	"github.com/aws/aws-sdk-go/aws/session"
)


func PopulateRecommendations(cfg config.Config) map[string]info.InstanceInfo {

	sess, _ := session.NewSession(&aws.Config{
		//Region: aws.String(cfg.Region),
	})
	client := computeoptimizer.New(sess)
	instancesInfo := make(map[string]info.InstanceInfo)

	// Hard-code accountID for now
	accountId := aws.String("xxx")
	params := computeoptimizer.GetEC2InstanceRecommendationsInput{
		AccountIds: aws.StringSlice([]string{*accountId}),
	}
	result, err := client.GetEC2InstanceRecommendations(&params)
	if err != nil {
		panic(err)
	}

	//fmt.Println(len(result.InstanceRecommendations))
	for _, res := range result.InstanceRecommendations {
		instanceId := strings.Split(*res.InstanceArn, "/")[1]
		region := strings.Split(*res.InstanceArn, ":")[3]

		var insInf info.InstanceInfo
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
		insInf.Region = region
		instancesInfo[instanceId] = insInf
	}
	//fmt.Println(len(instancesInfo))
	return instancesInfo
}
