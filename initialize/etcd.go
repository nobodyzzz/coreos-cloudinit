package initialize

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/coreos/coreos-cloudinit/system"
)

type EtcdEnvironment map[string]string

func (ec EtcdEnvironment) String() (out string) {
	public := os.Getenv("COREOS_PUBLIC_IPV4")
	private := os.Getenv("COREOS_PRIVATE_IPV4")

	out += "[Service]\n"

	for key, val := range ec {
		key = strings.ToUpper(key)
		key = strings.Replace(key, "-", "_", -1)

		if public != "" {
			val = strings.Replace(val, "$public_ipv4", public, -1)
		}

		if private != "" {
			val = strings.Replace(val, "$private_ipv4", private, -1)
		}

		out += fmt.Sprintf("Environment=\"ETCD_%s=%s\"\n", key, val)
	}
	return
}

// Write an EtcdEnvironment to the appropriate path on disk for etcd.service
func WriteEtcdEnvironment(env EtcdEnvironment, root string) error {
	file := system.File{
		Path: path.Join(root, "etc", "systemd", "system", "etcd.service.d", "20-cloudinit.conf"),
		RawFilePermissions: "0644",
		Content: env.String(),
	}

	return system.WriteFile(&file)
}