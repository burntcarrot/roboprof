package shipper

import (
	"bytes"
	"context"
	"net/http"

	"github.com/burntcarrot/roboprof/pkg/utils"
)

func Ship(ctx context.Context, buf *bytes.Buffer) error {
	addr := ""
	req, err := http.NewRequest(http.MethodPost, addr, buf)
	if err != nil {
		return err
	}

	r := req.WithContext(ctx)

	err = utils.SendReq(r)
	if err != nil {
		return err
	}

	return nil
}
