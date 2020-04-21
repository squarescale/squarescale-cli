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
	t.Run("nominal get statefull nodes", nominalCaseForGetStatefullNodes)

	t.Run("test unknown project", UnknownProjectOnGetStatefullNodes)

	t.Run("test HTTP error", UnexpectedErrorOnGetStatefullNodes)
	t.Run("test Internal Server error", HTTPErrorOnGetStatefullNodes)
	t.Run("test badly formed JSON error", CannotUnmarshalOnGetStatefullNodes)

	// AddStatefullNode
	t.Run("nominal add statefull node", nominalCaseForAddStatefullNode)
	t.Run("test unknown project", UnknownProjectOnAddStatefullNode)
	t.Run("test to create a duplicate statefull node", DuplicateNodeOnAddStatefullNode)

	t.Run("test HTTP error", UnexpectedErrorOnAddStatefullNode)
	t.Run("test Internal Server error", HTTPErrorOnAddStatefullNode)
	t.Run("test badly formed JSON error", CannotUnmarshalOnAddStatefullNode)
}

func nominalCaseForGetStatefullNodes(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var path string = r.URL.Path

		if path != "/projects/"+projectName+"/statefull_nodes" {
			t.Fatalf("Wrong token! Expected %s, got %s", "/projects/my-project/statefull_nodes", path)
		}

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

		if (r.Header.Get("Authorization")) != "bearer some-token" {
			t.Fatalf("Wrong token! Expected %s, got %s", "bearer some-token", r.Header.Get("Authorization"))
		}

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

		var path string = r.URL.Path

		if path != "/projects/"+projectName+"/statefull_nodes" {
			t.Fatalf("Wrong token! Expected %s, got %s", "/projects/unknown-project/statefull_nodes", path)
		}

		resBody := `
		{"error":"No project found for config name: unknown-project"}
		`

		w.Header().Set("Content-Type", "application/json")

		if (r.Header.Get("Authorization")) != "bearer some-token" {
			t.Fatalf("Wrong token! Expected %s, got %s", "bearer some-token", r.Header.Get("Authorization"))
		}

		w.WriteHeader(404)
		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	_, err := cli.GetStatefullNodes("unknown-project")

	expectedError := "Project 'unknown-project' does not exist"
	if err == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Expected error:\n`%s`\nGot:\n`%s`", expectedError, err)
	}
}

func CannotUnmarshalOnGetStatefullNodes(t *testing.T) {
	// given
	token := "some-token"
	projectName := "unknown-project"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var path string = r.URL.Path

		if path != "/projects/"+projectName+"/statefull_nodes" {
			t.Fatalf("Wrong token! Expected %s, got %s", "/projects/unknown-project/statefull_nodes", path)
		}

		resBody := `{]`

		w.Header().Set("Content-Type", "application/json")

		if (r.Header.Get("Authorization")) != "bearer some-token" {
			t.Fatalf("Wrong token! Expected %s, got %s", "bearer some-token", r.Header.Get("Authorization"))
		}

		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	_, err := cli.GetStatefullNodes("unknown-project")

	// then
	expectedError := "invalid character ']' looking for beginning of object key string"
	if err == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Expected error:\n`%s`\nGot:\n`%s`", expectedError, err)
	}
}

func HTTPErrorOnGetStatefullNodes(t *testing.T) {
	// given
	token := "some-token"
	projectName := "bad-project"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var path string = r.URL.Path

		if path != "/projects/"+projectName+"/statefull_nodes" {
			t.Fatalf("Wrong token! Expected %s, got %s", "/projects/bad-project/statefull_nodes", path)
		}

		resBody := `
		{"error":"Hu ho, dummy error"}
		`

		w.Header().Set("Content-Type", "application/json")

		if (r.Header.Get("Authorization")) != "bearer some-token" {
			t.Fatalf("Wrong token! Expected %s, got %s", "bearer some-token", r.Header.Get("Authorization"))
		}

		w.WriteHeader(500)
		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	_, err := cli.GetStatefullNodes("bad-project")

	// then
	expectedError := `1 error occurred:

* error: Hu ho, dummy error`
	if err == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Expected error:\n`%s`\nGot:\n`%s`", expectedError, err)
	}
}

