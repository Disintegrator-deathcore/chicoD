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

	// Меню выбора действия
	fmt.Print(
		"Выберите действие: \n",
		"1) Установить значок диска\n",
		"2) Удалить значок диска\n",
		"3) Скачать процессор бесплатно\n",
		"4) Скачать видеокарту бесплатно (пробовать на свой страх и риск)\n",
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
	case "4":
		InstallNewGPU()
	case "0":
		os.Exit(0)
	default:
		os.Exit(0)
	}

	fmt.Println("Нажмите Enter для выхода...")
	scanner.Scan()
}

// Функция проверки пути
func CheckPath(pathToIcon string) bool {
	if len(pathToIcon) < 2 {
		log.Println("Переменная не была задана...")
		return false
	} else {
		return true
	}
}

// Функция установки иконки диска
func SetIconToDisk() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Введите букву диска: ")
	scanner.Scan()
	disk := strings.TrimSpace(scanner.Text())

	fmt.Print("Введите полный путь для значка: ")
	scanner.Scan()
	pathToIcon := strings.TrimSpace(scanner.Text())

	if CheckPath(pathToIcon) {
		k, _, err := registry.CreateKey(
			registry.LOCAL_MACHINE,
			"SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Explorer\\DriveIcons\\"+disk+"\\DefaultIcon",
			registry.QUERY_VALUE|registry.SET_VALUE,
		)
		if err != nil {
			log.Fatal(err)
		}
		err = k.SetStringValue(``, pathToIcon)
		if err != nil {
			log.Println(err)
		}
		k.Close()
	}
}

// Функция удаления иконки с диска
func DeleteIconFromDisk() {
	var disk string

	fmt.Print("Введите букву диска: ")
	fmt.Scan(&disk)

	err := registry.DeleteKey(registry.LOCAL_MACHINE, "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Explorer\\DriveIcons\\"+disk+"\\DefaultIcon")
	if err != nil {
		log.Println(err)
	}
	err = registry.DeleteKey(registry.LOCAL_MACHINE, "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Explorer\\DriveIcons\\"+disk)
	if err != nil {
		log.Println(err)
	}
}

// Функция установки нового названия процессора
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
		log.Fatal(err)
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

// Функция проверки ключа драйвера
func checkInputDriver(scanner *bufio.Scanner) string {
	if strings.TrimSpace(scanner.Text()) == "" {
		return `{4d36e968-e325-11ce-bfc1-08002be10318}\0001`
	} else {
		return strings.TrimSpace(scanner.Text())
	}
}

// Функция проверки FriendlyName
func checkInputFriendlyName(scanner *bufio.Scanner) string {
	if strings.TrimSpace(scanner.Text()) == "" {
		return `Ебанутая видюха`
	} else {
		return strings.TrimSpace(scanner.Text())
	}
}

// Функция установки названия видеокарты
func InstallNewGPU() {
	pciEnumPath := `SYSTEM\CurrentControlSet\Enum\PCI`
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, pciEnumPath, registry.ENUMERATE_SUB_KEYS|registry.QUERY_VALUE)
	if err != nil {
		log.Fatal(err)
	}
	defer k.Close()

	found := false

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Введите ключ драйвера: ")
	scanner.Scan()

	driverCode := checkInputDriver(scanner)

	fmt.Print("Введите название видеокарты, которую хотите скачать: ")
	scanner.Scan()

	friendlyName := checkInputFriendlyName(scanner)

	searchKeyRecursive(k, pciEnumPath, &found, friendlyName, driverCode)

	if !found {
		fmt.Println("Ключ, содержащий указанный код драйвера, не найден.")
	}
}

// searchKeyRecursive рекурсивно обходит подключи.
func searchKeyRecursive(key registry.Key, path string, found *bool, friendlyName string, driverCode string) {
	// Получаем список подключей текущего ключа.
	subKeys, err := key.ReadSubKeyNames(-1)
	if err != nil {
		// Некоторые ключи могут быть недоступны для чтения — просто пропускаем их.
		return
	}

	for _, subKeyName := range subKeys {
		// Формируем полный путь к подключу.
		fullPath := path + "\\" + subKeyName

		// Открываем подключ. Если не удаётся (например, нет прав), пропускаем.
		subKey, err := registry.OpenKey(registry.LOCAL_MACHINE, fullPath, registry.QUERY_VALUE|registry.SET_VALUE|registry.ENUMERATE_SUB_KEYS)
		if err != nil {
			continue
		}

		// Проверяем, содержит ли текущий подключ (или его значения) искомый код драйвера.
		if keyContainsDriverCode(subKey, driverCode) {
			fmt.Printf("🔍 Найден ключ: %s\n", fullPath)

			// Добавляем или обновляем значение FriendlyName.
			err = subKey.SetStringValue("FriendlyName", friendlyName)
			if err != nil {
			} else {
				*found = true
			}
		}

		// Продолжаем рекурсивный обход.
		searchKeyRecursive(subKey, fullPath, found, friendlyName, driverCode)
		subKey.Close()
	}
}

// keyContainsDriverCode проверяет, встречается ли driverCode в имени ключа или в его строковых значениях.
func keyContainsDriverCode(key registry.Key, code string) bool {
	// 1. Проверяем имя самого ключа (в реестре имена подключей часто содержат код).
	//    Но у нас нет прямого доступа к имени ключа из registry.Key, поэтому
	//    эту проверку делает вызывающий код (сравнение subKeyName), см. ниже.

	// 2. Перебираем все строковые значения в ключе.
	valueNames, err := key.ReadValueNames(-1)
	if err != nil {
		return false
	}
	for _, name := range valueNames {
		val, _, err := key.GetStringValue(name)
		if err != nil {
			continue
		}
		if strings.Contains(val, code) {
			return true
		}
	}
	return false
}
