package squarescale_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/squarescale/squarescale-cli/squarescale"
)

func TestGetRedis(t *testing.T) {
	// GetRedisInfo
	t.Run("Nominal case on GetRedisInfo", nominalCaseOnGetRedisInfo)

	t.Run("Test redis not found on GetRedisInfo", UnknownRedisOnGetRedisInfo)
	t.Run("Test project not found on GetRedisInfo", UnknownProjectOnGetRedisInfo)

	// AddRedis
	t.Run("Nominal case on AddRedis", nominalCaseOnAddRedis)

	t.Run("Test duplicate redis name on AddRedis", DuplicateRedisErrorCaseOnAddRedis)
	t.Run("Test project not found on AddRedis", UnknownProjectOnAddRedis)

	// DeleteRedis
	t.Run("Nominal case on DeleteRedis", nominalCaseOnDeleteRedis)

	t.Run("Test redis not found on Deleteredis", UnknownRedisOnDeleteRedis)
	t.Run("Test project not found on Deleteredis", UnknownProjectOnDeleteRedis)

	// Error cases
	t.Run("Test HTTP client error on redis methods (get, add, delete and wait)", ClientHTTPErrorOnRedisMethods)
	t.Run("Test internal server error on redis methods (get, add, delete and wait)", InternalServerErrorOnRedisMethods)
	t.Run("Test badly JSON on redis methods (get)", CantUnmarshalOnRedisMethods)
	t.Run("redis", nominalCaseOnGetRedis)

}

func nominalCaseOnGetRedis(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkPath(t, "/projects/"+projectName+"/redis_databases", r.URL.Path)
		checkAuthorization(t, r.Header.Get("Authorization"), token)

		resBody := `{"redis_database_configs": [ { "name": "cache1" }, { "name": "cache2" } ]}`

		w.Header().Set("Content-Type", "application/json")

		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	redis, err := cli.GetRedis(projectName)

	// then
	var expectedInt int
	var expectedString string

	if err != nil {
		t.Fatalf("Expect no error, got `%s`", err)
	}

	expectedInt = 2
	if len(redis) != expectedInt {
		t.Fatalf("Expect redes to contain one element `%d`, but got actually `%d`", expectedInt, len(redis))
	}

	expectedString = "cache1"
	if redis[0].Name != expectedString {
		t.Errorf("Expect redis.Name `%s`, got `%s`", expectedString, redis[0].Name)
	}

	expectedString = "cache2"
	if redis[1].Name != expectedString {
		t.Errorf("Expect redis.Name `%s`, got `%s`", expectedString, redis[1].Name)
	}
}

func nominalCaseOnGetRedisInfo(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkPath(t, "/projects/"+projectName+"/redis_databases", r.URL.Path)
		checkAuthorization(t, r.Header.Get("Authorization"), token)

		resBody := `{"redis_database_configs":
		[
			{
				"name": "cache1"
			},
			{
				"name": "cache2"
			}
		]
		}`

		w.Header().Set("Content-Type", "application/json")

		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	redis, err := cli.GetRedisInfo(projectName, "cache1")

	// then
	var expectedString string

	if err != nil {
		t.Fatalf("Expect no error, got `%s`", err)
	}

	expectedString = "cache1"
	if redis.Name != expectedString {
		t.Errorf("Expect redis.Name `%s`, got `%s`", expectedString, redis.Name)
	}
}

func nominalCaseOnAddRedis(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkPath(t, "/projects/"+projectName+"/redis_databases", r.URL.Path)
		checkAuthorization(t, r.Header.Get("Authorization"), token)

		resBody := `
		{
			"name": "cache2"
		}
		`

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(201)

		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	err := cli.AddRedis(projectName, "cache2")

	// then
	if err != nil {
		t.Fatalf("Expect no error, got `%s`", err)
	}
}

func nominalCaseOnDeleteRedis(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"
	redisName := "my-redis"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkPath(t, "/projects/"+projectName+"/redis_databases/"+redisName, r.URL.Path)
		checkAuthorization(t, r.Header.Get("Authorization"), token)

		resBody := `
		null
		`

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)

		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	err := cli.DeleteRedis(projectName, redisName)

	// then
	if err != nil {
		t.Fatalf("Expect no error, got `%s`", err)
	}
}

