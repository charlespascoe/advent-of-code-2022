package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	lines, err := readLines("input.txt")
	if err != nil {
		log.Fatalf("Couldn't read input file: %s", err)
	}

	root := &Directory{}
	dir := root

	cdRe := regexp.MustCompile(`^\$ cd (.*)`)
	lsRe := regexp.MustCompile(`^\$ ls`)
	fileRe := regexp.MustCompile(`^(\d+)`)

	for _, line := range lines {
		cdMatch := cdRe.FindStringSubmatch(line)
		lsMatch := lsRe.MatchString(line)
		fileMatch := fileRe.FindStringSubmatch(line)

		switch {
		case cdMatch != nil:
			if cdMatch[1] == ".." {
				dir = dir.Up()
			} else if cdMatch[1] == "/" {
				dir = root
			} else {
				dir = dir.Down(cdMatch[1])
			}

		case lsMatch:
			// Reset Size to prevent multiple `ls` commands in a directory from
			// counting the content size twice
			dir.ContentSize = 0

		case fileMatch != nil:
			size, err := strconv.Atoi(fileMatch[1])
			if err != nil {
				// This should never happen
				panic(err)
			}

			dir.ContentSize += size
		}
	}

	root.CalcTotals()

	sum := 0

	root.VisitAll(func(dir *Directory) {
		if dir.TotalSize <= 100_000 {
			sum += dir.TotalSize
		}
	})

	fmt.Printf("Sum of directories <= 100kb: %d\n", sum)

	// Part 2

	diskSize := 70_000_000
	freeSpace := diskSize - root.TotalSize
	requiredSpace := 30_000_000

	best := root

	root.VisitAll(func(dir *Directory) {
		if (freeSpace+dir.TotalSize) > requiredSpace && dir.TotalSize < best.TotalSize {
			best = dir
		}
	})

	fmt.Printf("Best directory to delete: %s\nSpace freed: %d\n", best.Name, best.TotalSize)
}

func readLines(path string) ([]string, error) {
	if data, err := os.ReadFile(path); err != nil {
		return nil, err
	} else {
		return strings.Split(strings.TrimSpace(string(data)), "\n"), nil
	}
}

type Directory struct {
	Name     string
	Children map[string]*Directory
	Parent   *Directory

	// ContentSize is the size of the files immediately within this directory
	ContentSize int

	// TotalSize is the size of all files in the directory and child directories
	TotalSize int
}

func (d *Directory) Down(name string) *Directory {
	if d.Children == nil {
		d.Children = make(map[string]*Directory)
	}

	child := d.Children[name]

	if child == nil {
		child = &Directory{
			Name:   name,
			Parent: d,
		}

		d.Children[name] = child
	}

	return child
}

func (d *Directory) Up() *Directory {
	if d.Parent == nil {
		panic("can't go above root")
	}

	return d.Parent
}

func (d *Directory) CalcTotals() {
	d.TotalSize = d.ContentSize

	for _, child := range d.Children {
		child.CalcTotals()
		d.TotalSize += child.TotalSize
	}
}

func (d *Directory) VisitAll(f func(*Directory)) {
	f(d)

	for _, child := range d.Children {
		child.VisitAll(f)
	}
}
