package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/go-git/go-git/v5"
	"github.com/polisko/gitcommits"
)

// ErrRepoNotInTheReverseTree represents situation, when traversing current path there is not any local repo
var ErrRepoNotInTheReverseTree = fmt.Errorf("Git repository was not found in all tree")
var timeFormat = "2006-01-02 15:04:05"

func main() {
	var owner, repo, branch, commit string
	var logLevel int
	var useLocalRepo, shortPrint bool
	var r *git.Repository
	var err error
	var authToken string = "d25b223029bd4ed0dc604039a5d4c613b361a8df"

	// check, if auth token is there
	if a, ok := os.LookupEnv("AUTH_TOKEN"); ok {
		authToken = a
	}

	flag.StringVar(&owner, "o", "golang", "Repository owner")
	flag.StringVar(&repo, "r", "go", "Repository name")
	flag.StringVar(&branch, "b", "master", "Branch")
	flag.BoolVar(&useLocalRepo, "l", true, "Whether to try locate local git repository in the current tree")
	flag.BoolVar(&shortPrint, "s", true, "Whether to print shorter or longer output")
	flag.StringVar(&commit, "c", "", "commit hash (OID)")
	flag.IntVar(&logLevel, "logLevel", 4, "Loglevel, default 4=Info")
	flag.Parse()

	log.SetLevel(log.Level(logLevel))

	gitc, err := gitcommits.NewGitCommits(authToken)
	if err != nil {
		log.Fatalf("%s", err)
	}
	gitc.DefaultOwner = owner
	gitc.DefaultRepo = repo
	gitc.DefaultBranch = branch

	if useLocalRepo {
		r, err = findLocalRepo()
		if err != nil {
			switch err {
			case ErrRepoNotInTheReverseTree:
				pwd, _ := os.Getwd()
				log.Infof("Any git repository was found in the reverse tree of path %q", pwd)
			default:
				log.Errorf("Error: %v", err)
			}
		} else {
			pr, err := r.Head()
			if err != nil {
				log.Errorf("Error: %v", err)
			} else {
				// trying to setup from local HEAD
				commit = pr.Hash().String()
				gitc.DefaultBranch = pr.Name().Short()
				conf, err := r.Config()
				if err == nil {
					if len(conf.Remotes["origin"].URLs) > 0 {
						url := conf.Remotes["origin"].URLs[0]
						switch {
						case url[0:5] == "https":
							p := strings.Split(url, "/")
							gitc.DefaultRepo = p[len(p)-1]
							gitc.DefaultOwner = p[len(p)-2]
						case url[0:4] == "git@":
							p := strings.Split(strings.Split(url, ":")[1], "/")
							gitc.DefaultRepo = p[1][0:strings.Index(p[1], ".")]
							gitc.DefaultOwner = p[0]
						}
					}
				}
			}

		}

	}
	c, err := gitc.FindCommitWithCtx(context.Background(), commit)
	if err != nil {
		log.Fatal(err)
	}
	res, err := gitc.ListCommitsWithCtx(context.Background(), *c)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Owner: %s\nRepository: %s\nBranch: %s\nCommit: %s\nCommit's count: %d\n",
		gitc.DefaultOwner, gitc.DefaultRepo, gitc.DefaultBranch, commit, res.Repository.Ref.Target.Commit.TotalCount)
	fmt.Printf("%s\n", strings.Repeat("*", 20))
	for _, v := range res.Repository.Ref.Target.Commit.Nodes {
		if shortPrint {
			printShort(v)
		} else {
			printFull(v)

		}

	}
}

func printFull(v gitcommits.CommitFragment) {
	fmt.Printf("%s %v %s\n", v.OID[0:7], v.CommittedDate.Format(timeFormat), v.Author.Name)
	fmt.Printf("%s\n", v.Message)
	fmt.Printf("%s\n\n", strings.Repeat("*", 20))
}
func printShort(v gitcommits.CommitFragment) {
	fmt.Printf("%s %v %s\n", v.OID[0:7], v.CommittedDate.Format(timeFormat), v.Author.Name)
}
func findLocalRepo() (*git.Repository, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("Error getting `pwd`: %w", err)
	}
	p := strings.Count(dir, string(os.PathSeparator))
	for c := 0; c <= p; c++ {
		r, err := git.PlainOpen(dir)
		if err != nil {
			log.Debugf("Git repository not found in %s (%v)", dir, err)
		} else {
			return r, nil
		}
		i := strings.LastIndex(dir, string(os.PathSeparator))
		if c == p-1 {
			i = strings.Index(dir, string(os.PathSeparator))
			dir = dir[0 : i+1]
		} else {
			dir = dir[0:i]
		}
	}
	return nil, ErrRepoNotInTheReverseTree
}
