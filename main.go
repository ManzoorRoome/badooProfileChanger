package main

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

var moveUp bool

func main() {
	fmt.Println("Which account do you want to log in?")
	fmt.Println("1) Manzoor Roome")
	fmt.Println("2) Joe SleepyBiden")
	fmt.Println("3) Mansur 2")
	fmt.Print("Enter choice (1-3): ")
	var accountToLogIn int
	fmt.Scanln(&accountToLogIn)

	fmt.Printf("\nStarting automation...\n\n")

	for {
		xml := dumpUI()
		page := detectPage(xml)

		if page == 0 {
			if !waitForPage(10) {
				fmt.Println("No matching page found. Exiting.")
				return
			}
			continue
		}

		executePageAction(page, accountToLogIn)

		if page == 8 {
			fmt.Println("\nAutomation complete!")
			break
		}
	}
}

func dumpUI() string {
	exec.Command("adb", "shell", "uiautomator", "dump").Run()
	time.Sleep(200 * time.Millisecond)

	cmd := exec.Command("adb", "shell", "cat", "/sdcard/window_dump.xml")
	out, _ := cmd.Output()
	xmlStr := string(out)

	if strings.Contains(xmlStr, `com.badoo.mobile:id/ad_container`) {
		fmt.Println("Ad bar detected at bottom")
		moveUp = true
	} else {
		moveUp = false
	}

	return xmlStr
}

func detectPage(xml string) int {
	if strings.Contains(xml, "nav_bar_button_profile") && !strings.Contains(xml, "ownProfileRootView") {
		return 1
	}
	if strings.Contains(xml, "myProfileSettings") {
		return 2
	}
	if strings.Contains(xml, "Account,") && strings.Contains(xml, "@gmail.com") && !strings.Contains(xml, "Ready to log out?") {
		return 3
	}
	if strings.Contains(xml, "Log out") && !strings.Contains(xml, "Ready to log out?") {
		return 4
	}
	if strings.Contains(xml, "Ready to log out?") {
		return 5
	}
	if strings.Contains(xml, "Continue with other methods") {
		return 6
	}
	if strings.Contains(xml, "Continue With Google") {
		return 7
	}
	if strings.Contains(xml, "Choose an account") && strings.Contains(xml, "Badoo") {
		return 7
	}
	if strings.Contains(xml, "Create a passkey") && !strings.Contains(xml, "Ready to log out?") {
		return 8
	}
	return 0
}

func executePageAction(page int, accountToLogIn int) {
	switch page {
	case 1:
		if moveUp {
			tap(972, 2057)
		} else {
			tap(972, 2157)
		}
		fmt.Println("✓ Profile tapped")
	case 2:
		tap(701, 273)
		fmt.Println("✓ Settings tapped")
	case 3:
		tap(540, 615)
		fmt.Println("✓ Logging out from roome.manzoor@gmail.com")
	case 4:
		tap(540, 1520)
		fmt.Println("✓ Log out tapped")
	case 5:
		tap(540, 2067)
		fmt.Println("✓ Logged out")
		time.Sleep(time.Second)
	case 6:
		tap(540, 1639)
		fmt.Println("✓ Logging into other account")
	case 7:
		xml := dumpUI()
		if strings.Contains(xml, "Continue With Google") {
			tap(588, 776)
			fmt.Println("✓ Continue with Google tapped")
		} else {
			selectAccount(accountToLogIn)
			time.Sleep(time.Second)
		}
	case 8:
		tap(540, 2067)
		fmt.Println("✓ Logged into new account!")
	}
}

func selectAccount(choice int) {
	switch choice {
	case 1:
		tap(578, 865)
		fmt.Println("✓ Selected: Manzoor Roome")
		time.Sleep(time.Second)
	case 2:
		tap(548, 1609)
		fmt.Println("✓ Selected: Joe SleepyBiden")
		time.Sleep(time.Second)
	case 3:
		tap(462, 1287)
		fmt.Println("✓ Selected: Mansur 2")
		time.Sleep(time.Second)
	}
}

func waitForPage(seconds int) bool {
	retries := seconds / 2
	for i := 0; i < retries; i++ {
		time.Sleep(2 * time.Second)
		xml := dumpUI()
		if detectPage(xml) != 0 {
			return true
		}
	}
	return false
}

func tap(x, y int) {
	exec.Command("adb", "shell", "input", "tap", fmt.Sprintf("%d", x), fmt.Sprintf("%d", y)).Run()
}
