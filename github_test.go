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
	"errors"
	"reflect"
	"testing"

	"github.com/google/go-github/github"
)

func TestNewGithubProvider(t *testing.T) {
	type args struct {
		owner string
		repo  string
		token string
	}
	tests := []struct {
		name    string
		args    args
		want    *GithubProvider
		wantErr bool
	}{
		{
			name: "Succeeds with correct parameters",
			args: args{
				owner: "anowner",
				repo:  "arepo",
				token: "atoken",
			},
			want: &GithubProvider{
				owner: "anowner",
				repo:  "arepo",
			},
			wantErr: false,
		},
		{
			name: "Fails when a token is empty",
			args: args{
				owner: "anowner",
				repo:  "arepo",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewGithubProvider(tt.args.owner, tt.args.repo, tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewGithubProvider() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Go can't compare non nil functions
			if got != nil {
				got.latestFunc = nil
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGithubProvider() = %v, want %v", got, tt.want)
			}
		})
	}
}

func newLatestMockFunc(rel *github.RepositoryRelease, res *github.Response, err error) GetLatestFunc {
	return func(ctx context.Context, owner, repo string) (*github.RepositoryRelease, *github.Response, error) {
		return rel, res, err
	}
}

func newStringP(s string) *string { return &s }
func newBoolP(s bool) *bool       { return &s }

func TestGithubProviderLatest(t *testing.T) {
	type fields struct {
		latestFunc GetLatestFunc
		owner      string
		repo       string
	}
	tests := []struct {
		name    string
		fields  fields
		want    LatestResponse
		wantErr bool
	}{
		{
			name: "Latest() call succeeds",
			fields: fields{
				latestFunc: newLatestMockFunc(&github.RepositoryRelease{
					Name:    newStringP("99.99.99"),
					HTMLURL: newStringP("https://host/path/version"),
				}, nil, nil),
			},
			want: LatestResponse{
				Version:    newSemver("99.99.99"),
				URL:        "https://host/path/version",
				PreRelease: false,
			},
			wantErr: false,
		},
		{
			name: "Latest() call succeeds with Prerelease true",
			fields: fields{
				latestFunc: newLatestMockFunc(&github.RepositoryRelease{
					Name:    newStringP("99.99.99"),
					HTMLURL: newStringP("https://host/path/version"),
					Draft:   newBoolP(true),
				}, nil, nil),
			},
			want: LatestResponse{
				Version:    newSemver("99.99.99"),
				URL:        "https://host/path/version",
				PreRelease: true,
			},
			wantErr: false,
		},
		{
			name: "Latest() call succeeds with Prerelease true",
			fields: fields{
				latestFunc: newLatestMockFunc(&github.RepositoryRelease{
					Name:       newStringP("99.99.99"),
					HTMLURL:    newStringP("https://host/path/version"),
					Prerelease: newBoolP(true),
				}, nil, nil),
			},
			want: LatestResponse{
				Version:    newSemver("99.99.99"),
				URL:        "https://host/path/version",
				PreRelease: true,
			},
			wantErr: false,
		},
		{
			name: "Latest() call fails due wrong semver version",
			fields: fields{
				latestFunc: newLatestMockFunc(&github.RepositoryRelease{
					Name: newStringP("WHATISTHISVERSION"),
				}, nil, nil),
			},
			wantErr: true,
		},
		{
			name: "Latest() call fails due error returned by LatestFunc",
			fields: fields{
				latestFunc: newLatestMockFunc(nil, nil, errors.New("ERROR")),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &GithubProvider{
				latestFunc: tt.fields.latestFunc,
				owner:      tt.fields.owner,
				repo:       tt.fields.repo,
			}
			got, err := u.Latest()
			if (err != nil) != tt.wantErr {
				t.Errorf("GithubProvider.Latest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GithubProvider.Latest() = %v, want %v", got, tt.want)
			}
		})
	}
}
