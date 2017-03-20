// The client command demonstrates using ckan.Client to fetch a list of
// package names from the CKAN API for Denton, TX.
//
// Using the ckan.Client may be useful in some circumstances, especially to
// access API methods not supported by go-ckan.
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/tommyschaefer/go-ckan/ckan"
)

func main() {
	names, err := packageNames()
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not fetch packages: %v", err)
		os.Exit(2)
	}

	for _, name := range names {
		fmt.Println(name)
	}
}

func packageNames() ([]string, error) {
	var names []string

	cli, err := ckan.NewClient("http://data.cityofdenton.com/api/3/", nil)
	if err != nil {
		return names, fmt.Errorf("could not create new client: %v", err)
	}
	ctx := context.Background()

	req, err := cli.NewRequest("GET", "action/package_list?limit=10", nil)
	if err != nil {
		return names, fmt.Errorf("coud not create request: %v", err)
	}

	_, err = cli.Do(ctx, req, &names)
	if err != nil {
		return names, fmt.Errorf("could not perform request: %v", err)
	}

	return names, nil
}
