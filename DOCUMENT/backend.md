#Backend Config
- Run asset server first
  
  ```sh
  cd INDEX-WEB/STATIC/
  python ../DOCUMENT/cors_server.py
  ```

- To run Backend with various options, start with `CONFIG_PATH` variable.

  ```sh
  cd $GOPATH/src/github.com/stkim1/BACKEND/
  CONFIG_PATH="./config-dev.yaml" go run ./exec/backend/main.go
  ```
 
- If no options was provided, `config.yaml` is selected by default.