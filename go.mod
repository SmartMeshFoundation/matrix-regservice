module github.com/SmartMeshFoundation/matrix-regservice

require (
	github.com/SmartMeshFoundation/Photon v0.9.1
	github.com/ant0ine/go-json-rest v3.3.2+incompatible
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/ethereum/go-ethereum v1.8.17
	github.com/go-stack/stack v1.8.0 // indirect
	github.com/jinzhu/gorm v1.9.1
	github.com/jinzhu/inflection v0.0.0-20180308033659-04140366298a // indirect
	github.com/mattn/go-colorable v0.0.9
	github.com/mattn/go-isatty v0.0.4 // indirect
	github.com/mattn/go-sqlite3 v1.10.0 // indirect
	golang.org/x/crypto v0.0.0-20181106171534-e4dc69e5b2fd
	gopkg.in/urfave/cli.v1 v1.20.0
	gopkg.in/yaml.v2 v2.2.1
)

replace (
	github.com/ethereum/go-ethereum v1.8.17 => github.com/nkbai/go-ethereum v1.9.1
	golang.org/x/crypto v0.0.0-20181106171534-e4dc69e5b2fd => github.com/golang/crypto v0.0.0-20181106171534-e4dc69e5b2fd
)
