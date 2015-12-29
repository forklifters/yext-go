package yext

import (
	"fmt"
	"reflect"
	"testing"
)

var examplePhoto LocationPhoto = LocationPhoto{
	Url:         "http://www.google.com",
	Description: "An example image",
}

var exampleECL ECL = ECL{
	Id:       "ding",
	Name:     "ding",
	Title:    "ding",
	EclType:  "ding",
	Publish:  false,
	Currency: "ding",
	Sections: []ECLSection{},
}

var baseLocation Location = Location{
	Id:                     String("ding"),
	Name:                   String("ding"),
	CustomerId:             String("ding"),
	Address:                String("ding"),
	Address2:               String("ding"),
	DisplayAddress:         String("ding"),
	City:                   String("ding"),
	State:                  String("ding"),
	Zip:                    String("ding"),
	CountryCode:            String("ding"),
	Phone:                  String("ding"),
	LocalPhone:             String("ding"),
	AlternatePhone:         String("ding"),
	FaxPhone:               String("ding"),
	MobilePhone:            String("ding"),
	TollFreePhone:          String("ding"),
	TtyPhone:               String("ding"),
	SpecialOffer:           String("ding"),
	SpecialOfferUrl:        String("ding"),
	WebsiteUrl:             String("ding"),
	DisplayWebsiteUrl:      String("ding"),
	ReservationUrl:         String("ding"),
	Hours:                  String("ding"),
	AdditionalHoursText:    String("ding"),
	Description:            String("ding"),
	TwitterHandle:          String("ding"),
	FacebookPageUrl:        String("ding"),
	YearEstablished:        String("ding"),
	FolderId:               String("ding"),
	SuppressAddress:        Bool(false),
	IsPhoneTracked:         Bool(true),
	DisplayLat:             Float(1234.0),
	DisplayLng:             Float(1234.0),
	RoutableLat:            Float(1234.0),
	RoutableLng:            Float(1234.0),
	Keywords:               []string{"ding", "ding"},
	PaymentOptions:         []string{"ding", "ding"},
	VideoUrls:              []string{"ding", "ding"},
	Emails:                 []string{"ding", "ding"},
	Specialties:            []string{"ding", "ding"},
	Services:               []string{"ding", "ding"},
	Brands:                 []string{"ding", "ding"},
	Languages:              []string{"ding", "ding"},
	Logo:                   &examplePhoto,
	FacebookCoverPhoto:     &examplePhoto,
	FacebookProfilePicture: &examplePhoto,
	Photos:                 []LocationPhoto{examplePhoto, examplePhoto, examplePhoto},
	Lists:                  []ECL{exampleECL},
	Closed: &LocationClosed{
		IsClosed: "false",
	},
	CustomFields: map[string]interface{}{
		"1234": "ding",
	},
}

func TestCopy(t *testing.T) {
	secondLocation := *baseLocation.Copy()
	bV, sV := reflect.ValueOf(baseLocation), reflect.ValueOf(secondLocation)
	for i := 0; i < bV.NumField(); i++ {
		fieldName := reflect.TypeOf(baseLocation).Field(i).Name
		iBase, iSecond := bV.Field(i).Interface(), sV.Field(i).Interface()
		if !reflect.DeepEqual(iBase, iSecond) {
			t.Errorf("Copy not equal for field %v, expected %v got %v", fieldName, iBase, iSecond)
		}
	}
}

func TestDiffIdentical(t *testing.T) {
	secondLocation := *baseLocation.Copy()
	d, isDiff := baseLocation.Diff(secondLocation)
	if isDiff == true {
		t.Errorf("Expected diff to be false was true, diff result %v", d)
	} else if d != nil {
		t.Errorf("Expected an empty diff location, but got %v", d)
	}
}

type stringTest struct {
	baseValue          *string
	newValue           *string
	isDiff             bool
	expectedFieldValue *string
}

