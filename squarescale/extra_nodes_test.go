package squarescale_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/squarescale/squarescale-cli/squarescale"
)

func TestExtraNodes(t *testing.T) {

	// GetExtraNodes
	t.Run("Nominal case on GetExtraNodes", nominalCaseOnGetExtraNodes)

	t.Run("Test project not found on GetExtraNodes", UnknownProjectOnGetExtraNodes)

	// GetExtraNodeInfo
	t.Run("Nominal case on GetExtraNodeInfo", nominalCaseOnGetExtraNodeInfo)

	t.Run("Test extra-node not found on GetExtraNodeInfo", UnknownExtraNodeOnGetExtraNodeInfo)
	t.Run("Test project not found on GetExtraNodeInfo", UnknownProjectOnGetExtraNodeInfo)

	// AddExtraNode
	t.Run("Nominal case on AddExtraNode", nominalCaseOnAddExtraNode)

	t.Run("Test project not found on AddExtraNode", UnknownProjectOnAddExtraNode)
	t.Run("Test to create a duplicate on AddExtraNode", DuplicateNodeOnAddExtraNode)

	// DeleteExtraNode
	t.Run("Nominal case on DeleteExtraNode", nominalCaseOnDeleteExtraNode)

	t.Run("Test project not found on DeleteExtraNode", UnknownProjectOnDeleteExtraNode)
	t.Run("Test to delete a missing extra-node on DeleteExtraNode", UnknownExtraNodeOnDeleteExtraNode)
	t.Run("Test to delete a volume when deploy is in progress on DeleteExtraNode", DeployInProgressOnDeleteExtraNode)

	// WaitExtraNode
	t.Run("Nominal case on WaitExtraNode", nominalCaseOnWaitExtraNode)

	// BindVolumeOnExtraNode
	t.Run("Nominal case on BindVolumeOnExtraNode", nominalCaseOnBindVolumeOnExtraNode)

	t.Run("Test project not found on BindVolumeOnExtraNode", UnknownProjectOnBindVolumeOnExtraNode)
	t.Run("Test volume not found on BindVolumeOnExtraNode", UnknownVolumeOnBindVolumeOnExtraNode)
	t.Run("Test rebind volume already is on BindVolumeOnExtraNode", RebindVolumeOnBindVolumeOnExtraNode)

	// Error cases
	t.Run("Test HTTP client error on extra-nodes methods (get, add, delete)", ClientHTTPErrorOnExtraNodeMethods)
	t.Run("Test internal server error on extra-nodes methods (get, add, delete)", InternalServerErrorOnExtraNodeMethods)
	t.Run("Test badly JSON on extra-nodes methods (get, add)", CantUnmarshalOnExtraNodeMethods)
}

func nominalCaseOnGetExtraNodes(t *testing.T) {
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
	extraNodes, err := cli.GetExtraNodes(projectName)

	// then
	var expectedInt int
	var expectedString string

	if err != nil {
		t.Fatalf("Expect no error, got `%s`", err)
	}

	expectedInt = 2
	if len(extraNodes) != expectedInt {
		t.Fatalf("Expect extra-nodes to contain %d elements, but got actually %d", expectedInt, len(extraNodes))
	}

	expectedInt = 23
	if extraNodes[0].ID != expectedInt {
		t.Errorf("Expect extraNode.ID `%d`, got `%d`", expectedInt, extraNodes[0].ID)
	}

	expectedString = "nodeb"
	if extraNodes[0].Name != expectedString {
		t.Errorf("Expect extraNode.Name `%s`, got `%s`", expectedString, extraNodes[0].Name)
	}

	expectedString = "t2.micro"
	if extraNodes[0].NodeType != expectedString {
		t.Errorf("Expect extraNode.NodeType `%s`, got `%s`", expectedString, extraNodes[0].NodeType)
	}

	expectedString = "eu-west-1b"
	if extraNodes[0].Zone != expectedString {
		t.Errorf("Expect extraNode.Zone `%s`, got `%s`", expectedString, extraNodes[0].Zone)
	}

	expectedString = "provisionned"
	if extraNodes[0].Status != expectedString {
		t.Errorf("Expect extraNode.Status `%s`, got `%s`", expectedString, extraNodes[0].Status)
	}

	expectedInt = 22
	if extraNodes[1].ID != expectedInt {
		t.Errorf("Expect extraNode.ID `%d`, got `%d`", expectedInt, extraNodes[1].ID)
	}

	expectedString = "test1"
	if extraNodes[1].Name != expectedString {
		t.Errorf("Expect extraNode.Name `%s`, got `%s`", expectedString, extraNodes[1].Name)
	}

	expectedString = "t2.micro"
	if extraNodes[1].NodeType != expectedString {
		t.Errorf("Expect extraNode.NodeType `%s`, got `%s`", expectedString, extraNodes[1].NodeType)
	}

	expectedString = "eu-west-1a"
	if extraNodes[1].Zone != expectedString {
		t.Errorf("Expect extraNode.Zone `%s`, got `%s`", expectedString, extraNodes[1].Zone)
	}

	expectedString = "not_provisionned"
	if extraNodes[1].Status != expectedString {
		t.Errorf("Expect extraNode.Status `%s`, got `%s`", expectedString, extraNodes[1].Status)
	}
}

