package command

import "flag"

func isFlagPassed(name string, fs *flag.FlagSet) bool {
	found := false
	fs.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
			return
		}
	})
	return found
}
