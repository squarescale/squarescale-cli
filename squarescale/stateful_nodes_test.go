package squarescale_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/squarescale/squarescale-cli/squarescale"
)

func TestStatefulNodes(t *testing.T) {

	// GetStatefulNodes
	t.Run("Nominal case on GetStatefulNodes", nominalCaseOnGetStatefulNodes)

	t.Run("Test project not found on GetStatefulNodes", UnknownProjectOnGetStatefulNodes)

	// GetStatefulNodeInfo
	t.Run("Nominal case on GetStatefulNodeInfo", nominalCaseOnGetStatefulNodeInfo)

	t.Run("Test stateful-node not found on GetStatefulNodeInfo", UnknownStatefulNodeOnGetStatefulNodeInfo)
	t.Run("Test project not found on GetStatefulNodeInfo", UnknownProjectOnGetStatefulNodeInfo)

	// AddStatefulNode
	t.Run("Nominal case on AddStatefulNode", nominalCaseOnAddStatefulNode)

	t.Run("Test project not found on AddStatefulNode", UnknownProjectOnAddStatefulNode)
	t.Run("Test to create a duplicate on AddStatefulNode", DuplicateNodeOnAddStatefulNode)

	// DeleteStatefulNode
	t.Run("Nominal case on DeleteStatefulNode", nominalCaseOnDeleteStatefulNode)

	t.Run("Test project not found on DeleteStatefulNode", UnknownProjectOnDeleteStatefulNode)
	t.Run("Test to delete a missing stateful-node on DeleteStatefulNode", UnknownStatefulNodeOnDeleteStatefulNode)
	t.Run("Test to delete a volume when deploy is in progress on DeleteStatefulNode", DeployInProgressOnDeleteStatefulNode)

	// WaitStatefulNode
	t.Run("Nominal case on WaitStatefulNode", nominalCaseOnWaitStatefulNode)

	// BindVolumeOnStatefulNode
	t.Run("Nominal case on BindVolumeOnStatefulNode", nominalCaseOnBindVolumeOnStatefulNode)

	t.Run("Test project not found on BindVolumeOnStatefulNode", UnknownProjectOnBindVolumeOnStatefulNode)
	t.Run("Test volume not found on BindVolumeOnStatefulNode", UnknownVolumeOnBindVolumeOnStatefulNode)
	t.Run("Test rebind volume already is on BindVolumeOnStatefulNode", RebindVolumeOnBindVolumeOnStatefulNode)

	// Error cases
	t.Run("Test HTTP client error on stateful-nodes methods (get, add, delete)", ClientHTTPErrorOnStatefulNodeMethods)
	t.Run("Test internal server error on stateful-nodes methods (get, add, delete)", InternalServerErrorOnStatefulNodeMethods)
	t.Run("Test badly JSON on stateful-nodes methods (get, add)", CantUnmarshalOnStatefulNodeMethods)
}

