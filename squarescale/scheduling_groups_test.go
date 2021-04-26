package squarescale_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/squarescale/squarescale-cli/squarescale"
)

func TestSchedulingGroups(t *testing.T) {
	// GetSchedulingGroups
	t.Run("Nominal case on GetSchedulingGroups", nominalCaseOnGetSchedulingGroups)
	t.Run("Test project not found on GetSchedulingGroups", UnknownProjectOnGetSchedulingGroups)

	// GetSchedulingGroupInfo
	t.Run("Nominal case on GetSchedulingGroupInfo", nominalCaseOnGetSchedulingGroupInfo)
	t.Run("Test scheduling group not found on GetSchedulingGroupInfo", UnknownSchedulingGroupOnGetSchedulingGroupInfo)
	t.Run("Test project not found on GetSchedulingGroupInfo", UnknownProjectOnGetSchedulingGroupInfo)

	// AddSchedulingGroup
	t.Run("Nominal case on AddSchedulingGroup", nominalCaseOnAddSchedulingGroup)
	t.Run("Test project not found on AddSchedulingGroup", UnknownProjectOnAddSchedulingGroup)
	t.Run("Test to create a duplicate on AddSchedulingGroup", DuplicateNodeOnAddSchedulingGroup)

	// DeleteSchedulingGroup
	t.Run("Nominal case on DeleteSchedulingGroup", nominalCaseOnDeleteSchedulingGroup)
}

func nominalCaseOnGetSchedulingGroups(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkPath(t, "/projects/"+projectName+"/scheduling_groups", r.URL.Path)
		checkAuthorization(t, r.Header.Get("Authorization"), token)

		resBody := `
		[
			{
				"id": 22,
				"name": "group1"
			},
			{
				"id": 23,
				"name": "group2"
			}
		]
		`

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	schedulingGroups, err := cli.GetSchedulingGroups(projectName)

	// then
	var expectedInt int
	var expectedString string

	if err != nil {
		t.Fatalf("Expect no error, got `%s`", err)
	}

	expectedInt = 2
	if len(schedulingGroups) != expectedInt {
		t.Fatalf("Expect scheduling groups to contain %d elements, but got actually %d", expectedInt, len(schedulingGroups))
	}

	expectedInt = 22
	if schedulingGroups[0].ID != expectedInt {
		t.Errorf("Expect schedulingGroup.ID `%d`, got `%d`", expectedInt, schedulingGroups[0].ID)
	}

	expectedInt = 23
	if schedulingGroups[1].ID != expectedInt {
		t.Errorf("Expect schedulingGroup.ID `%d`, got `%d`", expectedInt, schedulingGroups[1].ID)
	}

	expectedString = "group1"
	if schedulingGroups[0].Name != expectedString {
		t.Errorf("Expect schedulingGroup.Name `%s`, got `%s`", expectedString, schedulingGroups[0].Name)
	}

	expectedString = "group2"
	if schedulingGroups[1].Name != expectedString {
		t.Errorf("Expect schedulingGroup.Name `%s`, got `%s`", expectedString, schedulingGroups[1].Name)
	}
}

func UnknownProjectOnGetSchedulingGroups(t *testing.T) {
	// given
	token := "some-token"
	projectName := "unknown-project"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resBody := `{"error":"No project found for config name: unknown-project"}`

		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(404)
		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	_, err := cli.GetSchedulingGroups(projectName)

	expectedError := "Project 'unknown-project' does not exist"
	if err == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, err)
	}
}

func nominalCaseOnGetSchedulingGroupInfo(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"
	schedulingGroupName := "group1"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkPath(t, "/projects/"+projectName+"/scheduling_groups", r.URL.Path)
		checkAuthorization(t, r.Header.Get("Authorization"), token)

		resBody := `
		[
			{
				"id": 23,
				"name": "group1"
			},
			{
				"id": 32,
				"name": "group2"
			}
		]
		`

		w.Header().Set("Content-Type", "application/json")

		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	theStatefulNode, err := cli.GetSchedulingGroupInfo(projectName, schedulingGroupName)

	// then
	var expectedInt int
	var expectedString string

	if err != nil {
		t.Fatalf("Expect no error, got `%s`", err)
	}

	expectedInt = 23
	if theStatefulNode.ID != expectedInt {
		t.Errorf("Expect schedulingGroup.ID `%d`, got `%d`", expectedInt, theStatefulNode.ID)
	}

	expectedString = "group1"
	if theStatefulNode.Name != expectedString {
		t.Errorf("Expect schedulingGroup.Name `%s`, got `%s`", expectedString, theStatefulNode.Name)
	}
}

