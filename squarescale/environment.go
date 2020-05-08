package squarescale

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Environment holds environment variables separated in 2 collections:
// - Project variables in a VariableGroup with Name Project and Variables which
// are defined project wide
// - Services variables in an array of VariableGroup with a service Name and
// variables defined in the service only for each
type Environment struct {
	Project  *VariableGroup
	Services []*VariableGroup
}

// NewEnvironment fetches environment variables from the API and returns a
// new Environment pointer if successful and an error otherwise.
func NewEnvironment(c *Client, project string) (*Environment, error) {
	code, body, err := c.get("/projects/" + project + "/environment")
	if err != nil {
		return nil, err
	}

	switch code {
	case http.StatusOK:
	case http.StatusNotFound:
		return nil, fmt.Errorf("Project '%s' not found", project)
	default:
		return nil, unexpectedHTTPError(code, body)
	}

	var env Environment
	if err = json.Unmarshal(body, &env); err != nil {
		return nil, err
	}

	return &env, nil
}

func valueInterface2String(val interface{}) string {
	switch v := val.(type) {
	case int:
		return fmt.Sprintf("%d", v)
	case float64:
		return fmt.Sprintf("%.0f", v)
	case string:
		return v
	default:
		panic(fmt.Sprintf("Don't know how to handle type %T!\n", v))
	}
	return ""
}

// UnmarshalJSON takes care of transforming the JSON representation of an
// Environment as sent by the API into a proper Environment struct.
func (env *Environment) UnmarshalJSON(body []byte) error {
	type Atom struct {
		Default map[string]interface{} `json:"default"`
		Custom  map[string]interface{} `json:"custom"`
	}

	type Global struct {
		Project    Atom            `json:"project"`
		PerService map[string]Atom `json:"per_service"`
	}

	var result Global
	if err := json.Unmarshal(body, &result); err != nil {
		return err
	}

	project := &VariableGroup{Name: "Project", Variables: []*Variable{}}
	for key, value := range result.Project.Default {
		newVar := &Variable{Key: key, Value: valueInterface2String(value), Predefined: true}
		project.Variables = append(project.Variables, newVar)
	}
	for key, value := range result.Project.Custom {
		newVar := &Variable{Key: key, Value: valueInterface2String(value), Predefined: false}
		project.Variables = append(project.Variables, newVar)
	}
	project.Variables = fold(project.Variables)
	env.Project = project

	for serviceName, serviceEnv := range result.PerService {
		service := &VariableGroup{Name: serviceName, Variables: []*Variable{}}

		for key, value := range serviceEnv.Default {
			newVar := &Variable{Key: key, Value: valueInterface2String(value), Predefined: true}
			service.Variables = append(service.Variables, newVar)
		}
		for key, value := range serviceEnv.Custom {
			newVar := &Variable{Key: key, Value: valueInterface2String(value), Predefined: false}
			service.Variables = append(service.Variables, newVar)
		}
		service.Variables = fold(merge(project.Variables, service.Variables))
		env.Services = append(env.Services, service)
	}

	return nil
}

// JSONObject takes care of creating a JSON representation of an Environment
// as expected by the API for a custom environment update.
func (env *Environment) JSONObject() (payload JSONObject, err error) {
	defer func() {
		if recover() != nil {
			err = fmt.Errorf("Cannot marshal environment into JSON")
		}
	}()

	global := make(map[string]string)
	for _, variable := range env.Project.Variables {
		if !variable.Predefined {
			global[variable.Key] = variable.Value
		}
	}

	perService := make(map[string]map[string]string)
	for _, service := range env.Services {
		perService[service.Name] = make(map[string]string)

		for _, variable := range service.Variables {
			if !variable.Predefined {
				perService[service.Name][variable.Key] = variable.Value
			}
		}
	}

	return JSONObject{"global": global, "per_service": perService}, nil
}

