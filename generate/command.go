package generate

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/FINTLabs/fint-model/common/config"
	"github.com/FINTLabs/fint-model/common/document"
	"github.com/FINTLabs/fint-model/common/github"
	"github.com/FINTLabs/fint-model/common/parser"
	"github.com/FINTLabs/fint-model/common/types"
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

	resource := c.Bool("resource")

	if c.String("lang") == "JAVA" {
		generateJavaCode(owner, repo, tag, filename, force, resource)
	}

	if c.String("lang") == "CS" {
		generateCSCode(owner, repo, tag, filename, force, resource)
	}

	if c.String("lang") == "ALL" {
		generateCSCode(owner, repo, tag, filename, force, resource)
		generateJavaCode(owner, repo, tag, filename, force, resource)
	}
}

func writeFile(path string, filename string, content []byte) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0777)
		if err != nil {
			return err
		}
	}
	return ioutil.WriteFile(path+"/"+filename, content, 0777)
}

func writeJavaClass(pkg string, class string, content []byte) error {
	path := fmt.Sprintf("%s/%s", config.JAVA_BASE_PATH, strings.Replace(pkg, ".", "/", -1))
	return writeFile(removeJavaPackagePathFromFilePath(path), class+".java", []byte(content))
}

func generateJavaCode(owner string, repo string, tag string, filename string, force bool, resource bool) {

	document.Get(owner, repo, tag, filename, force)
	fmt.Println("Generating Java code:")
	setupJavaDirStructure(owner, repo, tag, filename, force)
	classes, _, packageClassMap, _ := parser.GetClasses(owner, repo, tag, filename, force)
	for _, c := range classes {
		if resource {
			if c.Resource || len(c.Resources) > 0 || c.Identifiable {
				fmt.Printf("  > Creating resource class: %sResource.java\n", c.Name)
				class := GetJavaResourceClass(c)
				pkg := strings.Replace(c.Package, ".model.", ".model.resource.", -1)
				err := writeJavaClass(pkg, c.Name+"Resource", []byte(class))
				if err != nil {
					fmt.Printf("Unable to write file: %s", err)
				}

				fmt.Printf("  > Creating resources class: %sResources.java\n", c.Name)
				class = GetJavaResourcesClass(c)
				err = writeJavaClass(pkg, c.Name+"Resources", []byte(class))
				if err != nil {
					fmt.Printf("Unable to write file: %s", err)
				}
			}
		}

		fmt.Printf("  > Creating class: %s.java\n", c.Name)
		class := GetJavaClass(c)
		err := writeJavaClass(c.Package, c.Name, []byte(class))
		if err != nil {
			fmt.Printf("Unable to write file: %s", err)
		}
	}

	for p, cl := range packageClassMap {
		action := getAction(p, cl, tag)
		fmt.Printf("  > Creating action: %s.java\n", action.Name)
		actionEnum := GetJavaActionEnum(action)
		path := fmt.Sprintf("%s/%s", config.JAVA_BASE_PATH, strings.Replace(p, ".", "/", -1))
		err := writeFile(removeJavaPackagePathFromFilePath(path), action.Name+".java", []byte(actionEnum))
		if err != nil {
			fmt.Printf("Unable to write file: %s", err)
		}

	}

	fmt.Println("Finish generating Java code!")
}

func removeJavaPackagePathFromFilePath(path string) string {
	return strings.Replace(path, "no/fint/model/", "", -1)
}

func getAction(p string, cl []*types.Class, tag string) types.Action {
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

func writeCSClass(namespace string, class string, content []byte, resource bool) error {
	return writeFile(getCSPath(namespace, resource), class+".cs", []byte(content))
}

func generateCSCode(owner string, repo string, tag string, filename string, force bool, resource bool) {

	document.Get(owner, repo, tag, filename, force)
	fmt.Println("Generating CSharp code:")
	setupCSDirStructure(owner, repo, tag, filename, force)
	classes, _, _, packageClassMap := parser.GetClasses(owner, repo, tag, filename, force)
	for _, c := range classes {

		if resource {
			if c.Resource || len(c.Resources) > 0 {
				fmt.Printf("  > Creating resource class: %sResource.cs\n", c.Name)
				class := GetCSResourceClass(c)
				err := writeCSClass(c.Namespace, c.Name+"Resource", []byte(class), true)
				if err != nil {
					fmt.Printf("Unable to write file: %s", err)
				}

				fmt.Printf("  > Creating resources class: %sResources.cs\n", c.Name)
				class = GetCSResourcesClass(c)
				err = writeCSClass(c.Namespace, c.Name+"Resources", []byte(class), true)
				if err != nil {
					fmt.Printf("Unable to write file: %s", err)
				}
			}
		}

		fmt.Printf("  > Creating class: %s.cs\n", c.Name)

		class := GetCSClass(c)

		err := writeCSClass(c.Namespace, c.Name, []byte(class), false)
		if err != nil {
			fmt.Printf("Unable to write file: %s", err)
		}
	}

	for p, cl := range packageClassMap {
		action := getAction(p, cl, tag)
		fmt.Printf("  > Creating action: %s.cs\n", action.Name)
		actionEnum := GetCSActionEnum(action)
		err := writeCSClass(p, action.Name, []byte(actionEnum), false)
		if err != nil {
			fmt.Printf("Unable to write file: %s", err)
		}
	}

	fmt.Println("Finish generating CSharp code!")

}

func setupCSDirStructure(owner string, repo string, tag string, filename string, force bool) {
	fmt.Println("  > Setup directory structure.")
	os.RemoveAll(config.CS_BASE_PATH)
	err := os.MkdirAll(config.CS_BASE_PATH, 0777)
	if err != nil {
		fmt.Println("Unable to create base structure")
		fmt.Println(err)
	}

	/*
		for _, ns := range namespaces.DistinctNamespaceList(owner, repo, tag, filename, force) {
			path := getCSPath(ns)
			err := os.MkdirAll(path, 0777)
			if err != nil {
				fmt.Println("Unable to create namespace structure")
				fmt.Println(err)
			}
		}
	*/
}
func getCSPath(ns string, resource bool) string {
	base := getCSBase(resource)
	nsList := strings.Split(ns, ".")
	projectDir := fmt.Sprintf("%s.%s.%s", nsList[0], nsList[1], nsList[2])
	subDirs := strings.Replace(ns, projectDir, "", -1)
	subDirs = strings.Replace(subDirs, ".", "/", -1)
	path := fmt.Sprintf("%s/%s/%s", base, projectDir, subDirs)
	return path
}

func getCSBase(resource bool) string {
	if resource {
		return config.CS_BASE_PATH + "/resource"
	}
	return config.CS_BASE_PATH
}

func setupJavaDirStructure(owner string, repo string, tag string, filename string, force bool) {
	fmt.Println("  > Setup directory structure.")
	os.RemoveAll(config.JAVA_BASE_PATH)
	err := os.MkdirAll(config.JAVA_BASE_PATH, 0777)
	if err != nil {
		fmt.Println("Unable to create base structure")
		fmt.Println(err)
	}

	/*
		for _, pkg := range packages.DistinctPackageList(owner, repo, tag, filename, force) {
			path := fmt.Sprintf("%s/%s", config.JAVA_BASE_PATH, strings.Replace(pkg, ".", "/", -1))
			err := os.MkdirAll(removeJavaPackagePathFromFilePath(path), 0777)
			if err != nil {
				fmt.Println("Unable to create packages structure")
				fmt.Println(err)
			}
		}
	*/
}
