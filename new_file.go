package main

import (
    "os"
    "log"
    "io/ioutil"
    "encoding/json"
    "fmt"
)

type LOG struct {
    success     *log.Logger
    _err        *log.Logger
}

type FileInfo struct {
    dir             string
    _log            LOG
    all_files       []string
    to_delete       []string
    err_count       int
    success_count   int
    errors          []error
}

type JSONInfo struct {
    Dir             string      `json:"dir"`
    All_file        []string    `json:"all_files"`
    To_delete       []string    `json:"to_delete"`
    Success_count   int         `json:"success"`
    Err_count       int         `json:"failures"`
    Errors          []string    `json:"errors"`
}

func (info *FileInfo) success(msg string) {
    info._log.success.Println(msg)
    info.success_count++
}
func(info *FileInfo) _error(_e error, err string) {
    info._log._err.Println(err)
    info.err_count++
    info.errors = append(info.errors, _e)
}

func setup() *FileInfo {

    // Setting up loggers and default info for struct FileInfo
    l := LOG{}
    info := &FileInfo{err_count: 0}

    _, err := os.Stat("messages")
    info._log = l

    if err != nil {
        if os.IsNotExist(err) {
            _, e := os.Create("messages")

            if e != nil {
                log.Fatal(e)
            }
        }
    }

    file, _err := os.OpenFile("messages", os.O_WRONLY|os.O_CREATE, 0666)

    if _err != nil {
        log.Fatal(err)
    }

    l.success = log.New(file, "SUCCESS: ", log.Ldate|log.Ltime)
    l._err = log.New(file, "ERROR: ", log.Ldate|log.Ltime)

    dir, err := os.Getwd()

    if err != nil {
        info._error(err, err.Error())
    }

    info._log = l

    info.dir = dir
    info.success("Successfully grabbed main directory "+ info.dir)
    return info
}

func (info *FileInfo) read_dir() []string {
    var all_files []string

    files, err := ioutil.ReadDir(info.dir)

    if err != nil {
        info._error(err, "Couldn't read" + info.dir)
    }

    for _, f := range files {
        all_files = append(all_files, f.Name())
    }

    info.all_files = all_files

    return info.all_files
}

func (info *FileInfo) check_errors() {
    JI := JSONInfo{Dir: info.dir, All_file: info.all_files, To_delete: info.to_delete, Err_count: info.err_count, Success_count: info.success_count}

    for i := range info.errors {
        if i == len(info.errors) - 1 {
            break
        }
        JI.Errors = append(JI.Errors, info.errors[i].Error())
    }

    file, _ := json.MarshalIndent(JI, "", "\t")
    _ = ioutil.WriteFile("info.json", file, 0644)

    if info.err_count > info.success_count {
        panic(info.errors[info.err_count - 1])
    }
}

func (info *FileInfo) remove(filename string, size int) {
    err := os.Remove(filename)

    if err != nil {
        info._error(err, "Could not remove " + filename)
    }

    s := fmt.Sprintf("Successfully removed %s of size %dmb", filename, size / 1024 / 1024)
    info.success(s)
}

func (info *FileInfo) check_sizes() {
    for i := range info.all_files {
        file, err := os.Stat(info.all_files[i])

        if err != nil {
            info._error(err, "Error statting file " + info.all_files[i])
            continue
        }

        if file.Size()/1024/1024 > 1 {
            info.to_delete = append(info.to_delete, info.all_files[i])
            info.remove(info.to_delete[len(info.to_delete) - 1], int(file.Size()))
        }

        info.success("Successfully read " + info.all_files[i])
    }

    info.check_errors()
}

func main() {
    info := setup()
    info.read_dir()
    info.check_sizes()

    fmt.Println(info)
}
