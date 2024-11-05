package ssg

import (
	"fmt"
	"os"
	"testing"
)

func TestHandler(t *testing.T) {
	if err := testHandler(); err != nil {
		t.Fatal(err)
	}
}

func testHandler() error {
	target, err := os.MkdirTemp("", "")
	if err != nil {
		return fmt.Errorf("cannot make tempdir: %w", err)
	}
	bindings, err := GenerateSiteWithBindings(
		"test", target,
		"../../theme/lit", "algol_nu",
		"HEADER", "FOOTER",
		map[string]CustomPage{
			"/xyz": NewCustomPage("Xr0", "<x>Xr0</x>"),
			"/lmn": NewCustomPage("hello", "<b>hello, world</b>"),
		},
	)
	if err != nil {
		return fmt.Errorf("cannot generate: %w", err)
	}
	for _, url := range []string{
		"/",
		"/abc/def",
		"/nest/post",
		"/nest-no-ignore/README",
		"/nest-no-ignore/post",
		"/xyz",
		"/lmn",
	} {
		if file, ok := bindings[url]; !ok {
			return fmt.Errorf("%q not found", url)
		} else {
			fmt.Println(url, file)
		}
	}
	return nil
}

func readfile(file string) (string, error) {
	b, err := os.ReadFile(file)
	return string(b), err
}
