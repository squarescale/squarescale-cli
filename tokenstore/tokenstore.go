package tokenstore

import (
	"bytes"
	"io/ioutil"
	"os"

	"github.com/squarescale/go-netrc/netrc"
)

func netrcFile() string {
	return os.Getenv("HOME") + "/.netrc"
}

func initNetrcFileIfNotExist() {
	var _, err = os.Stat(netrcFile())
	if os.IsNotExist(err) {
		var file, _ = os.Create(netrcFile())
		defer file.Close()
	}
}

// GetToken retrieves the Squarescale token in the token store.
func GetToken(host string) (string, error) {
	initNetrcFileIfNotExist()
	n, err := netrc.ParseFile(netrcFile())
	if err != nil {
		return "", err
	}

	for _, m := range n.FindMachines(host) {
		if m.Login != "" || m.Account != "" {
			continue
		}
		return m.Password, nil
	}

	return "", nil
}

// SaveToken persists the Squarescale token for the given host in the token store.
func SaveToken(host, token string) error {
	n, err := netrc.ParseFile(netrcFile())
	if err != nil && os.IsNotExist(err) {
		n, err = netrc.Parse(bytes.NewReader(nil))
	}
	if err != nil {
		return err
	}

	for _, m := range n.FindMachines(host) {
		if m.Login != "" || m.Account != "" {
			continue
		}
		m.UpdatePassword(token)
		return saveNetrc(n)
	}

	n.NewMachine(host, "", token, "")
	return saveNetrc(n)
}

func saveNetrc(n *netrc.Netrc) error {
	text, err := n.MarshalText()
	if err != nil {
		return err
	}

	return ioutil.WriteFile(netrcFile(), text, 0600)
}
