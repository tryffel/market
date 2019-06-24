package util

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/tryffel/market"
	"github.com/tryffel/market/modules/request"
	"github.com/tryffel/market/modules/response"
)

func GetServerInfo(req request.Request, resp response.Response) {
	body := map[string]interface{}{}
	body["name"] = market.Name
	body["version"] = market.Version

	data, err := json.Marshal(body)
	if err != nil {
		logrus.Error(err)
	}
	err = resp.Write([]byte(data), response.StatusOk)
}
