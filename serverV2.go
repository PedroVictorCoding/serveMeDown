package main

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Purple = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"

var PORT string = ":9390"
var DEBUG bool = true

func main() {
	if runtime.GOOS == "windows" {
		Reset = ""
		Red = ""
		Green = ""
		Yellow = ""
		Blue = ""
		Purple = ""
		Cyan = ""
		Gray = ""
		White = ""
	}

	if DEBUG {
		fmt.Println("Debug Mode")
	}

	// Running twice, for vibe check
	batch_converter()
	batch_converter()

	fileServer := http.FileServer(http.Dir("./contents_html"))
	http.Handle("/", fileServer)

	assetsServer := http.FileServer((http.Dir("./assets")))
	http.Handle("/assets/", assetsServer)

	fmt.Printf("Starting server at port" + PORT + "\n")
	if err := http.ListenAndServe(PORT, nil); err != nil {
		log.Fatal(err)
	}
}

//////////////////////
//					//
// BATCH CONVERSION //
//					//
//////////////////////

func batch_converter() {
	fmt.Println("Starting Conversion")

	//TODO remove_filetype()

	filepath.WalkDir("./contents", func(s string, d fs.DirEntry, err error) error {
		var desired_path = "./contents_html/"
		var original_path = "./"
		var file_name = ""
		if err != nil {
			return err
		}
		if !d.IsDir() {
			var file_final_path = ""

			sz := len(s)
			if sz > 0 && s[sz-3] == '.' {
				file_final_path += s[:sz-3]
			}

			parts := strings.Split(file_final_path, "/")
			for part := range parts {
				if part != len(parts)-1 {
					desired_path += parts[part] + "/"
					original_path += parts[part] + "/"
				} else {
					original_path += parts[part] + ".md"
					file_name = parts[part]
				}
			}
			fmt.Println("original:\t\t", original_path)
			fmt.Println("desired file path:\t", desired_path+file_name)
			fmt.Println("desired folder path:\t", desired_path)

			os.Mkdir(desired_path, 0777)

			original_file_string, _ := os.ReadFile(original_path)

			var converted = mdToHTML(original_file_string, file_name)

			if err := os.WriteFile(desired_path+file_name, converted, 0666); err != nil {
				fmt.Println("\033[31m"+"ERR L2:\t\t\t", err, "\033[0m \n")
			}

		}
		return nil
	})

	// items, _ := os.ReadDir("./contents")
	// for _, item := range items {
	// 	fmt.Println(item.Name())
	// 	var file_string, _ = os.ReadFile("./contents/" + item.Name())
	// 	var file_converted = mdToHTML(file_string)
	// 	if err := os.WriteFile("./contents_html/"+item.Name(), file_converted, 0666); err != nil {
	// 		fmt.Println(err)
	// 	}

	// }

}

func mdToHTML(md []byte, title string) []byte {
	// create markdown parser with extensions
	extensions := (parser.CommonExtensions |
		parser.AutoHeadingIDs |
		parser.NoEmptyLineBeforeBlock |
		parser.Tables |
		parser.FencedCode |
		parser.Footnotes |
		parser.Autolink)
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// create HTML renderer with extensions
	htmlFlags := (html.CommonFlags |
		html.HrefTargetBlank |
		html.TOC | html.CompletePage |
		html.Smartypants |
		html.SmartypantsAngledQuotes |
		html.SmartypantsDashes |
		html.SmartypantsFractions |
		html.SmartypantsLatexDashes |
		html.LazyLoadImages)
	opts := html.RendererOptions{
		Title: title,
		CSS:   "/assets/global.css",
		Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	//var inject_html1 = ""
	//var final_code = []byte(inject_html1 + string(markdown.Render(doc, renderer)))
	return markdown.Render(doc, renderer)
}
