package storage

type Saver interface {
	SaveData(data *byte[]) bool
	SaveKeys(keys *map[string]interface{}) bool
}

type Loader interface {
	LoadData(data byte[]) bool
	LoadKeys(keys map[string]interface{}) bool
}

type LoadSaver interface {
	Loader
	Saver
}
