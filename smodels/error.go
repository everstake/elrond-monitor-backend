package smodels

type (
	Err interface {
		Error() string
		Code() int
		Message() string
	}
	Error struct {
		Err      string `json:"err"`
		Msg      string `json:"msg"`
		HttpCode int    `json:"http_code"`
	}
)

func (err Error) Error() string {
	return err.Err
}

func (err Error) Code() int {
	return err.HttpCode
}

func (err Error) Message() string {
	return err.Msg
}
