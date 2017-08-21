package main

import (
    "flag"

    "archive/tar"
    "compress/gzip"
    "path"
    "fmt"
    "io"
    "os"
)

var (
    filename string
    rel_path string
)

func init() {
    flag.StringVar(&filename, "f", "", "filename")
    flag.StringVar(&rel_path, "p", "", "filename")
}

func main() {
    flag.Parse()

    file, err := os.Open(filename)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    defer file.Close()

    // gzipの展開
    gzip_reader, err := gzip.NewReader(file)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    defer gzip_reader.Close()

    // tarの展開
    tar_reader := tar.NewReader(gzip_reader)

    i := 0
    for {
        header, err := tar_reader.Next()
        if err == io.EOF {
            break
        }

        name := header.Name
        name = path.Join(rel_path, name)

        switch header.Typeflag {
        case tar.TypeDir:
            continue
        case tar.TypeReg:
            fmt.Println("(", i, ")", "Name: ", name, header.Size)
        default:
            fmt.Printf("%s : %c %s %s\n",
                "Yikes! Unable to figure out type",
                header.Typeflag,
                "in file",
                name,
            )
        }
        i++
    }
}
