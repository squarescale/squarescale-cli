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

	// AddStatefullNode
	t.Run("Nominal case on AddStatefullNode", nominalCaseOnAddStatefullNode)

	t.Run("Test project not found on AddStatefullNode", UnknownProjectOnAddStatefullNode)
	t.Run("Test to create a duplicate on AddStatefullNode", DuplicateNodeOnAddStatefullNode)

	// Error cases
	t.Run("Test HTTP client error on statefull nodes methods (get, add)", ClientHTTPErrorOnStatefullNodeMethods)
	t.Run("Test internal server error on statefull nodes methods (get, add)", InternalServerErrorOnStatefullNodeMethods)
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

func nominalCaseOnAddStatefullNode(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"
	nodeName := "node1a"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkPath(t, "/projects/"+projectName+"/statefull_nodes", r.URL.Path)

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

		if (r.Header.Get("Authorization")) != "bearer some-token" {
			t.Fatalf("Wrong token! Expected %s, got %s", "bearer some-token", r.Header.Get("Authorization"))
		}

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

	// then
	if errOnGet == nil {
		t.Errorf("Error is not raised on GetStatefullNodes")
	}

	if errOnAdd == nil {
		t.Errorf("Error is not raised on AddStatefullNodes")
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