func UnknownSchedulingGroupOnGetSchedulingGroupInfo(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resBody := `
		[
			{
				"id": 23,
				"name": "group1"
			},
			{
				"id": 32,
				"name": "group2"
			}
		]
		`

		w.Header().Set("Content-Type", "application/json")

		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	_, err := cli.GetSchedulingGroupInfo(projectName, "missing-scheduling-group")

	// then
	expectedError := "Scheduling group 'missing-scheduling-group' not found for project 'my-project'"
	if err == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, err)
	}
}

func UnknownProjectOnGetSchedulingGroupInfo(t *testing.T) {
	// given
	token := "some-token"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resBody := `{"error":"No project found for config name: unknown-project"}`

		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(404)
		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	_, err := cli.GetStatefulNodeInfo("unknown-project", "missing-scheduling-group")

	// then
	expectedError := "Project 'unknown-project' does not exist"
	if err == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, err)
	}
}

func nominalCaseOnAddSchedulingGroup(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"
	nodeName := "group1"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkPath(t, "/projects/"+projectName+"/scheduling_groups", r.URL.Path)
		checkAuthorization(t, r.Header.Get("Authorization"), token)

		resBody := `
		{
			"id": 23,
			"name": "group1"
		}
		`

		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(201)
		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	newSchedulingGroup, err := cli.AddSchedulingGroup(projectName, nodeName)

	// then
	var expectedInt int
	var expectedString string

	if err != nil {
		t.Fatalf("Expect no error, got `%s`", err)
	}

	expectedInt = 23
	if newSchedulingGroup.ID != expectedInt {
		t.Errorf("Expect schedulingGroup.ID `%d`, got `%d`", expectedInt, newSchedulingGroup.ID)
	}

	expectedString = "group1"
	if newSchedulingGroup.Name != expectedString {
		t.Errorf("Expect schedulingGroup.Name `%s`, got `%s`", expectedString, newSchedulingGroup.Name)
	}
}

func UnknownProjectOnAddSchedulingGroup(t *testing.T) {
	// given
	token := "some-token"
	projectName := "unknown-project"
	nodeName := "group1"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resBody := `{"error":"No project found for config name: unknown-project"}`

		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(404)
		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	_, err := cli.AddSchedulingGroup(projectName, nodeName)

	// then
	expectedError := "Project 'unknown-project' does not exist"
	if err == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, err)
	}
}

func DuplicateNodeOnAddSchedulingGroup(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"
	nodeName := "group1"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resBody := `{"error":"PG::UniqueViolation: ERROR:  duplicate key value violates unique constraint \"index_scheduling_groups_on_cluster_id_and_name\"\nDETAIL:  Key (cluster_id, name)=(6163, nodeb) already exists.\n: INSERT INTO \"scheduling_groups\" (\"name\") VALUES ($1) RETURNING \"id\""}`

		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(409)
		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	_, err := cli.AddSchedulingGroup(projectName, nodeName)

	// then
	expectedError := "Scheduling group already exist on project 'my-project': group1"
	if err == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, err)
	}
}

func nominalCaseOnDeleteSchedulingGroup(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"
	schedulingGroupName := "group1"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkPath(t, "/projects/"+projectName+"/scheduling_groups/"+schedulingGroupName, r.URL.Path)
		checkAuthorization(t, r.Header.Get("Authorization"), token)

		resBody := `null`

		w.Header().Set("Content-Type", "application/json")

		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	err := cli.DeleteSchedulingGroup(projectName, schedulingGroupName)

	// then
	if err != nil {
		t.Fatalf("Expect no error, got `%s`", err)
	}
}
