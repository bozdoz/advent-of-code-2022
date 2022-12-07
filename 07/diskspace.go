package main

import (
	"fmt"
	"strings"

	"github.com/bozdoz/advent-of-code-2022/utils"
)

// a file entry
type entry struct {
	isFile, isDir bool
	name          string
	size          int
	entries       *map[string]*entry
}

// keeping track of cwd and parent directories
type dir []*entry

func parseHistory(data inType) (root *entry) {
	// initialize with "/"
	cur := newDirectory("/")

	var cwd = dir([]*entry{cur})

	for _, line := range data {
		switch {
		case strings.HasPrefix(line, "$ cd"):
			// check size of this dir, before we change
			updateDirSize(cur)

			// move to new directory
			cur = cwd.move(line[len("$ cd "):])
		case line == "$ ls":
			// ignore
			continue
		default:
			// we are updating entries in cwd
			addEntry(cur, line)
		}
	}

	// size remaining directories
	traverseAndSizeDirectories(cwd[0])

	// return root entry
	return cwd[0]
}

// current working DIRECTORY
func (d *dir) current() *entry {
	return (*d)[len(*d)-1]
}

func (d *dir) move(dir string) *entry {
	switch dir {
	case "/":
		// root entry (it only happens at the beginning)
		*d = (*d)[0:1]
	case "..":
		// move up/remove last entry
		old := *d
		n := len(old)
		*d = old[0 : n-1]
	default:
		// push new entry

		// does entry exist in current entries?
		ent, ok := (*d.current().entries)[dir]

		if ok {
			*d = append(*d, ent)
		} else {
			*d = append(*d, newDirectory(dir))
		}
	}

	return d.current()
}

func newDirectory(name string) *entry {
	return &entry{
		name:    name,
		isDir:   true,
		entries: &map[string]*entry{},
	}
}

func addEntry(cwd *entry, line string) {
	parts := strings.Fields(line)
	first, name := parts[0], parts[1]

	ent := &entry{
		name: name,
	}

	if first == "dir" {
		ent.isDir = true
		ent.entries = &map[string]*entry{}
	} else {
		ent.isFile = true
		ent.size = utils.ParseInt(first)
	}

	// add to tree
	(*cwd.entries)[name] = ent
}

func updateDirSize(cwd *entry) {
	if cwd.size != 0 {
		// ignore sized directories
		return
	}

	size := 0

	for _, ent := range *cwd.entries {
		if ent.isDir && ent.size == 0 {
			// assuming no empty dirs;
			// this dir hasn't been sized yet
			return
		}

		size += ent.size
	}

	cwd.size = size
}

// probably only one directory will remain unsized
// do to exiting the program on an 'ls' instead of a 'cd',
// which triggers the other directory sizing
// this also sizes the root directory, as an added bonus
func traverseAndSizeDirectories(cwd *entry) {
	for _, ent := range *cwd.entries {
		if ent.isDir && ent.size == 0 {
			// enter in, and size
			traverseAndSizeDirectories(ent)
		}
	}
	updateDirSize(cwd)
}

func traverseDirectories(cwd *entry, fun func(d *entry)) {
	fun(cwd)

	for _, ent := range *cwd.entries {
		if ent.isDir {
			traverseDirectories(ent, fun)
		}
	}
}

//
// string representations for debugging
//

func (e *entry) String() (out string) {
	return e.toString(0)
}

func (e *entry) toString(tabs int) (out string) {
	tab := strings.Repeat("  ", tabs)

	out += fmt.Sprint(tab, "- ", e.name)

	if e.isFile {
		out += fmt.Sprintf(" (file, size=%d)\n", e.size)
	} else {
		out += fmt.Sprintf(" (dir, size=%d)\n", e.size)
		for _, file := range *e.entries {
			out += file.toString(tabs + 1)
		}
	}

	return
}
