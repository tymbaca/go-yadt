# Yet Another Docx Templater

This `Go` package can generate `.docx` files from your template and `json` data.

> **Warning**
> This package requires an [page merger CLI tool](https://github.com/tymbaca/pagemerger).
> For now the only option is to build it from source. I hope I could publish an `.deb` package and homebrew tap soon.

## Install

To install package go to your `go` project directory and run:

```sh
go get github.com/tymbaca/go-yadt
```

Then you can import it by adding `github.com/tymbaca/go-yadt`.

## Usage
Use package via `yadt` shortcut.

### Input `docx` template
Package needs a *template* - `docx` file with placeholder fields. You need to specify that fields in curly braces in following format:
```
// template.docx

// docx text

Hello, {firstName} {lastName}! How do you do?

// more docx text
```

In this example we have two fields:
- `firstName`
- `lastName`

You will need to specify values for that field in `json` data.

### Input `json` data
For now this package only **eats raw `json` data in bytes**. You can set up parsing raw json from file or by your favorite web framework.

> For example in `gin` you can get raw json data by calling `c.GetRawData()` (where `c` is a `*gin.Context`)

Now let's set up our `json` into the file and name it `data.json`:

```json
[
    {
        "filename":"First file",
        "pages": [
            {
                "firstName": "Alex",
                "lastName": "Brown"
            },
            {
                "firstName": "John",
                "lastName": "White"
            }
        ]
    },
    {
        "filename":"Second file",
        "pages": [
            {
                "firstName": "Karina",
                "lastName": "Stone"
            },
            {
                "firstName": "Simon",
                "lastName": "Rock"
            }
        ]
    }
]
```

This structure is almost fully mandatory. Only items in `pages` can be deferent.

The input *rules* are simple:

- Whole json body is array
- Every item of that array is an object with two fields:
  - `filename` which contains result file name (**without extension**)
  - `pages` array. Its items are objects whose keys are **identical to keys in your template**.

> Notice that 'page' isn't necessarily means that size of your template needs to be 1 page. It can be more. Or less.

### Using in `Go`

```go
// main.go

import (
    "github.com/tymbaca/go-yadt"
)


// Create instance of FileGenerator
fileGenerator, err := yadt.NewFromFiles("template.docx", "data.json")
if err != nil {
    panic(err)
}

// Generate .docx files and pack to .zip
err = yadt.GenerateZip("output.zip")
if err != nil {
    panic(err)
}
```

On output you will get a `output.zip` archive which contains all files specified in `data.json`.

> Don't be shy, clone repository and run some tests (`./tests/` and `./utils/tests`).

## Roadmap
- [x] Core functionality
- [ ] Generate without packing to `zip`
- [ ] Changeable keys separator in `.docx` template
- [ ] Do something with page merging CLI dependency...
- [ ] More detailed error messages


## Questions
- [ ] What if some placeholder appear in `docx` template more than once? 
  - I think this question more related to github.com/lukasjarosch/go-docx... 
