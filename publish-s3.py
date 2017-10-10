#!/usr/bin/env python3
import boto3

bucket_name = "cli-releases"
executable_names = ["sqsc-linux-amd64", "sqsc-darwin-amd64"]

s3 = boto3.resource('s3')

for executable_name in executable_names:
    s3_key = "%s-latest" % executable_name
    print("Uploading %s to %s in bucket %s" % (executable_name, s3_key, bucket_name))
    s3.meta.client.upload_file(executable_name, bucket_name, s3_key)
    s3.meta.client.put_object_acl(ACL='public-read', Bucket=bucket_name, Key=s3_key)
