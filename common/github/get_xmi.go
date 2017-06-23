package github

import (
	"golang.org/x/net/context"
	"github.com/google/go-github/github"
	"fmt"
	"os"
	"io"
	"github.com/FINTprosjektet/fint-model/common/utils"
	"log"
	"strings"
	"io/ioutil"
	"gopkg.in/iconv.v1"
)

func GetXMIFile(tag string, force bool) string {

	outFileName := getFilePath(tag)

	if force {
		downloadFile(tag, outFileName)
		cleanFile(outFileName)
	} else if !utils.FileExists(outFileName) {
		downloadFile(tag, outFileName)
		cleanFile(outFileName)
	}

	return outFileName
}

func downloadFile(tag string, outFileName string) {
	client := github.NewClient(nil)
	ctx := context.Background()
	opt := &github.RepositoryContentGetOptions{
		Ref: tag,
	}
	out, err := client.Repositories.DownloadContents(ctx, GITHUB_OWNER, GITHUB_REPO, "FINT-informasjonsmodell.xml", opt)
	if err != nil {
		fmt.Printf("Unable to download XMI file from GitHub: %s", err)
	}
	outFile, err := os.Create(outFileName)
	defer outFile.Close()
	_, err = io.Copy(outFile, out)
	if err != nil {
		fmt.Printf("Unable to write XMI file: %s", err)
	}
}

func getFilePath(tag string) string {
	tmpDir := os.TempDir()
	dir := fmt.Sprintf("%s/fint-model", tmpDir)
	os.Mkdir(dir, 0777)
	outFileName := fmt.Sprintf("%s/%s.xml", dir, tag)

	return outFileName
}

func cleanFile(fileName string) {

	toUtf8(fileName)

	input, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")
	var newLines []string

	newLines = append(newLines, "<?xml version=\"1.0\" encoding=\"utf-8\"?>")
	keep := false
	for i, line := range lines {
		if strings.Contains(line, "<xmi:Extension extender=\"Enterprise Architect\" extenderID=\"6.5\">") {
			keep = true
		}
		if strings.Contains(line, "</xmi:XMI>") {
			keep = false
		}
		if keep {
			l := strings.Replace(lines[i], "uml:", "", -1)
			l = strings.Replace(l, "xmi:", "", -1)
			newLines = append(newLines, l)
		}
	}

	output := strings.Join(newLines, "\n")
	err = ioutil.WriteFile(fileName, []byte(output), 0777)
	if err != nil {
		log.Fatalln(err)
	}

}

func toUtf8(fileName string) {

	cd, err := iconv.Open("utf-8", "windows-1252")
	if err != nil {
		fmt.Println("iconv.Open failed!")
		return
	}
	defer cd.Close()

	input, err := os.Open(fileName)
	bufSize := 0 // default if zero
	r := iconv.NewReader(cd, input, bufSize)

	content, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(fileName, content, 0777)

	if err != nil {
		fmt.Println("\nio.Copy failed:", err)
	}
}

