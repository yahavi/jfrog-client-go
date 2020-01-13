package packages

import (
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"strings"
)

func CreatePath(packageStr string) (*Path, error) {
	parts := strings.Split(packageStr, "/")
	size := len(parts)
	if size != 3 {
		err := errorutils.NewError("Expecting an argument in the form of subject/repository/package")
		if err != nil {
			return nil, err
		}
	}
	return &Path{
		Subject: parts[0],
		Repo:    parts[1],
		Package: parts[2]}, nil
}
