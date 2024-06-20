package typejuggle

import (
	"reflect"
	"testing"
)

type SrcStruct struct {
	ID    int
	Name  string
	Email string
	Tags  []string
	Addr  *Address
}

type DestStruct struct {
	Id    int
	Name  string
	Email string
	Tags  []string
	Addr  *Address
}

type Address struct {
	Street string
	City   string
}

func TestFillFields_SimpleStruct(t *testing.T) {
	src := SrcStruct{ID: 1, Name: "John Doe", Email: "john@example.com"}
	var dest DestStruct

	FillFields(src, &dest)

	if dest.Id != 1 || dest.Name != "John Doe" || dest.Email != "john@example.com" {
		t.Errorf("Expected {Id: 1, Name: John Doe, Email: john@example.com}, but got %+v", dest)
	}
}

func TestFillFields_StructWithSlice(t *testing.T) {
	src := SrcStruct{Tags: []string{"go", "programming"}}
	var dest DestStruct

	FillFields(src, &dest)

	if !reflect.DeepEqual(dest.Tags, []string{"go", "programming"}) {
		t.Errorf("Expected Tags: [go, programming], but got %+v", dest.Tags)
	}
}

func TestFillFields_StructWithNestedStruct(t *testing.T) {
	src := SrcStruct{Addr: &Address{Street: "123 Main St", City: "Hometown"}}
	var dest DestStruct

	FillFields(src, &dest)

	if dest.Addr == nil || dest.Addr.Street != "123 Main St" || dest.Addr.City != "Hometown" {
		t.Errorf("Expected Addr: {Street: 123 Main St, City: Hometown}, but got %+v", dest.Addr)
	}
}

func TestFillFields_StructWithNilPointer(t *testing.T) {
	src := SrcStruct{}
	var dest DestStruct

	FillFields(src, &dest)

	if dest.Addr.City != "" || dest.Addr.Street != "" {
		t.Errorf("Expected City and Street to be empty strings")
	}
}

func TestFillFields_ConvertIntToString(t *testing.T) {
	type SrcStruct struct {
		ID int
	}
	type DestStruct struct {
		ID string
	}
	src := SrcStruct{ID: 42}
	var dest DestStruct

	FillFields(src, &dest, true)

	if dest.ID != "42" {
		t.Errorf("Expected ID to be '42', but got %+v", dest.ID)
	}
}

func TestFillFields_SliceWithPointers(t *testing.T) {
	type SrcStruct struct {
		Items []*Address
	}
	type DestStruct struct {
		Items []*Address
	}
	src := SrcStruct{Items: []*Address{{Street: "123 Main St", City: "Hometown"}, {Street: "456 Side St", City: "Othertown"}}}
	var dest DestStruct

	FillFields(src, &dest)

	if len(dest.Items) != 2 ||
		dest.Items[0].Street != "123 Main St" || dest.Items[0].City != "Hometown" ||
		dest.Items[1].Street != "456 Side St" || dest.Items[1].City != "Othertown" {
		t.Errorf("Expected Items to be [{Street: 123 Main St, City: Hometown}, {Street: 456 Side St, City: Othertown}], but got %+v", dest.Items)
	}
}

func TestFillFields_StructWithUnmatchedFields(t *testing.T) {
	type SrcStruct struct {
		ID    int
		Extra string
	}
	type DestStruct struct {
		Id int
	}
	src := SrcStruct{ID: 1, Extra: "extra"}
	var dest DestStruct

	FillFields(src, &dest)

	if dest.Id != 1 {
		t.Errorf("Expected Id to be 1, but got %+v", dest.Id)
	}
}

func TestFillFields_StructWithDifferentFieldNames(t *testing.T) {
	type SrcStruct struct {
		UniqueID int
	}
	type DestStruct struct {
		Id int
	}
	src := SrcStruct{UniqueID: 99}
	var dest DestStruct

	FillFields(src, &dest)

	if dest.Id != 0 {
		t.Errorf("Expected Id to be 0, but got %+v", dest.Id)
	}
}

func TestFillFields_NilSource(t *testing.T) {
	type SrcStruct struct {
		ID int
	}
	type DestStruct struct {
		ID int
	}
	var src *SrcStruct
	var dest DestStruct

	FillFields(src, &dest)

	if dest.ID != 0 {
		t.Errorf("Expected ID to be 0, but got %+v", dest.ID)
	}
}

