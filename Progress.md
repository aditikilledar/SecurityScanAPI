# Progress Tracker

TODO:
All done! :)
(for now)

DONE:
- Server startup
- Documentation README explaining the code & how to run it
- Modules
- DockerFile test if it works
- PARSE INDIVIDUAL PAYLOADS AND ADD TO DB
- ADD all fields to DB
- Connect to DB
- Store EACH payload in the file one by one with metadata
- \implement query 
    - Return all matched payloads (based on filter severity)
- Check Edge cases
- Tests for Core Logic (60%+) + Document how to run it
-   how to see if it's 60% ? - done
-   test if the correct result is returned for a sample file? - done

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

4. I am assuming that each file doesn't contain duplicates of CVE
