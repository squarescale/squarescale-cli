#!/usr/bin/python3

import os, requests, subprocess, sys, json

user_token = os.environ.get("GITHUB_USER_TOKEN")
if user_token is None:
    print("GITHUB_USER_TOKEN not defined.")
    sys.exit(1)

user  = user_token.split(":")[0]
token = user_token.split(":")[1]

# Get all releases
print("Retrieve all releases...")
url = "https://api.github.com/repos/squarescale/squarescale-cli/releases"
releases = requests.get(url, auth=(user, token))
if releases.status_code != 200:
    print("Error: GET releases " + str(releases.status_code))
    sys.exit(1)

git_branch = None
git_sha_1  = None

try:
    git_branch = subprocess.check_output(
        "git describe --exact-match",
        shell=True, universal_newlines=True).replace("\n", "")
except Exception:
    pass

# Find sha1
if git_branch:
    git_sha_1 = git_branch
else:
    git_sha_1 = subprocess.check_output(
        "git describe --always",
        shell=True, universal_newlines=True).replace("\n", "")

# Create new release
print("Create new release draft...")
headers = { "Content-Type": "application/json" }
data = {
    "tag_name": git_sha_1,
    "name": "cli latest release (" + git_sha_1 + ")",
    "draft": git_branch == None
}

created = requests.post(url, auth=(user, token), headers=headers, data=json.dumps(data))
if created.status_code != 201:
    print("Error: POST draft " + str(created.status_code))
    sys.exit(1)

release_id = str(created.json()["id"])

# Push executable contents
for executable in ["sqsc-linux-amd64", "sqsc-darwin-amd64"]:
    print("Push executable " + executable + " to release draft...")
    params = { "name": executable }
    headers = { "Content-Type": "application/octet-stream" }
    upload_url = "https://uploads.github.com/repos/squarescale/squarescale-cli/releases/" + release_id + "/assets"

    with open(executable, "rb") as f:
        data = f.read()

        uploaded = requests.post(
            upload_url, auth=(user, token),
            headers=headers, params=params, data=data)

        if uploaded.status_code != 201:
            print("Error: Upload " + executable + " " + str(uploaded.status_code))
            sys.exit(1)

# Clear all previous drafts
print("Remove old release drafts...")
for release in releases.json():
    if release["draft"]:
        delete_url = url + "/" + str(release["id"])
        deleted = requests.delete(delete_url, auth=(user, token))
        if deleted.status_code != 204:
            print("Error: DELETE draft " + str(deleted.status_code))
            sys.exit(1)

print("Done.")
