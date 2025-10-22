package main

import (
	"fmt"

	"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/bluez/profile/adapter"
)

// Returns the default Bluetooth adapter (usually hci0)
func Get_adapter() (*adapter.Adapter1, error) {
	a, err := api.GetDefaultAdapter()
	if err != nil {
		return nil, fmt.Errorf("failed to get default adapter: %v", err)
	}
	return a, nil
}

// Return adapter properties
func Prop_adapter(a *adapter.Adapter1) (name string, power string, discoverable string, err error) {
	props, err := a.GetProperties()
	if err != nil {
		return "", "", "", fmt.Errorf("failed to get adapter properties: %v", err)
	}

	name = string(a.Path())
	if props.Powered {
		power = "ON  🟢"
	} else {
		power = "OFF 🔴"
	}

	if props.Discoverable {
		discoverable = "YES 🌐"
	} else {
		discoverable = "NO  🚫"
	}

	return name, power, discoverable, nil
}

// Enable or Disable the adpater
func PowerOn_adapter(a *adapter.Adapter1, on bool) error {
	props, err := a.GetProperties()
	if err != nil {
		return fmt.Errorf("failed to get adapter properties: %v", err)
	}
	fmt.Println("Adapter:", props.Name, props.Address)
	fmt.Println("Powered:", props.Powered)

	if on {
		if !props.Powered {
			fmt.Println("Bluetooth adapter is OFF — powering on...")
			// sudo rfkill unblock bluetooth
			if err := a.SetProperty("Powered", true); err != nil {
				return fmt.Errorf("failed to power on adapter: %v (sudo rfkill unblock bluetooth)", err)
			}
			fmt.Println("Adapter powered ON ✅")
		} else {
			fmt.Println("Adapter already ON ✅")
		}
	} else {
		if props.Powered {
			fmt.Println("Bluetooth adapter is ON — powering off...")
			if err := a.SetProperty("Powered", false); err != nil {
				return fmt.Errorf("failed to power off adapter: %v", err)
			}
			fmt.Println("Adapter powered OFF ✅")
		} else {
			fmt.Println("Adapter already OFF ✅")
		}
	}

	return nil
}

// Makes the adapter visible to other devices
func Make_discoverable(a *adapter.Adapter1, on bool, timeout uint32) error {
	if on {
		fmt.Println("Making adapter discoverable...")

		if err := a.SetProperty("DiscoverableTimeout", timeout); err != nil {
			return fmt.Errorf("failed to set discoverable timeout: %v", err)
		}

		if err := a.SetProperty("Discoverable", true); err != nil {
			return fmt.Errorf("failed to make adapter discoverable: %v", err)
		}

		fmt.Printf("Adapter is now discoverable for %d seconds 🔍\n", timeout)
	} else {
		fmt.Println("Stopping discoverable mode...")

		if err := a.SetProperty("Discoverable", false); err != nil {
			return fmt.Errorf("failed to stop discoverable mode: %v", err)
		}

		fmt.Println("Adapter is no longer discoverable ❌")
	}

	return nil
}
