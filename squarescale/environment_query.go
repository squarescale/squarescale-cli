package squarescale

import (
	"bytes"
	"fmt"
	"github.com/olekukonko/tablewriter"
)

// QueryOptions holds environment variables filtering options to use with QueryVars
type QueryOptions struct {
	DisplayAll   bool
	ServiceName  string
	VariableName string
}

// QueryResult is an interface for objects returned by QueryVars when filtering
// environment variables.
// Its main purpose is to print the variables into a buffer for display.
type QueryResult interface {
	String() string
}

// QueryVars filters environment variables based on the given QueryOptions.
// It returns a QueryResult which can then be printed.
// If a requested service or variable don't exist, an error is returned.
func (env *Environment) QueryVars(options QueryOptions) (QueryResult, error) {
	if options.DisplayAll {
		return env, nil
	}

	var (
		result QueryResult
		err    error
	)

	if options.ServiceName != "" {
		result, err = env.GetServiceGroup(options.ServiceName)
		if err != nil {
			return nil, err
		}
	} else {
		result = env.Project
	}

	if options.VariableName != "" {
		result, err = result.(*VariableGroup).GetVariable(options.VariableName)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

// String returns a pretty string representation of the Project and Services
// variables with a title for each.
func (env *Environment) String() string {
	var lines = &bytes.Buffer{}

	prettyPrint(*env.Project, lines)
	for _, vg := range env.Services {
		prettyPrint(*vg, lines)
	}

	return lines.String()
}

// String returns a string representation of the VariableGroup's variables in a
// manner that is suitable for use in .env files.
func (vg *VariableGroup) String() string {
	if len(vg.Variables) == 0 {
		return "none\n"
	}

	var lines = &bytes.Buffer{}
	for _, variable := range vg.Variables {
		lines.WriteString(fmt.Sprintf("%s=%s\n", variable.Key, variable.Value))
	}
	return lines.String()
}

// String returns a string representation of the Variable's value.
func (v *Variable) String() string {
	return fmt.Sprintf("%s\n", v.Value)
}

func prettyPrint(vg VariableGroup, lines *bytes.Buffer) {
	lines.WriteString(fmt.Sprintf("%s\n", vg.Name))

	if len(vg.Variables) == 0 {
		lines.WriteString("  none\n")
	}

	data := make([][]string, 0, len(vg.Variables))
	for _, v := range vg.Variables {
		data = append(data, []string{v.Key, v.Value})
	}

	table := tablewriter.NewWriter(lines)
	table.AppendBulk(data)
	table.SetBorder(false)
	table.SetColumnSeparator("=")
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.Render()
}
