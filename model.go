package gitcommits

import "github.com/shurcooL/githubv4"

var viewer struct {
	Viewer struct {
		Login     githubv4.String
		CreatedAt githubv4.DateTime
	}
}

// CommitWithTS represents specific OID and its timestamp
type CommitWithTS struct {
	OID           githubv4.GitObjectID  `json:"oid,omitempty"`
	CommittedDate githubv4.GitTimestamp `json:"committed_date,omitempty"`
}

// query getCommitsByBranchByRepo($owner: String!, $repo: String!, $ref: String!, $since: GitTimestamp!) {
// 	repository(name: $repo, owner: $owner) {
// 	  ref(qualifiedName: $ref) {
// 		target {
// 		  ... on Commit {
// 			history(since: $since, first: 100) {
// 			  nodes {
// 				oid
// 				committedDate
// 				message
// 				author {
// 				  name
// 				  email
// 				}
// 			  }
// 			  totalCount
// 			  pageInfo {
// 				endCursor
// 				hasNextPage
// 			  }
// 			}
// 		  }
// 		}
// 	  }
// 	}
//   }

// var commits struct {
// }

type Result struct {
	Repository struct {
		Ref struct {
			Target struct {
				Commit struct {
					History `graphql:"history(since: $since, first: 100, after: $commitsCursor)" json:"history,omitempty`
				} `graphql:"... on Commit" json:"commit,omitempty"`
			} `json:"target,omitempty"`
		} `graphql:"ref(qualifiedName: $ref)" json:"ref,omitempty"`
	} `graphql:"repository(owner: $owner, name: $repo)" json:"repository,omitempty"`
}

type History struct {
	Nodes      []CommitFragment `json:"nodes,omitempty"`
	TotalCount githubv4.Int     `json:"total_count,omitempty"`
	PageInfo   PageInfo         `json:"page_info,omitempty"`
}

type CommitFragment struct {
	OID           githubv4.GitObjectID  `json:"oid,omitempty"`
	CommittedDate githubv4.GitTimestamp `json:"committed_date,omitempty"`
	Message       githubv4.String       `json:"message,omitempty"`
	Author        Author                `json:"author,omitempty"`
}

type PageInfo struct {
	HasNextPage githubv4.Boolean `json:"has_next_page,omitempty"`
	EndCursor   githubv4.String  `json:"end_cursor,omitempty"`
}

type Author struct {
	Name  githubv4.String `json:"name,omitempty"`
	Email githubv4.String `json:"email,omitempty"`
}
