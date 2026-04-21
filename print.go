/*
 *          Copyright 2026 Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package main

import (
	"fmt"
)

func printInfo() {
	message := "USAGE\n"
	message += "    file-cmd ( INFO | HELP | EXAMPLE | COMMAND )\n"
	message += "INFO\n"
	message += "    -v, --version     print version\n"
	message += "    -c, --copyright   print copyright\n"
	message += "HELP\n"
	message += "    ( -h | --help ) [COMMAND]\n"
	message += "               print this help, or print help for COMMAND\n"
	message += "EXAMPLE\n"
	message += "    ( -e | --example ) [COMMAND]\n"
	message += "               print available examples, or print expample for COMMAND\n"
	message += "COMMAND\n"
	message += "    split      split file into multiple files\n"
	message += "    concat     concatenate files\n"
	message += "    ls         list files\n"
	message += "    count      count files\n"
	message += "    cp         copy files\n"
	message += "    mv         move files\n"
	message += "    rm         delete files\n"
	message += "    clean      remove empty folders and files\n"
	message += "    text       generate random text"
	fmt.Println(message)
}

func printInfoExample() {
	message := "USAGE\n"
	message += "    file-cmd ( -e | --example ) [COMMAND]\n"
	message += "               print this help, or print expample for COMMAND\n"
	message += "COMMAND\n"
	message += "    split      split file into multiple files\n"
	message += "    concat     concatenate files\n"
	message += "    ls         list files\n"
	message += "    count      count files\n"
	message += "    cp         copy files\n"
	message += "    mv         move files\n"
	message += "    rm         delete files\n"
	message += "    clean      remove empty folders and files\n"
	message += "    text       generate random text"
	fmt.Println(message)
}

func printInfoSplit() {
	message := "USAGE\n"
	message += "    file-cmd split INPUT-FILE [OUTPUT-FILE] [OPTION]\n"
	message += "               split file into multiple files\n"
	message += "INPUT-FILE\n"
	message += "    [-i | --input] <file-path>\n"
	message += "OUTPUT-FILE\n"
	message += "    [-o | --output] <file-path>\n"
	message += "OPTION\n"
	message += "    -p=N       split file into N parts (chunks)\n"
	message += "    -b=N[U]    split file into N bytes per chunk, U = unit (k/K, m/M or g/G)\n"
	message += "    -l=N       split file into N lines per chunk\n"
	message += "    -w         overwrite output files"
	fmt.Println(message)
}

func printInfoConcat() {
	message := "USAGE\n"
	message += "    file-cmd concat INPUT-FILE [OUTPUT-FILE] [OPTION]\n"
	message += "               concatenate files into one file\n"
	message += "INPUT-FILE\n"
	message += "    [-i | --input] <file-path>\n"
	message += "               input file is only one file - the first one\n"
	message += "               concatenated files must be numbered (see: --example concat)\n"
	message += "OUTPUT-FILE\n"
	message += "    [-o | --output] <file-path>\n"
	message += "OPTION\n"
	message += "    -w         overwrite output file"
	fmt.Println(message)
}

func printInfoList(cmdStr string) {
	message := "USAGE\n"
	message += "    file-cmd " + cmdStr + " INPUT-DIR {OPTION} {FILTER}\n"
	message += "                      print files filtered by their content\n"
	message += "INPUT-DIR\n"
	message += "    [-i | --input] <directory-path>\n"
	message += "OPTION\n"
	message += "    --or              filter is OR (not AND)\n"
	message += "    -r, --recursive   recursive file iteration\n"
	message += "    -v, --verbose     output warnings to screen\n"
	message += "    -d=D              set delimiter D; default is comma, else space\n"
	message += "                      empty D is no delimiter\n"
	message += "    -t=N              use N threads\n"
	message += "FILTER\n"
	message += "    ( -f=W | W )      filter files by strings W; W is divided by delimiter"
	fmt.Println(message)
}

func printInfoCount(cmdStr string) {
	message := "USAGE\n"
	message += "    file-cmd " + cmdStr + " INPUT-DIR {OPTION} {FILTER}\n"
	message += "                      print number of files matched by filter\n"
	message += "INPUT-DIR\n"
	message += "    [-i | --input] <directory-path>\n"
	message += "OPTION\n"
	message += "    --or              filter is OR (not AND)\n"
	message += "    -r, --recursive   recursive file iteration\n"
	message += "    -v, --verbose     output warnings to screen\n"
	message += "    -t=N              use N threads\n"
	message += "FILTER\n"
	message += "    ( -f=W | W )      filter files by strings W; space or comma are separators"
	fmt.Println(message)
}

func printInfoCopy(cmdStr string) {
	message := "USAGE\n"
	message += "    file-cmd " + cmdStr + " INPUT-DIR OUTPUT-DIR {OPTION} {FILTER}\n"
	message += "                      copy files matched by filter\n"
	message += "INPUT-DIR\n"
	message += "    [-i | --input] <directory-path>\n"
	message += "OUTPUT-DIR\n"
	message += "    [-o | --output] <directory-path>\n"
	message += "OPTION\n"
	message += "    --or              filter is OR (not AND)\n"
	message += "    -w                overwrite output files\n"
	message += "    -r, --recursive   recursive file iteration\n"
	message += "    -v, --verbose     output warnings to screen\n"
	message += "    -t=N              use N threads\n"
	message += "FILTER\n"
	message += "    ( -f=W | W )      filter files by strings W; space or comma are separators"
	fmt.Println(message)
}

func printInfoMove(cmdStr string) {
	message := "USAGE\n"
	message += "    file-cmd " + cmdStr + " INPUT-DIR OUTPUT-DIR {OPTION} {FILTER}\n"
	message += "                      move files matched by filter\n"
	message += "INPUT-DIR\n"
	message += "    [-i | --input] <directory-path>\n"
	message += "OUTPUT-DIR\n"
	message += "    [-o | --output] <directory-path>\n"
	message += "OPTION\n"
	message += "    --or              filter is OR (not AND)\n"
	message += "    -w                overwrite output files\n"
	message += "    -r, --recursive   recursive file iteration\n"
	message += "    -v, --verbose     output warnings to screen\n"
	message += "    -t=N              use N threads\n"
	message += "FILTER\n"
	message += "    ( -f=W | W )      filter files by strings W; space or comma are separators"
	fmt.Println(message)
}

func printInfoRemove(cmdStr string) {
	message := "USAGE\n"
	message += "    file-cmd " + cmdStr + " INPUT-DIR {OPTION} {FILTER}\n"
	message += "                      delete files matched by filter\n"
	message += "INPUT-DIR\n"
	message += "    [-i | --input] <directory-path>\n"
	message += "OPTION\n"
	message += "    --or              filter is OR (not AND)\n"
	message += "    -r, --recursive   recursive file iteration\n"
	message += "    -v, --verbose     output warnings to screen\n"
	message += "FILTER\n"
	message += "    ( -f=W | W )      filter files by strings W; space or comma are separators"
	fmt.Println(message)
}

func printInfoClean(cmdStr string) {
	message := "USAGE\n"
	message += "    file-cmd " + cmdStr + " INPUT-DIR {OPTION}\n"
	message += "                      delete empty folders and empty regular files\n"
	message += "INPUT-DIR\n"
	message += "    [-i | --input] <directory-path>\n"
	message += "OPTION\n"
	message += "    -a, --all         delete also directories\n"
	message += "    -r, --recursive   recursive file iteration\n"
	message += "    -v, --verbose     output warnings to screen\n"
	message += "    -l, --list        print deleted"
	fmt.Println(message)
}

func printInfoText() {
	message := "USAGE\n"
	message += "    file-cmd text [OUTPUT-FILE] {OPTION}\n"
	message += "                      generate random text\n"
	message += "OUTPUT-FILE\n"
	message += "    [-o | --output] <file-path>\n"
	message += "                      standard output is default\n"
	message += "OPTION\n"
	message += "    -w                overwrite output file\n"
	message += "    -s=N[U]           size N of file, U = unit (k/K, m/M or g/G)\n"
	message += "    -e=E              use line terminator E (e.g. CRLF or LF or CR)\n"
	message += "    -d=D              use word delimiter D (space is default)\n"
	message += "    -f={a|l|u}        output format\n"
	message += "                        a = letters, only\n"
	message += "                        l = lower case, only\n"
	message += "                        u = upper case, only"
	fmt.Println(message)
}

func printVersion() {
	message := "0.1.0"
	fmt.Println(message)
}

func printExampleSplit() {
	message := "EXAMPLE A\n"
	message += "    split file tmp.txt into files: tmp.txt.0 and tmp.txt.1\n"
	message += "    $ file-cmd split -p=2 ./tmp.txt\n"
	message += "EXAMPLE B\n"
	message += "    split file tmp.txt into files: xyz.0 and xyz.1\n"
	message += "    $ file-cmd split -p=2 ./tmp.txt ./xyz\n"
	message += "EXAMPLE C\n"
	message += "    split file tmp.txt into files 10 kilobytes each\n"
	message += "    $ file-cmd split -b=10k ./tmp.txt\n"
	message += "EXAMPLE D\n"
	message += "    split file tmp.txt into files 10 kibibytes each\n"
	message += "    $ file-cmd split -b=10K ./tmp.txt\n"
	message += "EXAMPLE E\n"
	message += "    split file tmp.txt into files 100 bytes each\n"
	message += "    $ file-cmd split -b=100 ./tmp.txt"
	fmt.Println(message)
}

func printExampleConcat() {
	message := "EXAMPLE A\n"
	message += "    concatenate files tmp.txt.0 and tmp.txt.1 to tmp.txt\n"
	message += "    $ file-cmd concat ./tmp.txt\n"
	message += "EXAMPLE B\n"
	message += "    concatenate files tmp.txt.0 and tmp.txt.1 to tmp.txt\n"
	message += "    $ file-cmd concat ./tmp.txt.0\n"
	message += "EXAMPLE C\n"
	message += "    concatenate files tmp.txt.0 and tmp.txt.1 to xyz.txt\n"
	message += "    $ file-cmd concat ./tmp.txt ./xyz.txt"
	fmt.Println(message)
}

func printExampleList() {
	message := "EXAMPLE A\n"
	message += "    list files containing \"alice\" and \"bob\"\n"
	message += "    $ file-cmd ls . -f=alice,bob\n"
	message += "EXAMPLE B\n"
	message += "    list files containing \"alice\" and \"bob\"\n"
	message += "    $ file-cmd ls . alice bob\n"
	message += "EXAMPLE C\n"
	message += "    list files containing \"alice\" or \"bob\"\n"
	message += "    $ file-cmd ls . --or alice bob"
	fmt.Println(message)
}

func printExampleCount() {
	message := "EXAMPLE A\n"
	message += "    count files containing \"alice\" and \"bob\" in current directory\n"
	message += "    $ file-cmd count . -f=alice,bob\n"
	message += "EXAMPLE B\n"
	message += "    count files containing \"alice\" and \"bob\" in current directory\n"
	message += "    $ file-cmd count . alice bob\n"
	message += "EXAMPLE C\n"
	message += "    count files containing \"alice\" or \"bob\" in current directory\n"
	message += "    $ file-cmd count . --or alice bob\n"
	message += "EXAMPLE D\n"
	message += "    count all files in current directory and subdirectories\n"
	message += "    $ file-cmd count . -r"
	fmt.Println(message)
}

func printExampleCopy() {
	message := "EXAMPLE A\n"
	message += "    copy files containing \"alice\" and \"bob\" from ./a to ./b\n"
	message += "    $ file-cmd cp ./a ./b -f=alice,bob\n"
	message += "EXAMPLE B\n"
	message += "    copy files containing \"alice\" and \"bob\" from ./a to ./b\n"
	message += "    $ file-cmd cp ./a ./b alice bob\n"
	message += "EXAMPLE C\n"
	message += "    copy files containing \"alice\" or \"bob\" from ./a to ./b\n"
	message += "    $ file-cmd cp ./a ./b --or alice bob\n"
	message += "EXAMPLE D\n"
	message += "    copy files recursively containing \"2026-02-02\" from ./a to ./b\n"
	message += "    $ file-cmd cp ./a ./b -r 2026-02-02"
	fmt.Println(message)
}

func printExampleMove() {
	message := "EXAMPLE A\n"
	message += "    move files containing \"alice\" and \"bob\" from ./a to ./b\n"
	message += "    $ file-cmd mv ./a ./b -f=alice,bob\n"
	message += "EXAMPLE B\n"
	message += "    move files containing \"alice\" and \"bob\" from ./a to ./b\n"
	message += "    $ file-cmd mv ./a ./b alice bob\n"
	message += "EXAMPLE C\n"
	message += "    move files containing \"alice\" or \"bob\" from ./a to ./b\n"
	message += "    $ file-cmd mv ./a ./b --or alice bob\n"
	message += "EXAMPLE D\n"
	message += "    move files recursively containing \"2026-02-02\" from ./a to ./b\n"
	message += "    $ file-cmd mv ./a ./b -r 2026-02-02"
	fmt.Println(message)
}

func printExampleRemove() {
	message := "EXAMPLE A\n"
	message += "    delete files containing \"alice\" and \"bob\" from current directory\n"
	message += "    $ file-cmd rm . -f=alice,bob\n"
	message += "EXAMPLE B\n"
	message += "    delete files containing \"alice\" and \"bob\" from current directory\n"
	message += "    $ file-cmd rm . alice bob\n"
	message += "EXAMPLE C\n"
	message += "    delete files containing \"alice\" or \"bob\" from current directory\n"
	message += "    $ file-cmd rm . --or alice bob\n"
	message += "EXAMPLE D\n"
	message += "    delete files recursively containing \"2026-02-02\" from current directory\n"
	message += "    $ file-cmd rm . -r 2026-02-02"
	fmt.Println(message)
}

func printExampleClean() {
	message := "EXAMPLE\n"
	message += "    delete empty folders and empty regular files from current directory\n"
	message += "    $ file-cmd clean ."
	fmt.Println(message)
}

func printExampleText() {
	message := "EXAMPLE A\n"
	message += "    generate random text to file rnd.txt in current directory\n"
	message += "    $ file-cmd text rnd.txt -s=10k \n"
	message += "EXAMPLE B\n"
	message += "    generate random text to file with space as new line, i.e w/o new line\n"
	message += "    $ file-cmd text rnd.txt -e=\" \" -s=10k\n"
	message += "EXAMPLE C\n"
	message += "    generate random text of 100 bytes to standard output (console)\n"
	message += "    $ file-cmd text -s=100\n"
	message += "EXAMPLE D\n"
	message += "    generate random text of 100 bytes with lower case letters to console\n"
	message += "    $ file-cmd text -s=100 -f=al\n"
	fmt.Println(message)
}

func printCopyright() {
	message := "Copyright 2026, Vitali Baumtrok (vbsw@mailbox.org).\n"
	message += "Distributed under the Boost Software License, Version 1.0."
	fmt.Println(message)
}

func printWarning(command *tCommand, err error) {
	if command.verbose {
		fmt.Println("warning:", err.Error())
	}
}
