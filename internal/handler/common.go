package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"maps"
	"net/http"
	"strings"
	"time"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
	"github.com/ahleongzc/leetcode-live-backend/internal/util"
	"github.com/coder/websocket"

	"github.com/rs/zerolog/log"
)

func WriteJSONWebsocket(ctx context.Context, conn *websocket.Conn, payload util.JSONPayload) error {
	var marshalledPayload []byte
	var err error

	if payload != nil {
		marshalledPayload, err = json.Marshal(payload)
		if err != nil {
			return err
		}
	}

	conn.Write(ctx, websocket.MessageText, marshalledPayload)
	return nil
}

func WriteJSONHTTP(w http.ResponseWriter, payload util.JSONPayload, statusCode int, headers http.Header) error {
	var marshalledPayload []byte
	var err error

	if payload != nil {
		marshalledPayload, err = json.Marshal(payload)
		if err != nil {
			return err
		}
	}

	maps.Copy(w.Header(), headers)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(marshalledPayload)

	return nil
}

func ReadJSONBytes(data []byte, dst any) error {
	if dst == nil {
		return fmt.Errorf("dst cannot be nil when calling readJSON for websockets: %w", common.ErrInternalServerError)
	}

	maxBytes := common.INCOMING_PAYLOAD_MAX_BYTES
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()

	err := decodeJSON(dec, maxBytes, dst)
	if err != nil {
		return err
	}

	return nil
}

func ReadJSONHTTPReq(w http.ResponseWriter, r *http.Request, dst any) error {
	if dst == nil {
		return fmt.Errorf("dst cannot be nil when calling readJSON for http: %w", common.ErrInternalServerError)
	}

	maxBytes := common.INCOMING_PAYLOAD_MAX_BYTES
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := decodeJSON(dec, maxBytes, dst)
	if err != nil {
		return err
	}

	return nil
}

func decodeJSON(decoder *json.Decoder, maxBytes int, dst any) error {
	err := decoder.Decode(dst)

	if err != nil {
		// Client: provided invalid JSON syntax
		var syntaxError *json.SyntaxError
		// Client: JSON that is provided cannot be marshaled into go types
		var unmarshalTypeError *json.UnmarshalTypeError
		// Backend: When dst argument is non-nil, which is a dev mistake instead of a client error
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d): %w", syntaxError.Offset, common.ErrBadRequest)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return fmt.Errorf("body contains badly-formed JSON: %w", common.ErrBadRequest)

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q: %w", unmarshalTypeError.Field, common.ErrBadRequest)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d): %w", unmarshalTypeError.Offset, common.ErrBadRequest)

		case errors.Is(err, io.EOF):
			return fmt.Errorf("body must not be empty: %w", common.ErrBadRequest)

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s: %w", fieldName, common.ErrBadRequest)

		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes: %w", maxBytes, common.ErrBadRequest)

		case errors.As(err, &invalidUnmarshalError):
			return fmt.Errorf("unable to unmarshal json: %w", common.ErrInternalServerError)

		default:
			return common.ErrInternalServerError
		}
	}

	err = decoder.Decode(&struct{}{})
	if err != io.EOF {
		return fmt.Errorf("body must only contain a single JSON value: %w", common.ErrBadRequest)
	}

	return nil
}

func HandleErrorResponeWebsocket(ctx context.Context, conn *websocket.Conn, err error) {
	errorMessage := cleanErrorMessage(err)

	switch {
	case errors.Is(err, common.ErrNormalClientClosure):
		// Ignore normal client closure
	case errors.Is(err, common.ErrBadRequest):
		clientBadRequestErrorWebsocketResponse(ctx, conn, errors.New(errorMessage))
	case errors.Is(err, common.ErrInternalServerError):
		internalServerErrorWebsocketResponse(ctx, conn, errors.New(errorMessage))
	default:
		internalServerErrorWebsocketResponse(ctx, conn, errors.New(errorMessage))
	}
}

func HandleErrorResponseHTTP(w http.ResponseWriter, err error) {
	errorMessage := cleanErrorMessage(err)

	switch {
	case errors.Is(err, common.ErrBadRequest):
		clientBadRequestErrorHTTPResponse(w, errors.New(errorMessage))
	case errors.Is(err, common.ErrUnauthorized):
		clientUnauthorizedHTTPResponse(w)
	case errors.Is(err, common.ErrForbidden):
		clientForbiddenHTTPResponse(w)
	case errors.Is(err, common.ErrNotFound):
		dataNotFoundErrorHTTPResponse(w, errors.New(errorMessage))
	case errors.Is(err, common.ErrInternalServerError):
		internalServerErrorHTTPResponse(w, errors.New(errorMessage))
	default:
		internalServerErrorHTTPResponse(w, errors.New(errorMessage))
	}
}

func cleanErrorMessage(err error) string {
	errorMessage := err.Error()
	if idx := strings.LastIndex(errorMessage, ":"); idx != -1 {
		errorMessage = errorMessage[:idx]
	}

	return strings.TrimSpace(errorMessage)
}

func clientBadRequestErrorWebsocketResponse(ctx context.Context, conn *websocket.Conn, err error) {
	payload := util.NewJSONPayload()
	payload.Add("error", err.Error())

	WriteJSONWebsocket(ctx, conn, payload)
}

func internalServerErrorWebsocketResponse(ctx context.Context, conn *websocket.Conn, err error) {
	log.Error().Err(err).Msg("")

	payload := util.NewJSONPayload()
	payload.Add("error", "internal server error")

	WriteJSONWebsocket(ctx, conn, payload)
}

func dataNotFoundErrorHTTPResponse(w http.ResponseWriter, err error) {
	payload := util.NewJSONPayload()
	payload.Add("error", err.Error())

	WriteJSONHTTP(w, payload, http.StatusNotFound, nil)
}

func clientBadRequestErrorHTTPResponse(w http.ResponseWriter, err error) {
	payload := util.NewJSONPayload()
	payload.Add("error", err.Error())

	WriteJSONHTTP(w, payload, http.StatusBadRequest, nil)
}

func clientForbiddenHTTPResponse(w http.ResponseWriter) {
	payload := util.NewJSONPayload()
	payload.Add("error", "you are unauthorized to access this resource")

	WriteJSONHTTP(w, payload, http.StatusForbidden, nil)
}

func clientUnauthorizedHTTPResponse(w http.ResponseWriter) {
	payload := util.NewJSONPayload()
	payload.Add("error", "wrong credentials")

	WriteJSONHTTP(w, payload, http.StatusUnauthorized, nil)
}

func internalServerErrorHTTPResponse(w http.ResponseWriter, err error) {
	log.Error().Err(err).Msg("error")

	payload := util.NewJSONPayload()
	payload.Add("error", "internal server error")

	w.Header().Set("Connection", "close")
	WriteJSONHTTP(w, payload, http.StatusInternalServerError, nil)
}

func NewHTTPCookie(name, value string, expiryTimestampMS int64) *http.Cookie {
	var secure, httpOnly bool
	sameSite := http.SameSiteLaxMode

	if util.IsProdEnv() {
		secure = true
		httpOnly = true
		sameSite = http.SameSiteStrictMode
	}

	return &http.Cookie{
		Name:     name,
		Value:    value,
		Expires:  time.UnixMilli(expiryTimestampMS),
		Secure:   secure,
		HttpOnly: httpOnly,
		SameSite: sameSite,
	}
}
