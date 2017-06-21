package generate

import (
	"github.com/codegangsta/cli"
	"os"
	"fmt"
	"github.com/FINTprosjektet/fint-model/packages"
	"strings"
	"github.com/FINTprosjektet/fint-model/common/github"
	"github.com/FINTprosjektet/fint-model/common/document"
	"io/ioutil"
	"github.com/FINTprosjektet/fint-model/common/parser"
	"github.com/FINTprosjektet/fint-model/namespaces"
)

const javaBasePath = "java/src/main/java/"
const csharpBasePath = "net/"

func CmdGenerate(c *cli.Context) {

	var tag string
	if c.GlobalString("tag") == "latest" {
		tag = github.GetLatest()
	} else {
		tag = c.GlobalString("tag")
	}

	if c.String("lang") == "JAVA" {
		generateJavaCode(tag)
	}

	if c.String("lang") == "NET" {
		generateNetCode(tag)
	}

}

func generateJavaCode(tag string) {

	//document.GetFile(tag)
	document.Get(tag)
	fmt.Println("Generating Java code:")
	setupJavaDirStructure(tag)
	classes, impMap := parser.GetClasses(tag)
	for _, c := range classes {
		fmt.Printf("  > Creating class: %s.java\n", c.Name)
		var class string

		if len(c.Extends) > 0 && c.Abstract {
			class = GetAbstractExtendedJavaClass(c, impMap)
		} else if len(c.Extends) > 0 && c.Identifiable {
			class = GetExtendedJavaClassIdentifiable(c, impMap)
		} else if c.Identifiable {
			class = GetJavaClassIdentifiable(c, impMap)
		} else if len(c.Extends) > 0 {
			class = GetExtendedJavaClass(c, impMap)
		} else if c.Abstract {
			class = GetAbstractJavaClass(c, impMap)
		} else {
			class = GetJavaClass(c, impMap)
		}

		path := fmt.Sprintf("%s/%s/%s.java", javaBasePath, strings.Replace(c.Package, ".", "/", -1), c.Name)
		err := ioutil.WriteFile(path, []byte(class), 0777)
		if err != nil {
			fmt.Printf("Unable to write file: %s", err)
		}
	}

	fmt.Println("Finish generating Java code!")
}

func generateNetCode(tag string) {

	document.Get(tag)
	fmt.Println("Generating CSharp code:")
	setupCSharpDirStructure(tag)
	classes, impMap := parser.GetClasses(tag)
	for _, c := range classes {
		fmt.Printf("  > Creating class: %s.cs\n", c.Name)
		var class string

		if len(c.Extends) > 0 && c.Abstract {
			class = GetExtendedAbstractCSharpClass(c, impMap)
		} else if len(c.Extends) > 0 {
			class = GetExtendedCSharpClass(c, impMap)
		} else if c.Abstract {
			class = GetAbstractCSharpClass(c, impMap)
		} else {
			class = GetCSharpClass(c, impMap)
		}

		path := fmt.Sprintf("%s/%s.cs", getCSharpPath(c.Namespace), c.Name)
		err := ioutil.WriteFile(path, []byte(class), 0777)
		if err != nil {
			fmt.Printf("Unable to write file: %s", err)
		}
	}

	fmt.Println("Finish generating CSharp code!")

}

func setupCSharpDirStructure(tag string) {
	fmt.Println("  > Setup directory structure.")
	os.RemoveAll("net")
	err := os.MkdirAll(csharpBasePath, 0777)
	if err != nil {
		fmt.Println("Unable to create base structure")
		fmt.Println(err)
	}
	for _, ns := range namespaces.DistinctNamespaceList(tag) {
		path := getCSharpPath(ns)
		err := os.MkdirAll(path, 0777)
		if err != nil {
			fmt.Println("Unable to create namespace structure")
			fmt.Println(err)
		}
	}
}
func getCSharpPath(ns string) string {
	nsList := strings.Split(ns, ".")
	projectDir := fmt.Sprintf("%s.%s.%s", nsList[0], nsList[1], nsList[2])
	subDirs := strings.Replace(ns, projectDir, "", -1)
	subDirs = strings.Replace(subDirs, ".", "/", -1)
	path := fmt.Sprintf("%s/%s/%s", csharpBasePath, projectDir, subDirs)
	return path
}

func setupJavaDirStructure(tag string) {
	fmt.Println("  > Setup directory structure.")
	os.RemoveAll("java")
	err := os.MkdirAll(javaBasePath, 0777)
	if err != nil {
		fmt.Println("Unable to create base structure")
		fmt.Println(err)
	}
	for _, pkg := range packages.DistinctPackageList(tag) {
		path := fmt.Sprintf("%s/%s", javaBasePath, strings.Replace(pkg, ".", "/", -1))
		err := os.MkdirAll(path, 0777)
		if err != nil {
			fmt.Println("Unable to create packages structure")
			fmt.Println(err)
		}

	}
}
