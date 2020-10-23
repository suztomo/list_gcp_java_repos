package main

import (
	"context"
	"encoding/json"
	"net/http"
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
				//				repoURL := repo.URL
				name := *repo.Name
				lang := repo.GetLanguage()
				description := ""
				//				createdAt := repo.CreatedAt
				//				updatedAt := repo.UpdatedAt
				//				pushedAt := repo.PushedAt

				if repo.Description != nil {
					description = *repo.Description
				}

				exclusionWords := []string{"sample", "example", "Sample", "Example"}
				foundWord := false
				for _, word := range exclusionWords {
					if strings.Contains(name, word) || strings.Contains(description, word) {
						foundWord = true
					}
				}
				if foundWord {
					continue
				}

				if lang != "Java" {
					continue
				}

				if strings.Contains(description, "|") {
					println("Description contains selarator")
					os.Exit(1)

				}

				listContributors(tc, org, repo)

				// elements := []string{org, name, *repoURL, lang, description}
				//println(strings.Join(elements, "|"))
			}

			if resp.NextPage == 0 {
				break
			}
			opt.Page = resp.NextPage
		}

	}
}

func listContributors(client *http.Client, organizationName string, repo *github.Repository) error {

	response, err := client.Get(*repo.ContributorsURL)
	if err != nil {
		return err
	}

	defer response.Body.Close()
	var contributors []github.User

	json.NewDecoder(response.Body).Decode(&contributors)

	description := ""
	if repo.Description != nil {
		description = *repo.Description
	}

	for _, user := range contributors {
		elements := []string{
			organizationName,
			*repo.Name,
			*repo.URL,
			*repo.Language,
			description,
			*user.Login}
		println(strings.Join(elements, "|"))
	}

	return nil
}
