package common

import (
	"context"
	"errors"
	"log"
	"runtime"

	"api.mooody.me/db"
	"api.mooody.me/models"
)

func getCallerFrame() runtime.Frame {
	// We need the frame at index skipFrames+2, since we never want runtime.Callers and getFrame
	targetFrameIndex := 1 + 2

	// Set size to targetFrameIndex+2 to ensure we have room for one more caller than we need
	programCounters := make([]uintptr, targetFrameIndex+2)
	n := runtime.Callers(0, programCounters)

	frame := runtime.Frame{Function: "unknown"}
	if n > 0 {
		frames := runtime.CallersFrames(programCounters[:n])
		for more, frameIndex := true, 0; more && frameIndex <= targetFrameIndex; frameIndex++ {
			var frameCandidate runtime.Frame
			frameCandidate, more = frames.Next()
			if frameIndex == targetFrameIndex {
				frame = frameCandidate
			}
		}
	}

	return frame
}

func LogClient(client *models.APIClient, format string, msgs ...interface{}) {
	caller := getCallerFrame().Function
	if client == nil {
		log.Printf("INFO:  [%s][%s]: "+format, caller, "UNKNOWN CLIENT", msgs)
	} else {
		log.Printf("INFO:  [%s][%s]: "+format, caller, client.GetName(), msgs)
	}
}

func LogClientWithError(client *models.APIClient, err error) {
	caller := getCallerFrame().Function
	if client == nil {
		log.Printf("ERROR: [%s][%s]: (%s)", caller, "UNKNOWN CLIENT", err.Error())
	} else {
		log.Printf("ERROR: [%s][%s]: (%s)", caller, client.GetName(), err.Error())
	}
}

func GetClientFromAuth(ctx context.Context, auth *models.Auth, requirePrivileged bool) (*models.APIClient, error) {
	if auth == nil {
		return nil, errors.New("invalid Auth")
	}

	client, err := db.GetClientByUUID(ctx, auth.ClientUuid)
	if err != nil {
		return nil, errors.New("invalid client id")
	}

	if !client.GetEnabled() {
		return nil, errors.New("client is not enabled")
	}

	if requirePrivileged {
		if !client.GetPrivileged() {
			return client, errors.New("client isn't privileged as required")
		}
	}

	return client, nil
}