func UnknownProjectOnGetExtraNodes(t *testing.T) {
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
	_, err := cli.GetExtraNodes(projectName)

	expectedError := "Project 'unknown-project' does not exist"
	if err == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, err)
	}
}

func nominalCaseOnGetExtraNodeInfo(t *testing.T) {
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
	theExtraNode, err := cli.GetExtraNodeInfo(projectName, nodeName)

	// then
	var expectedInt int
	var expectedString string

	if err != nil {
		t.Fatalf("Expect no error, got `%s`", err)
	}

	expectedInt = 23
	if theExtraNode.ID != expectedInt {
		t.Errorf("Expect extraNode.ID `%d`, got `%d`", expectedInt, theExtraNode.ID)
	}

	expectedString = "node1a"
	if theExtraNode.Name != expectedString {
		t.Errorf("Expect extraNode.Name `%s`, got `%s`", expectedString, theExtraNode.Name)
	}

	expectedString = "t2.micro"
	if theExtraNode.NodeType != expectedString {
		t.Errorf("Expect extraNode.NodeType `%s`, got `%s`", expectedString, theExtraNode.NodeType)
	}

	expectedString = "eu-west-1c"
	if theExtraNode.Zone != expectedString {
		t.Errorf("Expect extraNode.Zone `%s`, got `%s`", expectedString, theExtraNode.Zone)
	}

	expectedString = "not_provisionned"
	if theExtraNode.Status != expectedString {
		t.Errorf("Expect extraNode.Status `%s`, got `%s`", expectedString, theExtraNode.Status)
	}
}

func UnknownExtraNodeOnGetExtraNodeInfo(t *testing.T) {
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
	_, err := cli.GetExtraNodeInfo(projectName, "missing-extra-node")

	// then
	expectedError := "extra-node 'missing-extra-node' not found for project 'my-project'"
	if err == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, err)
	}
}

func UnknownProjectOnGetExtraNodeInfo(t *testing.T) {
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
	_, err := cli.GetExtraNodeInfo("unknown-project", "missing-extra-node")

	// then
	expectedError := "Project 'unknown-project' does not exist"
	if err == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, err)
	}
}

func nominalCaseOnAddExtraNode(t *testing.T) {
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
	newExtraNode, err := cli.AddExtraNode(projectName, nodeName, "t2.micro", "eu-west-1a")

	// then
	var expectedInt int
	var expectedString string

	if err != nil {
		t.Fatalf("Expect no error, got `%s`", err)
	}

	expectedInt = 23
	if newExtraNode.ID != expectedInt {
		t.Errorf("Expect extraNode.ID `%d`, got `%d`", expectedInt, newExtraNode.ID)
	}

	expectedString = "node1a"
	if newExtraNode.Name != expectedString {
		t.Errorf("Expect extraNode.Name `%s`, got `%s`", expectedString, newExtraNode.Name)
	}

	expectedString = "t2.micro"
	if newExtraNode.NodeType != expectedString {
		t.Errorf("Expect extraNode.NodeType `%s`, got `%s`", expectedString, newExtraNode.NodeType)
	}

	expectedString = "eu-west-1c"
	if newExtraNode.Zone != expectedString {
		t.Errorf("Expect extraNode.Zone `%s`, got `%s`", expectedString, newExtraNode.Zone)
	}

	expectedString = "not_provisionned"
	if newExtraNode.Status != expectedString {
		t.Errorf("Expect extraNode.Status `%s`, got `%s`", expectedString, newExtraNode.Status)
	}
}

func UnknownProjectOnAddExtraNode(t *testing.T) {
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
	_, err := cli.AddExtraNode(projectName, nodeName, "t2.micro", "eu-west-1a")

	// then
	expectedError := "Project 'unknown-project' does not exist"
	if err == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, err)
	}
}

func DuplicateNodeOnAddExtraNode(t *testing.T) {
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
	_, err := cli.AddExtraNode(projectName, nodeName, "t2.micro", "eu-west-1a")

	// then
	expectedError := "extra-node already exist on project 'my-project': node1a"
	if err == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, err)
	}
}

func nominalCaseOnDeleteExtraNode(t *testing.T) {
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
	err := cli.DeleteExtraNode(projectName, nodeName)

	// then
	if err != nil {
		t.Fatalf("Expect no error, got `%s`", err)
	}
}

func UnknownProjectOnDeleteExtraNode(t *testing.T) {
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
	err := cli.DeleteExtraNode(projectName, nodeName)

	// then
	expectedError := "Project 'not-a-project' does not exist"
	if err == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, err)
	}
}

