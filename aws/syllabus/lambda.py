import os
import urllib.request
import jwt


def handler(event, _):
    domain = os.getenv("DOMAIN")
    jwt_secret = os.getenv("JWT_SECRET")

    try:
        record = event["Records"][0]
        key = record["s3"]["object"]["key"]

        # Algo must match Go server decode algo
        token = jwt.encode({}, jwt_secret, algorithm="HS256")

        url = f"{domain}/syllabi/{key}/sync"
        req = urllib.request.Request(
            url=url, method="GET", headers={"Authorization": f"Bearer {token}"}
        )

        with urllib.request.urlopen(req) as response:
            print(f"syllabus synced: {key}")
            print(f"response status: {response.status}")

    except Exception as e:
        print(f"error processing record: {e}")
