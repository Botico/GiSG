package main

import (
	"fmt"
	"bufio"
	"strings"
	"strconv"
	"os"
)

func main() {
	targetFile := ""

	if len(os.Args) != 2 {
		fmt.Println("GISG: Please supply your main file.")
		os.Exit(1)
	} else {
		targetFile = os.Args[1]
	}

	if _, err := os.Stat(targetFile); os.IsNotExist(err) {
		fmt.Println("GISG: Main file does not exist.")
                os.Exit(1)
	}

	fmt.Println("Please provide URLS for your dependencies (When you are done, type .)....")
	dependencies := GetList("Dependency")

	if len(dependencies) != 0 {
		fmt.Print("Dependencies: ")
		PrintList(dependencies)
	}

	fmt.Println("Please provide needed command line tools (When you are done, type .)....")
	commandLineTools := GetList("Command Line Tool")

	fmt.Print("Command line tools: ")

	if len(commandLineTools) != 0 {
		PrintList(commandLineTools)
	}

	var installToolsCommandsMac []string
	var installToolsCommandsLin []string
	var installToolsCommandsWin []string

	fmt.Println("Any custom commands needed to install tool (leave blank for default)?")
	for i := 0; i < 3; i++ {
		if i == 0 {
			fmt.Println("Mac OS X specifications: ")
		} else if i == 1 {
			fmt.Println("Linux specifications: ")
		} else if i == 2 {
			fmt.Println("Windows specifications: ")
		}

		for o := 0; o < len(commandLineTools); o++ {
			fmt.Print("\nCustom Command (" + commandLineTools[o] + ") --> ")
			reader := bufio.NewReader(os.Stdin)
			command, err := reader.ReadString('\n')
			Check(err)
			command = strings.Trim(command, "\n")

			if i == 0 {
				if command == "" {
					installToolsCommandsMac = append(installToolsCommandsMac,
					"brew install " + commandLineTools[o])
					continue
				}

				if strings.Contains(command, "brew") == false && strings.Contains(command, "port") == false {
					fmt.Println("Please use brew or port (macports)!" +
					" These are the only package managers supported for linux.")
					os.Exit(1)
				}

				installToolsCommandsMac = append(installToolsCommandsMac, command)
			} else if i == 1 {
				if command == "" {
					installToolsCommandsLin = append(installToolsCommandsLin,
					"apt install " + commandLineTools[o])
					continue
				}

				if strings.Contains(command, "apt") == false && strings.Contains(command, "dpkg") == false {
					fmt.Println("Please use apt or dpkg!" +
					" These are the only package managers supported for linux.")
					os.Exit(1)
				}

				installToolsCommandsLin = append(installToolsCommandsLin, command)
			} else if i == 2 {
				if command == "" {
					installToolsCommandsWin = append(installToolsCommandsWin,
					"choco install " + commandLineTools[o])
					continue
				}

				if strings.Contains(command, "choco") == false {
					fmt.Println("Please use choco (chocolate)!" +
					" This is the only package manager supported for windows.")
					os.Exit(1)
				}

				installToolsCommandsWin = append(installToolsCommandsWin, command)
			}
		}

		if len(installToolsCommandsMac) != 0 {
			fmt.Print("Mac OS X Command Line Tool Install Commands: ")
			PrintList(installToolsCommandsMac)
		}

		if len(installToolsCommandsLin) != 0 {
			fmt.Print("Linux Command Line Tool Install Commands: ")
			PrintList(installToolsCommandsLin)
		}

		if len(installToolsCommandsWin) != 0 {
			fmt.Print("Windows Command Line Tool Install Commands: ")
			PrintList(installToolsCommandsWin)
		}
	}

	fmt.Println("Creating files....")
	CreateBatchFile(targetFile, installToolsCommandsWin, dependencies, commandLineTools)
	CreateBashFileMac(targetFile, installToolsCommandsMac, dependencies, commandLineTools)
	CreateBashFileLin(targetFile, installToolsCommandsLin, dependencies, commandLineTools)

	fmt.Println("Please move all three files to the dir that houses your main go file.\nIf you are going to publish your program, it is recomended that you tell in the README about how to install.Keep in mind that when you run the script, you have to run it in the directory of your main go file.")
}

func GetList(item string) []string {
	iterations := 0
        var things []string
        reader := bufio.NewReader(os.Stdin)

	for {
                iterations++
                fmt.Print(item + " (" + strconv.Itoa(iterations) + ") --> ")
                thing, err := reader.ReadString('\n')
                Check(err)

                thing = strings.Trim(thing, "\n")

                if thing == "." {
                        break
                }

                things = append(things, thing)
        }

	return things
}

func PrintList(things []string) {
	for i := 0; i < len(things); i++ {
                if i == len(things) - 1 {
                        fmt.Print(things[i] + ".\n\n")
                } else {
                        fmt.Print(things[i] + ", ")
                }
        }
}

func Check(err error) {
	if err != nil {
		panic(err)
	}
}

