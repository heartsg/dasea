package logger

//File Writer

import (
    "os"
    "strconv"
)

type FileWriter struct {
    filepath        string // should be set to the actual filename
    rotate          bool
    sizelimit       int64
    rotatelimit     int
    currentrotate   int
    fp              *os.File
}

// Make a new FileWriter. Return nil if error occurs during setup.
func NewWriter(filepath string, rotate bool, sizelimit int64, rotatelimit int) *FileWriter {
    w := &FileWriter{filepath: filepath, rotate:rotate, sizelimit:sizelimit, rotatelimit: rotatelimit}
    if w.rotatelimit == 1 {
        w.rotatelimit = 2 //no point for one rotate, keep at least one backup
    }
    err := w.Init()
    if err != nil {
        return nil
    }
    return w
}

func (w *FileWriter) Init() (err error) {
    err = w.Close()
    if err != nil {
        return
    }
    
    w.currentrotate = 0
    
    var newfile string
    if !w.rotate {
        newfile = w.filepath
    } else {
        newfile = w.filepath + ".0"
    }
    w.fp, err = os.Create(newfile)
    
    return
}

func (w *FileWriter) Close() (err error) {
    if w.fp != nil {
        err = w.fp.Close()
        w.fp = nil
        if err != nil {
            return
        }
    }
    
    return
}

// Write satisfies the io.Writer interface.
func (w *FileWriter) Write(output []byte) (size int, err error) {
    if w.sizelimit > 0 {
        //check size limit before write
        var stat os.FileInfo
        stat, err = w.fp.Stat()
        if err != nil {
            return
        }
        if stat.Size() + int64(len(output)) > w.sizelimit {
            //oversize, try to put in one file
            //however, if len(output) is greater than sizelimit, have to chop anyway
            if int64(len(output)) > w.sizelimit {
                size, err = w.fp.Write(output[0:w.sizelimit - stat.Size()]) 
                if err != nil {
                    return
                }
                err = w.Rotate()
                if err != nil {
                    return
                }
                //make a recursive call, hope it won't be called very often
                var newsize int
                newsize, err = w.Write(output[size:len(output)])
                size = newsize + size
                return
            } else {
                err = w.Rotate()
                if err != nil {
                    return
                }
            }
        }
    }
    
    size, err = w.fp.Write(output)
    if err != nil {
        return
    }
    
    return
}

// Perform the actual act of rotating and reopening file.
func (w *FileWriter) Rotate() (err error) {
    err = w.Close()
    if err != nil {
        return
    }
    
    // Two different possibilities
    // 1. rotate == false, just delete backup, move current file to backup, and create new file
    // 2. rotate == true, if no rotatelimit, just incjrease and create new, if has rotatelimit, increase, and possibily delete old
    
    var delfile string
    if !w.rotate {
        //even no rotate, we intentionally keep at least one backup
        delfile = w.filepath + ".bak"
    } else if w.rotatelimit > 0 && w.currentrotate >= w.rotatelimit - 1 {
        delfile = w.filepath + "." + strconv.Itoa(w.currentrotate - w.rotatelimit + 1)
    }
    _, err = os.Stat(delfile)
    if err == nil {
        err = os.Remove(delfile)
        if err != nil {
            return
        }
    }
    
    if !w.rotate {
        //move current to ".bak"
        bakfile := w.filepath + ".bak"
        _, err = os.Stat(w.filepath)
        if err == nil {
            err = os.Rename(w.filepath, bakfile)
            if err != nil {
                return
            }
        }
   
    }

    // Create a file.
    var newfile string
    if !w.rotate {
        newfile = w.filepath
    } else {
        w.currentrotate++
        newfile = w.filepath + "." + strconv.Itoa(w.currentrotate)
    }
    w.fp, err = os.Create(newfile)
    return
}