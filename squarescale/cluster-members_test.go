package squarescale_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/squarescale/squarescale-cli/squarescale"
)

func TestClusterMembers(t *testing.T) {
	// GetClusterMembers
	t.Run("Nominal case on GetClusterMembers", nominalCaseOnGetClusterMembers)
	t.Run("Test project not found on GetClusterMembers", UnknownProjectOnGetClusterMembers)

	// GetClusterMemberInfo
	t.Run("Nominal case on GetClusterMemberInfo", nominalCaseOnGetClusterMemberInfo)
	t.Run("Test cluster member not found on GetClusterMemberInfo", UnknownClusterMemberOnGetClusterMemberInfo)
	t.Run("Test project not found on GetClusterMemberInfo", UnknownProjectOnGetClusterMemberInfo)

	// ConfigClusterMember
	t.Run("Nominal case on ConfigClusterMember", nominalCaseOnConfigClusterMember)
	t.Run("Test project not found on ConfigClusterMember", UnknownProjectOnConfigClusterMember)
}

func nominalCaseOnGetClusterMembers(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkPath(t, "/projects/"+projectName+"/cluster_members", r.URL.Path)
		checkAuthorization(t, r.Header.Get("Authorization"), token)

		resBody := `
		[
			{
				"id": 22,
				"name": "cm1",
				"public_ip": "127.0.0.1",
				"private_ip": "192.168.0.1",
				"nomad_status": "ready",
				"scheduling_group": "squarescale"
			},
			{
				"id": 23,
				"name": "cm2",
				"public_ip": "127.0.0.2",
				"private_ip": "192.168.0.2",
				"nomad_status": "ready",
				"scheduling_group": "sogilis"
			}
		]
		`

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	clusterMembers, err := cli.GetClusterMembers(projectName)

	// then
	var expectedInt int
	var expectedString string

	if err != nil {
		t.Fatalf("Expect no error, got `%s`", err)
	}

	expectedInt = 2
	if len(clusterMembers) != expectedInt {
		t.Fatalf("Expect cluster members to contain %d elements, but got actually %d", expectedInt, len(clusterMembers))
	}

	expectedInt = 22
	if clusterMembers[0].ID != expectedInt {
		t.Errorf("Expect clusterMember.ID `%d`, got `%d`", expectedInt, clusterMembers[0].ID)
	}

	expectedInt = 23
	if clusterMembers[1].ID != expectedInt {
		t.Errorf("Expect clusterMember.ID `%d`, got `%d`", expectedInt, clusterMembers[1].ID)
	}

	expectedString = "cm1"
	if clusterMembers[0].Name != expectedString {
		t.Errorf("Expect clusterMember.Name `%s`, got `%s`", expectedString, clusterMembers[0].Name)
	}

	expectedString = "cm2"
	if clusterMembers[1].Name != expectedString {
		t.Errorf("Expect clusterMember.Name `%s`, got `%s`", expectedString, clusterMembers[1].Name)
	}

	expectedString = "127.0.0.1"
	if clusterMembers[0].PublicIP != expectedString {
		t.Errorf("Expect clusterMember.PublicIP `%s`, got `%s`", expectedString, clusterMembers[0].PublicIP)
	}

	expectedString = "127.0.0.2"
	if clusterMembers[1].PublicIP != expectedString {
		t.Errorf("Expect clusterMember.PublicIP `%s`, got `%s`", expectedString, clusterMembers[1].PublicIP)
	}

	expectedString = "192.168.0.1"
	if clusterMembers[0].PrivateIP != expectedString {
		t.Errorf("Expect clusterMember.PrivateIP `%s`, got `%s`", expectedString, clusterMembers[0].PrivateIP)
	}

	expectedString = "192.168.0.2"
	if clusterMembers[1].PrivateIP != expectedString {
		t.Errorf("Expect clusterMember.PrivateIP `%s`, got `%s`", expectedString, clusterMembers[1].PrivateIP)
	}

	expectedString = "ready"
	if clusterMembers[0].Status != expectedString {
		t.Errorf("Expect clusterMember.Status `%s`, got `%s`", expectedString, clusterMembers[0].Status)
	}

	expectedString = "ready"
	if clusterMembers[1].Status != expectedString {
		t.Errorf("Expect clusterMember.Status `%s`, got `%s`", expectedString, clusterMembers[1].Status)
	}

	expectedString = "squarescale"
	if clusterMembers[0].SchedulingGroup != expectedString {
		t.Errorf("Expect clusterMember.SchedulingGroup `%s`, got `%s`", expectedString, clusterMembers[0].SchedulingGroup)
	}

	expectedString = "sogilis"
	if clusterMembers[1].SchedulingGroup != expectedString {
		t.Errorf("Expect clusterMember.SchedulingGroup `%s`, got `%s`", expectedString, clusterMembers[1].SchedulingGroup)
	}
}

