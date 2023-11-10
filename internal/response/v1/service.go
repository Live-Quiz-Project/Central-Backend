package v1

import (
	"time"

	l "github.com/Live-Quiz-Project/Backend/internal/live/v1"
)

type service struct {
	live l.Repository
	Repository
	timeout time.Duration
}

func NewService(lRepo l.Repository, rRepo Repository) Service {
	return &service{
		live:       lRepo,
		Repository: rRepo,
		timeout:    time.Duration(3) * time.Second,
	}
}

// ---------- Choice response related service methods ---------- //
// func (s *service) CreateChoiceResponse(ctx context.Context, req *CreateChoiceResponseRequest, uid uuid.UUID) (*CreateChoiceResponseResponse, error) {
// 	c, cancel := context.WithTimeout(ctx, s.timeout)
// 	defer cancel()

// 	r, err := s.Repository.CreateChoiceResponse(c, &ChoiceResponse{
// 		ParticipantID:  uid,
// 		OptionChoiceID: req.OptionChoiceID,
// 	})
// 	if err != nil {
// 		return &CreateChoiceResponseResponse{}, err
// 	}

// 	return &CreateChoiceResponseResponse{
// 		ID:             r.ID,
// 		ParticipantID:  r.ParticipantID,
// 		OptionChoiceID: r.OptionChoiceID,
// 	}, nil
// }
