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

package main

import (
	"fmt"
	"os"
	"log"
	"strings"
)

/* Presents a yes/no question. Returns 1 if yes, 0 otherwise. */
func ynQuestion(format string, va ...interface{}) bool {
	if *yesToAll {
		return true
	}

	if *noToAll {
		return false
	}

	for {
		fmt.Fprintf(os.Stderr, format, va...)
		fmt.Fprintf(os.Stderr, " (y/n): ")
		r := ""
		_, err := fmt.Scanf("%s", &r)
		if err != nil {
			log.Exit(err)
		}
		r = strings.TrimSpace(r)

		switch r {
		case "y", "Y":
			return true
		case "n", "N":
			return false
		default:
			fmt.Fprintf(os.Stderr, "hint: say y or n :)\n")
		}
	}
	return false
}
