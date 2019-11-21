/*
Copyright 2019 Isabella Pabón <isabella@chrysalix.org>
Copyright 2019 Luis Pabón <lpabon@chrysalix.org>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package skvdb

import (
	"context"
	"sync"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Skvdb is a simple interface to a secure kvdb
type Skvdb interface {
	Set(ctx context.Context, key, value string) error
	Get(ctx context.Context, key string) (string, error)
}

// The following is a very simple in-memory implementation that serves
// only as an example. For a better implementation, the library
// github.com/portworx/kvdb could be used.

// Element stores the user information
type Element struct {
	Value string
	Owner string
}

type SkvdbMem struct {
	lock sync.Mutex
	db   map[string]*Element
}

// New returns a new secure skvdb
func New() *SkvdbMem{
	return &SkvdbMem{
		db: make(map[string]*Element),
	}
}

func (s *SkvdbMem) Set(ctx context.Context, key, value string) error {

	// Get the key first to check if it exists and to check for permission
	_, err := s.Get(ctx, key)
	if !IsErrorNotFound(err) && err != nil {
		return err
	}

	// Get user information from context
	username, ok := GetUser(ctx)
	if !ok {
		return status.Errorf(codes.Internal, "Unable to determine user information")
	}

	// Create an element
	e := &Element{
		Value: value,
		Owner: username,
	}

	// Save element to database
	s.lock.Lock()
	defer s.lock.Unlock()
	s.db[key] = e

	return nil
}

func (s *SkvdbMem) Get(ctx context.Context, key string) (string, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	// Check permission
	username, ok := GetUser(ctx)
	if !ok {
		return "", status.Errorf(codes.Internal, "Unable to determine user information")
	}

	// Check if the key is there
	if e, ok := s.db[key]; !ok {
		// Key is not there
		return "", status.Errorf(codes.NotFound, "Key %s was not found", key)
	} else {
		// Check owner of the key
		if e.Owner != username {
			return "", status.Errorf(codes.PermissionDenied, "Access denied to key %s", key)
		}

		// Return the value
		return e.Value, nil
	}
}