// CommitEnvironment sends a request to the API to update the project's
// custom environment with the custom variables defined in the Environment.
func (env *Environment) CommitEnvironment(c *Client, project string) error {
	envJSON, err := env.JSONObject()
	if err != nil {
		msg := "An error occurred while setting the environment variables."
		msg += " The environment was left unchanged."
		return fmt.Errorf(msg)
	}
	requestBody := &JSONObject{"environment": envJSON, "format": "json"}

	code, body, err := c.put("/projects/"+project+"/environment/custom", requestBody)
	if err != nil {
		return err
	}

	switch code {
	case http.StatusNoContent:
		return nil
	case http.StatusNotFound:
		return fmt.Errorf("Project '%s' not found", project)
	case http.StatusUnprocessableEntity:
		return fmt.Errorf("%s", body)
	default:
		return unexpectedHTTPError(code, body)
	}
}

// GetServiceGroup finds a VariableGroup in the Environment Services field
// with the given name and returns it.
// If no services are found, nil is returned with an error.
func (env *Environment) GetServiceGroup(serviceName string) (*VariableGroup, error) {
	for _, vg := range env.Services {
		if vg.Name == serviceName {
			return vg, nil
		}
	}
	return nil, fmt.Errorf("Could not find container '%s'", serviceName)
}

// VariableGroup holds a Name and an array of Variable for a specific group
// like a service or the global part of a project environment variables.
type VariableGroup struct {
	Name      string
	Variables []*Variable
}

// GetVariable finds a Variable in the VariableGroup Variables field with the
// given name and returns it.
// If no variables are found, nil is returned with an error.
func (vg *VariableGroup) GetVariable(variableName string) (*Variable, error) {
	var finalVar *Variable

	for _, variable := range vg.Variables {
		if variable.Key == variableName {
			if finalVar == nil || !variable.Predefined {
				finalVar = variable
			}
		}
	}
	if finalVar == nil {
		err := fmt.Errorf("Could not find variable '%s' for container '%s'",
			variableName, vg.Name)
		return nil, err
	}

	return finalVar, nil
}

// SetVariable sets the value of the variable in the VariableGroup with the
// given Key to the given Value.
// We have 3 situations:
// - if the variable doesn't exist, add a new custom variable
// - if the variable exists and is predefined, add a new custom variable
// - if the variable exists and is custom, modify its value
func (vg *VariableGroup) SetVariable(key, value string) {
	for _, variable := range vg.Variables {
		if variable.Key == key && !variable.Predefined {
			variable.Value = value
			return
		}
	}

	custom := &Variable{Key: key, Value: value, Predefined: false}
	vg.Variables = append(vg.Variables, custom)
}

// RemoveVariable removes the variable with the given Key from the
// VariableGroup.
// We have 3 situations:
// - if the variable doesn't exist, do nothing
// - if the variable exists and has a custom definition, remove it
// - if the variable exists and has is predefined only, return an error
func (vg *VariableGroup) RemoveVariable(key string) error {
	var predefinedOnly = false

	variables := vg.Variables
	for i, variable := range variables {
		if variable.Key == key {
			predefinedOnly = variable.Predefined

			if !predefinedOnly {
				variables = remove(variables, i)
			}
		}
	}
	if predefinedOnly {
		return fmt.Errorf("Cannot remove predefined variable")
	}

	vg.Variables = variables
	return nil
}

// Variable represents an environment variable which has a Key, a Value and is
// Predefined or custom (i.e. not Predefined).
type Variable struct {
	Key        string
	Value      string
	Predefined bool
}

func merge(bases, overrides []*Variable) []*Variable {
	var variables []*Variable

	for _, variable := range bases {
		vCopy := *variable
		variables = append(variables, &vCopy)
	}

	for _, variable := range overrides {
		vCopy := *variable
		variables = append(variables, &vCopy)
	}

	return variables
}

func fold(originals []*Variable) []*Variable {
	variables := originals[:0]

	for i, variable := range originals {
		var found = false
		for _, override := range originals[i+1:] {
			if override.Key == variable.Key {
				found = true
				break
			}
		}
		if !found {
			variables = append(variables, variable)
		}
	}

	return variables
}

func remove(variables []*Variable, pos int) []*Variable {
	copy(variables[pos:], variables[pos+1:])
	variables[len(variables)-1] = nil

	return variables[:len(variables)-1]
}
