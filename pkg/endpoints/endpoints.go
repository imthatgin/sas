package endpoints

type HttpEndpointType uint

// AllEndpoints represents GET, PUT, POST, PATCH, DELETE.
const AllEndpoints = GET | PUT | POST | DELETE

const (
	GET HttpEndpointType = 1 << iota
	PUT
	POST
	PATCH // TODO: Distinguish between patch and put
	DELETE
)

func Has(flags HttpEndpointType, test HttpEndpointType) bool {
	return flags&test == test
}
