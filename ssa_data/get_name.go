package ssa_data

import (
	"errors"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

const entityTypeName string = "Name"
const entityTypeStats string = "Stats"
const filterNameEquals string = "Name = "
const filterGenderEquals string = "Gender ="
const filterRandomNumber string = "Random >="

type Gender string

const MaleFilter Gender = "M"
const FemaleFilter Gender = "F"

func GetName(ctx context.Context, name string, gender Gender) ([]*Name, error) {

	query := datastore.NewQuery(entityTypeName).
		Filter(filterNameEquals, name).
		Filter(filterGenderEquals, gender)

	return executeGetNameQuery(ctx, query)
}

func GetRandomName(ctx context.Context, gender Gender) (*Name, error) {
	query := datastore.NewQuery(entityTypeName).
		Filter(filterGenderEquals, gender).
		Filter(filterRandomNumber, randomFloat()).
		Limit(1)

	results, err := executeGetNameQuery(ctx, query)
	return results[0], err
}

func executeGetNameQuery(ctx context.Context, query *datastore.Query) ([]*Name, error) {
	results := []*Name{}
	keys, err := query.GetAll(ctx, &results)
	if err != nil {
		return nil, errors.New("Failed to get name: " + err.Error())
	}

	for index, key := range keys {
		statsName, err := getStatsForKey(ctx, results[index], key)
		results[index] = statsName
		if err != nil {
			return nil, err
		}
	}
	return results, nil
}

func getStatsForKey(ctx context.Context, name *Name, key *datastore.Key) (*Name, error) {

	statsQuery := datastore.NewQuery(entityTypeStats).Ancestor(key)
	stats := []*Stat{}
	_, err := statsQuery.GetAll(ctx, &stats)
	if err != nil {
		return nil, errors.New("Failed to retrieve stats: " + err.Error())
	}

	for _, stat := range stats {
		name.addStat(stat)
	}

	return name, nil
}
