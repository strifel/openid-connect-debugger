package main

import (
	"flag"
	"fmt"
	"os/exec"

	getcode "github.com/strifel/openid-connect-debugger/get_code"
	gettoken "github.com/strifel/openid-connect-debugger/get_token"
	"github.com/strifel/openid-connect-debugger/start_connection"
	userdata "github.com/strifel/openid-connect-debugger/user_data"
)

func main() {
	endpointUrl := flag.String("endpoint", "", "URL to instance")
	clientId := flag.String("clientid", "", "Client ID")
	secret := flag.String("secret", "", "Client Secret")
	verbosity := flag.Int("verbosity", 0, "Verbosity level (0 to 3). 2 is verbose, 3 shows the access token")
	macosOpen := flag.Bool("macosopen", false, "Open the browser on MacOS")

	flag.Parse()

	if *endpointUrl == "" {
		fmt.Println("No endpoint URL given")
		return
	}

	if *clientId == "" {
		fmt.Println("No client ID given")
		return
	}

	if *secret == "" {
		fmt.Println("No secret given")
		return
	}

	wellKnownURL := *endpointUrl + "/.well-known/openid-configuration"
	fmt.Println("Starting connection to", wellKnownURL)
	information, err := startconnection.Fetch(wellKnownURL)
	if err != nil {
		fmt.Println(err)
		return
	}

	if *verbosity >= 2 {
		fmt.Println("Issuer: ", information.Issuer)
		fmt.Println("Authorization endpoint: ", information.AuthorizationEndpoint)
		fmt.Println("Token endpoint: ", information.TokenEndpoint)
		fmt.Println("Userinfo endpoint: ", information.UserinfoEndpoint)
	}

	fmt.Println("Read data")

	fmt.Println("Starting authorization phase")
	fmt.Println("Please add http://localhost:8070/callback as valid redirect URI to your client")
	authUrl := information.AuthorizationEndpoint + "?client_id=" + *clientId + "&response_type=code&scope=openid&redirect_uri=http://localhost:8070/callback"
	fmt.Println("Then visit ", authUrl)

	if *macosOpen {
		exec.Command("open", authUrl).Run()
	}

	code := getcode.GetCode()

	if *verbosity >= 2 {
		fmt.Println("Got code: ", code)
	}

	fmt.Println("Completed authorization phase")

	fmt.Println("Starting token phase")
	token, err := gettoken.Fetch(information.TokenEndpoint, *clientId, *secret, code)
	if err != nil {
		fmt.Println(err)
		return
	}
	if *verbosity >= 3 {
		fmt.Println("Access token: ", token.AccessToken)
		fmt.Println("Refresh token: ", token.RefreshToken)
	}
	if *verbosity >= 2 {
		fmt.Println("Session state: ", token.SessionState)
		fmt.Println("Access token expires in: ", token.ExpiresIn)
		fmt.Println("Refresh token expires in: ", token.RefreshExpiresIn)
		fmt.Println("Not before policy: ", token.NotBeforePolicy)
		fmt.Println("Scope: ", token.Scope)
	}

	fmt.Println("Completed token phase")

	fmt.Println("Starting userinfo phase")
	userinfo, rawUser, err := userdata.Fetch(information.UserinfoEndpoint, token.AccessToken)
	if err != nil {
		fmt.Println(err)
		return
	}
	if *verbosity >= 1 {
		fmt.Println("Username: ", userinfo.PreferredUsername)
		fmt.Println("Email: ", userinfo.Email)
		fmt.Println("Email verified: ", userinfo.EmailVerified)
		fmt.Println("Name: ", userinfo.Name)
		fmt.Println("Family name: ", userinfo.FamilyName)
		fmt.Println("Given name: ", userinfo.GivenName)
		fmt.Println("Sub: ", userinfo.Sub)
	}
	if *verbosity >= 2 {
		fmt.Println("Raw User Data: ", rawUser)
	}

	fmt.Println("Completed userinfo phase")
}
