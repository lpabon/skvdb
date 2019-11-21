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
	"fmt"
	"sync"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Skvdb interface {
	Set(ctx context.Context, key []byte, value []byte) error
	Get(ctx context.Context, key []byte) ([]byte, error)
}

type Element struct {
	Value []byte
	Owner string
}

type skdvb struct {
	lock sync.Mutex
	db map[string]*Value
}

func New() Skvdb {
	return &skdvb{
		db: make(map[string]*Value),
	}
}

func (s *skdvb) Set(ctx context.Context, key, value string) error  {

	// Get the key first to check if it exists and to check for permission
	_, found, err := get(ctx, key) 
	if err != nil {
		return err
	}

	// Get user information from context
	username, ok := getUser(ctx)
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

func (s *skdvb) Get(ctx context.Context, key string) (string, error)  {
	// Save element to database
	s.lock.Lock()
	defer s.lock.Unlock()

		// Check permission
		username, ok := getUser(ctx)
		if !ok {
		return "", status.Errorf(codes.Internal, "Unable to determine user information")
		}

	if e, ok := s.db[key]; !ok {
		return "", status.Errorf(codes.NotFound, "Key %s was not found", key)
	} else {
		if e.Owner != username {
			return "", status.Errorf(codes.PermissionDenied, "Access denied to key %s", key)
		}
	}
	
	return nil, nil
}

func getUser(ctx context.Context) (string, bool) {
	return ctx.Value("username").(string)
}