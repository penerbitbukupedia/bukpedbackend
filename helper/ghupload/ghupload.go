package ghupload

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"

	"github.com/google/go-github/v59/github"

	"golang.org/x/oauth2"
)

// Function to calculate the SHA-256 hash of a file's content
func CalculateHash(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

// Function to upload file to GitHub with hashed filename
func GithubUpload(GitHubAccessToken, GitHubAuthorName, GitHubAuthorEmail string, fileContent []byte, githubOrg string, githubRepo string, pathFile string, replace bool) (content *github.RepositoryContentResponse, response *github.Response, err error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: GitHubAccessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	opts := &github.RepositoryContentFileOptions{
		Message: github.String("Upload file"),
		Content: fileContent,
		Branch:  github.String("main"),
		Author: &github.CommitAuthor{
			Name:  github.String(GitHubAuthorName),
			Email: github.String(GitHubAuthorEmail),
		},
	}

	content, response, err = client.Repositories.CreateFile(ctx, githubOrg, githubRepo, pathFile, opts)
	if (err != nil) && (replace) {
		currentContent, _, _, _ := client.Repositories.GetContents(ctx, githubOrg, githubRepo, pathFile, nil)
		opts.SHA = github.String(currentContent.GetSHA())
		content, response, err = client.Repositories.UpdateFile(ctx, githubOrg, githubRepo, pathFile, opts)
		return
	}

	return
}

// Function to get file content from GitHub repository
// Set header untuk mendownload file
// w.Header().Set("Content-Disposition", "attachment; filename=\"file.ext\"")
// w.Header().Set("Content-Type", "application/octet-stream")
// w.Header().Set("Content-Length", fmt.Sprint(len(fileContent)))
// // Tulis konten file ke response writer
// w.Write(fileContent)
func GithubGetFile(GitHubAccessToken, githubOrg, githubRepo, pathFile string) (fileContent []byte, err error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: GitHubAccessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// Get file content from the repository
	fileContentResponse, _, _, err := client.Repositories.GetContents(ctx, githubOrg, githubRepo, pathFile, nil)
	if err != nil {
		err = errors.New("error GetContents " + err.Error())
		return
	}

	// Decode the base64 encoded file content
	encodedContent, err := fileContentResponse.GetContent()
	if err != nil {
		err = errors.New("error fileContentResponse GetContents " + err.Error())
		return
	}
	fileContent, err = base64.StdEncoding.DecodeString(encodedContent)
	if err != nil {
		err = errors.New("error base64.StdEncoding.DecodeString " + err.Error())
		return
	}

	return
}
