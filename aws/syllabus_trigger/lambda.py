import os
import urllib.request
import jwt


def handler(event, _):
    domain = os.getenv("DOMAIN")
    jwt_secret = os.getenv("JWT_SECRET")

    if not domain or not jwt_secret:
        print("environment variable 'JWT_SECRET' or 'DOMAIN' are not set")
        return

    for record in event["Records"]:
        try:
            s3_info = record["s3"]
            syllabus_id = s3_info["object"]["key"]
            print("syllabus received:", syllabus_id)

            # Algo must match Go server decode algo
            token = jwt.encode({}, jwt_secret, algorithm="HS256")

            url = f"{domain}/syllabi/{syllabus_id}/sync"
            req = urllib.request.Request(
                url=url, method="GET", headers={"Authorization": f"Bearer {token}"}
            )

            with urllib.request.urlopen(req) as response:
                print(f"response status: {response.status}")

        except Exception as e:
            print(f"error processing record: {e}")