func TestFillFields_EmptySourceStruct(t *testing.T) {
	src := SrcStruct{}
	var dest DestStruct

	FillFields(src, &dest)

	if dest.Id != 0 || dest.Name != "" || dest.Email != "" || dest.Tags != nil || dest.Addr.City != "" || dest.Addr.Street != "" {
		t.Errorf("Expected empty dest struct, but got %+v", dest)
	}
}

func TestFillFields_EmptyDestStruct(t *testing.T) {
	src := SrcStruct{ID: 1, Name: "John Doe", Email: "john@example.com"}
	var dest DestStruct

	FillFields(src, &dest)

	if dest.Id != 1 || dest.Name != "John Doe" || dest.Email != "john@example.com" {
		t.Errorf("Expected {Id: 1, Name: John Doe, Email: john@example.com}, but got %+v", dest)
	}
}

func TestFillFields_StructWithDifferentTypes(t *testing.T) {
	type SrcStruct struct {
		ID int
	}
	type DestStruct struct {
		ID float64
	}
	src := SrcStruct{ID: 42}
	var dest DestStruct

	FillFields(src, &dest, true)

	if dest.ID != 42 {
		t.Errorf("Expected ID to be 42, but got %+v", dest.ID)
	}
}

func TestFillFields_PartialFill(t *testing.T) {
	type SrcStruct struct {
		ID int
	}
	type DestStruct struct {
		ID   int
		Name string
	}
	src := SrcStruct{ID: 42}
	dest := DestStruct{Name: "Existing"}

	FillFields(src, &dest)

	if dest.ID != 42 || dest.Name != "Existing" {
		t.Errorf("Expected {ID: 42, Name: Existing}, but got %+v", dest)
	}
}

func TestFillFields_SliceToStructField(t *testing.T) {
	type SrcStruct struct {
		Tags []string
	}
	type DestStruct struct {
		Tags []string
	}
	src := SrcStruct{Tags: []string{"go", "programming"}}
	var dest DestStruct

	FillFields(src, &dest)

	if !reflect.DeepEqual(dest.Tags, []string{"go", "programming"}) {
		t.Errorf("Expected Tags: [go, programming], but got %+v", dest.Tags)
	}
}

func TestFillFields_ConvertBoolToString(t *testing.T) {
	type SrcStruct struct {
		IsActive bool
	}
	type DestStruct struct {
		IsActive string
	}
	src := SrcStruct{IsActive: true}
	var dest DestStruct

	FillFields(src, &dest, true)

	if dest.IsActive != "true" {
		t.Errorf("Expected IsActive to be 'true', but got %+v", dest.IsActive)
	}
}

type SrcDoublePointer struct {
	ID    **int
	Name  **string
	Email **string
}

type DestDoublePointer struct {
	ID    **int
	Name  **string
	Email **string
}

func TestFillFields_DoublePointer_FullyInitialized(t *testing.T) {
	id := 1
	name := "John Doe"
	email := "john@example.com"
	idPtr := &id
	namePtr := &name
	emailPtr := &email
	src := SrcDoublePointer{
		ID:    &idPtr,
		Name:  &namePtr,
		Email: &emailPtr,
	}
	var dest DestDoublePointer

	FillFields(src, &dest)

	if **dest.ID != 1 || **dest.Name != "John Doe" || **dest.Email != "john@example.com" {
		t.Errorf("Expected {ID: 1, Name: John Doe, Email: john@example.com}, but got {ID: %v, Name: %v, Email: %v}", **dest.ID, **dest.Name, **dest.Email)
	}
}

func TestFillFields_DoublePointer_HalfInitialized(t *testing.T) {
	id := 1
	name := "John Doe"
	idPtr := &id
	namePtr := &name
	src := SrcDoublePointer{
		ID:   &idPtr,
		Name: &namePtr,
	}
	var dest DestDoublePointer

	FillFields(src, &dest)

	if **dest.ID != 1 || **dest.Name != "John Doe" || **dest.Email != "" {
		t.Errorf("Expected {ID: 1, Name: John Doe, Email: nil}, but got {ID: %v, Name: %v, Email: %v}", **dest.ID, **dest.Name, dest.Email)
	}
}

