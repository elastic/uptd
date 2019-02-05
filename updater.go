// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package uptd

import (
	"github.com/blang/semver"
)

// CheckResponse is given back by the Check method of an updater.
type CheckResponse struct {
	NeedsUpdate bool
	Current     semver.Version
	Latest      LatestResponse
}

// LatestResponse is given back by the Latest method of an update Provider.
type LatestResponse struct {
	Version    semver.Version
	URL        string
	PreRelease bool
}

// Provider represents the remote artifact repository where the updater checks
// if there's a newer version.
type Provider interface {
	// Latest queries the remote provider to check for a newer version of the
	// executable available in the remote artifact repository.
	Latest() (LatestResponse, error)
}

// Uptd checks if the current version is the latest one and thus up to date.
type Uptd struct {
	provider Provider
	current  semver.Version
}

// New instantiates a new Uptd from the sent Provider and current
// version. If the version is not semver compatible an error will be returned.
func New(provider Provider, current string) (Uptd, error) {
	version, err := semver.New(current)
	if err != nil {
		return Uptd{}, err
	}
	return Uptd{
		provider: provider,
		current:  *version,
	}, nil
}

// Check checks if there's a newer release available in the remote and
// returns a response that instructs whether or not the client should download
// the newer version if any is available.
func (u Uptd) Check() (CheckResponse, error) {
	latest, err := u.provider.Latest()
	if err != nil {
		return CheckResponse{}, err
	}

	return CheckResponse{
		NeedsUpdate: latest.Version.GT(u.current) && !latest.PreRelease,
		Current:     u.current,
		Latest:      latest,
	}, nil
}
