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

// GetServiceResource returns a Service Resource
func GetServiceResource() schema.Resource {
	return schema.Resource{
		ResourceType:       collector.ECSService,
		ResourceTypeName:   "ECS Service",
		ResourceGroupType:  constant.CONTAINER,
		Desc:               `https://docs.aws.amazon.com/AmazonECS/latest/APIReference/API_DescribeServices.html`,
		ResourceDetailFunc: GetServiceDetail,
		RowField: schema.RowField{
			ResourceId:   "$.Service.ServiceArn",
			ResourceName: "$.Service.ServiceName",
		},
		Dimension: schema.Regional,
	}
}

// ServiceDetail aggregates all information for a single ECS service.
type ServiceDetail struct {
	Service types.Service
}

// GetServiceDetail fetches the details for all ECS services in a region.
func GetServiceDetail(ctx context.Context, service schema.ServiceInterface, res chan<- any) error {
	client := service.(*collector.Services).ECS

	clusterArns, err := listClusters(ctx, client)
	if err != nil {
		log.CtxLogger(ctx).Error("failed to list ecs clusters for services", zap.Error(err))
		return err
	}

	var wg sync.WaitGroup
	tasks := make(chan struct{ cluster, service string }, 100) // Buffered channel

	// Start workers
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range tasks {
				describedServices, err := describeServices(ctx, client, task.cluster, []string{task.service})
				if err != nil {
					log.CtxLogger(ctx).Warn("failed to describe ecs service", zap.String("cluster", task.cluster), zap.String("service", task.service), zap.Error(err))
					continue
				}
				if len(describedServices) > 0 {
					res <- ServiceDetail{Service: describedServices[0]}
				}
			}
		}()
	}

	// Collect all service ARNs from all clusters
	for _, clusterArn := range clusterArns {
		serviceArns, err := listServices(ctx, client, clusterArn)
		if err != nil {
			log.CtxLogger(ctx).Warn("failed to list ecs services for cluster", zap.String("cluster", clusterArn), zap.Error(err))
			continue
		}
		for _, serviceArn := range serviceArns {
			tasks <- struct{ cluster, service string }{cluster: clusterArn, service: serviceArn}
		}
	}
	close(tasks)

	wg.Wait()

	return nil
}

// listServices retrieves all ECS service ARNs in a cluster.
func listServices(ctx context.Context, c *ecs.Client, clusterArn string) ([]string, error) {
	var serviceArns []string
	paginator := ecs.NewListServicesPaginator(c, &ecs.ListServicesInput{Cluster: &clusterArn})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		serviceArns = append(serviceArns, page.ServiceArns...)
	}
	return serviceArns, nil
}

// describeServices retrieves the details for a list of services.
func describeServices(ctx context.Context, c *ecs.Client, clusterArn string, serviceArns []string) ([]types.Service, error) {
	output, err := c.DescribeServices(ctx, &ecs.DescribeServicesInput{Cluster: &clusterArn, Services: serviceArns, Include: []types.ServiceField{types.ServiceFieldTags}})
	if err != nil {
		return nil, err
	}
	return output.Services, nil
}
