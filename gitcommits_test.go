package gitcommits

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/shurcooL/githubv4"
)

// Tests are using production GitHub

func TestGitCommits_FindCommitWithCtx(t *testing.T) {
	// Needs authToken to be set in ENV variables
	authToken := os.Getenv("AUTH_TOKEN")
	if authToken == "" {
		fmt.Println("Please set up 'AUTH_TOKEN' to run the tests.")
		t.FailNow()
		return
	}
	cli, err := NewGitCommits("")
	if err != ErrMissingOrBadAuthToken {
		t.Errorf("Must fail due to missing authToken")
	}
	cli, err = NewGitCommits(authToken)
	if err != nil {
		t.Errorf("Must not fail or check if authToken is valid")
	}
	cli.DefaultBranch = "master"
	cli.DefaultOwner = "golang"
	cli.DefaultRepo = "go"
	type args struct {
		ctx context.Context
		oid string
	}
	tests := []struct {
		name    string
		args    args
		want    *CommitWithTS
		wantErr bool
	}{
		{name: "Commit does exist", args: args{ctx: context.Background(), oid: "cda8ee0"},
			want: &CommitWithTS{OID: "cda8ee095e487951eab5a53a097e2b8f400f237d",
				CommittedDate: githubv4.GitTimestamp{Time: time.Date(2021, 02, 26, 20, 49, 57, 0, time.UTC)}}, wantErr: false},
		{name: "Commit does NOT exist", args: args{ctx: context.Background(), oid: "NONEXISTED"},
			want: nil, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gc := cli
			got, err := gc.FindCommitWithCtx(tt.args.ctx, tt.args.oid)
			if (err != nil) != tt.wantErr {
				t.Errorf("GitCommits.FindCommitWithCtx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GitCommits.FindCommitWithCtx() = %v, want %v", got, tt.want)
			}
		})
	}
}
