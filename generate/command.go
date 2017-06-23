package generate

import (
	"github.com/codegangsta/cli"
	"os"
	"fmt"
	"github.com/FINTprosjektet/fint-model/packages"
	"strings"
	"github.com/FINTprosjektet/fint-model/common/github"
	"github.com/FINTprosjektet/fint-model/common/document"
	"github.com/FINTprosjektet/fint-model/common/parser"
	"github.com/FINTprosjektet/fint-model/namespaces"
	"io/ioutil"
	"github.com/FINTprosjektet/fint-model/common/config"
)

func CmdGenerate(c *cli.Context) {

	var tag string
	if c.GlobalString("tag") == config.DEFAULT_TAG {
		tag = github.GetLatest()
	} else {
		tag = c.GlobalString("tag")
	}

	if c.String("lang") == "JAVA" {
		generateJavaCode(tag)
	}

	if c.String("lang") == "CS" {
		generateCSCode(tag)
	}

}

func generateJavaCode(tag string) {

	//document.GetFile(tag)
	document.Get(tag)
	fmt.Println("Generating Java code:")
	setupJavaDirStructure(tag)
	classes, _ := parser.GetClasses(tag)
	for _, c := range classes {
		fmt.Printf("  > Creating class: %s.java\n", c.Name)
		class := GetJavaClass(c)

		path := fmt.Sprintf("%s/%s/%s.java", config.JAVA_BASE_PATH, strings.Replace(c.Package, ".", "/", -1), c.Name)
		err := ioutil.WriteFile(path, []byte(class), 0777)
		if err != nil {
			fmt.Printf("Unable to write file: %s", err)
		}

	}

	fmt.Println("Finish generating Java code!")
}

func generateCSCode(tag string) {

	document.Get(tag)
	fmt.Println("Generating CSharp code:")
	setupCSDirStructure(tag)
	classes, _ := parser.GetClasses(tag)
	for _, c := range classes {
		fmt.Printf("  > Creating class: %s.cs\n", c.Name)

		class := GetCSClass(c)

		path := fmt.Sprintf("%s/%s.cs", getCSPath(c.Namespace), c.Name)
		err := ioutil.WriteFile(path, []byte(class), 0777)
		if err != nil {
			fmt.Printf("Unable to write file: %s", err)
		}

	}

	fmt.Println("Finish generating CSharp code!")

}

func setupCSDirStructure(tag string) {
	fmt.Println("  > Setup directory structure.")
	os.RemoveAll("net")
	err := os.MkdirAll(config.CS_BASE_PATH, 0777)
	if err != nil {
		fmt.Println("Unable to create base structure")
		fmt.Println(err)
	}
	for _, ns := range namespaces.DistinctNamespaceList(tag) {
		path := getCSPath(ns)
		err := os.MkdirAll(path, 0777)
		if err != nil {
			fmt.Println("Unable to create namespace structure")
			fmt.Println(err)
		}
	}
}
func getCSPath(ns string) string {
	nsList := strings.Split(ns, ".")
	projectDir := fmt.Sprintf("%s.%s.%s", nsList[0], nsList[1], nsList[2])
	subDirs := strings.Replace(ns, projectDir, "", -1)
	subDirs = strings.Replace(subDirs, ".", "/", -1)
	path := fmt.Sprintf("%s/%s/%s", config.CS_BASE_PATH, projectDir, subDirs)
	return path
}

func setupJavaDirStructure(tag string) {
	fmt.Println("  > Setup directory structure.")
	os.RemoveAll("java")
	err := os.MkdirAll(config.JAVA_BASE_PATH, 0777)
	if err != nil {
		fmt.Println("Unable to create base structure")
		fmt.Println(err)
	}
	for _, pkg := range packages.DistinctPackageList(tag) {
		path := fmt.Sprintf("%s/%s", config.JAVA_BASE_PATH, strings.Replace(pkg, ".", "/", -1))
		err := os.MkdirAll(path, 0777)
		if err != nil {
			fmt.Println("Unable to create packages structure")
			fmt.Println(err)
		}

	}
}
