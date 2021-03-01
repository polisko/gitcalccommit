package gitcommits

import (
	"context"
	"errors"
	"fmt"

	"github.com/shurcooL/githubv4"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

// GitCommits provides core functionality, needed by homework
type GitCommits struct {
	token         string
	DefaultOwner  string
	DefaultRepo   string
	DefaultBranch string
	client        *githubv4.Client
}

//ErrMissingOrBadAuthToken represents
var ErrMissingOrBadAuthToken = errors.New("Missing or wrong Git authorization token")

// NewGitCommits return new GitCommits object, based on authorization token.
// Since github GraphQL API needs it
func NewGitCommits(token string) (*GitCommits, error) {
	if token == "" {
		return nil, ErrMissingOrBadAuthToken
	}
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	ctx := context.Background()
	httpClient := oauth2.NewClient(ctx, src)

	client := githubv4.NewClient(httpClient)

	//test token
	err := client.Query(ctx, &viewer, nil)
	if err != nil {
		return nil, ErrMissingOrBadAuthToken
	}
	log.Tracef("Viewer identity: %+v", viewer)
	return &GitCommits{client: client, token: token, DefaultBranch: "main"}, nil
}

// FindCommitWithCtx searches for SHA (OID) (shortened too) in given repository
// if it finds, return full SHA
// if not, nil and corresponding error
func (gc *GitCommits) FindCommitWithCtx(ctx context.Context, oid string) (*CommitWithTS, error) {
	// query getCommitsSHA($owner: String!, $repo: String!, $expression: String!) {
	// 	repository(owner: $owner, name: $repo) {
	// 	  object(expression: $expression) {
	// 		... on Commit {
	// 		  oid
	// 		  messageHeadline
	// 		  committedDate
	// 		}
	// 	  }
	// 	}
	//   }
	var commit struct {
		Repository struct {
			Object struct {
				Commit struct {
					OID           githubv4.GitObjectID
					CommittedDate githubv4.GitTimestamp
				} `graphql:"... on Commit"`
			} `graphql:"object(expression: $expression)"`
		} `graphql:"repository(owner: $owner, name: $repo)"`
	}
	params := map[string]interface{}{
		"owner":      githubv4.String(gc.DefaultOwner),
		"repo":       githubv4.String(gc.DefaultRepo),
		"expression": githubv4.String(oid)}

	err := gc.client.Query(ctx, &commit, params)
	if err != nil {
		return nil, fmt.Errorf("Error querying %+v with param %+v: %w", commit, params, err)
	}

	if commit.Repository.Object.Commit.OID == "" {
		return nil, fmt.Errorf("Commit %q not found", oid)
	}

	return &CommitWithTS{OID: commit.Repository.Object.Commit.OID,
		CommittedDate: commit.Repository.Object.Commit.CommittedDate}, nil
}

// ListCommitsWithCtx returns list of commits, after specific timestamp in the specific branch
func (gc *GitCommits) ListCommitsWithCtx(ctx context.Context, after CommitWithTS) (*Result, error) {
	var full []CommitFragment
	var res Result
	params := map[string]interface{}{
		"owner":         githubv4.String(gc.DefaultOwner),
		"repo":          githubv4.String(gc.DefaultRepo),
		"ref":           githubv4.String("refs/heads/" + gc.DefaultBranch),
		"since":         after.CommittedDate,
		"commitsCursor": (*githubv4.String)(nil),
	}
	for {
		err := gc.client.Query(ctx, &res, params)
		if err != nil {
			return nil, fmt.Errorf("Error querying %w", err)
		}
		full = append(full, res.Repository.Ref.Target.Commit.Nodes...)
		if !res.Repository.Ref.Target.Commit.PageInfo.HasNextPage {
			break
		}
		params["commitsCursor"] = res.Repository.Ref.Target.Commit.PageInfo.EndCursor
	}
	res.Repository.Ref.Target.Commit.Nodes = full
	return &res, nil
}
