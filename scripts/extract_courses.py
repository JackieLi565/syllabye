"""
extract_courses.py

This script scrapes Toronto Metropolitan University's academic calendar to extract
course related data. The scraped data is saved as a JSON file.

Usage:
    python extract_courses.py

Arguments:
    year (optional): The academic year range to scrape (e.g., "2024-2025").
        Default: "2024-2025"

    path (optional): The output path of the JSON payload.
        Default: "./out"

    verbose (optional): Script verbose mode.
        Default: False

Outputs:
    - JSON file saved in ~/[syllabye_project]/out/courses.json containing the scraped course data.
"""

from datetime import datetime
import json
import requests
from bs4 import BeautifulSoup
import argparse

parser = argparse.ArgumentParser(description="Syllabye Course Scraper Script")

parser.add_argument(
    "--year", type=str, default="2024-2025", help="Website year range iteration"
)
parser.add_argument(
    "--path",
    type=str,
    default="./out",
    help="Directory path of output destination",
)
parser.add_argument("-v", "--verbose", action="store_true", help="Enable verbose mode")

args = parser.parse_args()


def v_print(*stm):
    """
    Helper function to print log message if verbose argument is specified.
    """
    if args.verbose:
        print(*(datetime.now().strftime("%Y-%m-%d %H:%M:%S "), *stm))


domain = "https://www.torontomu.ca"
courses_uri = f"/calendar/{args.year}/courses/"
# Request courses page
res = requests.get(domain + courses_uri)
if res.status_code != 200:
    v_print(
        f"[request] request to {courses_uri} failed with error code {res.status_code}"
    )
    exit(f"Request Error: HTTP request to {courses_uri} failed")
else:
    v_print(f"[request] {courses_uri} request complete")

soup = BeautifulSoup(res.content, "html.parser")

# Extract courses rows from the table while removing the header
table_rows = soup.find_all("tr")[1::]

courses = []

# Iterate through each program and extract the course data source link
# Request the course payload and extract course specific data from the payload
data_endpoint = "/jcr:content/content/rescalendarcoursestack.data.1.json"
for row in table_rows:
    tds = row.find_all("td")
    a_tag = tds[0].find("a")

    course_category = str(a_tag.text).split("(")[0].strip()
    # Remove html page extension for json data endpoint
    course_uri = a_tag.get("href").split(".html")[0] + data_endpoint

    res = requests.get(domain + course_uri)
    if res.status_code != 200:
        v_print(
            f"[request] request to {course_uri} failed with error code {res.status_code}"
        )
        exit(f"Request Error: HTTP request to {course_uri} failed")
    else:
        v_print(f"[request] {course_uri} request complete")

    payload = res.json()
    # Iterate through each course and request the detail payload
    for course in payload["data"]:
        course_data_url = course["dataURL"]
        res = requests.get(domain + course_data_url)
        if res.status_code != 200:
            v_print(
                f"[request] request to {course_data_url} failed with error code {res.status_code}"
            )
            exit(f"Request Error: HTTP request to {course_data_url} failed")
        else:
            v_print(f"[request] {course_data_url} request complete")

        payload: dict = res.json()
        courses.append(
            {
                "category": course_category,
                "title": payload["longTitle"],
                "description": str(payload["courseDescription"])
                .encode("utf-8")
                .decode("unicode_escape"),
                "uri": "/" + "/".join(course["page"].split("/")[-3:]),
                "course": payload["courseCode"],
                "alpha": payload.get("courseAlphaCode", None),
                "code": payload.get("courseNumberCode", None),
            }
        )

    v_print(f"[application] category {course_category} requests complete")

output_file = f"{args.path}/courses.json"
with open(output_file, "w") as file:
    json.dump(courses, file)

v_print(f"[file] courses written to {output_file}")
