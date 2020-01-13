package fspatterns

import (
	"bytes"
	"fmt"
	serviceutils "github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
	"github.com/jfrog/jfrog-client-go/utils/io/fileutils/checksum"
	"os"
	"regexp"
	"strings"
)

// Return all the existing paths of the provided root path
func GetPaths(rootPath string, isRecursive, includeDirs, isSymlink bool) ([]string, error) {
	var paths []string
	var err error
	if isRecursive {
		paths, err = fileutils.ListFilesRecursiveWalkIntoDirSymlink(rootPath, !isSymlink)
	} else {
		paths, err = fileutils.ListFiles(rootPath, includeDirs)
	}
	if err != nil {
		return paths, err
	}
	return paths, nil
}

// Transform to regexp and prepare Exclude patterns to be used
func PrepareExcludePathPattern(params serviceutils.FileGetter) string {
	excludePathPattern := ""
	if len(params.GetExcludePatterns()) > 0 {
		for _, singleExcludePattern := range params.GetExcludePatterns() {
			if len(singleExcludePattern) > 0 {
				singleExcludePattern = utils.ReplaceTildeWithUserHome(singleExcludePattern)
				singleExcludePattern = utils.PrepareLocalPathForUpload(singleExcludePattern, params.IsRegexp())
				if params.IsRecursive() && strings.HasSuffix(singleExcludePattern, fileutils.GetFileSeparator()) {
					singleExcludePattern += "*"
				}
				excludePathPattern += fmt.Sprintf(`(%s)|`, singleExcludePattern)
			}
		}
		if len(excludePathPattern) > 0 {
			excludePathPattern = excludePathPattern[:len(excludePathPattern)-1]
		}
	}
	return excludePathPattern
}

// Return only subpaths of the provided by the user path that matched to the provided regexp.
// Subpaths that matched to an exclude pattern won't returned
func PrepareAndFilterPaths(path, excludePathPattern string, preserveSymlinks, includeDirs bool, regexp *regexp.Regexp) (matches []string, isDir, isSymlinkFlow bool, err error) {
	isDir, err = fileutils.IsDirExists(path, false)
	if err != nil {
		return
	}

	excludedPath, err := IsPathExcluded(path, excludePathPattern)
	if err != nil {
		return
	}

	if excludedPath {
		return
	}
	isSymlinkFlow = preserveSymlinks && fileutils.IsPathSymlink(path)

	if isDir && !includeDirs && !isSymlinkFlow {
		return
	}
	matches = regexp.FindStringSubmatch(path)
	return
}

func GetSingleFileToUpload(rootPath, targetPath string, flat, preserveSymLink bool) (utils.Artifact, error) {
	symlinkPath, err := GetFileSymlinkPath(rootPath)
	if err != nil {
		return utils.Artifact{}, err
	}

	var uploadPath string
	if !strings.HasSuffix(targetPath, "/") {
		uploadPath = targetPath
	} else {
		// If not preserving symlinks and symlink target is valid, use symlink target for upload
		if !preserveSymLink && symlinkPath != "" {
			rootPath = symlinkPath
		}
		if flat {
			uploadPath, _ = fileutils.GetFileAndDirFromPath(rootPath)
			uploadPath = targetPath + uploadPath
		} else {
			uploadPath = targetPath + rootPath
			uploadPath = utils.TrimPath(uploadPath)
		}
	}

	return utils.Artifact{LocalPath: rootPath, TargetPath: uploadPath, Symlink: symlinkPath}, nil
}

func IsPathExcluded(path string, excludePathPattern string) (excludedPath bool, err error) {
	if len(excludePathPattern) > 0 {
		excludedPath, err = regexp.MatchString(excludePathPattern, path)
	}
	return
}

// If filePath is path to a symlink we should return the link content e.g where the link points
func GetFileSymlinkPath(filePath string) (string, error) {
	fileInfo, e := os.Lstat(filePath)
	if e != nil {
		return "", errorutils.WrapError(e)
	}
	var symlinkPath = ""
	if fileutils.IsFileSymlink(fileInfo) {
		symlinkPath, e = os.Readlink(filePath)
		if e != nil {
			return "", errorutils.WrapError(e)
		}
	}
	return symlinkPath, nil
}

// Get the local root path, from which to start collecting artifacts to be uploaded to Artifactory.
// If path dose not exist error will be returned.
func GetRootPath(pattern string, isRegexp, preserveSymLink bool) (string, error) {
	rootPath := utils.GetRootPath(pattern, isRegexp)

	if !fileutils.IsPathExists(rootPath, preserveSymLink) {
		return "", errorutils.NewError("Path does not exist: " + rootPath)
	}

	return rootPath, nil
}

// When handling symlink we want to simulate the creation of empty file
func CreateSymlinkFileDetails() (*fileutils.FileDetails, error) {
	checksumInfo, err := checksum.Calc(bytes.NewBuffer([]byte(fileutils.SYMLINK_FILE_CONTENT)))
	if err != nil {
		return nil, err
	}

	details := new(fileutils.FileDetails)
	details.Checksum.Md5 = checksumInfo[checksum.MD5]
	details.Checksum.Sha1 = checksumInfo[checksum.SHA1]
	details.Checksum.Sha256 = checksumInfo[checksum.SHA256]
	details.Size = int64(0)
	return details, nil
}
