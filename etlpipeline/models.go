package etlpipeline

type APIBuilding struct {
	GeomSource    string       "json:\"geomsource\""
	MplutoBbl     string       "json:\"mpluto_bbl\""
	BaseBbl       string       "json:\"base_bbl\""
	ShapeLen      string       "json:\"shape_len\""
	ShapeArea     string       "json:\"shape_area\""
	GroundElev    string       "json:\"groundelev\""
	FeatCode      string       "json:\"feat_code\""
	HeightRoof    string       "json:\"heightroof\""
	DoittID       string       "json:\"doitt_id\""
	Lststatype    string       "json:\"lststatype\""
	Lstmoddate    string       "json:\"lstmoddate\""
	Name          string       "json:\"name\""
	ConstructYear string       "json:\"cnstrct_yr\""
	Bin           string       "json:\"bin\""
	TheGeom       MultiPolygon "json:\"the_geom\""
}

type DBBuilding struct {
	GeomSource      string
	MplutoBbl       string
	BaseBbl         string
	ShapeLen        float64
	ShapeArea       float64
	GroundElev      int64
	FeatCode        int64
	HeightRoof      float64
	DoittID         string
	Lststatype      string
	Lstmodyear      string
	Lstmodmonth     string
	Lstmodday       string
	Name            string
	ConstructYear   string
	Bin             string
	GeomType        string
	GeomCoordinates [][]float64
}

type MultiPolygon struct {
	Type        string          "json:\"type\""
	Coordinates [][][][]float64 "json:\"coordinates\""
}
