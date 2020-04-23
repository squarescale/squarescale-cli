package squarescale_test

import "testing"

func checkAuthorization(t *testing.T, header string, token string) {
	expectedHeader := "bearer " + token
	if header != expectedHeader {
		t.Fatalf("Wrong token! Expected `%s`, got `%s`", expectedHeader, header)
	}
}

func checkPath(t *testing.T, expectedPath string, currentPath string) {
	if currentPath != expectedPath {
		t.Fatalf("Wrong path! Expected `%s`, got `%s`", expectedPath, currentPath)
	}
}
