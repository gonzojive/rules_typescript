// Package runfiles that provides utility helpers for resolving Bazel runfiles within Go.
package runfiles

import (
	"fmt"

	"github.com/bazelbuild/rules_go/go/tools/bazel"
)

// Runfile returns the base directory to the bazel runfiles
func Runfile(_base string, manifestPath string) (string, error) {
	got, err := bazel.Runfile(manifestPath)
	if err != nil {
		closest, closestErr := findClosestRunfileEntry(manifestPath)
		if closestErr == nil && closest != nil {
			return "", fmt.Errorf("manifest file %s not found, did you mean %s?", manifestPath, closest.ShortPath)
		}
		runfilesDir, runfilesErr := bazel.RunfilesPath()
		if runfilesErr == nil {
			return "", fmt.Errorf("manifest file %q not found in runfile directory %q", manifestPath, runfilesDir)
		}
		return "", err
	}
	return got, nil
}

func findClosestRunfileEntry(manifestPath string) (*bazel.RunfileEntry, error) {
	rfes, err := bazel.ListRunfiles()
	if err != nil {
		return nil, err
	}
	var best *bazel.RunfileEntry
	longestSharedSuffix := ""
	for _, rfe := range rfes {
		if suf := sharedSuffix(manifestPath, rfe.Path); len(suf) > len(longestSharedSuffix) {
			longestSharedSuffix = suf
			rfe := rfe
			best = &rfe
			if len(suf) == len(manifestPath) {
				break
			}
		}
	}
	if len(manifestPath) > 0 && (float64(len(longestSharedSuffix))/float64(len(manifestPath))) > .40 {
		return best, nil
	}
	return nil, nil
}

func sharedSuffix(a, b string) string {
	posA, posB := len(a), len(b)
	for {
		if posA == 0 || posB == 0 {
			break
		}
		if a[posA-1] != b[posB-1] {
			break
		}
		posA--
		posB--
	}
	return a[posA:]
}
