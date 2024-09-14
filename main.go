package main

import (
	"flag"
	"os/exec"

	"github.com/strifel/openid-connect-debugger/flag_parsing"
	getcode "github.com/strifel/openid-connect-debugger/get_code"
	gettoken "github.com/strifel/openid-connect-debugger/get_token"
	startconnection "github.com/strifel/openid-connect-debugger/start_connection"
	userdata "github.com/strifel/openid-connect-debugger/user_data"

	logging "github.com/strifel/openid-connect-debugger/logging"
)

func main() {
	endpointUrl := flag.String("endpoint", "", "URL to instance")
	clientId := flag.String("clientid", "", "Client ID")
	secret := flag.String("secret", "", "Client Secret")
	var scopes flag_parsing.ScopeFlags
	flag.Var(&scopes, "scope", "Specify scopes (e.g. -scope profile -scope openid)")
	verbosity := flag.Int("verbosity", 0, "Verbosity level (0 to 3). 2 is verbose, 3 shows the access token")
	macosOpen := flag.Bool("macosopen", false, "Open the browser on MacOS")
	callbackEndpoint := flag.String("callbackurl", "", "Callback endpoint. Default http://localhost")

	flag.Parse()

	if *endpointUrl == "" {
		logging.Error("No endpoint URL given")
		return
	}

	if *clientId == "" {
		logging.Error("No client ID given")
		return
	}

	if *secret == "" {
		logging.Error("No secret given")
		return
	}

	if *callbackEndpoint == "" {
		*callbackEndpoint = "localhost"
	}

	wellKnownURL := *endpointUrl + "/.well-known/openid-configuration"
	logging.Info("Starting connection to %s", wellKnownURL)
	information, err := startconnection.Fetch(wellKnownURL)
	if err != nil {
		logging.Error("%v", err)
		return
	}

	if *verbosity >= 2 {
		logging.Debug("Issuer:                 %s", information.Issuer)
		logging.Debug("Authorization endpoint: %s", information.AuthorizationEndpoint)
		logging.Debug("Token endpoint:         %s", information.TokenEndpoint)
		logging.Debug("Userinfo endpoint:      %s", information.UserinfoEndpoint)
	}

	logging.Info("Read data")

	logging.Info("Starting authorization phase")
	callbackURL := "http://" + *callbackEndpoint + ":8070/callback"
	logging.Info("Please add " + callbackURL + " as valid redirect URI to your client")
	authUrl := information.AuthorizationEndpoint + "?client_id=" + *clientId + "&response_type=code&scope=" + scopes.String() + "&redirect_uri=" + callbackURL
	logging.Info("Then visit %s", authUrl)

	if *macosOpen {
		exec.Command("open", authUrl).Run()
	}

	code := getcode.GetCode(*callbackEndpoint)

	if *verbosity >= 2 {
		logging.Debug("Query Client code %s", code)
	}

	logging.Success("Completed authorization phase")

	logging.Info("Starting token phase")
	token, err := gettoken.Fetch(information.TokenEndpoint, *clientId, *secret, code, *callbackEndpoint)
	if err != nil {
		logging.Error("%v", err)
		return
	}
	if *verbosity >= 3 {
		logging.Debug("Access token:\n%s", token.AccessToken)
		logging.Debug("Refresh token:\n%s", token.RefreshToken)
		logging.Debug("ID Token:\n%s", token.IDToken)
	}
	if *verbosity >= 2 {
		logging.Debug("Session state:	%s", token.SessionState)
		logging.Debug("Access token TTL:	%d", token.ExpiresIn)
		logging.Debug("Refresh token TTL:	%d", token.RefreshExpiresIn)
		logging.Debug("Not before policy:	%d", token.NotBeforePolicy)
		logging.Debug("Scope:		%s", token.Scope)
	}

	logging.Success("Completed token phase")

	logging.Info("Starting userinfo phase")
	userinfo, rawUser, err := userdata.Fetch(information.UserinfoEndpoint, token.AccessToken)
	if err != nil {
		logging.Error("%v", err)
		return
	}
	if *verbosity >= 1 {
		logging.Info("Username:		%s", userinfo.PreferredUsername)
		logging.Info("Email:		%s", userinfo.Email)
		logging.Info("Email verified:	%t", userinfo.EmailVerified)
		logging.Info("Name:		%s", userinfo.Name)
		logging.Info("Family name:		%s", userinfo.FamilyName)
		logging.Info("Given name:		%s", userinfo.GivenName)
		logging.Info("Sub:			%s", userinfo.Sub)
	}
	if *verbosity >= 2 {
		logging.Debug("Raw User Data: %s", rawUser)
	}

	logging.Success("Completed userinfo phase")
}
