package tokenstore

import (
	"bytes"
	"io/ioutil"
	"os"

	"github.com/squarescale/go-netrc/netrc"
)

func NetrcFile() string {
	return os.Getenv("HOME") + "/.netrc"
}

func GetToken(host string) (string, error) {
	n, err := netrc.ParseFile(NetrcFile())
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

func SaveToken(host, token string) error {
	n, err := netrc.ParseFile(NetrcFile())
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

	return ioutil.WriteFile(NetrcFile(), text, 0600)
}
