# OpenID Connect Debugger
A simple Debugger for OpenID Connect IdPs.

## Why
You ever needed to quickly check if the IdP works? What the user endpoint returns? Get the token?<br>
Then this is a simple command line utility for that use case.

There are some online tools for that but they either need a server or use cors.
This tool is a simple go binary that just runs on your computer.

## Usage
Run the binary. It will show the help.
It will work something like that:
```
openidcheck -clientid test -endpoint https://idp.example.com -secret 1234 -verbosity 3
```

It will automatically use the well-known definition of the server.

### Usage with Keycloak
For Keycloak you might need to specify the realm.
Use it something like:
```
openidcheck -clientid test -endpoint https://idp.example.com/realms/yourRealm -secret 1234 -verbosity 3
```