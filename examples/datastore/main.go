// The datastore command demonstrates using the datastore service to fetch the
// datastore table metadata from the CKAN API for Denton, TX.
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

	var m []ckan.DataStoreMetadata

	for l, o := 100, 0; ; o += l {
		metadata, _, err := cli.DataStore.TableMetadata(ctx, &ckan.ListOptions{Limit: l, Offset: o})
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not fetch table metadata: %v", err)
			os.Exit(2)
		}
		m = append(m, metadata...)

		if len(metadata) < l {
			break
		}
	}

	fmt.Println(len(m))
}
