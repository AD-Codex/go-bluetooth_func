package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/godbus/dbus/v5"
	"github.com/muka/go-bluetooth/bluez/profile/adapter"
	"github.com/muka/go-bluetooth/bluez/profile/agent"
	"github.com/muka/go-bluetooth/bluez/profile/device"
)

// Convert MAC
func formatMAC(mac string) string {
	// Example: 68:5F:4A:5C:10:31 -> 68_5F_4A_5C_10_31
	return strings.ReplaceAll(mac, ":", "_")
}

// RegisterAgent creates and registers a pairing agent
func registerAgent(agentPath dbus.ObjectPath) error {
	agen, err := agent.NewAgent1("com.example.agent", agentPath)
	if err != nil {
		return fmt.Errorf("failed to create agent: %v", err)
	}

	manager, err := agent.NewAgentManager1()
	if err != nil {
		return fmt.Errorf("failed to get agent manager: %v", err)
	}

	if err := manager.RegisterAgent(agen.Path(), "KeyboardDisplay"); err != nil {
		return fmt.Errorf("failed to register agent: %v", err)
	}

	if err := manager.RequestDefaultAgent(agen.Path()); err != nil {
		return fmt.Errorf("failed to request default agent: %v", err)
	}

	fmt.Println("✅ Bluetooth agent registered")
	return nil
}

// Pair with given MAC address
func Pair_device(adapter *adapter.Adapter1, mac string, timeout time.Duration) error {
	// Create device path from MAC
	devPath := dbus.ObjectPath(string(adapter.Path()) + "/dev_" + formatMAC(mac))
	agentPath := dbus.ObjectPath("/com/example/agent" + formatMAC(mac))
	// Ensure agent is registered
	if err := registerAgent(agentPath); err != nil {
		return err
	}

	// Get the device object
	dev, err := device.NewDevice1(devPath)
	if err != nil {
		return fmt.Errorf("failed to get device object: %v", err)
	}

	// Check if already paired
	props, err := dev.GetProperties()
	if err != nil {
		return fmt.Errorf("failed to get properties: %v", err)
	}
	if props.Paired {
		fmt.Printf("Device %s is already paired ✅\n", mac)
		return nil
	}

	fmt.Printf("Pairing with device %s ...\n", mac)
	err = dev.Pair()
	if err != nil {
		return fmt.Errorf("pairing failed: %v", err)
	}

	// Wait until device is paired
	start := time.Now()
	for {
		props, _ = dev.GetProperties()
		if props.Paired {
			fmt.Printf("Device %s paired successfully ✅\n", mac)
			break
		}
		if time.Since(start) > timeout {
			return fmt.Errorf("pairing timeout after %s", timeout)
		}
		time.Sleep(500 * time.Millisecond)
	}

	return nil

}

func unregisterAgent(agentPath dbus.ObjectPath) error {
	manager, err := agent.NewAgentManager1()
	if err != nil {
		return fmt.Errorf("failed to get agent manager: %v", err)
	}

	if err := manager.UnregisterAgent(agentPath); err != nil {
		return fmt.Errorf("failed to unregister agent: %v", err)
	}

	fmt.Println("✅ Bluetooth agent unregistered")
	return nil
}

// Remove the given MAC from pair list
func Remove_device(adapter *adapter.Adapter1, mac string) error {
	// Create device path from MAC
	devPath := dbus.ObjectPath(string(adapter.Path()) + "/dev_" + formatMAC(mac))
	agentPath := dbus.ObjectPath("/com/example/agent" + formatMAC(mac))

	// Ensure agent is unregistered
	if err := unregisterAgent(agentPath); err != nil {
		return err
	}

	// Remove the device from BlueZ
	if err := adapter.RemoveDevice(devPath); err != nil {
		return fmt.Errorf("failed to forget device: %v", err)
	}

	fmt.Printf("Device %s forgotten successfully ✅\n", mac)
	return nil
}
