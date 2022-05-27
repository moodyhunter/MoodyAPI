#!/bin/bash

python -m venv .venv
source ./.venv/bin/activate
pip install -r ./requirements.txt
pip install pip-upgrader

pip-upgrade
deactivate
