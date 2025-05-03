# 🧠 Anomaly Detection System – Take-Home Assignment

## 📋 Overview

Your task is to build a system that ingests simulated time-series data, detects anomalies using configurable rules, and provides a frontend to visualize and resolve these anomalies.
This mirrors real challenges we face determining when there is an issue with a data collection pipeline, and how to resolve it.

You are encouraged to make and declare assumptions, however feel free to ask as many questions as necessary to understand the requirements and design your solution.

You are welcome to use any technologies.

---

## ✅ Core Requirements

### 1. Simulated Data Ingestion
- Approach ingestion however you want, we don't care how the data is delivered to the system - just that we can test functionality.
- We've provided 4 sample JSONL files (as well as their schema) to get you started but feel free to create sample data of your own to show functionality. Data will be "delivered" by running a script and passing a path to a local jsonl file.

### 2. Anomaly Detection
- Propose a set of alerting conditions/rules that can be used to detect anomalies. You should implement at least three of the rules you propose, including one for which the base metric is the count of null values by column:
- Define a schema for anomalies and write them to a database.

### 3. Configurable Rules 
Create a simple UI that:
- Visualizes incoming data and active anomalies
- Displays and manages rule configurations per source
- CRUD for detection rules (thresholds, windows)

### 4. API
- this can be stubbed, however the API should be callable with new datasets in whatever format you choose

---

## 📦 Submission

1. A zipped folder containing all code & data for the project, or invitation to a repo containing the code
2. A README.md file containing:
   - a brief description of the project
   - Quickstart/notes on running/deployment/test
   - how you might test this
   - future features, challenges that might surface at scale, limitations, things you chose not to do, tradeoffs etc
   - any other relevant information
3. Shell scripts/executables that can be used to run/test your system*
4. (Optional) a video/Loom or demo link

*The executables for testing the system should support passing a path to a local jsonl file, which can be ingested into the system and used for testing.
You may choose to require any additional arguments for things like labeling the source of the data, or specifying a different config file.

## ⏱ Time

Estimated completion time: **10–14 focused hours**

If you're unable to complete the project in what you feel is a reasonable time,
submit what you have and write next steps/what you would have done in
a README.md. If you need more than 1 week for any reason let
us know and we can work something out!
