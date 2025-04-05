import json
import urllib.request


def handler(event, _):
    for record in event["Records"]:
        try:
            body = json.loads(record["body"])
            print("received message from request:", body["requestId"])

            req = urllib.request.Request(
                url=body["url"],
                method=body["method"],
                data=body["payload"],
                headers=body["headers"],
            )

            with urllib.request.urlopen(req) as response:
                print(f"request response status: {response.status}")

        except Exception as e:
            print(f"error processing message: {e}")
