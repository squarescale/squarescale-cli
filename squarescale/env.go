package squarescale

import (
	"log"
	"os"
	"strconv"
)

// AllowNodeSizeIds exports the value of the environment variable
// SQSC_CLI_ALLOW_NODE_SIZE_IDS
var AllowNodeSizeIds = boolEnv("SQSC_CLI_ALLOW_NODE_SIZE_IDS", false)

func boolEnv(varname string, defVal bool) bool {
	val := os.Getenv(varname)
	if val == "" {
		return defVal
	}
	res, err := strconv.ParseBool(val)
	if err != nil {
		log.Fatalf("Error parsing %s: %v", varname, err)
	}
	return res
}