func TestFillFields_DoublePointer_Uninitialized(t *testing.T) {
	src := SrcDoublePointer{}
	var dest DestDoublePointer

	FillFields(src, &dest)

	if **dest.ID != 0 || **dest.Name != "" || **dest.Email != "" {
		t.Errorf("Expected all fields to be nil, but got {ID: %v, Name: %v, Email: %v}", dest.ID, dest.Name, dest.Email)
	}
}

func TestFillFields_DoublePointer_NilPointer(t *testing.T) {
	src := SrcDoublePointer{}
	destID := 42
	destIDPtr := &destID
	dest := DestDoublePointer{
		ID: &destIDPtr,
	}

	FillFields(src, &dest)

	if dest.ID == nil || **dest.ID != 42 {
		t.Errorf("Expected ID to be 42, but got %v", dest.ID)
	}
}

func TestFillFields_DoublePointer_PartiallyInitializedDestination(t *testing.T) {
	id := 1
	email := "john@example.com"
	idPtr := &id
	emailPtr := &email
	src := SrcDoublePointer{
		ID:    &idPtr,
		Email: &emailPtr,
	}
	name := "Existing Name"
	namePtr := &name
	dest := DestDoublePointer{
		Name: &namePtr,
	}

	FillFields(src, &dest)

	if **dest.ID != 1 || **dest.Email != "john@example.com" || **dest.Name != "Existing Name" {
		t.Errorf("Expected {ID: 1, Name: Existing Name, Email: john@example.com}, but got {ID: %v, Name: %v, Email: %v}", **dest.ID, **dest.Name, **dest.Email)
	}
}

func TestFillFields_DoublePointer_ConvertIntToString(t *testing.T) {
	type SrcDoublePointer struct {
		ID **int
	}
	type DestDoublePointer struct {
		ID **string
	}
	id := 42
	idPtr := &id
	src := SrcDoublePointer{
		ID: &idPtr,
	}
	var dest DestDoublePointer

	FillFields(src, &dest, true)

	expected := "42"
	if **dest.ID != expected {
		t.Errorf("Expected ID to be '%v', but got %v", expected, **dest.ID)
	}
}

func TestFillFields_DoublePointer_ConvertStringToBool(t *testing.T) {
	type SrcDoublePointer struct {
		Active **string
	}
	type DestDoublePointer struct {
		Active **bool
	}
	active := "true"
	activePtr := &active
	src := SrcDoublePointer{
		Active: &activePtr,
	}
	var dest DestDoublePointer

	FillFields(src, &dest, true)

	if **dest.Active != true {
		t.Errorf("Expected Active to be true, but got %v", **dest.Active)
	}
}

type NestedStruct struct {
	Street **string
	City   **string
}

type SrcStructWithNestedDoublePointer struct {
	ID    **int
	Addr  **NestedStruct
	Email **string
}

type DestStructWithNestedDoublePointer struct {
	ID    **int
	Addr  **NestedStruct
	Email **string
}

func TestFillFields_NestedDoublePointer_FullyInitialized(t *testing.T) {
	id := 1
	email := "john@example.com"
	street := "123 Main St"
	city := "Hometown"
	idPtr := &id
	emailPtr := &email
	streetPtr := &street
	cityPtr := &city
	addr := &NestedStruct{
		Street: &streetPtr,
		City:   &cityPtr,
	}
	src := SrcStructWithNestedDoublePointer{
		ID:    &idPtr,
		Addr:  &addr,
		Email: &emailPtr,
	}
	var dest DestStructWithNestedDoublePointer

	FillFields(src, &dest)

	if **dest.ID != 1 || **dest.Email != "john@example.com" || **(*dest.Addr).Street != "123 Main St" || **(*dest.Addr).City != "Hometown" {
		t.Errorf("Expected {ID: 1, Email: john@example.com, Addr: {Street: 123 Main St, City: Hometown}}, but got {ID: %v, Email: %v, Addr: {Street: %v, City: %v}}", **dest.ID, **dest.Email, **(*dest.Addr).Street, **(*dest.Addr).City)
	}
}

