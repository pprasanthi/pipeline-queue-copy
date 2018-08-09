package client

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/plouc/go-gitlab-client/gitlab"
)

// GitLabClient - Simplified interface for a GitLab client to wrap
type GitLabClient interface {
	ProjectPipelines(string, *gitlab.PipelinesOptions) (*gitlab.PipelineCollection, *gitlab.ResponseMeta, error)
	ProjectPipeline(string, string) (*gitlab.PipelineWithDetails, *gitlab.ResponseMeta, error)
}

// Client - Wrapper struct for a GitLab client
type Client struct {
	Client GitLabClient
}

// ListRunningPipelines - Query the GitLab API for a Project's list of Pipelines, sorted
//						  in ascending order (oldest to newest)
func (c *Client) ListRunningPipelines(projectID string) ([]*gitlab.PipelineWithDetails, error) {
	options := &gitlab.PipelinesOptions{
		Scope:      "running",
		Status:     "running",
		YamlErrors: false,
		SortOptions: gitlab.SortOptions{
			OrderBy: "id",
			Sort:    gitlab.SortDirectionAsc,
		},
	}

	pipelines, _, err := c.Client.ProjectPipelines(projectID, options)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return nil, err
	}

	result := make([]*gitlab.PipelineWithDetails, len(pipelines.Items))

	for index, pipeline := range pipelines.Items {
		strID := strconv.Itoa(pipeline.Id)
		details, _, err := c.Client.ProjectPipeline(projectID, strID)
		if err != nil {
			fmt.Printf("Error: %v", err)
			return nil, err
		}

		result[index] = details
	}

	return result, nil
}

// IndexOfPipeline - Given a sorted list of Pipelines, determine if the given Pipeline ID is
//					 first in the list (0th index).
func (c *Client) IndexOfPipeline(pipelines []*gitlab.PipelineWithDetails, pipelineID string) (int, error) {
	targetID, err := strconv.Atoi(pipelineID)
	if err != nil {
		return -1, err
	}

	for i, pipeline := range pipelines {
		if pipeline.Id == targetID {
			return i, nil
		}
	}

	return -1, fmt.Errorf("PipelineID %v not found in collection", pipelineID)
}

// SortByUpdated - Given a slice of PipelineWithDetails, sort them by the updated_at attribute
//				   Sorts the slice in-place (e.g., modifies the slice passed in)
func (c *Client) SortByUpdated(pipelines []*gitlab.PipelineWithDetails) []*gitlab.PipelineWithDetails {
	sort.Slice(pipelines, func(i, j int) bool {
		timeI, _ := time.Parse(time.RFC3339, pipelines[i].UpdatedAt)
		timeJ, _ := time.Parse(time.RFC3339, pipelines[j].UpdatedAt)

		return timeI.Before(timeJ)
	})

	return pipelines
}

// DetermineIfFirst - Checks to see if the Job's Pipeline is the oldest running pipeline
//                    for the Project.
func (c *Client) DetermineIfFirst(projectID string, pipelineID string) (bool, error) {
	pipelines, err := c.ListRunningPipelines(projectID)
	if err != nil {
		return false, err
	}

	c.SortByUpdated(pipelines)

	position, err := c.IndexOfPipeline(pipelines, pipelineID)
	if err != nil {
		return false, err
	}

	return position == 0, err
}

// New - Factory method for creating a new GitLab client wrapper
func New(desiredClient GitLabClient, hostname string, token string) (*Client, error) {
	var gitlabClient GitLabClient

	if desiredClient != nil {
		gitlabClient = desiredClient
	} else {
		gitlabClient = gitlab.NewGitlab(hostname, "/api/v4", token)
	}

	return &Client{
		Client: gitlabClient,
	}, nil

}