var stringTests = []stringTest{
	stringTest{String("ding"), String("ding"), false, nil},
	stringTest{String("ding"), String("dong"), true, String("dong")},
	stringTest{nil, String("dong"), true, String("dong")},
	stringTest{nil, String(""), true, String("")},
	stringTest{String(""), nil, false, nil},
	stringTest{String(""), String("dong"), true, String("dong")},
	{nil, nil, false, nil},
}

func formatStringPtr(s *string) string {
	if s == nil {
		return "nil"
	} else if *s == "" {
		return "empty string"
	} else {
		return *s
	}
}

func (t stringTest) formatErrorBase(index int) string {
	bv := formatStringPtr(t.baseValue)
	nv := formatStringPtr(t.newValue)
	return fmt.Sprintf("Failure with example %v:\n\tbase: '%v'\n\tnew: '%v'", index, bv, nv)
}

func TestStringDiffs(t *testing.T) {
	a, b := *new(Location), *new(Location)

	for i, data := range stringTests {
		a.Name, b.Name = data.baseValue, data.newValue
		d, isDiff := a.Diff(b)
		if isDiff != data.isDiff {
			t.Errorf("%vExpected diff to be %v\nbut was %v\ndiff struct was %v\n", data.formatErrorBase(i), data.isDiff, isDiff, d)
		}
		if d == nil && data.expectedFieldValue == nil {
			continue
		} else if d == nil && data.expectedFieldValue != nil {
			t.Errorf("%v\ndelta was nil but expected %v\n", data.formatErrorBase(i), formatStringPtr(data.expectedFieldValue))
		} else if d != nil && data.expectedFieldValue == nil {
			t.Errorf("%v\ndiff was%v\n", data.formatErrorBase(i), d)
		} else if *d.Name != *data.expectedFieldValue {
			t.Errorf("%v\ndiff was%v\n", data.formatErrorBase(i), d)
		}
	}
}

type boolTest struct {
	baseValue          *bool
	newValue           *bool
	isDiff             bool
	expectedFieldValue *bool
}

var boolTests = []boolTest{
	{Bool(false), Bool(false), false, nil},
	{Bool(true), Bool(true), false, nil},
	{Bool(false), Bool(true), true, Bool(true)},
	{Bool(true), Bool(false), true, Bool(false)},
	{nil, Bool(false), true, Bool(false)},
	{nil, Bool(true), true, Bool(true)},
	{Bool(false), nil, false, nil},
	{Bool(true), nil, false, nil},
	{nil, nil, false, nil},
}

func formatBoolPtr(b *bool) string {
	if b == nil {
		return "nil"
	} else {
		return fmt.Sprintf("%v", *b)
	}
}

func (t boolTest) formatErrorBase(index int) string {
	bv := formatBoolPtr(t.baseValue)
	nv := formatBoolPtr(t.newValue)
	return fmt.Sprintf("Failure with example %v:\n\tbase: '%v'\n\tnew: '%v'", index, bv, nv)
}

func TestBoolDiffs(t *testing.T) {
	a, b := *new(Location), *new(Location)
	for i, data := range boolTests {
		a.SuppressAddress, b.SuppressAddress = data.baseValue, data.newValue
		d, isDiff := a.Diff(b)
		if isDiff != data.isDiff {
			t.Errorf("%vExpected diff to be %v\nbut was %v\ndiff struct was %v\n", data.formatErrorBase(i), data.isDiff, isDiff, d)
		}
		if d == nil && data.expectedFieldValue == nil {
			continue
		} else if d == nil && data.expectedFieldValue != nil {
			t.Errorf("%v\ndelta was nil but expected %v\n", data.formatErrorBase(i), formatBoolPtr(data.expectedFieldValue))
		} else if d != nil && data.expectedFieldValue == nil {
			t.Errorf("%v\ndiff was%v\n", data.formatErrorBase(i), d)
		} else if *d.SuppressAddress != *data.expectedFieldValue {
			t.Errorf("%v\ndiff was%v\n", data.formatErrorBase(i), d)
		}
	}
}

type stringArrayTest struct {
	baseValue          []string
	newValue           []string
	isDiff             bool
	expectedFieldValue []string
}

