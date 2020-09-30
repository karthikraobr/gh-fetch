package gh

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-github/v32/github"
)

func NewTestClient(q interface{}, err error) *github.Client {
	var body []byte
	status := http.StatusOK
	if err != nil {
		body, _ = json.Marshal(err.Error())
		status = http.StatusInternalServerError
	} else {
		body, _ = json.Marshal(q)
	}
	fakeHttpClient := NewFakeHttpClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: status,
			// Send response to be tested
			Body: ioutil.NopCloser(bytes.NewReader(body)),
			// Must be set to non-nil value or it panics
			Header: make(http.Header),
		}
	})
	return github.NewClient(fakeHttpClient)
}

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func NewFakeHttpClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func TestClient_ListRepositories(t *testing.T) {
	repo1 := Repository{
		ID:        1,
		CreatedAt: time.Now(),
		Name:      "blog",
		NodeID:    "1",
		Owner:     "me",
	}
	type fields struct {
		client *github.Client
		log    *log.Logger
	}
	type args struct {
		ctx      context.Context
		username string
		opt      *github.RepositoryListOptions
	}
	tests := map[string]struct {
		fields  fields
		args    args
		want    []*Repository
		wantErr bool
	}{
		"valid": {
			fields: fields{
				client: NewTestClient(mapToRepository(repo1), nil),
				log:    &log.Logger{},
			},
			args: args{
				ctx:      context.Background(),
				opt:      &github.RepositoryListOptions{Type: "public"},
				username: "me",
			},
			want: []*Repository{&repo1},
		},
		"error": {
			fields: fields{
				client: NewTestClient(nil, errors.New("not found")),
				log:    &log.Logger{},
			},
			args: args{
				ctx:      context.Background(),
				opt:      &github.RepositoryListOptions{Type: "public"},
				username: "me",
			},
			wantErr: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			g := &Client{
				client: tt.fields.client,
				log:    tt.fields.log,
			}
			got, err := g.ListRepositories(tt.args.ctx, tt.args.username, tt.args.opt)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ListRepositories() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("Client.ListRepositories() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_ListCommits(t *testing.T) {
	commit := Commit{
		Author:      "me",
		CommentsURL: "url",
		NodeID:      "nodeid",
		SHA:         "sha",
	}
	type fields struct {
		client *github.Client
		log    *log.Logger
	}
	type args struct {
		ctx      context.Context
		username string
		reponame string
		opt      *github.CommitsListOptions
	}
	tests := map[string]struct {
		fields  fields
		args    args
		want    []*Commit
		wantErr bool
	}{
		"valid": {
			fields: fields{
				client: NewTestClient(mapToCommit(&commit), nil),
				log:    &log.Logger{},
			},
			args: args{
				ctx:      context.Background(),
				opt:      &github.CommitsListOptions{},
				username: "me",
				reponame: "repo",
			},
			want: []*Commit{&commit},
		},
		"error": {
			fields: fields{
				client: NewTestClient(nil, errors.New("not found")),
				log:    &log.Logger{},
			},
			args: args{
				ctx:      context.Background(),
				opt:      &github.CommitsListOptions{},
				username: "me",
				reponame: "repo",
			},
			wantErr: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			g := &Client{
				client: tt.fields.client,
				log:    tt.fields.log,
			}
			got, err := g.ListCommits(tt.args.ctx, tt.args.username, tt.args.reponame, tt.args.opt)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Commits() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("Client.Commits() = %v, want %v", got, tt.want)
			}
		})
	}
}
