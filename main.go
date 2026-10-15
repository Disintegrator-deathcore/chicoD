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

func checkInputDriver(scanner *bufio.Scanner) string {
	if strings.TrimSpace(scanner.Text()) == "" {
		log.Println("Значение не было введено")
		return `{4d36e968-e325-11ce-bfc1-08002be10318}\0001`
	} else {
		log.Println("Значение было введено")
		return strings.TrimSpace(scanner.Text())
	}
}

func checkInputFriendlyName(scanner *bufio.Scanner) string {
	if strings.TrimSpace(scanner.Text()) == "" {
		log.Println("Значение не было введено")
		return `Ебанутая видюха`
	} else {
		log.Println("Значение было введено")
		return strings.TrimSpace(scanner.Text())
	}
}

func InstallNewGPU() {
	searchBase := `SYSTEM\CurrentControlSet\Enum\PCI`
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, searchBase, registry.ENUMERATE_SUB_KEYS|registry.QUERY_VALUE)
	if err != nil {
		log.Println(err)
	}
	defer k.Close()

	found := false

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Введите ключ драйвера: ")
	scanner.Scan()

	driverCode := checkInputDriver(scanner)
	strings.TrimSpace(scanner.Text())

	fmt.Print("Введите название видеокарты, которую хотите скачать: ")
	scanner.Scan()

	friendlyName := checkInputFriendlyName(scanner)

	searchKeyRecursive(k, searchBase, &found, friendlyName, driverCode)

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
		if keyContainsDriverCode(subKey, driverCode) == true {
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
