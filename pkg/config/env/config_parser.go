package env

import (
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/pechorka/gostdlib/pkg/errs"
)

/*
Config should look like this:
type Config struct {
	Db DbConfig `env:"DB"`
	Kafka KafkaConfig `env:"KAFKA"`
	Timeout CustomDuration `env:"TIMEOUT" default:"10d"`
}

type DbConfig struct {
	Host string `env:"HOST" required:"true"`
	Port int    `env:"PORT" default:"5432"`
}

type KafkaConfig struct {
	Brokers []string `env:"BROKERS"` // default sep is comma
}

type CustomDuration time.Duration

func (d *CustomDuration) UnmarshalEnv(v string) error {
	if strings.HasSuffix(v, "d") {
		v = strings.TrimSuffix(v, "d")
		days, err := strconv.Atoi(v)
		if err != nil {
			return errs.Wrap(err, "failed to parse duration")
		}
		*d = CustomDuration(time.Duration(days) * 24 * time.Hour)
		return nil
	}
	// default duration parsing
	duration, err := time.ParseDuration(v)
	if err != nil {
		return errs.Wrap(err, "failed to parse duration")
	}
	*d = CustomDuration(duration)
	return nil
}

Expected env file:
DB_HOST=localhost
DB_PORT=5432

KAFKA_BROKERS=localhost:9092,localhost:9093

TIMEOUT=10d
*/

func ParseConfig(cfg any) error {
	// Get the reflect value and type of the config struct
	v := reflect.ValueOf(cfg)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return errs.New("config must be a non-nil pointer")
	}

	v = v.Elem()
	t := v.Type()

	if t.Kind() != reflect.Struct {
		return errs.New("config must be a struct")
	}

	// Iterate through all fields in the struct
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		if err := parseField(fieldValue, field); err != nil {
			return errs.Wrapf(err, "failed to parse field %s", field.Name)
		}
	}

	return nil
}

func parseField(v reflect.Value, field reflect.StructField) error {
	prefix := field.Tag.Get("env")
	if prefix == "" {
		prefix = field.Name
	}

	t := v.Type()

	// Handle struct fields recursively
	if t.Kind() == reflect.Struct {
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			fieldValue := v.Field(i)

			envKey := prefix + "_" + field.Tag.Get("env")
			if field.Tag.Get("env") == "" {
				envKey = prefix + "_" + field.Name
			}

			envValue, err := getEnvValue(envKey, field.Tag)
			if err != nil {
				return errs.Wrapf(err, "failed to get env value for field %s", field.Name)
			}

			if err := setFieldValue(fieldValue, envValue); err != nil {
				return errs.Wrapf(err, "failed to set field %s", field.Name)
			}
		}
		return nil
	}

	envValue, err := getEnvValue(prefix, field.Tag)
	if err != nil {
		return errs.Wrapf(err, "failed to get env value for field %s", field.Name)
	}

	// Handle non-struct fields
	return setFieldValue(v, envValue)
}

func getEnvValue(key string, tag reflect.StructTag) (string, error) {
	envValue := os.Getenv(strings.ToUpper(key))
	if envValue == "" {
		if defaultVal := tag.Get("default"); defaultVal != "" {
			envValue = defaultVal
		}
	}
	if envValue == "" && tag.Get("required") == "true" {
		return "", errs.Newf("required environment variable %s is not set", key)
	}
	return envValue, nil
}

type EnvUnmarshaler interface {
	UnmarshalEnv(string) error
}

func setFieldValue(v reflect.Value, envValue string) error {
	// Check if the field implements UnmarshalEnv
	if unmarshaler, ok := v.Addr().Interface().(EnvUnmarshaler); ok {
		return unmarshaler.UnmarshalEnv(envValue)
	}

	// Handle different types
	switch v.Kind() {
	case reflect.String:
		v.SetString(envValue)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val, err := strconv.ParseInt(envValue, 10, 64)
		if err != nil {
			return errs.Wrapf(err, "failed to parse int value")
		}
		v.SetInt(val)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val, err := strconv.ParseUint(envValue, 10, 64)
		if err != nil {
			return errs.Wrapf(err, "failed to parse uint value")
		}
		v.SetUint(val)

	case reflect.Float32, reflect.Float64:
		val, err := strconv.ParseFloat(envValue, 64)
		if err != nil {
			return errs.Wrapf(err, "failed to parse float value")
		}
		v.SetFloat(val)

	case reflect.Bool:
		val, err := strconv.ParseBool(envValue)
		if err != nil {
			return errs.Wrapf(err, "failed to parse bool value")
		}
		v.SetBool(val)

	case reflect.Slice:
		parts := strings.Split(envValue, ",")
		slice := reflect.MakeSlice(v.Type(), len(parts), len(parts))
		for i, part := range parts {
			if err := setFieldValue(slice.Index(i), part); err != nil {
				return err
			}
		}
		v.Set(slice)

	default:
		return errs.Newf("unsupported type %s", v.Kind())
	}

	return nil
}
