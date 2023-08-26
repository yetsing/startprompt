package main

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

/*
自动给当前提交打上下一个版本号的 tag
*/

func toInt(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		log.Fatal(err)
	}
	return n
}

type Version struct {
	major int
	minor int
	patch int
}

func (v *Version) format() string {
	return fmt.Sprintf("v%d.%d.%d", v.major, v.minor, v.patch)
}

func main() {
	out, err := exec.Command("git", "tag").Output()
	if err != nil {
		log.Fatal(err)
	}
	var latest Version
	tags := strings.Split(string(out), "\n")
	for _, tag := range tags {
		if !strings.HasPrefix(tag, "v") {
			continue
		}
		parts := strings.Split(tag[1:], ".")
		if len(parts) != 3 {
			log.Fatalf("invalid version tag: %q", tag)
		}
		ver := Version{
			major: toInt(parts[0]),
			minor: toInt(parts[1]),
			patch: toInt(parts[2]),
		}
		if ver.major > latest.major {
			latest = ver
		} else if ver.minor > latest.minor {
			latest = ver
		} else if ver.patch > latest.patch {
			latest = ver
		}
	}
	latest.patch++
	fmt.Printf("next: %s\n", latest.format())
	out, err = exec.Command("git", "tag", latest.format()).Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(out))
}
