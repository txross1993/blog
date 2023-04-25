// You can edit this code!
// Click here and start typing.
package main

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"time"
)

func main() {
	var Contacts = []Person{
		{"Allison", "DuMonte", "08-01-1984"},
		{"Ben", "Ridley", "11-21-1992"},
		{"Christina", "Chandler", "02-19-1987"},
		{"Daniel", "Pulaski", "07-09-1989"},
	}

	var templateFns = map[string]any{
		"age": func(birthday string) int {
			const dateLayout = "01-02-2006"
			const hoursInYr = 8766
			t, _ := time.Parse(birthday, dateLayout)
			age := time.Since(t)
			return int(age.Hours()) / hoursInYr
		},
	}

	var Greeting = `Happy Birthday, {{ .FirstName }} {{ .LastName }}!
You are {{ age .Birthday }} years old today!`

	tpl, err := template.New("greeting").Funcs(templateFns).Parse(Greeting)
	if err != nil {
		fmt.Println("error parsing template: ", err)
		os.Exit(1)
	}

	dest := bytes.NewBuffer([]byte{})
	for _, c := range Contacts {
		err = tpl.Execute(dest, c)
		if err != nil {
			fmt.Println("error executing template: ", err)
			os.Exit(1)
		}
		fmt.Println(dest.String())
		dest.Reset()
	}
}

type Person struct {
	FirstName string
	LastName  string
	Birthday  string
}

func (p Person) Age() int {
	const dateLayout = "01-02-2006"
	const hoursInYr = 8766
	t, _ := time.Parse(p.Birthday, dateLayout)
	age := time.Since(t)
	return int(age.Hours()) / hoursInYr
}

func yearsOld(birthday string) int {
	const dateLayout = "01-02-2006"
	const hoursInYr = 8766
	t, _ := time.Parse(birthday, dateLayout)
	age := time.Since(t)
	return int(age.Hours()) / hoursInYr
}
