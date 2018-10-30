package fs

import (
	"fmt"
	"io/ioutil"
	"os"
	"syscall"
	"unicode"
)

func isDigits(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}

// GetInode returns Inode for file
func GetInode(file string) (uint64, error) {
	fileinfo, err := os.Stat(file)
	if err != nil {
		return 0, fmt.Errorf("error stat file: %+v", err)
	}
	stat, ok := fileinfo.Sys().(*syscall.Stat_t)
	if !ok {
		return 0, fmt.Errorf("not a stat_t")
	}
	return stat.Ino, nil
}

// FindFileInProc Traverse /proc/<pid>/<suffix> files,
// compare their inodes with inode parameter and returns file if inode matches
// use FindProcInode(xxx, "/ns/net") for example
func FindFileInProc(inode uint64, suffix string) (string, error) {
	files, err := ioutil.ReadDir("/proc")
	if err != nil {
		return "", fmt.Errorf("can't read /proc directory: %+v", err)
	}

	for _, f := range files {
		name := f.Name()
		if isDigits(name) {
			filename := "/proc/" + name + suffix
			tryInode, err := GetInode(filename)
			if err == nil {
				if tryInode == inode {
					return filename, nil
				}
			}
		}
	}

	return "", fmt.Errorf("not found")
}
