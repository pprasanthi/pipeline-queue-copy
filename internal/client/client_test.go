package client_test

import (
	"testing"

	"github.com/plouc/go-gitlab-client/gitlab"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gitlab.com/fenrirunbound/pipeline-queue/internal/client"
)

type ClientTestSuite struct {
	suite.Suite
}

type MockGitlabClient struct {
	mock.Mock
}

func (m *MockGitlabClient) ProjectPipelines(projectId string, opts *gitlab.PipelinesOptions) (*gitlab.PipelineCollection, *gitlab.ResponseMeta, error) {
	args := m.Called(projectId)
	var collection *gitlab.PipelineCollection = args.Get(0).(*gitlab.PipelineCollection)

	return collection, nil, args.Error(2)
}

func (suite *ClientTestSuite) TestCreateClient() {
	assert := assert.New(suite.T())
	fakeClient := new(MockGitlabClient)
	c, err := client.New(fakeClient, "", "")

	assert.Nil(err, "Got error: %v", err)
	assert.NotNil(c, "Client is Nil")
}

func (suite *ClientTestSuite) TestListPipelines() {
	assert := assert.New(suite.T())
	fakeClient := new(MockGitlabClient)
	c, err := client.New(fakeClient, "", "")

	testPipeline := &gitlab.Pipeline{
		Id: 1027,
	}
	testCollection := &gitlab.PipelineCollection{
		Items: []*gitlab.Pipeline{
			testPipeline,
		},
	}
	fakeClient.On("ProjectPipelines", "103").Return(testCollection, nil, nil)

	pipelines, err := c.ListRunningPipelines("103")

	assert.Nil(err, "Got error: %v", err)
	assert.Len(pipelines.Items, 1, "Pipeline count was not accurate: %v", pipelines)
	assert.Equal(pipelines.Items[0], testPipeline, "Unexpected pipeline: %v", pipelines.Items[0])
}

func (suite *ClientTestSuite) TestIndexOfPipeline() {
	assert := assert.New(suite.T())
	fakeClient := new(MockGitlabClient)
	c, _ := client.New(fakeClient, "", "")

	testCollection := &gitlab.PipelineCollection{
		Items: []*gitlab.Pipeline{
			&gitlab.Pipeline{
				Id: 0,
			},
			&gitlab.Pipeline{
				Id: 333,
			},
			&gitlab.Pipeline{
				Id: 1027,
			},
			&gitlab.Pipeline{
				Id: 1234,
			},
		},
	}

	targetIndex, err := c.IndexOfPipeline(testCollection, "1027")

	assert.Nil(err, "Got error: %v", err)
	assert.Equal(targetIndex, 2, "Expected index 2, got %v", targetIndex)
}

func (suite *ClientTestSuite) TestDetermineIfFirst() {
	assert := assert.New(suite.T())
	fakeClient := new(MockGitlabClient)
	c, _ := client.New(fakeClient, "", "")

	testCollection := &gitlab.PipelineCollection{
		Items: []*gitlab.Pipeline{
			&gitlab.Pipeline{
				Id: 1027,
			},
			&gitlab.Pipeline{
				Id: 1234,
			},
		},
	}
	fakeClient.On("ProjectPipelines", "987").Return(testCollection, nil, nil)

	isFirst, err := c.DetermineIfFirst("987", "1027")

	assert.Nil(err, "Got error: %v", err)
	assert.True(isFirst, "Expected to be first, but was: %v", isFirst)
}

func TestClienTestSuite(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}
