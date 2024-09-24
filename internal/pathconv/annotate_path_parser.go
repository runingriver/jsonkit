package pathconv

import (
	"strconv"
	"strings"

	"github.com/jsonkit/jkerr"
)

type KeyType int

const (
	StrKeyType = 1
	NumKeyType = 2
	AllKeyType = 3
)

const (
	Asterisk = "*"
)

type Path struct {
	KeyType KeyType
	Key     string
	Idx     int
}

func NewPath(kt KeyType, k string, i int) *Path {
	return &Path{KeyType: kt, Key: k, Idx: i}
}

func DotPathToSlice(path string) []*Path {
	hierarchy := strings.Split(path, ".")
	pathList := make([]*Path, 0, len(hierarchy))
	for _, pathSeg := range hierarchy {
		if pathSeg == Asterisk {
			pathList = append(pathList, NewPath(AllKeyType, pathSeg, 0))
			continue
		}
		if idx, err := strconv.Atoi(pathSeg); err == nil {
			pathList = append(pathList, NewPath(NumKeyType, pathSeg, idx))
			continue
		}
		pathList = append(pathList, NewPath(StrKeyType, pathSeg, 0))
	}
	return pathList
}

func ValidNormalPath(path string) error {
	if err := ValidPath(path); err != nil {
		return err
	}
	if strings.HasSuffix(path, ".*") {
		return jkerr.New(jkerr.PathIllegalErr, "path cannot end with *")
	}
	return nil
}

func ValidPath(path string) error {
	if path == "" {
		return jkerr.New(jkerr.PathIllegalErr, "path cannot be empty")
	}
	if strings.HasPrefix(path, "*.") {
		return jkerr.New(jkerr.PathIllegalErr, "path cannot start with *")
	}
	if strings.Contains(path, ".*.") {
		return jkerr.New(jkerr.PathIllegalErr, "path search un-support *")
	}
	return nil
}
