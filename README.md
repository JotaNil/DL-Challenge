# DL-Challenge

## Environment setup
For the API to run the following configurations must be made for the correct connection to the DataBase:
* Be sure to have a postgresSQL or a SQL DB up with the [IP2Proxy data](https://lite.ip2location.com/database/px7-ip-proxytype-country-region-city-isp-domain-usagetype-asn) already imported.
* Change the connection info `(host,port,user,dbName)` in `./cmd/services/sql.go` to match yours. Default data it's ready for a default postgresSQL installation.
* If needed, change `ipdataSchemaTableName` in `./cmd/api/ipdata/dao.go` to match your schema and table name. The default used is: `proxydata.ip2location`.
* Set the environment variable `DL_CHALLENGE_DBPASS` with the password of the user defined in the connection info.
> if there is any error with the configuration the error message should be enough to correct them. This error will be given on the API startup, it will be present in a panic.

## Running the project

To run the project you may use your preferred IDE, or in the case you want to run it with a terminal just go to `./main` and execute `go run *.go`.
The app will run in `localhost:8000`. 
> If for some reason the 8000 port its already in use, it can be changed in `./main/routing.go`.

## Endpoints

### Get Ip count by country name
This endpoint returns the count of all ips present in the database of the given country.

Url:
> /ipdata/count/ip/{country_name}

Params:
> country_name: must be an ISO 3166 complaint string, capitalization is expected.

Response body: 
```
{
   "country_name":"Argentina",
   "ip_count":79012
}
```

cURL:
> curl 127.0.0.1:8000/ipdata/count/ip/Argentina -H "Accept: application/json"

### Get top ISPs from Switzerland
This endpoint returns a list of the top 10 ISPs from Switzerland with his respective ip count.

Url:
> /ipdata/top10/Switzerland

Response body: 
```
[
   {
      "isp":"Ivan Bulavkin",
      "ip_count":147
   },
   ...
]
```

cURL:
> curl 127.0.0.1:8000/ipdata/top10/switzerland -H "Accept: application/json"

### Get data by IP
This endpoint returns all the data available in the database of the given ip.

Url:
> /ipdata/{ip}

Params:
> ip: must be a valid IPv4. 

Response body: 
```
{
   "ip_from":95781810,
   "ip_to":95781817,
   "proxy_type":"PUB",
   "country_code":"GB",
   "county_name":"United Kingdom of Great Britain and Northern Ireland",
   "region_name":"England",
   "city_name":"Saint Albans",
   "isp":"IPXO Limited",
   "ip_string":"5.181.131.180"
}
```

cURL:
> curl 127.0.0.1:8000/ipdata/5.181.131.180 -H "Accept: application/json"

## Error handling

All endpoints will return the appropriate status code for the request. 
In the case that the response has a different from 200 status code, a string must be unmarshalled to parse the error that is returned in plain txt as the http go pkg marshals it.

