package gee

import (
	"errors"
	"fmt"
	"testing"

	gomock "github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
)

// 使用testify
// https://juejin.cn/post/6917956015132672007
// 这个测试的逻辑如下：
// Sam在getFromPeer获取失败，只能从sinker中获取，所以sinker和pickpeer各命中一次
// 对于不是Sam的从getFromPeer获取成功，pickpeer解peergetter.Get各命中2次
func Test_GroupGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	NOT_EXISTED := errors.New("not exists")
	sinker := NewMockSinker(ctrl)
	sinker.EXPECT().Get(gomock.Any()).DoAndReturn(func(key string) ([]byte, error) {
		return []byte(key), nil
	}).Times(1)

	peers := NewMockPeerPicker(ctrl)

	name := "testgroup"
	g := NewGroup(name, 2<<10, sinker)
	getter := NewMockPeerGetter(ctrl)
	//peers.EXPECT().PickPeer(gomock.Any()).Return(getter, nil).Times(10)

	getter.EXPECT().Get(gomock.Eq(name), gomock.Any()).DoAndReturn(
		func(group, key string) ([]byte, error) {
			return []byte(key), nil
		}).Times(2)
	peers.EXPECT().PickPeer(gomock.Not("Sam")).Return(getter, nil).Times(2)     // 如果不设置times就只执行一次
	peers.EXPECT().PickPeer(gomock.Eq("Sam")).Return(nil, NOT_EXISTED).Times(1) // 如果不设置times就只执行一次
	g.RegisterPeers(peers)

	cases := []struct {
		//name string
		key   string
		value []byte
		err   error
	}{
		{"Bill", []byte("Bill"), nil},
		{"Sam", []byte("Sam"), nil}, // 此处应该会命中Sinker
		{"Tsing", []byte("Tsing"), nil},
	}
	// goconvey vs testify
	// https://knapsackpro.com/testing_frameworks/difference_between/goconvey/vs/go-testify
	// https://stackshare.io/stackups/goconvey-vs-testify
	/*
			Some of the features offered by GoConvey are:

		--Directly integrates with go test
		--Fully-automatic web UI (works with native Go tests, too)
		--Huge suite of regression tests
		On the other hand, Testify provides the following key features:
		--Easy assertions
		--Mocking
		--Testing suite interfaces and functions
	*/
	for _, tt := range cases {
		t.Run(tt.key, func(t *testing.T) {
			value, err := g.Get(tt.key)
			fmt.Println(string(value))
			fmt.Println(err)
			//assert.Equal(t, err, tt.err)
			//assert.Equal(t, value, tt.value)
			Convey("test Group.Get", t, func() {
				So(err, ShouldEqual, tt.err)
				So(string(value), ShouldEqual, string(tt.value)) //[]byte还无法should equal
			})

		})

	}

}

// 还可以测试GetFromPeer
func Test_getFromPeer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sinker := NewMockSinker(ctrl)
	peers := NewMockPeerPicker(ctrl)

	name := "testgroup"
	g := NewGroup(name, 2<<10, sinker)
	g.RegisterPeers(peers)
	getter := NewMockPeerGetter(ctrl)
	peers.EXPECT().PickPeer(gomock.Any()).Return(getter, nil).Times(3)
	getter.EXPECT().Get(gomock.Eq(name), gomock.Any()).DoAndReturn(
		func(group, key string) ([]byte, error) {
			return []byte(key), nil
		}).Times(3)

	cases := []struct {
		//name string
		key   string
		value []byte
		err   error
	}{
		{"Bill", []byte("Bill"), nil},
		{"Sam", []byte("Sam"), nil}, // 此处应该会命中Sinker
		{"Tsing", []byte("Tsing"), nil},
	}

	for _, tt := range cases {
		t.Run(tt.key, func(t *testing.T) {
			value, err := g.getFromPeer(tt.key)
			assert.Equal(t, err, tt.err)
			assert.Equal(t, value, tt.value)
		})
	}
}
