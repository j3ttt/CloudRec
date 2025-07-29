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

// GetTaskResource returns a Task Resource
func GetTaskResource() schema.Resource {
	return schema.Resource{
		ResourceType:       collector.ECSTask,
		ResourceTypeName:   "ECS Task",
		ResourceGroupType:  constant.CONTAINER,
		Desc:               `https://docs.aws.amazon.com/AmazonECS/latest/APIReference/API_DescribeTasks.html`,
		ResourceDetailFunc: GetTaskDetail,
		RowField: schema.RowField{
			ResourceId:   "$.Task.TaskArn",
			ResourceName: "$.Task.TaskArn", // No friendly name
		},
		Dimension: schema.Regional,
	}
}

// TaskDetail aggregates all information for a single ECS task.
type TaskDetail struct {
	Task types.Task
}

// GetTaskDetail fetches the details for all ECS tasks in a region.
func GetTaskDetail(ctx context.Context, service schema.ServiceInterface, res chan<- any) error {
	client := service.(*collector.Services).ECS

	clusterArns, err := listClusters(ctx, client)
	if err != nil {
		log.CtxLogger(ctx).Error("failed to list ecs clusters for tasks", zap.Error(err))
		return err
	}

	var wg sync.WaitGroup
	tasksChan := make(chan struct{ cluster, task string }, 100) // Buffered channel

	// Start workers
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range tasksChan {
				describedTasks, err := describeTasks(ctx, client, task.cluster, []string{task.task})
				if err != nil {
					log.CtxLogger(ctx).Warn("failed to describe ecs task", zap.String("cluster", task.cluster), zap.String("task", task.task), zap.Error(err))
					continue
				}
				if len(describedTasks) > 0 {
					res <- TaskDetail{Task: describedTasks[0]}
				}
			}
		}()
	}

	// Collect all task ARNs from all clusters
	for _, clusterArn := range clusterArns {
		taskArns, err := listTasks(ctx, client, clusterArn)
		if err != nil {
			log.CtxLogger(ctx).Warn("failed to list ecs tasks for cluster", zap.String("cluster", clusterArn), zap.Error(err))
			continue
		}
		for _, taskArn := range taskArns {
			tasksChan <- struct{ cluster, task string }{cluster: clusterArn, task: taskArn}
		}
	}
	close(tasksChan)

	wg.Wait()

	return nil
}

// listTasks retrieves all ECS task ARNs in a cluster.
func listTasks(ctx context.Context, c *ecs.Client, clusterArn string) ([]string, error) {
	var taskArns []string
	paginator := ecs.NewListTasksPaginator(c, &ecs.ListTasksInput{Cluster: &clusterArn})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		taskArns = append(taskArns, page.TaskArns...)
	}
	return taskArns, nil
}

// describeTasks retrieves the details for a list of tasks.
func describeTasks(ctx context.Context, c *ecs.Client, clusterArn string, taskArns []string) ([]types.Task, error) {
	output, err := c.DescribeTasks(ctx, &ecs.DescribeTasksInput{Cluster: &clusterArn, Tasks: taskArns, Include: []types.TaskField{types.TaskFieldTags}})
	if err != nil {
		return nil, err
	}
	return output.Tasks, nil
}