func nominalCaseOnGetStatefulNodes(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkPath(t, "/projects/"+projectName+"/statefull_nodes", r.URL.Path)
		checkAuthorization(t, r.Header.Get("Authorization"), token)

		resBody := `
		[
			{
				"id": 23,
				"name": "nodeb",
				"node_type": "t2.micro",
				"zone": "eu-west-1b",
				"status": "provisionned"
			},
			{
				"id": 22,
				"name": "test1",
				"node_type": "t2.micro",
				"zone": "eu-west-1a",
				"status": "not_provisionned"
			}
		]
		`

		w.Header().Set("Content-Type", "application/json")

		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	statefulNodes, err := cli.GetStatefulNodes(projectName)

	// then
	var expectedInt int
	var expectedString string

	if err != nil {
		t.Fatalf("Expect no error, got `%s`", err)
	}

	expectedInt = 2
	if len(statefulNodes) != expectedInt {
		t.Fatalf("Expect stateful-nodes to contain %d elements, but got actually %d", expectedInt, len(statefulNodes))
	}

	expectedInt = 23
	if statefulNodes[0].ID != expectedInt {
		t.Errorf("Expect statefulNode.ID `%d`, got `%d`", expectedInt, statefulNodes[0].ID)
	}

	expectedString = "nodeb"
	if statefulNodes[0].Name != expectedString {
		t.Errorf("Expect statefulNode.Name `%s`, got `%s`", expectedString, statefulNodes[0].Name)
	}

	expectedString = "t2.micro"
	if statefulNodes[0].NodeType != expectedString {
		t.Errorf("Expect statefulNode.NodeType `%s`, got `%s`", expectedString, statefulNodes[0].NodeType)
	}

	expectedString = "eu-west-1b"
	if statefulNodes[0].Zone != expectedString {
		t.Errorf("Expect statefulNode.Zone `%s`, got `%s`", expectedString, statefulNodes[0].Zone)
	}

	expectedString = "provisionned"
	if statefulNodes[0].Status != expectedString {
		t.Errorf("Expect statefulNode.Status `%s`, got `%s`", expectedString, statefulNodes[0].Status)
	}

	expectedInt = 22
	if statefulNodes[1].ID != expectedInt {
		t.Errorf("Expect statefulNode.ID `%d`, got `%d`", expectedInt, statefulNodes[1].ID)
	}

	expectedString = "test1"
	if statefulNodes[1].Name != expectedString {
		t.Errorf("Expect statefulNode.Name `%s`, got `%s`", expectedString, statefulNodes[1].Name)
	}

	expectedString = "t2.micro"
	if statefulNodes[1].NodeType != expectedString {
		t.Errorf("Expect statefulNode.NodeType `%s`, got `%s`", expectedString, statefulNodes[1].NodeType)
	}

	expectedString = "eu-west-1a"
	if statefulNodes[1].Zone != expectedString {
		t.Errorf("Expect statefulNode.Zone `%s`, got `%s`", expectedString, statefulNodes[1].Zone)
	}

	expectedString = "not_provisionned"
	if statefulNodes[1].Status != expectedString {
		t.Errorf("Expect statefulNode.Status `%s`, got `%s`", expectedString, statefulNodes[1].Status)
	}
}

func UnknownProjectOnGetStatefulNodes(t *testing.T) {
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
	_, err := cli.GetStatefulNodes(projectName)

	expectedError := "Project 'unknown-project' does not exist"
	if err == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, err)
	}
}

func nominalCaseOnGetStatefulNodeInfo(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"
	nodeName := "node1a"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkPath(t, "/projects/"+projectName+"/statefull_nodes", r.URL.Path)
		checkAuthorization(t, r.Header.Get("Authorization"), token)

		resBody := `
		[
			{
				"id": 23,
				"name": "node1a",
				"node_type": "t2.micro",
				"zone": "eu-west-1c",
				"status": "not_provisionned"
			},
			{
				"id": 32,
				"name": "node1b",
				"node_type": "t2.micro",
				"zone": "eu-west-1b",
				"status": "provisionned"
			}
		]
		`

		w.Header().Set("Content-Type", "application/json")

		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	theStatefulNode, err := cli.GetStatefulNodeInfo(projectName, nodeName)

	// then
	var expectedInt int
	var expectedString string

	if err != nil {
		t.Fatalf("Expect no error, got `%s`", err)
	}

	expectedInt = 23
	if theStatefulNode.ID != expectedInt {
		t.Errorf("Expect statefulNode.ID `%d`, got `%d`", expectedInt, theStatefulNode.ID)
	}

	expectedString = "node1a"
	if theStatefulNode.Name != expectedString {
		t.Errorf("Expect statefulNode.Name `%s`, got `%s`", expectedString, theStatefulNode.Name)
	}

	expectedString = "t2.micro"
	if theStatefulNode.NodeType != expectedString {
		t.Errorf("Expect statefulNode.NodeType `%s`, got `%s`", expectedString, theStatefulNode.NodeType)
	}

	expectedString = "eu-west-1c"
	if theStatefulNode.Zone != expectedString {
		t.Errorf("Expect statefulNode.Zone `%s`, got `%s`", expectedString, theStatefulNode.Zone)
	}

	expectedString = "not_provisionned"
	if theStatefulNode.Status != expectedString {
		t.Errorf("Expect statefulNode.Status `%s`, got `%s`", expectedString, theStatefulNode.Status)
	}
}

