package command

import (
	"testing"

	"github.com/squarescale/squarescale-cli/squarescale"
)

func TestEmptySchedulingGroupArg(t *testing.T) {
	sched := getSchedulingGroupsArray("uuid", &squarescale.Client{}, "")
	if len(sched) != 0 {
		t.Error("Scheduling group list must be empty by default")
	}
}
