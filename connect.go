package main

import (
	"fmt"
	"time"

	"github.com/godbus/dbus/v5"
	"github.com/muka/go-bluetooth/bluez/profile/adapter"
	"github.com/muka/go-bluetooth/bluez/profile/device"
)

// Connect to given MAC address
func Connect_device(adapter *adapter.Adapter1, mac string) error {
	fmt.Printf("🔗 Connecting to device %s...\n", mac)

	// Create device path from MAC
	devPath := dbus.ObjectPath(string(adapter.Path()) + "/dev_" + formatMAC(mac))

	// Get the device object
	dev, err := device.NewDevice1(devPath)
	if err != nil {
		return fmt.Errorf("failed to get device object: %v", err)
	}

	// Connect to the device
	err = dev.Connect()
	if err != nil {
		return fmt.Errorf("connection failed: %v", err)
	}

	// Wait for the connection to establish
	time.Sleep(2 * time.Second)

	props := dev.Properties
	fmt.Printf("✅ Connected to %s (%s)\n", props.Name, props.Address)

	return nil
}

// Disconnect from given MAC address
func Disconnect_device(adapter *adapter.Adapter1, mac string) error {
	// Create device path from MAC
	devPath := dbus.ObjectPath(string(adapter.Path()) + "/dev_" + formatMAC(mac))

	// Get the device object
	dev, err := device.NewDevice1(devPath)
	if err != nil {
		return fmt.Errorf("failed to get device object: %v", err)
	}

	// Check if connected
	props, _ := dev.GetProperties()
	if !props.Connected {
		fmt.Printf("Device %s is already disconnected ✅\n", mac)
		return nil
	}

	// Disconnect
	if err := dev.Disconnect(); err != nil {
		return fmt.Errorf("failed to disconnect: %v", err)
	}

	fmt.Printf("Device %s disconnected successfully ✅\n", mac)

	return nil
}