var stringArrayTests = []stringArrayTest{
	{[]string{"ding", "dong"}, []string{"ding", "dong"}, false, nil},
	{[]string{"ding", "dong"}, []string{"ding", "dong", "dang"}, true, []string{"ding", "dong", "dang"}},
	{[]string{"ding", "dong", "dang"}, []string{"ding", "dong"}, true, []string{"ding", "dong"}},
	{[]string{}, []string{}, false, nil},
	{[]string{}, []string{"ding"}, true, []string{"ding"}},
	{[]string{}, nil, false, nil},
	{nil, []string{}, true, []string{}},
	{nil, nil, false, nil},
	{[]string{"ding"}, []string{}, true, []string{}},
	{[]string{"ding"}, nil, false, nil},
}

func (t stringArrayTest) formatErrorBase(index int) string {
	return fmt.Sprintf("Failure with example %v:\n\tbase: '%v'\n\tnew: '%v'", index, t.baseValue, t.newValue)
}

func TestStringArrayDiffs(t *testing.T) {
	a, b := *new(Location), *new(Location)
	for i, data := range stringArrayTests {
		a.PaymentOptions, b.PaymentOptions = data.baseValue, data.newValue
		d, isDiff := a.Diff(b)
		if isDiff != data.isDiff {
			t.Errorf("%vExpected diff to be %v\nbut was %v\ndiff struct was %v\n", data.formatErrorBase(i), data.isDiff, isDiff, d)
		}
		if d == nil && data.expectedFieldValue == nil {
			continue
		} else if d == nil && data.expectedFieldValue != nil {
			t.Errorf("%v\ndelta was nil but expected %v\n", data.formatErrorBase(i), data.expectedFieldValue)
		} else if len(d.PaymentOptions) != len(data.expectedFieldValue) {
			t.Errorf("%v\ndiff was%v\n", data.formatErrorBase(i), d)
		} else {
			for i := 0; i < len(d.PaymentOptions); i++ {
				if d.PaymentOptions[i] != data.expectedFieldValue[i] {
					t.Errorf("%v\ndiff was%v\n", data.formatErrorBase(i), d)
				}
			}
		}
	}
}

type floatTest struct {
	baseValue          *float64
	newValue           *float64
	isDiff             bool
	expectedFieldValue *float64
}

var floatTests = []floatTest{
	{Float(1234.0), Float(1234.0), false, nil},
	{Float(1234.0), nil, false, nil},
	{Float(0), nil, false, nil},
	{nil, nil, false, nil},
	{Float(0), Float(0), false, nil},
	{Float(0), Float(9876.0), true, Float(9876.0)},
	{Float(1234.0), Float(9876.0), true, Float(9876.0)},
	{nil, Float(9876.0), true, Float(9876.0)},
	{nil, Float(0), true, Float(0)},
}

func formatFloatPtr(b *float64) string {
	if b == nil {
		return "nil"
	} else {
		return fmt.Sprintf("%v", *b)
	}
}

func (t floatTest) formatErrorBase(index int) string {
	return fmt.Sprintf("Failure with example %v:\n\tbase: '%v'\n\tnew: '%v'", index, formatFloatPtr(t.baseValue), formatFloatPtr(t.newValue))
}

func TestFloatDiffs(t *testing.T) {
	a, b := *new(Location), *new(Location)
	for i, data := range floatTests {
		a.DisplayLat, b.DisplayLat = data.baseValue, data.newValue
		d, isDiff := a.Diff(b)
		if isDiff != data.isDiff {
			t.Errorf("%vExpected diff to be %v\nbut was %v\ndiff struct was %v\n", data.formatErrorBase(i), data.isDiff, isDiff, d)
		}
		if d == nil && data.expectedFieldValue == nil {
			continue
		} else if d == nil && data.expectedFieldValue != nil {
			t.Errorf("%v\ndelta was nil but expected %v\n", data.formatErrorBase(i), formatFloatPtr(data.expectedFieldValue))
		} else if d != nil && data.expectedFieldValue == nil {
			t.Errorf("%v\ndiff was%v\n", data.formatErrorBase(i), d)
		} else if *d.DisplayLat != *data.expectedFieldValue {
			t.Errorf("%v\ndiff was%v\n", data.formatErrorBase(i), d)
		}
	}
}

