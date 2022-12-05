package gormupdatemap

import (
	"fmt"
	"reflect"
	"strings"
	"time"
	"unicode"

	log "github.com/sirupsen/logrus"
)

// CreateUpdateMap creates a map to use in the gorm.Updates method. Takes in a
// struct which defines the fields to add to the map. Adds all non-nil fields to the update map.
// Using this map even default values are updated, like false booleans and 0
// values for ints and floats. The isAdmin parameter determines whether the update map will contain
// fields from the given updateRequestStruct that have the 'admin_only' tag. The returned *string
// value is a validation error string, which returns any validation errors, like for instance trying
// to update admin fields as a non-admin
func CreateUpdateMap(updateRequestStruct interface{}, isAdmin bool) (map[string]interface{}, *string) {
	rType := reflect.TypeOf(updateRequestStruct)
	rVal := reflect.ValueOf(updateRequestStruct)
	n := rType.NumField()

	// Keep track of the fields we need to update
	fieldsToSet := make(map[string]interface{})

	for i := 0; i < n; i++ {
		fType := rType.Field(i)
		fVal := rVal.Field(i)

		var val reflect.Value
		if fVal.Kind() == reflect.Ptr {
			// Ignore fields that are nil
			if fVal.IsNil() {
				log.Debugf("Not adding field '%s' to update map as it is nil.", fType.Name)
				continue
			}

			// Grab value of non-nil field
			val = fVal.Elem()
		} else {
			// Ignore non-pointer fields
			log.Debugf("Not adding field '%s' to update map as it is not a pointer.", fType.Name)
			continue
		}

		// Get json field tag value
		JSONTag := fType.Tag.Get("json")

		// Skip unexported fields
		if fType.PkgPath != "" {
			log.Debugf("Not adding field '%s' to update map as it is not exported.", fType.Name)
			continue
		}

		// Skip json omitted fields
		if JSONTag == "-" {
			continue
		}

		// If no tag is set, use the field name
		if JSONTag == "" {
			JSONTag = fType.Name
		}

		// Check for multiple values in field
		JSONTagValues := strings.Split(JSONTag, ",")
		if len(JSONTagValues) > 0 {
			// Grab the first value
			JSONTag = JSONTagValues[0]
		}

		// Check if updating this field requires admin permissions
		adminTag := fType.Tag.Get("admin_only")
		if adminTag != "" && !isAdmin {
			// Return validation error
			return nil, strPtr(fmt.Sprintf("Only admins can update field '%s'.", JSONTag))
		}

		// If the tag contains an underscore followed by a number, remove the underscore
		// This is to comply with the GORM column naming strategy
		underscoreSegments := strings.Split(JSONTag, "_")
		JSONTag = ""
		for i, seg := range underscoreSegments {
			// If segment is followed by another segment
			if i+1 < len(underscoreSegments) {
				// If first character in next segment is a digit, add no underscore
				if unicode.IsDigit(rune(underscoreSegments[i+1][0])) {
					JSONTag += seg
				} else {
					// Add an underscore
					JSONTag += seg + "_"
				}
			} else {
				// Add last segment
				JSONTag += seg
			}
		}

		// Make all characters lowercase
		JSONTag = strings.ToLower(JSONTag)

		switch val.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			fieldsToSet[JSONTag] = val.Int()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			fieldsToSet[JSONTag] = val.Uint()
		case reflect.Float32, reflect.Float64:
			fieldsToSet[JSONTag] = val.Float()
		case reflect.String:
			fieldsToSet[JSONTag] = val.String()
		case reflect.Bool:
			fieldsToSet[JSONTag] = val.Bool()
		case reflect.Struct:
			if val.Type().String() == "time.Time" {
				if v, ok := val.Interface().(time.Time); ok {
					fieldsToSet[JSONTag] = v
				} else {
					log.Debugf("Not adding field '%s' to update map as it is not a valid time.Time value.", fType.Name)
				}
			}
		}
	}

	return fieldsToSet, nil
}
