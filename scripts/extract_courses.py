"""
extract_courses.py

This script scrapes Toronto Metropolitan University's academic calendar to extract
course related data. The scraped data is saved as a JSON file.

Usage:
    python extract_courses.py [year_range]

Arguments:
    year_range (optional): The academic year range to scrape (e.g., "2024-2025").
                           Defaults to "2024-2025" if not provided.

Outputs:
    - JSON file saved in ./out/[year_range]_courses.json containing the scraped program data.

Example:
    python extract_courses.py 2023-2024
"""

import sys
import json
import requests
from bs4 import BeautifulSoup

arguments = sys.argv

year_range = arguments[1] if len(arguments) == 2 else "2024-2025"

domain = "https://www.torontomu.ca"
courses_uri = f"/calendar/{year_range}/courses/"

# Request courses page
res = requests.get(domain + courses_uri)
if res.status_code != 200:
    exit(f"Request error code: {res.status_code}")

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
        exit(f"Request error code: {res.status_code}\n\nurl: {res.url}")

    payload = res.json()
    # Iterate through each course and request the detail payload
    for course in payload["data"]:
        res = requests.get(domain + course["dataURL"])
        if res.status_code != 200:
            exit(f"Request error code: {res.status_code}\n\nurl: {res.url}")

        payload: dict = res.json()
        courses.append(
            {
                "category": course_category,
                "title": payload["longTitle"],
                "description": str(payload["courseDescription"])
                .encode("utf-8")
                .decode("unicode_escape"),
                "course": payload["courseCode"],
                "alpha": payload.get("courseAlphaCode", None),
                "code": payload.get("courseNumberCode", None),
            }
        )

    print(f"{course_category} - completed")


output_file = f"./out/{year_range}_courses.json"
with open(output_file, "w") as file:
    json.dump(courses, file)

print(f"Courses successfully saved to {output_file}")
