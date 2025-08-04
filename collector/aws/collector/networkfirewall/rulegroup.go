// Licensed to the Apache Software Foundation (ASF) under one or more
// contributor license agreements.  See the NOTICE file distributed with
// this work for additional information regarding copyright ownership.
// The ASF licenses this file to You under the Apache License, Version 2.0
// (the "License"); you may not use this file except in compliance with
// the License.  You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package networkfirewall

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/networkfirewall"
	"github.com/aws/aws-sdk-go-v2/service/networkfirewall/types"
	"github.com/cloudrec/aws/collector"
	"github.com/core-sdk/constant"
	"github.com/core-sdk/log"
	"github.com/core-sdk/schema"
	"go.uber.org/zap"
)


// GetRuleGroupResource returns AWS Network Firewall Rule Group resource definition
func GetRuleGroupResource() schema.Resource {
	return schema.Resource{
		ResourceType:       collector.NetworkFirewall,
		ResourceTypeName:   "Network Firewall Rule Group",
		ResourceGroupType:  constant.NET,
		Desc:               "https://docs.aws.amazon.com/network-firewall/latest/APIReference/API_RuleGroup.html",
		ResourceDetailFunc: GetRuleGroupDetail,
		RowField: schema.RowField{
			ResourceId:   "$.RuleGroup.RuleGroupResponse.Arn",
			ResourceName: "$.RuleGroup.RuleGroupResponse.Name",
		},
		Dimension: schema.Regional,
	}
}

// RuleGroupDetail aggregates all information for a single Network Firewall Rule Group.
type RuleGroupDetail struct {
	RuleGroup *networkfirewall.DescribeRuleGroupOutput
	Tags      map[string]string
}

// GetRuleGroupDetail fetches the details for all Network Firewall Rule Groups in a region.
func GetRuleGroupDetail(ctx context.Context, service schema.ServiceInterface, res chan<- any) error {
	client := service.(*collector.Services).NetworkFirewall

	ruleGroups, err := listRuleGroups(ctx, client)
	if err != nil {
		log.CtxLogger(ctx).Error("failed to list Network Firewall Rule Groups", zap.Error(err))
		return err
	}

	var wg sync.WaitGroup
	tasks := make(chan types.RuleGroupMetadata, len(ruleGroups))

	// Start worker goroutines
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for ruleGroup := range tasks {
				detail := describeRuleGroupDetail(ctx, client, ruleGroup)
				if detail != nil {
					res <- detail
				}
			}
		}()
	}

	// Add tasks
	for _, ruleGroup := range ruleGroups {
		tasks <- ruleGroup
	}
	close(tasks)

	wg.Wait()
	return nil
}

// listRuleGroups retrieves all Network Firewall Rule Groups in a region.
func listRuleGroups(ctx context.Context, c *networkfirewall.Client) ([]types.RuleGroupMetadata, error) {
	var ruleGroups []types.RuleGroupMetadata
	input := &networkfirewall.ListRuleGroupsInput{
		MaxResults: aws.Int32(100),
	}

	paginator := networkfirewall.NewListRuleGroupsPaginator(c, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		ruleGroups = append(ruleGroups, page.RuleGroups...)
	}
	return ruleGroups, nil
}

// describeRuleGroupDetail fetches all details for a single rule group.
func describeRuleGroupDetail(ctx context.Context, client *networkfirewall.Client, ruleGroup types.RuleGroupMetadata) *RuleGroupDetail {
	// Get detailed rule group information
	describeInput := &networkfirewall.DescribeRuleGroupInput{
		RuleGroupArn: ruleGroup.Arn,
	}
	describeOutput, err := client.DescribeRuleGroup(ctx, describeInput)
	if err != nil {
		log.CtxLogger(ctx).Error("failed to describe Network Firewall Rule Group", zap.String("arn", *ruleGroup.Arn), zap.Error(err))
		return nil
	}

	var tags map[string]string

	// Get tags - Network Firewall Rule Group doesn't have a direct API to list tags
	// but we can extract relevant information from the rule group itself
	tags = extractRuleGroupTags(describeOutput)

	return &RuleGroupDetail{
		RuleGroup: describeOutput,
		Tags:      tags,
	}
}

// extractRuleGroupTags extracts relevant information from a rule group as tags
func extractRuleGroupTags(ruleGroup *networkfirewall.DescribeRuleGroupOutput) map[string]string {
	tags := make(map[string]string)

	// Extract some key information from the rule group as tags
	if ruleGroup.RuleGroup != nil {
		if ruleGroup.RuleGroup.RuleVariables != nil {
			tags["HasRuleVariables"] = "true"
		} else {
			tags["HasRuleVariables"] = "false"
		}
		
		// Add rule count information
		if ruleGroup.RuleGroup.RulesSource != nil {
			if ruleGroup.RuleGroup.RulesSource.RulesString != nil {
				tags["RulesSourceType"] = "RulesString"
			} else if ruleGroup.RuleGroup.RulesSource.StatefulRules != nil {
				tags["RulesSourceType"] = "StatefulRules"
				tags["StatefulRulesCount"] = string(rune(len(ruleGroup.RuleGroup.RulesSource.StatefulRules)))
			} else if ruleGroup.RuleGroup.RulesSource.StatelessRulesAndCustomActions != nil {
				tags["RulesSourceType"] = "StatelessRulesAndCustomActions"
				if ruleGroup.RuleGroup.RulesSource.StatelessRulesAndCustomActions.StatelessRules != nil {
					tags["StatelessRulesCount"] = string(rune(len(ruleGroup.RuleGroup.RulesSource.StatelessRulesAndCustomActions.StatelessRules)))
				}
			}
		}
		
		// Add capacity information
		if ruleGroup.RuleGroupResponse.Capacity != nil {
			tags["Capacity"] = string(rune(*ruleGroup.RuleGroupResponse.Capacity))
		}
		
		// Add type information
		tags["Type"] = string(ruleGroup.RuleGroupResponse.Type)
	}

	return tags
}