func UnknownStatefulNodeOnGetStatefulNodeInfo(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resBody := `
		[
			{
				"id": 23,
				"name": "node1a",
				"node_type": "t2.micro",
				"zone": "eu-west-1c",
				"status": "not_provisionned"
			},
			{
				"id": 32,
				"name": "node1b",
				"node_type": "t2.micro",
				"zone": "eu-west-1b",
				"status": "provisionned"
			}
		]
		`

		w.Header().Set("Content-Type", "application/json")

		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	_, err := cli.GetStatefulNodeInfo(projectName, "missing-stateful-node")

	// then
	expectedError := "Stateful-node 'missing-stateful-node' not found for project 'my-project'"
	if err == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, err)
	}
}

func UnknownProjectOnGetStatefulNodeInfo(t *testing.T) {
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
	_, err := cli.GetStatefulNodeInfo("unknown-project", "missing-stateful-node")

	// then
	expectedError := "Project 'unknown-project' does not exist"
	if err == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, err)
	}
}

func nominalCaseOnAddStatefulNode(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"
	nodeName := "node1a"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkPath(t, "/projects/"+projectName+"/statefull_nodes", r.URL.Path)
		checkAuthorization(t, r.Header.Get("Authorization"), token)

		resBody := `
		{
			"id": 23,
			"name": "node1a",
			"node_type": "t2.micro",
			"zone": "eu-west-1c",
			"status": "not_provisionned"
		}
		`

		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(201)
		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	newStatefulNode, err := cli.AddStatefulNode(projectName, nodeName, "t2.micro", "eu-west-1a")

	// then
	var expectedInt int
	var expectedString string

	if err != nil {
		t.Fatalf("Expect no error, got `%s`", err)
	}

	expectedInt = 23
	if newStatefulNode.ID != expectedInt {
		t.Errorf("Expect statefulNode.ID `%d`, got `%d`", expectedInt, newStatefulNode.ID)
	}

	expectedString = "node1a"
	if newStatefulNode.Name != expectedString {
		t.Errorf("Expect statefulNode.Name `%s`, got `%s`", expectedString, newStatefulNode.Name)
	}

	expectedString = "t2.micro"
	if newStatefulNode.NodeType != expectedString {
		t.Errorf("Expect statefulNode.NodeType `%s`, got `%s`", expectedString, newStatefulNode.NodeType)
	}

	expectedString = "eu-west-1c"
	if newStatefulNode.Zone != expectedString {
		t.Errorf("Expect statefulNode.Zone `%s`, got `%s`", expectedString, newStatefulNode.Zone)
	}

	expectedString = "not_provisionned"
	if newStatefulNode.Status != expectedString {
		t.Errorf("Expect statefulNode.Status `%s`, got `%s`", expectedString, newStatefulNode.Status)
	}
}

func UnknownProjectOnAddStatefulNode(t *testing.T) {
	// given
	token := "some-token"
	projectName := "unknown-project"
	nodeName := "node1a"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resBody := `{"error":"No project found for config name: unknown-project"}`

		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(404)
		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	_, err := cli.AddStatefulNode(projectName, nodeName, "t2.micro", "eu-west-1a")

	// then
	expectedError := "Project 'unknown-project' does not exist"
	if err == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, err)
	}
}

func DuplicateNodeOnAddStatefulNode(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"
	nodeName := "node1a"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resBody := `{"error":"PG::UniqueViolation: ERROR:  duplicate key value violates unique constraint \"index_statefull_nodes_on_cluster_id_and_name\"\nDETAIL:  Key (cluster_id, name)=(6163, nodeb) already exists.\n: INSERT INTO \"statefull_nodes\" (\"name\", \"node_type\", \"zone\", \"cluster_id\", \"status\") VALUES ($1, $2, $3, $4, $5) RETURNING \"id\""}`

		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(409)
		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	_, err := cli.AddStatefulNode(projectName, nodeName, "t2.micro", "eu-west-1a")

	// then
	expectedError := "Stateful-node already exist on project 'my-project': node1a"
	if err == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, err)
	}
}

