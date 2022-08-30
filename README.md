<h1 align="center">
  <div>
    <img src="https://raw.githubusercontent.com/adrg/adrg.github.io/master/assets/projects/frontmatter/logo.png" width="120px" alt="frontmatter logo"/>
  </div>
  frontmatter
</h1>

<h3 align="center">Go library for detecting and decoding various content front matter formats.</h3>

<p align="center">
    <a href="https://pkg.go.dev/github.com/aisbergg/go-frontmatter">
        <img alt="pkg.go.dev documentation" src="https://pkg.go.dev/badge/github.com/aisbergg/go-frontmatter">
    </a>
    <a href="https://opensource.org/licenses/MIT" rel="nofollow">
        <img alt="MIT License" src="https://img.shields.io/github/license/adrg/frontmatter"/>
    </a>
    <br />
    <a href="https://goreportcard.com/report/github.com/aisbergg/go-frontmatter">
        <img alt="Go report card" src="https://goreportcard.com/badge/github.com/aisbergg/go-frontmatter?style=flat" />
    </a>
    <a href="https://github.com/aisbergg/go-frontmatter/graphs/contributors">
        <img alt="GitHub contributors" src="https://img.shields.io/github/contributors/adrg/frontmatter" />
    </a>
    <a href="https://github.com/aisbergg/go-frontmatter/issues?q=is%3Aopen+is%3Aissue">
        <img alt="GitHub open issues" src="https://img.shields.io/github/issues-raw/adrg/frontmatter">
    </a>
    <a href="https://github.com/aisbergg/go-frontmatter/issues?q=is%3Aissue+is%3Aclosed">
        <img alt="GitHub closed issues" src="https://img.shields.io/github/issues-closed-raw/adrg/frontmatter" />
    </a>
</p>

## Installation

```sh
go get github.com/aisbergg/go-frontmatter
```

## Usage

**Default usage.**

```go
package main

import (
	"fmt"
	"strings"

	"github.com/aisbergg/go-frontmatter/pkg/frontmatter"
)

var input = `
---json
{
  "name": "frontmatter",
  "tags": ["foo", "bar", "baz"]
}
---
rest of the content
`

func main() {
	var matter struct {
		Name string   `json:"name"`
		Tags []string `json:"tags"`
	}

	body, err := frontmatter.Parse(strings.NewReader(input), &matter)
	if err != nil {
		panic(err)
	}
	// NOTE: If a front matter must be present in the input data, use
	//       frontmatter.MustParse instead.

	fmt.Printf("%+v\n", matter)
	fmt.Println(string(body))

	// Output:
	// {Name:frontmatter Tags:[foo bar baz]}
	// rest of the content
}
```

**Bring your own formats.**

> This library includes only a JSON format by default. This removes the need for external dependencies and gives your the freedom to choose whatever unmarshaller library you want.

If you like to use any other formats than JSON, you can easily add them. Here is how:

```go
package main

import (
	"fmt"
	"strings"

	"github.com/aisbergg/go-frontmatter"
	"gopkg.in/yaml.v3"
)

var input = `
---
name: "frontmatter"
"tags": ["foo", "bar", "baz"]
...
rest of the content
`

func main() {
	var matter struct {
		Name string   `yaml:"name"`
		Tags []string `yaml:"tags"`
	}

	formats := []*frontmatter.Format{
		frontmatter.NewFormat("---", "...", yaml.Unmarshal),
	}

	rest, err := frontmatter.Parse(strings.NewReader(input), &matter, formats...)
	if err != nil {
		// Treat error.
	}
	// NOTE: If a front matter must be present in the input data, use
	//       frontmatter.MustParse instead.

	fmt.Printf("%+v\n", matter)
	fmt.Println(string(rest))

	// Output:
	// {Name:frontmatter Tags:[foo bar baz]}
	// rest of the content
}
```

Full documentation can be found at: https://pkg.go.dev/github.com/aisbergg/go-frontmatter.

## Contributing

Contributions in the form of pull requests, issues or just general feedback, are always welcome. See [CONTRIBUTING.md](CONTRIBUTING.md).

## Acknowledgements

This project is a fork of [github.com/adrg/frontmatter](https://github.com/adrg/frontmatter) developed by [Adrian-George Bostan](https://github.com/adrg). This fork removes the external dependencies and improves performance.

## Licence

This project is under the MIT Licence. See the [LICENCE](LICENCE) file for the full license text.
