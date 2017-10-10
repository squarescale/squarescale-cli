#!/usr/bin/env python3
import sys, subprocess
import boto3

git_branch = None

try:
    git_branch = subprocess.check_output(
        "git rev-parse --abbrev-ref HEAD",
        shell=True, universal_newlines=True).replace("\n", "")
except Exception:
    print("Not on a branch. Exiting.")
    sys.exit()

if git_branch != "master":
    print("Not on master branch. Exiting.")
    sys.exit()

bucket_name = "cli-releases"
executable_names = ["sqsc-linux-amd64", "sqsc-darwin-amd64"]

s3 = boto3.resource('s3')

for executable_name in executable_names:
    s3_key = "%s-latest" % executable_name
    print("Uploading %s to %s in bucket %s" % (executable_name, s3_key, bucket_name))
    s3.meta.client.upload_file(executable_name, bucket_name, s3_key)
    s3.meta.client.put_object_acl(ACL='public-read', Bucket=bucket_name, Key=s3_key)
