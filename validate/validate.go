package validate

import (
	"reflect"
	"strings"

	enlocale "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	lib "github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/translations/en"

	"github.com/jkaveri/goservice/check"
)

var validatorInstance = NewValidator()

type Validator struct {
	*lib.Validate
	translator ut.Translator
}

func (v *Validator) Translate(err error) string {
	var (
		errs          = err.(lib.ValidationErrors)
		fieldMessages = errs.Translate(v.translator)
		sb            strings.Builder
		i             int
	)

	for d, msg := range fieldMessages {
		if i > 0 {
			_, _ = sb.WriteString("; ")
		}

		_, _ = sb.WriteString(d)
		_, _ = sb.WriteString(": ")
		_, _ = sb.WriteString(msg)
		i++
	}

	return sb.String()
}

func NewValidator() *Validator {
	v := lib.New(
		lib.WithRequiredStructEnabled(),
	)

	// register validators
	check.PanicIfError(
		v.RegisterValidation("shortuuid", validateShortUUID),
	)

	// register translations
	translator := registerTranslations(v)

	return &Validator{
		Validate:   v,
		translator: translator,
	}
}

func validateShortUUID(fl lib.FieldLevel) bool {
	field := fl.Field()
	k := field.Kind()

	switch k {
	case reflect.String:
		n := field.Interface().(string)

		return len(n) == 22 || len(n) == 0
	default:
		panic("shortuuid validation can only be applied to string")
	}
}

func registerTranslations(v *lib.Validate) ut.Translator {
	universalTranslator := ut.New(enlocale.New())
	enTranslator, _ := universalTranslator.GetTranslator("en")

	// register default ranslations
	_ = en.RegisterDefaultTranslations(v, enTranslator)

	check.PanicIfError(v.RegisterTranslation(
		"shortuuid",
		enTranslator,
		func(ut ut.Translator) error {
			return ut.Add("shortuuid", "{0} is not valid short uuid", false)
		},
		func(ut ut.Translator, fe lib.FieldError) string {
			str, _ := ut.T(fe.Tag(), fe.Field())
			return str
		},
	))

	check.PanicIfError(v.RegisterTranslation(
		"required_without",
		enTranslator,
		func(ut ut.Translator) error {
			return ut.Add(
				"required_without",
				"{0} and {1} are both empty",
				true,
			)
		},
		func(ut ut.Translator, fe lib.FieldError) string {
			t, _ := ut.T(fe.Tag(), fe.Field(), fe.Param())

			return t
		},
	))

	return enTranslator
}