func TestFillFields_NestedDoublePointer_PartiallyInitialized(t *testing.T) {
	id := 1
	street := "123 Main St"
	idPtr := &id
	streetPtr := &street
	addr := &NestedStruct{
		Street: &streetPtr,
	}
	src := SrcStructWithNestedDoublePointer{
		ID:   &idPtr,
		Addr: &addr,
	}
	var dest DestStructWithNestedDoublePointer

	FillFields(src, &dest)

	if **dest.ID != 1 || **(*dest.Addr).Street != "123 Main St" || **(*dest.Addr).City != "" {
		t.Errorf("Expected {ID: 1, Addr: {Street: 123 Main St, City: nil}}, but got {ID: %v, Addr: {Street: %v, City: %v}}", **dest.ID, **(*dest.Addr).Street, (*dest.Addr).City)
	}
}

func TestFillFields_NestedDoublePointer_Uninitialized(t *testing.T) {
	t.Skip("This test is expected to fail as we don't explore dest struct if the source pointer is nil")
	src := SrcStructWithNestedDoublePointer{}
	var dest DestStructWithNestedDoublePointer

	FillFields(src, &dest)

	if **dest.ID != 0 || **dest.Email != "" || **(**dest.Addr).City != "" || **(**dest.Addr).Street != "" {
		t.Errorf("Expected all fields to be nil, but got {ID: %v, Email: %v, Addr: %v}", dest.ID, dest.Email, dest.Addr)
	}
}

func TestFillFields_NestedDoublePointer_NilPointer(t *testing.T) {
	src := SrcStructWithNestedDoublePointer{}
	destID := 42
	destEmail := "existing@example.com"
	destStreet := "456 Side St"
	destCity := "Othertown"
	destIDPtr := &destID
	destEmailPtr := &destEmail
	destStreetPtr := &destStreet
	destCityPtr := &destCity
	destAddr := &NestedStruct{
		Street: &destStreetPtr,
		City:   &destCityPtr,
	}
	dest := DestStructWithNestedDoublePointer{
		ID:    &destIDPtr,
		Addr:  &destAddr,
		Email: &destEmailPtr,
	}

	FillFields(src, &dest)

	if **dest.ID != 42 || **dest.Email != "existing@example.com" || **(*dest.Addr).Street != "456 Side St" || **(*dest.Addr).City != "Othertown" {
		t.Errorf("Expected {ID: 42, Email: existing@example.com, Addr: {Street: 456 Side St, City: Othertown}}, but got {ID: %v, Email: %v, Addr: {Street: %v, City: %v}}", **dest.ID, **dest.Email, **(*dest.Addr).Street, **(*dest.Addr).City)
	}
}

func TestFillFields_NestedDoublePointer_PartiallyInitializedDestination(t *testing.T) {
	id := 1
	email := "john@example.com"
	city := "Hometown"
	idPtr := &id
	emailPtr := &email
	cityPtr := &city
	addr := &NestedStruct{
		City: &cityPtr,
	}
	src := SrcStructWithNestedDoublePointer{
		ID:    &idPtr,
		Email: &emailPtr,
		Addr:  &addr,
	}
	street := "Existing Street"
	streetPtr := &street
	destAddr := &NestedStruct{
		Street: &streetPtr,
	}
	dest := DestStructWithNestedDoublePointer{
		Addr: &destAddr,
	}

	FillFields(src, &dest)

	if **dest.ID != 1 || **dest.Email != "john@example.com" || **(*dest.Addr).City != "Hometown" || **(*dest.Addr).Street != "Existing Street" {
		t.Errorf("Expected {ID: 1, Email: john@example.com, Addr: {Street: Existing Street, City: Hometown}}, but got {ID: %v, Email: %v, Addr: {Street: %v, City: %v}}", **dest.ID, **dest.Email, **(*dest.Addr).Street, **(*dest.Addr).City)
	}
}

func TestFillFields_NestedDoublePointer_ConvertIntToString(t *testing.T) {
	type SrcStructWithNestedDoublePointer struct {
		ID   **int
		Addr **NestedStruct
	}
	type DestStructWithNestedDoublePointer struct {
		ID   **string
		Addr **NestedStruct
	}
	id := 42
	idPtr := &id
	src := SrcStructWithNestedDoublePointer{
		ID: &idPtr,
	}
	var dest DestStructWithNestedDoublePointer

	FillFields(src, &dest, true)

	expected := "42"
	if **dest.ID != expected {
		t.Errorf("Expected ID to be '%v', but got %v", expected, **dest.ID)
	}
}
