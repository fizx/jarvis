package main

import (
	"bytes"
	"io"
	"log"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/iancoleman/strcase"

	"github.com/alexkappa/mustache"
	"github.com/fizx/jarvis/generated/assets"
	"github.com/shurcooL/httpfs/vfsutil"
)

//go:generate go run gen/main.go

func apply(template string, data map[string]string) string {
	m := mustache.New()
	err := m.ParseString(template)
	if err != nil {
		log.Fatalln(err)
	}
	s, err := m.RenderString(data)
	if err != nil {
		log.Fatalln(err)
	}
	return s
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("Jarvis expects a project name as the sole argument!")
	}
	data := map[string]string{}
	u, err := url.Parse(os.Args[1])
	if err != nil {
		u, err = url.Parse("https://" + os.Args[1])
	}
	if err != nil {
		log.Fatal(err)
	}

	segments := strings.Split(u.EscapedPath(), "/")
	name := segments[len(segments)-1]
	owner := segments[len(segments)-2]
	root := strcase.ToSnake(name)
	data["owner"] = owner
	data["project"] = root
	data["class"] = strcase.ToCamel(name)
	data["package"] = u.Hostname() + u.EscapedPath()
	fs := assets.Assets
	walkFn := func(templatePath string, fi os.FileInfo, r io.ReadSeeker, err error) error {
		realPath := apply(templatePath, data)
		switch fi.IsDir() {
		case false:
			localPath := path.Join(root, realPath)
			os.MkdirAll(path.Dir(localPath), 0755)
			buf := new(bytes.Buffer)
			buf.ReadFrom(r)
			s := buf.String()
			contents := apply(s, data)
			f, err := os.OpenFile(localPath, os.O_RDWR|os.O_CREATE, 0644)
			if err != nil {
				log.Fatal(err)
				return err
			}
			if err != nil {
				log.Fatal(err)
			}
			if _, err := f.Write([]byte(contents)); err != nil {
				f.Close()
				log.Fatal(err)
			}
			if err := f.Close(); err != nil {
				log.Fatal(err)
			}
		case true:
			localPath := path.Join(root, realPath)
			os.MkdirAll(localPath, 0755)
			return nil
		}
		return nil
	}
	err = vfsutil.WalkFiles(fs, "/", walkFn)
	if err != nil {
		log.Fatalln(err)
	}
}
