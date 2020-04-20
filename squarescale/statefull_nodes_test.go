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
			t.Fatalf("Wrong path ! Expected %s, got %s", "/projects/my-project/statefull_nodes", path)
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
			t.Fatalf("Wrong path ! Expected %s, got %s", "bearer some-token", r.Header.Get("Authorization"))
		}

		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	statefullNodes, err := cli.GetStatefullNodes(projectName)

	if err != nil {
		t.Fatalf("Expect no error, got %s", err)
	}

	if len(statefullNodes) != 2 {
		t.Fatalf("Expect statefull_nodes to contain one element %d, but got actually %d", 2, len(statefullNodes))
	}

	if statefullNodes[0].ID != 23 {
		t.Errorf("Expect statefullNodeID `%d`, got `%d`", 23, statefullNodes[0].ID)
	}

	if statefullNodes[0].Name != "nodeb" {
		t.Errorf("Expect statefullNodeName `%s`, got `%s`", "nodeb", statefullNodes[0].Name)
	}

	if statefullNodes[0].NodeType != "t2.micro" {
		t.Errorf("Expect statefullNodeNodeType `%s`, got `%s`", "t2.micro", statefullNodes[0].NodeType)
	}

	if statefullNodes[0].Zone != "eu-west-1b" {
		t.Errorf("Expect statefullNodeZone `%s`, got `%s`", "eu-west-1b", statefullNodes[0].Zone)
	}

	if statefullNodes[0].Status != "provisionned" {
		t.Errorf("Expect statefullNodeStatus `%s`, got `%s`", "provisionned", statefullNodes[0].Status)
	}

	if statefullNodes[1].ID != 22 {
		t.Errorf("Expect statefullNodeID `%d`, got `%d`", 22, statefullNodes[1].ID)
	}

	if statefullNodes[1].Name != "test1" {
		t.Errorf("Expect statefullNodeName `%s`, got `%s`", "test1", statefullNodes[1].Name)
	}

	if statefullNodes[1].NodeType != "t2.micro" {
		t.Errorf("Expect statefullNodeNodeType `%s`, got `%s`", "t2.micro", statefullNodes[1].NodeType)
	}

	if statefullNodes[1].Zone != "eu-west-1a" {
		t.Errorf("Expect statefullNodeZone `%s`, got `%s`", "eu-west-1a", statefullNodes[1].Zone)
	}

	if statefullNodes[1].Status != "not_provisionned" {
		t.Errorf("Expect statefullNodeStatus `%s`, got `%s`", "not_provisionned", statefullNodes[1].Status)
	}
}

func UnknownProjectOnGetStatefullNodes(t *testing.T) {
	// given
	token := "some-token"
	projectName := "unknown-project"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var path string = r.URL.Path

		if path != "/projects/"+projectName+"/statefull_nodes" {
			t.Fatalf("Wrong path ! Expected %s, got %s", "/projects/unknown-project/statefull_nodes", path)
		}

		resBody := `
		{"error":"No project found for config name: unknown-project"}
		`

		w.Header().Set("Content-Type", "application/json")

		if (r.Header.Get("Authorization")) != "bearer some-token" {
			t.Fatalf("Wrong path ! Expected %s, got %s", "bearer some-token", r.Header.Get("Authorization"))
		}

		w.WriteHeader(404)
		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	_, err := cli.GetStatefullNodes("unknown-project")

	if err == nil {
		t.Fatalf("Error is not raised %s", err)
	}

}

func CannotUnmarshalOnGetStatefullNodes(t *testing.T) {
	// given
	token := "some-token"
	projectName := "unknown-project"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var path string = r.URL.Path

		if path != "/projects/"+projectName+"/statefull_nodes" {
			t.Fatalf("Wrong path ! Expected %s, got %s", "/projects/unknown-project/statefull_nodes", path)
		}

		resBody := `{]`

		w.Header().Set("Content-Type", "application/json")

		if (r.Header.Get("Authorization")) != "bearer some-token" {
			t.Fatalf("Wrong path ! Expected %s, got %s", "bearer some-token", r.Header.Get("Authorization"))
		}

		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	_, err := cli.GetStatefullNodes("unknown-project")

	// then
	if err == nil {
		t.Fatalf("Error is not raised %s", err)
	}

}

