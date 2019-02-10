package entities

type UserWithLocation struct {
	User             User
	Location         Location
	DistanceInMeters uint16
}
