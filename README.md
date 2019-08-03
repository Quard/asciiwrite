# ASCII Write

Render text with [FIGlet](http://www.figlet.org/) fonts (ASCII Art font)

```  ___   _____  _____  _____  _____   _    _        _  _         
 / _ \ /  ___|/  __ \|_   _||_   _| | |  | |      (_)| |        
/ /_\ \\ `--. | /  \/  | |    | |   | |  | | _ __  _ | |_   ___ 
|  _  | `--. \| |      | |    | |   | |/\| || '__|| || __| / _ \
| | | |/\__/ /| \__/\ _| |_  _| |_  \  /\  /| |   | || |_ |  __/
\_| |_/\____/  \____/ \___/  \___/   \/  \/ |_|   |_| \__| \___|
```

*my second educational project on Go*

## How to run

At first you need to save Firebase credentials json to `config/credentials/asciiwrite-firebase.json`

And then build and run docker image

`docker build -f build/Docker.local -t asciiwrite .`

`docker run -e AUTHTOKEN=<token for uploading fonts> asciiwrite`

or you could use `build/Dockerfile` to prepare an image and deploy to Google Cloud Run

## API

`api/openapi.yaml` — swagger schema 

`third_party/redoc-static.html` — rendered swagger schema info html