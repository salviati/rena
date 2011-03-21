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

/*
  rena(1) - Rename Archive. Tool crafted to rename files in
  anime/manga archives/dirs.
  Also can be used to list missing/duplicate episodes.
*/
package main

import (
	"path/filepath"
	"os"
	"regexp"
	"flag"
	"fmt"
	"sort"
	"strings"
	"log"
)

const (
	pkg, version, author, about, usage string = "rena", "20101231", "Utkan Güngördü",
		"rena(1) - Rename Archive. Tool crafted to rename files in anime/manga archives/dirs. Also can be used to list missing/duplicate episodes.",
		"rena [options] file1/dir1 [file2/dir2 ... fileN/dirN]"
)

var templ = flag.String("t", "{n}", "Name template. Ex: Naruto-{N}. {N} will be replaced with episode number(s) (uses template.Parse). The extension of the old file will be prepended.")
var recurse = flag.Bool("r", false, "Recurse into subdirectories.")
var eseperator = flag.String("s", "-", "Episode seperator.")
var nDigits = flag.Int("N", 0, "Number of digits for (zero padded) episode numbers.")
var yesToAll = flag.Bool("y", false, "Yes to all.")
var noToAll = flag.Bool("n", false, "No to all.")
var showVersion = flag.Bool("version", false, "Show version info and quit.")
var showHelp = flag.Bool("h", false, "Display this message.")
var chopRegexp = flag.String("C", "", "Crop what matches to given regexp. Ex. usage: ^FMA2 to get rid of the troublesome numerical prefix 'FMA2'.")


const MAXEPS = 1e4

var filters []*regexp.Regexp
var episodes = make([]*Episode, 0, MAXEPS)
var epsFormat = "%d" // Format string used to convert episode numbers into strings

var chop = []string{
	strings.Repeat("[0-9A-Fa-f]", 8), // CRC32
	"\\[" + strings.Repeat("[0-9A-Fa-f]", 8) + "\\]", // CRC32
	"360[pP]",
	"480[pP]",
	"576[pP]",
	"720[pP]",
	"1080[pP]",
	"4320[pP]",
}

func init() {
	flag.Parse()
	chop = append(chop, *chopRegexp)

	if *showVersion {
		printVersion(pkg, version, author)
		os.Exit(0)
	}
	if *showHelp || flag.NArg() == 0 {
		printHelp(pkg, version, about, usage)
		os.Exit(0)
	}

	if *nDigits != 0 {
		epsFormat = fmt.Sprintf("%%0%dd", *nDigits)
	}

	if *yesToAll && *noToAll {
		log.Fatal("Yes or No, make up your mind!\n")
	}

	filters = make([]*regexp.Regexp, len(chop))
	for i := 0; i < len(chop); i++ {
		filters[i] = regexp.MustCompile(chop[i])
	}
}


type walkEnt struct{}

func (*walkEnt) VisitDir(dname string, d *os.FileInfo) bool {
	return *recurse
}

func (*walkEnt) VisitFile(fpath string, d *os.FileInfo) {
	episodes = append(episodes, NewEpisode(fpath))
}

/* Recursively renames files under dirName. */
func rena(dirName string) {
	_, err := os.Stat(dirName)
	if err != nil {
		log.Println(err)
		return
	}

	v := new(walkEnt)
	ech := make(chan os.Error)
	go func() { filepath.Walk(dirName, v, ech); close(ech) }()
	for e := range ech {
		log.Println(e)
	}
}


func main() {
	for i := 0; i < flag.NArg(); i++ {
		rena(flag.Arg(i))
	}

	alleps := make([]int, 0, MAXEPS)

	// Do renaming for all episodes
	for _, e := range episodes {
		newpath := e.String()
		oldpath := e.OldPath()
		if newpath == "" {
			log.Println("[IGNORING] ", oldpath)
			continue
		}

		episodeNumbers := e.GetEpisodeNumbers()
		for e := range episodeNumbers {
			alleps = append(alleps, e)
		}

		log.Println("[INFO]", oldpath, "-> ", newpath)
		if oldpath == newpath {
			log.Println("[LOG] Target and source are identical, no action needed.")
			continue
		}

		if ynQuestion(fmt.Sprintln("Rename", oldpath, "-> ", newpath)) {
			log.Printf("[LOG] Renaming %s -> %s\n", oldpath, newpath)
			err := os.Rename(oldpath, newpath)
			if err != nil {
				log.Fatal(err, "\n")
			}
		}

	}

	sort.SortInts(alleps)
	nmissing := 0
	for i := 0; i < len(alleps)-1; i++ {
		if alleps[i+1]-alleps[i] == 1 {
			continue
		}
		if alleps[i+1] == alleps[i] {
			log.Println("[EPSSCAN] Duplicates exists for episode:", alleps[i])
			continue
		}
		for j := alleps[i] + 1; j < alleps[i+1]; j++ {
			log.Println("[EPSSCAN] Missing episode", j)
			nmissing++
		}
	}

	log.Println("Missing", nmissing, "episodes in total.")
}
