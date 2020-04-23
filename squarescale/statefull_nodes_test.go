package squarescale_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/squarescale/squarescale-cli/squarescale"
)

func TestStatefullNodes(t *testing.T) {

	// GetStatefullNodes
	t.Run("Nominal case on GetStatefullNodes", nominalCaseOnGetStatefullNodes)

	t.Run("Test project not found on GetStatefullNodes", UnknownProjectOnGetStatefullNodes)

	// GetStatefullNodeInfo
	t.Run("Nominal case on GetStatefullNodeInfo", nominalCaseOnGetStatefullNodeInfo)

	t.Run("Test statefull node not found on GetStatefullNodeInfo", StatefullNodeNotFoundOnGetStatefullNodeInfo)
	t.Run("Test project not found on GetStatefullNodeInfo", ProjectNotFoundOnGetStatefullNodeInfo)

	// AddStatefullNode
	t.Run("Nominal case on AddStatefullNode", nominalCaseOnAddStatefullNode)

	t.Run("Test project not found on AddStatefullNode", UnknownProjectOnAddStatefullNode)
	t.Run("Test to create a duplicate on AddStatefullNode", DuplicateNodeOnAddStatefullNode)

	// DeleteStatefullNode
	t.Run("Nominal case on DeleteStatefullNode", nominalCaseOnDeleteStatefullNode)

	t.Run("Test project not found on DeleteStatefullNode", UnknownProjectOnDeleteStatefullNode)
	t.Run("Test to delete a missing statefull node on DeleteStatefullNode", NodeNotFoundOnDeleteStatefullNode)
	t.Run("Test to delete a volume when deploy is in progress on DeleteStatefullNode", DeployInProgressOnDeleteStatefullNode)

	// WaitStatefullNode
	t.Run("Nominal case on WaitStatefullNode", nominalCaseOnWaitStatefullNode)

	// Error cases
	t.Run("Test HTTP client error on statefull nodes methods (get, add, delete)", ClientHTTPErrorOnStatefullNodeMethods)
	t.Run("Test internal server error on statefull nodes methods (get, add, delete)", InternalServerErrorOnStatefullNodeMethods)
	t.Run("Test badly JSON on statefull nodes methods (get, add)", CantUnmarshalOnStatefullNodeMethods)
}

func nominalCaseOnGetStatefullNodes(t *testing.T) {
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
	statefullNodes, err := cli.GetStatefullNodes(projectName)

	// then
	var expectedInt int
	var expectedString string

	if err != nil {
		t.Fatalf("Expect no error, got `%s`", err)
	}

	expectedInt = 2
	if len(statefullNodes) != expectedInt {
		t.Fatalf("Expect statefull_nodes to contain %d elements, but got actually %d", expectedInt, len(statefullNodes))
	}

	expectedInt = 23
	if statefullNodes[0].ID != expectedInt {
		t.Errorf("Expect statefullNode.ID `%d`, got `%d`", expectedInt, statefullNodes[0].ID)
	}

	expectedString = "nodeb"
	if statefullNodes[0].Name != expectedString {
		t.Errorf("Expect statefullNode.Name `%s`, got `%s`", expectedString, statefullNodes[0].Name)
	}

	expectedString = "t2.micro"
	if statefullNodes[0].NodeType != expectedString {
		t.Errorf("Expect statefullNode.NodeType `%s`, got `%s`", expectedString, statefullNodes[0].NodeType)
	}

	expectedString = "eu-west-1b"
	if statefullNodes[0].Zone != expectedString {
		t.Errorf("Expect statefullNode.Zone `%s`, got `%s`", expectedString, statefullNodes[0].Zone)
	}

	expectedString = "provisionned"
	if statefullNodes[0].Status != expectedString {
		t.Errorf("Expect statefullNode.Status `%s`, got `%s`", expectedString, statefullNodes[0].Status)
	}

	expectedInt = 22
	if statefullNodes[1].ID != expectedInt {
		t.Errorf("Expect statefullNode.ID `%d`, got `%d`", expectedInt, statefullNodes[1].ID)
	}

	expectedString = "test1"
	if statefullNodes[1].Name != expectedString {
		t.Errorf("Expect statefullNode.Name `%s`, got `%s`", expectedString, statefullNodes[1].Name)
	}

	expectedString = "t2.micro"
	if statefullNodes[1].NodeType != expectedString {
		t.Errorf("Expect statefullNode.NodeType `%s`, got `%s`", expectedString, statefullNodes[1].NodeType)
	}

	expectedString = "eu-west-1a"
	if statefullNodes[1].Zone != expectedString {
		t.Errorf("Expect statefullNode.Zone `%s`, got `%s`", expectedString, statefullNodes[1].Zone)
	}

	expectedString = "not_provisionned"
	if statefullNodes[1].Status != expectedString {
		t.Errorf("Expect statefullNode.Status `%s`, got `%s`", expectedString, statefullNodes[1].Status)
	}
}