type photoTest struct {
	baseValue          *LocationPhoto
	newValue           *LocationPhoto
	isDiff             bool
	expectedFieldValue *LocationPhoto
}

func formatPhoto(b *LocationPhoto) string {
	if b == nil {
		return "nil"
	} else {
		return fmt.Sprintf("%v", *b)
	}
}

func (t photoTest) formatErrorBase(index int) string {
	return fmt.Sprintf("Failure with example %v:\n\tbase: '%v'\n\tnew: '%v'", index, formatPhoto(t.baseValue), formatPhoto(t.newValue))
}

var photoTests = []photoTest{
	{&LocationPhoto{Url: "ding", Description: "dong"}, &LocationPhoto{Url: "ding", Description: "dong"}, false, nil},
	{&LocationPhoto{Url: "ding", Description: "dong"}, nil, false, nil},
	{nil, &LocationPhoto{Url: "ding", Description: "dong"}, true, &LocationPhoto{Url: "ding", Description: "dong"}},
	{&LocationPhoto{Url: "ding"}, &LocationPhoto{Url: "ding", Description: "dong"}, true, &LocationPhoto{Url: "ding", Description: "dong"}},
	{&LocationPhoto{Description: "dong"}, &LocationPhoto{Url: "ding", Description: "dong"}, true, &LocationPhoto{Url: "ding", Description: "dong"}},
}

func TestPhotoDiffs(t *testing.T) {
	a, b := *new(Location), *new(Location)
	for i, data := range photoTests {
		a.FacebookCoverPhoto, b.FacebookCoverPhoto = data.baseValue, data.newValue
		d, isDiff := a.Diff(b)
		if isDiff != data.isDiff {
			t.Errorf("%vExpected diff to be %v\nbut was %v\ndiff struct was %v\n", data.formatErrorBase(i), data.isDiff, isDiff, d)
		}
		if d == nil && data.expectedFieldValue == nil {
			continue
		} else if d == nil && data.expectedFieldValue != nil {
			t.Errorf("%v\ndelta was nil but expected %v\n", data.formatErrorBase(i), formatPhoto(data.expectedFieldValue))
		} else if d != nil && data.expectedFieldValue == nil {
			t.Errorf("%v\ndiff was%v\n", data.formatErrorBase(i), d)
		} else if *d.FacebookCoverPhoto != *data.expectedFieldValue {
			t.Errorf("%v\ndiff was%v\n", data.formatErrorBase(i), d)
		}
	}
}

type photoArrayTest struct {
	baseValue          []LocationPhoto
	newValue           []LocationPhoto
	isDiff             bool
	expectedFieldValue []LocationPhoto
}

var photoArrayTests = []photoArrayTest{
	{nil, []LocationPhoto{}, true, []LocationPhoto{}},
	{nil, nil, false, nil},
	{[]LocationPhoto{LocationPhoto{Url: "ding", Description: "dong"}}, []LocationPhoto{}, true, []LocationPhoto{}},
	{[]LocationPhoto{LocationPhoto{Url: "ding", Description: "dong"}}, nil, false, nil},
	{[]LocationPhoto{LocationPhoto{Url: "ding", Description: "dong"}}, []LocationPhoto{LocationPhoto{Url: "dong", Description: "ding"}}, true, []LocationPhoto{LocationPhoto{Url: "dong", Description: "ding"}}},
	{[]LocationPhoto{LocationPhoto{Url: "ding", Description: "dong"}}, []LocationPhoto{LocationPhoto{Url: "ding", Description: "dong"}, LocationPhoto{Url: "ding", Description: "dong"}}, true, []LocationPhoto{LocationPhoto{Url: "ding", Description: "dong"}, LocationPhoto{Url: "ding", Description: "dong"}}},
}

func (t photoArrayTest) formatErrorBase(index int) string {
	return fmt.Sprintf("Failure with example %v:\n\tbase: '%v'\n\tnew: '%v'", index, t.baseValue, t.newValue)
}

