package handler

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"poi-service/cmd/data"
	"testing"
)

func Test_poiHandler_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mongoMock := NewMockDbHandler(ctrl)
	handlerToTest := NewPoiHandler(mongoMock)

	t.Run("handler not nil", func(t *testing.T) {
		assert.NotNil(t, handlerToTest)
	})

	t.Run("poi nil", func(t *testing.T) {
		id, err := handlerToTest.Create(nil)
		assert.NotNil(t, err)
		assert.Empty(t, id)
	})

	t.Run("longitude out of range", func(t *testing.T) {
		id, err := handlerToTest.Create(&data.Poi{Longitude: 82})
		assert.NotNil(t, err)
		assert.Empty(t, id)

		id, err = handlerToTest.Create(&data.Poi{Longitude: -200})
		assert.NotNil(t, err)
		assert.Empty(t, id)
	})

	t.Run("latitude out of range", func(t *testing.T) {
		id, err := handlerToTest.Create(&data.Poi{Latitude: 100})
		assert.NotNil(t, err)
		assert.Empty(t, id)

		id, err = handlerToTest.Create(&data.Poi{Latitude: -100})
		assert.NotNil(t, err)
		assert.Empty(t, id)
	})

	t.Run("create entry in db", func(t *testing.T) {
		mongoMock.EXPECT().AddPoi(gomock.Any()).Return("abc", nil)
		data := &data.Poi{
			Name:      "abc",
			Latitude:  90,
			Longitude: 20,
		}
		id, err := handlerToTest.Create(data)
		assert.Nil(t, err)
		assert.Equal(t, "abc", id)

		mongoMock.EXPECT().AddPoi(gomock.Any()).Return("", errors.New("Some error"))
		id, err = handlerToTest.Create(data)
		assert.NotNil(t, err)
	})

	t.Run("create entry in db", func(t *testing.T) {
		mongoMock.EXPECT().AddPoi(gomock.Any()).Return("abc", nil)
		data := &data.Poi{
			Name:      "abc",
			Latitude:  90,
			Longitude: 20,
		}
		id, err := handlerToTest.Create(data)
		assert.Nil(t, err)
		assert.Equal(t, "abc", id)

		mongoMock.EXPECT().AddPoi(gomock.Any()).Return("", errors.New("Some error"))
		id, err = handlerToTest.Create(data)
		assert.NotNil(t, err)
	})
}

func Test_poiHandler_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mongoMock := NewMockDbHandler(ctrl)
	handlerToTest := NewPoiHandler(mongoMock)

	t.Run("handler not nil", func(t *testing.T) {
		mongoMock.EXPECT().UpdatePoi("abc", gomock.Any()).Return(errors.New("Some error"))
		err := handlerToTest.Update(data.Id("abc"), &data.Poi{})
		assert.NotNil(t, err)

		mongoMock.EXPECT().UpdatePoi("abc", gomock.Any()).Return(nil)
		err = handlerToTest.Update(data.Id("abc"), &data.Poi{})
		assert.Nil(t, err)
	})
}

func Test_poiHandler_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mongoMock := NewMockDbHandler(ctrl)
	handlerToTest := NewPoiHandler(mongoMock)

	t.Run("handler not nil", func(t *testing.T) {
		mongoMock.EXPECT().GetPoi("abc").Return(PoiDbEntry{
			Id:       "abc",
			Name:     "mc donalds",
			Location: NewLocation(23, 25),
		}, nil)
		data, err := handlerToTest.Get(data.Id("abc"))
		assert.Nil(t, err)
		assert.Equal(t, "mc donalds", data.Name)
		assert.Equal(t, float64(23), data.Longitude)
		assert.Equal(t, float64(25), data.Latitude)
	})
}

func Test_poiHandler_Search(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mongoMock := NewMockDbHandler(ctrl)
	handlerToTest := NewPoiHandler(mongoMock)

	t.Run("get by radius", func(t *testing.T) {
		resp := PoiDbEntries{}
		resp = append(resp, PoiDbEntry{
			Id:       "abc",
			Name:     "mc donalds",
			Location: NewLocation(23, 25),
		})
		resp = append(resp, PoiDbEntry{
			Id:       "abcd",
			Name:     "burger king",
			Location: NewLocation(23, 25),
		})

		mongoMock.EXPECT().SearchByRadius(gomock.Any(), uint64(20)).Return(resp, nil)
		data, err := handlerToTest.Search(data.SearchArea{RadiusInMeter: 20})
		assert.Nil(t, err)
		assert.Equal(t, 2, len(data))
	})

	t.Run("get all", func(t *testing.T) {
		resp := PoiDbEntries{}
		resp = append(resp, PoiDbEntry{
			Id:       "abc",
			Name:     "mc donalds",
			Location: NewLocation(23, 25),
		})
		resp = append(resp, PoiDbEntry{
			Id:       "abcd",
			Name:     "burger king",
			Location: NewLocation(23, 25),
		})

		mongoMock.EXPECT().GetAllPois().Return(resp, nil)
		data, err := handlerToTest.Search(data.SearchArea{})
		assert.Nil(t, err)
		assert.Equal(t, 2, len(data))

	})
}
