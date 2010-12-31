/*
   Copyright (c) Utkan Güngördü <utkan@freeconsole.org>

   This program is free software; you can redistribute it and/or modify
   it under the terms of the GNU General Public License as
   published by the Free Software Foundation; either version 3 or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details


   You should have received a copy of the GNU General Public
   License along with this program; if not, write to the
   Free Software Foundation, Inc.,
   51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
*/

// Uses templ (for templating the new name) and eseperator (separation string
// for episode numbers in the output) globals.

package main

import (
	"path"
	"strconv"
	"regexp"
	"template"
	"fmt"
	"bytes"
)

// Extract episode numbers from filename, starting from the first numeral
// and ending at the last consecutive number.
// Ex: "Shippuuden_-_178-179-200.mkv" -> []int{178, 179}
func getEpisodeNumbers(filename string) []int {
	r := regexp.MustCompile("[0-9]+") //FIXME: make this a global
	matches := r.FindAllString(filename, -1)
	var eps = make([]int, 0, len(matches))

	for _, m := range matches {
		newe, _ := strconv.Atoi(m)
		if len(eps) > 1 && eps[len(eps)-1]+1 != newe {
			break
		}
		eps = append(eps, newe)
	}

	return eps
}

type Episode struct {
	dir              string // Directory of the original file.
	oldname, newname string
	episodeNumbers   []int // Array of episode numbers
}

func NewEpisode(filepath string) *Episode {
	e := new(Episode)
	e.dir, e.oldname = path.Split(filepath)
	e.Rename()
	return e
}

func (e *Episode) String() string {
	if e.newname == "" {
		return ""
	}
	return path.Join(e.dir, e.newname)
}


type nameTemplate struct {
	n string
}

func (e *Episode) Rename() {
	e.newname = e.oldname
	for _, r := range filters {
		e.newname = r.ReplaceAllString(e.newname, "")
	}

	e.episodeNumbers = getEpisodeNumbers(e.newname)

	if len(e.episodeNumbers) == 0 {
		e.newname = ""
		return
	}

	estr := ""
	for i, en := range e.episodeNumbers {
		estr += fmt.Sprintf(epsFormat, en)
		if len(e.episodeNumbers) > i+1 {
			estr += *eseperator
		}
	}

	t := template.MustParse(*templ, nil)
	w := bytes.NewBufferString("")
	t.Execute(&nameTemplate{n: estr}, w)
	e.newname = string(w.Bytes()) + path.Ext(e.oldname)
}

func (e *Episode) GetEpisodeNumbers() (r []int) {
	copy(r, e.episodeNumbers) //FIXME
	return
}

func (e *Episode) OldPath() string {
	return path.Join(e.dir, e.oldname)
}