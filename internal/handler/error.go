package handler

import (
	"context"
	"github.com/rs/zerolog/log"
	"runtime/debug"
)

type skillError struct {
	cause error
	msg   string
}

type ErrorHandler struct {
	errorChanel chan skillError
}

func NewErrorHandler() *ErrorHandler {
	errors := make(chan skillError)
	return &ErrorHandler{errors}
}

func (h *ErrorHandler) HandleError(err error) {
	h.errorChanel <- skillError{cause: err}
}

func (h *ErrorHandler) HandleErrorWithMsg(err error, msg string) {
	h.errorChanel <- skillError{cause: err, msg: msg}
}

func (h *ErrorHandler) Do(ctx context.Context) {
	for {
		select {

		case <-ctx.Done():
			return

		case err, ok := <-h.errorChanel:
			if !ok {
				return
			}
			log.Error().Err(err.cause).Stack().Msgf("Error skill err: %v, msg: %s,\n stack: %s", err.cause, err.msg, string(debug.Stack()))
		}
	}
}
