package ckan

import (
	"context"
	"net/http"
)

// PackagesService handles communication with the package related methods of
// the CKAN API.
type PackagesService service

// List all of the available packages names for a given CKAN instance.
//
// CKAN API docs: http://docs.ckan.org/en/latest/api/#ckan.logic.action.get.package_list
func (s *PackagesService) List(ctx context.Context, opt *ListOptions) ([]string, *http.Response, error) {
	path := "action/package_list"
	path, err := addOptions(path, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	var packages []string
	resp, err := s.client.Do(ctx, req, &packages)
	if err != nil {
		return packages, resp, err
	}

	return packages, resp, nil
}
