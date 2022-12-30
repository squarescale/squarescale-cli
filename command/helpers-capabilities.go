package command

import (
	"fmt"
	"strings"
)

func contains( a []string, s string) bool {
	for _, t := range a {
		if t == s {
			return true
		}
	}
	return false
}

func getDockerCapabilitiesArray(capabilities string) ([]string, error) {
	var authorizedCapabilities = []string{"AUDIT_WRITE","CHOWN","DAC_OVERRIDE","FOWNER","FSETID","KILL","MKNOD","NET_BIND_SERVICE","NET_RAW","SETFCAP","SETGID","SETPCAP","SETUID","SYS_CHROOT","AUDIT_CONTROL","AUDIT_READ","BLOCK_SUSPEND","BPF","CHECKPOINT_RESTORE","DAC_READ_SEARCH","IPC_LOCK","IPC_OWNER","LEASE","LINUX_IMMUTABLE","MAC_ADMIN","MAC_OVERRIDE","NET_ADMIN","NET_BROADCAST","PERFMON","SYS_ADMIN","SYS_BOOT","SYS_MODULE","SYS_NICE","SYS_PACCT","SYS_PTRACE","SYS_RAWIO","SYS_RESOURCE","SYS_TIME","SYS_TTY_CONFIG","SYSLOG","WAKE_ALARM"}

	var dockerCapabilitiesArray = strings.Split(capabilities, ",")

	for _, capability := range dockerCapabilitiesArray {
		if !contains(authorizedCapabilities, capability)	{
			return nil, fmt.Errorf("No such capability : %s", capability)
		}
	}

	return dockerCapabilitiesArray, nil
}
