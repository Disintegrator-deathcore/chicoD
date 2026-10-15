package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/schollz/progressbar"
	"golang.org/x/sys/windows/registry"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print(
		"Выберите действие: \n",
		"1) Установить значок диска\n",
		"2) Удалить значок диска\n",
		"3) Скачать процессор бесплатно\n",
		"0) Выйти\n",
		": ",
	)
	scanner.Scan()
	ans := strings.TrimSpace(scanner.Text())

	switch ans {
	case "1":
		SetIconToDisk()
	case "2":
		DeleteIconFromDisk()
	case "3":
		InstallNewCPU()
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

	CheckPath(pathToIcon)

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

func InstallNewCPU() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Введите название процессора для установки: ")
	scanner.Scan()
	procName := strings.TrimSpace(scanner.Text())
	k, _, err := registry.CreateKey(
		registry.LOCAL_MACHINE,
		`HARDWARE\DESCRIPTION\System\CentralProcessor\0`,
		registry.ALL_ACCESS,
	)
	if err != nil {
		log.Println(err)
		k.Close()
	}

	bar := progressbar.New(100)
	for i := 0; i < 100; i++ {
		bar.Add(1)
		time.Sleep(40 * time.Millisecond)
	}
	fmt.Println()

	k.SetStringValue(`ProcessorNameString`, procName)
	k.Close()

	fmt.Println("Новый процессор установлен")
}

// func spinner(done <-chan struct{}) {
// 	frames := []rune{'|', '/', '-', '\\'}
// 	i := 0
// 	for {
// 		select {
// 		case <-done:
// 			fmt.Print("\r\033[K") // очистить строку
// 			return
// 		default:
// 			fmt.Printf("\r%c Обработка...", frames[i%len(frames)])
// 			i++
// 			time.Sleep(100 * time.Millisecond)
// 		}
// 	}
// }

// func longTask() {
// 	time.Sleep(3 * time.Second)
// }