func UnknownProjectOnGetStatefullNodes(t *testing.T) {
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
	_, err := cli.GetStatefullNodes(projectName)

	expectedError := "Project 'unknown-project' does not exist"
	if err == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Expected error:\n`%s`\nGot:\n`%s`", expectedError, err)
	}
}

func nominalCaseOnGetStatefullNodeInfo(t *testing.T) {
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
	theStatefullNode, err := cli.GetStatefullNodeInfo(projectName, nodeName)

	// then
	var expectedInt int
	var expectedString string

	if err != nil {
		t.Fatalf("Expect no error, got `%s`", err)
	}

	expectedInt = 23
	if theStatefullNode.ID != expectedInt {
		t.Errorf("Expect statefullNode.ID `%d`, got `%d`", expectedInt, theStatefullNode.ID)
	}

	expectedString = "node1a"
	if theStatefullNode.Name != expectedString {
		t.Errorf("Expect statefullNode.Name `%s`, got `%s`", expectedString, theStatefullNode.Name)
	}

	expectedString = "t2.micro"
	if theStatefullNode.NodeType != expectedString {
		t.Errorf("Expect statefullNode.NodeType `%s`, got `%s`", expectedString, theStatefullNode.NodeType)
	}

	expectedString = "eu-west-1c"
	if theStatefullNode.Zone != expectedString {
		t.Errorf("Expect statefullNode.Zone `%s`, got `%s`", expectedString, theStatefullNode.Zone)
	}

	expectedString = "not_provisionned"
	if theStatefullNode.Status != expectedString {
		t.Errorf("Expect statefullNode.Status `%s`, got `%s`", expectedString, theStatefullNode.Status)
	}
}

func StatefullNodeNotFoundOnGetStatefullNodeInfo(t *testing.T) {
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
	_, err := cli.GetStatefullNodeInfo(projectName, "missing-statefull-node")

	// then
	expectedError := "Statefull node 'missing-statefull-node' not found for project 'my-project'"
	if err == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Expected error:\n`%s`\nGot:\n`%s`", expectedError, err)
	}
}

func ProjectNotFoundOnGetStatefullNodeInfo(t *testing.T) {
	// given
	token := "some-token"
	projectName := "unknown-project"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkPath(t, "/projects/"+projectName+"/statefull_nodes", r.URL.Path)
		checkAuthorization(t, r.Header.Get("Authorization"), token)

		resBody := `{"error":"No project found for config name: unknown-project"}`

		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(404)
		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	_, err := cli.GetStatefullNodeInfo("unknown-project", "missing-statefull-node")

	// then
	expectedError := "Project 'unknown-project' does not exist"
	if err == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Expected error:\n`%s`\nGot:\n`%s`", expectedError, err)
	}
}