func HTTPErrorOnGetStatefullNodes(t *testing.T) {
	// given
	token := "some-token"
	projectName := "bad-project"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var path string = r.URL.Path

		if path != "/projects/"+projectName+"/statefull_nodes" {
			t.Fatalf("Wrong path ! Expected %s, got %s", "/projects/bad-project/statefull_nodes", path)
		}

		resBody := `
		{"error":"Hu ho, dummy error"}
		`

		w.Header().Set("Content-Type", "application/json")

		if (r.Header.Get("Authorization")) != "bearer some-token" {
			t.Fatalf("Wrong path ! Expected %s, got %s", "bearer some-token", r.Header.Get("Authorization"))
		}

		w.WriteHeader(500)
		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	_, err := cli.GetStatefullNodes("bad-project")

	// then
	if err == nil {
		t.Fatalf("Error is not raised %s", err)
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
			t.Fatalf("Wrong path ! Expected %s, got %s", "/projects/my-project/statefull_nodes", path)
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
			t.Fatalf("Wrong path ! Expected %s, got %s", "bearer some-token", r.Header.Get("Authorization"))
		}

		w.WriteHeader(201)
		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	newStatefullNode, err := cli.AddStatefullNode(projectName, nodeName, "t2.micro", "eu-west-1a")

	if err != nil {
		t.Fatalf("Expect no error, got `%s`", err)
	}

	if newStatefullNode.ID != 23 {
		t.Errorf("Expect statefullNodeID `%d`, got `%d`", 23, newStatefullNode.ID)
	}

	if newStatefullNode.Name != "node1a" {
		t.Errorf("Expect statefullNodeName `%s`, got `%s`", "node1a", newStatefullNode.Name)
	}

	if newStatefullNode.NodeType != "t2.micro" {
		t.Errorf("Expect statefullNodeNodeType `%s`, got `%s`", "t2.micro", newStatefullNode.NodeType)
	}

	if newStatefullNode.Zone != "eu-west-1c" {
		t.Errorf("Expect statefullNodeZone `%s`, got `%s`", "eu-west-1c", newStatefullNode.Zone)
	}

	if newStatefullNode.Status != "not_provisionned" {
		t.Errorf("Expect statefullNodeStatus `%s`, got `%s`", "not_provisionned", newStatefullNode.Status)
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
			t.Fatalf("Wrong path ! Expected %s, got %s", "/projects/unknown-project/statefull_nodes", path)
		}

		resBody := `
		{"error":"No project found for config name: unknown-project"}
		`

		w.Header().Set("Content-Type", "application/json")

		if (r.Header.Get("Authorization")) != "bearer some-token" {
			t.Fatalf("Wrong path ! Expected %s, got %s", "bearer some-token", r.Header.Get("Authorization"))
		}

		w.WriteHeader(404)
		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	_, err := cli.AddStatefullNode("unknown-project", nodeName, "t2.micro", "eu-west-1a")

	if err == nil {
		t.Fatalf("Error is not raised")
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
			t.Fatalf("Wrong path ! Expected %s, got %s", "/projects/my-project/statefull_nodes", path)
		}

		resBody := `{"error":"PG::UniqueViolation: ERROR:  duplicate key value violates unique constraint \"index_statefull_nodes_on_cluster_id_and_name\"\nDETAIL:  Key (cluster_id, name)=(6163, nodeb) already exists.\n: INSERT INTO \"statefull_nodes\" (\"name\", \"node_type\", \"zone\", \"cluster_id\", \"status\") VALUES ($1, $2, $3, $4, $5) RETURNING \"id\""}`

		w.Header().Set("Content-Type", "application/json")

		if (r.Header.Get("Authorization")) != "bearer some-token" {
			t.Fatalf("Wrong path ! Expected %s, got %s", "bearer some-token", r.Header.Get("Authorization"))
		}

		w.WriteHeader(409)
		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	_, err := cli.AddStatefullNode(projectName, nodeName, "t2.micro", "eu-west-1a")

	if err == nil {
		t.Fatalf("Error is not raised")
	}

	if fmt.Sprintf("%s", err) != "Statefull node already exist on project 'my-project': node1a" {
		t.Fatalf("Error raised is node `Statefull node already exist on project 'my-project': node1a`: `%s`", err)
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
			t.Fatalf("Wrong path ! Expected %s, got %s", "/projects/bad-project/statefull_nodes", path)
		}

		resBody := `
		{"error":"Hu ho, dummy error"}
		`

		w.Header().Set("Content-Type", "application/json")

		if (r.Header.Get("Authorization")) != "bearer some-token" {
			t.Fatalf("Wrong path ! Expected %s, got %s", "bearer some-token", r.Header.Get("Authorization"))
		}

		w.WriteHeader(500)
		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	_, err := cli.AddStatefullNode(projectName, nodeName, "t2.micro", "eu-west-1a")

	// then
	if err == nil {
		t.Fatalf("Error is not raised %s", err)
	}

}

func CannotUnmarshalOnAddStatefullNode(t *testing.T) {
	// given
	token := "some-token"
	projectName := "unknown-project"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var path string = r.URL.Path

		if path != "/projects/"+projectName+"/statefull_nodes" {
			t.Fatalf("Wrong path ! Expected %s, got %s", "/projects/unknown-project/statefull_nodes", path)
		}

		resBody := `{]`

		w.Header().Set("Content-Type", "application/json")

		if (r.Header.Get("Authorization")) != "bearer some-token" {
			t.Fatalf("Wrong path ! Expected %s, got %s", "bearer some-token", r.Header.Get("Authorization"))
		}

		w.WriteHeader(201)
		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	_, err := cli.AddStatefullNode(projectName, "node-in-error", "t2.micro", "eu-west-1a")

	// then
	if err == nil {
		t.Fatalf("Error is not raised")
	}
}
