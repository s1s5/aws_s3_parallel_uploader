package main

import (
    "flag"

    "archive/tar"
    "compress/gzip"

    "bytes"
    "fmt"
    "io"
    "os"
)

var (
    filename string
)

func init() {
    flag.StringVar(&filename, "f", "", "filename")
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
        // const (
        // TypeReg           = '0'    // regular file
        // TypeRegA          = '\x00' // regular file
        // TypeLink          = '1'    // hard link
        // TypeSymlink       = '2'    // symbolic link
        // TypeChar          = '3'    // character device node
        // TypeBlock         = '4'    // block device node
        // TypeDir           = '5'    // directory
        // TypeFifo          = '6'    // fifo node
        // TypeCont          = '7'    // reserved
        // TypeXHeader       = 'x'    // extended header
        // TypeXGlobalHeader = 'g'    // global extended header
        // TypeGNULongName   = 'L'    // Next file has a long name
        // TypeGNULongLink   = 'K'    // Next file symlinks to a file w/ a long name
        // TypeGNUSparse     = 'S'    // sparse file
        // )

        // type Header struct {
        // Name       string    // name of header file entry
        // Mode       int64     // permission and mode bits
        // Uid        int       // user id of owner
        // Gid        int       // group id of owner
        // Size       int64     // length in bytes
        // ModTime    time.Time // modified time
        // Typeflag   byte      // type of header entry
        // Linkname   string    // target name of link
        // Uname      string    // user name of owner
        // Gname      string    // group name of owner
        // Devmajor   int64     // major number of character or block device
        // Devminor   int64     // minor number of character or block device
        // AccessTime time.Time // access time
        // ChangeTime time.Time // status change time
        // Xattrs     map[string]string
        // }

        switch header.Typeflag {
        case tar.TypeDir:
            continue
        case tar.TypeReg:
            bb := new(bytes.Buffer)
            _, err := io.Copy(bb, tar_reader)
            if err != nil {
                fmt.Println("ExtractTarGz: Copy() failed: %s", err.Error())
            }
            fmt.Println("(", i, ")", "Name: ", name, header.Size, bb.Len())
            
            // fmt.Println(" --- ")
            // io.Copy(os.Stdout, tar_reader)
            // fmt.Println(" --- ")
            // os.Exit(0)
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


// if err != nil {
//  log.Fatalf("ExtractTarGz: Next() failed: %s", err.Error())
// }

// switch header.Typeflag {
//  case tar.TypeDir:
//  if err := os.Mkdir(header.Name, 0755); err != nil {
//   log.Fatalf("ExtractTarGz: Mkdir() failed: %s", err.Error())
//  }
//  case tar.TypeReg:
//  outFile, err := os.Create(header.Name)
//  if err != nil {
//   log.Fatalf("ExtractTarGz: Create() failed: %s", err.Error())
//  }
//  defer outFile.Close()
//  if _, err := io.Copy(outFile, tarReader); err != nil {
//   log.Fatalf("ExtractTarGz: Copy() failed: %s", err.Error())
//  }
//  default:
//  log.Fatalf(
//   "ExtractTarGz: uknown type: %s in %s",
//   header.Typeflag,
//   header.Name)
// }
