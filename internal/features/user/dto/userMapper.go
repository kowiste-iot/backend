package dto

import "backend/internal/features/user/domain"


  func ToDTO(user *domain.User, roles []string)*UserDTO {

    return &UserDTO{
      ID: user.ID(),
      AuthID: user.AuthID(),
      FirstName: user.FirstName(),
      LastName: user.LastName(),
      Email: user.Email(),
      Roles: roles,
      UpdatedAt: user.UpdatedAt(),
    }
  }

