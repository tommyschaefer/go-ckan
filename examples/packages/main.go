// The packages command demonstrates using the packages service to fetch a
// list of package names from the CKAN API for Denton, TX.
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/tommyschaefer/go-ckan/ckan"
)

func main() {
	cli, err := ckan.NewClient("http://data.cityofdenton.com/api/3/", nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not initialize client: %v", err)
		os.Exit(2)
	}
	ctx := context.Background()

	packages, _, err := cli.Packages.List(ctx, &ckan.ListOptions{Limit: 10})
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not fetch packages: %v", err)
		os.Exit(2)
	}

	for _, pkg := range packages {
		fmt.Println(pkg)
	}
}
