---
title: How to Use Text Templates in Golang with Complext Data Inputs
---

Have you ever used Go's templating methods but got stuck when your input data structure contained something less trivial than top-level fields? If so, this blog post is for you! I will walk through simple examples using a struct with built-in property types. Then I will show you pipelining functions and how to deal with nilable data. 

# Getting Started: A Simple Greeting
If you're not familiar, Golang's `text/tempalte` [package](https://pkg.go.dev/text/template) provides ways to format text data using objects declared in a program. The simplest use case is accessing data from a simple structure containing top level fields. Let's look at an example, where `Person` has a `FirstName`, `LastName`, and a `Birthday` as string properties. 

```golang
type Person struct {
    FirstName   string
    LastName    string
    Birthday    string
}
```
Imagine you are building a list of your contacts and want to automatically construct a birthday greeting card for each person.

Let's break down this code step by step. 

```golang
var Contacts = []Person{
    {"Allison", "DuMonte", "08-01-1984"},
    {"Ben", "Ridley", "11-21-1992"},
    {"Christina", "Chandler", "02-19-1987"},
    {"Daniel", "Pulaski", "07-09-1989"},
}
```

First, we declare a list of contacts - people that we know and their birthdays in text form. This acts as a set of input data.

{% raw %}
```golang
var Greeting = var Greeting = `Happy Birthday, {{ .FirstName }} {{ .LastName }}!
You were born on {{ .Birthday }}.`
```
{% endraw %}

## Syntax notes

The first thing you'll notice in the greeting string may be the curly brace syntax. The {% raw %}`{{` opening and `}}`{% endraw %} closing braces are called **delimeters** and signal to the text template library where it should replace text with that of your input data. You may change these delimeters but this blog post will not go into detail about that. 

Next you may notice the syntax `.FirstName` - a **dot** followed by the exact name of a field of our `Person` struct. This dot-notated field name will be eerily familiar to you if you have ever used [jq](https://stedolan.github.io/jq/). If you haven't used `jq`, you can think of the `.` as an accessor to the root of your input data. So the struct itself is `Person`, and you access the field `FirstName` of Person using the accessor `.FirstName`. 

In plain english, the templating package will read from left to right of the input text and at each point it encounters a **delimeter**, it will attempt to extract the value it should replace the delimited object with. If an error occurs, the templater will stop and return to you the string it has parsed so far and any errors. Let's define our template using the template package.

```golang
tpl, err := template.New("greeting").Parse(Greeting)
if err != nil {
    fmt.Println("error parsing template: ", err)
    os.Exit(1)
}
```

The code above defines our new template processing object, `tpl` by parsing our input `Greeting`. The `New()` method accepts a string, which is the name of the template you're defining. In practice, I have never found the name of the tempalte to be important in code, so you can name it whatever. If the input to `Parse()` is invalid, the method will return an error and we will stop processing. 

Finally, let's write the data out. I won't go into the details of what an `io.Writer` is here. Just know that when `Execute` is called, it will read and process input from the beginning until `EOF` or an error is encountered.

```golang
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
```
Output: 
```bash
Happy Birthday, Allison DuMonte!
You were born on 08-01-1984.
Happy Birthday, Ben Ridley!
You were born on 11-21-1992.
Happy Birthday, Christina Chandler!
You were born on 02-19-1987.
Happy Birthday, Daniel Pulaski!
You were born on 07-09-1989.
```

# Getting Fancy: Custom Functions

The greeting example is not very personable, since we simply state the person's name and birthday. Let's make it more personal by calculating the age from the input birthday string by providing our own custom function to the templater. 

There are two ways you can do this: by providing a method on the struct `Person` or by defining a function that will take arguments and return a single value or return a value and an error. 

## Map of Template Functions with Args
{% raw %}
```golang
var templateFns = map[string]any{
    "age": func(birthday string) int {
        const dateLayout = "01-02-2006"
	    const hoursInYr = 8766
	    t, _ := time.Parse(birthday, dateLayout)
	    age := time.Since(t)
	    return int(age.Hours()) / hoursInYr
    },
}

tpl, err := template
    .New("greeting")
    .Funcs(templateFns)
    .Parse(Greeting)

if err != nil {
    fmt.Println("error parsing template: ", err)
    os.Exit(1)
}

var Greeting = `Happy Birthday, {{ .FirstName }} {{ .LastName }}!
You are {{ age .Birthday }} years old today!`

```
{% endraw %}

The `map[string]any` is the type of the arg that can be provided to your template before you call `Parse()`. The `string` is the name of the function, `age` is used inside if {% raw %}`{{}}`{%endraw%}, as you can see referenced in the template string as {%raw%}`{{ age .Birthday }}`{%endraw%}

## Function as a Method on a Struct
{% raw %}
```golang
func (p Person) Age() int {
    const dateLayout = "01-02-2006"
	const hoursInYr = 8766
	t, _ := time.Parse(p.Birthday, dateLayout)
	age := time.Since(t)
	return int(age.Hours()) / hoursInYr
}

var Greeting = `Happy Birthday, {{ .FirstName }} {{ .LastName }}!
You are {{ .Age }} years old today!`
```
{% endraw %}

Like properties of a struct, methods defined on a struct may be used inside of the template function as well. They are referenced just like properties and must be exported (start with a capitalized letter) for them to be used in templating functions.
That means if we defined the function {% highlight go %}`func (p Person) age() { ... }`{%endhighlight%} it would not be available in templating functions. Since `age` does not begin with a capitalized letter, that instructs the compiler that it is not an exported function.

# Dealing with Complex Data

Input data isn't always as simple as a struct with top level attributes that are built-in types. Often you'll have structs that have properties that are structs themselves, or even pointers to structs, slices, or maps. That means you also have the dreaded possibility of <span style="color:red">*NIL POINTER DEREFERENCE*!</span>

Moreover, you may not always be the person who writes the templating functions, you may only be able to write the template. Luckily, the template library provides clever ways to interact with complex data.  

Let's look at a more complex data type with maps, slices, and pointers to structs.

```golang
type Node struct {
    Left, Right *Node
    Data map[string]interface{}
}

type Tree struct {
    Nodes []*Node
}

var (
    left = &Node{
        Data: map[string]interface{}{
            "key1": "left1",
            "key2": "left2",
        },
    }
    right = &Node{
        Data: map[string]interface{}{
            "key1": "right1",
        },
    }
    root = &Node{Left: left, Right: right}
    tree = &Tree{Nodes: []*Node{root, left, right}}
)
```

Let's say we want to template out this tree and print all the data points for `key1`. 

{% raw %}
```golang
var templateStr = `All the key1 keys: [
{{- range $node := .Nodes -}}
    {{- if $node.Data }} {{ $node.Data.key1 -}} {{ end -}}
{{- end -}}
]`
```
{% endraw %}

```bash
# output
All the key1 keys: [ left1  right1 ]
```

## Syntax Notes
The `range` function will let us range over the `Tree`'s nodes and find all the `key1`.
Some things to note:
  - Like the template function `"age"` we defined in the `map[string]any` above, the built-in `range` function uses the `function <assignment> arg` syntax.
  - The **assignment** syntax is `$varName :=`, where `$varName` can be anything you want, just preceded with a `$` symbol.
  - The {% raw %}`{{ if $node.Data }}`{% endraw %} will only evaluate to `true` if the pipelined `.Data` field is not null. In the case of a map that means the map is not empty. If the `Data` field were a pointer, the `if` statement would evaluate to if the pointer was non-nil. If it were a string, it would evaluate if the string was non-empty. If it were an integer, it would evaluate if it was not `0`. And so on...
  - The two {% raw %}`{{ end }}`{% endraw %} are used to end the {% raw %}`{{ if }}`{% endraw %} and {% raw %}`{{ range }}`{% endraw %} blocks respectively.

What if you only wanted to find the first instance of `key1` and use that value? Well just like how regular code has flow control features, so does the text templating engine. You can use:
  - {% raw %}`{{ if ... }} {{ else }} {{ end }}`{% endraw %}
  - {% raw %}`{{ with ... }} {{ else }} {{ end }}`{% endraw %}
  - {% raw %}`{{ range ... }} {{ continue }} {{ break }} {{ end }}`{% endraw %}

## Pipeline Examples
Let's see some examples:

1. I want to print all the left node keys (prefixed with string `left`):

{% raw %}
  ```golang
  var templateStr = `Keys left: {{ range $node := .Nodes -}}
    {{- range $key := $node.Data -}}
        {{- with $keyPrefix := (slice $key 0 4) -}}
            {{- if eq $keyPrefix "left" }} {{- $key }} {{ end -}}
        {{- end -}}
    {{- end -}}
{{- end -}}`
  ```
  {% endraw %}

  ```bash
  #output
  Keys left: left1 left2 
  ```

 - **Note**: The hypen syntax `-` inside delimiters {% raw %}`{{- ... -}}`{% endraw %} instructs the templating engine to remove spaces/tabs/newlines/returns on whichever side the `-` is.<br><br>
 - **Note**: The parenthetical `(slice $key 0 4)` instructs the templating engine to evaluate the value first before evaluating the rest of the with statement. `slice $key 0 4` is equivalent to slicing the string to the 4th index non-inclusive. For keys that start with `right`, that would mean the function returns `righ`. For keys that start with `left`, it will return `left`.<br><br>
 - **Note**: For each flow control statement, you must end it with {% raw %}`{{ end }}`{% endraw %} or else your template will fail to parse.  

{:start="2"}
2. I want to print only the first key prefixed `left` and stop:

{% raw %}
```golang
var templateStr = `First left key: {{ range $node := .Nodes -}}
    {{- range $key := $node.Data -}}
        {{- with $keyPrefix := (slice $key 0 4) -}}
            {{- if eq $keyPrefix "left" }} {{- $key }} {{ break }} {{ end -}}
        {{- end -}}
    {{- end -}}
{{- end -}}`
```
{% endraw %}

```bash
#output
First left key: left1 
```

There is more that you can do with text templates, and if you'd like to explore more please check out the documentation for the [text/template](https://pkg.go.dev/text/template@go1.20.3) package. 

For now, I hope you've learned enough to create something to fit your needs. 

# Source Code
If you need more details see the code examples in the repo:
  - [Simple Person Birthday example]({{ site.repository_url }}/blob/master/_examples/go/template_simple.go)
  - [Pipelined Tree example]({{ site.repository_url }}/blob/master/_examples/go/template_pipeline.go)