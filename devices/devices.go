package devices

// Information about a Bluetooth device
type Discovered_data struct {
	Name   string   // Device name
	UUIDs  []string // BLE service UUIDs (empty for Classic devices)
	MAC    string   // Device address
	Signal int16    // RSSI (signal strength, BLE only; 0 for Classic)
	Class  uint32   // Device class (for Classic devices)
	Type   string   // "BLE" or "Classic"
}
