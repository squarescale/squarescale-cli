package command

import (
	"strings"
)

// Todo : add check for existing capabilities provided by the back-end

func getDockerCapabilitiesArray(capabilities string) []string {
	if capabilities == "" {
		return []string{}
	}
	return strings.Split(capabilities, ",")
}

func getDefaultDockerCapabilitiesArray() []string {
	return strings.Split("AUDIT_WRITE,CHOWN,DAC_OVERRIDE,FOWNER,FSETID,KILL,MKNOD,NET_BIND_SERVICE,SETFCAP,SETGID,SETPCAP,SETUID,SYS_CHROOT", ",")
}
