package crypto

import (
	"github.com/cvetkovski98/zvax-slots/internal/dto"
	"github.com/golang-jwt/jwt/v4"
)

const SECRET = "secret"

func SignReservation(reservation *dto.Reservation) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"slotId":        reservation.SlotID,
		"reservationId": reservation.ReservationID,
		"validUntil":    reservation.ValidUntil,
	})

	return token.SignedString([]byte(SECRET))
}
