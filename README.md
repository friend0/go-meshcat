# meshcat-go
A golang interface to meshcat, a small application for running `three.js` applications.
The motivation for this project is to demonstrate networked simulations for robotics, using the NATS.io
message middleware as the communication layer. NATS has clients in just about every major language, and also interfaces with MQTT.
It is performant, and has top tier tooling to support development. 


This work is heavily inspired by `meshcat` and `meshcat-python`, and depends on the bootstrapped three.js
application contained therein. It implements a Go backend for serving the application bundle in production, 
as well as the NATS server. The server maintains routines


## Getting Started
Watch this section for updates on how to run the app. 


## Development

This project 