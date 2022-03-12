package valix

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	tagOpenApi             = "oas"
	tagOpenApiDescription  = "description"
	tagOpenApiDescription2 = "desc"
	tagOpenApiTitle        = "title"
	tagOpenApiFormat       = "format"
	tagOpenApiFormat2      = "fmt"
	tagOpenApiExample      = "example"
	tagOpenApiExample2     = "eg"
	tagOpenApiDeprecated   = "deprecated"
)

const (
	msgOasPrefix            = "tag " + tagOpenApi + " - "
	msgOasUnknownTokenInTag = msgOasPrefix + "unknown token '%s'"
	msgOasUnexpectedColon   = msgOasPrefix + "unexpected ':' colon after token '%s'"
	msgOasExpectedString    = msgOasPrefix + "expected string (enclosed with \"\" or '') after token '%s'"
)

type OasInfo struct {
	Description string
	Title       string
	Format      string
	Example     string
	Deprecated  bool
}

func (pv *PropertyValidator) processOasTag(fld reflect.StructField) error {
	if tag, ok := fld.Tag.Lookup(tagOpenApi); ok {
		tagItems, err := parseCommas(tag)
		if err != nil {
			return err
		}
		return pv.processOasTagItems(tagItems)
	}
	return nil
}

func (pv *PropertyValidator) processOasTagItems(tagItems []string) error {
	for _, ti := range tagItems {
		tagItem := strings.Trim(ti, " ")
		if tagItem != "" {
			if err := pv.addOasTagItem(tagItem); err != nil {
				return err
			}
		}
	}
	return nil
}

func (pv *PropertyValidator) addOasTagItem(tagItem string) (result error) {
	if pv.OasInfo == nil {
		pv.OasInfo = &OasInfo{}
	}
	result = nil
	tagToken := tagItem
	tagValue := ""
	tagValueIsStr := false
	hasColon := false
	if cAt := firstValidColonAt(tagItem); cAt != -1 {
		hasColon = true
		tagToken = strings.Trim(tagItem[0:cAt], " ")
		tagValue = strings.Trim(tagItem[cAt+1:], " ")
		tagValueIsStr = (strings.HasPrefix(tagValue, "\"") && strings.HasSuffix(tagValue, "\"")) ||
			(strings.HasPrefix(tagValue, "'") && strings.HasSuffix(tagValue, "'"))
		if tagValueIsStr {
			tagValue = tagValue[1 : len(tagValue)-1]
		}
	}
	colonErr := false
	strErr := false
	switch tagToken {
	case tagOpenApiDescription, tagOpenApiDescription2:
		strErr = !tagValueIsStr
		if !strErr {
			pv.OasInfo.Description = tagValue
		}
		break
	case tagOpenApiTitle:
		strErr = !tagValueIsStr
		if !strErr {
			pv.OasInfo.Title = tagValue
		}
		break
	case tagOpenApiFormat, tagOpenApiFormat2:
		strErr = !tagValueIsStr
		if !strErr {
			pv.OasInfo.Format = tagValue
		}
		break
	case tagOpenApiExample, tagOpenApiExample2:
		strErr = !tagValueIsStr
		if !strErr {
			pv.OasInfo.Example = tagValue
		}
		break
	case tagOpenApiDeprecated:
		colonErr = hasColon
		pv.OasInfo.Deprecated = true
		break
	default:
		result = fmt.Errorf(msgOasUnknownTokenInTag, tagToken)
	}
	if strErr {
		result = fmt.Errorf(msgOasExpectedString, tagToken)
	} else if colonErr {
		result = fmt.Errorf(msgOasUnexpectedColon, tagToken)
	}
	return
}
