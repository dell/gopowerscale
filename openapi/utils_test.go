/*
Copyright (c) 2025 Dell Inc, or its subsidiaries.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package openapi

import (
	"encoding/json"
	"testing"
	"time"
)

func TestPtrBool(t *testing.T) {
	val := true
	ptr := PtrBool(val)
	if *ptr != val {
		t.Errorf("PtrBool() = %v, want %v", *ptr, val)
	}
}

func TestPtrInt(t *testing.T) {
	val := 42
	ptr := PtrInt(val)
	if *ptr != val {
		t.Errorf("PtrInt() = %v, want %v", *ptr, val)
	}
}

func TestPtrInt32(t *testing.T) {
	v := int32(42)
	ptr := PtrInt32(v)
	if *ptr != v {
		t.Errorf("PtrInt32(%d) = %d; want %d", v, *ptr, v)
	}
}

func TestPtrInt64(t *testing.T) {
	v := int64(42)
	ptr := PtrInt64(v)
	if *ptr != v {
		t.Errorf("PtrInt64(%d) = %d; want %d", v, *ptr, v)
	}
}

func TestPtrFloat32(t *testing.T) {
	v := float32(42.0)
	ptr := PtrFloat32(v)
	if *ptr != v {
		t.Errorf("PtrFloat32(%f) = %f; want %f", v, *ptr, v)
	}
}

func TestPtrFloat64(t *testing.T) {
	v := float64(42.0)
	ptr := PtrFloat64(v)
	if *ptr != v {
		t.Errorf("PtrFloat64(%f) = %f; want %f", v, *ptr, v)
	}
}

func TestPtrString(t *testing.T) {
	v := "hello"
	ptr := PtrString(v)
	if *ptr != v {
		t.Errorf("PtrString(%s) = %s; want %s", v, *ptr, v)
	}
}

func TestPtrTime(t *testing.T) {
	v := time.Now()
	ptr := PtrTime(v)
	if *ptr != v {
		t.Errorf("PtrTime(%v) = %v; want %v", v, *ptr, v)
	}
}

func TestNullableBool(t *testing.T) {
	val := true
	nb := NewNullableBool(&val)

	if !nb.IsSet() {
		t.Errorf("NewNullableBool().IsSet() = %v, want %v", nb.IsSet(), true)
	}

	if *nb.Get() != val {
		t.Errorf("NewNullableBool().Get() = %v, want %v", *nb.Get(), val)
	}

	nb.Unset()
	if nb.IsSet() {
		t.Errorf("NullableBool.Unset().IsSet() = %v, want %v", nb.IsSet(), false)
	}

	jsonVal, err := json.Marshal(nb)
	if err != nil {
		t.Errorf("NullableBool.MarshalJSON() error = %v", err)
	}
	if string(jsonVal) != "null" {
		t.Errorf("NullableBool.MarshalJSON() = %v, want %v", string(jsonVal), "null")
	}
}

func TestNullableBool_UnmarshalJSON(t *testing.T) {
	var nb NullableBool
	jsonData := []byte("true")
	err := nb.UnmarshalJSON(jsonData)
	if err != nil {
		t.Errorf("NullableBool.UnmarshalJSON() error = %v", err)
	}
	if !nb.IsSet() {
		t.Errorf("NullableBool.UnmarshalJSON().IsSet() = %v, want %v", nb.IsSet(), true)
	}
	if *nb.Get() != true {
		t.Errorf("NullableBool.UnmarshalJSON().Get() = %v, want %v", *nb.Get(), true)
	}
}

func TestNullableInt(t *testing.T) {
	val := 42
	ni := NullableInt{}
	ni.Set(&val)

	if !ni.IsSet() {
		t.Errorf("NullableInt.Set().IsSet() = %v, want %v", ni.IsSet(), true)
	}

	if *ni.Get() != val {
		t.Errorf("NullableInt.Get() = %v, want %v", *ni.Get(), val)
	}

	ni.Unset()
	if ni.IsSet() {
		t.Errorf("NullableInt.Unset().IsSet() = %v, want %v", ni.IsSet(), false)
	}

	jsonVal, err := json.Marshal(ni)
	if err != nil {
		t.Errorf("NullableInt.MarshalJSON() error = %v", err)
	}
	if string(jsonVal) != "null" {
		t.Errorf("NullableInt.MarshalJSON() = %v, want %v", string(jsonVal), "null")
	}

	err = json.Unmarshal([]byte("42"), &ni)
	if err != nil {
		t.Errorf("NullableInt.UnmarshalJSON() error = %v", err)
	}
	if !ni.IsSet() {
		t.Errorf("NullableInt.UnmarshalJSON().IsSet() = %v, want %v", ni.IsSet(), true)
	}
	if *ni.Get() != 42 {
		t.Errorf("NullableInt.UnmarshalJSON().Get() = %v, want %v", *ni.Get(), 42)
	}
}

func TestNullableBool_Set(t *testing.T) {
	var nb NullableBool

	// Test setting a true value
	trueVal := true
	nb.Set(&trueVal)
	if nb.value == nil || *nb.value != trueVal {
		t.Errorf("Expected value to be %v, got %v", trueVal, nb.value)
	}
	if !nb.isSet {
		t.Errorf("Expected isSet to be true, got %v", nb.isSet)
	}

	// Test setting a false value
	falseVal := false
	nb.Set(&falseVal)
	if nb.value == nil || *nb.value != falseVal {
		t.Errorf("Expected value to be %v, got %v", falseVal, nb.value)
	}
	if !nb.isSet {
		t.Errorf("Expected isSet to be true, got %v", nb.isSet)
	}

	// Test setting a nil value
	nb.Set(nil)
	if nb.value != nil {
		t.Errorf("Expected value to be nil, got %v", nb.value)
	}
	if !nb.isSet {
		t.Errorf("Expected isSet to be true, got %v", nb.isSet)
	}
}

func TestNullableInt32(t *testing.T) {
	val := int32(42)
	ni32 := NullableInt32{}
	ni32.Set(&val)

	if !ni32.IsSet() {
		t.Errorf("NullableInt32.Set().IsSet() = %v, want %v", ni32.IsSet(), true)
	}

	if *ni32.Get() != val {
		t.Errorf("NullableInt32.Get() = %v, want %v", *ni32.Get(), val)
	}

	ni32.Unset()
	if ni32.IsSet() {
		t.Errorf("NullableInt32.Unset().IsSet() = %v, want %v", ni32.IsSet(), false)
	}

	jsonVal, err := json.Marshal(ni32)
	if err != nil {
		t.Errorf("NullableInt32.MarshalJSON() error = %v", err)
	}
	if string(jsonVal) != "null" {
		t.Errorf("NullableInt32.MarshalJSON() = %v, want %v", string(jsonVal), "null")
	}

	err = json.Unmarshal([]byte("42"), &ni32)
	if err != nil {
		t.Errorf("NullableInt32.UnmarshalJSON() error = %v", err)
	}
	if !ni32.IsSet() {
		t.Errorf("NullableInt32.UnmarshalJSON().IsSet() = %v, want %v", ni32.IsSet(), true)
	}
	if *ni32.Get() != 42 {
		t.Errorf("NullableInt32.UnmarshalJSON().Get() = %v, want %v", *ni32.Get(), 42)
	}
}

func TestNullableInt64(t *testing.T) {
	val := int64(42)
	ni64 := NullableInt64{}
	ni64.Set(&val)

	if !ni64.IsSet() {
		t.Errorf("NullableInt64.Set().IsSet() = %v, want %v", ni64.IsSet(), true)
	}

	if *ni64.Get() != val {
		t.Errorf("NullableInt64.Get() = %v, want %v", *ni64.Get(), val)
	}

	ni64.Unset()
	if ni64.IsSet() {
		t.Errorf("NullableInt64.Unset().IsSet() = %v, want %v", ni64.IsSet(), false)
	}

	jsonVal, err := json.Marshal(ni64)
	if err != nil {
		t.Errorf("NullableInt64.MarshalJSON() error = %v", err)
	}
	if string(jsonVal) != "null" {
		t.Errorf("NullableInt64.MarshalJSON() = %v, want %v", string(jsonVal), "null")
	}

	err = json.Unmarshal([]byte("42"), &ni64)
	if err != nil {
		t.Errorf("NullableInt64.UnmarshalJSON() error = %v", err)
	}
	if !ni64.IsSet() {
		t.Errorf("NullableInt64.UnmarshalJSON().IsSet() = %v, want %v", ni64.IsSet(), true)
	}
	if *ni64.Get() != 42 {
		t.Errorf("NullableInt64.UnmarshalJSON().Get() = %v, want %v", *ni64.Get(), 42)
	}
}

func TestNullableFloat32(t *testing.T) {
	val := float32(42.0)
	nf32 := NullableFloat32{}
	nf32.Set(&val)

	if !nf32.IsSet() {
		t.Errorf("NullableFloat32.Set().IsSet() = %v, want %v", nf32.IsSet(), true)
	}

	if *nf32.Get() != val {
		t.Errorf("NullableFloat32.Get() = %v, want %v", *nf32.Get(), val)
	}

	nf32.Unset()
	if nf32.IsSet() {
		t.Errorf("NullableFloat32.Unset().IsSet() = %v, want %v", nf32.IsSet(), false)
	}

	jsonVal, err := json.Marshal(nf32)
	if err != nil {
		t.Errorf("NullableFloat32.MarshalJSON() error = %v", err)
	}
	if string(jsonVal) != "null" {
		t.Errorf("NullableFloat32.MarshalJSON() = %v, want %v", string(jsonVal), "null")
	}

	err = json.Unmarshal([]byte("42.0"), &nf32)
	if err != nil {
		t.Errorf("NullableFloat32.UnmarshalJSON() error = %v", err)
	}
	if !nf32.IsSet() {
		t.Errorf("NullableFloat32.UnmarshalJSON().IsSet() = %v, want %v", nf32.IsSet(), true)
	}
	if *nf32.Get() != 42.0 {
		t.Errorf("NullableFloat32.UnmarshalJSON().Get() = %v, want %v", *nf32.Get(), 42.0)
	}
}

func TestNullableFloat64(t *testing.T) {
	val := float64(42.0)
	nf64 := NullableFloat64{}
	nf64.Set(&val)

	if !nf64.IsSet() {
		t.Errorf("NullableFloat64.Set().IsSet() = %v, want %v", nf64.IsSet(), true)
	}

	if *nf64.Get() != val {
		t.Errorf("NullableFloat64.Get() = %v, want %v", *nf64.Get(), val)
	}

	nf64.Unset()
	if nf64.IsSet() {
		t.Errorf("NullableFloat64.Unset().IsSet() = %v, want %v", nf64.IsSet(), false)
	}

	jsonVal, err := json.Marshal(nf64)
	if err != nil {
		t.Errorf("NullableFloat64.MarshalJSON() error = %v", err)
	}
	if string(jsonVal) != "null" {
		t.Errorf("NullableFloat64.MarshalJSON() = %v, want %v", string(jsonVal), "null")
	}

	err = json.Unmarshal([]byte("42.0"), &nf64)
	if err != nil {
		t.Errorf("NullableFloat64.UnmarshalJSON() error = %v", err)
	}
	if !nf64.IsSet() {
		t.Errorf("NullableFloat64.UnmarshalJSON().IsSet() = %v, want %v", nf64.IsSet(), true)
	}
	if *nf64.Get() != 42.0 {
		t.Errorf("NullableFloat64.UnmarshalJSON().Get() = %v, want %v", *nf64.Get(), 42.0)
	}
}

func TestNullableBool_MarshalJSON(t *testing.T) {
	trueValue := true
	nullableBool := NewNullableBool(&trueValue)

	data, err := nullableBool.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON failed: %v", err)
	}

	expected := `true`
	if string(data) != expected {
		t.Errorf("Expected JSON %s, got %s", expected, string(data))
	}

	// Test null value
	nullableBool.Unset()
	data, err = nullableBool.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON failed: %v", err)
	}

	expected = `null`
	if string(data) != expected {
		t.Errorf("Expected JSON %s, got %s", expected, string(data))
	}
}

func TestNullableInt_Get(t *testing.T) {
	value := 42
	nullableInt := NewNullableInt(&value)

	got := nullableInt.Get()
	if got == nil || *got != value {
		t.Errorf("Expected Get() to return %d, got %v", value, got)
	}

	nullableInt.Unset()
	got = nullableInt.Get()
	if got != nil {
		t.Errorf("Expected Get() to return nil after Unset(), got %v", *got)
	}
}

func TestNullableInt_Set(t *testing.T) {
	value := 42
	nullableInt := NullableInt{}

	nullableInt.Set(&value)
	got := nullableInt.Get()
	if got == nil || *got != value {
		t.Errorf("Expected Set() to set value %d, got %v", value, got)
	}

	if !nullableInt.IsSet() {
		t.Errorf("Expected IsSet() to return true after Set(), got false")
	}
}

func TestNullableInt_IsSet(t *testing.T) {
	nullableInt := NullableInt{}

	if nullableInt.IsSet() {
		t.Errorf("Expected IsSet() to return false for unset NullableInt, got true")
	}

	value := 42
	nullableInt.Set(&value)

	if !nullableInt.IsSet() {
		t.Errorf("Expected IsSet() to return true after Set(), got false")
	}
}

func TestNullableInt_Unset(t *testing.T) {
	value := 42
	nullableInt := NewNullableInt(&value)

	nullableInt.Unset()
	if nullableInt.IsSet() {
		t.Errorf("Expected IsSet() to return false after Unset(), got true")
	}

	if nullableInt.Get() != nil {
		t.Errorf("Expected Get() to return nil after Unset(), got %v", nullableInt.Get())
	}
}

func TestNullableInt32_Get(t *testing.T) {
	value := int32(42)
	nullableInt32 := NewNullableInt32(&value)

	got := nullableInt32.Get()
	if got == nil || *got != value {
		t.Errorf("Expected Get() to return %d, got %v", value, got)
	}

	nullableInt32.Unset()
	got = nullableInt32.Get()
	if got != nil {
		t.Errorf("Expected Get() to return nil after Unset(), got %v", got)
	}
}

func TestNullableInt32_Set(t *testing.T) {
	value := int32(42)
	nullableInt32 := NullableInt32{}

	nullableInt32.Set(&value)
	got := nullableInt32.Get()
	if got == nil || *got != value {
		t.Errorf("Expected Set() to set value %d, got %v", value, got)
	}

	if !nullableInt32.IsSet() {
		t.Errorf("Expected IsSet() to return true after Set(), got false")
	}
}

func TestNullableInt32_IsSet(t *testing.T) {
	nullableInt32 := NullableInt32{}

	if nullableInt32.IsSet() {
		t.Errorf("Expected IsSet() to return false for unset NullableInt32, got true")
	}

	value := int32(42)
	nullableInt32.Set(&value)

	if !nullableInt32.IsSet() {
		t.Errorf("Expected IsSet() to return true after Set(), got false")
	}
}

func TestNullableInt32_Unset(t *testing.T) {
	value := int32(42)
	nullableInt32 := NewNullableInt32(&value)

	nullableInt32.Unset()
	if nullableInt32.IsSet() {
		t.Errorf("Expected IsSet() to return false after Unset(), got true")
	}

	if nullableInt32.Get() != nil {
		t.Errorf("Expected Get() to return nil after Unset(), got %v", nullableInt32.Get())
	}
}

func TestNullableInt32_MarshalJSON(t *testing.T) {
	value := int32(42)
	nullableInt32 := NewNullableInt32(&value)

	data, err := nullableInt32.MarshalJSON()
	if err != nil {
		t.Fatalf("Unexpected error during MarshalJSON: %v", err)
	}

	expected := "42"
	if string(data) != expected {
		t.Errorf("Expected MarshalJSON output to be %s, got %s", expected, string(data))
	}
}

func TestNullableInt32_UnmarshalJSON(t *testing.T) {
	data := []byte("42")
	nullableInt32 := NullableInt32{}

	err := nullableInt32.UnmarshalJSON(data)
	if err != nil {
		t.Fatalf("Unexpected error during UnmarshalJSON: %v", err)
	}

	got := nullableInt32.Get()
	if got == nil || *got != 42 {
		t.Errorf("Expected UnmarshalJSON to set value to 42, got %v", got)
	}

	if !nullableInt32.IsSet() {
		t.Errorf("Expected IsSet() to return true after UnmarshalJSON, got false")
	}
}

func TestNullableInt64_Get(t *testing.T) {
	value := int64(64)
	nullableInt64 := NewNullableInt64(&value)

	got := nullableInt64.Get()
	if got == nil || *got != value {
		t.Errorf("Expected Get() to return %d, got %v", value, got)
	}

	nullableInt64.Unset()
	got = nullableInt64.Get()
	if got != nil {
		t.Errorf("Expected Get() to return nil after Unset(), got %v", got)
	}
}

func TestNullableInt64_Set(t *testing.T) {
	value := int64(64)
	nullableInt64 := NullableInt64{}

	nullableInt64.Set(&value)
	got := nullableInt64.Get()
	if got == nil || *got != value {
		t.Errorf("Expected Set() to set value %d, got %v", value, got)
	}

	if !nullableInt64.IsSet() {
		t.Errorf("Expected IsSet() to return true after Set(), got false")
	}
}

func TestNullableInt64_IsSet(t *testing.T) {
	nullableInt64 := NullableInt64{}

	if nullableInt64.IsSet() {
		t.Errorf("Expected IsSet() to return false for unset NullableInt64, got true")
	}

	value := int64(64)
	nullableInt64.Set(&value)

	if !nullableInt64.IsSet() {
		t.Errorf("Expected IsSet() to return true after Set(), got false")
	}
}

func TestNullableInt64_Unset(t *testing.T) {
	value := int64(64)
	nullableInt64 := NewNullableInt64(&value)

	nullableInt64.Unset()
	if nullableInt64.IsSet() {
		t.Errorf("Expected IsSet() to return false after Unset(), got true")
	}

	if nullableInt64.Get() != nil {
		t.Errorf("Expected Get() to return nil after Unset(), got %v", nullableInt64.Get())
	}
}

func TestNewNullableInt64(t *testing.T) {
	val := int64(42)
	nullableInt := NewNullableInt64(&val)

	if nullableInt == nil {
		t.Fatal("Expected NewNullableInt64 to return a non-nil value")
	}
	if *nullableInt.value != val {
		t.Fatalf("Expected %d, got %d", val, *nullableInt.value)
	}
	if !nullableInt.isSet {
		t.Fatal("Expected isSet to be true")
	}
}

// Test MarshalJSON for NullableInt64
func TestNullableInt64_MarshalJSON(t *testing.T) {
	val := int64(42)
	nullableInt := NewNullableInt64(&val)
	data, err := nullableInt.MarshalJSON()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expected := "42"
	if string(data) != expected {
		t.Fatalf("Expected %s, got %s", expected, string(data))
	}
}

// Test UnmarshalJSON for NullableInt64
func TestNullableInt64_UnmarshalJSON(t *testing.T) {
	nullableInt := &NullableInt64{}
	data := []byte("42")
	err := nullableInt.UnmarshalJSON(data)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expected := int64(42)
	if *nullableInt.value != expected {
		t.Fatalf("Expected %d, got %d", expected, *nullableInt.value)
	}
	if !nullableInt.isSet {
		t.Fatal("Expected isSet to be true")
	}
}

// Test MarshalJSON for NullableFloat32
func TestNullableFloat32_MarshalJSON(t *testing.T) {
	val := float32(42.42)
	nullableFloat := NewNullableFloat32(&val)
	data, err := nullableFloat.MarshalJSON()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expected := "42.42"
	if string(data) != expected {
		t.Fatalf("Expected %s, got %s", expected, string(data))
	}
}

// Test UnmarshalJSON for NullableFloat32
func TestNullableFloat32_UnmarshalJSON(t *testing.T) {
	nullableFloat := &NullableFloat32{}
	data := []byte("42.42")
	err := nullableFloat.UnmarshalJSON(data)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expected := float32(42.42)
	if *nullableFloat.Get() != expected {
		t.Fatalf("Expected %f, got %f", expected, *nullableFloat.Get())
	}
	if !nullableFloat.IsSet() {
		t.Fatal("Expected IsSet to be true")
	}
}

// Test MarshalJSON for NullableFloat64
func TestNullableFloat64_MarshalJSON(t *testing.T) {
	val := float64(42.42)
	nullableFloat := NewNullableFloat64(&val)
	data, err := nullableFloat.MarshalJSON()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expected := "42.42"
	if string(data) != expected {
		t.Fatalf("Expected %s, got %s", expected, string(data))
	}
}

// Test UnmarshalJSON for NullableFloat64
func TestNullableFloat64_UnmarshalJSON(t *testing.T) {
	nullableFloat := &NullableFloat64{}
	data := []byte("42.42")
	err := nullableFloat.UnmarshalJSON(data)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expected := float64(42.42)
	if *nullableFloat.Get() != expected {
		t.Fatalf("Expected %f, got %f", expected, *nullableFloat.Get())
	}
	if !nullableFloat.IsSet() {
		t.Fatal("Expected IsSet to be true")
	}
}

func TestNewNullableString(t *testing.T) {
	val := "Hello, World!"
	nullableString := NewNullableString(&val)

	if nullableString == nil {
		t.Fatal("Expected NewNullableString to return a non-nil value")
	}
	if *nullableString.value != val {
		t.Fatalf("Expected %s, got %s", val, *nullableString.value)
	}
	if !nullableString.isSet {
		t.Fatal("Expected isSet to be true")
	}
}

// Test MarshalJSON for NullableString
func TestNullableString_MarshalJSON(t *testing.T) {
	val := "Hello, World!"
	nullableString := NewNullableString(&val)
	data, err := nullableString.MarshalJSON()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expected := `"Hello, World!"`
	if string(data) != expected {
		t.Fatalf("Expected %s, got %s", expected, string(data))
	}
}

// Test UnmarshalJSON for NullableString
func TestNullableString_UnmarshalJSON(t *testing.T) {
	nullableString := &NullableString{}
	data := []byte(`"Hello, World!"`)
	err := nullableString.UnmarshalJSON(data)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expected := "Hello, World!"
	if *nullableString.value != expected {
		t.Fatalf("Expected %s, got %s", expected, *nullableString.value)
	}
	if !nullableString.isSet {
		t.Fatal("Expected isSet to be true")
	}
}

// Test NewNullableTime constructor
func TestNewNullableTime(t *testing.T) {
	val := time.Now()
	nullableTime := NewNullableTime(&val)

	if nullableTime == nil {
		t.Fatal("Expected NewNullableTime to return a non-nil value")
	}
	if nullableTime.value == nil {
		t.Fatal("Expected value to be set")
	}
	if !nullableTime.isSet {
		t.Fatal("Expected isSet to be true")
	}
}

// Test MarshalJSON for NullableTime
func TestNullableTime_MarshalJSON(t *testing.T) {
	// Get current time and truncate nanoseconds
	val := time.Now().Truncate(time.Second) // Truncate to second precision
	nullableTime := NewNullableTime(&val)
	data, err := nullableTime.MarshalJSON()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Format the time to second precision
	expected := `"` + val.Format("2006-01-02T15:04:05-07:00") + `"`
	// ignore suffix since this can fail depending on env.
	// in GH actions we see: expected ends in +00:00 but string(data) ends in Z
	if string(data)[0:20] != expected[0:20] {
		t.Fatalf("Expected %s, got %s", expected[0:20], string(data)[0:20])
	}
}

// Test UnmarshalJSON for NullableTime
func TestNullableTime_UnmarshalJSON(t *testing.T) {
	nullableTime := &NullableTime{}
	data := []byte(`"2025-01-01T00:00:00Z"`)
	err := nullableTime.UnmarshalJSON(data)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expected := "2025-01-01T00:00:00Z"
	if nullableTime.value.Format(time.RFC3339) != expected {
		t.Fatalf("Expected %s, got %s", expected, nullableTime.value.Format(time.RFC3339))
	}
	if !nullableTime.isSet {
		t.Fatal("Expected isSet to be true")
	}
}

// Test IsNil function
func TestIsNil(t *testing.T) {
	var nilPtr *int
	if !IsNil(nilPtr) {
		t.Fatal("Expected IsNil(nilPtr) to return true")
	}

	nonNilPtr := new(int)
	if IsNil(nonNilPtr) {
		t.Fatal("Expected IsNil(nonNilPtr) to return false")
	}

	var emptySlice []int
	if !IsNil(emptySlice) {
		t.Fatal("Expected IsNil(emptySlice) to return true")
	}

	nonEmptySlice := []int{1}
	if IsNil(nonEmptySlice) {
		t.Fatal("Expected IsNil(nonEmptySlice) to return false")
	}
}

func TestNullableString_Set(t *testing.T) {
	var ns NullableString

	// Test setting a non-nil value
	str := "hello"
	ns.Set(&str)
	if ns.value == nil || *ns.value != str {
		t.Errorf("Expected value to be %v, got %v", str, ns.value)
	}
	if !ns.isSet {
		t.Errorf("Expected isSet to be true, got %v", ns.isSet)
	}

	// Test setting a nil value
	ns.Set(nil)
	if ns.value != nil {
		t.Errorf("Expected value to be nil, got %v", ns.value)
	}
	if !ns.isSet {
		t.Errorf("Expected isSet to be true, got %v", ns.isSet)
	}
}

func TestNullableString_Get(t *testing.T) {
	var ns NullableString

	// Test getting a non-nil value
	str := "hello"
	ns.Set(&str)
	if ns.Get() == nil || *ns.Get() != str {
		t.Errorf("Expected Get() to return %v, got %v", str, ns.Get())
	}

	// Test getting a nil value
	ns.Unset()
	if ns.Get() != nil {
		t.Errorf("Expected Get() to return nil, got %v", ns.Get())
	}
}

func TestNullableString_IsSet(t *testing.T) {
	var ns NullableString

	// Test IsSet when value is set
	str := "hello"
	ns.Set(&str)
	if !ns.IsSet() {
		t.Errorf("Expected IsSet() to be true, got %v", ns.IsSet())
	}

	// Test IsSet when value is unset
	ns.Unset()
	if ns.IsSet() {
		t.Errorf("Expected IsSet() to be false, got %v", ns.IsSet())
	}
}

func TestNullableString_Unset(t *testing.T) {
	var ns NullableString

	// Test Unset
	str := "hello"
	ns.Set(&str)
	ns.Unset()
	if ns.value != nil {
		t.Errorf("Expected value to be nil, got %v", ns.value)
	}
	if ns.isSet {
		t.Errorf("Expected isSet to be false, got %v", ns.isSet)
	}
}

func TestNullableTime_Set(t *testing.T) {
	var nt NullableTime

	// Test setting a non-nil value
	now := time.Now()
	nt.Set(&now)
	if nt.value == nil || *nt.value != now {
		t.Errorf("Expected value to be %v, got %v", now, nt.value)
	}
	if !nt.isSet {
		t.Errorf("Expected isSet to be true, got %v", nt.isSet)
	}

	// Test setting a nil value
	nt.Set(nil)
	if nt.value != nil {
		t.Errorf("Expected value to be nil, got %v", nt.value)
	}
	if !nt.isSet {
		t.Errorf("Expected isSet to be true, got %v", nt.isSet)
	}
}

func TestNullableTime_Get(t *testing.T) {
	var nt NullableTime

	// Test getting a non-nil value
	now := time.Now()
	nt.Set(&now)
	if nt.Get() == nil || *nt.Get() != now {
		t.Errorf("Expected Get() to return %v, got %v", now, nt.Get())
	}

	// Test getting a nil value
	nt.Unset()
	if nt.Get() != nil {
		t.Errorf("Expected Get() to return nil, got %v", nt.Get())
	}
}

func TestNullableTime_IsSet(t *testing.T) {
	var nt NullableTime

	// Test IsSet when value is set
	now := time.Now()
	nt.Set(&now)
	if !nt.IsSet() {
		t.Errorf("Expected IsSet() to be true, got %v", nt.IsSet())
	}

	// Test IsSet when value is unset
	nt.Unset()
	if nt.IsSet() {
		t.Errorf("Expected IsSet() to be false, got %v", nt.IsSet())
	}
}

func TestNullableTime_Unset(t *testing.T) {
	var nt NullableTime

	// Test Unset
	now := time.Now()
	nt.Set(&now)
	nt.Unset()
	if nt.value != nil {
		t.Errorf("Expected value to be nil, got %v", nt.value)
	}
	if nt.isSet {
		t.Errorf("Expected isSet to be false, got %v", nt.isSet)
	}
}
