package gh

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/google/go-github/v32/github"
)

// Client represents the gihub client which connects to the github api
type Client struct {
	client *github.Client
	log    *log.Logger
}

// New initializes a github client.
func New(client *http.Client, log *log.Logger) *Client {
	return &Client{
		client: github.NewClient(client),
		log:    log,
	}
}

//Fetcher represents github fetch operations
type Fetcher interface {
	ListRepositories(ctx context.Context, username string, opt *github.RepositoryListOptions) ([]*Repository, error)
	ListCommits(ctx context.Context, username, repoName string, opt *github.CommitsListOptions) ([]*Commit, error)
}

// ListRepositories lists all the public repositories of a user
func (g *Client) ListRepositories(ctx context.Context, username string, opt *github.RepositoryListOptions) ([]*Repository, error) {
	res, _, err := g.client.Repositories.List(ctx, username, opt)
	if err != nil {
		return nil, err
	}
	return mapFromRepository(res...), nil
}

// ListCommits lists the commits of a repository
func (g *Client) ListCommits(ctx context.Context, username, repoName string, opt *github.CommitsListOptions) ([]*Commit, error) {
	res, _, err := g.client.Repositories.ListCommits(ctx, username, repoName, opt)
	if err != nil {
		return nil, err
	}
	return mapFromCommit(res...), nil
}

type Repository struct {
	ID         int64 `gorm:"primaryKey"`
	NodeID     string
	Owner      string `gorm:"index"`
	Name       string
	CreatedAt  time.Time
	LastAccess time.Time
}

func mapFromRepository(in ...*github.Repository) []*Repository {
	var res []*Repository
	for _, v := range in {
		repo := Repository{}
		if v.Owner != nil {
			repo.Owner = *v.Owner.Login
		}
		if v.CreatedAt != nil {
			repo.CreatedAt = v.CreatedAt.Time
		}
		if v.ID != nil {
			repo.ID = *v.ID
		}
		if v.NodeID != nil {
			repo.NodeID = *v.NodeID
		}

		if v.Name != nil {
			repo.Name = *v.Name
		}
		res = append(res, &repo)
	}
	return res
}
func mapToRepository(in ...Repository) []*github.Repository {
	var res []*github.Repository
	for _, v := range in {
		res = append(res, &github.Repository{
			ID:        &v.ID,
			Owner:     &github.User{Login: &v.Owner},
			CreatedAt: &github.Timestamp{Time: v.CreatedAt},
			NodeID:    &v.NodeID,
			Name:      &v.Name,
		})
	}
	return res
}

type Commit struct {
	NodeID      string
	SHA         string `gorm:"primaryKey"`
	Author      string
	CommentsURL string
}

func mapFromCommit(in ...*github.RepositoryCommit) []*Commit {
	var res []*Commit
	for _, v := range in {
		commit := Commit{}
		if v.CommentsURL != nil {
			commit.CommentsURL = *v.CommentsURL
		}
		if v.SHA != nil {
			commit.SHA = *v.SHA
		}
		if v.NodeID != nil {
			commit.NodeID = *v.NodeID
		}
		if v.Author != nil && v.Author.Login != nil {
			commit.Author = *v.Author.Login
		}
		res = append(res, &commit)
	}
	return res
}

func mapToCommit(in ...*Commit) []*github.RepositoryCommit {
	var res []*github.RepositoryCommit
	for _, v := range in {
		commit := github.RepositoryCommit{
			NodeID:      &v.NodeID,
			SHA:         &v.SHA,
			CommentsURL: &v.CommentsURL,
		}
		commit.Author = &github.User{Login: &v.Author}
		res = append(res, &commit)
	}
	return res
}
