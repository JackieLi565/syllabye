import json
import urllib.request
import urllib.parse
import os


def handler(event, _):
    env = os.getenv("ENV")
    domain = os.getenv("DOMAIN")

    for record in event["Records"]:
        try:
            body = json.loads(record["body"])
            request_id = body["requestId"]
            print("received message from request:", request_id)

            parsed_url = urllib.parse.urlparse(body["url"])
            if not all([parsed_url.scheme, parsed_url.netloc]):
                print("invalid url received:", body["url"])
                continue

            url = domain + parsed_url.path if env == "development" else body["url"]
            req = urllib.request.Request(
                url=url,
                method=body["method"],
            )

            payload = body.get("payload")
            if payload:
                req.data = json.dumps(payload).encode("utf-8")

            headers = body.get("headers")
            if headers:
                req.headers = headers

            with urllib.request.urlopen(req) as response:
                print(f"request response status: {response.status}")

        except Exception as e:
            print(f"error processing message: {e}")
