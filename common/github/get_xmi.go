package github

import (
	"golang.org/x/net/context"
	"github.com/google/go-github/github"
	"fmt"
	"os"
	"io"
)

func GetXMIFile(ref string) string {
	client := github.NewClient(nil)
	ctx := context.Background()

	opt := &github.RepositoryContentGetOptions{
		Ref: ref,
	}
	out, err := client.Repositories.DownloadContents(ctx, GITHUB_OWNER, GITHUB_REPO, "FINT-informasjonsmodell.xml", opt)

	if err != nil {
		fmt.Printf("Unable to download XMI file from GitHub: %s", err)
	}

	outFileName := fmt.Sprintf("%s.xml", ref)
	outFile, err := os.Create(outFileName)
	// handle err
	defer outFile.Close()
	_, err = io.Copy(outFile, out)

	if err != nil {
		fmt.Printf("Unable to write XMI file: %s", err)
	}

	return outFileName
}