func UnexpectedErrorOnGetStatefullNodes(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	_, err := cli.GetStatefullNodes(projectName)

	// then
	if err == nil {
		t.Fatalf("Error is not raised")
	}
}

func nominalCaseForAddStatefullNode(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"
	nodeName := "node1a"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var path string = r.URL.Path

		if path != "/projects/"+projectName+"/statefull_nodes" {
			t.Fatalf("Wrong token! Expected %s, got %s", "/projects/my-project/statefull_nodes", path)
		}

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

		var path string = r.URL.Path

		if path != "/projects/"+projectName+"/statefull_nodes" {
			t.Fatalf("Wrong token! Expected %s, got %s", "/projects/unknown-project/statefull_nodes", path)
		}

		resBody := `
		{"error":"No project found for config name: unknown-project"}
		`

		w.Header().Set("Content-Type", "application/json")

		if (r.Header.Get("Authorization")) != "bearer some-token" {
			t.Fatalf("Wrong token! Expected %s, got %s", "bearer some-token", r.Header.Get("Authorization"))
		}

		w.WriteHeader(404)
		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	_, err := cli.AddStatefullNode("unknown-project", nodeName, "t2.micro", "eu-west-1a")

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

		var path string = r.URL.Path

		if path != "/projects/"+projectName+"/statefull_nodes" {
			t.Fatalf("Wrong token! Expected %s, got %s", "/projects/my-project/statefull_nodes", path)
		}

		resBody := `{"error":"PG::UniqueViolation: ERROR:  duplicate key value violates unique constraint \"index_statefull_nodes_on_cluster_id_and_name\"\nDETAIL:  Key (cluster_id, name)=(6163, nodeb) already exists.\n: INSERT INTO \"statefull_nodes\" (\"name\", \"node_type\", \"zone\", \"cluster_id\", \"status\") VALUES ($1, $2, $3, $4, $5) RETURNING \"id\""}`

		w.Header().Set("Content-Type", "application/json")

		if (r.Header.Get("Authorization")) != "bearer some-token" {
			t.Fatalf("Wrong token! Expected %s, got %s", "bearer some-token", r.Header.Get("Authorization"))
		}

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

func UnexpectedErrorOnAddStatefullNode(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"
	nodeName := "node1a"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	_, err := cli.AddStatefullNode(projectName, nodeName, "t2.micro", "eu-west-1a")

	// then
	if err == nil {
		t.Fatalf("Error is not raised")
	}
}

func HTTPErrorOnAddStatefullNode(t *testing.T) {
	// given
	token := "some-token"
	projectName := "bad-project"
	nodeName := "node1a"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var path string = r.URL.Path

		if path != "/projects/"+projectName+"/statefull_nodes" {
			t.Fatalf("Wrong token! Expected %s, got %s", "/projects/bad-project/statefull_nodes", path)
		}

		resBody := `
		{"error":"Hu ho, dummy error"}
		`

		w.Header().Set("Content-Type", "application/json")

		if (r.Header.Get("Authorization")) != "bearer some-token" {
			t.Fatalf("Wrong token! Expected %s, got %s", "bearer some-token", r.Header.Get("Authorization"))
		}

		w.WriteHeader(500)
		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	_, err := cli.AddStatefullNode(projectName, nodeName, "t2.micro", "eu-west-1a")

	// then
	expectedError := `1 error occurred:

* error: Hu ho, dummy error`
	if err == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Expected error:\n`%s`\nGot:\n`%s`", expectedError, err)
	}
}

func CannotUnmarshalOnAddStatefullNode(t *testing.T) {
	// given
	token := "some-token"
	projectName := "unknown-project"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var path string = r.URL.Path

		if path != "/projects/"+projectName+"/statefull_nodes" {
			t.Fatalf("Wrong token! Expected %s, got %s", "/projects/unknown-project/statefull_nodes", path)
		}

		resBody := `{]`

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
	_, err := cli.AddStatefullNode(projectName, "node-in-error", "t2.micro", "eu-west-1a")

	// then
	expectedError := "invalid character ']' looking for beginning of object key string"
	if err == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Expected error:\n`%s`\nGot:\n`%s`", expectedError, err)
	}
}
