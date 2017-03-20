package ckan

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestRegions_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/action/package_list", func(w http.ResponseWriter, r *http.Request) {
		if m := "GET"; m != r.Method {
			t.Errorf("Request method = %v, want %v", r.Method, m)
		}
		fmt.Fprint(w, `{"success":true,"result":["a","b","c"]}`)
	})

	packages, _, err := client.Packages.List(context.Background(), nil)
	if err != nil {
		t.Errorf("Packages.List returned error: %v", err)
	}

	expected := []string{"a", "b", "c"}
	if !reflect.DeepEqual(packages, expected) {
		t.Errorf("Packages.List returned %+v, expected %+v", packages, expected)
	}
}

func TestRegions_List_withListOptions(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/action/package_list", func(w http.ResponseWriter, r *http.Request) {
		if l := r.FormValue("limit"); l != "10" {
			t.Errorf("Expected limit to be 10, but got %v", l)
		}
		if o := r.FormValue("offset"); o != "5" {
			t.Errorf("Expected offset to be 5, but got %v", o)
		}

		fmt.Fprint(w, `{"success":true,"result":[]}`)
	})

	opts := &ListOptions{Limit: 10, Offset: 5}
	_, _, err := client.Packages.List(context.Background(), opts)
	if err != nil {
		t.Errorf("Packages.List returned error: %v", err)
	}
}
