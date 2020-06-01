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
	"context"
	"strings"

	"github.com/blang/semver"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var _ Provider = new(GithubProvider)

// GetLatestFunc is used by the GithubProvider to discover newer releases
type GetLatestFunc func(ctx context.Context, owner, repo string) (*github.RepositoryRelease, *github.Response, error)

// GithubProvider is the implementation of an update Provider that is used
// by the Updater struct
type GithubProvider struct {
	latestFunc GetLatestFunc
	owner      string
	repo       string
}

// NewGithubProvider constructs a new GithubProvider from its parameters.
func NewGithubProvider(owner, repo, token string) (*GithubProvider, error) {
	if token == "" {
		return nil, errPersonalAccessTokenIsMissing
	}

	client := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	))

	return &GithubProvider{
		latestFunc: github.NewClient(client).Repositories.GetLatestRelease,
		owner:      owner,
		repo:       repo,
	}, nil
}

// Latest queries the remote Github repository to check for a newer version of
// the executable available.
func (u *GithubProvider) Latest() (LatestResponse, error) {
	release, _, err := u.latestFunc(context.Background(), u.owner, u.repo)
	if err != nil {
		return LatestResponse{}, err
	}

	version, err := semver.Parse(strings.TrimPrefix(release.GetTagName(), "v"))
	if err != nil {
		return LatestResponse{}, err
	}

	return LatestResponse{
		Version:    version,
		URL:        release.GetHTMLURL(),
		PreRelease: release.GetPrerelease() || release.GetDraft(),
	}, err
}
