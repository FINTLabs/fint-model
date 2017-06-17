package document

import (
	"os"
	"fmt"
	"github.com/FINTprosjektet/fint-model/github"
	"gopkg.in/iconv.v1"
	"io/ioutil"
	"log"
	"strings"
	"github.com/antchfx/xquery/xml"
)

func Get(tag string) *xmlquery.Node {

	fileName := github.GetXMIFile(tag)
	cleanFile(fileName)
	defer os.Remove(fileName)
	//fileName := "feature-141087757.xml"

	f, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
	}
	doc, err := xmlquery.Parse(f)
	if err != nil {
		fmt.Println(err)
	}
	return doc

}

func GetFile(tag string) {
	fileName := github.GetXMIFile(tag)
	cleanFile(fileName)
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
