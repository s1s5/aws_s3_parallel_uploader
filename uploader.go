package main

import (
    "flag"
    "github.com/docker/goamz/aws"
    "github.com/docker/goamz/s3"
    "archive/tar"
    "compress/gzip"
    "path"
    "bytes"
    "fmt"
    "io"
    "os"
)

type file_and_bb struct {
    name string
    buffer *bytes.Buffer
}

var (
    num_workers int = 4
    filename string
    rel_path string = "/"
    bucket_name string
)

func init() {
    flag.StringVar(&bucket_name, "b", "", "Bucket Name")
    flag.StringVar(&filename, "f", "", "filename")
    flag.StringVar(&rel_path, "p", "", "filename")
}

func consumer(bucket *s3.Bucket, message <-chan *file_and_bb, done chan<- bool) {
    for bb := range message {
        // fmt.Println("got bb", bb.name, ":", bb.buffer.Len())

        // err = bucket.Put(bb.name, bb.buffer, "application/octet-stream", s3.BucketOwnerFull, s3.Options{})
        err := bucket.Put(bb.name, bb.buffer.Bytes(), "application/octet-stream", s3.Private, s3.Options{})
        if err != nil {
            fmt.Println("failed ", bb.name, ":", bb.buffer.Len(), err, err.Error())
            continue
        }
        fmt.Println("uploaded", bb.name, ":", bb.buffer.Len())
    }
    done <- true
}

func main() {
    flag.Parse()

    auth, err := aws.EnvAuth()
    if err != nil {
        panic(err.Error())
    }

    // Open Bucket
    s := s3.New(auth, aws.APNortheast)
    // s := s3.New(auth, aws.USWest2)
    bucket := s.Bucket(bucket_name)

    message := make(chan *file_and_bb, 10)
    done := make(chan bool)

    for i := 0; i < num_workers; i++ {
        go consumer(bucket, message, done)
    }


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
            fb := &file_and_bb{name, new(bytes.Buffer)}
            _, err := io.Copy(fb.buffer, tar_reader)
            if err != nil {
                fmt.Println("ExtractTarGz: Copy() failed: %s", err.Error())
            } else {
                message <- fb
            }
        default:
            fmt.Printf("%s : %c %s %s\n", "Yikes! Unable to figure out type",
                header.Typeflag, "in file", name,)
        }
        i++
    }

    close(message)
    for i := 0; i < num_workers; i++ {
        <- done
    }
}
