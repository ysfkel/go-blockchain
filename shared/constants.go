package shared

const (
	LatestHashKey = "latesthash"
	Genesis       = "Genesis"
	DbPath        = "./tmp/blocks"
	DbFile        = "./tmp/blocks/MANIFEST" //constant to check if database exists , badger creates this when db is initialized
)