func nominalCaseOnDeleteStatefulNode(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"
	nodeName := "node1a"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkPath(t, "/projects/"+projectName+"/statefull_nodes/"+nodeName, r.URL.Path)
		checkAuthorization(t, r.Header.Get("Authorization"), token)

		resBody := `null`

		w.Header().Set("Content-Type", "application/json")

		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	err := cli.DeleteStatefulNode(projectName, nodeName)

	// then
	if err != nil {
		t.Fatalf("Expect no error, got `%s`", err)
	}
}

func UnknownProjectOnDeleteStatefulNode(t *testing.T) {
	// given
	token := "some-token"
	projectName := "not-a-project"
	nodeName := "node1a"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resBody := `{"error":"No project found for config name: not-a-project"}`

		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(404)
		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	err := cli.DeleteStatefulNode(projectName, nodeName)

	// then
	expectedError := "Project 'not-a-project' does not exist"
	if err == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, err)
	}
}

func UnknownStatefulNodeOnDeleteStatefulNode(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"
	nodeName := "missing-node"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resBody := `{"error":"Couldn't find StatefulNode with [WHERE \"statefull_nodes\".\"cluster_id\" = $1 AND \"statefull_nodes\".\"name\" = $2]"}`

		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(404)
		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	err := cli.DeleteStatefulNode(projectName, nodeName)

	// then
	expectedError := "Stateful-node 'missing-node' does not exist"
	if err == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, err)
	}
}

func DeployInProgressOnDeleteStatefulNode(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"
	nodeName := "my-node"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(400)
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	err := cli.DeleteStatefulNode(projectName, nodeName)

	// then
	expectedError := "Deploy probably in progress"
	if err == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, err)
	}
}

func nominalCaseOnWaitStatefulNode(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"
	nodeName := "my-node"
	httptestCount := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkPath(t, "/projects/"+projectName+"/statefull_nodes", r.URL.Path)
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
				"node_type": "t2.micro",
				"zone": "eu-west-1b",
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
	_, err := cli.WaitStatefulNode(projectName, nodeName, 0)

	// then
	if err != nil {
		t.Fatalf("Expect no error, got `%s`", err)
	}
}

func nominalCaseOnBindVolumeOnStatefulNode(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"
	nodeName := "my-node"
	volumeName := "vol1c"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkPath(t, "/projects/"+projectName+"/statefull_nodes/"+nodeName, r.URL.Path)
		checkAuthorization(t, r.Header.Get("Authorization"), token)

		resBody := `null`

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)

		w.Write([]byte(resBody))

	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	err := cli.BindVolumeOnStatefulNode(projectName, nodeName, volumeName)

	// then
	if err != nil {
		t.Fatalf("Expect no error, got `%s`", err)
	}
}

func UnknownProjectOnBindVolumeOnStatefulNode(t *testing.T) {
	// given
	token := "some-token"
	nodeName := "my-node"
	volumeName := "vol1c"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resBody := `{"error":"not_found"}`

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(404)

		w.Write([]byte(resBody))

	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	err := cli.BindVolumeOnStatefulNode("unknown-project", nodeName, volumeName)

	// then
	expectedError := "Project 'unknown-project' does not exist"
	if err == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, err)
	}
}

func UnknownVolumeOnBindVolumeOnStatefulNode(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"
	nodeName := "my-node"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resBody := `{"error":"Couldn't find Volume with [WHERE \"volumes\".\"cluster_id\" = $1 AND \"volumes\".\"name\" = $2]"}`

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(404)

		w.Write([]byte(resBody))

	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	err := cli.BindVolumeOnStatefulNode(projectName, nodeName, "unknown-volume")

	// then
	expectedError := "Volume 'unknown-volume' does not exist"
	if err == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, err)
	}
}

