// The client command demonstrates using ckan.Client to fetch a list of
// package names from the CKAN API for Denton, TX.
//
// Using the ckan.Client may be useful in some circumstances, especially to
// access API methods not supported by go-ckan.
package main

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/tommyschaefer/go-ckan/ckan"
)

func main() {
	u, err := url.Parse("http://data.cityofdenton.com/api/3/")
	if err != nil {
		log.Fatal(err)
	}

	cli := ckan.NewClient(u, nil)
	ctx := context.Background()

	req, err := cli.NewRequest("GET", "action/package_list?limit=10", nil)
	if err != nil {
		log.Fatal(err)
	}

	var names []string
	_, err = cli.Do(ctx, req, &names)
	if err != nil {
		log.Fatal(err)
	}

	for _, name := range names {
		fmt.Println(name)
	}
}
