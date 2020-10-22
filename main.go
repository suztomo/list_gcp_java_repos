package main

import (
	"context"
	"os"
	"strings"

	"github.com/google/go-github/github" // with go modules disabled
	"golang.org/x/oauth2"
)

func main() {
	ctx := context.Background()

	githubToken := os.Getenv("GITHUB_TOKEN")
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	orgs := []string{
		"cloudendpoints",
		"androidthings", "androidx", "aosp-mirror", "apigee", "bumptech",
		"cdapio", "cloudspannerecosystem", "dialogflow", "firebase",
		"firebaseextended", "google",
		"googleapis", "google-business-communications", "google-cloudsearch",
		"google-pay", "googleads", "GoogleCloudDataproc",
		"GoogleCloudPlatform",
		"googlecontainertools", "googlegsa", "grafeas", "grpc-ecosystem",
		"jspecify", "protocolbuffers", "rosjava", "stackdriver",
	}

	for _, org := range orgs {
		// list all repositories for the authenticated user
		opt := &github.RepositoryListByOrgOptions{Type: "public"}
		for {
			repos, resp, err := client.Repositories.ListByOrg(ctx, org, opt)
			if err != nil {
				println("Failed to list repositories")
				return
			}

			for _, repo := range repos {
				repoUrl := repo.URL
				lang := repo.GetLanguage()
				description := ""
				createdAt := repo.CreatedAt
				updatedAt := repo.UpdatedAt
				pushedAt := repo.PushedAt

				if repo.Description != nil {
					description = *repo.Description
				}

				if lang != "Java" {
					continue
				}

				elements := []string{org, *repo.Name, *repoUrl, lang, description, createdAt.Format("2006-01-02"),
					updatedAt.Format("2006-01-02"), pushedAt.Format("2006-01-02")}
				println(strings.Join(elements, "|"))
			}

			if resp.NextPage == 0 {
				break
			}
			opt.Page = resp.NextPage
		}

	}
}
