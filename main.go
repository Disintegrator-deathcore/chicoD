package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/sys/windows/registry"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Выберите действие: \n1) Установить значок диска\n2) Удалить значок диска\n0) Выйти\n: ")
	scanner.Scan()
	ans := strings.TrimSpace(scanner.Text())

	switch ans {
	case "1":
		SetIconToDisk()
	case "2":
		DeleteIconFromDisk()
	case "0":
		os.Exit(0)
	default:
		os.Exit(0)
	}

	fmt.Println("Нажмите Enter для выхода...")
	scanner.Scan()
}

func CheckPath(pathToIcon string) {
	if len(pathToIcon) < 2 {
		log.Println("Переменная не была задана...")

		SetIconToDisk()
	} else {
		log.Println("Переменная была задана")
	}
}

func SetIconToDisk() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Введите букву диска: ")
	scanner.Scan()
	disk := strings.TrimSpace(scanner.Text())

	fmt.Print("Введите путь для значка: ")
	scanner.Scan()
	pathToIcon := strings.TrimSpace(scanner.Text())

	pathToIcon = CheckPath(pathToIcon)

	k, _, err := registry.CreateKey(
		registry.LOCAL_MACHINE,
		"SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Explorer\\DriveIcons\\"+disk+"\\DefaultIcon",
		registry.RESOURCE_REQUIREMENTS_LIST,
	)
	if err != nil {
		log.Fatal(err)
	}

	k.SetStringValue(``, pathToIcon)
	k.Close()
}

func DeleteIconFromDisk() {
	var disk string

	fmt.Print("Введите букву диска: ")
	fmt.Scan(&disk)

	k := registry.DeleteKey(registry.LOCAL_MACHINE, "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Explorer\\DriveIcons\\"+disk+"\\DefaultIcon")
	if k != nil {
		log.Print(k.Error())
	}
	k = registry.DeleteKey(registry.LOCAL_MACHINE, "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Explorer\\DriveIcons\\"+disk)
}
