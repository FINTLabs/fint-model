package generate

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/FINTprosjektet/fint-model/common/config"
	"github.com/FINTprosjektet/fint-model/common/document"
	"github.com/FINTprosjektet/fint-model/common/github"
	"github.com/FINTprosjektet/fint-model/common/parser"
	"github.com/FINTprosjektet/fint-model/common/types"
	"github.com/FINTprosjektet/fint-model/namespaces"
	"github.com/FINTprosjektet/fint-model/packages"
	"github.com/codegangsta/cli"
)

func CmdGenerate(c *cli.Context) {

	var tag string
	if c.GlobalString("tag") == config.DEFAULT_TAG {
		tag = github.GetLatest(c.GlobalString("owner"), c.GlobalString("repo"))
	} else {
		tag = c.GlobalString("tag")
	}
	force := c.GlobalBool("force")
	owner := c.GlobalString("owner")
	repo := c.GlobalString("repo")
	filename := c.GlobalString("filename")

	if c.String("lang") == "JAVA" {
		generateJavaCode(owner, repo, tag, filename, force)
	}

	if c.String("lang") == "CS" {
		generateCSCode(owner, repo, tag, filename, force)
	}

	if c.String("lang") == "ALL" {
		generateCSCode(owner, repo, tag, filename, force)
		generateJavaCode(owner, repo, tag, filename, force)
	}
}

func generateJavaCode(owner string, repo string, tag string, filename string, force bool) {

	document.Get(owner, repo, tag, filename, force)
	fmt.Println("Generating Java code:")
	setupJavaDirStructure(owner, repo, tag, filename, force)
	classes, _, packageClassMap, _ := parser.GetClasses(owner, repo, tag, filename, force)
	for _, c := range classes {
		fmt.Printf("  > Creating class: %s.java\n", c.Name)
		class := GetJavaClass(c)

		path := fmt.Sprintf("%s/%s/%s.java", config.JAVA_BASE_PATH, strings.Replace(c.Package, ".", "/", -1), c.Name)
		err := ioutil.WriteFile(removeJavaPackagePathFromFilePath(path), []byte(class), 0777)
		if err != nil {
			fmt.Printf("Unable to write file: %s", err)
		}

	}

	for p, cl := range packageClassMap {
		action := getAction(p, cl, tag)
		fmt.Printf("  > Creating action: %s.java\n", action.Name)
		actionEnum := GetJavaActionEnum(action)
		path := fmt.Sprintf("%s/%s/%s.java", config.JAVA_BASE_PATH, strings.Replace(p, ".", "/", -1), action.Name)
		err := ioutil.WriteFile(removeJavaPackagePathFromFilePath(path), []byte(actionEnum), 0777)
		if err != nil {
			fmt.Printf("Unable to write file: %s", err)
		}

	}

	fmt.Println("Finish generating Java code!")
}

func removeJavaPackagePathFromFilePath(path string) string {
	return strings.Replace(path, "no/fint/model/", "", -1)
}

func getAction(p string, cl []types.Class, tag string) types.Action {
	var action types.Action

	packageList := strings.Split(p, ".")
	pkg := packageList[len(packageList)-1]
	action.Name = fmt.Sprintf("%sActions", strings.Title(pkg))

	action.Package = p
	action.Namespace = p
	action.GitTag = tag

	for _, c := range cl {
		if c.Identifiable && !c.Abstract {
			action.Classes = append(action.Classes, strings.ToUpper(c.Name))
		}
	}
	return action
}

func generateCSCode(owner string, repo string, tag string, filename string, force bool) {

	document.Get(owner, repo, tag, filename, force)
	fmt.Println("Generating CSharp code:")
	setupCSDirStructure(owner, repo, tag, filename, force)
	classes, _, _, packageClassMap := parser.GetClasses(owner, repo, tag, filename, force)
	for _, c := range classes {
		fmt.Printf("  > Creating class: %s.cs\n", c.Name)

		class := GetCSClass(c)

		path := fmt.Sprintf("%s/%s.cs", getCSPath(c.Namespace), c.Name)
		err := ioutil.WriteFile(path, []byte(class), 0777)
		if err != nil {
			fmt.Printf("Unable to write file: %s", err)
		}

	}

	for p, cl := range packageClassMap {
		action := getAction(p, cl, tag)
		fmt.Printf("  > Creating action: %s.cs\n", action.Name)
		actionEnum := GetCSActionEnum(action)
		path := fmt.Sprintf("%s/%s.cs", getCSPath(p), action.Name)
		err := ioutil.WriteFile(path, []byte(actionEnum), 0777)
		if err != nil {
			fmt.Printf("Unable to write file: %s", err)
		}

	}

	fmt.Println("Finish generating CSharp code!")

}

func setupCSDirStructure(owner string, repo string, tag string, filename string, force bool) {
	fmt.Println("  > Setup directory structure.")
	os.RemoveAll("net")
	err := os.MkdirAll(config.CS_BASE_PATH, 0777)
	if err != nil {
		fmt.Println("Unable to create base structure")
		fmt.Println(err)
	}
	for _, ns := range namespaces.DistinctNamespaceList(owner, repo, tag, filename, force) {
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

func setupJavaDirStructure(owner string, repo string, tag string, filename string, force bool) {
	fmt.Println("  > Setup directory structure.")
	os.RemoveAll("java")
	err := os.MkdirAll(config.JAVA_BASE_PATH, 0777)
	if err != nil {
		fmt.Println("Unable to create base structure")
		fmt.Println(err)
	}

	for _, pkg := range packages.DistinctPackageList(owner, repo, tag, filename, force) {
		path := fmt.Sprintf("%s/%s", config.JAVA_BASE_PATH, strings.Replace(pkg, ".", "/", -1))
		err := os.MkdirAll(removeJavaPackagePathFromFilePath(path), 0777)
		if err != nil {
			fmt.Println("Unable to create packages structure")
			fmt.Println(err)
		}
	}

}
