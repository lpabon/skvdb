/*
Copyright 2019 Isabella Pabón <isabella@chrysalix.org>

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
	"google.golang.org/grpc/codes"
	"context"
	"testing"

	"github.com/lpabon/lputils/tests"
)

func TestSkvdbNew(t *testing.T) {
	a := New()
	tests.Assert(t, a != nil)
	tests.Assert(t, len(a.db) == 0)
}

func TestSkvdb(t *testing.T) {
	ctx := SetUser(context.Background(), "user1")
	badctx := SetUser(context.Background(), "user2")

	a := New()
	err := a.Set(ctx, "key", "val")
	tests.Assert(t, err == nil)
	v, err := a.Get(ctx, "key")
	tests.Assert(t, err == nil)
	tests.Assert(t, v == "val")

	err = a.Set(ctx, "key", "anotherval")
	tests.Assert(t, err == nil)
	v, err = a.Get(ctx, "key")
	tests.Assert(t, err == nil)
	tests.Assert(t, v == "anotherval")

	err = a.Set(badctx, "key", "aaa")
	tests.Assert(t, err != nil)
	s := FromError(err)
	tests.Assert(t, s.Code() == codes.PermissionDenied)

	v, err = a.Get(badctx, "key")
	tests.Assert(t, err != nil)
	s = FromError(err)
	tests.Assert(t, s.Code() == codes.PermissionDenied)

}
