package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"os/exec"
)

// executeCommand takes the command and its args and run it or return an error
func executeCommand(command string, args []string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

// Change MAC address using `ip` instead of `ifconfig`
func changeMacAddressIpCommand(interfaceName string) error {
	// generate a random mac address
	newMacAddress := generateRandomMac()

	// Bring the interface down
	if err := executeCommand("sudo ip link set",
		[]string{interfaceName, "down"}); err != nil {
		return fmt.Errorf("failed to bring interface down: %v", err)
	}

	// Change the MAC address
	if err := executeCommand("sudo ip link set",
		[]string{interfaceName, "address", newMacAddress}); err != nil {
		return fmt.Errorf("failed to change MAC address: %v", err)
	}

	// Bring the interface up
	if err := executeCommand("sudo ip link set",
		[]string{interfaceName, "up"}); err != nil {
		return fmt.Errorf("failed to bring interface up: %v", err)
	}

	return nil
}

// generateRandomMac generates a random MAC address.
func generateRandomMac() string {
	// Generate 6 random bytes for the MAC address
	mac := make([]byte, 6)

	// First byte should have the least significant bit clear (uni-cast address)
	mac[0] = byte(randInt(0x00, 0xFE)) // Random value from 0x00 to 0xFE

	// Next 5 bytes can be completely random
	for i := 1; i < 6; i++ {
		mac[i] = byte(randInt(0, 0xFF)) // Random value from 0x00 to 0xFF
	}

	// Format the MAC address as a string "XX:XX:XX:XX:XX:XX"
	return fmt.Sprintf("%02X:%02X:%02X:%02X:%02X:%02X", mac[0], mac[1], mac[2], mac[3], mac[4], mac[5])
}

// randInt generates a random integer in the range [min, max).
func randInt(min, max int) int {
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(max-min)))
	return int(n.Int64()) + min
}

// changeMacAddress changes the MAC address of the given interface.
func changeMacAddressIfConfig(interfaceName string) error {
	// generate a random mac address
	newMacAddress := generateRandomMac()

	// Bring the interface down
	if err := executeCommand("sudo ifconfig",
		[]string{interfaceName, "down"}); err != nil {
		return fmt.Errorf("failed to bring interface down: %v", err)
	}

	// Change the MAC address
	if err := executeCommand("sudo ifconfig",
		[]string{interfaceName, "hw", "ether", newMacAddress}); err != nil {
		return fmt.Errorf("failed to change MAC address: %v", err)
	}

	// Bring the interface up

	if err := executeCommand("sudo ifconfig",
		[]string{interfaceName, "up"}); err != nil {
		return fmt.Errorf("failed to bring interface up: %v", err)
	}

	return nil
}

func main() {

}