func TestPhotoArrayDiffs(t *testing.T) {
	a, b := *new(Location), *new(Location)
	for i, data := range photoArrayTests {
		a.Photos, b.Photos = data.baseValue, data.newValue
		d, isDiff := a.Diff(b)
		if isDiff != data.isDiff {
			t.Errorf("%vExpected diff to be %v\nbut was %v\ndiff struct was %v\n", data.formatErrorBase(i), data.isDiff, isDiff, d)
		}
		if d == nil && data.expectedFieldValue == nil {
			continue
		} else if d == nil && data.expectedFieldValue != nil {
			t.Errorf("%v\ndelta was nil but expected %v\n", data.formatErrorBase(i), data.expectedFieldValue)
		} else if len(d.Photos) != len(data.expectedFieldValue) {
			t.Errorf("%v\ndiff was%v\n", data.formatErrorBase(i), d)
		} else {
			for i := 0; i < len(d.Photos); i++ {
				if d.Photos[i] != data.expectedFieldValue[i] {
					t.Errorf("%v\ndiff was%v\n", data.formatErrorBase(i), d)
				}
			}
		}
	}
}

type eclTest struct {
	baseValue          []ECL
	newValue           []ECL
	isDiff             bool
	expectedFieldValue []ECL
}

var eclTests = []eclTest{
	{nil, []ECL{}, true, []ECL{}},
	{[]ECL{}, nil, false, nil},
	{nil, nil, false, nil},
	{
		[]ECL{
			ECL{
				Id:       "1",
				Name:     "ding",
				Title:    "dong",
				EclType:  "EVENTS",
				Currency: "USD",
				Sections: []ECLSection{},
			},
		},
		[]ECL{
			ECL{
				Id:       "1",
				Name:     "ding",
				Title:    "dong",
				EclType:  "EVENTS",
				Currency: "USD",
				Sections: []ECLSection{},
			},
		},
		false,
		nil,
	},
	{
		[]ECL{
			ECL{
				Id:       "1",
				Name:     "Dang",
				Title:    "dong",
				EclType:  "EVENTS",
				Currency: "USD",
				Sections: []ECLSection{},
			},
		},
		[]ECL{
			ECL{
				Id:       "1",
				Name:     "ding",
				Title:    "dong",
				EclType:  "EVENTS",
				Currency: "USD",
				Sections: []ECLSection{},
			},
		},
		true,
		[]ECL{
			ECL{
				Id:       "1",
				Name:     "ding",
				Title:    "dong",
				EclType:  "EVENTS",
				Currency: "USD",
				Sections: []ECLSection{},
			},
		},
	},
	{
		[]ECL{
			ECL{
				Id:       "1",
				Name:     "ding",
				Title:    "dong",
				EclType:  "EVENTS",
				Currency: "USD",
				Sections: []ECLSection{},
			},
			ECL{
				Id:       "1",
				Name:     "ding",
				Title:    "dong",
				EclType:  "EVENTS",
				Currency: "USD",
				Sections: []ECLSection{},
			},
		},
		[]ECL{
			ECL{
				Id:       "1",
				Name:     "ding",
				Title:    "dong",
				EclType:  "EVENTS",
				Currency: "USD",
				Sections: []ECLSection{},
			},
		},
		true,
		[]ECL{
			ECL{
				Id:       "1",
				Name:     "ding",
				Title:    "dong",
				EclType:  "EVENTS",
				Currency: "USD",
				Sections: []ECLSection{},
			},
		},
	},
	{
		[]ECL{
			ECL{
				Id:       "1",
				Name:     "ding",
				Title:    "dong",
				EclType:  "EVENTS",
				Currency: "USD",
				Sections: []ECLSection{},
			},
		},
		[]ECL{
			ECL{
				Id:       "1",
				Name:     "ding",
				Title:    "dong",
				EclType:  "EVENTS",
				Currency: "USD",
				Sections: []ECLSection{},
			},
			ECL{
				Id:       "1",
				Name:     "ding",
				Title:    "dong",
				EclType:  "EVENTS",
				Currency: "USD",
				Sections: []ECLSection{},
			},
		},
		true,
		[]ECL{
			ECL{
				Id:       "1",
				Name:     "ding",
				Title:    "dong",
				EclType:  "EVENTS",
				Currency: "USD",
				Sections: []ECLSection{},
			},
			ECL{
				Id:       "1",
				Name:     "ding",
				Title:    "dong",
				EclType:  "EVENTS",
				Currency: "USD",
				Sections: []ECLSection{},
			},
		},
	},
}