func UnknownProjectOnGetClusterMembers(t *testing.T) {
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
	_, err := cli.GetClusterMembers(projectName)

	expectedError := "Project 'unknown-project' does not exist"
	if err == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, err)
	}
}

func nominalCaseOnGetClusterMemberInfo(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"
	clusterMemberName := "cm1"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkPath(t, "/projects/"+projectName+"/cluster_members", r.URL.Path)
		checkAuthorization(t, r.Header.Get("Authorization"), token)

		resBody := `
		[
			{
				"id": 22,
				"name": "cm1",
				"public_ip": "127.0.0.1",
				"private_ip": "192.168.0.1",
				"nomad_status": "ready",
				"scheduling_group": "squarescale"
			},
			{
				"id": 23,
				"name": "cm2",
				"public_ip": "127.0.0.2",
				"private_ip": "192.168.0.2",
				"nomad_status": "ready",
				"scheduling_group": "sogilis"
			}
		]
		`

		w.Header().Set("Content-Type", "application/json")

		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	theClusterMember, err := cli.GetClusterMemberInfo(projectName, clusterMemberName)

	// then
	var expectedInt int
	var expectedString string

	if err != nil {
		t.Fatalf("Expect no error, got `%s`", err)
	}

	expectedInt = 22
	if theClusterMember.ID != expectedInt {
		t.Errorf("Expect clusterMember.ID `%d`, got `%d`", expectedInt, theClusterMember.ID)
	}

	expectedString = "cm1"
	if theClusterMember.Name != expectedString {
		t.Errorf("Expect clusterMember.Name `%s`, got `%s`", expectedString, theClusterMember.Name)
	}
}

func UnknownClusterMemberOnGetClusterMemberInfo(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resBody := `
		[
			{
				"id": 22,
				"name": "cm1",
				"public_ip": "127.0.0.1",
				"private_ip": "192.168.0.1",
				"nomad_status": "ready",
				"scheduling_group": "squarescale"
			},
			{
				"id": 23,
				"name": "cm2",
				"public_ip": "127.0.0.2",
				"private_ip": "192.168.0.2",
				"nomad_status": "ready",
				"scheduling_group": "sogilis"
			}
		]
		`

		w.Header().Set("Content-Type", "application/json")

		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	_, err := cli.GetClusterMemberInfo(projectName, "missing-cluster-member")

	// then
	expectedError := "Cluster node 'missing-cluster-member' not found for project 'my-project'"
	if err == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, err)
	}
}

func UnknownProjectOnGetClusterMemberInfo(t *testing.T) {
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
	_, err := cli.GetClusterMemberInfo("unknown-project", "missing-cluster-member")

	// then
	expectedError := "Project 'unknown-project' does not exist"
	if err == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, err)
	}
}

func nominalCaseOnConfigClusterMember(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"

	clusterMember := squarescale.ClusterMember{
		ID:              1,
		Name:            "cm1",
		PublicIP:        "10.0.0.1",
		PrivateIP:       "192.168.0.1",
		SchedulingGroup: "squarescale",
	}
	newClusterMember := squarescale.ClusterMember{
		SchedulingGroup: "sogilis",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkPath(t, fmt.Sprintf("/projects/%s/cluster_members/%d", projectName, clusterMember.ID), r.URL.Path)
		checkAuthorization(t, r.Header.Get("Authorization"), token)
		resBody := `null`
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(resBody))
	}))
	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	err := cli.ConfigClusterMember(projectName, clusterMember, newClusterMember)

	// then
	if err != nil {
		t.Fatalf("Expect no error, got `%s`", err)
	}
}

func UnknownProjectOnConfigClusterMember(t *testing.T) {
	// given
	token := "some-token"

	clusterMember := squarescale.ClusterMember{
		ID:              1,
		Name:            "cm1",
		PublicIP:        "10.0.0.1",
		PrivateIP:       "192.168.0.1",
		SchedulingGroup: "squarescale",
	}
	newClusterMember := squarescale.ClusterMember{
		ID:              1,
		Name:            "cm1",
		PublicIP:        "10.0.0.1",
		PrivateIP:       "192.168.0.1",
		SchedulingGroup: "sogilis",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resBody := `{"error": "Project 'project-not-found' not found"}`
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(404)
		w.Write([]byte(resBody))
	}))
	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	err := cli.ConfigClusterMember("project-not-found", clusterMember, newClusterMember)

	// then
	expectedError := "Project 'project-not-found' not found"

	if err == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, err)
	}
}
