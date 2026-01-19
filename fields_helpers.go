package logx

import "time"

// Err creates a logging field for an error.
//
// The error is converted to its string representation.
func Err(err error) Field {
	return Field{
		Key:   "error",
		Value: err.Error(),
	}
}

// String creates a string logging field.
func String(key, value string) Field {
	return Field{
		Key:   key,
		Value: value,
	}
}

// Int creates an integer logging field.
func Int(key string, value int) Field {
	return Field{
		Key:   key,
		Value: value,
	}
}

// Int64 creates an int64 logging field.
func Int64(key string, value int64) Field {
	return Field{
		Key:   key,
		Value: value,
	}
}

// Duration creates a duration logging field.
func Duration(key string, value time.Duration) Field {
	return Field{
		Key:   key,
		Value: value,
	}
}