func nominalCaseOnAddStatefullNode(t *testing.T) {
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
	newStatefullNode, err := cli.AddStatefullNode(projectName, nodeName, "t2.micro", "eu-west-1a")

	// then
	var expectedInt int
	var expectedString string

	if err != nil {
		t.Fatalf("Expect no error, got `%s`", err)
	}

	expectedInt = 23
	if newStatefullNode.ID != expectedInt {
		t.Errorf("Expect statefullNode.ID `%d`, got `%d`", expectedInt, newStatefullNode.ID)
	}

	expectedString = "node1a"
	if newStatefullNode.Name != expectedString {
		t.Errorf("Expect statefullNode.Name `%s`, got `%s`", expectedString, newStatefullNode.Name)
	}

	expectedString = "t2.micro"
	if newStatefullNode.NodeType != expectedString {
		t.Errorf("Expect statefullNode.NodeType `%s`, got `%s`", expectedString, newStatefullNode.NodeType)
	}

	expectedString = "eu-west-1c"
	if newStatefullNode.Zone != expectedString {
		t.Errorf("Expect statefullNode.Zone `%s`, got `%s`", expectedString, newStatefullNode.Zone)
	}

	expectedString = "not_provisionned"
	if newStatefullNode.Status != expectedString {
		t.Errorf("Expect statefullNode.Status `%s`, got `%s`", expectedString, newStatefullNode.Status)
	}
}

func UnknownProjectOnAddStatefullNode(t *testing.T) {
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
	_, err := cli.AddStatefullNode(projectName, nodeName, "t2.micro", "eu-west-1a")

	// then
	expectedError := "Project 'unknown-project' does not exist"
	if err == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Expected error:\n`%s`\nGot:\n`%s`", expectedError, err)
	}
}

func DuplicateNodeOnAddStatefullNode(t *testing.T) {
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
	_, err := cli.AddStatefullNode(projectName, nodeName, "t2.micro", "eu-west-1a")

	// then
	expectedError := "Statefull node already exist on project 'my-project': node1a"
	if err == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Expected error:\n`%s`\nGot:\n`%s`", expectedError, err)
	}
}

func nominalCaseOnDeleteStatefullNode(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"
	nodeName := "node1a"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkPath(t, "/projects/"+projectName+"/statefull_nodes"+nodeName, r.URL.Path)
		checkAuthorization(t, r.Header.Get("Authorization"), token)

		resBody := `null`

		w.Header().Set("Content-Type", "application/json")

		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	err := cli.DeleteStatefullNode(projectName, nodeName)

	// then
	if err != nil {
		t.Fatalf("Expect no error, got `%s`", err)
	}
}

func UnknownProjectOnDeleteStatefullNode(t *testing.T) {
	// given
	token := "some-token"
	projectName := "not-a-project"
	nodeName := "node1a"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkPath(t, "/projects/"+projectName+"/statefull_nodes"+nodeName, r.URL.Path)

		resBody := `{"error":"No project found for config name: not-a-project"}`

		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(404)
		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	err := cli.DeleteStatefullNode(projectName, nodeName)

	// then
	expectedError := "Project 'not-a-project' does not exist"
	if err == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Expected error:\n`%s`\nGot:\n`%s`", expectedError, err)
	}
}

func NodeNotFoundOnDeleteStatefullNode(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"
	nodeName := "missing-node"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkPath(t, "/projects/"+projectName+"/statefull_nodes"+nodeName, r.URL.Path)

		resBody := `{"error":"Couldn't find StatefullNode with [WHERE \"statefull_nodes\".\"cluster_id\" = $1 AND \"statefull_nodes\".\"name\" = $2]"}`

		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(404)
		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	err := cli.DeleteStatefullNode(projectName, nodeName)

	// then
	expectedError := "Statefull node 'missing-node' does not exist"
	if err == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Expected error:\n`%s`\nGot:\n`%s`", expectedError, err)
	}
}

func DeployInProgressOnDeleteStatefullNode(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"
	nodeName := "my-node"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkPath(t, "/projects/"+projectName+"/statefull_nodes"+nodeName, r.URL.Path)

		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(400)
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	err := cli.DeleteStatefullNode(projectName, nodeName)

	// then
	expectedError := "Deploy probably in progress"
	if err == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Expected error:\n`%s`\nGot:\n`%s`", expectedError, err)
	}
}

