package tsh

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Version contains the tsh Version
type Version struct {
	Major     int
	Minor     int
	Patch     int
	GoVersion string
}

// NewVersion create a Version from string
func NewVersion(s string) (*Version, error) {
	split := strings.Split(s, " ")
	if len(split) < 2 {
		return nil, fmt.Errorf("not enough Version string")
	}

	// ensure the Version tag has `v`
	if !strings.Contains(split[1], "v") {
		return nil, fmt.Errorf("invalid Version")
	}

	tag := strings.Replace(split[1], "v", "", -1)

	// ignores aplha, beta & rc
	compile, err := regexp.Compile("-(alpha|beta|rc).*")
	if err != nil {
		return nil, err
	}
	tag = compile.ReplaceAllString(tag, "")
	numbers := strings.Split(tag, ".")
	if len(numbers) < 3 {
		return nil, fmt.Errorf("%s is invalid", tag)
	}

	v := &Version{
		Major: atoi(numbers[0]),
		Minor: atoi(numbers[1]),
		Patch: atoi(numbers[2]),
	}
	return v, nil
}

func atoi(s string) (i int) {
	i, _ = strconv.Atoi(s)
	return
}

// Strings return string formatted
func (v *Version) Strings() string {
	return fmt.Sprintf("Teleport v%d.%d.%d", v.Major, v.Minor, v.Patch)
}

// IsSupported return weather the current Version is supported
func (v *Version) IsSupported(cv *Version) bool {
	if cv.Major > v.Major {
		return true
	}
	if cv.Major == v.Major && cv.Minor > v.Minor {
		return true
	}
	return cv.Major == v.Major && cv.Minor == v.Minor && cv.Patch >= v.Patch
}
