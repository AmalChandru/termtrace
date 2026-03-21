package replay

import "fmt"

func Run(path string) error {
	if path != "" {
		fmt.Printf("replay %q: not yet implemented\n", path)
	} else {
		fmt.Println("replay: not yet implemented (no workflow file specified)")
	}
	return nil
}
