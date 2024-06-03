# meshcat-go
A golang interface to meshcat, a small application for running `three.js` applications.
The motivation for this project is to demonstrate networked simulations for robotics, using the NATS.io
message middleware as the communication layer. NATS has clients in just about every major language, and also interfaces with MQTT.
It is performant, and has top tier tooling to support development. 
`go-meshcat` implements a Go backend for serving the application bundle in production, 
as well as the NATS server. The basic workflow is that the backend `Echo` serves the client JS bundle on the main route `/`.
Echo also runs a Websocket server that the client connects to when it's loaded. In addition, the backend hosts subscriptions to `meshcat` topics
that map to the commands implemented in the `meshcat` frontend. When the nats service for a particular topic gets a new message, it uses the 
web socket connection to proxy the request to the client. Note that it may be possible in a future implementation for most of the nats processing 
to happen on the client. 

There are some rudimentary development flows in place with hot reloading on air (Go) and the webpack (frontend) sides of the house.
I will add deployment tooling eventually. There are some basic stack commands provided in the Makefile.

This work is heavily inspired by `meshcat` and `meshcat-python`, and depends on the bootstrapped three.js
application contained therein.

## Getting Started
Watch this section for updates on how to run the app. 
- Basic dev requirements: docker, go
- Install Air, Webpack


## Development
Use `nats-server` to launch a local NATS server if you're not running this against the production NATS server on `sunset`.
Run the backend with `air` at the root of this repo.
Run the frontend development server with `npx webpack serve`, in the `web/meshcat` directory.
The server should now be running at `http://localhost:8081`. Note that this is where the webpack dev server is hosting the frontend bundle. The backend is actually
running on `http://localhost:8080`. The former is the live bundle, and proxies to the backend via settings configured in webpack. The latter is the backend, which will serve the bundle that was built at runtime.

 TODO: 
 - reorganize around a `frontend`/`backend` structure.
 - refactor frontend to be more modular. It's currently one long main file.
 - implement htmx view window with some input features

