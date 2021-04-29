package squarescale_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/squarescale/squarescale-cli/squarescale"
)

func TestExternalNodes(t *testing.T) {
	// GetExternalNodes
	t.Run("Nominal case on GetExternalNodes", nominalCaseOnGetExternalNodes)
	t.Run("Test project not found on GetExternalNodes", UnknownProjectOnGetExternalNodes)

	// GetExternalNodeInfo
	t.Run("Nominal case on GetExternalNodeInfo", nominalCaseOnGetExternalNodeInfo)
	t.Run("Test external node not found on GetExternalNodeInfo", UnknownSchedulingGroupOnGetExternalNodeInfo)
	t.Run("Test project not found on GetExternalNodeInfo", UnknownProjectOnGetExternalNodeInfo)

	// AddExternalNode
	t.Run("Nominal case on AddExternalNode", nominalCaseOnAddExternalNode)
	t.Run("Test project not found on AddExternalNode", UnknownProjectOnAddExternalNode)

	// WaitStatefulNode
	t.Run("Nominal case on WaitStatefulNode", nominalCaseOnWaitExternalNode)

	// Error cases
	// t.Run("Test HTTP client error on stateful-nodes methods (get, add, delete)", ClientHTTPErrorOnSchedulingGroupMethods)
	// t.Run("Test internal server error on stateful-nodes methods (get, add, delete)", InternalServerErrorOnSchedulingGroupMethods)
	// t.Run("Test badly JSON on stateful-nodes methods (get, add)", CantUnmarshalOnSchedulingGroupMethods)
}

func nominalCaseOnGetExternalNodes(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkPath(t, "/projects/"+projectName+"/external_nodes", r.URL.Path)
		checkAuthorization(t, r.Header.Get("Authorization"), token)

		resBody := `
		[
			{
				"id": 22,
				"name": "externalNode1",
				"public_ip": "127.0.0.1"
			},
			{
				"id": 23,
				"name": "externalNode2",
				"public_ip": "127.0.0.1"
			}
		]
		`

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	externalNodes, err := cli.GetExternalNodes(projectName)

	// then
	var expectedInt int
	var expectedString string

	if err != nil {
		t.Fatalf("Expect no error, got `%s`", err)
	}

	expectedInt = 2
	if len(externalNodes) != expectedInt {
		t.Fatalf("Expect scheduling groups to contain %d elements, but got actually %d", expectedInt, len(externalNodes))
	}

	expectedInt = 22
	if externalNodes[0].ID != expectedInt {
		t.Errorf("Expect externalNode.ID `%d`, got `%d`", expectedInt, externalNodes[0].ID)
	}

	expectedInt = 23
	if externalNodes[1].ID != expectedInt {
		t.Errorf("Expect externalNode.ID `%d`, got `%d`", expectedInt, externalNodes[1].ID)
	}

	expectedString = "externalNode1"
	if externalNodes[0].Name != expectedString {
		t.Errorf("Expect externalNode.Name `%s`, got `%s`", expectedString, externalNodes[0].Name)
	}

	expectedString = "externalNode2"
	if externalNodes[1].Name != expectedString {
		t.Errorf("Expect externalNode.Name `%s`, got `%s`", expectedString, externalNodes[1].Name)
	}
}

func UnknownProjectOnGetExternalNodes(t *testing.T) {
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
	_, err := cli.GetExternalNodes(projectName)

	expectedError := "Project 'unknown-project' does not exist"
	if err == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, err)
	}
}

func nominalCaseOnGetExternalNodeInfo(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"
	externalNodeName := "external-node-1"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkPath(t, "/projects/"+projectName+"/external_nodes", r.URL.Path)
		checkAuthorization(t, r.Header.Get("Authorization"), token)

		resBody := `
		[
			{
				"id": 23,
				"name": "external-node-1",
				"public_ip": "127.0.0.1"
			},
			{
				"id": 32,
				"name": "external-node-2",
				"public_ip": "127.0.0.1"
			}
		]
		`

		w.Header().Set("Content-Type", "application/json")

		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	theSchedulingGroup, err := cli.GetExternalNodeInfo(projectName, externalNodeName)

	// then
	var expectedInt int
	var expectedString string

	if err != nil {
		t.Fatalf("Expect no error, got `%s`", err)
	}

	expectedInt = 23
	if theSchedulingGroup.ID != expectedInt {
		t.Errorf("Expect externalNode.ID `%d`, got `%d`", expectedInt, theSchedulingGroup.ID)
	}

	expectedString = "external-node-1"
	if theSchedulingGroup.Name != expectedString {
		t.Errorf("Expect externalNode.Name `%s`, got `%s`", expectedString, theSchedulingGroup.Name)
	}

	expectedString = "127.0.0.1"
	if theSchedulingGroup.PublicIP != expectedString {
		t.Errorf("Expect externalNode.PublicIP `%s`, got `%s`", expectedString, theSchedulingGroup.Name)
	}
}

