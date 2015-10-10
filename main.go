package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"golang.org/x/net/html"
)

var (
	errNoTag = errors.New("no tag")
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		PrintPatches(os.Stdin)
	} else {
		for _, file := range args {
			f, err := os.Open(file)
			if err != nil {
				log.Fatalln("os.Open:", err)
			}
			PrintPatches(f)
			f.Close()
		}
	}
}

func PrintPatches(r io.Reader) {
	patches, err := ReadPatches(r)
	if err != nil {
		log.Fatalln("ReadPatches:", err)
	}
	for _, p := range patches {
		if !p.IsPlan9() {
			continue
		}
		fmt.Printf("%s %s", p.Name, strings.ToLower(p.Status))
		for _, dep := range p.Deps {
			fmt.Printf(" %s", dep)
		}
		fmt.Printf("\n")
	}
}

func ReadPatches(r io.Reader) ([]*Patch, error) {
	p := &Tokenizer{
		z: html.NewTokenizer(r),
	}
	patches := make([]*Patch, 0, 100)
	for {
		ok := p.LookupOptionalTag("table", "body")
		if !ok {
			break
		}
		ok = p.LookupOptionalTag("th", "table")
		if !ok {
			// skip a table that is not contained a patch.
			continue
		}
		for {
			if ok = p.LookupOptionalTag("tr", "table"); !ok {
				break
			}
			patch, err := ReadPatch(p)
			if err != nil {
				return nil, err
			}
			patches = append(patches, patch)
		}
	}
	if p.err != nil {
		return nil, p.err
	}
	return patches, nil
}

func ReadPatch(p *Tokenizer) (*Patch, error) {
	var patch Patch

	p.LookupTag("td", "tr")
	p.LookupTag("a", "td")
	patch.URL = p.Attr("href")
	patch.Name = p.Text()
	p.LookupTag("td", "tr")
	patch.Description = p.Text()  // TODO: can't handle a text including tags.
	p.LookupTag("td", "tr")
	patch.Author = p.Text()
	p.LookupTag("td", "tr")
	patch.Status = p.Text()
	p.LookupTag("td", "tr")
	patch.Deps = make([]string, 0, 5)
	ok := p.LookupOptionalTag("a", "td")
	for ok {
		patch.Deps = append(patch.Deps, p.Text())
		ok = p.LookupOptionalTag("a", "td")
	}
	if p.err != nil {
		return nil, p.err
	}
	return &patch, nil
}

type Patch struct {
	Name        string
	URL         string
	Description string
	Author      string
	Status      string
	Deps        []string
}

func (p *Patch) IsPlan9() bool {
	return strings.HasPrefix(p.URL, "9legacy/patch")
}

func (p *Patch) IsP9P() bool {
	return strings.HasPrefix(p.URL, "p9p/patch")
}

type Tokenizer struct {
	z   *html.Tokenizer
	err error
}

func (p *Tokenizer) LookupTag(tagName, innerTagName string) {
	if p.err != nil {
		return
	}
	ok := p.LookupOptionalTag(tagName, innerTagName)
	if !ok && p.err == nil {
		p.err = errNoTag
	}
}

func (p *Tokenizer) LookupOptionalTag(tagName, innerTagName string) bool {
	if p.err != nil {
		return false
	}
	var tag []byte
	for {
		t := p.z.Next()
		switch t {
		case html.ErrorToken:
			p.err = p.z.Err()
			return false
		case html.EndTagToken:
			tag, _ = p.z.TagName()
			if string(tag) == innerTagName {
				return false
			}
		case html.StartTagToken:
			tag, _ = p.z.TagName()
			if string(tag) == tagName {
				return true
			}
		}
	}
}

func (p *Tokenizer) Attr(name string) string {
	if p.err != nil {
		return ""
	}
	ok := true
	for ok {
		var key, val []byte
		key, val, ok = p.z.TagAttr()
		if string(key) == name {
			return string(val)
		}
	}
	return ""
}

func (p *Tokenizer) Text() string {
	if p.err != nil {
		return ""
	}
	t := p.z.Next()
	switch t {
	case html.ErrorToken:
		p.err = p.z.Err()
		return ""
	case html.TextToken:
		return string(p.z.Text())
	}
	p.err = errNoTag
	return ""
}
