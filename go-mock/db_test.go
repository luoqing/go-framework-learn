package main
import(
	"testing"
    gomock "github.com/golang/mock/gomock"
    "errors"
)

// https://golangrepo.com/repo/golang-mock-go-testing-frameworks
func TestGetFromDB(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish() // 断言 DB.Get() 方法是否被调用

	m := NewMockDB(ctrl)
	m.EXPECT().Get(gomock.Eq("Tom")).Return(100, errors.New("not exist"))
	m.EXPECT().Get(gomock.Any()).Return(630, nil)
	
	/*打桩stubs, 有明确的参数和返回值， 也可以动态设置返回值也经常使用
	m.EXPECT().Get(gomock.Not("Sam")).Return(0, nil) 
	m.EXPECT().Get(gomock.Nil()).Return(0, errors.New("nil")) 
	m.EXPECT().Get(gomock.Not("Sam")).Return(0, nil)
	m.EXPECT().Get(gomock.Any()).Do(func(key string) {
		t.Log(key)
	})
	m.EXPECT().Get(gomock.Any()).DoAndReturn(func(key string) (int, error) {
		if key == "Sam" {
			return 630, nil
		}
		return 0, errors.New("not exist")
	})
	*/

	if v := GetFromDB(m, "Tom"); v != -1 {
		t.Fatal("expected -1, but got", v)
	}
}

func TestTimesGetFromDB(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := NewMockDB(ctrl)
	m.EXPECT().Get(gomock.Not("Sam")).Return(0, nil).Times(2)
	GetFromDB(m, "ABC")
	GetFromDB(m, "DEF")
}

func TestOrderGetFromDB(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish() // 断言 DB.Get() 方法是否被调用

	m := NewMockDB(ctrl)
	o1 := m.EXPECT().Get(gomock.Eq("Tom")).Return(0, errors.New("not exist"))
	o2 := m.EXPECT().Get(gomock.Eq("Sam")).Return(630, nil)
	gomock.InOrder(o1, o2)
	GetFromDB(m, "Tom")
	GetFromDB(m, "Sam")
}