func nominalCaseOnWaitStatefullNode(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"
	nodeName := "my-node"
	httptestCount := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkPath(t, "/projects/"+projectName+"/statefull_nodes", r.URL.Path)
		checkAuthorization(t, r.Header.Get("Authorization"), token)

		var statefullNodeStatus string

		httptestCount++

		if httptestCount < 2 {
			statefullNodeStatus = "not_provisionned"
		} else {
			statefullNodeStatus = "provisionned"
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
		`, statefullNodeStatus)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)

		w.Write([]byte(resBody))

	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	_, err := cli.WaitStatefullNode(projectName, nodeName, 0)

	// then
	if err != nil {
		t.Fatalf("Expect no error, got `%s`", err)
	}
}

func ClientHTTPErrorOnStatefullNodeMethods(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"
	nodeName := "node1a"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	_, errOnGet := cli.GetStatefullNodes(projectName)
	_, errOnAdd := cli.AddStatefullNode(projectName, nodeName, "t2.micro", "eu-west-1a")
	errOnDelete := cli.DeleteStatefullNode(projectName, nodeName)
	_, errOnGetInfo := cli.GetStatefullNodeInfo(projectName, nodeName)
	_, errOnWait := cli.WaitStatefullNode(projectName, nodeName, 0)

	// then
	if errOnGet == nil {
		t.Errorf("Error is not raised on GetStatefullNodes")
	}

	if errOnAdd == nil {
		t.Errorf("Error is not raised on AddStatefullNodes")
	}

	if errOnDelete == nil {
		t.Errorf("Error is not raised on DeleteStatefullNodes")
	}

	if errOnGetInfo == nil {
		t.Errorf("Error is not raised on GetStatefullNodeInfo")
	}

	if errOnWait == nil {
		t.Errorf("Error is not raised on WaitStatefullNode")
	}
}

func InternalServerErrorOnStatefullNodeMethods(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"
	nodeName := "node1a"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	_, errOnGet := cli.GetStatefullNodes(projectName)
	_, errOnAdd := cli.AddStatefullNode(projectName, nodeName, "t2.micro", "eu-west-1a")
	errOnDelete := cli.DeleteStatefullNode(projectName, nodeName)
	_, errOnGetInfo := cli.GetStatefullNodeInfo(projectName, nodeName)
	_, errOnWait := cli.WaitStatefullNode(projectName, nodeName, 0)

	// then
	expectedError := "An unexpected error occurred (code: 500)"
	if errOnGet == nil {
		t.Errorf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", errOnGet) != expectedError {
		t.Errorf("Expected error:\n`%s`\nGot:\n`%s`", expectedError, errOnGet)
	}

	if errOnAdd == nil {
		t.Errorf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", errOnAdd) != expectedError {
		t.Errorf("Expected error:\n`%s`\nGot:\n`%s`", expectedError, errOnAdd)
	}

	if errOnDelete == nil {
		t.Errorf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", errOnDelete) != expectedError {
		t.Errorf("Expected error:\n`%s`\nGot:\n`%s`", expectedError, errOnDelete)
	}

	if errOnGetInfo == nil {
		t.Errorf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", errOnGetInfo) != expectedError {
		t.Errorf("Expected error:\n`%s`\nGot:\n`%s`", expectedError, errOnGetInfo)
	}

	if errOnWait == nil {
		t.Errorf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", errOnWait) != expectedError {
		t.Errorf("Expected error:\n`%s`\nGot:\n`%s`", expectedError, errOnWait)
	}
}

func CantUnmarshalOnStatefullNodeMethods(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resBody := `{]`

		w.Header().Set("Content-Type", "application/json")

		if (r.Header.Get("Authorization")) != "bearer some-token" {
			t.Fatalf("Wrong token! Expected `%s`, got `%s`", "bearer some-token", r.Header.Get("Authorization"))
		}

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
	_, errOnGet := cli.GetStatefullNodes(projectName)
	_, errOnAdd := cli.AddStatefullNode(projectName, "node-in-error", "t2.micro", "eu-west-1a")

	// then
	expectedError := "invalid character ']' looking for beginning of object key string"
	if errOnGet == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", errOnGet) != expectedError {
		t.Fatalf("Expected error:\n`%s`\nGot:\n`%s`", expectedError, errOnGet)
	}

	if errOnAdd == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", errOnAdd) != expectedError {
		t.Fatalf("Expected error:\n`%s`\nGot:\n`%s`", expectedError, errOnAdd)
	}
}
