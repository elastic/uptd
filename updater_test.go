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
	"errors"
	"reflect"
	"testing"

	"github.com/blang/semver"
)

func newSemver(s string) semver.Version {
	v, _ := semver.Parse(s)
	return v
}

type nilProvider struct{}

// Latest queries the remote provider to check for a newer version of the
// executable available in the remote artifact repository.
func (p *nilProvider) Latest() (LatestResponse, error) { return LatestResponse{}, nil }

func TestNew(t *testing.T) {
	type args struct {
		provider Provider
		current  string
	}
	tests := []struct {
		name    string
		args    args
		want    Uptd
		wantErr bool
	}{
		{
			name: "Succeeds",
			args: args{
				provider: new(nilProvider),
				current:  "99.99.99",
			},
			want: Uptd{
				provider: new(nilProvider),
				current:  newSemver("99.99.99"),
			},
			wantErr: false,
		},
		{
			name: "Fails due invalid version",
			args: args{
				provider: new(nilProvider),
				current:  "WHATISTHISVERSION",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.provider, tt.args.current)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

type mockProvider struct {
	res LatestResponse
	err error
}

func (p *mockProvider) Latest() (LatestResponse, error) { return p.res, p.err }

func TestUptdCheck(t *testing.T) {
	type fields struct {
		provider Provider
		current  semver.Version
	}
	tests := []struct {
		name    string
		fields  fields
		want    CheckResponse
		wantErr bool
	}{
		{
			name: "Newer remote version prompts to update",
			fields: fields{
				provider: &mockProvider{
					res: LatestResponse{
						Version:    newSemver("99.99.99"),
						URL:        "http://host/path/version",
						PreRelease: false,
					},
					err: nil,
				},
				current: newSemver("99.99.98"),
			},
			want: CheckResponse{
				NeedsUpdate: true,
				Current:     newSemver("99.99.98"),
				Latest: LatestResponse{
					Version:    newSemver("99.99.99"),
					URL:        "http://host/path/version",
					PreRelease: false,
				},
			},
			wantErr: false,
		},
		{
			name: "Fails due Latest() call returning an error",
			fields: fields{
				provider: &mockProvider{
					err: errors.New("ERROR"),
				},
				current: newSemver("99.99.98"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := Uptd{
				provider: tt.fields.provider,
				current:  tt.fields.current,
			}
			got, err := u.Check()
			if (err != nil) != tt.wantErr {
				t.Errorf("Uptd.Check() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Uptd.Check() = %v, want %v", got, tt.want)
			}
		})
	}
}
