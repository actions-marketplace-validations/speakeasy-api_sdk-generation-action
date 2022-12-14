package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/v48/github"
	"golang.org/x/oauth2"
)

func createRelease(sdkVersion, commitHash, openAPIDocPath, openAPIDocVersion, speakeasyCLIVersion, accessToken string) error {
	fmt.Println("Creating release")

	repoPath := os.Getenv("GITHUB_REPOSITORY")
	parts := strings.Split(repoPath, "/")

	tag := "v" + sdkVersion

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	_, _, err := github.NewClient(tc).Repositories.CreateRelease(context.Background(), os.Getenv("GITHUB_REPOSITORY_OWNER"), parts[len(parts)-1], &github.RepositoryRelease{
		TagName:         github.String(tag),
		TargetCommitish: github.String(commitHash),
		Name:            github.String(fmt.Sprintf("%s - %s", tag, time.Now().Format("2006-01-02 15:04:05"))),
		Body: github.String(fmt.Sprintf(`# Generated by Speakeasy CLI
Based on:
- OpenAPI Doc %s %s
- Speakeasy CLI %s https://github.com/speakeasy-api/speakeasy`, openAPIDocVersion, openAPIDocPath, speakeasyCLIVersion)),
	})
	if err != nil {
		return fmt.Errorf("failed to create release: %w", err)
	}

	return nil
}
