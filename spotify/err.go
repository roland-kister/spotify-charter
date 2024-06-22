package spotify

import (
	"errors"
	"fmt"
	"io"
)

type AuthErrResp struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func authErrRespToErr(body *io.ReadCloser) error {
	resp, err := decodeResp[AuthErrResp](body)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("auth err (%s):\n%s", resp.Error, resp.ErrorDescription)

	return errors.New(msg)
}

type RegErrResp struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func regErrRespToErr(body *io.ReadCloser) error {
	resp, err := decodeResp[RegErrResp](body)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("api err [%d]: %s", resp.Status, resp.Message)

	return errors.New(msg)
}
