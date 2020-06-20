<a href="http://www.gnu.org/licenses/gpl-3.0.html">
    <img src="https://img.shields.io/badge/License-GPL%20v3-blue.svg" alt="License" />
</a>
<a href="https://gitlab.com/bobymcbobs/go-rpi-gpio-api/releases">
    <img src="https://img.shields.io/badge/version-0.0.1-brightgreen.svg" />
</a>

# go-rpi-gpio-api

> A simple cloud-native golang api webserver for controlling a Raspberry Pi's GPIO pins

## Features
- Simple bearer token authorization support
- TLS/SSL support
- \>5MB container image

## API endpoints
| Endpoint                    | Purpose                                       |
| -                           | -                                             |
| `/api/pin`                  | List all pins (1-40) and their current states |
| `/api/pin/{number}`         | Get a single pin and return it's status       |
| `/api/pin/{number}/{state}` | Modify a pin's state (0 for low, 1 for high)  |

### Example request responses

#### All pins
```json
{
  "metadata": {
    "selfLink": "/api/pin",
    "version": "0.0.1",
    "timestamp": 1581673372,
    "response": "Fetched all pins"
  },
  "spec": [
    {
      "number": 1,
      "state": 1
    },
    {
      "number": 2,
      "state": 1
    },
```
...
```
    {
      "number": 39,
      "state": 1
    },
    {
      "number": 40,
      "state": 0
    }
  ]
}
```

### A single pin
```json
{
  "metadata": {
    "selfLink": "/api/pin/39",
    "version": "0.0.1",
    "timestamp": 1581673438,
    "response": "Fetched pin state"
  },
  "spec": {
    "number": 39,
    "state": 1
  }
}
```

## Local usage
```bash
docker run -it --rm --privileged -p 8080:8080 registry.gitlab.com/bobymcbobs/go-rpi-gpio-api:latest
```

## Building
```bash
docker build -t registry.gitlab.com/bobymcbobs/go-rpi-gpio-api:latest .
```

Note: since golang supports cross compilation, this container (Linux+arm) can be built from any supported platform!

## Deployment in k8s
Make sure you update the values in the yaml files
```bash
kubectl apply -f k8s-manifests/
```

## Environment variables

| Name                   | Purpose                                          | Defaults     |
| -                      | -                                                | -            |
| `APP_AUTH_SECRET`      | require a value in Authorization bearer header   |              |
| `APP_PORT`             | the port and interface which the app serves from | `:8080`      |
| `APP_PORT_TLS`         | the port and interface which the app serves from | `:4433`      |
| `APP_USE_TLS`          | run the app with TLS enabled                     | `false`      |
| `APP_TLS_PUBLIC_CERT`  | the public certificate for the app to use        | `server.crt` |
| `APP_TLS_PRIVATE_CERT` | the private cert for the app to use              | `server.tls` |

## Notes
Communication to GPIO pins requires either root privileges or preferably the user to be in the `gpio` group.

To make sure that GPIO access is configured correctly, run:
```sh
addgroup --gid 997 gpio
chown root.gpio /dev/gpiomem
chmod g+rw /dev/gpiomem
```

## License
Copyright 2019 Caleb Woodbine.
This project is licensed under the [GPL-3.0](http://www.gnu.org/licenses/gpl-3.0.html) and is [Free Software](https://www.gnu.org/philosophy/free-sw.en.html).
This program comes with absolutely no warranty.