func RebindVolumeOnBindVolumeOnStatefulNode(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"
	nodeName := "my-node"
	volumeName := "vol1c"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resBody := `{"error":"vol02c already bound with node2c"}`

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)

		w.Write([]byte(resBody))

	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	err := cli.BindVolumeOnStatefulNode(projectName, nodeName, volumeName)

	// then
	expectedError := "Volume vol1c already bound with my-node"
	if err == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, err)
	}
}

func ClientHTTPErrorOnStatefulNodeMethods(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"
	nodeName := "node1a"
	volumeName := "vol1c"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	_, errOnGet := cli.GetStatefulNodes(projectName)
	_, errOnAdd := cli.AddStatefulNode(projectName, nodeName, "t2.micro", "eu-west-1a")
	errOnDelete := cli.DeleteStatefulNode(projectName, nodeName)
	_, errOnGetInfo := cli.GetStatefulNodeInfo(projectName, nodeName)
	_, errOnWait := cli.WaitStatefulNode(projectName, nodeName, 0)
	errOnBindVolume := cli.BindVolumeOnStatefulNode(projectName, nodeName, volumeName)

	// then
	if errOnGet == nil {
		t.Errorf("Error is not raised on GetStatefulNodes")
	}

	if errOnAdd == nil {
		t.Errorf("Error is not raised on AddStatefulNodes")
	}

	if errOnDelete == nil {
		t.Errorf("Error is not raised on DeleteStatefulNodes")
	}

	if errOnGetInfo == nil {
		t.Errorf("Error is not raised on GetStatefulNodeInfo")
	}

	if errOnWait == nil {
		t.Errorf("Error is not raised on WaitStatefulNode")
	}

	if errOnBindVolume == nil {
		t.Errorf("Error is not raised on BindVolumeOnStatefulNode")
	}
}

func InternalServerErrorOnStatefulNodeMethods(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"
	nodeName := "node1a"
	volumeName := "vol1c"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	_, errOnGet := cli.GetStatefulNodes(projectName)
	_, errOnAdd := cli.AddStatefulNode(projectName, nodeName, "t2.micro", "eu-west-1a")
	errOnDelete := cli.DeleteStatefulNode(projectName, nodeName)
	_, errOnGetInfo := cli.GetStatefulNodeInfo(projectName, nodeName)
	_, errOnWait := cli.WaitStatefulNode(projectName, nodeName, 0)
	errOnBindVolume := cli.BindVolumeOnStatefulNode(projectName, nodeName, volumeName)

	// then
	expectedError := "An unexpected error occurred (code: 500)"
	if errOnGet == nil {
		t.Errorf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", errOnGet) != expectedError {
		t.Errorf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, errOnGet)
	}

	if errOnAdd == nil {
		t.Errorf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", errOnAdd) != expectedError {
		t.Errorf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, errOnAdd)
	}

	if errOnDelete == nil {
		t.Errorf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", errOnDelete) != expectedError {
		t.Errorf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, errOnDelete)
	}

	if errOnGetInfo == nil {
		t.Errorf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", errOnGetInfo) != expectedError {
		t.Errorf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, errOnGetInfo)
	}

	if errOnWait == nil {
		t.Errorf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", errOnWait) != expectedError {
		t.Errorf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, errOnWait)
	}

	if errOnBindVolume == nil {
		t.Errorf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", errOnBindVolume) != expectedError {
		t.Errorf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, errOnBindVolume)
	}
}

func CantUnmarshalOnStatefulNodeMethods(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resBody := `{]`

		w.Header().Set("Content-Type", "application/json")

		if r.Method == "POST" { // Manage Add
			w.WriteHeader(201)
		} else {
			w.WriteHeader(200)
		}

		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	_, errOnGet := cli.GetStatefulNodes(projectName)
	_, errOnAdd := cli.AddStatefulNode(projectName, "node-in-error", "t2.micro", "eu-west-1a")

	// then
	expectedError := "invalid character ']' looking for beginning of object key string"
	if errOnGet == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", errOnGet) != expectedError {
		t.Fatalf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, errOnGet)
	}

	if errOnAdd == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", errOnAdd) != expectedError {
		t.Fatalf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, errOnAdd)
	}
}
