package tools

import (
	"fmt"
	"log"
	"os"
)

func EmbedResource(fileName, resourceName string) (n int) {
	newResource, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatalln(err)
	}

	resourceFile, err := os.OpenFile("../resources/resources.go", os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModeAppend)
	if err != nil {
		log.Fatalln(err)
	}

	n, err = resourceFile.WriteString(fmt.Sprintf("\nvar %v = %#v", resourceName, newResource))
	if err != nil {
		log.Fatalln(err)
	}
	return
}
