//go:build functional

package functional

import (
	"net/http"

	"starter-boilerplate/internal/user/transport/dto"
)

func (s *FunctionalSuite) TestGetUser_AdminAccessesAnyUser() {
	token := s.IssueAccessToken("usr-admin-001", "admin")
	resp := s.DoAuthRequest(http.MethodGet, "/api/v1/users/usr-user-001", token, "")
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	var body struct {
		User dto.UserDTO `json:"user"`
	}
	s.ReadJSON(resp, &body)
	s.Assert().Equal("usr-user-001", body.User.ID)
	s.Assert().Equal("user@example.com", body.User.Email)
	s.Assert().Equal("user", body.User.Role)
}

func (s *FunctionalSuite) TestGetUser_UserAccessesSelf() {
	token := s.IssueAccessToken("usr-user-001", "user")
	resp := s.DoAuthRequest(http.MethodGet, "/api/v1/users/usr-user-001", token, "")
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	var body struct {
		User dto.UserDTO `json:"user"`
	}
	s.ReadJSON(resp, &body)
	s.Assert().Equal("usr-user-001", body.User.ID)
}

func (s *FunctionalSuite) TestGetUser_UserDeniedAccessToOther() {
	token := s.IssueAccessToken("usr-user-001", "user")
	resp := s.DoAuthRequest(http.MethodGet, "/api/v1/users/usr-user-002", token, "")
	defer resp.Body.Close()
	s.Assert().Equal(http.StatusForbidden, resp.StatusCode)
}

func (s *FunctionalSuite) TestGetUser_NonexistentUser() {
	token := s.IssueAccessToken("usr-admin-001", "admin")
	resp := s.DoAuthRequest(http.MethodGet, "/api/v1/users/nonexistent-id", token, "")
	defer resp.Body.Close()
	s.Assert().Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *FunctionalSuite) TestGetUser_NoAuthHeader() {
	resp := s.DoRequest(http.MethodGet, "/api/v1/users/usr-user-001", "", nil)
	defer resp.Body.Close()
	s.Assert().Equal(http.StatusUnauthorized, resp.StatusCode)
}

func (s *FunctionalSuite) TestGetUser_InvalidToken() {
	resp := s.DoAuthRequest(http.MethodGet, "/api/v1/users/usr-user-001", "invalid.token.here", "")
	defer resp.Body.Close()
	s.Assert().Equal(http.StatusUnauthorized, resp.StatusCode)
}
