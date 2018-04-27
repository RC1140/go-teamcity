package teamcity_test

import (
	"testing"

	teamcity "github.com/cvbarros/go-teamcity-sdk/pkg/teamcity"
	"github.com/stretchr/testify/assert"
)

func TestCreateBuildTypeForProject(t *testing.T) {
	client := setup()
	newProject := getTestProjectData(testBuildTypeProjectId)
	_, err := client.Projects.Create(newProject)

	if err != nil {
		t.Fatalf("Failed to create project for buildType: %s", err)
	}
	newBuildType := getTestBuildTypeData(testBuildTypeProjectId)

	actual, err := client.BuildTypes.Create(testBuildTypeProjectId, newBuildType)

	if err != nil {
		t.Fatalf("Failed to CreateBuildType: %s", err)
	}

	if actual == nil {
		t.Fatalf("CreateBuildType did not return a valid instance")
	}

	cleanUpProject(t, client, testBuildTypeProjectId)

	assert.NotEmpty(t, actual.ID)
	assert.Equal(t, newBuildType.ProjectID, actual.ProjectID)
	assert.Equal(t, newBuildType.Name, actual.Name)
}

func TestAttachVcsRoot(t *testing.T) {
	client := setup()
	newProject := getTestProjectData(testBuildTypeProjectId)

	if _, err := client.Projects.Create(newProject); err != nil {
		t.Fatalf("Failed to create project for buildType: %s", err)
	}

	newBuildType := getTestBuildTypeData(testBuildTypeProjectId)

	createdBuildType, err := client.BuildTypes.Create(testBuildTypeProjectId, newBuildType)
	if err != nil {
		t.Fatalf("Failed to CreateBuildType: %s", err)
	}

	newVcsRoot := getTestVcsRootData(testBuildTypeProjectId)

	vcsRootCreated, err := client.VcsRoots.Create(testBuildTypeProjectId, newVcsRoot)

	if err != nil {
		t.Fatalf("Failed to create vcs root: %s", err)
	}

	err = client.BuildTypes.AttachVcsRoot(createdBuildType.ID, vcsRootCreated)
	if err != nil {
		t.Fatalf("Failed to attach vcsRoot '%s' to buildType '%s': %s", createdBuildType.ID, vcsRootCreated.ID, err)
	}

	actual, err := client.BuildTypes.GetById(createdBuildType.ID)
	if err != nil {
		t.Fatalf("Failed to get buildType '%s' for asserting: %s", createdBuildType.ID, err)
	}

	assert.NotEqualf(t, actual.VcsRootEntries.Count, 0, "Expected VcsRootEntries to contain at least one element")
	vcsEntries := idMapVcsRootEntries(actual.VcsRootEntries)
	assert.Containsf(t, vcsEntries, vcsRootCreated.ID, "Expected VcsRootEntries to contain the VcsRoot with id = %s, but it did not", vcsRootCreated.ID)

	cleanUpProject(t, client, testBuildTypeProjectId)
}

func idMapVcsRootEntries(v *teamcity.VcsRootEntries) map[string]string {
	out := make(map[string]string)
	for _, item := range v.Items {
		out[item.VcsRoot.ID] = item.Id
	}

	return out
}

func getTestBuildTypeData(projectId string) *teamcity.BuildType {

	return &teamcity.BuildType{
		Name:        "Pull Request",
		Description: "Inspection Build",
		ProjectID:   projectId,
	}
}

func cleanUpBuildType(t *testing.T, c *teamcity.Client, id string) {
	if err := c.BuildTypes.Delete(id); err != nil {
		t.Errorf("Unable to delete build type with id = '%s', err: %s", id, err)
		return
	}

	deleted, err := c.BuildTypes.GetById(id)

	if deleted != nil {
		t.Errorf("Build type not deleted during cleanup.")
		return
	}

	if err == nil {
		t.Errorf("Expected 404 Not Found error when getting Deleted Build Type, but no error returned.")
	}
}
