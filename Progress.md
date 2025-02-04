# Progress Tracker

TODO:
- Documentation README explaining the code
- \implement query 
    - Return all matched payloads (based on filter severity)
- Tests for Core Logic (60%+)

DONE:
- Server startup
- Modules
- PARSE INDIVIDUAL PAYLOADS AND ADD TO DB
- ADD all fields to DB
- Connect to DB
- Store EACH payload in the file one by one with metadata

CAN DO (Future)
- store more metadata like vulnerabilities etc

Random Questions / Doubts:
1. Can I assume that all the json files that will be scanned will be of a similar format? 
    I think YES for this example.

2. "Concurrency: Process ≥ 3 files in parallel."
    Does the problem statement mean concurrency or parallelism??
    I have assumed Concurrency, and used goroutines.

3. "Error Handling: Retry failed GitHub API calls (2 attempts)."
    Total number of attempts should be 3 or 2?
    I have implemented as 1 attempt and 2 reattempt.