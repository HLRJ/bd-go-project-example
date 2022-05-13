package cotroller

import (
	"strconv"

	"github.com/Moonlight-Zhao/go-project-example/service"
)

type PageData struct {
	Code int64       `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func QueryPageInfo(topicIdStr string) *PageData {
	topicId, err := strconv.ParseInt(topicIdStr, 10, 64)
	if err != nil {
		return &PageData{
			Code: -1,
			Msg:  err.Error(),
		}
	}
	pageInfo, err := service.QueryPageInfo(topicId)
	if err != nil {
		return &PageData{ //返回错误信息
			Code: -1,
			Msg:  err.Error(),
		}
	}
	return &PageData{ // 返回成功和正确的信息
		Code: 0,
		Msg:  "success",
		Data: pageInfo,
	}

}