func (t eclTest) formatErrorBase(index int) string {
	return fmt.Sprintf("Failure with example %v:\n\tbase: '%v'\n\tnew: '%v'", index, t.baseValue, t.newValue)
}

func TestECLDiffs(t *testing.T) {
	a, b := *new(Location), *new(Location)
	for i, data := range eclTests {
		a.Lists, b.Lists = data.baseValue, data.newValue
		d, isDiff := a.Diff(b)
		if isDiff != data.isDiff {
			t.Errorf("%vExpected diff to be %v\nbut was %v\ndiff struct was %v\n", data.formatErrorBase(i), data.isDiff, isDiff, d)
		}
		if d == nil && data.expectedFieldValue == nil {
			continue
		} else if d == nil && data.expectedFieldValue != nil {
			t.Errorf("%v\ndelta was nil but expected %v\n", data.formatErrorBase(i), data.expectedFieldValue)
		} else if len(d.Lists) != len(data.expectedFieldValue) {
			t.Errorf("%v\ndiff was%v\n", data.formatErrorBase(i), d)
		} else {
			for i := 0; i < len(d.Lists); i++ {
				good := true
				dL, eL := d.Lists[i], data.expectedFieldValue[i]
				if dL.Id != eL.Id {
					good = false
				} else if dL.Name != eL.Name {
					good = false
				} else if dL.Title != eL.Title {
					good = false
				} else if dL.EclType != eL.EclType {
					good = false
				} else if dL.Currency != eL.Currency {
					good = false
				} else if len(dL.Sections) != len(eL.Sections) {
					good = false
				} else {
					for j := 0; j < len(dL.Sections); j++ {
						dS, eS := dL.Sections[j], eL.Sections[i]
						if dS.Id != eS.Id {
							good = false
						} else if dS.Name != eS.Name {
							good = false
						} else if dS.Description != eS.Description {
							good = false
						}
					}
				}

				if !good {
					t.Errorf("%v\ndiff was%v\n", data.formatErrorBase(i), d)
				}
			}
		}
	}
}

type customFieldsTest struct {
	baseValue          map[string]interface{}
	newValue           map[string]interface{}
	isDiff             bool
	expectedFieldValue map[string]interface{}
}

var baseCustomFields = map[string]interface{}{
	"62150": []LocationPhoto{
		LocationPhoto{
			ClickThroughURL: "https://locations.yext.com",
			Description:     "This is the caption",
			Url:             "http://a.mktgcdn.com/p-sandbox/gRcmaehu-FoJtL3Ld6vNjYHpbZxmPSYZ1cTEF_UU7eY/1247x885.png",
		},
	},
	"62151": LocationPhoto{
		ClickThroughURL: "https://locations.yext.com",
		Description:     "This is a caption on a single!",
		Url:             "http://a.mktgcdn.com/p-sandbox/bSZ_mKhfFYGih6-ry5mtbwB_JbKu930kFxHOaQRwZC4/1552x909.png",
	},
	"62152": []string{
		"this",
		"is",
		"a",
		"textlist",
	},
	"62153": "This is a\r\nmulti\r\nline\r\ntext",
	"62154": "This is a single line text.",
	// Hours CustomField Type not really working in the model right now
	// "62144": {
	//   "holidayTimes": [
	//     {
	//       "date": "2015-12-14",
	//       "time": "9:00"
	//     },
	//     {
	//       "date": "2015-12-15",
	//       "time": "0:-1"
	//     }
	//   ],
	//   "dailyTimes": "2:18:00,3:18:00,4:18:00,5:18:00,6:18:00"
	// },
	"62155": "https://locations.yext.com/this-is-a-url",
	"62145": "12/14/2015",
	"62156": "true",
	// Hours CustomField Type not really working in the model right now
	// "62146": {
	//   "additionalHoursText": "We have wacky hours!",
	//   "hours": "2:9:00:18:00,3:19:00:22:00,3:9:00:18:00,4:0:00:0:00,5:9:00:18:00,6:9:00:18:00",
	//   "holidayHours": [
	//     {
	//       "date": "2015-12-12",
	//       "hours": ""
	//     },
	//     {
	//       "date": "2015-12-13",
	//       "hours": "10:00:13:00"
	//     },
	//     {
	//       "date": "2015-12-14",
	//       "hours": "0:00:0:00"
	//     },
	//     {
	//       "date": "2015-12-15",
	//       "hours": "10:00:17:00"
	//     }
	//   ]
	// },
	"62157": map[string]interface{}{"url": "http://www.youtube.com/watch?v=sYMYktsKmSk"},
	"62147": "10",
	"62148": []string{
		"27348",
		"27349",
	},
}

