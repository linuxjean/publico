#!/bin/bash

# Step 1: Update the package list and install required packages
sudo apt update
sudo apt install -y python3-venv python3-pip

# Step 2: Create a virtual environment
python3 -m venv venv

# Step 3: Activate the virtual environment
source venv/bin/activate

# Step 4: Install requirements using torsocks
if [ -f requirements.txt ]; then
    torsocks pip install -r requirements.txt
else
    echo "requirements.txt not found in the current directory. Exiting."
    deactivate
    exit 1
fi

# Step 5: Run the Python script
if [ -f bip84.py ]; then
    python3 bip84.py
else
    echo "bip84.py not found in the current directory. Exiting."
fi

# Deactivate the virtual environment
deactivate

