/*
  pin
    state management of pins
*/

package pin

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/gorilla/mux"
	"github.com/stianeikeland/go-rpio"
	"gitlab.com/bobymcbobs/go-rpi-gpio-api/src/types"
)

// GetPin ...
// returns the state of a given pin
func GetPin(num int) (pin types.Pin, err error) {
	err = ValidatePinNumber(num)
	if err != nil {
		return pin, err
	}
	err = OpenGPIOpins()
	if err != nil {
		return pin, err
	}
	defer rpio.Close()
	pinSelect := rpio.Pin(num)
	state := pinSelect.Read()

	pin = types.Pin{
		Number: num,
		State:  int(state),
	}
	return pin, err
}

// WritePin ...
// writes the state to a given pin
func WritePin(num int, state int, mode int) (pin types.Pin, err error) {
	err = ValidatePinNumber(num)
	if err != nil {
		return pin, err
	}
	err = OpenGPIOpins()
	if err != nil {
		return pin, err
	}
	defer rpio.Close()
	pinSelect := rpio.Pin(num)
	switch mode {
	case 0:
		pinSelect.Mode(rpio.Input)
	case 1:
		pinSelect.Mode(rpio.Output)
	default:
		return types.Pin{}, fmt.Errorf("Invalid mode - valid options: 0 (input), 1 (output)")
	}
	switch state {
	case 0:
		pinSelect.Write(rpio.Low)
		// pinSelect.Low()
	case 1:
		pinSelect.Write(rpio.High)
		// pinSelect.High()
	default:
		return types.Pin{}, fmt.Errorf("Invalid pin number - valid options: 0 (low), 1 (high)")
	}

	pin = types.Pin{
		Number: num,
		State:  state,
	}
	return pin, err
}

// ListPins ...
// returns the state of all pins
func ListPins() (pinList types.PinList, err error) {
	for num := 1; num <= 40; num++ {
		pinState, err := GetPin(num)
		if err != nil {
			return types.PinList{}, err
		}
		pinList = append(pinList, pinState)
	}
	return pinList, err
}

// OpenGPIOpins ...
// makes sure that the GPIO pins can be access
func OpenGPIOpins() (err error) {
	if err := rpio.Open(); err != nil {
		return err
	}
	return err
}

// ValidatePinNumber ...
// ensures that only pin numbers between a range can be accessed
func ValidatePinNumber(num int) (err error) {
	if !(num <= 40 && num >= 0) {
		return fmt.Errorf("Invalid pin number")
	}
	return err
}

// CheckGPIOpinForGround ...
// validates a pin number, checking it for no being ground
func CheckGPIOpinForGround(r *http.Request) (matches bool) {
	vars := mux.Vars(r)
	pinIdstr := vars["pin"]
	matches, _ = regexp.MatchString(`\b(([1-5])|(([7-8]))|(([0-1]|1[0-3]))|(([0-1]|1[5-9]))|(([2]|2[1-4]))|(([2]|2[6-9]))|(([2]|2[7-8]))|(([3]|3[1-3]))|(([3]|3[5-8]))|40)\b`, pinIdstr)
	return !matches
}

// CheckForValidState ...
// validates pin state
func CheckForValidState(r *http.Request) (matches bool) {
	vars := mux.Vars(r)
	pinIdstr := vars["state"]
	matches, _ = regexp.MatchString(`\b(0|1)\b`, pinIdstr)
	return matches
}