func CreateBashFileMac(targetFile string, installToolsCommands []string, dependencies []string, commandLineTools []string) {
	_, err := os.Create(strings.Split(targetFile, ".")[0] + "_mac-os-x_install.sh")
	Check(err)

	f, err := os.OpenFile(strings.Split(targetFile, ".")[0] + "_mac-os-x_install.sh", os.O_APPEND|os.O_WRONLY, 0600)
	Check(err)

	_, err = f.WriteString("echo 'Installer made with GISC'" + "\n")
	Check(err)


	_, err = f.WriteString(InstallProgramBashCode("git",
	"xcode-select --install") + "\n")
	Check(err)

	usesBrew := false
	usesPorts := false

	for i := 0; i < len(installToolsCommands); i++ {
		if strings.Contains(installToolsCommands[i], "brew") {
			usesBrew = true
		} else if strings.Contains(installToolsCommands[i], "port") {
			usesPorts = true
		}
	}

	if usesBrew && usesPorts {
		_, err = f.WriteString("\n" + InstallProgramBashCode("brew",
                "/usr/bin/ruby -e \"$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)\"") + " || exit 1\n")
		_, err = f.WriteString("\n" + InstallProgramBashCode("port",
        "mkdir -p /opt/mports\n\tcd /opt/mports\n\tgit clone https://github.com/macports/macports-base.git\n\tgit checkout v2.5.4\n\t" +
        "cd /opt/mports/macports-base\n\t./configure --enable-readline\n\tmake\n\tmake install\n\tmake distclean\n\t") + " || exit 1\n")
	} else if usesBrew {
		_, err = f.WriteString("\n" + InstallProgramBashCode("brew",
                "/usr/bin/ruby -e \"$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)\"") + " || exit 1\n")
	} else if usesPorts {
		_, err = f.WriteString("\n" + InstallProgramBashCode("port",
	"mkdir -p /opt/mports\n\tcd /opt/mports\n\tgit clone https://github.com/macports/macports-base.git\n\tgit checkout v2.5.4\n\t" +
	"cd /opt/mports/macports-base\n\t./configure --enable-readline\n\tmake\n\tmake install\n\tmake distclean\n\t") + " || exit 1\n")
	}


	if usesBrew {
		_, err = f.WriteString(InstallProgramBashCode("go", "brew install go") + "\n")
		Check(err)
	} else if usesPorts {
		_, err = f.WriteString(InstallProgramBashCode("go", "port install go") + "\n")
                Check(err)
	}

	for i := 0; i < len(installToolsCommands); i++ {
		_, err := f.WriteString("\n" + InstallProgramBashCode(commandLineTools[i], installToolsCommands[i]) + " || exit 1\n")
		Check(err)
	}

	for i := 0; i < len(dependencies); i++ {
		_, err = f.WriteString("\n" + "go get " + dependencies[i] + " || exit 1\n")
		Check(err)
	}

	_, err = f.WriteString("\n" + "go build || exit 1\n")
        Check(err)
}

func CreateBashFileLin(targetFile string, installToolsCommands []string, dependencies []string, commandLineTools []string) {
	_, err := os.Create(strings.Split(targetFile, ".")[0] + "_linux_install.sh")
	Check(err)
	f, err := os.OpenFile(strings.Split(targetFile, ".")[0] + "_linux_install.sh", os.O_APPEND|os.O_WRONLY, 0600)
	Check(err)

	_, err = f.WriteString("echo 'Installer made with GISC'\n")
	Check(err)

	_, err = f.WriteString(InstallProgramBashCode("go",
	"apt install golang-go") + "\n")
	Check(err)

	_, err = f.WriteString(InstallProgramBashCode("git",
	"apt install git") + "\n")
        Check(err)

	for i := 0; i < len(installToolsCommands); i++ {
		_, err := f.WriteString("\n" + InstallProgramBashCode(commandLineTools[i], installToolsCommands[i]) + " || exit 1\n")
		Check(err)
	}

	for i := 0; i < len(dependencies); i++ {
		_, err = f.WriteString("\n" + "go get " + dependencies[i] + " || exit 1\n")
		Check(err)
	}

	_, err = f.WriteString("\n" + "go build || exit 1\n")
        Check(err)
}

func CreateBatchFile(targetFile string, installToolsCommands []string, dependencies []string, commandLineTools []string) {
	_, err := os.Create(strings.Split(targetFile, ".")[0] + "_install.bat")
	Check(err)

	f, err := os.OpenFile(strings.Split(targetFile, ".")[0] + "_install.bat", os.O_APPEND|os.O_WRONLY, 0600)
	Check(err)

	_, err = f.WriteString("ECHO 'Installer made with GISG'\n")
	Check(err)

	_, err = f.WriteString(ProgramExistsBatchCode("choco") + "\n")
	Check(err)

	_, err = f.WriteString(ProgramExistsBatchCode("go") + "\n")
	Check(err)

	_, err = f.WriteString(ProgramExistsBatchCode("git") + "\n")
        Check(err)

	for i := 0; i < len(installToolsCommands); i++ {
		_, err := f.WriteString("\n" + InstallProgramBatchCode(commandLineTools[i], installToolsCommands[i]) + " || exit /b\n")
		Check(err)
	}

	for i := 0; i < len(dependencies); i++ {
		_, err = f.WriteString("\n" + "go get " + dependencies[i] + " || exit /b\n")
		Check(err)
	}

	_, err = f.WriteString("\n" + "go build || exit 1\n")
        Check(err)
}

func ProgramExistsBashCode(command string) string {
	return "if ! [ -x \"$(command -v " + command + ")\" ]; then\n\techo 'Error: " + command + " is not installed. Please install it.' >&2\n\texit 1\nfi"
}

func InstallProgramBashCode(command string, installCode string) string {
	return "if ! [ -x \"$(command -v " + command + ")\" ]; then\n\t " + installCode + ">&2\n\texit 1\nfi"
}

func ProgramExistsBatchCode(command string) string {
	return "\nWHERE " + command + "\nIF %ERRORLEVEL% NEQ 0 ECHO " + command + " is not installed. Please install it."
}

func InstallProgramBatchCode(command string, installCode string) string {
	return "\nWHERE " + command + "\nIF %ERRORLEVEL% NEQ 0 " + installCode
}
