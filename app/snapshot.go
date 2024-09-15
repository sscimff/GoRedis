package main

import (
	"flag"
)

type Config struct {
	Dir        string
	DBFileName string
}

var config = Config{
	Dir:        "/data",
	DBFileName: "dump.rdb",
}


func init() {
    flag.StringVar(&config.Dir, "dir", "/data", "Directory for RDB files")
    flag.StringVar(&config.DBFileName, "dbfilename", "dump.rdb", "RDB filename")
    flag.Parse()
}