func UnknownExtraNodeOnDeleteExtraNode(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"
	nodeName := "missing-node"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resBody := `{"error":"Couldn't find ExtraNode with [WHERE \"statefull_nodes\".\"cluster_id\" = $1 AND \"statefull_nodes\".\"name\" = $2]"}`

		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(404)
		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	err := cli.DeleteExtraNode(projectName, nodeName)

	// then
	expectedError := "extra-node 'missing-node' does not exist"
	if err == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, err)
	}
}

func DeployInProgressOnDeleteExtraNode(t *testing.T) {
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
	err := cli.DeleteExtraNode(projectName, nodeName)

	// then
	expectedError := "Deploy probably in progress"
	if err == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, err)
	}
}

func nominalCaseOnWaitExtraNode(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"
	nodeName := "my-node"
	httptestCount := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkPath(t, "/projects/"+projectName+"/statefull_nodes", r.URL.Path)
		checkAuthorization(t, r.Header.Get("Authorization"), token)

		var extraNodeStatus string

		httptestCount++

		if httptestCount < 2 {
			extraNodeStatus = "not_provisionned"
		} else {
			extraNodeStatus = "provisionned"
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
		`, extraNodeStatus)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)

		w.Write([]byte(resBody))

	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	_, err := cli.WaitExtraNode(projectName, nodeName, 0)

	// then
	if err != nil {
		t.Fatalf("Expect no error, got `%s`", err)
	}
}

func nominalCaseOnBindVolumeOnExtraNode(t *testing.T) {
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
	err := cli.BindVolumeOnExtraNode(projectName, nodeName, volumeName)

	// then
	if err != nil {
		t.Fatalf("Expect no error, got `%s`", err)
	}
}

func UnknownProjectOnBindVolumeOnExtraNode(t *testing.T) {
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
	err := cli.BindVolumeOnExtraNode("unknown-project", nodeName, volumeName)

	// then
	expectedError := "Project 'unknown-project' does not exist"
	if err == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, err)
	}
}

func UnknownVolumeOnBindVolumeOnExtraNode(t *testing.T) {
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
	err := cli.BindVolumeOnExtraNode(projectName, nodeName, "unknown-volume")

	// then
	expectedError := "Volume 'unknown-volume' does not exist"
	if err == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, err)
	}
}

func RebindVolumeOnBindVolumeOnExtraNode(t *testing.T) {
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
	err := cli.BindVolumeOnExtraNode(projectName, nodeName, volumeName)

	// then
	expectedError := "Volume vol1c already bound with my-node"
	if err == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, err)
	}
}

func ClientHTTPErrorOnExtraNodeMethods(t *testing.T) {
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
	_, errOnGet := cli.GetExtraNodes(projectName)
	_, errOnAdd := cli.AddExtraNode(projectName, nodeName, "t2.micro", "eu-west-1a")
	errOnDelete := cli.DeleteExtraNode(projectName, nodeName)
	_, errOnGetInfo := cli.GetExtraNodeInfo(projectName, nodeName)
	_, errOnWait := cli.WaitExtraNode(projectName, nodeName, 0)
	errOnBindVolume := cli.BindVolumeOnExtraNode(projectName, nodeName, volumeName)

	// then
	if errOnGet == nil {
		t.Errorf("Error is not raised on GetExtraNodes")
	}

	if errOnAdd == nil {
		t.Errorf("Error is not raised on AddExtraNodes")
	}

	if errOnDelete == nil {
		t.Errorf("Error is not raised on DeleteExtraNodes")
	}

	if errOnGetInfo == nil {
		t.Errorf("Error is not raised on GetExtraNodeInfo")
	}

	if errOnWait == nil {
		t.Errorf("Error is not raised on WaitExtraNode")
	}

	if errOnBindVolume == nil {
		t.Errorf("Error is not raised on BindVolumeOnExtraNode")
	}
}

func InternalServerErrorOnExtraNodeMethods(t *testing.T) {
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
	_, errOnGet := cli.GetExtraNodes(projectName)
	_, errOnAdd := cli.AddExtraNode(projectName, nodeName, "t2.micro", "eu-west-1a")
	errOnDelete := cli.DeleteExtraNode(projectName, nodeName)
	_, errOnGetInfo := cli.GetExtraNodeInfo(projectName, nodeName)
	_, errOnWait := cli.WaitExtraNode(projectName, nodeName, 0)
	errOnBindVolume := cli.BindVolumeOnExtraNode(projectName, nodeName, volumeName)

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

func CantUnmarshalOnExtraNodeMethods(t *testing.T) {
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
	_, errOnGet := cli.GetExtraNodes(projectName)
	_, errOnAdd := cli.AddExtraNode(projectName, "node-in-error", "t2.micro", "eu-west-1a")

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
