package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/bluez/profile/adapter"
	"github.com/muka/go-bluetooth/bluez/profile/device"

	"bluetooth_v2/devices"
)

// Scane available all (BLE + Classic) devices
func Scan_allDevices(a *adapter.Adapter1, duration time.Duration) ([]devices.Discovered_data, error) {
	fmt.Println("Starting Bluetooth (Classic + BLE) scan...")

	// Start discovery
	// discoveryChan → channel of discovered devices (*adapter.DeviceDiscovered).
	// cancel → function to stop discovery.
	// err → any errors starting discovery.
	discoveryChan, cancel, err := api.Discover(a, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to start discovery: %v", err)
	}
	defer cancel()

	// stores devices by MAC address, avoids duplicates
	devicesMap := make(map[string]devices.Discovered_data)
	var mu sync.Mutex
	done := make(chan struct{})

	// Start goroutine to process discovered devices
	go detected_devices(discoveryChan, devicesMap, &mu, done)

	// Wait for the duration
	time.Sleep(duration)
	// stop discovery
	cancel()
	// wait for goroutine to finish processing
	<-done

	// Convert map to slice
	mu.Lock()
	devices := make([]devices.Discovered_data, 0, len(devicesMap))
	for _, d := range devicesMap {
		devices = append(devices, d)
	}
	mu.Unlock()

	fmt.Println("Scan complete ✅")
	return devices, nil
}

// Goroutine to process discovered devices
func detected_devices(discoveryChan <-chan *adapter.DeviceDiscovered, devicesMap map[string]devices.Discovered_data, mu *sync.Mutex, done chan struct{}) {
	// Continuously reads discovered devices until the channel is closed
	for ev := range discoveryChan {
		dev, err := device.NewDevice1(ev.Path)
		if err != nil {
			fmt.Printf("Failed to get device %s: %s", ev.Path, err)
			continue
		}
		props := dev.Properties

		devType := "Unknown"
		if len(props.UUIDs) > 0 {
			devType = "BLE"
		} else if props.Class != 0 {
			devType = "Classic"
		}

		device := devices.Discovered_data{
			Name:   props.Name,
			UUIDs:  props.UUIDs,
			MAC:    props.Address,
			Signal: props.RSSI,
			Class:  props.Class,
			Type:   devType,
		}

		// Store in map to avoid duplicates
		mu.Lock()
		if _, exists := devicesMap[device.MAC]; !exists {
			devicesMap[device.MAC] = device
			fmt.Printf("🚀 [%s] Name: %s, MAC: %s, Signal: %d, Class: %d\n", device.Type, device.Name, device.MAC, device.Signal, device.Class)

		}
		mu.Unlock()
	}
	// Sends a signal to the main routine
	done <- struct{}{}
}
