package utils

import (
	"fmt"
	"reflect"
	"spmb/back-end/models"
	"strings"
	"time"
)

func FormatSettings(settings []models.Settings) []models.SettingResponse {
	var formatted []models.SettingResponse

	for _, s := range settings {
		formatted = append(formatted, models.SettingResponse{
			SettingID:      s.SettingID,
			AppName:        s.AppName,
			Modal:          s.Modal,
			ModalHeader:    s.ModalHeader,
			ModalBody:      s.ModalBody,
			SchedulePic:    s.SchedulePic,
			RequirementPic: s.RequirementPic,
			HowToRegPic:    s.HowToRegPic,
			HowToRegVid:    s.HowToRegVid,
			IsActive:       s.IsActive,
			CreatedBy:      s.CreatedBy.String(),
			CreatedAt:      s.UpdatedAt.Format("2006-01-02 15:04:05"),
			UpdatedBy:      s.UpdatedBy.String(),
			UpdatedAt:      s.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return formatted
}
func FormatFAQ(faq []models.FAQ) []models.FAQResponse {
	var formatted []models.FAQResponse

	for _, s := range faq {
		formatted = append(formatted, models.FAQResponse{
			FaqID:     s.FaqID,
			Question:  s.Question,
			Answer:    s.Answer,
			IsActive:  s.IsActive,
			CreatedBy: s.CreatedBy,
			CreatedAt: s.UpdatedAt.Format("2006-01-02 15:04:05"),
			UpdatedBy: s.UpdatedBy,
			UpdatedAt: s.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return formatted
}

func FormatStructSlice(slice interface{}, layout string) ([]map[string]interface{}, error) {
	// Ensure that the input is a slice.
	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Slice {
		return nil, fmt.Errorf("FormatStructSlice expects a slice, got %s", v.Kind())
	}

	result := make([]map[string]interface{}, 0, v.Len())

	// Iterate over the slice elements.
	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i)
		// If element is a pointer, dereference it.
		if elem.Kind() == reflect.Ptr {
			elem = elem.Elem()
		}

		// Ensure the element is a struct.
		if elem.Kind() != reflect.Struct {
			return nil, fmt.Errorf("FormatStructSlice expects a slice of structs, element %d is %s", i, elem.Kind())
		}

		formatted := make(map[string]interface{})
		t := elem.Type()

		for j := 0; j < elem.NumField(); j++ {
			field := t.Field(j)
			value := elem.Field(j)

			// Get the JSON tag; fallback to the field name.
			jsonTag := field.Tag.Get("json")
			if jsonTag == "" || jsonTag == "-" {
				jsonTag = field.Name
			} else {
				// In case the tag has options (e.g. "question,omitempty"), take only the key.
				jsonTag = strings.Split(jsonTag, ",")[0]
			}

			// If the field is time.Time, format it.
			if value.Type() == reflect.TypeOf(time.Time{}) {
				tVal := value.Interface().(time.Time)
				formatted[jsonTag] = tVal.Format(layout)
			} else {
				formatted[jsonTag] = value.Interface()
			}
		}
		result = append(result, formatted)
	}
	return result, nil
}
