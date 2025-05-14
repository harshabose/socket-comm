package main

import (
	"encoding/json"
	"fmt"
	"time"
)

// TestStruct contains different time and duration fields
type TestStruct struct {
	DefaultTime     time.Time     `json:"default_time"`
	CustomTime      time.Time     `json:"custom_time,omitempty"`
	DefaultDuration time.Duration `json:"default_duration"`
	StringDuration  time.Duration `json:"string_duration"`
}

// CustomTime wraps time.Time with custom JSON marshaling
type CustomTime struct {
	time.Time
}

func (ct CustomTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(ct.Time.Unix())
}

func (ct *CustomTime) UnmarshalJSON(data []byte) error {
	var timestamp int64
	if err := json.Unmarshal(data, &timestamp); err != nil {
		return err
	}
	ct.Time = time.Unix(timestamp, 0)
	return nil
}

// DurationString wraps time.Duration with string representation
type DurationString struct {
	time.Duration
}

func (ds DurationString) MarshalJSON() ([]byte, error) {
	return json.Marshal(ds.String())
}

func (ds *DurationString) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	d, err := time.ParseDuration(str)
	if err != nil {
		return err
	}
	ds.Duration = d
	return nil
}

func testDefaultTimeFormat() {
	fmt.Println("=== Testing Default time.Time and time.Duration ===")

	// Test various time values
	now := time.Now()
	epoch := time.Unix(0, 0)
	microTime := time.Unix(0, 1234567890) // Unix timestamp with microseconds

	tests := []struct {
		name string
		data TestStruct
	}{
		{
			name: "Current Time",
			data: TestStruct{
				DefaultTime:     now,
				DefaultDuration: 5*time.Second + 250*time.Millisecond,
			},
		},
		{
			name: "Epoch",
			data: TestStruct{
				DefaultTime:     epoch,
				DefaultDuration: time.Hour + 30*time.Minute,
			},
		},
		{
			name: "Microseconds",
			data: TestStruct{
				DefaultTime:     microTime,
				DefaultDuration: 123*time.Microsecond + 456*time.Nanosecond,
			},
		},
	}

	for _, test := range tests {
		fmt.Printf("\n--- %s ---\n", test.name)

		// Marshal to JSON
		jsonData, err := json.Marshal(test.data)
		if err != nil {
			fmt.Printf("Marshal error: %v\n", err)
			continue
		}
		fmt.Printf("JSON: %s\n", string(jsonData))

		// Unmarshal back
		var result TestStruct
		if err := json.Unmarshal(jsonData, &result); err != nil {
			fmt.Printf("Unmarshal error: %v\n", err)
			continue
		}

		// Verify
		fmt.Printf("Original Time: %v\n", test.data.DefaultTime)
		fmt.Printf("Decoded Time:  %v\n", result.DefaultTime)
		fmt.Printf("Time Equal: %v\n", test.data.DefaultTime.Equal(result.DefaultTime))
		fmt.Printf("Original Duration: %v\n", test.data.DefaultDuration)
		fmt.Printf("Decoded Duration:  %v\n", result.DefaultDuration)
		fmt.Printf("Duration Equal: %v\n", test.data.DefaultDuration == result.DefaultDuration)
	}
}

func testCustomFormats() {
	fmt.Println("\n\n=== Testing Custom Time and Duration Formats ===")

	type CustomStruct struct {
		UnixTime    CustomTime     `json:"unix_time"`
		StringDur   DurationString `json:"string_duration"`
		NanoSeconds int64          `json:"nanoseconds"`
	}

	now := time.Now()

	data := CustomStruct{
		UnixTime:    CustomTime{now},
		StringDur:   DurationString{45*time.Minute + 30*time.Second},
		NanoSeconds: now.UnixNano(),
	}

	// Marshal
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("Marshal error: %v\n", err)
		return
	}
	fmt.Printf("Custom JSON: %s\n", string(jsonData))

	// Unmarshal
	var result CustomStruct
	if err := json.Unmarshal(jsonData, &result); err != nil {
		fmt.Printf("Unmarshal error: %v\n", err)
		return
	}

	fmt.Printf("Original Unix Time: %d\n", data.UnixTime.Unix())
	fmt.Printf("Decoded Unix Time:  %d\n", result.UnixTime.Unix())
	fmt.Printf("Original Duration: %s\n", data.StringDur.String())
	fmt.Printf("Decoded Duration:  %s\n", result.StringDur.String())
}