func copyCustomFields(cf map[string]interface{}) map[string]interface{} {
	n := map[string]interface{}{}
	for key, value := range cf {
		n[key] = value
	}
	return n
}

func appendJunkToCustomFields(cf map[string]interface{}) map[string]interface{} {
	n := copyCustomFields(cf)
	n["guy"] = "random junk"
	return n
}

func deleteKeyFromCustomField(cf map[string]interface{}) map[string]interface{} {
	n := copyCustomFields(cf)
	delete(n, "62148")
	return n
}

func modifyCF(cf map[string]interface{}) map[string]interface{} {
	n := copyCustomFields(cf)
	n["62153"] = "This is a\r\nMODIFIED multi\r\nline\r\ntext"
	return n
}

func zeroCFKEy(cf map[string]interface{}, key string) map[string]interface{} {
	n := copyCustomFields(cf)

	if value, ok := n[key]; ok {
		n[key] = reflect.Zero(reflect.TypeOf(value)).Interface()
	}

	return n
}

var (
	copyOfBase = copyCustomFields(baseCustomFields)
	appendedCF = appendJunkToCustomFields(baseCustomFields)
	trimmedCF  = deleteKeyFromCustomField(baseCustomFields)
	modifiedCF = modifyCF(baseCustomFields)
)

var customFieldsTests = []customFieldsTest{
	{nil, nil, false, nil},
	{map[string]interface{}{}, nil, false, nil},
	{map[string]interface{}{}, map[string]interface{}{}, false, nil},
	{nil, map[string]interface{}{}, true, map[string]interface{}{}},
	{baseCustomFields, copyOfBase, false, nil},
	{baseCustomFields, appendedCF, true, map[string]interface{}{"guy": "random junk"}},
	{baseCustomFields, trimmedCF, false, nil},
	{baseCustomFields, modifiedCF, true, map[string]interface{}{"62153": "This is a\r\nMODIFIED multi\r\nline\r\ntext"}},
}

func addZeroTests() {
	for key, val := range baseCustomFields {
		z := zeroCFKEy(baseCustomFields, key)
		zeroForKey := reflect.Zero(reflect.TypeOf(val)).Interface()
		test := customFieldsTest{baseCustomFields, z, true, map[string]interface{}{key: zeroForKey}}
		customFieldsTests = append(customFieldsTests, test)
	}
}

func (t customFieldsTest) formatErrorBase(index int) string {
	return fmt.Sprintf("Failure with example %v:\n\tbase: '%v'\n\tnew: '%v'", index, t.baseValue, t.newValue)
}

