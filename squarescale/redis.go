package squarescale

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Redis describes a project container as returned by the SquareScale API
type Redis struct {
	RedisDatabaseConfigs []struct {
		Name string `json:"name"`
	} `json:"redis_database_configs"`
}

// RedisDbConfig hold config for one redis
type RedisDbConfig struct {
	Name string "json:\"name\""
}

// {"redis_database_configs":[{"name":"plop1"}]}

// GetRedis gets all the redis attached to a Project
func (c *Client) GetRedis(projectUUID string) ([]RedisDbConfig, error) {
	code, body, err := c.get("/projects/" + projectUUID + "/redis_databases")
	if err != nil {
		return []RedisDbConfig{}, err
	}

	switch code {
	case http.StatusOK:
	case http.StatusNotFound:
		return []RedisDbConfig{}, fmt.Errorf("Project '%s' does not exist", projectUUID)
	default:
		return []RedisDbConfig{}, unexpectedHTTPError(code, body)
	}

	var redesDb Redis
	var redesDbConfig []RedisDbConfig

	// log.Printf("redis json respond : %s\n", body)

	if err := json.Unmarshal(body, &redesDb); err != nil {
		return []RedisDbConfig{}, err
	}
	for _, redis := range redesDb.RedisDatabaseConfigs {
		redesDbConfig = append(redesDbConfig, redis)
	}
	return redesDbConfig, nil

}

// GetRedisInfo gets the redis of a project based on its name.
func (c *Client) GetRedisInfo(projectUUID, name string) (RedisDbConfig, error) {
	redes, err := c.GetRedis(projectUUID)
	if err != nil {
		return RedisDbConfig{}, err
	}

	// log.Printf("GetRedisInfo iterator decode : %s\n", redes)

	for _, redis := range redes {
		if redis.Name == name {
			// if redis.Name == name {
			return redis, nil
		}
	}

	return RedisDbConfig{}, fmt.Errorf("Redis '%s' not found for project '%s'", name, projectUUID)
}

// DeleteRedis delete redis of a project based on its name.
func (c *Client) DeleteRedis(projectUUID string, name string) error {
	url := fmt.Sprintf("/projects/%s/redis_databases/%s", projectUUID, name)
	code, body, err := c.delete(url)
	if err != nil {
		return err
	}

	switch code {
	case http.StatusOK:
		return nil
	case http.StatusNotFound:
		if fmt.Sprintf("%s", body) == `{"error":"Couldn't find Redis with [WHERE \"redis_databases\".\"cluster_id\" = $1 AND \"redis_databases\".\"name\" = $2]"}` {
			return fmt.Errorf("{\"error\":\"No redis found for name: %s\"}", name)
		}
		return fmt.Errorf("%s", body)
	default:
		return unexpectedHTTPError(code, body)
	}
}

// AddRedis add redis of a project based on its name.
func (c *Client) AddRedis(projectUUID string, name string) error {
	payload := JSONObject{
		"name": name,
	}

	url := fmt.Sprintf("/projects/%s/redis_databases", projectUUID)
	code, body, err := c.post(url, &payload)
	if err != nil {
		return err
	}

	switch code {
	case http.StatusCreated:
		return nil
	case http.StatusConflict:
		return fmt.Errorf("Redis already exist on project '%s': %s", projectUUID, name)
	case http.StatusNotFound:
		return unexpectedHTTPError(code, body)
	default:
		return unexpectedHTTPError(code, body)
	}
}
