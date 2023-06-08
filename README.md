# go-yadt

This `Go` package can generate `.docx` files from your template and `json` data.

> **Warning**
> This package requires an [page merger CLI tool](https://github.com/tymbaca/pagemerger).
> For now the only option is to build it from source. I hope could publish an `.deb` package of binary soon.

## Install

To install package go to your `go` project directory and run:

```sh
go get github.com/tymbaca/go-yadt
```

Then you can import it by adding `github.com/tymbaca/go-yadt`.

## Usage
Use package via `yadt` shortcut.

### Input data
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

// Read bytes from json file
bodyBytes, _ = os.ReadFile("data.json")

// Create instance of FileGenerator
fileGenerator, err := yadt.New("template.docx")
if err != nil {
    panic(err)
}

// Generate .docx files and pack to .zip
err = yadt.GenerateZip("output.zip")
if err != nil {
    panic(err)
}
```


## Roadmap
- [x] Core functionality
- [ ] Generate without packing to `zip`
- [ ] Do something with page merging CLI dependency...
- [ ] More detailed error messages
