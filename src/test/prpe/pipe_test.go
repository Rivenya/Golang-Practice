package pipe

import (
	"strconv"
	"strings"
	"testing"
)

type Request interface{}

type Response interface{}

type error interface{}

type Filter interface {
	Process(data Request) (Response, error)
}

type SplitFilter struct {
	character string
}

type AddFilter struct {
}

func (p *SplitFilter) Process(data Request) (Response, error) {
	ret, ok := data.(string)
	if !ok {
		return nil, "Excepted Type,you should input string type"
	}
	return strings.Split(ret, p.character), nil
}

func (p *AddFilter) Process(data Request) (Response, error) {
	ret, ok := data.([]string)
	if !ok {
		return nil, "Excepted Type,you should input Array string type"
	}
	tmp := 0
	for _, val := range ret {
		ret, err := strconv.Atoi(val)
		if err != nil {
			return nil, "Excepted Type,you should input Array string type"
		}
		tmp += int(ret)
	}
	return tmp, nil
}

func pipeFilter(val Request, filter ...Filter) (Response, error) {
	for _, i := range filter {
		ret, err := i.Process(val)
		if err != nil {
			return nil, err
		}
		val = ret
	}
	return val, nil
}

func Test_Pipe(t *testing.T) {
	tmp := []int{1, 2, 3}
	splitFn := SplitFilter{character: ","}
	addFn := AddFilter{}
	ret, err := pipeFilter(tmp, &splitFn, &addFn)
	if err != nil {
		t.Log(err)
	} else {
		t.Log(ret)
	}
}