func UnknownSchedulingGroupOnGetExternalNodeInfo(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resBody := `
		[
			{
				"id": 23,
				"name": "external-node-1",
				"public_ip": "127.0.0.1"
			},
			{
				"id": 32,
				"name": "external-node-2",
				"public_ip": "127.0.0.1"
			}
		]
		`

		w.Header().Set("Content-Type", "application/json")

		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	_, err := cli.GetExternalNodeInfo(projectName, "missing-external-node")

	// then
	expectedError := "External node 'missing-external-node' not found for project 'my-project'"
	if err == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, err)
	}
}

func UnknownProjectOnGetExternalNodeInfo(t *testing.T) {
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
	_, err := cli.GetExternalNodeInfo("unknown-project", "missing-external-node")

	// then
	expectedError := "Project 'unknown-project' does not exist"
	if err == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, err)
	}
}

func nominalCaseOnAddExternalNode(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"
	externalNodeName := "external-node-1"
	public_ip := "127.0.0.1"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkPath(t, "/projects/"+projectName+"/external_nodes", r.URL.Path)
		checkAuthorization(t, r.Header.Get("Authorization"), token)

		resBody := `
		{
			"id": 23,
			"name": "external-node-1",
			"public_ip": "127.0.0.1"
		}
		`

		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(201)
		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	newSchedulingGroup, err := cli.AddExternalNode(projectName, externalNodeName, public_ip)

	// then
	var expectedInt int
	var expectedString string

	if err != nil {
		t.Fatalf("Expect no error, got `%s`", err)
	}

	expectedInt = 23
	if newSchedulingGroup.ID != expectedInt {
		t.Errorf("Expect externalNode.ID `%d`, got `%d`", expectedInt, newSchedulingGroup.ID)
	}

	expectedString = "external-node-1"
	if newSchedulingGroup.Name != expectedString {
		t.Errorf("Expect externalNode.Name `%s`, got `%s`", expectedString, newSchedulingGroup.Name)
	}

	expectedString = "127.0.0.1"
	if newSchedulingGroup.PublicIP != expectedString {
		t.Errorf("Expect externalNode.PublicIP `%s`, got `%s`", expectedString, newSchedulingGroup.Name)
	}
}

func UnknownProjectOnAddExternalNode(t *testing.T) {
	// given
	token := "some-token"
	projectName := "unknown-project"
	externalNodeName := "external-node-1"
	public_ip := "127.0.0.1"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resBody := `{"error":"No project found for config name: unknown-project"}`

		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(404)
		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	_, err := cli.AddExternalNode(projectName, externalNodeName, public_ip)

	// then
	expectedError := "Project 'unknown-project' does not exist"
	if err == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, err)
	}
}

func nominalCaseOnWaitExternalNode(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"
	nodeName := "my-node"
	httptestCount := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkPath(t, "/projects/"+projectName+"/external_nodes", r.URL.Path)
		checkAuthorization(t, r.Header.Get("Authorization"), token)

		var statefulNodeStatus string

		httptestCount++

		if httptestCount < 2 {
			statefulNodeStatus = "not_provisionned"
		} else {
			statefulNodeStatus = "provisionned"
		}
		resBody := fmt.Sprintf(`
		[
			{
				"id": 23,
				"name": "my-node",
				"public_ip": "127.0.0.1",
				"status": "%s"
			}
		]
		`, statefulNodeStatus)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)

		w.Write([]byte(resBody))

	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	_, err := cli.WaitExternalNode(projectName, nodeName, 0)

	// then
	if err != nil {
		t.Fatalf("Expect no error, got `%s`", err)
	}
}
