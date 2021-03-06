package release

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/google/go-github/v32/github"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
	"golang.org/x/xerrors"
)

type githubOpt struct {
	Version    string
	From       string
	Attach     []string
	GithubRepo string
	BodyFile   string
}

func githubRelease(opt *githubOpt) error {
	token := os.Getenv("GITHUB_APITOKEN")
	if token == "" {
		return xerrors.New("GITHUB_APITOKEN is mandatory")
	}
	if !strings.Contains(opt.GithubRepo, "/") {
		return xerrors.Errorf("invalid repo name: %s", opt.GithubRepo)
	}
	ver, err := semver.NewVersion(opt.Version)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}
	preRelease := false
	if ver.Prerelease() != "" {
		preRelease = true
	}

	body := ""
	if _, err := os.Lstat(opt.BodyFile); !os.IsNotExist(err) {
		b, err := ioutil.ReadFile(opt.BodyFile)
		if err != nil {
			return xerrors.Errorf(": %w", err)
		}
		body = string(b)
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(context.Background(), ts)
	client := github.NewClient(tc)

	r := strings.Split(opt.GithubRepo, "/")
	owner, repo := r[0], r[1]

	release, res, err := client.Repositories.GetReleaseByTag(context.Background(), owner, repo, opt.Version)
	if err != nil && res == nil {
		return xerrors.Errorf(": %w", err)
	}

	// Create new release
	attachedFiles := make(map[string]struct{})
	if release != nil {
		for _, v := range release.Assets {
			attachedFiles[v.GetName()] = struct{}{}
		}
	} else {
		branch, _, err := client.Repositories.GetBranch(context.Background(), owner, repo, opt.From)
		if err != nil {
			return xerrors.Errorf(": %w", err)
		}
		if branch == nil {
			return xerrors.Errorf("branch(%s) is not found", opt.From)
		}
		fmt.Printf("Get commit hash %s\n", branch.Commit.GetSHA())

		r, _, err := client.Repositories.CreateRelease(context.Background(), owner, repo, &github.RepositoryRelease{
			TagName:         github.String(opt.Version),
			TargetCommitish: branch.Commit.SHA,
			Body:            github.String(body),
			Prerelease:      github.Bool(preRelease),
		})
		if err != nil {
			return xerrors.Errorf(": %w", err)
		}
		fmt.Printf("Created release id=%d TagName=%s\n", r.GetID(), r.GetTagName())
		release = r
	}

	for _, v := range opt.Attach {
		if _, err := os.Stat(v); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "%s is not found", v)
			continue
		}
		f, err := os.Open(v)
		if err != nil {
			return xerrors.Errorf(": %w", err)
		}

		filename := filepath.Base(f.Name())
		if _, ok := attachedFiles[filename]; ok {
			fmt.Printf("%s is already exist. skip uploading this file", filename)
			continue
		}

		assets, _, err := client.Repositories.UploadReleaseAsset(context.Background(), owner, repo, release.GetID(), &github.UploadOptions{Name: filename}, f)
		if err != nil {
			return xerrors.Errorf(": %w", err)
		}
		if assets == nil {
			fmt.Fprintf(os.Stderr, "Failed upload an asset")
		}
		fmt.Printf("Upload asset %s\n", filename)
	}

	return nil
}

func GitHub(rootCmd *cobra.Command) {
	opt := githubOpt{}

	ghRelease := &cobra.Command{
		Use:   "github",
		Short: "Create GitHub Release",
		RunE: func(_ *cobra.Command, _ []string) error {
			return githubRelease(&opt)
		},
	}
	ghRelease.Flags().StringVar(&opt.Version, "version", "", "")
	ghRelease.Flags().StringVar(&opt.From, "from", "", "")
	ghRelease.Flags().StringArrayVar(&opt.Attach, "attach", []string{}, "")
	ghRelease.Flags().StringVar(&opt.GithubRepo, "repo", "", "")
	ghRelease.Flags().StringVar(&opt.BodyFile, "body", "", "Release body")

	rootCmd.AddCommand(ghRelease)
}
