# connet

## to test manually
0. push a couple test apps if you don't already have them:

  ```
  (cd src/acceptance/example-apps/proxy && cf push test1 && cf push test2)
  ```

0. launch the policy server

  ```
  (cd src/policy-server && go run main.go -configFile <( echo '{ "listen_address": "127.0.0.1:5555" }' ))
  ```

0. then in a separate terminal try out the cf cli plugin

  ```
  go install cf-cli-plugin && CF_TRACE=true cf uninstall-plugin connet; cf install-plugin -f bin/cf-cli-plugin && cf plugins

  cf net-allow test1 test2
  cf net-list
  cf net-disallow test1 test2
  cf net-list
  ```
