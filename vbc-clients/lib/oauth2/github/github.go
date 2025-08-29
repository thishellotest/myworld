// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package github provides constants for using OAuth2 to access Github.
package github // import "vbc/lib/oauth2/github"

import (
	"vbc/lib/oauth2/endpoints"
)

// Endpoint is Github's OAuth 2.0 endpoint.
var Endpoint = endpoints.GitHub
