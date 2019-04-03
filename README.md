# BuildingLoader

This repository contains the code for loading the Building footprints dataset into a local MongoDB database using an ETL process written in Go.  
It then provides a locally hosted API for getting the data from the database in JSON format with options for filtering the data in many ways.  
It also provides summary statistics like mean and median for numerical data and maximum value counts for non-numerical data.

## Installation

This project has the following requirements:  

- Go Programming Language distribution (download and setup as per the instructions from here - <https://golang.org/dl/)>
- Local MongoDB server (download and install as per the instructions from here - <https://www.mongodb.com/download-center/community)>
- Gorilla Mux library and Mongo go driver for Go
  - After installing Go run the following command from the terminal `go get github.com/gorilla/mux go.mongodb.org/mongo-driver/mongo`  
- You can then download this project by running `go get github.com/nitishp1812/buildingloader`  
This downloads the project to the `$GOPATH/src/nitishp1812/buildingloader` directory

## Running

After downloading the repository, you can build the repository with the following command from the terminal from the
`$GOPATH` directory:  
`go build github.com/nitishp1812/buildingloader`  
This generates the executable in the same directory as the user called it from  
Alternatively, running `go install github.com/nitishp1812/buildingloader` generates the executable in the `$GOPATH/bin` directory.  

Running this executable builds the database by calling the API and then sets up a local server at `localhost:5000/`

## Results

`localhost:5000/` gives the JSON output for all the entries in the database.  
This output can be filtered by using `localhost:5000/filter/{parameter}{comparison}{parameter}`  
For example, this can be `localhost:5000/filter/construct_year>2009`  
We can compare different kinds of parameters together as well. For example,  
`localhost:5000/filter/base_bbl=mpluto_bbl`  
The full list of parameters which can be used to filter the output are given in the Details section  

Some summary statistics of the data can also be seen by going to `localhost:5000/summary/`.  
This shows the summary statistics of the whole data. This can also be filtered down to show the summary statistics for some filtered data like `localhost:5000/summary/filter/{parameter}{comparison}{parameter}`  

The comparison operators that are supported are:

- = (check for equality)
- <, > (check for greater than, less than)
- ! (check for inequality)

## Details

The parameters which can be used to filter the output are:

Field name | Description | Type
---------- | ----------- | ----
base_bbl | Borough, block, and lot number for the tax lot that the footprint is physically located within. | text
bin | Building Identification Number. A number assigned by City Planning and used by Dept. of Buildings to reference information pertaining to an individual building. The first digit is a borough code (1 = Manhattan, 2 = The Bronx, 3 = Brooklyn, 4 = Queens, 5 = Staten Island). The remaining 6 digits are unique for buildings within that borough. In some cases where these 6 digits are all zeros (e.g. 1000000, 2000000, etc.) the BIN is unassigned or unknown. | text
construct_year | The year construction of the building was completed. Records where this is zero or NULL mean that this information was not available. | text
doitt_id | Unique identifier assigned by DOITT(Department of Information Technology and Telecommunications) | text
feat_code | Type of Building. List of values:<br>2100 = Building<br>5100 = Building Under Construction<br>5110 = Garage<br>2110 = Skybridge<br>1001 = Gas Station Canopy<br>1002 = Storage Tank<br>1003 = Placeholder (triangle for permitted bldg)<br>1004 = Auxiliary Structure (non-addressable, not garage)<br>1005 = Temporary Structure (e.g. construction trailer) | text
geom_source | Indicates the reference source used to add or update the feature. Photogrammetric means the feature was added or updated using photogrammetric stereo-compilation methodology. This is the most accurate update method and should conform to the ASPRS accuracy standards. Other (Manual) means the feature was added or updated by heads-up digitizing from orthophotos or approximated from a plan drawing. These features will be generally be less accurate and may not conform to the ASPRS accuracy standards. | text
ground_elev | Lowest Elevation at the building ground level. Calculated from LiDAR or photogrammetrically. | number
height_roof | Building Height is calculated as the difference from the building elevation from the Elevation point feature class and the elevation in the interpolated TIN model. This is the height of the roof above the ground elevation, NOT its height above sea level. Records where this is zero or NULL mean that this information was not available. | number
last_status_type | Feature last status type (Demolition, Alteration, Geometry, Initialization, Correction, Marked for Construction, Marked For Demolition, Constructed) | text
mpluto_bbl | Borough, block, and lot number to be used for joining the building footprints data to DCP's MapPLUTO data, which aggregates data for condominium buildings using DOF's billing BBL. For non-condominium buildings the billing BBL is the same as the BASE_BBL. For condominium buildings the billing BBL may be the same for multiple buildings on different physical tax lots if they are part of the same billing unit for DOF purposes. | text
shape_area | The area of the shape | number
shape_len | The length of the shape | number

There is another field in the database but this cannot be used to filter the output

Field Name | Description
---------- | -----------
last_mod_date | Feature last modified date
