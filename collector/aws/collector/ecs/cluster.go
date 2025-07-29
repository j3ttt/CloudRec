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

package ecs

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/cloudrec/aws/collector"
	"github.com/core-sdk/constant"
	"github.com/core-sdk/log"
	"github.com/core-sdk/schema"
	"go.uber.org/zap"
	"sync"
)

const (
	maxWorkers = 10
)

// GetClusterResource returns a Cluster Resource
func GetClusterResource() schema.Resource {
	return schema.Resource{
		ResourceType:       collector.ECSCluster,
		ResourceTypeName:   "ECS Cluster",
		ResourceGroupType:  constant.CONTAINER,
		Desc:               `https://docs.aws.amazon.com/AmazonECS/latest/APIReference/API_DescribeClusters.html`,
		ResourceDetailFunc: GetClusterDetail,
		RowField: schema.RowField{
			ResourceId:   "$.Cluster.ClusterArn",
			ResourceName: "$.Cluster.ClusterName",
		},
		Dimension: schema.Regional,
	}
}

// ClusterDetail aggregates all information for a single ECS cluster.
type ClusterDetail struct {
	Cluster types.Cluster
}

// GetClusterDetail fetches the details for all ECS clusters in a region.
func GetClusterDetail(ctx context.Context, service schema.ServiceInterface, res chan<- any) error {
	client := service.(*collector.Services).ECS

	clusterArns, err := listClusters(ctx, client)
	if err != nil {
		log.CtxLogger(ctx).Error("failed to list ecs clusters", zap.Error(err))
		return err
	}

	var wg sync.WaitGroup
	tasks := make(chan string, len(clusterArns))

	// Start workers
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for arn := range tasks {
				describedClusters, err := describeClusters(ctx, client, []string{arn})
				if err != nil {
					log.CtxLogger(ctx).Warn("failed to describe ecs cluster", zap.String("arn", arn), zap.Error(err))
					continue
				}
				if len(describedClusters) > 0 {
					res <- ClusterDetail{Cluster: describedClusters[0]}
				}
			}
		}()
	}

	// Add tasks to the queue
	for _, arn := range clusterArns {
		tasks <- arn
	}
	close(tasks)

	wg.Wait()

	return nil
}

// listClusters retrieves all ECS cluster ARNs in a region.
func listClusters(ctx context.Context, c *ecs.Client) ([]string, error) {
	var clusterArns []string
	paginator := ecs.NewListClustersPaginator(c, &ecs.ListClustersInput{})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		clusterArns = append(clusterArns, page.ClusterArns...)
	}
	return clusterArns, nil
}

// describeClusters retrieves the details for a list of clusters.
func describeClusters(ctx context.Context, c *ecs.Client, clusterArns []string) ([]types.Cluster, error) {
	output, err := c.DescribeClusters(ctx, &ecs.DescribeClustersInput{Clusters: clusterArns, Include: []types.ClusterField{types.ClusterFieldTags, types.ClusterFieldSettings}})
	if err != nil {
		return nil, err
	}
	return output.Clusters, nil
}
