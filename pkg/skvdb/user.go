/*
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
)

// InterceptorContextKey is the type used for keys in the context
type InterceptorContextKey string

const (
	// ContextUsernameKey is the key value used to store the username in the context
	ContextUsernameKey InterceptorContextKey = "username"
)

// GetUser returns the user name saved in the context
func GetUser(ctx context.Context) (string, bool) {
	// This value was saved in the context by the auth interceptor
	u, ok := ctx.Value(ContextUsernameKey).(string)
	return u, ok
}

// SetUser saves the context in the context
func SetUser(ctx context.Context, username string) context.Context{
	return context.WithValue(ctx, ContextUsernameKey, username)
}