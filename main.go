package main

import (
    //"fmt"
    "io/ioutil"
    "log"
    "os"
    "math"
)

const no_permission = 0000
const r_w_e_o_o = 0700 // Read, Write, Execute owner only
const r_w_e_o_a_g = 0770 // Read, Write, Execute owner and group
const r_w_e_o_g_o = 0777 // Read, Write, Execute owner, group and others
const execute = 0111
const write = 0222
const write_and_execute = 0333
const read = 0444
const read_and_execute = 0555
const read_and_write = 0666
const special = 0740 // Only owner can read and write. owner and group can execute. Others have no permission

type charInfo struct {
    invalidChars    []byte
    validChars      []byte
}

type fileInfo struct {
    filename    string
    data        []byte
    filesize    int64
    size_in_mb  float32
    c_i         []charInfo
    index       int
}

func setup(filename string, data []byte) *fileInfo {
    info := fileInfo{ filename: filename, data: data }
    
    _, e := os.Stat(info.filename)
    
    if e != nil {
        if os.IsExist(e) {
            log.Fatal("File already exists.")
            os.Exit(1)
        }
    }
    
    file, err := os.Create(info.filename)
    
    if err != nil {
        log.Fatal("Error writing to the file.")
    }

    size, err := file.Write(info.data)
    file.Close()
    
    info.filesize = int64(size)
    info.size_in_mb = float32((info.filesize / 1024) / 1024)
    return &info
}

func (info *fileInfo) new_info(filename string) *fileInfo {

    if filename != "" {
        info.filename = filename
    }

    file, err := os.Stat(info.filename)
    info.filesize = file.Size()
    
    if err != nil {
        log.Fatal("Error with openeing and reading the file")
    }
    
    if int64(file.Size()) == info.filesize {
        data, err := ioutil.ReadFile(info.filename)

        if err != nil {
            //log.Fatal("Error reading the data")
            log.Fatal(err)
        }

        info.size_in_mb = float32(math.Round(float64((info.filesize / 1024) / 1024)))

        return &fileInfo{filename: info.filename, filesize: info.filesize, data: data, size_in_mb: info.size_in_mb}
    }
    return &fileInfo{filename: info.filename, filesize: 0, data: nil}
}

func (info *fileInfo) gather_all() {
    c_i := charInfo{ }
    for i := 0; i < int(info.filesize); i++ {
        if info.data[i] < 0x41 || info.data[i] > 0x7a {
            if info.data[i] == 0x20 || info.data[i] == 0x0a {
                c_i.validChars = append(c_i.validChars, info.data[i])
            }
            c_i.invalidChars = append(c_i.invalidChars, info.data[i])
        } else {
            c_i.validChars = append(c_i.validChars, info.data[i])
        }
    }

    if len(c_i.invalidChars) > 0 {
        info.c_i = append(info.c_i, c_i)
        info.index += 1
    }
}

func (info fileInfo) write_valid(toFile string) {
    if toFile == "" {
        log.Fatal("Cannot write to a undefined file.")
    }
    file, err := os.Stat(info.filename)

    if err != nil {
        log.Fatal(err)
    }

    if file.Size() == info.filesize {
        info.filename = toFile
        _, err := os.Stat(info.filename)

        if err != nil {
            if os.IsNotExist(err) {
                _,err := os.Create(info.filename)

                if err != nil {
                    log.Fatal(err)
                }
            }
        }
        
        file, err := os.OpenFile(info.filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
        defer file.Close()

        if err != nil {
            log.Fatal(err)
        }

        file.WriteString(string(info.c_i[info.index - 1].validChars))
    }
}

func (info fileInfo) write_invalid(toFile string) {
    if toFile == "" {
        log.Fatal("Cannot write to a undefined file.")
    }
    file, err := os.Stat(info.filename)

    if err != nil {
        log.Fatal(err)
    }

    if file.Size() == info.filesize {
        info.filename = toFile
        _, err := os.Stat(info.filename)

        if err != nil {
            if os.IsNotExist(err) {
                _,err := os.Create(info.filename)

                if err != nil {
                    log.Fatal(err)
                }
            }
        }

        file, err := os.OpenFile(info.filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
        defer file.Close()

        if err != nil {
            log.Fatal(err)
        }

        file.WriteString(string(info.c_i[info.index - 1].invalidChars))
    }
}

func main() { 
    str := "Hey bro"
    f_info := setup("man", []byte(str))
    f_info = f_info.new_info("main")

    //fmt.Print(string(f_info.data), " -> \n\t\t'" + f_info.filename + "' has ", f_info.filesize * 8, " bits, which is ", f_info.filesize, " in bytes which is")
    //fmt.Printf(" %.1f MB!", f_info.size_in_mb)

    f_info.gather_all()
    f_info.write_valid("man")
    f_info.write_invalid("man2")
}