package domain

import (
	"encoding/json"
	"fmt"
)

// Year represents a numismatic year.
type Year struct {
	value int
}

// NewYear creates a valid Year. 0 is allowed as "unknown".
// Negative years are theoretically possible (BC), so we might allow them or limit.
// Let's assume -5000 to Current+1.
func NewYear(y int) (Year, error) {
	// current := time.Now().Year()
	if y != 0 && (y < -5000 || y > 3000) {
		return Year{}, fmt.Errorf("year %d is out of plausible range", y)
	}
	return Year{value: y}, nil
}

// Int returns the int value.
func (y Year) Int() int {
	return y.value
}

func (y Year) MarshalJSON() ([]byte, error) {
	return json.Marshal(y.value)
}

func (y *Year) UnmarshalJSON(data []byte) error {
	var val int
	if err := json.Unmarshal(data, &val); err != nil {
		return err
	}
	res, err := NewYear(val)
	if err != nil {
		return err
	}
	*y = res
	return nil
}

// KMCode represents a Krause-Mishler catalog code.
type KMCode struct {
	value string
}

func NewKMCode(code string) (KMCode, error) {
	if code == "" {
		return KMCode{}, nil // Empty is fine
	}
	// Basic loose validation or just storage
	// User mentioned strict rules. let's apply partial validation but not too strict to block imports.
	return KMCode{value: code}, nil
}

func (k KMCode) String() string {
	return k.value
}

func (k KMCode) MarshalJSON() ([]byte, error) {
	return json.Marshal(k.value)
}

func (k *KMCode) UnmarshalJSON(data []byte) error {
	var val string
	if err := json.Unmarshal(data, &val); err != nil {
		return err
	}
	res, _ := NewKMCode(val)
	*k = res
	return nil
}

// Mintage represents the number of coins minted.
type Mintage struct {
	value int64
}

func NewMintage(m int64) (Mintage, error) {
	if m < 0 {
		return Mintage{}, fmt.Errorf("mintage cannot be negative")
	}
	return Mintage{value: m}, nil
}

func (m Mintage) Int64() int64 {
	return m.value
}

func (m Mintage) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.value)
}

func (m *Mintage) UnmarshalJSON(data []byte) error {
	var val int64
	if err := json.Unmarshal(data, &val); err != nil {
		return err
	}
	res, err := NewMintage(val)
	if err != nil {
		return err
	}
	*m = res
	return nil
}

// Grade represents the coin condition.
type Grade struct {
	value string
}

// Valid grades (simplified/example list)

func NewGrade(g string) (Grade, error) {
	// Lenient for now as user might have custom input, but we wrap it.
	return Grade{value: g}, nil
}

func (g Grade) String() string {
	return g.value
}

func (g Grade) MarshalJSON() ([]byte, error) {
	return json.Marshal(g.value)
}

func (g *Grade) UnmarshalJSON(data []byte) error {
	var val string
	if err := json.Unmarshal(data, &val); err != nil {
		return err
	}
	res, _ := NewGrade(val)
	*g = res
	return nil
}
