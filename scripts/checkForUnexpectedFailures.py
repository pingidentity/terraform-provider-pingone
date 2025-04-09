import sys
import os
import re
import requests
from collections import deque

# Check if the script is being run with the correct number of arguments
if len(sys.argv) != 2:
    print("Usage: python checkForUnexpectedFailures.py <test_output_file>")
    sys.exit(1)

# Check if required environment variables are set
if "GITHUB_REPOSITORY" not in os.environ:
    print("GITHUB_REPOSITORY environment variable is not set.")
    sys.exit(1)
if "GITHUB_TOKEN" not in os.environ:
    print("GITHUB_TOKEN environment variable is not set.")
    sys.exit(1)

testOutputFile = sys.argv[1]

# Verify the output file exists
if not os.path.exists(testOutputFile):
    print(f"Test output file {testOutputFile} does not exist.")
    sys.exit(1)

# Find all the tests that failed and their error messages
failedRegex = re.compile(r'^--- FAIL: (\S+) \(\S+s\)$')
errorsByTestName = {}
with open(rf'{testOutputFile}') as testOutput:
    reversedOutput = list(reversed(testOutput.readlines()))
    for i, line in enumerate(reversedOutput):
        match = failedRegex.match(line)
        # Find failed tests
        if match:
            testName = match.group(1)
            # Get the corresponding output for the test, continuing to read in the file in reverse until hitting the same test name
            testErrorLines = deque()
            for j in range(i+1, len(reversedOutput)):
                line = reversedOutput[j]
                # If we see the test name again, we have reached the end of the test output
                if testName in line:
                    break
                # Reset the error lines when we hit a go error progress line - all the test output will be in one continuous block
                if line.startswith("---") or line.startswith("==="):
                    testErrorLines = []
                else:
                    testErrorLines.insert(0, line)
            if len(testErrorLines) > 0:
                errorsByTestName[testName] = "".join(testErrorLines)

# Get issues in the repo that have the label "type/test-failure"
githubRepo = os.environ.get("GITHUB_REPOSITORY")
githubToken = os.environ.get("GITHUB_TOKEN")
ghIssuesUrl = f"https://api.github.com/repos/{githubRepo}/issues"
requestParams = {
    "labels": "type/test-failure"
}
response = requests.get(ghIssuesUrl, params=requestParams)
if response.status_code != 200:
    print(f"Failed to fetch issues from GitHub: {response.status_code}")
    print(response.json())
    sys.exit(1)

# Check for the individual test in the corresponding github issues
issues = response.json()
testNotFound = {name for name in errorsByTestName}
for issue in issues:
    for testName in errorsByTestName:
        if testName in issue['body'] or testName in issue['title']:
            #TODO: look at expected output in issue and compare, in case the failure is different.
            # Will need to handle randomized output like resource ids.
            # For now, just the presence of the test name in the issue is enough.
            # Test was found in issue
            testNotFound.remove(testName)

# If any of the tests are not in the issues, notify
if len(testNotFound) > 0:
    print("The following tests failed but are not known in any existing Github test-failure issues:")
    for testName in testNotFound:
        print(f"FAIL: {testName}")
        print(errorsByTestName[testName])
    #TODO send notification
    # Exit with a non-zero status to indicate failure
    sys.exit(1)
