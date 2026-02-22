//go:build functional

package functional

import (
	"fmt"
	"net/http"

	"starter-boilerplate/internal/user/transport/handler"
)

// --- Login tests ---

func (s *FunctionalSuite) TestLogin_Success() {
	body := `{"email":"admin@example.com","password":"P@ssw0rd123"}`
	resp := s.DoRequest(http.MethodPost, "/api/v1/auth/login", body, nil)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	var tok handler.TokenBody
	s.ReadJSON(resp, &tok)
	s.Assert().NotEmpty(tok.AccessToken)
	s.Assert().NotEmpty(tok.RefreshToken)
}

func (s *FunctionalSuite) TestLogin_WrongPassword() {
	body := `{"email":"admin@example.com","password":"WrongPassword123"}`
	resp := s.DoRequest(http.MethodPost, "/api/v1/auth/login", body, nil)
	defer resp.Body.Close()
	s.Assert().Equal(http.StatusUnauthorized, resp.StatusCode)
}

func (s *FunctionalSuite) TestLogin_NonexistentEmail() {
	body := `{"email":"nonexistent@example.com","password":"P@ssw0rd123"}`
	resp := s.DoRequest(http.MethodPost, "/api/v1/auth/login", body, nil)
	defer resp.Body.Close()
	s.Assert().Equal(http.StatusUnauthorized, resp.StatusCode)
}

func (s *FunctionalSuite) TestLogin_InvalidInput_MissingEmail() {
	body := `{"password":"P@ssw0rd123"}`
	resp := s.DoRequest(http.MethodPost, "/api/v1/auth/login", body, nil)
	defer resp.Body.Close()
	s.Assert().Equal(http.StatusUnprocessableEntity, resp.StatusCode)
}

func (s *FunctionalSuite) TestLogin_InvalidInput_ShortPassword() {
	body := `{"email":"admin@example.com","password":"12345"}`
	resp := s.DoRequest(http.MethodPost, "/api/v1/auth/login", body, nil)
	defer resp.Body.Close()
	s.Assert().Equal(http.StatusUnprocessableEntity, resp.StatusCode)
}

func (s *FunctionalSuite) TestLogin_InvalidInput_BadEmailFormat() {
	body := `{"email":"not-an-email","password":"P@ssw0rd123"}`
	resp := s.DoRequest(http.MethodPost, "/api/v1/auth/login", body, nil)
	defer resp.Body.Close()
	s.Assert().Equal(http.StatusUnprocessableEntity, resp.StatusCode)
}

// --- Refresh tests ---

func (s *FunctionalSuite) TestRefresh_Success() {
	rt := s.IssueRefreshToken("usr-admin-001", "admin")
	body := fmt.Sprintf(`{"refresh_token":"%s"}`, rt)
	resp := s.DoRequest(http.MethodPost, "/api/v1/auth/refresh", body, nil)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	var tok handler.TokenBody
	s.ReadJSON(resp, &tok)
	s.Assert().NotEmpty(tok.AccessToken)
	s.Assert().NotEmpty(tok.RefreshToken)
}

func (s *FunctionalSuite) TestRefresh_InvalidToken() {
	body := `{"refresh_token":"garbage.token.here"}`
	resp := s.DoRequest(http.MethodPost, "/api/v1/auth/refresh", body, nil)
	defer resp.Body.Close()
	s.Assert().Equal(http.StatusUnauthorized, resp.StatusCode)
}

func (s *FunctionalSuite) TestRefresh_AccessTokenUsedAsRefresh() {
	at := s.IssueAccessToken("usr-admin-001", "admin")
	body := fmt.Sprintf(`{"refresh_token":"%s"}`, at)
	resp := s.DoRequest(http.MethodPost, "/api/v1/auth/refresh", body, nil)
	defer resp.Body.Close()
	s.Assert().Equal(http.StatusUnauthorized, resp.StatusCode)
}

// --- Full flow tests ---

func (s *FunctionalSuite) TestLoginThenGetUser_FullFlow() {
	// Login
	loginBody := `{"email":"user@example.com","password":"P@ssw0rd123"}`
	loginResp := s.DoRequest(http.MethodPost, "/api/v1/auth/login", loginBody, nil)
	s.Require().Equal(http.StatusOK, loginResp.StatusCode)

	var tok handler.TokenBody
	s.ReadJSON(loginResp, &tok)

	// GetUser with access token
	userResp := s.DoAuthRequest(http.MethodGet, "/api/v1/users/usr-user-001", tok.AccessToken, "")
	s.Require().Equal(http.StatusOK, userResp.StatusCode)

	var u handler.GetUserBody
	s.ReadJSON(userResp, &u)
	s.Assert().Equal("usr-user-001", u.ID)
	s.Assert().Equal("user@example.com", u.Email)
}

func (s *FunctionalSuite) TestLoginThenRefreshThenGetUser_FullFlow() {
	// Login
	loginBody := `{"email":"admin@example.com","password":"P@ssw0rd123"}`
	loginResp := s.DoRequest(http.MethodPost, "/api/v1/auth/login", loginBody, nil)
	s.Require().Equal(http.StatusOK, loginResp.StatusCode)

	var tok handler.TokenBody
	s.ReadJSON(loginResp, &tok)

	// Refresh
	refreshBody := fmt.Sprintf(`{"refresh_token":"%s"}`, tok.RefreshToken)
	refreshResp := s.DoRequest(http.MethodPost, "/api/v1/auth/refresh", refreshBody, nil)
	s.Require().Equal(http.StatusOK, refreshResp.StatusCode)

	var newTok handler.TokenBody
	s.ReadJSON(refreshResp, &newTok)

	// GetUser with new access token
	userResp := s.DoAuthRequest(http.MethodGet, "/api/v1/users/usr-admin-001", newTok.AccessToken, "")
	s.Require().Equal(http.StatusOK, userResp.StatusCode)

	var u handler.GetUserBody
	s.ReadJSON(userResp, &u)
	s.Assert().Equal("usr-admin-001", u.ID)
	s.Assert().Equal("admin@example.com", u.Email)
	s.Assert().Equal("admin", u.Role)
}