func TestCustomFieldsDiff(t *testing.T) {
	addZeroTests()
	a, b := *new(Location), *new(Location)
	for i, data := range customFieldsTests {
		a.CustomFields, b.CustomFields = data.baseValue, data.newValue
		d, isDiff := a.Diff(b)
		if isDiff != data.isDiff {
			t.Errorf("%vExpected diff to be %v\nbut was %v\ndiff struct was %v\n", data.formatErrorBase(i), data.isDiff, isDiff, d)
		}
		if d == nil && data.expectedFieldValue == nil {
			continue
		} else if d == nil && data.expectedFieldValue != nil {
			t.Errorf("%v\ndelta was nil but expected %v\n", data.formatErrorBase(i), data.expectedFieldValue)
		} else if !reflect.DeepEqual(data.expectedFieldValue, d.CustomFields) {
			t.Errorf("%v\ndiff was%v\n", data.formatErrorBase(i), d)
		}
	}
}

type closedTest struct {
	baseValue          *LocationClosed
	newValue           *LocationClosed
	isDiff             bool
	expectedFieldValue *LocationClosed
}

var closedTests = []closedTest{
	{nil, nil, false, nil},
	{&LocationClosed{}, nil, false, nil},
	{&LocationClosed{}, &LocationClosed{}, false, nil},
	{
		nil,
		&LocationClosed{
			IsClosed:   "true",
			ClosedDate: "1/1/2001",
		},
		true,
		&LocationClosed{
			IsClosed:   "true",
			ClosedDate: "1/1/2001",
		},
	},
	{
		&LocationClosed{
			IsClosed:   "true",
			ClosedDate: "1/1/2001",
		},
		&LocationClosed{
			IsClosed:   "true",
			ClosedDate: "1/1/2001",
		},
		false,
		nil,
	},
	{
		&LocationClosed{
			IsClosed:   "true",
			ClosedDate: "1/1/2001",
		},
		nil,
		false,
		nil,
	},
	{
		&LocationClosed{
			IsClosed:   "true",
			ClosedDate: "1/1/2001",
		},
		&LocationClosed{
			IsClosed:   "false",
			ClosedDate: "1/1/2001",
		},
		true,
		&LocationClosed{
			IsClosed:   "false",
			ClosedDate: "1/1/2001",
		},
	},
	{
		&LocationClosed{
			IsClosed:   "true",
			ClosedDate: "1/1/2001",
		},
		&LocationClosed{
			IsClosed:   "true",
			ClosedDate: "1/1/2002",
		},
		true,
		&LocationClosed{
			IsClosed:   "true",
			ClosedDate: "1/1/2002",
		},
	},
}

func formatClosed(b *LocationClosed) string {
	if b == nil {
		return "nil"
	} else {
		return fmt.Sprintf("%v", *b)
	}
}

func (t closedTest) formatErrorBase(index int) string {
	return fmt.Sprintf("Failure with example %v:\n\tbase: '%v'\n\tnew: '%v'", index, formatClosed(t.baseValue), formatClosed(t.newValue))
}

func TestClosedDiffs(t *testing.T) {
	a, b := *new(Location), *new(Location)
	for i, data := range closedTests {
		a.Closed, b.Closed = data.baseValue, data.newValue
		d, isDiff := a.Diff(b)
		if isDiff != data.isDiff {
			t.Errorf("%vExpected diff to be %v\nbut was %v\ndiff struct was %v\n", data.formatErrorBase(i), data.isDiff, isDiff, d)
		}
		if d == nil && data.expectedFieldValue == nil {
			continue
		} else if d == nil && data.expectedFieldValue != nil {
			t.Errorf("%v\ndelta was nil but expected %v\n", data.formatErrorBase(i), formatClosed(data.expectedFieldValue))
		} else if d != nil && data.expectedFieldValue == nil {
			t.Errorf("%v\ndelta was not nil but expected nil\n diff:%v\n", data.formatErrorBase(i), d)
		} else if d.Closed.IsClosed != data.expectedFieldValue.IsClosed {
			t.Errorf("%v\ndiff was%v\n", data.formatErrorBase(i), d)
		} else if d.Closed.ClosedDate != data.expectedFieldValue.ClosedDate {
			t.Errorf("%v\ndiff was%v\n", data.formatErrorBase(i), d)
		}
	}
}
