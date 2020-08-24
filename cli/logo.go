package cli

import (
    "fmt"
    "io/ioutil"
    "github.com/fatih/color"
)

func logoPrint() {
    file, err := ioutil.ReadFile("cli/logo.txt")
    if err != nil {
        fmt.Println(err)
    }
    col := color.New(color.FgHiCyan).Add(color.Bold)
	col.Println(string(file))
}