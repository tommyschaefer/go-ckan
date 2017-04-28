package ckan

import (
	"context"
	"net/http"
)

// DataStoreService handles communication with the datastore methods
// of the CKAN API.
type DataStoreService service

// DataStoreMetadata holds the table metadata for a given CKAN datastore.
//
// CKAN API docs: http://docs.ckan.org/en/latest/maintaining/datastore.html#db-internals
type DataStoreMetadata struct {
	ID      string `json:"_id"`
	AliasOf string `json:"alias_of"`
	Name    string `json:"name"`
	OID     int    `json:"oid"`
}

// TeableMetadata fetches the datastore table metadata for a given CKAN
// instance.
//
// CKAN API docs: http://docs.ckan.org/en/latest/maintaining/datastore.html#db-internals
func (s *DataStoreService) TableMetadata(ctx context.Context, opt *ListOptions) ([]DataStoreMetadata, *http.Response, error) {
	path := "action/datastore_search?resource_id=_table_metadata"
	path, err := addOptions(path, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	var result struct {
		Records []DataStoreMetadata `json:"records"`
	}
	resp, err := s.client.Do(ctx, req, &result)
	if err != nil {
		return result.Records, resp, err
	}

	return result.Records, resp, nil
}
