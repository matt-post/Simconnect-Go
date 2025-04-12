package simconnect

import (
	"errors"
	"fmt"
	"time"
	"strings"

	simconnect_data "github.com/JRascagneres/Simconnect-Go/simconnect-data"
)

func derefDataType(fieldType string) (uint32, error) {
	switch fieldType {
	case "int32", "bool":
		return simconnect_data.DATATYPE_INT32, nil
	case "int64":
		return simconnect_data.DATATYPE_INT64, nil
	case "float32":
		return simconnect_data.DATATYPE_FLOAT32, nil
	case "float64":
		return simconnect_data.DATATYPE_FLOAT64, nil
	case "[8]byte":
		return simconnect_data.DATATYPE_STRING8, nil
	case "[32]byte":
		return simconnect_data.DATATYPE_STRING32, nil
	case "[64]byte":
		return simconnect_data.DATATYPE_STRING64, nil
	case "[128]byte":
		return simconnect_data.DATATYPE_STRING128, nil
	case "[256]byte":
		return simconnect_data.DATATYPE_STRING256, nil
	case "[260]byte":
		return simconnect_data.DATATYPE_STRING260, nil
	default:
		// Handle dynamic [N]byte
		var size int
		_, err := fmt.Sscanf(fieldType, "[%d]byte", &size)
		if err == nil {
			switch size {
			case 8:
				return simconnect_data.DATATYPE_STRING8, nil
			case 32:
				return simconnect_data.DATATYPE_STRING32, nil
			case 64:
				return simconnect_data.DATATYPE_STRING64, nil
			case 128:
				return simconnect_data.DATATYPE_STRING128, nil
			case 256:
				return simconnect_data.DATATYPE_STRING256, nil
			case 260:
				return simconnect_data.DATATYPE_STRING260, nil
			default:
				// Best effort fallback: treat as STRING256 if not a known size
				return simconnect_data.DATATYPE_STRING256, nil
			}
		}
	}

	return 0, fmt.Errorf("DATATYPE not implemented: %s", fieldType)
}

func retryFunc(maxRetryCount int, waitDuration time.Duration, dataFunc func() (bool, error)) error {
	numAttempts := 1

	for {
		shouldRetry, _ := dataFunc()
		if !shouldRetry {
			return nil
		}

		numAttempts++

		if numAttempts >= maxRetryCount {
			return errors.New("timeout exceeded err")
		}

		time.Sleep(waitDuration)
	}
}
