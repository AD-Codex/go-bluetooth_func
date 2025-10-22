package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.PanicLevel)

	// Initialize adapter
	a, err := Get_adapter()
	if err != nil {
		log.Fatalf("Error getting adapter: %v", err)
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		adapterName, powerState, discoverableState, err := Prop_adapter(a)
		if err != nil {
			fmt.Printf("⚠️  Failed to get adapter properties: %v\n", err)
			continue
		}
		fmt.Println("\n====================================================")
		fmt.Println("          Bluetooth Control Terminal ")
		fmt.Println("====================================================")
		fmt.Printf("Adapter: %s | Power: %s | Discoverable: %s\n", adapterName, powerState, discoverableState)
		fmt.Println("-----------------------------------")
		fmt.Println("1. Power ON Adapter		| 2. Power OFF Adapter")
		fmt.Println("3. Make Discoverable		| 4. Stop Discoverable")
		fmt.Println("5. Scan Devices			| 6. Connect Devic")
		fmt.Println("7. Disconnect Device		| 8. Pair Device")
		fmt.Println("9. Remove Device		| 10. Exit")
		fmt.Print("Enter choice: ")

		choiceStr, _ := reader.ReadString('\n')
		choiceStr = strings.TrimSpace(choiceStr)
		choice, err := strconv.Atoi(choiceStr)
		if err != nil {
			fmt.Println("❌ Invalid input. Enter a number.")
			continue
		}

		switch choice {
		case 1:
			err = PowerOn_adapter(a, true)
			if err != nil {
				fmt.Printf("⚠️  Error powering ON: %v\n", err)
			} else {
				fmt.Println("✅ Adapter powered ON")
			}

		case 2:
			err = PowerOn_adapter(a, false)
			if err != nil {
				fmt.Printf("⚠️  Error powering OFF: %v\n", err)
			} else {
				fmt.Println("✅ Adapter powered OFF")
			}

		case 3:
			err = Make_discoverable(a, true, 30)
			if err != nil {
				fmt.Printf("⚠️  Error enabling discoverable: %v\n", err)
			} else {
				fmt.Println("✅ Adapter is discoverable for 30 seconds")
			}

		case 4:
			err = Make_discoverable(a, false, 0)
			if err != nil {
				fmt.Printf("⚠️  Error disabling discoverable: %v\n", err)
			} else {
				fmt.Println("✅ Discoverable mode stopped")
			}

		case 5:
			fmt.Print("Enter scan duration (seconds): ")
			durStr, _ := reader.ReadString('\n')
			durStr = strings.TrimSpace(durStr)
			dur, _ := strconv.Atoi(durStr)

			fmt.Println("🔍 Scanning for devices...")
			devices, err := Scan_allDevices(a, time.Duration(dur)*time.Second)
			if err != nil {
				fmt.Printf("⚠️  Scan error: %v\n", err)
			} else {
				fmt.Printf("✅ %d devices found\n", len(devices))
				for addr, dev := range devices {
					if len(dev.Name) > 0 {
						fmt.Printf("   %d | Name: %s | MAC: %s\n", addr, dev.Name, dev.MAC)
					}
				}
			}

		case 6:
			fmt.Print("Enter MAC to connect: ")
			mac, _ := reader.ReadString('\n')
			mac = strings.TrimSpace(mac)
			if err := Connect_device(a, mac); err != nil {
				fmt.Printf("⚠️  Connect error: %v\n", err)
			} else {
				fmt.Println("✅ Connected successfully")
			}

		case 7:
			fmt.Print("Enter MAC to disconnect: ")
			mac, _ := reader.ReadString('\n')
			mac = strings.TrimSpace(mac)
			if err := Disconnect_device(a, mac); err != nil {
				fmt.Printf("⚠️  Disconnect error: %v\n", err)
			} else {
				fmt.Println("✅ Disconnected successfully")
			}

		case 8:
			fmt.Print("Enter MAC to pair: ")
			mac, _ := reader.ReadString('\n')
			mac = strings.TrimSpace(mac)
			if err := Pair_device(a, mac, 20); err != nil {
				fmt.Printf("⚠️  Pair error: %v\n", err)
			} else {
				fmt.Println("✅ Paired successfully")
			}

		case 9:
			fmt.Print("Enter MAC to remove: ")
			mac, _ := reader.ReadString('\n')
			mac = strings.TrimSpace(mac)
			if err := Pair_device(a, mac, 20); err != nil {
				fmt.Printf("⚠️  Pair error: %v\n", err)
			} else {
				fmt.Println("✅ Paired successfully")
			}
			if err := Remove_device(a, mac); err != nil {
				fmt.Printf("⚠️  Remove error: %v\n", err)
			} else {
				fmt.Println("✅ Removed successfully")
			}

		case 10:
			fmt.Println("👋 Exiting program.")
			return

		default:
			fmt.Println("❌ Invalid choice. Try again.")
		}

		fmt.Print("\nPress Enter to continue...")
		reader.ReadString('\n')
	}
}