func testDifferentDurationFormats() {
	fmt.Println("\n\n=== Testing Different Duration Formats ===")

	durationTests := []time.Duration{
		// Nanoseconds
		1 * time.Nanosecond,
		123 * time.Nanosecond,
		// Microseconds
		1 * time.Microsecond,
		456 * time.Microsecond,
		// Milliseconds
		1 * time.Millisecond,
		789 * time.Millisecond,
		// Seconds
		1 * time.Second,
		30 * time.Second,
		// Minutes
		1 * time.Minute,
		5 * time.Minute,
		// Hours
		1 * time.Hour,
		24 * time.Hour,
		// Complex combinations
		1*time.Hour + 23*time.Minute + 45*time.Second + 678*time.Millisecond,
		365 * 24 * time.Hour, // 1 year
	}

	for _, dur := range durationTests {
		fmt.Printf("\nDuration: %v\n", dur)
		fmt.Printf("Nanoseconds: %d\n", dur.Nanoseconds())
		fmt.Printf("Microseconds: %d\n", dur.Microseconds())
		fmt.Printf("Milliseconds: %d\n", dur.Milliseconds())
		fmt.Printf("Seconds: %f\n", dur.Seconds())
		fmt.Printf("Minutes: %f\n", dur.Minutes())
		fmt.Printf("Hours: %f\n", dur.Hours())

		// Test JSON marshaling
		type DurWrapper struct {
			Duration time.Duration `json:"duration"`
		}

		wrapper := DurWrapper{Duration: dur}
		jsonData, err := json.Marshal(wrapper)
		if err != nil {
			fmt.Printf("Marshal error: %v\n", err)
			continue
		}
		fmt.Printf("JSON: %s\n", string(jsonData))

		// Test unmarshaling
		var result DurWrapper
		if err := json.Unmarshal(jsonData, &result); err != nil {
			fmt.Printf("Unmarshal error: %v\n", err)
			continue
		}
		fmt.Printf("Roundtrip successful: %v\n", dur == result.Duration)
	}
}

func testTimeZones() {
	fmt.Println("\n\n=== Testing Time Zones ===")

	locations := []*time.Location{
		time.UTC,
		time.Local,
		time.FixedZone("EST", -5*60*60),
		time.FixedZone("PST", -8*60*60),
	}

	baseTime := time.Date(2024, 3, 15, 14, 30, 45, 123456789, time.UTC)

	for _, loc := range locations {
		localTime := baseTime.In(loc)

		type TimeWrapper struct {
			Time time.Time `json:"time"`
		}

		wrapper := TimeWrapper{Time: localTime}
		jsonData, err := json.Marshal(wrapper)
		if err != nil {
			fmt.Printf("Marshal error: %v\n", err)
			continue
		}

		fmt.Printf("\nLocation: %s\n", loc)
		fmt.Printf("Time: %s\n", localTime)
		fmt.Printf("JSON: %s\n", string(jsonData))

		var result TimeWrapper
		if err := json.Unmarshal(jsonData, &result); err != nil {
			fmt.Printf("Unmarshal error: %v\n", err)
			continue
		}

		fmt.Printf("Decoded Time: %s\n", result.Time)
		fmt.Printf("Equal: %v\n", localTime.Equal(result.Time))
		fmt.Printf("Location preserved: %v\n", localTime.Location() == result.Time.Location())
	}
}

func main() {
	testDefaultTimeFormat()
	testCustomFormats()
	testDifferentDurationFormats()
	testTimeZones()
}
