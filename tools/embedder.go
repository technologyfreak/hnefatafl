package tools

import (
	"fmt"
	"os"
)

func EmbedResource(fileName, resourceName string) (err error) {
	newResource, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	resourceFile, err := os.OpenFile("../resources/resources.go", os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModeAppend)
	if err != nil {
		return err
	}

	_, err = resourceFile.WriteString(fmt.Sprintf("\nvar %v = %#v", resourceName, newResource))
	if err != nil {
		return err
	}
	return
}
