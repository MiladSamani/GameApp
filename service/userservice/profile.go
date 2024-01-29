package userservice

import "gameAppProject/pkg/richerror"
import "gameAppProject/param"

func (s Service) Profile(req param.ProfileRequest) (param.ProfileResponse, error) {
	const op = "userservice.Profile"

	user, err := s.repo.GetUserByID(req.UserID)
	if err != nil {
		return param.ProfileResponse{}, richerror.New(op).WithErr(err).
			WithMeta(map[string]interface{}{"req": req})
	}

	return param.ProfileResponse{Name: user.Name}, nil
}
