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
	"flag"
)

func printVersion(pkg string, version string, author string) {
	fmt.Println(pkg, version)
	fmt.Println("Copyright (C) 2010", author)
	fmt.Println("This program is free software; you may redistribute it under the terms of")
	fmt.Println("the GNU General Public License version 3 or (at your option) a later version.")
	fmt.Println("This program has absolutely no warranty.")
	fmt.Println("Report bugs to bug@freeconsole.org")
}

func printHelp(pkg string, version string, about string, usage string) {
	fmt.Println(pkg, version, "\n")
	fmt.Println(about)
	fmt.Println("Usage:")
	fmt.Println("\t", usage)
	fmt.Println("Options:")
	flag.PrintDefaults()
}
