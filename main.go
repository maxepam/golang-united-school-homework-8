package homework8

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

type Arguments map[string]string

type Item struct {
	Id  string			`json:"id"`
	Email string	  `json:"email"`
	Age int			`json:"age"`
}

func Perform(args Arguments, writer io.Writer) error {
	operation, ok := args["operation"]
	if (!ok || len(operation) == 0) {
		return fmt.Errorf("-operation flag has to be specified")
	}
	if (!isAllowedOperation(operation)) {
		return fmt.Errorf("Operation %s not allowed!", operation)
	}
	
	filename, ok := args["fileName"]
	if (!ok || len(filename) == 0) {
		return fmt.Errorf("-fileName flag has to be specified")
	}
	
	if (operation == "list") {
		return listItems(filename, writer)
	}

	if (operation == "add") {
		item, ok := args["item"]
		if (!ok || len(item) == 0) {
			return fmt.Errorf("-item flag has to be specified")
		}
		return addItem(filename, item, writer)
	}

	id, ok := args["id"]
	if (!ok || len(id) == 0) {
		return fmt.Errorf("-id flag has to be specified")
	}

	if (operation == "remove") {
		return removeItem(filename, id, writer)
	} else {
		return findItem(filename, id, writer)
	}
}

//this is a handler for operation=list
func listItems(fileName string, writer io.Writer) error {
	fileContent, err := getFileContent(fileName)
	if err != nil {
		return err
	}
	writer.Write(fileContent)
	return nil
}

//this is a handler for operation=add
func addItem(fileName string, item string, writer io.Writer) error {
	itemByte := []byte(item)
	if !json.Valid(itemByte) {
		return fmt.Errorf("Item is not valid json object")
	}
	var newItem Item
	err := json.Unmarshal(itemByte, &newItem)
	if err != nil {
		return fmt.Errorf("Item doesn't have right fields: id, email, age. Error is %w", err)
	}
	
	fileContent, err := getFileContent(fileName)
	if err != nil {
		return err
	}
	
	var items []Item
	err = json.Unmarshal(fileContent, &items)
	if err != nil && len(fileContent) > 0 {
		return fmt.Errorf("Json file doesn't have json object. Error is %w", err)
	}
	
	foundId := findIdOfElement(items, newItem.Id)
	if foundId != -1 {
		fmt.Fprintf(writer, "Item with id %s already exists", newItem.Id)
		return nil
	} else {
		items = append(items, newItem)
	}
	
	newFileContent, err := json.Marshal(items)
	if err != nil {
		return fmt.Errorf("Cannot Marshal new object with json string. Error is %w", err)
	}
	ioutil.WriteFile(fileName, newFileContent, 0755)

	return nil
}

//this is a handler for operation=remove
func removeItem(fileName string, id string, writer io.Writer) error {
	fileContent, err := getFileContent(fileName)
	if err != nil {
		return err
	}
	
	var items []Item
	err = json.Unmarshal(fileContent, &items)
	if err != nil && len(fileContent) > 0 {
		return fmt.Errorf("Json file doesn't have json object. Error is %w", err)
	}
	
	foundId := findIdOfElement(items, id)
	if foundId == -1 {
		fmt.Fprintf(writer, "Item with id %s not found", id)
		return nil
	}
	newItems := make([]Item, 0)
	newItems = append(newItems, items[:foundId]...)
	newItems = append(newItems, items[foundId+1:]...)
	
	newFileContent, err := json.Marshal(newItems)
	if err != nil {
		return fmt.Errorf("Cannot Marshal new object with json string. Error is %w", err)
	}
	ioutil.WriteFile(fileName, newFileContent, 0755)

	return nil
}

//this is a handler for operation=findById
func findItem(fileName string, id string, writer io.Writer) error {
	fileContent, err := getFileContent(fileName)
	if err != nil {
		return err
	}
	
	var items []Item
	err = json.Unmarshal(fileContent, &items)
	if err != nil && len(fileContent) > 0 {
		return fmt.Errorf("Json file doesn't have json object. Error is %w", err)
	}
	
	foundId := findIdOfElement(items, id)
	if foundId == -1 {
		fmt.Fprintf(writer, "")
		return nil
	}
	
	newItemByte, err := json.Marshal(items[foundId])
	if err != nil {
		return fmt.Errorf("Cannot Marshal new object with json string. Error is %w", err)
	}
	writer.Write(newItemByte)

	return nil
}

//get file content
func getFileContent(fileName string) ([]byte, error) {
	jsonFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil, fmt.Errorf("Cannot open file with filename %s and get an error %w", fileName, err)
	}
	defer jsonFile.Close()
  
	content, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, fmt.Errorf("Cannot read file with filename %s and get an error %w", fileName, err)
	}

	return content, nil
}

//find Item element in a slice by id
func findIdOfElement(items []Item, id string) int {
	for index, item := range items {
		if item.Id == id {
			return index
		}
	}
	return -1
}

//check if operation is allowed or no
func isAllowedOperation(op string) bool {
	return op == "add" || op == "list" || op == "findById" || op == "remove"
}

//transform command line arguments into Arguments map
func parseArgs() Arguments {
	cmds := os.Args[1:]
	args := make(Arguments, len(cmds) / 2)

	l := len(cmds);
	for i := 0; i < l; i+=2 {
		if (l > i+1) {
			id := strings.ReplaceAll(cmds[i], "\"", "")
			id = strings.ReplaceAll(id, "-", "")
			args[id] = cmds[i+1]
		}
	  }

	return args
}

func main() {
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}
