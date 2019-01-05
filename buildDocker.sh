#!/bin/bash
sudo docker build . -t idea_evolver:latest
sudo docker run -p 3000:3000 --name idea_container idea_evolver:latest