func UnknownRedisOnGetRedisInfo(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resBody := `{"redis_database_configs":
		[
			{
				"name": "cache1"
			},
			{
				"name": "cache2"
			}
		]
		}`

		w.Header().Set("Content-Type", "application/json")

		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	_, err := cli.GetRedisInfo(projectName, "missing-redis")

	// then
	expectedError := "Redis 'missing-redis' not found for project 'my-project'"
	if err == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, err)
	}
}

func UnknownProjectOnGetRedisInfo(t *testing.T) {
	// given
	token := "some-token"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resBody := `
		{
			"error": "No project found for config name: not-a-project"
		}
		`

		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(404)
		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	_, err := cli.GetRedisInfo("not-a-project", "a-redis")

	// then
	expectedError := "Project 'not-a-project' does not exist"
	if err == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, err)
	}
}

func DuplicateRedisErrorCaseOnAddRedis(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resBody := `
		{"error":"PG::UniqueViolation: ERROR:  duplicate key value violates unique constraint \"index_redis_on_cluster_id_and_name\"\nDETAIL:  Key (cluster_id, name)=(6163, cache1) already exists.\n: INSERT INTO \"redis_databases\" (\"name\", \"cluster_id\") VALUES ($1, $2) RETURNING \"id\""}
		`

		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(409)
		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	err := cli.AddRedis(projectName, "cache1")

	// then
	expectedError := "Redis already exist on project 'my-project': cache1"
	if err == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, err)
	}
}

func UnknownProjectOnAddRedis(t *testing.T) {
	// given
	token := "some-token"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resBody := `
		{"error":"No project found for config name: another-project"}
		`

		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(404)
		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	err := cli.AddRedis("another-project", "cache1")

	// then
	expectedError := `Error: No project found for config name: another-project`
	if err == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, err)
	}
}

func UnknownProjectOnDeleteRedis(t *testing.T) {
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
	err := cli.DeleteRedis("unknown-project", "cache1")

	// then
	expectedError := `{"error":"No project found for config name: unknown-project"}`
	if err == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, err)
	}
}

func UnknownRedisOnDeleteRedis(t *testing.T) {
	// given
	token := "some-token"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resBody := `{"error":"Couldn't find Redis with [WHERE \"redis_databases\".\"cluster_id\" = $1 AND \"redis_databases\".\"name\" = $2]"}`

		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(404)
		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	err := cli.DeleteRedis("my-project", "cache1")

	// then
	expectedError := `{"error":"No redis found for name: cache1"}`
	if err == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, err)
	}
}

func ClientHTTPErrorOnRedisMethods(t *testing.T) {
	// given
	token := "some-token"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	_, errOnGet := cli.GetRedis("another-project")
	errOnAdd := cli.AddRedis("another-project", "cache1")
	errOnDelete := cli.DeleteRedis("my-project", "cache1")

	// then
	if errOnGet == nil {
		t.Errorf("Error is not raised on GetRedis")
	}

	if errOnAdd == nil {
		t.Errorf("Error is not raised on AddRedis")
	}

	if errOnDelete == nil {
		t.Errorf("Error is not raised on DeleteRedis")
	}
}

func InternalServerErrorOnRedisMethods(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	_, errOnGetRedis := cli.GetRedis(projectName)
	errOnAddRedis := cli.AddRedis(projectName, "cache1")
	errOnDeleteRedis := cli.DeleteRedis(projectName, "cache1")

	// then
	expectedError := "An unexpected error occurred (code: 500)"
	if errOnGetRedis == nil {
		t.Errorf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", errOnGetRedis) != expectedError {
		t.Errorf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, errOnGetRedis)
	}

	if errOnAddRedis == nil {
		t.Errorf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", errOnAddRedis) != expectedError {
		t.Errorf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, errOnAddRedis)
	}

	if errOnDeleteRedis == nil {
		t.Errorf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", errOnDeleteRedis) != expectedError {
		t.Errorf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, errOnDeleteRedis)
	}
}

func CantUnmarshalOnRedisMethods(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resBody := `{]`

		w.Header().Set("Content-Type", "application/json")

		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	_, errOnGetRedis := cli.GetRedis(projectName)

	// then
	expectedError := "invalid character ']' looking for beginning of object key string"
	if errOnGetRedis == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", errOnGetRedis) != expectedError {
		t.Fatalf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, errOnGetRedis)
	}
}
