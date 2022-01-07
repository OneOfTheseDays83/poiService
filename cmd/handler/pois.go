package handler

import (
	"errors"
	"github.com/google/uuid"
	"poi-service/cmd/data"
)

// PoiHandler provide abstraction to manage pois
type PoiHandler interface {
	Create(poi *data.Poi) (uniqueId string, err error)
	Update(idToUpdate data.Id, updatedPoi *data.Poi) (err error)
	Get(id data.Id) (resp data.Poi, err error)
	Delete(id data.Id) (err error)
	Search(pos data.SearchArea) (resp data.Pois, err error)
}

func NewPoiHandler(dbHandler DbHandler) PoiHandler {
	if dbHandler == nil {
		return nil
	}
	return &poiHandler{dbHandler: dbHandler}
}

type poiHandler struct {
	dbHandler DbHandler
}

func (p *poiHandler) Create(poi *data.Poi) (uniqueId string, err error) {
	if poi == nil {
		return "", errors.New("poi is nil")
	}
	if poi.Longitude < -180 || poi.Longitude > 80 {
		return "", errors.New("longitude out of range")
	}
	if poi.Latitude < -90 || poi.Latitude > 90 {
		return "", errors.New("latitude out of range")
	}

	uniqueId, err = p.dbHandler.AddPoi(PoiDbEntry{
		Id:       uuid.New().String(),
		Name:     poi.Name,
		Location: NewLocation(poi.Latitude, poi.Longitude),
	})
	return uniqueId, err
}

func (p *poiHandler) Update(idToUpdate data.Id, updatedPoi *data.Poi) error {
	return p.dbHandler.UpdatePoi(string(idToUpdate), PoiDbEntry{
		Id:       string(idToUpdate),
		Name:     updatedPoi.Name,
		Location: NewLocation(updatedPoi.Latitude, updatedPoi.Longitude),
	})
}

func (p *poiHandler) Get(id data.Id) (resp data.Poi, err error) {
	result, err := p.dbHandler.GetPoi(string(id))
	if err != nil {
		return
	}

	resp.Latitude = result.Location.Coordinates[0]
	resp.Longitude = result.Location.Coordinates[1]
	resp.Name = result.Name
	return
}

func (p *poiHandler) Delete(id data.Id) error {
	return p.dbHandler.DeletePoi(string(id))
}

func (p *poiHandler) Search(pos data.SearchArea) (resp data.Pois, err error) {
	var pois PoiDbEntries

	if pos.RadiusInMeter == 0 {
		// TODO a proper solution would use paging - but that is something to be adder later
		pois, err = p.dbHandler.GetAllPois()
	} else {
		pois, err = p.dbHandler.SearchByRadius(NewLocation(pos.Latitude, pos.Longitude), pos.RadiusInMeter)
	}

	if err != nil {
		return
	}

	for _, entry := range pois {
		resp = append(resp, data.Poi{
			Name:      entry.Name,
			Latitude:  entry.Location.Coordinates[0],
			Longitude: entry.Location.Coordinates[1],
		})
	}

	return resp, nil
}
