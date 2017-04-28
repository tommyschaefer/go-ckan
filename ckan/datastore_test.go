package ckan

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestDataStore_TableMetadata(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/action/datastore_search", func(w http.ResponseWriter, r *http.Request) {
		if m := "GET"; m != r.Method {
			t.Errorf("Request method = %v, want %v", r.Method, m)
		}
		fmt.Fprint(w, `{"success":true,"result":{"records":[{"_id":"1","alias_of":"2","name":"3","oid":4},{"_id":"5","alias_of":"6","name":"7","oid":8}]}}`)
	})

	packages, _, err := client.DataStore.TableMetadata(context.Background(), nil)
	if err != nil {
		t.Errorf("DataStore.TableMetadata returned error: %v", err)
	}

	expected := []DataStoreMetadata{
		{ID: "1", AliasOf: "2", Name: "3", OID: 4},
		{ID: "5", AliasOf: "6", Name: "7", OID: 8},
	}
	if !reflect.DeepEqual(packages, expected) {
		t.Errorf("DataStore.TableMetadata returned %+v, expected %+v", packages, expected)
	}
}

func TestDataStore_TableMetadata_withListOptions(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/action/datastore_search", func(w http.ResponseWriter, r *http.Request) {
		if id := r.FormValue("resource_id"); id != "_table_metadata" {
			t.Errorf("Expected resource ID to be _table_metadata, but got %v", id)
		}
		if l := r.FormValue("limit"); l != "10" {
			t.Errorf("Expected limit to be 10, but got %v", l)
		}
		if o := r.FormValue("offset"); o != "5" {
			t.Errorf("Expected offset to be 5, but got %v", o)
		}

		fmt.Fprint(w, `{"success":true,"result":{"records":[]}}`)
	})

	opts := &ListOptions{Limit: 10, Offset: 5}
	_, _, err := client.DataStore.TableMetadata(context.Background(), opts)
	if err != nil {
		t.Errorf("DataStore.TableMetadata returned error: %v", err)
	}
}
